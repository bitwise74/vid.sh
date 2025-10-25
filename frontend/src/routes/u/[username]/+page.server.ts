import { PUBLIC_BASE_URL, PUBLIC_CDN_URL } from '$env/static/public'
import type { Video } from '$lib/api-v2/Files'
import type { PageLoad } from '../../settings/$types'

export type ProfileVideo = {
    file_key: string
    name: string
    duration: number
    created_at: number
    thumbnail_url: string
    video_url: string
}

export type ProfileData = {
    avatarHash: string
    public: boolean
    username: string
    videos: ProfileVideo[]
}

export const load: PageLoad = async ({ params, fetch, url }) => {
    const username = params['username']
    if (username === 'placeholder.svg') return {}
    const res = await fetch(`${PUBLIC_BASE_URL}/api/profile/${username}`)

    const body = await res.json()

    if (body.videos) {
        for (const vid of body.videos as Video[]) {
            vid.thumbnail_url = `${PUBLIC_CDN_URL}/${vid.file_key.replace('.mp4', '.webp')}`
            vid.video_url = `${PUBLIC_CDN_URL}/${vid.file_key}`
        }
    }

    return body
}
