import { writable } from 'svelte/store'

export const selectedVideos = writable<Array<string>>([])

export const dashboardView = writable<string>('grid')

// Filtering options
export const perPage = writable('20')
