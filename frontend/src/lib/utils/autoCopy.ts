import { PUBLIC_BASE_URL, PUBLIC_CDN_URL } from '$env/static/public'
import { toastStore } from '$lib/stores/ToastStore'

// If optAutoCopy is enabled automatically copies the URL to the saved video
// to the user's clipboard. Also works with the optRichEmbeds option
export async function AutoCopy(file_key: string) {
    if (localStorage.getItem('optAutoCopy') !== 'true') return

    let url = ''

    if (localStorage.getItem('optRichEmbeds') === 'true') {
        url = `${PUBLIC_BASE_URL}/v/${file_key}`
    } else {
        url = `${PUBLIC_CDN_URL}/${file_key}`
    }

    await navigator.clipboard.writeText(url)
        .catch(() => {
            if ("Notification" in window && Notification.permission === 'granted') {
                new Notification('Failed to copy link to clipboard', {
                    body: 'This could be because you switched tabs or because you denied clipboard access',
                    icon: '/favicon.svg'
                })
            }
        })

        toastStore.info({
        title: 'Link copied to clipboard',
        message: 'Managed by auto copy URLs',
        desktopNotification: true
    })
}
