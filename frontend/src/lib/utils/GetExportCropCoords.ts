import type { CropOpts } from '$lib/api-v2/Files'
import { cropH, cropW, cropX, cropY, isCroppingEnabled, videoH, videoW } from '$lib/stores/EditOptions'
import { get } from 'svelte/store'

export function GetExportCropCoords(): CropOpts {
    if (!get(isCroppingEnabled))
        return {
            x: 0,
            y: 0,
            w: 0,
            h: 0
        }

    const videoWidth = get(videoW)
    const videoHeight = get(videoH)

    return {
        x: Math.round(get(cropX) * videoWidth),
        y: Math.round(get(cropY) * videoHeight),
        w: Math.round(get(cropW) * videoWidth),
        h: Math.round(get(cropH) * videoHeight)
    }
}
