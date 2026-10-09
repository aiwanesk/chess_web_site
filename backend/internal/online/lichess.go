package online

import (
	"bufio"
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// lichessGame : les champs utiles d'une ligne du flux ndjson de
// /api/games/user/{pseudo}.
type lichessGame struct {
	ID         string `json:"id"`
	Rated      bool   `json:"rated"`
	Variant    string `json:"variant"`
	Speed      string `json:"speed"`
	CreatedAt  int64  `json:"createdAt"`
	Status     string `json:"status"`
	Winner     string `json:"winner"`
	Moves      string `json:"moves"`
	InitialFen string `json:"initialFen"`
	Players    struct {
		White lichessPlayer `json:"white"`
		Black lichessPlayer `json:"black"`
	} `json:"players"`
	Opening struct {
		ECO  string `json:"eco"`
		Name string `json:"name"`
	} `json:"opening"`
}

type lichessPlayer struct {
	User struct {
		Name string `json:"name"`
	} `json:"user"`
	Rating int `json:"rating"`
}

// fetchLichess lit le flux en continu : Lichess envoie une partie par ligne, à
// environ vingt par seconde pour une requête anonyme. Cinq cents parties
// prennent donc une demi-minute — c'est pour ça que le chargement tourne en
// arrière-plan plutôt que dans la requête du navigateur.
func (c *Client) fetchLichess(ctx context.Context, q Query, onGame func(Game)) error {
	v := url.Values{}
	v.Set("max", strconv.Itoa(q.Max))
	v.Set("moves", "true")
	v.Set("opening", "true")
	v.Set("clocks", "false")
	v.Set("evals", "false")
	if q.Rated {
		v.Set("rated", "true")
	}
	if !q.Since.IsZero() {
		v.Set("since", strconv.FormatInt(q.Since.UnixMilli(), 10))
	}
	// Lichess filtre lui-même par cadence : on ne télécharge pas ce qu'on jette.
	var perfs []string
	for _, sp := range Speeds {
		if q.wants(sp) {
			perfs = append(perfs, lichessPerfs[sp]...)
		}
	}
	if len(q.Speeds) > 0 {
		v.Set("perfType", strings.Join(perfs, ","))
	}
	resp, err := c.get(ctx, c.LichessURL+"/api/games/user/"+url.PathEscape(q.Username)+"?"+v.Encode(),
		"application/x-ndjson")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var lg lichessGame
		if err := json.Unmarshal(line, &lg); err != nil {
			continue
		}
		if g, ok := lichessToGame(lg); ok {
			onGame(g)
		}
	}
	return sc.Err()
}

// lichessPerfs : nos cadences vers les perfType de Lichess.
var lichessPerfs = map[string][]string{
	"bullet":         {"ultraBullet", "bullet"},
	"blitz":          {"blitz"},
	"rapid":          {"rapid"},
	"classical":      {"classical"},
	"correspondence": {"correspondence"},
}

func lichessToGame(lg lichessGame) (Game, bool) {
	if lg.Variant != "" && lg.Variant != "standard" {
		return Game{}, false
	}
	if lg.InitialFen != "" && lg.InitialFen != startFEN {
		return Game{}, false
	}
	switch lg.Status {
	case "created", "started", "aborted", "noStart":
		return Game{}, false // pas une vraie partie
	}
	san := strings.Fields(lg.Moves)
	if len(san) == 0 {
		return Game{}, false
	}
	uci, err := toUCI(san)
	if err != nil {
		return Game{}, false
	}
	res := 0
	switch lg.Winner {
	case "white":
		res = 1
	case "black":
		res = -1
	}
	speed := lg.Speed
	if speed == "ultraBullet" {
		speed = "bullet"
	}
	g := Game{
		Source:   Lichess,
		URL:      "https://lichess.org/" + lg.ID,
		White:    nameOr(lg.Players.White.User.Name, "Anonyme"),
		Black:    nameOr(lg.Players.Black.User.Name, "Anonyme"),
		WhiteElo: lg.Players.White.Rating,
		BlackElo: lg.Players.Black.Rating,
		Speed:    speed,
		Rated:    lg.Rated,
		Result:   res,
		ECO:      lg.Opening.ECO,
		Opening:  lg.Opening.Name,
		SAN:      san,
		UCI:      uci,
	}
	if lg.CreatedAt > 0 {
		g.Date = time.UnixMilli(lg.CreatedAt).UTC().Format("2006-01-02")
		g.Played = lg.CreatedAt / 1000
	}
	return g, true
}

func nameOr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
