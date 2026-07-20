import type { RouteRecord } from 'vite-react-ssg'
import { Layout } from './components/Layout'
import { blogStaticPaths } from './lib/blogSlugs'
import { categoryStaticPaths } from './lib/categories'

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
      { path: 'en/pricing', lazy: () => import('./pages/Tarifs') },
      { path: 'en/contact', lazy: () => import('./pages/Contact') },

      // Info pages
      { path: 'a-propos', lazy: () => import('./pages/APropos') },
      { path: 'resultats', lazy: () => import('./pages/Resultats') },
      { path: 'tarifs', lazy: () => import('./pages/Tarifs') },
      { path: 'contact', lazy: () => import('./pages/Contact') },

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

      // Concrete /404 so the SSG build emits dist/404.html (served by the Go
      // backend on any unmatched path). The splat handles client-side misses.
      { path: '404', lazy: () => import('./pages/NotFound') },
      { path: '*', lazy: () => import('./pages/NotFound') },
    ],
  },
]
