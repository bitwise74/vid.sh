import { PUBLIC_BASE_URL } from '$env/static/public'
import { FFmpegMonitorProgress, StartFFmpegJob } from './FFmpeg'
import type { UserStats } from './User'

export type Video = {
    id: string
    file_key: string
    state: string
    name: string
    private: boolean
    format: string
    size: number
    version: number
    duration: number
    created_at: number
    expires_at?: number

    // Variables not send by the server
    thumbnail_url?: string
    video_url?: string
}

export type BulkFetchOpts = {
    page: number
    limit: number
    sort: 'newest' | 'oldest' | 'az' | 'za' | 'size-asc' | 'size-desc'
    tags: string // TODO
}

export type VideoProcessingOpts = {
    trimStart: number
    trimEnd: number
    targetSize: number
    losslessExport: boolean
    cropX: number
    cropY: number
    cropW: number
    cropH: number
}

export type VideoUpdateOpts = {
    processing_options?: VideoProcessingOpts
    name?: string
    private?: boolean
}

export type SearchOpts = {
    query: string
    page: number
    limit: number
}

export type CropOpts = {
    x: number
    y: number
    w: number
    h: number
}

/**
 * Checks if the logged in user owns a file
 * @param id ID of the file
 * @returns ownership status
 */
export async function CheckFileOwnership(id: string): Promise<boolean> {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/files/${id}/owns`, { credentials: 'include' })
    const body = await req.json()

    if (!req.ok) {
        console.error(`[Files/CheckFileOwnership]: Request failed, requestID: ${body.requestID || 'Unknown'}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return body.owns
}

/**
 * Fetches a single file from a server by id
 * @param id File ID to fetch
 * @returns Video details
 */
export async function FetchFile(id: string): Promise<Video> {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/files/${id}`, {
        credentials: 'include',
        method: 'GET'
    })
    const body = await req.json()

    if (!req.ok) {
        console.error(`[Files/FetchFile]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    // Discard the user
    return body.file
}

/**
 * Fetches multiple files of a user
 * @param o Fetch options
 */
export async function FetchFiles(o: BulkFetchOpts): Promise<Array<Video>> {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/files/bulk`, {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify(o)
    })
    const body = await req.json()

    if (!req.ok) {
        console.error(`[Files/FetchFiles]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    // Discard the user
    return body
}

/**
 * Uploads a file without any editing options
 * @param f File to be uploaded
 */
export async function UploadFile(f: File): Promise<Video> {
    const form = new FormData()

    form.append('file', f)

    const req = await fetch(`${PUBLIC_BASE_URL}/api/files`, {
        credentials: 'include',
        method: 'POST',
        body: form
    })
    const body = await req.json()

    if (!req.ok) {
        console.error(`[Files/UploadFile]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    // Discard the user
    return body
}

/**
 * Updates a file with processing options
 * @param id File ID to update
 * @param o Editing options
 */
export async function UpdateFile(id: string, o: VideoUpdateOpts): Promise<Video> {
    console.log(o)

    const req = await fetch(`${PUBLIC_BASE_URL}/api/files/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(o),
        credentials: 'include'
    })

    const body = await req.json()

    if (!req.ok) {
        console.error(`[Files/UpdateFile]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return body
}

/**
 * Deletes multiple files from the server
 * @param ids IDS of files to remove
 */
export async function DeleteFiles(ids: Array<string>): Promise<UserStats> {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/files`, {
        credentials: 'include',
        method: 'DELETE',
        body: JSON.stringify({
            ids: ids
        })
    })
    const body = await req.json()

    if (!req.ok) {
        console.error(`[Files/DeleteFile]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return body
}

/**
 * Searches for videos matching a search query
 * @param q
 */
export async function SearchFiles(o: SearchOpts): Promise<Array<Video>> {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/files/search`, {
        credentials: 'include',
        method: 'DELETE',
        body: JSON.stringify(o)
    })
    const body = await req.json()

    if (!req.ok) {
        console.error(`[Files/SearchFiles]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return body
}

/**
 * Uploads a video for processing and downloads or saves it to the cloud
 * @param f File to process
 * @param o Processing options
 * @param token Turnstile token
 * @param saveToCloud Should the video be saved to cloud?
 */
export async function ExportVideo(f: File, o: VideoProcessingOpts, token: string, saveToCloud: false): Promise<Blob>
export async function ExportVideo(f: File, o: VideoProcessingOpts, token: string, saveToCloud: true): Promise<Video>
export async function ExportVideo(f: File, o: VideoProcessingOpts, token: string, saveToCloud = false): Promise<Blob | Video> {
    const form = new FormData()

    form.append('file', f)
    form.append('trimStart', `${o.trimStart}`)
    form.append('trimEnd', `${o.trimEnd}`)
    form.append('targetSize', `${o.targetSize}`)
    form.append('losslessExport', `${o.losslessExport}`)
    form.append('saveToCloud', `${saveToCloud}`)
    form.append('crop[x]', `${o.cropX}`)
    form.append('crop[y]', `${o.cropY}`)
    form.append('crop[w]', `${o.cropW}`)
    form.append('crop[h]', `${o.cropH}`)

    const jobID = await StartFFmpegJob()
    FFmpegMonitorProgress(jobID)
    const req = await fetch(`${PUBLIC_BASE_URL}/api/ffmpeg/process?jobID=${jobID}`, {
        method: 'POST',
        credentials: 'include',
        body: form,
        headers: {
            TurnstileToken: token
        }
    })

    if (!req.ok) {
        const body = await req.json()
        console.error(`[Files/ExportVideo]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return saveToCloud ? await req.json() : await req.blob()
}
