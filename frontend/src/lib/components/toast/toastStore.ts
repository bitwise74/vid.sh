import { writable } from "svelte/store"

export interface Toast {
  id: string
  type: "success" | "error" | "warning" | "info"
  title: string
  message: string
  duration: number
  dismissible: boolean
}

type ToastType = "success" | "error" | "warning" | "info"

interface ToastOpts {
        title: string
        message?: string
        duration?: number
        dismissible?: boolean
}

function createToastStore() {
  const { subscribe, update } = writable<Toast[]>([])

  return {
    subscribe,
    add: (toast: ToastOpts, type: ToastType) => {
      const id = Math.random().toString(36).substr(2, 9)
      const newToast: Toast = {
        id,
        duration: toast.duration || 5000,
        dismissible: toast.dismissible === undefined ? true : toast.dismissible,
        message: toast.message || '',
        title: toast.title,
        type: type
      }

      update((toasts) => [...toasts, newToast])

      if (newToast.duration && newToast.duration > 0) {
        setTimeout(() => {
          toastStore.remove(id)
        }, newToast.duration)
      }

      return id
    },
    remove: (id: string) => {
      update(t => t.filter(t => t.id !== id))
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
  }
}

export const toastStore = createToastStore()
