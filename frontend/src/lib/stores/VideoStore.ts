import type { Video } from '$lib/api/Files'
import { writable } from 'svelte/store'

function uniqueArrayStore<T>(key = 'id') {
    const { subscribe, update } = writable([] as any[])

    return {
        subscribe,
        upsert: (item: T) =>
            update((arr) => {
                const index = arr.findIndex((x) => x[key] === item[key])
                if (index === -1) {
                    return [...arr, item]
                } else {
                    const copy = [...arr]
                    copy[index] = item
                    return copy
                }
            }),
        delete: (item: T) =>
            update((arr) => {
                return arr.filter((x) => x[key] !== item[key])
            }),
        set: (items: Array<T>) =>
            update(() => {
                const uniqueItems = [] as any[]
                const seenKeys = new Set()
                for (const item of items) {
                    if (!seenKeys.has(item[key])) {
                        seenKeys.add(item[key])
                        uniqueItems.push(item)
                    }
                }
                return uniqueItems
            }),
        addFront: (item: T) =>
            update((arr) => {
                const index = arr.findIndex((x) => x[key] === item[key])
                if (index === -1) {
                    return [item, ...arr]
                } else {
                    return arr
                }
            }),
        get: (id: string): T | null => {
            let foundItem = null
            update((arr) => {
                const index = arr.findIndex((x) => x[key] === id)
                if (index !== -1) {
                    foundItem = arr[index]
                }
                return arr
            })
            return foundItem
        },
        clear: () => update(() => []),
        update: (fn: (arr: Array<T>) => Array<T>) => update(fn),
        idx: (keyValue: any) => {
            let foundIndex = -1
            update((arr) => {
                const index = arr.findIndex((x) => x[key] === keyValue)
                foundIndex = index
                return arr
            })
            return foundIndex
        },
        replace: (idx: number, item: T) =>
            update((arr) => {
                const copy = [...arr]
                copy[idx] = item
                return copy
            })
    }
}

export const currentTime = writable(0)

export const currentVideoURL = writable('')
export const videos = uniqueArrayStore<Video>('id')
