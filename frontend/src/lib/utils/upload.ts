import { goto } from '$app/navigation'
import { PUBLIC_CDN_URL, PUBLIC_FILE_SIZE_LIMIT } from '$env/static/public'
import { UploadFile, type Video } from '$lib/api/Files'
import { user } from '$lib/stores/AppVars'
import { toastStore } from '$lib/stores/ToastStore'
import { videos } from '$lib/stores/VideoStore'

const ALLOWED_FORMATS = ['video/mp4', 'video/quicktime', 'video/x-matroska']

/**
 * Full implementation of the upload file button. Support multiple files and does everything
 * automatically like updating states and copying things.
 * @param richEmbed
 */
export async function UploadFileButton(richEmbed = false) {
    if (!window.location.href.endsWith('/dashboard')) {
        await goto('/dashboard')
    }

    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'video/*'
    input.multiple = true
    input.onchange = onChange
    input.click()
}

async function onChange(e: Event) {
    const target = e.target as HTMLInputElement
    if (!target.files || target.files.length === 0) return

    const MAX_FILE_SIZE = parseInt(PUBLIC_FILE_SIZE_LIMIT)
    let files = Array.from(target.files)

    if (files.some((f) => f.size > MAX_FILE_SIZE)) {
        toastStore.error({
            title: 'File too large',
            message: `One or more files exceed the size limit of ${Math.floor(MAX_FILE_SIZE / (1024 * 1024))} MB. Upload cancelled`,
            duration: 10000
        })
        return
    }

    if (files.some((f) => !ALLOWED_FORMATS.includes(f.type))) {
        toastStore.error({
            title: 'Invalid file format',
            message: 'One or more files are not in a supported format (mp4, mov, mkv). Upload cancelled.',
            duration: 10000
        })
        return
    }

    if (files.some((f) => f.type === 'video/x-matroska')) {
        toastStore.info({
            title: 'MKV support is experimental',
            message: 'Due to limited browser support the file will be re-muxed to the MP4 format. This may take a while depending on the file size.',
            duration: 15000
        })
    }

    const placeholders: Array<Video> = files.map((file, idx) => ({
        id: `${Date.now()}${idx}`,
        file_key: `${Date.now()}${idx}`,
        state: 'processing',
        name: file.name,
        private: false,
        format: file.type,
        size: file.size,
        version: 1,
        duration: 0,
        created_at: Date.now() / 1000
    }))

    for (const p of placeholders) videos.fPush([p])

    const queue = [] as any[]

    for (const file of files) {
        queue.push(() => {
            return new Promise<void>(async (resolve) => {
                const placeholder = placeholders.shift()!

                const video = await UploadFile(file).catch((err) => {
                    toastStore.error({
                        title: 'Failed to upload file',
                        message: `${err.message}`,
                        duration: 10000
                    })
                    return null
                })

                if (video) {
                    video.thumbnail_url = `${PUBLIC_CDN_URL}/${video.file_key.replace('.mp4', '.webp')}`
                    video.video_url = `${PUBLIC_CDN_URL}/${video.file_key}`

                    videos.delete('file_key', placeholder.file_key)
                    videos.fPush([video])

                    user.update((u) => {
                        u.stats.uploadedFiles += 1
                        u.stats.usedStorage += video.size
                        return u
                    })
                } else {
                    videos.delete('file_key', placeholder.file_key)
                }

                resolve()
            })
        })
    }

    await Promise.allSettled(queue.map((f) => f()))
}
