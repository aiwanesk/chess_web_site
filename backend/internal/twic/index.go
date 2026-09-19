package twic

import (
	"context"
	"io"
	"runtime"
	"sync"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// Construction de mega.db depuis un PGN.
//
// C'est le même travail que l'import hebdomadaire, à l'échelle du dessus : la
// MegaBase fait des millions de parties au lieu de sept mille. D'où le seul
// écart de forme — la traduction SAN→UCI est répartie sur tous les cœurs, parce
// que c'est elle qui coûte, et une seule transaction ne peut pas porter dix
// millions de lignes.
//
// Le point qui compte n'est pas la vitesse, c'est que ce soit CE code : la
// même lecture de PGN, le même générateur de coups, la même empreinte que
// l'importeur TWIC. Un indexeur écrit à côté, dans un autre langage, serait
// une seconde implémentation de l'empreinte — et le jour où elle diverge d'un
// espace ou d'un accent, le rattrapage réimporte tout en double sans rien dire.

// IndexStats compte ce qui est passé. Variants et Rejected sont séparés pour la
// même raison que dans l'import : le premier est du tri voulu, le second une
// alerte.
type IndexStats struct {
	Read     int64
	Added    int64
	Skipped  int64
	Variants int64
	Rejected int64
}

// IndexOptions règle le débit. Les valeurs nulles prennent des défauts sains.
type IndexOptions struct {
	Workers int // traducteurs SAN→UCI en parallèle ; défaut : nombre de cœurs
	Batch   int // parties par transaction ; défaut : 20 000
	// Progress est appelée de temps en temps depuis la goroutine d'écriture.
	// Elle doit être rapide : elle retient l'insertion.
	Progress func(IndexStats)
	Every    int // fréquence d'appel de Progress, en parties lues ; défaut : 100 000
}

// Index lit un flux PGN et remplit la base.
//
// Rejouable : l'index UNIQUE sur l'empreinte fait qu'un second passage sur le
// même fichier ne compte que des doublons. C'est aussi la reprise après
// interruption — on relance, et seules les parties manquantes entrent.
func Index(ctx context.Context, r io.Reader, w *games.Writer, opt IndexOptions) (IndexStats, error) {
	if opt.Workers <= 0 {
		opt.Workers = runtime.NumCPU()
	}
	if opt.Batch <= 0 {
		opt.Batch = 20000
	}
	if opt.Every <= 0 {
		opt.Every = 100000
	}

	raws := make(chan rawGame, opt.Workers*64)
	type result struct {
		game    games.ImportGame
		variant bool
		bad     bool
	}
	done := make(chan result, opt.Workers*64)

	// Lecture : une seule goroutine, séquentielle par nature.
	go func() {
		defer close(raws)
		scanGames(r, func(g rawGame) bool {
			select {
			case raws <- g:
				return true
			case <-ctx.Done():
				return false
			}
		})
	}()

	// Traduction : c'est là que passe le temps, donc c'est là qu'on parallélise.
	var wg sync.WaitGroup
	for i := 0; i < opt.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for g := range raws {
				conv, err := convert(g)
				switch {
				case err == errVariant:
					done <- result{variant: true}
				case err != nil:
					done <- result{bad: true}
				default:
					done <- result{game: conv}
				}
			}
		}()
	}
	go func() { wg.Wait(); close(done) }()

	// Écriture : une seule goroutine. Le Writer garde des caches de joueurs et
	// d'événements sans verrou — ils n'ont qu'un propriétaire, et c'est ici.
	var st IndexStats
	tx, err := w.Begin()
	if err != nil {
		return st, err
	}
	inBatch := 0
	for res := range done {
		st.Read++
		switch {
		case res.variant:
			st.Variants++
		case res.bad:
			st.Rejected++
		default:
			added, err := tx.Insert(res.game)
			if err != nil {
				st.Rejected++
				break
			}
			if added {
				st.Added++
			} else {
				st.Skipped++
			}
			inBatch++
		}
		if inBatch >= opt.Batch {
			if err := tx.Commit(); err != nil {
				return st, err
			}
			if tx, err = w.Begin(); err != nil {
				return st, err
			}
			inBatch = 0
		}
		if opt.Progress != nil && st.Read%int64(opt.Every) == 0 {
			opt.Progress(st)
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return st, err
	}
	if opt.Progress != nil {
		opt.Progress(st)
	}
	return st, ctx.Err()
}
