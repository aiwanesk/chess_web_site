/**
 * Calendrier des tournois — la source de vérité, éditée à la main.
 *
 * ▶ POUR AJOUTER UN TOURNOI : ajoute une entrée dans TOURNAMENTS ci-dessous.
 *   L'ordre n'a pas d'importance, le calendrier trie tout seul. Les dates sont
 *   au format AAAA-MM-JJ. Pour un tournoi d'une seule journée, mets la même
 *   date en `start` et en `end`.
 *
 * ▶ POUR RELIER UN CARNET : une fois l'article publié dans /content/blog,
 *   renseigne `slugFr` (et `slugEn` s'il est traduit) avec le nom du fichier
 *   .md sans son extension. La barre du tournoi devient alors cliquable et
 *   mène directement au carnet.
 *
 * Ce fichier est lu au build : les tournois se retrouvent dans le HTML
 * pré-rendu, donc lisibles par Google sans exécuter de JavaScript.
 */
import type { CalEvent } from './calendar'

export const TOURNAMENTS: CalEvent[] = [
  {
    id: 'pontevedra-2026',
    start: '2026-07-25',
    end: '2026-07-30',
    name: 'Open de Pontevedra',
    nameEn: 'Pontevedra Open',
    location: 'Pontevedra, Espagne',
    locationEn: 'Pontevedra, Spain',
    format: '9 rondes',
    formatEn: '9 rounds',
    result: '4.5/9',
    slugFr: 'open-pontevedra-2026',
    slugEn: 'open-pontevedra-2026',
  },
  {
    id: 'badalona-2026',
    start: '2026-08-02',
    end: '2026-08-10',
    name: '50e Open de Badalona',
    nameEn: '50th Badalona Open',
    location: 'Badalona, Espagne',
    locationEn: 'Badalona, Spain',
    format: '9 rondes · 90 min + 30 s',
    formatEn: '9 rounds · 90 min + 30 s',
    result: '4.5/8',
    slugFr: 'open-badalona-2026',
    slugEn: 'open-badalona-2026',
  },
  {
    id: 'match-equipe-2026-08',
    start: '2026-08-22',
    end: '2026-08-23',
    name: 'Match par équipes',
    nameEn: 'Team match',
  },
  {
    id: 'blitz-geneve-2026',
    start: '2026-08-30',
    end: '2026-08-30',
    name: 'Tournoi de blitz de Genève',
    nameEn: 'Geneva blitz tournament',
    location: 'Genève, Suisse',
    locationEn: 'Geneva, Switzerland',
  },
  {
    id: 'grand-prix-monthey-2026',
    start: '2026-09-04',
    end: '2026-09-06',
    name: 'Grand Prix de Monthey',
    nameEn: 'Monthey Grand Prix',
    location: 'Monthey, Suisse',
    locationEn: 'Monthey, Switzerland',
  },
  {
    id: 'match-n1-2026-09',
    start: '2026-09-12',
    end: '2026-09-13',
    name: 'Match de Nationale 1',
    nameEn: 'National League 1 match',
  },
  {
    id: 'match-n1-2026-10',
    start: '2026-10-10',
    end: '2026-10-11',
    name: 'Match de Nationale 1',
    nameEn: 'National League 1 match',
  },
  {
    id: 'championnat-suisse-rapide-blitz-2026',
    start: '2026-09-26',
    end: '2026-09-27',
    name: 'Championnat suisse de blitz et de rapide',
    nameEn: 'Swiss blitz & rapid championship',
    location: 'Suisse',
    locationEn: 'Switzerland',
  },
]
