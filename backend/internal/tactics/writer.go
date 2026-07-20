package tactics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WeekFile is the published shape read by the frontend.
type WeekFile struct {
	Week    string   `json:"week"` // e.g. "2026-S30"
	Puzzles []Puzzle `json:"puzzles"`
}

// WritePuzzles serialises the week's puzzles to dir/<week>.json.
//
// Defense-in-depth: before writing, it refuses if any `forbidden` string (the
// configured usernames) appears in the output — so even a future bug can never
// publish an identifier.
func WritePuzzles(dir, week string, puzzles []Puzzle, forbidden []string) (string, error) {
	data, err := json.MarshalIndent(WeekFile{Week: week, Puzzles: puzzles}, "", "  ")
	if err != nil {
		return "", err
	}
	low := strings.ToLower(string(data))
	for _, f := range forbidden {
		if f != "" && strings.Contains(low, strings.ToLower(f)) {
			return "", fmt.Errorf("refus d'écrire : un identifiant apparaît dans la sortie du batch")
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, week+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
