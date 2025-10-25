import { writable } from 'svelte/store'

export const selectedVideos = writable<Array<string>>([])
