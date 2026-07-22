/**
 * Anti-spam form token (self-made, no third party). Every form POST carries a
 * short-lived HMAC token obtained from GET /api/form-token, so a bot POSTing
 * straight to the API — with no token — is rejected. The token is cached and
 * transparently refreshed once if the server reports it stale (403).
 */
let cached = ''
let inflight: Promise<string> | null = null

async function fetchToken(): Promise<string> {
  if (typeof fetch === 'undefined') return ''
  try {
    const r = await fetch('/api/form-token', { headers: { Accept: 'application/json' } })
    if (!r.ok) return ''
    const b = (await r.json()) as { token?: string }
    return b.token || ''
  } catch {
    return ''
  }
}

/** Returns a form token, fetching (and caching) it if needed. */
export async function getFormToken(force = false): Promise<string> {
  if (cached && !force) return cached
  if (!inflight) inflight = fetchToken().then((t) => ((cached = t), (inflight = null), t))
  return inflight
}

/** POST a JSON payload with the anti-spam token, retrying once on a stale token. */
export async function postWithToken(url: string, payload: Record<string, unknown>): Promise<Response> {
  const send = (token: string) =>
    fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...payload, token }),
    })

  let res = await send(await getFormToken())
  if (res.status === 403) res = await send(await getFormToken(true)) // refresh once
  return res
}
