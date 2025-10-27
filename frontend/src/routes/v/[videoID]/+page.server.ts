import { PUBLIC_BASE_URL } from '$env/static/public'
import type { PageServerLoad } from './$types'

export const load: PageServerLoad = async ({ params, fetch }) => {
    const videoResp = await fetch(`${PUBLIC_BASE_URL}/api/files/${params.videoID}`)

    if (videoResp.status !== 200) {
        return { video: undefined }
    }

    const body = await videoResp.json()
    return {
        video: body
    }
}
