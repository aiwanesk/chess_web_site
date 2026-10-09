package online

import (
	"context"
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/twic"
)

// Chess.com n'a pas de flux par joueur : il publie une archive par mois. On
// remonte donc les mois du plus récent au plus ancien jusqu'à avoir assez de
// parties, ou jusqu'à passer la date de début.

type chessComArchives struct {
	Archives []string `json:"archives"`
}

type chessComMonth struct {
	Games []chessComGame `json:"games"`
}

type chessComGame struct {
	URL          string         `json:"url"`
	PGN          string         `json:"pgn"`
	TimeClass    string         `json:"time_class"`
	Rated        bool           `json:"rated"`
	Rules        string         `json:"rules"`
	EndTime      int64          `json:"end_time"`
	InitialSetup string         `json:"initial_setup"`
	White        chessComPlayer `json:"white"`
	Black        chessComPlayer `json:"black"`
}

type chessComPlayer struct {
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Result   string `json:"result"`
}

func (c *Client) fetchChessCom(ctx context.Context, q Query, onGame func(Game)) error {
	base := c.ChessComURL + "/pub/player/" + url.PathEscape(strings.ToLower(q.Username))
	resp, err := c.get(ctx, base+"/games/archives", "application/json")
	if err != nil {
		return err
	}
	var arch chessComArchives
	err = json.NewDecoder(resp.Body).Decode(&arch)
	resp.Body.Close()
	if err != nil {
		return err
	}

	n := 0
	// Les mois arrivent du plus ancien au plus récent : on les prend à l'envers.
	for i := len(arch.Archives) - 1; i >= 0 && n < q.Max; i-- {
		monthURL := arch.Archives[i]
		// L'adresse vient de la réponse de Chess.com : on n'accepte que ses
		// propres archives, jamais un hôte arbitraire.
		if !strings.HasPrefix(monthURL, "https://api.chess.com/pub/player/") {
			continue
		}
		if !q.Since.IsZero() && monthEnd(monthURL).Before(q.Since) {
			break
		}
		resp, err := c.get(ctx, c.ChessComURL+strings.TrimPrefix(monthURL, "https://api.chess.com"), "application/json")
		if err != nil {
			return err
		}
		var m chessComMonth
		err = json.NewDecoder(resp.Body).Decode(&m)
		resp.Body.Close()
		if err != nil {
			return err
		}
		for j := len(m.Games) - 1; j >= 0 && n < q.Max; j-- {
			cg := m.Games[j]
			if !q.Since.IsZero() && time.Unix(cg.EndTime, 0).Before(q.Since) {
				continue
			}
			g, ok := chessComToGame(cg)
			if !ok || !q.wants(g.Speed) || (q.Rated && !g.Rated) {
				continue
			}
			onGame(g)
			n++
		}
	}
	return nil
}

// monthEnd : la fin du mois d'une archive « …/games/2026/09 ». Une adresse
// illisible renvoie l'instant présent : on ne la prend jamais pour trop vieille.
func monthEnd(u string) time.Time {
	parts := strings.Split(strings.TrimRight(u, "/"), "/")
	if len(parts) < 2 {
		return time.Now()
	}
	t, err := time.Parse("2006/01", parts[len(parts)-2]+"/"+parts[len(parts)-1])
	if err != nil {
		return time.Now()
	}
	return t.AddDate(0, 1, 0)
}

// Les codes de nulle de Chess.com. Tout autre résultat que « win » chez l'un
// veut dire que l'autre a gagné.
var chessComDraws = map[string]bool{
	"agreed": true, "repetition": true, "stalemate": true, "insufficient": true,
	"50move": true, "timevsinsufficient": true,
}

func chessComToGame(cg chessComGame) (Game, bool) {
	if cg.Rules != "" && cg.Rules != "chess" {
		return Game{}, false // Chess960, bughouse…
	}
	if cg.InitialSetup != "" && cg.InitialSetup != startFEN {
		return Game{}, false
	}
	san := twic.MainLine(cg.PGN)
	if len(san) == 0 {
		return Game{}, false
	}
	uci, err := toUCI(san)
	if err != nil {
		return Game{}, false
	}
	res := 0
	switch {
	case cg.White.Result == "win":
		res = 1
	case cg.Black.Result == "win":
		res = -1
	case chessComDraws[cg.White.Result]:
		res = 0
	}
	speed := cg.TimeClass
	if speed == "daily" {
		speed = "correspondence"
	}
	g := Game{
		Source:   ChessCom,
		URL:      cg.URL,
		White:    cg.White.Username,
		Black:    cg.Black.Username,
		WhiteElo: cg.White.Rating,
		BlackElo: cg.Black.Rating,
		Speed:    speed,
		Rated:    cg.Rated,
		Result:   res,
		ECO:      pgnTag(cg.PGN, "ECO"),
		SAN:      san,
		UCI:      uci,
	}
	if cg.EndTime > 0 {
		g.Date = time.Unix(cg.EndTime, 0).UTC().Format("2006-01-02")
	}
	g.Opening = openingName(pgnTag(cg.PGN, "ECOUrl"))
	return g, true
}

// pgnTag lit un en-tête PGN, ou "" s'il est absent.
func pgnTag(pgn, key string) string {
	prefix := "[" + key + " \""
	for _, line := range strings.Split(pgn, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSuffix(strings.TrimPrefix(line, prefix), "\"]")
		}
	}
	return ""
}

// moveInSlug repère le début d'une suite de coups dans une fiche d'ouverture :
// « …-Variation-5...Nf6-6.O-O » ou « …-Variation...9.Nc3 ».
var moveInSlug = regexp.MustCompile(`(\.\.\.|-)\d+\.`)

// openingName tire un nom lisible de l'URL de la fiche d'ouverture Chess.com,
// « …/openings/Sicilian-Defense-Kan-Variation...5.Nc3-Qc7 » → « Sicilian
// Defense Kan Variation ». La suite de coups collée au nom est coupée : elle
// noie l'information dans la liste des parties.
func openingName(ecoURL string) string {
	if ecoURL == "" {
		return ""
	}
	name := ecoURL[strings.LastIndex(ecoURL, "/")+1:]
	if loc := moveInSlug.FindStringIndex(name); loc != nil {
		name = name[:loc[0]]
	}
	name = strings.TrimRight(name, ".-")
	return strings.ReplaceAll(name, "-", " ")
}
