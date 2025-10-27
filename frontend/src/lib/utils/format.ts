/**
 * Formats bytes into a human-readable string.
 * @param bytes Number of bytes
 * @returns Formatted size string
 */
export function FormatSize(bytes: number): string {
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
        if (bytes === 0) return '0 B'
        const i = Math.floor(Math.log(bytes) / Math.log(1024))
        return parseFloat((bytes / Math.pow(1024, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * Formats duration in seconds to HH:MM:SS format.
 * @param seconds Duration in seconds
 * @returns Formatted duration string
 */
export function FormatDuration(seconds: number): string {
        const hrs = Math.floor(seconds / 3600)
        const mins = Math.floor((seconds % 3600) / 60)
        const secs = Math.floor(seconds % 60)

        const parts: Array<string> = []
        if (hrs > 0) parts.push(hrs.toString().padStart(2, '0'))
        parts.push(mins.toString().padStart(2, '0'))
        parts.push(secs.toString().padStart(2, '0'))

        return parts.join(':')
}

/**
 * Formats a UNIX timestamp to a human-readable date string.
 * @param unix UNIX timestamp in seconds
 * @returns Formatted date string
 */
export function FormatDate(unix: number): string {
        const date = new Date(unix * 1000)
        return date.toLocaleDateString(undefined, {
                year: 'numeric',
                month: 'short',
                day: 'numeric'
        })
}
