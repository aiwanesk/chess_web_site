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

/**
 * Hôte du référent, et UNIQUEMENT pour la première vue du chargement.
 *
 * Le serveur ne peut pas lire cette information : `/api/hit` est un appel
 * same-origin, donc son en-tête `Referer` désigne toujours le site lui-même —
 * c'est même ce que le serveur vérifie pour écarter les appels directs. L'hôte
 * externe n'existe que dans `document.referrer`, côté navigateur.
 *
 * Il faut le remonter une seule fois : `document.referrer` ne change pas
 * pendant la vie du document, donc l'envoyer à chaque changement de route
 * compterait cinq provenances pour un visiteur qui lit cinq pages.
 *
 * On n'envoie que l'hôte, jamais l'URL : savoir qu'une visite vient de
 * chatgpt.com est utile, savoir quelle conversation l'a produite ne l'est pas.
 */
let refSent = false

function referrerHost(): string | undefined {
  if (refSent) return undefined
  refSent = true
  try {
    const raw = document.referrer
    if (!raw) return undefined
    const host = new URL(raw).hostname.toLowerCase()
    return host && host !== window.location.hostname ? host : undefined
  } catch {
    return undefined
  }
}

/** Signale une page vue. Silencieux en cas d'échec : une statistique perdue ne
 * doit jamais peser sur la navigation. */
export function sendHit(path: string): void {
  if (typeof window === 'undefined') return
  const ref = referrerHost()
  const body = JSON.stringify(ref ? { path, ref } : { path })
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
