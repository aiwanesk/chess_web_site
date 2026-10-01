package server

import (
	"archive/zip"
	"compress/flate"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// Export de mega.db en PGN, depuis l'onglet Bases de /admin.
//
// Un ZIP et pas un .pgn.gz : Windows l'ouvre sans rien installer, et c'est ce
// que megaindex relit déjà. Il est produit À LA VOLÉE — rien n'est écrit sur le
// disque du conteneur, qui n'a pas la place de dupliquer une base de 11 Go.

// exportPath est servi sans passer par le compresseur global : du Brotli
// par-dessus du Deflate ne gagnerait rien, et surtout son ResponseWriter ne
// sait pas lever la limite d'écriture (voir noLimit).
const exportPath = "/admin/parties/export"

func (s *Server) handlePartiesExport(w http.ResponseWriter, r *http.Request) {
	if !s.exportMu.TryLock() {
		// Deux exports simultanés, c'est deux fois tout le disque à lire pour
		// le même fichier — sans doute un double clic.
		s.corpusError(w, "un export est déjà en cours", http.StatusConflict)
		return
	}
	defer s.exportMu.Unlock()

	if _, err := os.Stat(s.cfg.GamesDB); err != nil {
		s.corpusError(w, "aucune base de parties en ligne", http.StatusServiceUnavailable)
		return
	}
	// Une connexion À SOI plutôt que celle de l'explorateur : un téléversement
	// ou la fin d'un import TWIC la ferment pour la rouvrir, ce qui couperait
	// l'export en plein milieu. Celle-ci garde son instantané jusqu'au bout.
	g, err := games.Open(s.cfg.GamesDB)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer g.Close()

	noLimit(w)
	privateHeaders(w)
	name := "mega-" + time.Now().Format("2006-01-02")
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`.zip"`)

	zw := zip.NewWriter(w)
	// BestSpeed : le goulot est le processeur du conteneur, pas la bande
	// passante. Le PGN se compresse assez bien pour que les niveaux plus lents
	// ne gagnent que quelques pour cent, au prix de minutes de plus.
	zw.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestSpeed)
	})
	f, err := zw.CreateHeader(&zip.FileHeader{
		Name: name + ".pgn", Method: zip.Deflate, Modified: time.Now(),
	})
	if err != nil {
		return
	}
	start := time.Now()
	n, err := g.Export(r.Context(), f)
	if err != nil {
		// Les en-têtes sont partis, on ne peut plus changer le statut. Ne PAS
		// fermer le zip : sans répertoire central, l'archive tronquée est
		// refusée à l'ouverture au lieu de passer pour une base complète.
		slog.Warn("export PGN interrompu", "parties", n, "err", err)
		return
	}
	if err := zw.Close(); err != nil {
		slog.Warn("export PGN : fin d'archive", "err", err)
		return
	}
	slog.Info("export PGN terminé", "parties", n, "duree", time.Since(start).Round(time.Second))
}

// noLimit lève la limite d'écriture du serveur (WriteTimeout, 30 s) pour CETTE
// réponse seulement. Elle protège le site d'un client qui lit au compte-gouttes,
// mais elle couperait un export de plusieurs minutes — silencieusement, en
// laissant au navigateur un fichier qui a l'air complet jusqu'à ce qu'on
// l'ouvre.
func noLimit(w http.ResponseWriter) {
	if err := http.NewResponseController(w).SetWriteDeadline(time.Time{}); err != nil {
		slog.Warn("limite d'écriture non levée : la réponse sera coupée à 30 s", "err", err)
	}
}
