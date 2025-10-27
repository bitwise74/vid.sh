import { DeleteFiles } from '$lib/api/Files'
import { selectedVideos } from '$lib/stores/appControl'
import { user } from '$lib/stores/AppVars'
import { toastStore } from '$lib/stores/ToastStore'
import { videos } from '$lib/stores/VideoStore'
import { derived } from 'svelte/store'

export const videoSelectionToast = derived(selectedVideos, ($selectedVideos) => {
    if ($selectedVideos.length <= 0) {
        toastStore.remove('video-selection-notification')
        return null
    }

    const config = {
        id: 'video-selection-notification',
        title: `${$selectedVideos.length} video${$selectedVideos.length > 1 ? 's' : ''} selected`,
        dismissible: false,
        buttons: [
            {
                text: 'Clear',
                action: () => selectedVideos.set([]),
                class: 'btn btn-outline-primary'
            },
            {
                text: 'Delete',
                action: async () => {
                    try {
                        const newStats = await DeleteFiles($selectedVideos)
                        user.update((u) => {
                            if (!u) return u
                            return { ...u, stats: newStats }
                        })
                        toastStore.remove('video-selection-notif')
                        toastStore.success({
                            title: 'Videos deleted',
                            message: `${$selectedVideos.length} video${$selectedVideos.length > 1 ? 's' : ''} deleted successfully`
                        })
                        for (const id of $selectedVideos) {
                            videos.delete(videos.get(id)!)
                        }
                        selectedVideos.set([])
                    } catch (error) {
                        console.error('Error deleting videos:', error)
                        toastStore.remove('video-selection-notif')
                        toastStore.error({
                            title: 'Error',
                            message: error.message
                        })
                    }
                },
                class: 'btn btn-outline-danger'
            }
        ]
    }

    if (toastStore.exists('video-selection-notification')) {
        toastStore.update('video-selection-notification', config)
        return null
    } else {
        toastStore.question(config)
    }
})
