import { derived, writable } from "svelte/store";
import { loadedVideosCount, user } from "./AppVars";

export const shouldRefetch = writable(true)
export const selectedVideos = writable<Array<string>>([])

export const dashboardView = writable<'grid' | 'list'>('grid')

// Filtering options
export const perPage = writable('20')