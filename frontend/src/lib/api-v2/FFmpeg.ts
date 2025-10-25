import { PUBLIC_BASE_URL } from '$env/static/public'
import { jobStats } from '$lib/stores/AppVars'

/**
 * Starts a new FFmpeg job
 * @returns Job ID of the started job
 */
export async function StartFFmpegJob(): Promise<string> {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/ffmpeg/start`, {
        credentials: 'include'
    })

    const body = await req.json()

    if (!req.ok) {
        console.error(`[FFmpeg/StartFFmpegJob]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return body.jobID
}

/**
 * Monitors the progress of an FFmpeg job writing it's progress to the jobStats store
 * @param jobID ID of the job to monitor
 */
export async function FFmpegMonitorProgress(jobID: string) {
    const source = new EventSource(`${PUBLIC_BASE_URL}/api/ffmpeg/progress?jobID=${jobID}`, {
        withCredentials: true
    })

    source.onmessage = (e) => {
        const v = e.data.split('|')
        let [progress, state] = [v[0], v[1]]

        if (progress <= 0) return

        jobStats.set({
            progress,
            state
        })
    }

    source.onerror = (e) => {
        source.close()
    }
}
