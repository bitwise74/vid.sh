import { PUBLIC_BASE_URL } from '$env/static/public'
import type { Video } from './Files'

export type UserStats = {
    uploadedFiles: number
    maxStorage: number
    usedStorage: number
}

export type User = {
    avatarHash: string
    username: string
    stats: UserStats
    videos: Array<Video>
    publicProfileEnabled: boolean
}

export type UserUpdate = {
    avatarHash?: string
    avatar?: File | null
    username?: string
    publicProfileEnabled?: boolean
    email?: string
}
/**
 * Fetches the current user
 * @returns User object
 */
export async function GetUser(fetch: typeof window.fetch): Promise<User> {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/users`, {
        credentials: typeof window !== 'undefined' ? 'include' : undefined
    })
    const body = await req.json()

    if (!req.ok) {
        console.error(`[User/GetUser]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return body
}

/**
 * Updates the user profile
 * @param data Partial user data to update
 * @returns Updated user data
 */
export async function UpdateUser(data: Omit<UserUpdate, 'avatarHash'>): Promise<Omit<UserUpdate, 'avatar'>> {
    const form = new FormData()

    if (data.avatar) {
        form.append('avatar', data.avatar)
    }

    if (data.username) {
        form.append('username', data.username)
    }

    if (data.email) {
        form.append('email', data.email)
    }

    if (data.publicProfileEnabled !== undefined) {
        form.append('publicProfileEnabled', `${data.publicProfileEnabled}`)
    }

    const req = await fetch(`${PUBLIC_BASE_URL}/api/users/update`, {
        credentials: 'include',
        method: 'PATCH',
        body: form
    })
    const body = await req.json()

    if (req.status === 409) {
        throw new Error('Username already taken')
    }

    if (!req.ok) {
        console.error(`[User/UpdateUser]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return body
}
