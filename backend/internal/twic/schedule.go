package twic

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"
)

// Le planificateur.
//
// Il vit DANS le serveur, en goroutine, plutôt que dans le cron du panneau
// Virtuozzo. Trois raisons, et la troisième est la vraie :
//
//   - le cron de l'hébergeur ne vit pas dans le dépôt : il se perd au prochain
//     changement de machine, et rien ne le rappelle ;
//   - il faudrait lui donner un moyen de s'authentifier sur /admin, donc écrire
//     le jeton dans une crontab, donc le sortir du seul endroit qui le garde ;
//   - et il se teste. `nextRun` est une fonction pure : la logique de « mardi
//     prochain » se vérifie sans attendre mardi.
//
// Le prix à payer, c'est qu'un conteneur éteint ne déclenche rien. D'où le
// rattrapage au démarrage : le curseur est un NUMÉRO, pas une date, donc un
// serveur rallumé après trois semaines reprend simplement au 1663.

// runHour est l'heure UTC du réveil hebdomadaire — soit 9 h en Suisse l'été,
// 8 h l'hiver. TWIC paraît le lundi soir heure britannique ; viser le mardi
// matin laisse la publication se poser.
const runHour = 7

// bootGrace : au démarrage, on ne retente que si la dernière tentative date de
// plus de ça. Le serveur redémarre à chaque publication d'article, et sans ce
// garde-fou une après-midi de corrections ferait dix requêtes à un site tenu
// par une personne.
const bootGrace = 12 * time.Hour

// Start lance la boucle. Elle s'arrête avec le contexte.
func (im *Importer) Start(ctx context.Context) {
	if !im.Enabled() {
		return
	}
	go func() {
		if im.dueAtBoot() {
			im.runLogged(ctx, "démarrage")
		}
		for {
			next := nextRun(time.Now())
			t := time.NewTimer(time.Until(next))
			slog.Info("TWIC : prochain passage", "quand", next.Format(time.RFC3339))
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				im.runLogged(ctx, "mardi")
			}
		}
	}()
}

// nextRun donne le prochain mardi à runHour, strictement après now.
//
// Pure et sans dépendance à la base de fuseaux : le conteneur n'embarque pas
// forcément tzdata, et une heure locale absente ferait silencieusement basculer
// tout le calcul en UTC — autant l'assumer.
func nextRun(now time.Time) time.Time {
	n := now.UTC()
	t := time.Date(n.Year(), n.Month(), n.Day(), runHour, 0, 0, 0, time.UTC)
	days := (int(time.Tuesday) - int(t.Weekday()) + 7) % 7
	t = t.AddDate(0, 0, days)
	if !t.After(n) {
		t = t.AddDate(0, 0, 7)
	}
	return t
}

// dueAtBoot dit s'il faut tenter un rattrapage maintenant. Elle ouvre la base
// pour lire l'horodatage de la dernière tentative — et sert au passage de
// vérification d'existence : tant que mega.db n'a pas été téléversée, il n'y a
// rien à mettre à jour, et surtout rien à créer.
func (im *Importer) dueAtBoot() bool {
	if _, err := os.Stat(im.Path); err != nil {
		slog.Info("TWIC : base de parties absente, mise à jour en sommeil", "path", im.Path)
		return false
	}
	last, err := im.lastChecked()
	if err != nil {
		slog.Warn("TWIC : état illisible, on tente quand même", "err", err)
		return true
	}
	if last.IsZero() {
		return true
	}
	return time.Since(last) > bootGrace
}

func (im *Importer) lastChecked() (time.Time, error) {
	w, err := openExisting(im.Path)
	if err != nil {
		return time.Time{}, err
	}
	defer w.Close()
	v, err := w.MetaGet(metaChecked)
	if err != nil || v == "" {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, v)
}

func (im *Importer) runLogged(ctx context.Context, cause string) {
	rep, err := im.Run(ctx)
	switch {
	case errors.Is(err, ErrBusy):
		return
	case err != nil:
		// Un échec ne fait pas tomber le serveur et ne consomme pas la semaine :
		// le curseur n'a avancé que pour les numéros réellement entrés.
		slog.Error("TWIC : mise à jour échouée", "cause", cause, "err", err,
			"curseur", rep.Last)
		return
	}
	if len(rep.Issues) == 0 {
		slog.Info("TWIC : déjà à jour", "cause", cause, "curseur", rep.Last)
		return
	}
	slog.Info("TWIC : mise à jour terminée", "cause", cause,
		"numeros", len(rep.Issues), "ajoutees", rep.Added,
		"doublons", rep.Skipped, "curseur", rep.Last)
}

// Status résume l'état pour /admin sans rien déclencher.
type Status struct {
	Enabled bool      `json:"enabled"`
	Present bool      `json:"present"`
	Last    int       `json:"last"`
	Checked time.Time `json:"checked"`
	Next    time.Time `json:"next"`
	Running bool      `json:"running"`
	// Updated est la date du dernier numéro réellement entré, Issue son bilan.
	// Tous deux vides tant que l'importeur n'a rien fait entrer — une base
	// construite sur le PC n'en dit rien.
	Updated time.Time `json:"updated"`
	Issue   *Issue    `json:"issue,omitempty"`
	// Error est l'échec du dernier passage, vide s'il a réussi.
	Error string `json:"error,omitempty"`
}

func (im *Importer) Status() Status {
	st := Status{Enabled: im.Enabled(), Last: DefaultLast, Next: nextRun(time.Now())}
	if !st.Enabled {
		return st
	}
	im.mu.Lock()
	st.Running = im.running
	im.mu.Unlock()
	if _, err := os.Stat(im.Path); err != nil {
		return st
	}
	st.Present = true
	w, err := openExisting(im.Path)
	if err != nil {
		return st
	}
	defer w.Close()
	if v, _ := w.MetaGet(MetaLast); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			st.Last = n
		}
	}
	if v, _ := w.MetaGet(metaChecked); v != "" {
		st.Checked, _ = time.Parse(time.RFC3339, v)
	}
	if v, _ := w.MetaGet(metaUpdated); v != "" {
		st.Updated, _ = time.Parse(time.RFC3339, v)
	}
	if v, _ := w.MetaGet(metaIssue); v != "" {
		var is Issue
		if json.Unmarshal([]byte(v), &is) == nil {
			st.Issue = &is
		}
	}
	st.Error, _ = w.MetaGet(metaError)
	return st
}
