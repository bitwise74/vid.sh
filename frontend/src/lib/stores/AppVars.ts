import type { User } from '$lib/api-v2/User'
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

export const user = writable<User>()

// Data loading controller. If set to true anywhere, layout will refetch user data and set the data store
export const shouldRefetch = writable(true)
