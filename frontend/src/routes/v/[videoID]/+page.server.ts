import { PUBLIC_BASE_URL, PUBLIC_CDN_URL } from '$env/static/public'
import type { PageServerLoad } from './$types'

export const load: PageServerLoad = async ({ params, fetch }) => {
    const videoResp = await fetch(`${PUBLIC_BASE_URL}/api/files/${params.videoID}`);

    if (videoResp.status !== 200) {
        return {
            og: {
                title: "No such file",
                description: "Either the file doesn't exist, was deleted, or is private",
                url: "https://bitwise0x.dev"
            }
        };
    }

    const body = await videoResp.json();

    return {
        og: {
            title: body.file?.name ?? "No such file",
            description: body.file ? `Video by @${body.user.username}` : "Either the file doesn't exist, was deleted, or is private",
            url: `${PUBLIC_BASE_URL}/v/${params.videoID}`,
            video_url: body.file ? `${PUBLIC_CDN_URL}/${body.file.file_key}` : undefined,
            site_name: body.file ? `Video by @${body.user.username}` : undefined
        },
        video: body
    };
};
