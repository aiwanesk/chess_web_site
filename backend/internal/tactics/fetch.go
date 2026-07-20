package tactics

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Game is one of Alexandre's games, pulled for analysis. None of this is ever
// published — only the mirrored puzzle derived from it is.
type Game struct {
	PGN       string
	MyColor   string // "white" | "black"
	TimeClass string // bullet | blitz | rapid | classical
	Source    string // "chesscom" | "lichess"
	EndTime   time.Time
}

const userAgent = "iwanesko-tactics-bot/0.1 (+https://iwanesko.ch)"

// FetchWeek pulls Alexandre's games from the last `days` days on both platforms.
// Bullet is intentionally INCLUDED (he plays a lot of it).
func FetchWeek(ctx context.Context, cfg Config, days int) ([]Game, error) {
	since := time.Now().AddDate(0, 0, -days)
	var games []Game

	if cfg.ChessComUser != "" {
		g, err := fetchChessCom(ctx, cfg.ChessComUser, since)
		if err != nil {
			return nil, fmt.Errorf("chess.com: %w", err)
		}
		games = append(games, g...)
	}
	if cfg.LichessUser != "" {
		g, err := fetchLichess(ctx, cfg.LichessUser, cfg.LichessToken, since)
		if err != nil {
			return nil, fmt.Errorf("lichess: %w", err)
		}
		games = append(games, g...)
	}
	return games, nil
}

// --- chess.com (Published-Data API) ----------------------------------------

type chessComArchive struct {
	Games []struct {
		PGN       string `json:"pgn"`
		TimeClass string `json:"time_class"`
		EndTime   int64  `json:"end_time"`
		Rules     string `json:"rules"`
		White     struct {
			Username string `json:"username"`
		} `json:"white"`
		Black struct {
			Username string `json:"username"`
		} `json:"black"`
	} `json:"games"`
}

func fetchChessCom(ctx context.Context, user string, since time.Time) ([]Game, error) {
	// A 7-day window can straddle a month boundary → fetch this month and last.
	now := time.Now()
	months := []time.Time{now, now.AddDate(0, -1, 0)}
	seen := map[time.Time]bool{}

	var out []Game
	for _, m := range months {
		key := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.UTC)
		if seen[key] {
			continue
		}
		seen[key] = true

		url := fmt.Sprintf("https://api.chess.com/pub/player/%s/games/%04d/%02d",
			strings.ToLower(user), m.Year(), int(m.Month()))
		var arch chessComArchive
		if err := getJSON(ctx, url, "", &arch); err != nil {
			return nil, err
		}
		for _, g := range arch.Games {
			if g.Rules != "" && g.Rules != "chess" {
				continue // skip variants
			}
			end := time.Unix(g.EndTime, 0)
			if end.Before(since) {
				continue
			}
			color := "white"
			if strings.EqualFold(g.Black.Username, user) {
				color = "black"
			}
			out = append(out, Game{
				PGN: g.PGN, MyColor: color, TimeClass: g.TimeClass,
				Source: "chesscom", EndTime: end,
			})
		}
	}
	return out, nil
}

// --- lichess (games export, NDJSON) ----------------------------------------

type lichessGame struct {
	Speed     string `json:"speed"`
	CreatedAt int64  `json:"createdAt"`
	PGN       string `json:"pgn"`
	Players   struct {
		White struct {
			User struct {
				Name string `json:"name"`
			} `json:"user"`
		} `json:"white"`
		Black struct {
			User struct {
				Name string `json:"name"`
			} `json:"user"`
		} `json:"black"`
	} `json:"players"`
}

func fetchLichess(ctx context.Context, user, token string, since time.Time) ([]Game, error) {
	url := fmt.Sprintf("https://lichess.org/api/games/user/%s?since=%d&pgnInJson=true&clocks=false&evals=false&opening=false",
		user, since.UnixMilli())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/x-ndjson")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var out []Game
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var g lichessGame
		if err := json.Unmarshal([]byte(line), &g); err != nil {
			continue
		}
		color := "white"
		if strings.EqualFold(g.Players.Black.User.Name, user) {
			color = "black"
		}
		out = append(out, Game{
			PGN: g.PGN, MyColor: color, TimeClass: g.Speed,
			Source: "lichess", EndTime: time.UnixMilli(g.CreatedAt),
		})
	}
	return out, sc.Err()
}

func getJSON(ctx context.Context, url, token string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil // no games that month
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d for %s", resp.StatusCode, url)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
