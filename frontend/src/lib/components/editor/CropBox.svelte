<script lang="ts">
    import { cropH, cropW, cropX, cropY, isCroppingEnabled, videoH, videoW } from '$lib/stores/EditOptions'
    import { onDestroy, onMount } from 'svelte'
    import { get, writable } from 'svelte/store'

    let cropBox: HTMLElement
    let { video }: { video: HTMLVideoElement } = $props()

    let isDragging = writable(false)
    let isResizing = writable(false)

    onMount(() => {
        video?.addEventListener('loadedmetadata', applyCropBox)
        window.addEventListener('resize', applyCropBox)

        window.addEventListener('mouseup', () => {
            if (get(isDragging)) isDragging.set(false)
            if (get(isResizing)) isResizing.set(false)

            updateCropStores()
        })

        cropBox.addEventListener('mousedown', handleDragStart)

        const handles = cropBox.querySelectorAll('.handle')
        handles.forEach((handle) => {
            handle.addEventListener('mousedown', (e) => handleResizeStart(e as MouseEvent))
        })
    })

    onDestroy(() => {
        video?.removeEventListener('loadedmetadata', applyCropBox)
        window.removeEventListener('resize', applyCropBox)
    })
    $effect(() => {
        if (!$isCroppingEnabled) {
            cropBox.classList.add('d-none')
        } else {
            cropBox.classList.remove('d-none')
        }
    })

    function applyCropBox() {
        const rect = video.getBoundingClientRect()

        videoW.set(video.videoWidth)
        videoH.set(video.videoHeight)

        const defaultWidth = video.videoWidth * 0.5
        const defaultHeight = video.videoHeight * 0.5

        if (get(cropW) === 0) cropW.set(defaultWidth / video.videoWidth)
        if (get(cropH) === 0) cropH.set(defaultHeight / video.videoHeight)

        if (get(cropX) === 0) cropX.set(0.5 - get(cropW) / 2)
        if (get(cropY) === 0) cropY.set(0.5 - get(cropH) / 2)

        cropBox.style.left = get(cropX) * rect.width + 'px'
        cropBox.style.top = get(cropY) * rect.height + 'px'
        cropBox.style.width = get(cropW) * rect.width + 'px'
        cropBox.style.height = get(cropH) * rect.height + 'px'
    }

    function updateCropStores() {
        const rect = video.getBoundingClientRect()

        const left = cropBox.offsetLeft
        const top = cropBox.offsetTop
        const width = cropBox.offsetWidth
        const height = cropBox.offsetHeight

        cropX.set(left / rect.width)
        cropY.set(top / rect.height)
        cropW.set(width / rect.width)
        cropH.set(height / rect.height)
    }

    function handleDragStart(e: MouseEvent) {
        if ((e.target as HTMLElement).classList.contains('handle')) return

        isDragging.set(true)
        const startX = e.clientX
        const startY = e.clientY
        const startLeft = cropBox.offsetLeft
        const startTop = cropBox.offsetTop

        function onMouseMove(ev: MouseEvent) {
            if (!get(isDragging)) return

            const dx = ev.clientX - startX
            const dy = ev.clientY - startY

            const rect = video.getBoundingClientRect()
            const boxW = cropBox.offsetWidth
            const boxH = cropBox.offsetHeight

            const newLeft = clamp(startLeft + dx, 0, rect.width - boxW)
            const newTop = clamp(startTop + dy, 0, rect.height - boxH)

            cropBox.style.left = newLeft + 'px'
            cropBox.style.top = newTop + 'px'
        }

        function onMouseUp() {
            isDragging.set(false)
            updateCropStores()
            window.removeEventListener('mousemove', onMouseMove)
            window.removeEventListener('mouseup', onMouseUp)
        }

        window.addEventListener('mousemove', onMouseMove)
        window.addEventListener('mouseup', onMouseUp)
    }

    function handleResizeStart(e: MouseEvent) {
        e.stopPropagation()
        isResizing.set(true)

        const handle = e.target as HTMLElement
        const handleType = handle.dataset.pos
        const startX = e.clientX
        const startY = e.clientY
        const startLeft = cropBox.offsetLeft
        const startTop = cropBox.offsetTop
        const startWidth = cropBox.offsetWidth
        const startHeight = cropBox.offsetHeight
        const rect = video.getBoundingClientRect()

        function onMouseMove(ev: MouseEvent) {
            if (!get(isResizing)) return

            const dx = ev.clientX - startX
            const dy = ev.clientY - startY

            let newLeft = startLeft
            let newTop = startTop
            let newWidth = startWidth
            let newHeight = startHeight

            if (handleType?.includes('right')) newWidth = clamp(startWidth + dx, 20, rect.width - startLeft)
            if (handleType?.includes('bottom')) newHeight = clamp(startHeight + dy, 20, rect.height - startTop)
            if (handleType?.includes('left')) {
                newWidth = clamp(startWidth - dx, 20, startLeft + startWidth)
                newLeft = clamp(startLeft + dx, 0, startLeft + startWidth - 20)
            }
            if (handleType?.includes('top')) {
                newHeight = clamp(startHeight - dy, 20, startTop + startHeight)
                newTop = clamp(startTop + dy, 0, startTop + startHeight - 20)
            }

            cropBox.style.left = newLeft + 'px'
            cropBox.style.top = newTop + 'px'
            cropBox.style.width = newWidth + 'px'
            cropBox.style.height = newHeight + 'px'
        }

        function onMouseUp() {
            isResizing.set(false)
            updateCropStores()
            window.removeEventListener('mousemove', onMouseMove)
            window.removeEventListener('mouseup', onMouseUp)
        }

        window.addEventListener('mousemove', onMouseMove)
        window.addEventListener('mouseup', onMouseUp)
    }

    function clamp(val: number, min: number, max: number) {
        return Math.min(Math.max(val, min), max)
    }
</script>

<div id="crop-box" bind:this={cropBox}>
    <div class="handle top-left" data-pos="top-left"></div>
    <div class="handle top-right" data-pos="top-right"></div>
    <div class="handle bottom-left" data-pos="bottom-left"></div>
    <div class="handle bottom-right" data-pos="bottom-right"></div>
</div>

<style>
    #crop-box {
        position: absolute;
        top: 50px;
        left: 50px;
        width: 200px;
        height: 200px;
        border: 2px solid white;
        background: rgba(255, 255, 255, 0.2);
        cursor: move;
        box-sizing: border-box;
    }

    .handle {
        position: absolute;
        width: 10px;
        height: 10px;
        background: white;

        border: 1px solid black;
    }

    .top-left {
        top: -5px;
        left: -5px;
        cursor: nw-resize;
    }
    .top-right {
        top: -5px;
        right: -5px;
        cursor: ne-resize;
    }
    .bottom-left {
        bottom: -5px;
        left: -5px;
        cursor: sw-resize;
    }
    .bottom-right {
        bottom: -5px;
        right: -5px;
        cursor: se-resize;
    }
</style>
