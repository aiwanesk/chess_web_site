import type { RouteRecord } from 'vite-react-ssg'
import { Layout } from './components/Layout'
import { blogStaticPaths, blogStaticPathsEN } from './lib/blogSlugs'
import { categoryStaticPaths, categoryStaticPathsEN } from './lib/categories'
import { weekStaticPaths } from './lib/tactics'

/**
 * Route table — one entry per indexable URL. Pages are lazy-loaded so each
 * route ships its own JS chunk (money pages don't pay for blog/Markdown code).
 * Keep this aligned with backend/internal/content/site.go (sitemap/llms.txt).
 */
export const routes: RouteRecord[] = [
  {
    path: '/',
    element: <Layout />,
    children: [
      { index: true, lazy: () => import('./pages/Home') },

      // Money pages
      { path: 'cours-echecs-adultes-geneve', lazy: () => import('./pages/CoursAdultesGeneve') },
      { path: 'preparation-tournoi-echecs', lazy: () => import('./pages/PreparationTournoi') },
      { path: 'cours-echecs-en-ligne', lazy: () => import('./pages/CoursEnLigne') },
      { path: 'cours-echecs-groupe-geneve', lazy: () => import('./pages/CoursGroupe') },
      { path: 'cours-echecs-ados-competition', lazy: () => import('./pages/CoursAdos') },
      { path: 'stages-echecs-geneve', lazy: () => import('./pages/Stages') },
      { path: 'conferences-echecs-entreprise', lazy: () => import('./pages/Conferences') },
      { path: 'team-building-echecs-geneve', lazy: () => import('./pages/TeamBuilding') },

      // English (EN slugs). Same components, locale detected from the URL.
      { path: 'en', lazy: () => import('./pages/Home') },
      { path: 'en/adult-chess-lessons-geneva', lazy: () => import('./pages/CoursAdultesGeneve') },
      { path: 'en/tournament-preparation', lazy: () => import('./pages/PreparationTournoi') },
      { path: 'en/online-chess-lessons', lazy: () => import('./pages/CoursEnLigne') },
      { path: 'en/group-chess-lessons-geneva', lazy: () => import('./pages/CoursGroupe') },
      { path: 'en/junior-chess-coaching', lazy: () => import('./pages/CoursAdos') },
      { path: 'en/chess-camps-geneva', lazy: () => import('./pages/Stages') },
      { path: 'en/corporate-chess-talks', lazy: () => import('./pages/Conferences') },
      { path: 'en/chess-team-building-geneva', lazy: () => import('./pages/TeamBuilding') },
      { path: 'en/about', lazy: () => import('./pages/APropos') },
      { path: 'en/results', lazy: () => import('./pages/Resultats') },
      { path: 'en/pricing', lazy: () => import('./pages/Tarifs') },
      { path: 'en/tactics', lazy: () => import('./pages/Tactiques') },
      { path: 'en/contact', lazy: () => import('./pages/Contact') },

      // Info pages
      { path: 'a-propos', lazy: () => import('./pages/APropos') },
      { path: 'resultats', lazy: () => import('./pages/Resultats') },
      { path: 'tarifs', lazy: () => import('./pages/Tarifs') },
      { path: 'contact', lazy: () => import('./pages/Contact') },
      { path: 'confidentialite', lazy: () => import('./pages/Confidentialite') },
      { path: 'en/privacy', lazy: () => import('./pages/Confidentialite') },

      // Tactiques de la semaine — index + article hebdo (/tactiques/JJ-MM-AA)
      { path: 'tactiques', lazy: () => import('./pages/Tactiques') },
      { path: 'tactiques/:date', lazy: () => import('./pages/TactiquesWeek'), getStaticPaths: weekStaticPaths },
      {
        path: 'en/tactics/:date',
        lazy: () => import('./pages/TactiquesWeek'),
        getStaticPaths: () => weekStaticPaths().map((p) => '/en/tactics/' + p.split('/').pop()),
      },

      // Blog
      { path: 'blog', lazy: () => import('./pages/Blog') },
      {
        path: 'blog/categorie/:cat',
        lazy: () => import('./pages/BlogCategory'),
        getStaticPaths: categoryStaticPaths,
      },
      {
        path: 'blog/:slug',
        lazy: () => import('./pages/BlogPost'),
        getStaticPaths: blogStaticPaths,
      },

      // Blog (EN)
      { path: 'en/blog', lazy: () => import('./pages/Blog') },
      {
        path: 'en/blog/category/:cat',
        lazy: () => import('./pages/BlogCategory'),
        getStaticPaths: categoryStaticPathsEN,
      },
      {
        path: 'en/blog/:slug',
        lazy: () => import('./pages/BlogPost'),
        getStaticPaths: blogStaticPathsEN,
      },

      // Concrete /404 so the SSG build emits dist/404.html (served by the Go
      // backend on any unmatched path). The splat handles client-side misses.
      { path: '404', lazy: () => import('./pages/NotFound') },
      { path: '*', lazy: () => import('./pages/NotFound') },
    ],
  },
]
