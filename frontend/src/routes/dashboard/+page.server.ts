import { PUBLIC_BASE_URL } from '$env/static/public'
import { redirect, type ServerLoad } from '@sveltejs/kit'

export const load: ServerLoad = async ({ cookies, fetch }) => {
    const loggedIn = cookies.get('logged_in')
    const userID = cookies.get("user_id")
    if (loggedIn !== '1') {
            throw redirect(302, '/login')
    }

    // Server side things need to use the built in fetch
    const resp = await fetch(`${PUBLIC_BASE_URL}/api/validate`, {
            method: "GET",
            credentials: 'include'
    })

    if (resp.status === 200) return {}
    const body = await resp.json()
    switch (body.error) {
            case "account_not_verified": {
                    throw redirect(302, '/verify?userID=' + userID)
            }
            default: {
                await fetch(`${PUBLIC_BASE_URL}/api/users/logout`, { method: "POST" })
                throw redirect(302, '/login')
            }
    }
}
