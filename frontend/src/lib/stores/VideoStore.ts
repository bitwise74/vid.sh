import type { Video } from '$lib/api/Files'
import { writable } from 'svelte/store'

function extendedWritable<T>() {
    const { subscribe, update } = writable([] as Array<T>)

    return {
        subscribe,
        update,
        // Pushes item(s) to the front of the array
        fPush: (items: Array<T>) => {
            update((arr) => {
                arr.unshift(...items)
                return arr
            })
        },
        // Pushes item(s) to the back
        push: (items: Array<T>) => {
            update((arr) => {
                arr.push(...items)
                return arr
            })
        },
        // Replaces all items by new ones
        set: (items: Array<T>) => {
            update((arr) => {
                arr = items
                return arr
            })
        },
        // Deletes an item by it's key value pair
        delete: (key: string, val: string) => {
            update((arr) => {
                const idx = arr.findIndex((x) => x[key] === val)

                if (idx != -1) {
                    arr.splice(idx, 1)
                }

                return arr
            })
        },
        // Replaces an item by file_key
        replace: (key: string, item: T) => {
            update((arr) => {
                const idx = arr.findIndex((x) => x['file_key'] === key)

                // Hacks svelte's reactivity system since a simple splice
                // wouldn't update the store and would leave stale data
                if (idx != -1) {
                    arr[idx] = { ...item }
                }
                return arr
            })
        }
    }
}

export const currentTime = writable(0)

export const currentVideoURL = writable('')
export const videos = extendedWritable<Video>()
