export type AuditEvent = {
    id: string
    type: 'login' | 'logout' | 'password_change' | 'profile_update' | 'video_upload' | 'video_delete'
    timestamp: string
    ip: string
    userAgent: string
    details?: Record<string, unknown>
}

export type AuditLogResponse = {
    events: AuditEvent[]
    page: number
    results: number
    total: number
}

/**
 * Fetches the audit log for the current user.
 * @param page
 * @param results
 * @returns
 */
export async function FetchAuditLog(page: number, results: number) {
    throw new Error('Not implemented')
}
