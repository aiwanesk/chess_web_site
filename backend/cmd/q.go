package main

import (
	"fmt"
	"os"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

func main() {
	s, err := games.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer s.Close()
	for _, q := range []string{"carlsen", "nguyen", "grunfeld"} {
		ps, _ := s.SearchPlayers(q, 3)
		for _, p := range ps {
			fmt.Printf("%-10s → %-32s (%d parties)\n", q, p.Name, p.Games)
		}
	}
	ps, _ := s.SearchPlayers("carlsen, m", 1)
	if len(ps) > 0 {
		tr, err := s.OpeningTree(games.Filter{PlayerIDs: []int64{ps[0].ID}, Colour: games.White},
			games.TreeOptions{MaxDepth: 1, MinGames: 5})
		fmt.Println("err:", err)
		for _, n := range tr {
			fmt.Printf("  %-5s %4d parties  %.0f%%\n", n.SAN, n.Games,
				100*float64(2*n.Wins+n.Draws)/float64(2*n.Games))
		}
	}
}
