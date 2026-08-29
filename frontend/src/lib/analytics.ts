/**
 * Balise de fréquentation.
 *
 * Le comptage se faisait au serveur, sur chaque réponse HTML : tout crawler qui
 * se présentait avec un user-agent de navigateur était donc compté comme un
 * visiteur humain. D'où un tableau de bord dominé par des datacenters
 * américains et asiatiques, très loin du public réel d'un coach genevois.
 *
 * Les robots n'exécutent quasiment jamais de JavaScript : déclencher le
 * comptage depuis la page filtre l'essentiel du bruit sans rien changer aux
 * garanties de confidentialité — pas de cookie, pas de tiers, aucune IP
 * stockée, pays déduit hors ligne côté serveur.
 */
const ENDPOINT = '/api/hit'

/** Signale une page vue. Silencieux en cas d'échec : une statistique perdue ne
 * doit jamais peser sur la navigation. */
export function sendHit(path: string): void {
  if (typeof window === 'undefined') return
  const body = JSON.stringify({ path })
  try {
    // sendBeacon survit à la fermeture de l'onglet et ne retarde jamais le rendu.
    if (typeof navigator.sendBeacon === 'function') {
      navigator.sendBeacon(ENDPOINT, new Blob([body], { type: 'application/json' }))
      return
    }
    void fetch(ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body,
      keepalive: true,
    }).catch(() => {})
  } catch {
    /* navigation privée, extension de blocage, hors ligne : on laisse passer */
  }
}
