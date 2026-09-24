import type { ReactNode } from 'react'
import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { SITE } from '../lib/site'
import { Container } from '../components/Container'
import { Section, Eyebrow, CtaLink } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { ChessMotif } from '../components/ChessMotif'
import { personSchema, breadcrumbSchema, type Crumb } from '../lib/schema'
import { TOURNAMENTS } from '../lib/tournaments'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

type Field = { label: string; value: string }

const COPY: Record<Locale, {
  path: string; title: string; description: string; eyebrow: string; crumb: string
  heroTitle: string; lead: ReactNode; primaryCta: string; secondaryCta: string
  /** `season` est la liste des opens 2026 avec leur score, lue dans tournaments.ts. */
  prose: (season: string) => ReactNode
  retrouvezEyebrow: string; retrouvezP: string; bottomCta: string
  photoAlt: string
  fields: Field[]
}> = {
  fr: {
    path: '/a-propos',
    title: 'À propos d’Alexandre Iwanesko, Maître FIDE',
    description:
      'Parcours d’Alexandre Iwanesko, Maître FIDE et coach d’échecs à Genève : titre, Nyon 1 en LNA, saison 2026 et méthode d’enseignement.',
    eyebrow: 'À propos',
    crumb: 'À propos',
    heroTitle: 'Alexandre Iwanesko, Maître FIDE & coach d’échecs à Genève',
    lead: (
      <>
        {SITE.person.description} Mon objectif&nbsp;: rendre votre progression{' '}
        <strong>concrète et mesurable</strong>.
      </>
    ),
    primaryCta: 'Me contacter',
    secondaryCta: 'Voir les résultats',
    prose: (season) => (
      <>
        <h2>Parcours et titre FIDE</h2>
        <p>
          <strong>Maître FIDE</strong> est un titre international décerné à vie par la Fédération
          internationale des échecs. Mon classement et la liste de mes tournois sont publics&nbsp;:
          ils vivent sur ma{' '}
          <a href="https://ratings.fide.com/profile/682136" rel="noopener" target="_blank">
            fiche FIDE officielle
          </a>
          , qui reste à jour sans que j’y touche. C’est la seule pièce qui compte, et n’importe qui
          peut la vérifier.
        </p>
        {/* TODO(alex) — l’année d’obtention du titre de Maître FIDE et les tournois où les normes
            ont été faites n’apparaissent nulle part dans le dépôt. Deux phrases de ta part
            complètent cette section. */}
        <p>
          J’ai <strong>arrêté la compétition pendant six ans</strong>, et je m’y suis remis en{' '}
          <strong>août 2025</strong>. Le bilan de cette première année est public&nbsp;:{' '}
          <strong>188 parties classées FIDE</strong> et <strong>139 points Elo perdus</strong> en
          partie lente — mes 2289 décrivaient le joueur que j’étais en 2018, pas celui qui se
          rasseyait. Je l’ai raconté en détail, chiffres et courbe à l’appui, dans{' '}
          <Link to="/blog/reprendre-les-echecs-apres-une-pause">mon bilan de reprise</Link>. Si vous
          cherchez un coach qui ne parle que de ses victoires, vous n’êtes pas à la bonne adresse.
        </p>

        <h2>Palmarès et club</h2>
        <p>
          Je joue en <strong>Ligue nationale A</strong>, la première division suisse par équipes,
          avec <strong>Nyon&nbsp;1</strong>. La saison 2026 s’est jouée sur le maintien, obtenu lors
          des rondes 6 et 7 — l’objectif de l’année, rempli en un week-end.
        </p>
        <p>
          Côté opens, la saison 2026 compte&nbsp;: {season}. Chaque tournoi a son carnet sur le blog,
          parties commentées et erreurs comprises&nbsp;; le{' '}
          <Link to="/calendrier">calendrier de compétition</Link> les réunit.
        </p>
        {/* TODO(alex) — le palmarès d’avant l’arrêt (2019 et avant) manque entièrement : titres,
            podiums, championnats, meilleur classement atteint. Confirme aussi depuis quelle
            saison tu joues à Nyon, et dans quel club tu as été formé. */}

        <h2>Enseigner les échecs</h2>
        <p>
          J’enseigne aux <strong>adultes de 1200 à 2200 Elo</strong> et aux{' '}
          <strong>adolescents en compétition</strong>, en présentiel à Genève et en ligne dans toute
          la Suisse romande. Ce sont des cours de perfectionnement&nbsp;: je ne prends pas de grands
          débutants, parce que ce n’est ni mon métier ni ce que je fais de mieux.
        </p>
        <p>
          Le résultat le plus net à ce jour est celui de Flavien&nbsp;: deux mois de travail sur la
          compréhension positionnelle, et un classement passé de{' '}
          <strong>2000 à 2100 Elo FIDE</strong>. Les autres cas suivis sont décrits un par un sur la
          page <Link to="/resultats">résultats</Link> — sans étoiles ni notes, parce que personne ne
          m’en a donné.
        </p>
        {/* TODO(alex) — l’expérience d’enseignement est le trou le plus visible de cette page :
            depuis quelle année tu enseignes, combien d’élèves tu as suivis, et si tu as encadré
            en club, en école ou pour la fédération. */}

        <h2>Ma méthode</h2>
        <p>
          Pas de recettes toutes faites. On part de vos parties classées, on isole les deux ou trois
          erreurs qui coûtent le plus cher, et on construit un plan sur huit à douze semaines avec
          des objectifs vérifiables. On progresse par la compréhension, pas par la mémorisation — et
          ma propre reprise me l’a rappelé de la pire façon&nbsp;: le jugement revient vite, le
          calcul demande du travail.
        </p>
        <p>
          Concrètement, ça donne des{' '}
          <Link to="/cours-echecs-adultes-geneve">cours d’échecs pour adultes à Genève</Link>, des{' '}
          <Link to="/cours-echecs-en-ligne">séances en ligne</Link>, et à l’approche d’un open de la{' '}
          <Link to="/preparation-tournoi-echecs">préparation de tournoi</Link>.
        </p>
      </>
    ),
    retrouvezEyebrow: 'Retrouvez-moi',
    retrouvezP:
      'Fiches officielles et profils publics : de quoi vérifier le titre, le classement et les tournois joués, sans avoir à me croire sur parole.',
    bottomCta: 'Réserver un premier cours',
    photoAlt: 'Alexandre Iwanesko, Maître FIDE et coach d’échecs à Genève',
    fields: [
      { label: 'Titre', value: 'Maître FIDE' },
      { label: 'Équipe', value: 'Nyon 1 (LNA)' },
      { label: 'Public', value: '1200–2200 Elo' },
      { label: 'Lieu', value: 'Genève / en ligne' },
      { label: 'Langues', value: 'FR · EN' },
      { label: 'ID FIDE', value: '682136' },
    ],
  },
  en: {
    path: '/en/about',
    title: 'About Alexandre Iwanesko, FIDE Master',
    description:
      'The background of Alexandre Iwanesko, FIDE Master and chess coach in Geneva: title, Nyon 1 in the top division, 2026 season and teaching method.',
    eyebrow: 'About',
    crumb: 'About',
    heroTitle: 'Alexandre Iwanesko, FIDE Master & chess coach in Geneva',
    lead: (
      <>
        FIDE Master and chess coach in Geneva, specialising in the progress of adults (1200–2200
        Elo) and competitive teenagers. My goal: to make your progress{' '}
        <strong>concrete and measurable</strong>.
      </>
    ),
    primaryCta: 'Get in touch',
    secondaryCta: 'See the results',
    prose: (season) => (
      <>
        <h2>Background and FIDE title</h2>
        <p>
          <strong>FIDE Master</strong> is an international title awarded for life by the
          International Chess Federation. My rating and the list of my tournaments are public: they
          live on my{' '}
          <a href="https://ratings.fide.com/profile/682136" rel="noopener" target="_blank">
            official FIDE profile
          </a>
          , which stays current without my touching it. It is the only document that counts, and
          anyone can check it.
        </p>
        {/* TODO(alex) — the year the FIDE Master title was awarded, and the tournaments where the
            norms were made, appear nowhere in the repository. Two sentences from you complete
            this section. */}
        <p>
          I <strong>stopped competing for six years</strong> and came back in{' '}
          <strong>August 2025</strong>. The first year is on the record:{' '}
          <strong>188 FIDE-rated games</strong> and <strong>139 rating points lost</strong> at
          classical — my 2289 described the player I was in 2018, not the one sitting back down. I
          wrote it all up, numbers and curve included, in{' '}
          <Link to="/en/blog/returning-to-chess-after-a-break">my return-to-chess review</Link>. If
          you are looking for a coach who only talks about his wins, this is the wrong address.
        </p>

        <h2>Results and club</h2>
        <p>
          I play in the <strong>Swiss National League A</strong>, the top division of the Swiss team
          championship, with <strong>Nyon&nbsp;1</strong>. The 2026 season came down to survival,
          secured in rounds 6 and 7 — the objective of the year, met in a single weekend.
        </p>
        <p>
          On the open circuit, the 2026 season reads: {season}. Every tournament has its diary on the
          blog, annotated games and mistakes included; the{' '}
          <Link to="/en/calendar">competition calendar</Link> gathers them.
        </p>
        {/* TODO(alex) — everything before the break (2019 and earlier) is missing: titles, podiums,
            championships, peak rating. Also confirm which season you joined Nyon, and the club
            where you learned. */}

        <h2>Teaching chess</h2>
        <p>
          I coach <strong>adults from 1200 to 2200 Elo</strong> and{' '}
          <strong>competitive teenagers</strong>, in person in Geneva and online across
          French-speaking Switzerland. These are improvement lessons: I do not take complete
          beginners, because it is neither my job nor what I do best.
        </p>
        <p>
          The clearest result so far is Flavien’s: two months of work on positional understanding,
          and a rating that went from <strong>2000 to 2100 FIDE</strong>. The other coached cases are
          described one by one on the <Link to="/en/results">results page</Link> — with no stars and
          no ratings, because nobody gave me any.
        </p>
        {/* TODO(alex) — teaching experience is the most visible gap on this page: since which year
            you have been teaching, how many students you have coached, and whether you have worked
            for a club, a school or the federation. */}

        <h2>My method</h2>
        <p>
          No ready-made recipes. We start from your rated games, isolate the two or three mistakes
          that cost the most, and build a plan over eight to twelve weeks with checkable goals. You
          improve through understanding, not memorisation — and my own comeback reminded me of that
          the hard way: judgement returns quickly, calculation takes work.
        </p>
        <p>
          In practice that means{' '}
          <Link to="/en/adult-chess-lessons-geneva">adult chess lessons in Geneva</Link>,{' '}
          <Link to="/en/online-chess-lessons">online sessions</Link>, and, with an open coming up,{' '}
          <Link to="/en/tournament-preparation">tournament preparation</Link>.
        </p>
      </>
    ),
    retrouvezEyebrow: 'Find me',
    retrouvezP:
      'Official records and public profiles: enough to verify the title, the rating and the tournaments played, without taking my word for it.',
    bottomCta: 'Book a first lesson',
    photoAlt: 'Alexandre Iwanesko, FIDE Master and chess coach in Geneva',
    fields: [
      { label: 'Title', value: 'FIDE Master' },
      { label: 'Team', value: 'Nyon 1 (top division)' },
      { label: 'Audience', value: '1200–2200 Elo' },
      { label: 'Location', value: 'Geneva / online' },
      { label: 'Languages', value: 'FR · EN' },
      { label: 'FIDE ID', value: '682136' },
    ],
  },
}

/**
 * Les opens de la saison, lus dans tournaments.ts plutôt que recopiés ici.
 *
 * Le filtre garde ce qui a un score ET un lieu : un match par équipes n'a pas de
 * lieu dans tournaments.ts, un tournoi encore à jouer n'a pas de score. La page
 * ne peut donc pas annoncer un résultat que le calendrier contredit.
 */
function seasonResults(year: string, locale: Locale): string {
  return TOURNAMENTS.filter((tn) => tn.result && tn.location && tn.start.startsWith(year))
    .slice()
    .sort((a, b) => (a.start < b.start ? -1 : 1))
    .map((tn) => `${(locale === 'en' ? tn.nameEn : tn.name) ?? tn.name} (${tn.result})`)
    .join(', ')
}

export function Component() {
  const locale = useLocale()
  const c = COPY[locale]
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: c.crumb, path: c.path },
  ]

  return (
    <>
      <Seo
        title={c.title}
        description={c.description}
        path={c.path}
        jsonLd={[personSchema(), breadcrumbSchema(crumbs)]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow={c.eyebrow}
        title={c.heroTitle}
        lead={c.lead}
        primaryCta={{ to: pathFor('contact', locale), label: c.primaryCta }}
        secondaryCta={{ to: pathFor('resultats', locale), label: c.secondaryCta }}
      />

      <Section>
        <Container className="grid gap-12 lg:grid-cols-[1.5fr_1fr] lg:items-start">
          <div>
          <div className="prose">
            {c.prose(seasonResults('2026', locale))}
          </div>

          <div className="mt-10 rounded-2xl border border-ink-200/80 bg-cream-100 p-7 shadow-soft">
            <Eyebrow>{c.retrouvezEyebrow}</Eyebrow>
            <p className="text-ink-600">
              {c.retrouvezP}
            </p>
            <ul className="mt-5 flex flex-wrap gap-3 text-sm font-medium">
              {SITE.person.sameAs.map((url) => (
                <li key={url}>
                  <a
                    href={url}
                    rel="me noopener"
                    target="_blank"
                    className="inline-flex rounded-full border border-ink-300 bg-paper px-4 py-2 text-ink-700 transition-colors hover:border-gold-400 hover:text-gold-700"
                  >
                    {new URL(url).hostname.replace('www.', '')}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div className="mt-10">
            <CtaLink to={pathFor('reserver', locale)} variant="primary">
              {c.bottomCta}
            </CtaLink>
          </div>
          </div>

          {/* Visual aside — portrait (ou motif de marque) + repères vérifiables. */}
          <aside className="lg:sticky lg:top-24">
            <div className="overflow-hidden rounded-3xl border border-ink-200/80 bg-paper p-6 shadow-card">
              {/* TODO(alex) — la photo manque, et c'est la page E-E-A-T du site : un visiteur
                  qui hésite à réserver veut voir à qui il parle. Dépose un portrait dans
                  public/ (carré ou paysage, ≥ 800 px) et renseigne SITE.person.image : il
                  s'affiche ici, et personSchema l'émet en Person.image. Le motif de marque
                  tient la place en attendant — un logo ne peut pas faire office de portrait. */}
              {SITE.person.image ? (
                <img
                  src={SITE.person.image}
                  alt={c.photoAlt}
                  width={512}
                  height={512}
                  className="mx-auto w-full rounded-2xl object-cover"
                />
              ) : (
                <ChessMotif className="mx-auto w-full max-w-[16rem] text-ink-900" />
              )}
              <dl className="mt-6 grid grid-cols-2 gap-x-4 gap-y-4">
                {c.fields.map((f) => (
                  <div key={f.label}>
                    <dt className="text-[0.68rem] font-semibold uppercase tracking-[0.12em] text-ink-500">
                      {f.label}
                    </dt>
                    <dd className="mt-1 font-display font-bold text-ink-900">{f.value}</dd>
                  </div>
                ))}
              </dl>
            </div>
          </aside>
        </Container>
      </Section>
    </>
  )
}
