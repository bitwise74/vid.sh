import { PUBLIC_BASE_URL, PUBLIC_CDN_URL } from '$env/static/public'
import { toastStore } from '$lib/stores/ToastStore'

// If optAutoCopy is enabled automatically copies the URL to the saved video
// to the user's clipboard. Also works with the optRichEmbeds option
export async function AutoCopy(id: string, file_key: string) {
    if (localStorage.getItem('optAutoCopy') !== 'true') return

    let url = ''

    if (localStorage.getItem('optRichEmbeds') === 'true') {
        url = `${PUBLIC_BASE_URL}/v/${id}`
    } else {
        url = `${PUBLIC_CDN_URL}/${file_key}`
    }

    await navigator.clipboard.writeText(url)
    toastStore.info({
        title: 'Link copied to clipboard',
        message: 'Managed by auto copy URLs'
    })
}
