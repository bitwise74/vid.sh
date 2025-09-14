import { derived, writable } from 'svelte/store'

export const isCompressingDisabled = writable(false)
export const targetSize = writable(0)
export const losslessExport = writable(false)
export const exportFormat = writable('mp4')
export const exportFps = writable("Don't change")

export const trimStart = writable(0)
export const trimEnd = writable(0)
export const videoDuration = writable(0)
export const videoName = writable('')

export const videoSource = writable('')
export const selectedFile = writable<File | null>(null)

export const isExporting = writable(false)
export const isSaving = writable(false)

export const isCroppingEnabled = writable(false)

export const cropX = writable(0)
export const cropY = writable(0)
export const cropW = writable(0)
export const cropH = writable(0)

export const videoW = writable(0)
export const videoH = writable(0)

export const settingsUnchanged = derived(
    [trimStart, trimEnd, targetSize, exportFps, exportFormat, videoDuration],
    ([$trimStart, $trimEnd, $targetSize, $exportFps, $exportFormat, $videoDuration]) =>
        $trimStart === 0 && $trimEnd === $videoDuration && $targetSize === 0 && $exportFps === "Don't change" && $exportFormat === 'mp4'
)
