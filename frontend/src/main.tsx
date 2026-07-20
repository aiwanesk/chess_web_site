import './styles.css'
import { ViteReactSSG } from 'vite-react-ssg'
import { routes } from './routes'

// vite-react-ssg pre-renders every route to static HTML at build, then hydrates
// on the client. createRoot is the entry it drives.
export const createRoot = ViteReactSSG({ routes })
