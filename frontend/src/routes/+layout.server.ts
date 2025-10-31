import { PUBLIC_CDN_URL } from '$env/static/public'
import { GetUser } from '$lib/api/User'
import { toastStore } from '$lib/stores/ToastStore'
import type { LayoutServerLoad } from './$types'

export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
    if (cookies.get('logged_in') !== '1') return null

    const data = await GetUser(fetch).catch((err) => {
        throw toastStore.error({
            title: 'Failed to fetch user data',
            message: err.message
        })
    })

    if (!data) return null

    const vids = data.videos || []

    for (let i = 0; i < vids.length; i++) {
        const v = vids[i]

        vids[i].thumbnail_url = `${PUBLIC_CDN_URL}/${v.file_key.split('.')[0]}.webp`
        vids[i].video_url = `${PUBLIC_CDN_URL}/${v.file_key}${v.version > 1 ? `?v=${v.version}` : ''}`
    }

    data.videos = vids
    return data
}
