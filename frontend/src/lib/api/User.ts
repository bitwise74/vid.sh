import { PUBLIC_BASE_URL } from '$env/static/public'
import type { Video } from './Files'

export type UserStats = {
    maxStorage: number
    usedStorage: number
    uploadedFiles: number
    totalViews: number
    totalWatchTime: number
}

export interface InitialData {
    videos: Video[]
    stats: UserStats
}

export async function LoadInitialData(): Promise<InitialData> {
    const resp = await fetch(`${PUBLIC_BASE_URL}/api/users`, {
        method: 'GET',
        credentials: 'include'
    })

    const body = await resp.json()

    // Only possible codes are 200 and 500
    if (resp.status !== 200) throw new Error(body.error)
    return body
}

// TODO: finish
export async function VerifyUser(userID: string, token: string) {
    const resp = await fetch(`${PUBLIC_BASE_URL}/api/users/${userID}/verify`, {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({ token: token })
    })
}
