import { PUBLIC_BASE_URL, PUBLIC_CDN_URL } from '$env/static/public'
import type { Video } from '$lib/api/Files'
import type { PageServerLoad } from './$types'

export type ProfileVideo = {
    file_key: string
    name: string
    duration: number
    created_at: number
    thumbnail_url: string
    video_url: string
    size: number
}

export type ProfileData = {
    avatarHash: string
    username: string
    videos: ProfileVideo[]
    found: boolean
}

export type ProfileResponse = {
        found: boolean
}

export const load: PageServerLoad = async ({ params, fetch }) => {
    const username = params['username']
    if (username === 'placeholder.svg') return { found: false };

    const res = await fetch(`${PUBLIC_BASE_URL}/api/profile/${username}`)
    if (res.status === 404) {
        return { found: false }
    }

    console.log(res)

    const body = await res.json().catch(err => {
        console.error(err)
    })
    body.found = true

    if (body.videos) {
        for (const vid of body.videos as Video[]) {
            vid.thumbnail_url = `${PUBLIC_CDN_URL}/${vid.file_key.replace('.mp4', '.webp')}`
            vid.video_url = `${PUBLIC_CDN_URL}/${vid.file_key}`
        }
    }

    return body
}
