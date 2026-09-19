# Chess Web Site — Claude Code Memory

## Project Structure
- **Framework**: React (Vite/SSR)
- **Blog articles FR**: `content/blog/*.md`
- **Blog articles EN**: `content/blog/en/*.md`
- **Diagrams script**: `scripts/gen-diagrams.mjs` — generates SVG chess diagrams from FEN using Lichess cburnett piece set
- **Diagram output**: `public/images/blog/<tournament>/`
- **i18n**: FR (default) + EN, routes defined in `frontend/src/lib/i18n.tsx`
- **Blog slugs**: `frontend/src/lib/blogSlugs.ts` — imports from `content/blog/*.md` (FR) and `content/blog/en/*.md` (EN)
- **Content processing**: `frontend/src/lib/content.ts`
- **Routes**: `frontend/src/routes.tsx` — FR at `/blog/:slug`, EN at `/en/blog/:slug`
- **altSlug**: frontmatter field to link FR ↔ EN versions (enables hreflang)

## Thème clair / sombre
- **Interrupteur** : bouton lune/soleil dans l'en-tête (`frontend/src/components/ThemeToggle.tsx`). Le choix est mémorisé dans `localStorage`; sans choix, on suit `prefers-color-scheme`.
- **L'état vit sur `<html data-theme>`**, posé par le script inline de `frontend/index.html` **avant le premier rendu** — indispensable sur un site pré-rendu, sinon la page clignote en blanc avant de basculer.
- **Comment ça bascule** : aucune classe `dark:` dans les composants. Les utilitaires Tailwind v4 compilent vers `var(--color-…)`, donc il suffit de **redéfinir les tokens** sous `html[data-theme='dark']` dans `styles.css`. L'échelle `ink` et l'échelle `cream` sont inversées ; les ors 200→500 et la famille `slab` ne bougent pas.
- **Trois tokens de rôle à connaître** — la palette ne peut plus servir deux rôles à la fois :
  - `paper` = ce qui était codé en dur `bg-white` (cartes, en-tête, sections). **Ne jamais réintroduire `bg-white`** pour une surface.
  - `slab-*` = les panneaux qui restent **sombres dans les deux thèmes** (pied de page, bandeaux CTA) et le texte posé dessus (`text-slab-300/400`). Un `text-ink-*` sur un slab devient illisible en mode nuit.
  - `on-gold` = le texte d'un bouton doré, qui doit rester quasi noir en permanence.
- **Le bouton n'a aucun état React** : les deux icônes sont dans le DOM et le CSS choisit laquelle sort. C'est ce qui évite le saut d'hydratation.
- **Attention en test** : Chrome headless annonce `prefers-color-scheme: dark` par défaut — une capture « claire » sort sombre si on ne force pas le thème.

## Tournament Calendar
- **Tournament list (edit this)**: `frontend/src/lib/tournaments.ts` — the single source of truth, edited by hand when publishing. Baked in at build, so every tournament ships inside the pre-rendered HTML (no runtime fetch, no backend involved).
- **Page**: `/calendrier` (FR) + `/en/calendar` (EN) — `frontend/src/pages/Calendrier.tsx`
- **Grid component**: `frontend/src/components/TournamentCalendar.tsx` — month grid where a multi-day tournament renders as one continuous bar, plus « previous / next » cards and a full-season list below
- **Grid/date logic**: `frontend/src/lib/calendar.ts` (`buildMonth`, `formatRange`, `statusOf`)
- **Linking a diary**: set `slugFr` / `slugEn` to the article file name (no `/blog/` prefix). A past tournament with a slug becomes clickable straight to its diary.
- **Adding a tournament**: one entry in `TOURNAMENTS`; order doesn't matter, the grid sorts itself. `start === end` for a one-day event.
- **`__BUILD_DATE__`**: injected by `define` in `vite.config.ts`. The calendar's first render uses it so SSG output and hydration match; the visitor's real date takes over in an effect.

## Student Results (`/resultats`)
- Presents **coached cases** (what was worked on → what it produced), described by Alexandre — **not** reviews written by students.
- **No star ratings, no `aggregateRatingSchema`** on this page. Declaring ratings nobody gave is fabricated structured data: manual-action risk with Google, and misleading advertising under Swiss LCD art. 3 / the EU Omnibus directive.
- The builder in `lib/schema.ts` is kept for the day a student sends a **real written review** — add the quote, then wire the schema back.
- Every figure on the page must be traceable to something real (Flavien's +100 FIDE Elo, the tournament count read from `tournaments.ts`).

## Blog Writing Style (Tournament Diaries)
- **VOICE — read first**: Before writing OR rewriting any tournament diary (`content/blog/**`, FR and EN), read [`docs/voix-carnet-tournoi.md`](docs/voix-carnet-tournoi.md) and apply it on the FIRST draft (don't write "clean" then fix later). It is the canonical tone reference: bans reporting connectors ("À partir de là", "Le problème c'est que", "Vient alors le moment-clé", "Nouveau carrefour", "Résultat :"…), requires showing emotions via concrete detail rather than declaring them, favours self-deprecation, reader address and lived imagery. Gold standard: `content/blog/open-pontevedra-2026.md`. Does NOT apply to SEO/informational articles (see `docs/plan-editorial-blog.md`).
- **Tone**: Self-deprecating humor, storytelling, honest about mistakes
- **Framework per round**:
  1. How you felt before the game (fatigue, motivation, mental state)
  2. How the opening went
  3. The moment you felt the advantage (or the pressure)
  4. What was going through your head at key moments
  5. Why this game mattered (or didn't)
  6. Self-deprecation — never take yourself too seriously
  7. Funny expressions when they fit naturally
- **Chess notation**: French style (Fou=F, Cavalier=C, Tour=T, Dame=D, Roi=R) for FR articles, English style (B, N, R, Q, K) for EN articles
- **Played moves vs variations**: a move actually played is written `` `20.Tc2` ``; a move only calculated is written `` *`23.g4`* `` — it renders as `<em><code>` and is greyed with no background (rule `.prose em code` in `frontend/src/styles.css`). Never let a reader mistake a variation for the game.
- **Engine verdicts**: when a diary separates the player's own calculations from the engine's, the engine goes in a grey callout labelled « L'ordinateur, après coup » (`<div class="engine-note">` + `<p class="engine-note-label">` ; see `championnat-suisse-equipes-2026.md`), never inline in the prose.
- **Diagrams**: Multiple per round at narrative turning points, with witty captions. Markup: `<div class="diagram-container">` + `<img …>` + `<p class="diagram-caption">`. **Rien à écrire pour la performance** : `enhanceDiagrams()` (`content.ts`) ajoute `loading="lazy"` et `decoding="async"` à toute image de `/images/blog/`, et `aspect-ratio` sur `.diagram-container img` réserve la hauteur (sans quoi différer décaleraient l'ancre `#parties`). Une balise qui porte déjà `loading=` n'est jamais touchée. **Never inline styles** — a style attribute cannot follow the light/dark theme, and these blocks stayed white rectangles in night mode until they were moved to classes. Same for the summary table: rows use `class="row-head"` / `row-a` / `row-b`, cells `cell-muted` / `cell-warn` / `cell-alt`.
- **Lessons**: Short, concrete, honest — not generic advice
- **Author**: Alexandre Iwanesko, FM (FIDE Master), 33 years old

## Fréquentation (/admin)
- **Les humains sont comptés par une balise JS**, pas par le serveur : `frontend/src/lib/analytics.ts` envoie `POST /api/hit` à chaque changement de route, `handleHit` (`backend/internal/server/analytics.go`) l'enregistre.
- **Pourquoi** : le comptage serveur enregistrait toute réponse HTML, donc tous les crawlers déguisés en navigateur. Le tableau de bord affichait des milliers de « visites humaines » venues de datacenters US/DE/NL/SG, pour 3 clics réels dans Search Console. Les robots n'exécutent presque jamais de JS.
- Le middleware serveur ne compte plus que les **robots qui s'annoncent** (regex user-agent), ce qui garde la colonne « bots » lisible.
- Garde-fous sur `/api/hit` : même origine exigée (`Origin`, sinon `Referer`), user-agent de robot rejeté, chemin validé par `cleanHitPath` (interne, pas `/admin` ni `/newsletter`). Testé dans `server_test.go`.
- Rien ne change côté vie privée : aucune IP stockée, aucun cookie, pays déduit hors ligne, empreinte de visiteur tournante à la journée.

## Tableaux dans les articles
- **Écrire les tableaux en Markdown**, pas en HTML. `enhanceTables()` (`frontend/src/lib/content.ts`) les enveloppe dans un `.table-wrap` et colore les cellules qui ne contiennent qu'un nombre signé (`+3,6` en vert, `−59,0` en rouge).
- Le style vit dans `styles.css`, **scopé sous `.table-wrap`** : en-tête presque noir, lignes alternées, chiffres tabulaires, dernière colonne alignée à droite. Rien de tout ça ne touche les tableaux HTML écrits à la main dans les carnets, qui gardent leur mise en forme.
- Au-delà de 768 px, un tableau **déborde volontairement de la colonne de texte** (68 caractères) et se recentre sur la page : six colonnes n'y entrent pas, et c'est toujours la dernière — la plus utile — qui se retrouvait hors champ.
- **Pas de drapeaux emoji** : Windows ne les rend pas et les affiche en paires de lettres (« CH », « FR »).

## PGN Game Viewer
- **Games live in**: `content/games/<id>.pgn` — one file per game, standard PGN with headers. Paste the score, done: no build step, no diagram generation.
- **Embed in an article**: put `[[pgn:<id>]]` alone on its own line in the Markdown (FR or EN). `splitChunks()` in `frontend/src/lib/content.ts` splits the compiled HTML on that marker and `BlogPost.tsx` renders a `<PgnViewer>` in its place.
- **Several games, one board**: `[[pgn:<id-a>,<id-b>]]` renders a single viewer with a picker (tab per game, labelled `R<ronde> · <adversaire>`). Preferred for a team-match diary: one board at the foot of the article rather than one per round.
- **Anchor**: the first `[[pgn:…]]` block of a post gets `id="parties"`, so the intro can link down to it with `[texte](#parties)`.
- **Component**: `frontend/src/components/PgnViewer.tsx` — board, first/prev/play/next/last, **flip button**, clickable move list, arrow-key navigation. Read-only.
- **Rules engine**: `frontend/src/lib/chess.ts` (legal-move generation, only used to resolve SAN into from/to squares) + `frontend/src/lib/pgn.ts` (PGN → plies, main line only).
- **Main line only**: comments `{...}`, sidelines `(...)` and NAGs are stripped. A sideline worth telling goes in the prose, not in the viewer.
- **Notation**: the viewer renders French notation (C, F, T, D, R) on FR pages and English on EN pages automatically — always paste the PGN in **English** SAN.
- **Board orientation**: defaults to Alexandre's side (detected from the `White`/`Black` headers), overridable with a `[Orientation "black"]` header. The reader can flip it anyway.
- **Piece set**: `frontend/src/components/chessPieces.tsx` — shared with `PuzzleBoard.tsx` (Lichess cburnett, inlined SVG).

## Newsletter
- **Quand partent les mails** : l'annonceur (`backend/internal/server/announce.go`) tourne **une seule fois, au démarrage du serveur**. Pas de cron. Publier = redéployer, et le redémarrage envoie tout ce qui n'est pas encore dans la table `notified`. Il faut `DB_PATH` **et** le SMTP configurés, sinon il ne fait rien.
- Au **tout premier lancement**, il marque tout le catalogue comme déjà notifié : déployer la fonctionnalité n'envoie jamais d'e-mail rétroactif.
- Ce qu'il annonce : les **articles des deux langues** (`content/blog/*.md` → `blog:<slug>`, `content/blog/en/*.md` → `blogEN:<slug>`), les **tactiques hebdo** et les **événements** de `content/events.json`. Seuls les **confirmés** reçoivent quelque chose, un envoi toutes les 300 ms.
- **Chaque article porte sa langue** et ne part qu'aux abonnés de cette langue. Conséquence assumée : un article FR **sans traduction** n'atteint pas les abonnés EN — mieux que de leur envoyer un lien français.
- **Élargir la collecte à un nouveau lot exige un amorçage**, sinon le redémarrage poste tout l'arriéré d'un coup. `seedKind()` marque l'existant comme vu sans rien envoyer, une seule fois, derrière une sentinelle (`__seeded:blogEN__`). Faire pareil pour tout futur type d'item.
- **Onglet Newsletter du `/admin`** : liste, statut, langue, dates, taux de confirmation. Les **jetons ne sont jamais affichés** — un `unsub_token` est une capacité, il désabonne sans autre preuve, et un tableau de bord se photographie.
- Se désinscrire **supprime la ligne** : un ancien abonné ne laisse aucune trace dans le tableau.

## Silo FR/EN des pages argent
- Une page argent a **deux objets de configuration**, `FR` et `EN`. `cluster` n'était posé que sur `FR` : les pages EN n'ont donc jamais affiché d'articles liés, même quand la traduction existait. **En ajouter un à l'une, l'ajouter à l'autre.**
- `CoursEnLigne.tsx` n'a de `cluster` sur aucune des deux : le cluster `en-ligne` du plan éditorial n'a encore aucun article.

## Cartes de partage (Open Graph) et JSON-LD
- **`scripts/gen-og.mjs`** produit les PNG 1200×630. Deux gabarits : `card()` pour un **carnet** (un vrai diagramme d'échiquier imbriqué depuis `public/images/blog/<dir>/`, donc la position montrée est celle de l'article) et `chartCard()` pour un **article de fond** (la courbe Elo, **lue directement dans le Markdown** — carte et graphique ne peuvent pas diverger).
- Le SVG reste à côté du PNG : c'est la source éditable. **C'est le PNG qui va dans le front-matter** — les réseaux sociaux n'affichent pas d'aperçu SVG. Carnets : `/images/blog/<dir>/og-*.png`. Articles de fond : `/public/og/<slug>.png`.
- Sans `image:` en front-matter, la page retombe sur `og/default.png` — silencieusement. `check-content-links.mjs` vérifie qu'une `image:` déclarée existe, **pour tout article** et plus seulement pour les carnets.
- **JSON-LD** (`frontend/src/lib/schema.ts`, `articleSchema`) : une page d'article ne porte aucun nœud `Person` ni `Organization` (ils ne sont émis que sur l'accueil et `/a-propos`). L'`@id` seul y serait une **référence pendante** — le nom voyage donc avec. `inLanguage` suit la locale de la page (un article EN se déclarait `fr`), et `image` a toujours une valeur.

## Espace privé : mise à jour TWIC de la base de parties
- **Le planificateur est en Go, dans le processus** (`backend/internal/twic/schedule.go`) — pas dans le cron de l'hébergeur. Il réveille l'import **tous les mardis à 7 h UTC**, et **rattrape au démarrage** (garde-fou : pas deux tentatives à moins de 12 h d'écart, sinon une après-midi de redéploiements harcèlerait un site tenu par une personne). `nextRun` est une fonction pure, donc « mardi prochain » se teste sans attendre mardi.
- **Deux garde-fous indépendants contre le doublon.** Le curseur `meta.twic_last` **dans mega.db** dit où on en est (à défaut : reprise au **1640**) ; l'index UNIQUE sur `game.hash` absorbe tout ce qui repasserait malgré lui. Le second n'est utile que si les lignes déjà présentes portent l'empreinte : **l'indexeur qui construit mega.db doit remplir `hash` avec la formule de `games.Key`** (documentée en Python dans `import.go`), sinon un rattrapage réimporte.
- Mettre le curseur **dans la base** et pas à côté : téléverser une mega.db reconstruite sur le PC remplace aussi le curseur, donc rien à régler après coup.
- **Le Chess960 est écarté explicitement**, pas compté comme illisible (`errVariant`). Une partie qui ne part pas de la position initiale ne se raccroche à aucune branche de l'arbre d'ouvertures — mélangée aux vraies, elle ferait croire à de la théorie là où il n'y en a pas. Le compteur `rejected`, lui, doit rester à **zéro** : chaque unité est une partie que le lecteur de SAN n'a pas su rejouer.
- **L'importeur ne crée JAMAIS la base.** Sans mega.db téléversée il se met en sommeil : `games.OpenWriter` poserait sinon un schéma vide, et l'explorateur répondrait « aucune partie » au lieu de « aucune base ».
- Déclenchement manuel : bouton dans l'onglet **Bases** de `/admin` (`POST /admin/parties/twic`) — pour la veille d'un tournoi. Il tourne sur `context.Background()`, pas sur la requête : fermer l'onglet n'interrompt pas un import.
- **Test facultatif contre le vrai site** : `TWIC_LIVE=1640-1662 go test ./internal/twic/ -run Live -v`. C'est la seule vérification que l'adresse des archives n'a pas bougé et que le lecteur encaisse du PGN réel (174 919 parties rejouées, 0 illisible au 19.09.2026).

## Explorateur de parties (`/admin/parties/`)
- **L'arbre se prolonge à la demande.** Le serveur en renvoie 14 demi-coups à la fois (sept coups) ; quand on atteint le bout d'une branche, le client recharge depuis le chemin courant (`/api/tree?path=…`) et greffe. C'est ce qui lève la limite sans jamais transférer un gros arbre — descendu à 34 demi-coups en test.
- **Une ligne de la liste s'ouvre en partie entière.** `GET /api/game?id=` renvoie la partie avec **le FEN après chaque demi-coup** : le `san` et l'`uci` complets sont déjà en base, et les positions sont calculées côté serveur. Le navigateur n'a toujours aucune règle du jeu à connaître — c'est le même principe que pour l'arbre, la règle vit d'un seul côté.
- La partie devient une **chaîne de nœuds à un seul enfant** : la navigation de l'arbre (avancer / reculer / cliquer la ligne) fonctionne dessus sans code séparé. L'échiquier s'oriente du côté du joueur cherché.
- **L'import TWIC ne ferme PAS la base en lecture.** Le mode WAL est fait pour ça. La fermer rendait l'explorateur muet (503) pendant les deux minutes du rattrapage initial — pile au moment où on vient de téléverser et où on veut vérifier. Seul `uploadMu` est pris, pour qu'un téléversement ne bascule pas le fichier en plein import ; `reopen` remet une connexion neuve après coup.

## Construire mega.db (`backend/cmd/megaindex`)
- **L'indexeur est en Go, pas en Python.** Il réutilise le lecteur de PGN, le générateur SAN→UCI et l'empreinte `games.Key` de l'import TWIC : un indexeur écrit à côté serait une **seconde implémentation de l'empreinte**, et le jour où elle diverge d'un espace ou d'un accent, le rattrapage hebdomadaire réimporte tout en double sans rien dire.
- Il tourne **sur le PC** : `go build -o megaindex.exe ./cmd/megaindex` puis `megaindex.exe -out mega.db "D:ases\*.zip"`. Accepte `.pgn` et `.zip`, développe les jokers lui-même (cmd.exe ne le fait pas).
- **Rejouable** : relancer sur le même fichier ne compte que des doublons. C'est aussi la reprise après un Ctrl+C — on relance, seul ce qui manque entre.
- Mesuré : **~8 900 parties/s** (23 291 parties en 3 s, 0 illisible). Une MegaBase de 11 millions se construit en une vingtaine de minutes.
- **Ne pas poser `-twic-last`** sauf certitude : la clé écrite empêche le rattrapage automatique du serveur depuis 1640.
- `scanGames` lit **en flux** et décode le Latin-1 **ligne par ligne** : décider sur le fichier entier transformerait tous les vrais caractères UTF-8 d'un fichier presque propre en charabia.

## SAN → UCI (`backend/internal/corpus/san.go`)
- Générateur de coups **légaux** écrit pour l'import PGN : TWIC n'écrit que du SAN, et « Cbd2 » ne devient « b1d2 » qu'en sachant quels cavaliers peuvent vraiment y aller.
- **La validation vient de `corpus.db`**, qui stocke pour chaque arête **le SAN ET l'UCI** (produits par python-chess) : le test rejoue le graphe et compare. 47 931 coups, 1241 roques, 21 prises en passant, 1147 désambiguïsations. C'est cet oracle qui a trouvé le bug du `scan` de pièces glissantes — un `return` au premier obstacle **toutes directions confondues** laissait bouger une pièce clouée.
- Le filtre de légalité n'est pas du luxe : c'est lui qui rend « Ne2 » non ambigu quand l'autre cavalier est cloué. Un simple filtre pseudo-légal échoue sur la ronde 4 d'un tournoi sur deux.

## Diagram Generation
- Run `node scripts/gen-diagrams.mjs` after modifying FEN positions
- Each diagram entry: `{ file, fen, lastMove?, flip?, dir? }`
- `dir` picks the tournament folder under `public/images/blog/` (default `pontevedra-2026`); CSE 2026 uses `dir: 'cse-2026'`
- `flip: true` for games played as Black
- `lastMove` format: `'e2e4'` (from-to squares)

## Tournament Schedule (2026)
- Pontevedra: July 25–30, 2026 (completed, 4.5/9)
- Badalona: next tournament (directly after Pontevedra)
