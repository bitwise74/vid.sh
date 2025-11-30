import type { User } from '$lib/api/User'
import { writable } from 'svelte/store'

type JobStats = {
    progress: number
    state?: string
}

export const isLoggedIn = writable(false)
export const loadedVideosCount = writable(0)

export const jobStats = writable<JobStats>({
    progress: 0,
    state: 'Uploading video to server...'
})

export const user = writable<User>({
    avatarHash: '',
    publicProfileEnabled: false,
    stats: {
        maxStorage: 0,
        usedStorage: 0,
        uploadedFiles: 0
    },
    username: 'Unknown',
    videos: [],
    defaultPrivateVideos: false
})
