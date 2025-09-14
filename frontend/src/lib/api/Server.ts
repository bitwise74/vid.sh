import { PUBLIC_BASE_URL } from '$env/static/public'

// Checks if the server is alive
export async function Heartbeat() {
    const resp = await fetch(`${PUBLIC_BASE_URL}/api/heartbeat`, {
        method: 'HEAD'
    })

    if (resp.status !== 200) throw new Error('Backend server offline')
}
