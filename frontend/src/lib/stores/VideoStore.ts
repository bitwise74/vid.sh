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
        // Returns an item by it's file_key
        get: (file_key: string): T | null => {
            let found: T | null = null

            update((arr) => {
                const idx = arr.findIndex((x) => x['file_key'] === file_key)
                if (idx != -1) {
                    found = arr[idx]
                }
                return arr
            })

            return found
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
                if (idx != -1) { arr[idx] = { ...item } }
                return arr
            })
        }
    }
}

export const currentTime = writable(0)

export const currentVideoURL = writable('')
export const videos = extendedWritable<Video>()
