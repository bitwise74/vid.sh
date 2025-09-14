import { PUBLIC_BASE_URL } from '$env/static/public'
import { jobProgress } from '../../routes/editor/Logic'
import type { UserStats } from './User'

export interface Video {
    id: number
    name: string
    format: string
    views: number
    size: number
    duration: number
    created_at: number
    file_key: string
    thumb_key: string
    state: string
}

export interface UpdateVideoOpts {
    name?: string
    processing_options?: ProcessingOpts
    save_to_cloud?: boolean
}

export interface ProcessingOpts {
    targetSize: number
    trimStart: number
    trimEnd: number
    losslessExport: boolean
    fps: number
    format: string
    saveToCloud: boolean
    crop: CropOpts
}

export interface ExportOpts {
    file: File
    processingOpts: ProcessingOpts
}

export interface CropOpts {
        x: number
        y: number
        w: number
        h: number
}

export async function LoadVideos(page: number, limit: number): Promise<Video[]> {
    if (page < 0) {
        throw new Error("Page can't be smaller than 0")
    }

    if (limit < 1 || limit > 250) {
        throw new Error('Limit must be between 1-250')
    }

    const resp = await fetch(`${PUBLIC_BASE_URL}/api/files/bulk?limit=${limit}&page=${page}`, {
        method: 'GET',
        credentials: 'include'
    })

    const body = await resp.json()
    if (resp.status === 200) return body

    console.error(`[${body.requestID}] LoadVideos request failed: ${body.error}`)
    throw new Error(body.error)
}

export async function LoadVideo(id: string | number): Promise<Video | undefined> {
    const resp = await fetch(`${PUBLIC_BASE_URL}/api/files/${id}`, {
        method: 'GET',
        credentials: 'include'
    })

    const body = await resp.json()

    if (resp.status === 200) return body

    console.error(`[${body.requestID}] LoadVideos request failed: ${body.error}`)
    throw new Error(body.error)
}

export async function SearchVideos(search: string, limit: number): Promise<Video[]> {
    if (search === '') {
        return LoadVideos(limit, 0)
    }

    const resp = await fetch(`${PUBLIC_BASE_URL}/api/files/search?query=${encodeURIComponent(search)}&limit=${limit}`, {
        credentials: 'include'
    })

    const body = await resp.json()

    if (resp.status === 200) return body

    console.error(`[${body.requestID}] LoadVideos request failed: ${body.error}`)
    throw new Error(body.error)
}

export async function DeleteVideo(videoID: number | string): Promise<UserStats> {
    const resp = await fetch(`${PUBLIC_BASE_URL}/api/files/${videoID}`, {
        method: 'DELETE',
        credentials: 'include'
    })

    const body = await resp.json()

    if (resp.status === 200) {
        return body
    }

    throw new Error(body.error)
}

export async function OwnsVideo(videoID: number | string): Promise<boolean> {
    const resp = await fetch(`${PUBLIC_BASE_URL}/api/files/${videoID}/owns`, {
        method: 'GET',
        credentials: 'include'
    })

    const body = await resp.json()
    if (resp.status === 200) return body.owns

    throw new Error(body.error)
}

export async function UpdateVideo(videoID: number | string, opts: UpdateVideoOpts): Promise<Video | undefined> {
    if (opts.processing_options) {
        // Get job ID
        const jobIDResp = await fetch(`${PUBLIC_BASE_URL}/api/ffmpeg/start`, {
            method: 'GET',
            credentials: 'include'
        })

        if (jobIDResp.status != 200) {
            const body = await jobIDResp.json()

            console.error('Failed to initialize ffmpeg job', body.requestID, body.error)
            throw new Error(body.error)
        }

        const jobID = (await jobIDResp.json()).jobID

        const source = new EventSource(`${PUBLIC_BASE_URL}/api/ffmpeg/progress?jobID=${jobID}`, {
            withCredentials: true
        })

        source.onmessage = function (event) {
            const v = parseFloat(event.data)

            if (v >= 100) {
                source.close()
            }

            jobProgress.set(v)
        }
    }

    const resp = await fetch(`${PUBLIC_BASE_URL}/api/files/${videoID}`, {
        method: 'PATCH',
        credentials: 'include',
        body: JSON.stringify(opts)
    })

    const body = await resp.json()

    if (resp.status === 200) return body

    console.error(`[${body.requestID}] UpdateVideo request failed: ${body.error}`)
    throw new Error(body.error)
}

export async function ExportVideo(o: ExportOpts, token: string) {
    const { trimStart, trimEnd, targetSize, losslessExport, saveToCloud } = o.processingOpts

    const form = new FormData()
    form.append('file', o.file)

    if (trimStart > 0) form.append('trimStart', `${trimStart}`)
    if (trimEnd > 0) form.append('trimEnd', `${trimEnd}`)
    if (targetSize > 0) form.append('targetSize', `${targetSize}`)
    if (saveToCloud) form.append('saveToCloud', `${saveToCloud}`)
    form.append('losslessExport', `${losslessExport}`)
    if (o.processingOpts.crop) form.append("crop", JSON.stringify(o.processingOpts.crop))

    // Get job ID
    const jobIDResp = await fetch(`${PUBLIC_BASE_URL}/api/ffmpeg/start`, {
        method: 'GET',
        credentials: 'include'
    })

    const body = await jobIDResp.json()
    if (jobIDResp.status != 200) {
        console.error('Failed to initialize ffmpeg job', body.requestID, body.error)
        throw new Error(body.error)
    }

    const jobID = body.jobID

    const source = new EventSource(`${PUBLIC_BASE_URL}/api/ffmpeg/progress`, {
        withCredentials: true
    })

    source.onmessage = function (event) {
        const v = parseFloat(event.data)

        if (v >= 100) {
            source.close()
        }

        jobProgress.set(v)
    }

    const resp = await fetch(`${PUBLIC_BASE_URL}/api/ffmpeg/process?jobID=${jobID}`, {
        method: 'POST',
        headers: {
            TurnstileToken: token
        },
        credentials: 'include',
        body: form
    })

    if (resp.status === 200) {
        return await resp.blob()
    }

    const bodyResp = await resp.json()
    throw new Error(bodyResp.error)
}

export async function UploadVideo(f: File): Promise<Video | undefined> {
    const form = new FormData()

    form.append('file', f)

    const resp = await fetch(`${PUBLIC_BASE_URL}/api/files`, {
        method: 'POST',
        credentials: 'include',
        body: form
    })

    const body = await resp.json()

    if (resp.status === 200) return body
    throw new Error(body.error)
}
