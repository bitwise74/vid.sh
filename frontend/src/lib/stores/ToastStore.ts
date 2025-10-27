import { writable } from 'svelte/store'

type ToastButton = {
    text: string
    action: () => void
    class?: string
}

export type ToastTypes = 'error' | 'success' | 'info' | 'warning' | 'loading' | 'question'

export interface Toast {
    id: string
    type: ToastTypes
    title: string
    message: string
    duration: number
    dismissible: boolean
    buttons: Array<ToastButton>
}

interface ToastOpts {
    id?: string
    title: string
    message?: string
    duration?: number
    dismissible?: boolean
    buttons?: Array<ToastButton>
}

function createToastStore() {
    const { subscribe, update } = writable<Toast[]>([])

    return {
        subscribe,
        add: (toast: ToastOpts, type: ToastTypes) => {
            if (!toast.id) {
                toast.id = Math.random().toString(36).substr(2, 9)
            }

            const newToast: Toast = {
                id: toast.id,
                duration: toast.duration || 5000,
                dismissible: type === 'loading' ? false : (toast.dismissible ?? true),
                buttons: toast.buttons || [],
                message: toast.message || '',
                title: toast.title,
                type: type
            }

            update((toasts) => [...toasts, newToast])
            return toast.id
        },
        remove: (id: string) => {
            update((t) => t.filter((t) => t.id !== id))
        },
        clear: () => {
            update(() => [])
        },
        success: (o: ToastOpts) => {
            return toastStore.add(o, 'success')
        },
        error: (o: ToastOpts) => {
            return toastStore.add(o, 'error')
        },
        warning: (o: ToastOpts) => {
            return toastStore.add(o, 'warning')
        },
        info: (o: ToastOpts) => {
            return toastStore.add(o, 'info')
        },
        loading: (o: ToastOpts) => {
            return toastStore.add(o, 'loading')
        },
        question: (o: ToastOpts) => {
            toastStore.add(o, 'question')
        },
        // Checks if a toast with the given ID exists
        exists: (id: string) => {
            let exists = false
            update((toasts) => {
                exists = toasts.some((t) => t.id === id)
                return toasts
            })
            return exists
        },
        update: (id: string, o: Partial<ToastOpts> & { type?: ToastTypes }) => {
            update((toasts) =>
                toasts.map((t) => {
                    if (t.id === id) {
                        return {
                            ...t,
                            ...o,
                            type: o.type || t.type
                        }
                    }
                    return t
                })
            )
        }
    }
}

export const toastStore = createToastStore()
