<script lang="ts">
    import { goto } from '$app/navigation'
    import { page } from '$app/state'
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import { TokenValidationResult, ValidateToken } from '$lib/api/Auth'
    import { ExportVideo, LoadVideo, OwnsVideo, UpdateVideo, UploadVideo } from '$lib/api/Files'
    import Header from '$lib/components/Header.svelte'
    import { toastStore } from '$lib/components/toast/toastStore'
    import VideoPlayer from '$lib/components/video/Player.svelte'
    import VideoUpload from '$lib/components/video/Upload.svelte'
    import { currentTime } from '$lib/stores/VideoStore'
    import { getCookie } from '$lib/utils/Cookies'
    import { onDestroy, onMount } from 'svelte'
    import RangeSlider from 'svelte-range-slider-pips'
    import { Turnstile } from 'svelte-turnstile'
    import { jobProgress } from './Logic'

    const videoID = page.url.searchParams.get('id')

    // Video metadata
    let videoName = $state('')
    let videoSize = $state(0)
    let videoDuration = $state(0)

    // The video is either a cloudfront url or an attached file
    let videoSrc = $state('')
    let selectedVideo: File | null = $state(null)

    let isExporting = $state(false)
    let isSaving = $state(false)
    let isEditing = $state(false)
    let isCompressingDisabled = $state(false)

    // Export options
    let targetSize: number = $state(0)
    let trimStart: number = $state(0)
    let trimEnd: number = $state(0)
    let exportFormat = $state('mp4')
    let exportFps = $state("Don't change")
    let losslessExport = $state(false)

    // Other things needed
    let isLoggedIn = $state(false)
    let turnstileToken = ''

    const settingsUnchanged = (): boolean => {
        return (
            targetSize === -1 &&
            trimStart === 0 &&
            trimEnd === videoDuration &&
            exportFormat === 'mp4' &&
            exportFps === "Don't change"
        )
    }

    // Mainly checks for when the video ID query is present
    onMount(async () => {
        if (getCookie('logged_in') == '1') {
            isLoggedIn = (await ValidateToken()) == TokenValidationResult.VALID
        }

        if (!videoID) return

        // Reject if not logged in but trying to edit existing video
        if (!isLoggedIn) {
            toastStore.error({
                title: 'You must be logged in to edit videos',
                dismissible: false,
                duration: 10000
            })
            return
        }

        const owns = await OwnsVideo(videoID).catch((err) => {
            console.error(err)
        })

        if (!owns) {
            toastStore.error({
                title: "You don't own this video",
                duration: 10000
            })
            return
        }

        isEditing = true

        // Load video if logged in and owns video
        const videoData = await LoadVideo(videoID).catch((err) => {
            toastStore.error({
                title: 'Failed to load video',
                message: err.message
            })
            console.error(err)
        })

        if (!videoData) return

        videoName = videoData.name
        videoSize = parseFloat((videoData.size / (1024 * 1024)).toFixed(2))
        videoDuration = videoData.duration
        videoSrc = `${PUBLIC_CDN_URL}/${videoData?.file_key}`
        trimEnd = videoDuration
    })

    onDestroy(() => {
        if (videoSrc) URL.revokeObjectURL(videoSrc)
    })

    const handleVideoSelect = (f: File) => {
        if (videoSrc) URL.revokeObjectURL(videoSrc)
        videoSrc = URL.createObjectURL(f)

        // This assignment changes what the processing function does. That is
        // if a video is not selected the function will only send edit instructions
        // to the server for a file in the S3 server.
        selectedVideo = f

        if (f.type !== 'video/mp4') {
            toastStore.warning({
                title: 'Video will be remuxed',
                message:
                    'Files other than mp4 are not officially supported. Editing tools should work fine, but the video itself will be remuxed into an mp4 file',
                duration: 30000
            })
        }

        if (f.type === 'video/x-matroska') {
            toastStore.warning({
                title: '.mkv file detected',
                message:
                    "Due to limitations from the browser playback of .mkv files doesn't work. Your settings will still apply to them and they'll be converted to mp4 files upon export/upload, but you won't be able to preview them here. Some settings may be entirely broken.",
                duration: 30000
            })

            toastStore.info({
                title: 'Possible solution',
                message:
                    "Upload your video to the dashboard first to convert it into an mp4 file, then click on it's dropdown menu and choose 'edit'",
                duration: 30000
            })
        }

        videoSize = parseFloat((f.size / (1024 * 1024)).toFixed(2))
        videoName = f.name
        isCompressingDisabled = videoSize < 5

        const video = document.createElement('video')

        const onMeta = () => {
            videoDuration = video.duration
            trimEnd = video.duration

            video.removeEventListener('loadedmetadata', onMeta)
            video.src = ''
            video.load()
        }

        video.addEventListener('loadedmetadata', onMeta)
        video.src = videoSrc
    }

    const handleVideoClear = () => {
        videoSrc = ''
        selectedVideo = null
        videoDuration = 0
        currentTime.set(0)
        trimStart = 0
        trimEnd = 0
        videoSize = 0
        exportFormat = 'mp4'
        exportFps = "Don't change"
        if (videoSrc) URL.revokeObjectURL(videoSrc)
    }

    const handleCompressChange = (e: CustomEvent<{ value: number }>) => {
        targetSize = e.detail.value
    }

    function handleTrimChange(e: CustomEvent<{ values: number[] }>) {
        const [start, end] = e.detail.values
        trimStart = start
        currentTime.set(start)
        trimEnd = end
    }

    async function handleExport() {
        isExporting = true

        try {
            if (selectedVideo) {
                // Handle for attached files (send to server)

                const processedBlob = await ExportVideo(
                    {
                        file: selectedVideo,
                        processingOpts: {
                            format: exportFormat,
                            fps: parseInt(exportFps),
                            losslessExport,
                            targetSize,
                            trimEnd,
                            trimStart,
                            saveToCloud: false
                        }
                    },
                    turnstileToken
                )

                downloadBlob(processedBlob, 'edited_' + videoName)
            } else {
                // Handle send instructions

                const video = await UpdateVideo(videoID!, {
                    processing_options: {
                        format: exportFormat,
                        fps: parseInt(exportFps),
                        losslessExport,
                        trimEnd,
                        trimStart,
                        targetSize,
                        saveToCloud: false
                    }
                })
                if (!video) return
                window.location.href = '/dashboard'
            }
        } catch (error) {
            toastStore.error({
                title: 'Failed to edit video',
                message: error.message,
                duration: 10000
            })
        } finally {
            isExporting = false
        }
    }

    async function handleSaveToCloud() {
        if (!selectedVideo) return

        isSaving = true

        try {
            if (settingsUnchanged()) {
                // Uploads can only happen with attached video files
                await UploadVideo(selectedVideo)
                goto('/dashboard')
                return
            }

            await ExportVideo(
                {
                    file: selectedVideo,
                    processingOpts: {
                        format: exportFormat,
                        fps: parseInt(exportFps),
                        losslessExport,
                        targetSize,
                        trimEnd,
                        trimStart,
                        saveToCloud: true
                    }
                },
                turnstileToken
            )
            goto('/dashboard')
        } catch (error) {
            toastStore.error({
                title: 'Failed to save video to cloud',
                message: error.message,
                duration: 10000
            })
            console.error(error)
        } finally {
            isSaving = false
        }
    }

    function downloadBlob(blob: Blob, filename: string) {
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = filename
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        URL.revokeObjectURL(url)
    }

    function formatTime(time: number) {
        const minutes = Math.floor(time / 60)
        const seconds = Math.floor(time % 60)
        const milliseconds = Math.floor((time % 1) * 100)
        return `${minutes}:${seconds.toString().padStart(2, '0')}.${milliseconds.toString().padStart(2, '0')}`
    }

    function parseTime(str: string): number | null {
        const parts = str.split(':').map(Number)
        if (parts.some(isNaN)) return null
        if (parts.length === 2) return parts[0] * 60 + parts[1]
        return parts[0]
    }

    function clampTrimValues(start: number, end: number): [number, number] {
        const clampedStart = Math.max(0, Math.min(start, videoDuration))
        const clampedEnd = Math.max(clampedStart, Math.min(end, videoDuration))
        return [clampedStart, clampedEnd]
    }

    function handleTrimStartChange(e: FocusEvent) {
        const target = e.target as HTMLSpanElement
        const val = parseTime(target.innerText)
        if (val !== null) {
            ;[trimStart, trimEnd] = clampTrimValues(val, trimEnd)
        }
        target.innerText = formatTime(trimStart)
    }

    function handleTrimEndChange(e: FocusEvent) {
        const target = e.target as HTMLSpanElement
        const val = parseTime(target.innerText)
        if (val !== null) {
            ;[trimStart, trimEnd] = clampTrimValues(trimStart, val)
        }
        target.innerText = formatTime(trimEnd)
    }

    function handleTargetSizeChange(e: FocusEvent) {
        const target = e.target as HTMLSpanElement
        let val = parseInt(target.innerText)

        if (isNaN(val)) {
            val = 0
        }

        target.innerText = `${val}`
    }

    function onTimeUpdate(video: HTMLVideoElement) {
        if (video.currentTime >= trimEnd) {
            video.pause()
            video.currentTime = trimEnd
        }
    }
</script>

<svelte:head>
    <title>Video Editor - vid.sh</title>
</svelte:head>

<div class="min-vh-100 bg-light">
    <Header title="Editor" page="editor" />

    <Turnstile
        siteKey="0x4AAAAAABkH5R_4hvXLiZqn"
        appearance="interaction-only"
        on:turnstile-callback={(e: CustomEvent<{ token: string }>) => {
            turnstileToken = e.detail.token
        }} />

    <main class="container py-4">
        <div class="row g-4">
            {#if videoSrc}
                <div class="col-lg-8">
                    <div class="mb-4">
                        <VideoUpload
                            videoSelected
                            {videoName}
                            {videoSize}
                            onClear={handleVideoClear} />
                    </div>
                    <VideoPlayer
                        src={videoSrc}
                        duration={videoDuration}
                        {trimEnd}
                        {trimStart}
                        {onTimeUpdate} />
                </div>
            {:else}
                <div class="mb-4">
                    <VideoUpload
                        onVideoSelect={handleVideoSelect}
                        videoSelected={false}
                        onClear={handleVideoClear} />
                </div>
            {/if}

            {#if videoSrc}
                <div class="col-lg-4">
                    <div class="card shadow">
                        <div class="card-body">
                            <!-- Tabs -->
                            <ul class="nav nav-tabs mb-4" role="tablist">
                                <li class="nav-item" role="presentation">
                                    <button
                                        class="nav-link active"
                                        data-bs-toggle="tab"
                                        data-bs-target="#trim-tab"
                                        type="button">
                                        <i class="bi-scissors me-1"></i> Trim
                                    </button>
                                </li>
                                <li class="nav-item" role="presentation">
                                    <button
                                        class="nav-link"
                                        data-bs-toggle="tab"
                                        data-bs-target="#crop-tab"
                                        disabled
                                        type="button">
                                        <i class="bi-crop me-1"></i> Crop
                                    </button>
                                </li>
                                <li class="nav-item" role="presentation">
                                    <button
                                        class="nav-link"
                                        data-bs-toggle="tab"
                                        data-bs-target="#compress-tab"
                                        type="button">
                                        <i class="bi-file-earmark-arrow-down me-1"></i> Compress
                                    </button>
                                </li>
                                <li class="nav-item" role="presentation">
                                    <button
                                        class="nav-link"
                                        data-bs-toggle="tab"
                                        data-bs-target="#export-tab"
                                        type="button">
                                        <i class="bi-box-arrow-up me-1"></i> Export
                                    </button>
                                </li>
                            </ul>

                            <div class="tab-content">
                                <!-- Trim Tab -->
                                <div class="tab-pane show active" id="trim-tab">
                                    <p>
                                        From:
                                        <span
                                            class="trim-time editable"
                                            contenteditable="true"
                                            role="textbox"
                                            tabindex="0"
                                            onblur={handleTrimStartChange}
                                            onkeydown={(e: KeyboardEvent) => {
                                                if (e.key === 'Enter') {
                                                    ;(e.target as HTMLSpanElement).blur()
                                                    e.preventDefault()
                                                }
                                            }}>
                                            {formatTime(trimStart)}
                                        </span>
                                        to:
                                        <span
                                            class="trim-time editable"
                                            contenteditable="true"
                                            role="textbox"
                                            tabindex="0"
                                            onblur={handleTrimEndChange}
                                            onkeydown={(e: KeyboardEvent) => {
                                                if (e.key === 'Enter') {
                                                    ;(e.target as HTMLSpanElement).blur()
                                                    e.preventDefault()
                                                }
                                            }}>
                                            {formatTime(trimEnd)}
                                        </span>
                                    </p>

                                    <div class="mb-3">
                                        <RangeSlider
                                            spring={false}
                                            values={[trimStart, trimEnd]}
                                            min={0}
                                            max={videoDuration}
                                            step={0.01}
                                            float
                                            range
                                            draggy={true}
                                            on:stop={handleTrimChange}
                                            formatter={formatTime}
                                            style="--range-range: #1463DA">
                                        </RangeSlider>

                                        <div class="d-flex justify-content-between">
                                            <button
                                                aria-label="Set Trim Start to Current Time"
                                                class="btn btn-outline-primary"
                                                title="Set Trim Start to Current Time"
                                                onclick={() => (trimStart = $currentTime)}>
                                                <i class="bi bi-skip-start-fill"></i>
                                            </button>
                                            <button
                                                aria-label="Set Trim End to Current Time"
                                                class="btn btn-outline-primary"
                                                title="Set Trim End to Current Time"
                                                onclick={() => (trimEnd = $currentTime)}>
                                                <i class="bi bi-skip-end-fill"></i>
                                            </button>
                                        </div>
                                    </div>

                                    <div class="small text-muted">
                                        New duration: {formatTime(trimEnd - trimStart)}<br />
                                        Tip: Click on the timestamps to edit
                                    </div>
                                </div>

                                <!-- Compress Tab -->
                                <div class="tab-pane" id="compress-tab">
                                    <div class="mb-3">
                                        <label
                                            for="target-size"
                                            class={isCompressingDisabled ? 'text-muted' : ''}>
                                            Target Size:
                                            {#if targetSize <= 0 || targetSize >= videoSize}
                                                Not set
                                            {:else}
                                                <span
                                                    class="target-size editable"
                                                    contenteditable="true"
                                                    role="textbox"
                                                    tabindex="0"
                                                    onblur={handleTargetSizeChange}
                                                    onkeydown={(e: KeyboardEvent) => {
                                                        if (e.key === 'Enter') {
                                                            ;(e.target as HTMLSpanElement).blur()
                                                            e.preventDefault()
                                                        }
                                                    }}>
                                                    {targetSize}
                                                </span>MB
                                            {/if}
                                        </label>

                                        <RangeSlider
                                            spring={false}
                                            value={0}
                                            disabled={isCompressingDisabled}
                                            min={0}
                                            max={videoSize}
                                            step={1}
                                            float
                                            on:change={handleCompressChange}
                                            formatter={(val) => {
                                                if (val == 0 || val == videoSize) {
                                                    targetSize = -1
                                                    return 'Not set'
                                                }

                                                return val.toString() + ' MB'
                                            }}></RangeSlider>
                                    </div>
                                    <div class="mb-3">
                                        <label class="form-label" for="processing-speed"
                                            >Other processing options</label>
                                        <p class="small text-muted">
                                            Lossless exports produce a better video quality, but take longer to finish and don't optimize the video size.
                                        </p>
                                        <input
                                            class="form-check-input"
                                            type="checkbox"
                                            id="losslessExport"
                                            bind:checked={losslessExport} />
                                        <label class="form-check-label small" for="losslessExport"
                                            >Lossless export</label>
                                    </div>
                                </div>

                                <!-- Export Tab -->
                                <div class="tab-pane" id="export-tab">
                                    <div class="mb-3">
                                        <label class="form-label" for="export-format">Format</label>
                                        <p class="small text-muted">
                                            Some formats may take longer to process
                                        </p>
                                        <select
                                            class="form-select"
                                            id="export-format"
                                            disabled
                                            bind:value={exportFormat}>
                                            <option value="mp4">mp4</option>
                                            <option value="mov">mov</option>
                                            <option value="avi">avi</option>
                                            <option value="gif">gif</option>
                                        </select>
                                    </div>
                                    <div class="mb-3">
                                        <label class="form-label" for="export-framerate">
                                            Framerate
                                        </label>
                                        <p class="small text-muted">
                                            Reducing the framerate will also reduce the file size
                                        </p>

                                        <select
                                            class="form-select"
                                            id="export-framerate"
                                            disabled
                                            bind:value={exportFps}>
                                            <option value="Don't change">Don't change</option>
                                            <option value="60">60</option>
                                            <option value="30">30</option>
                                            <option value="24">24</option>
                                        </select>
                                    </div>
                                </div>
                            </div>

                            <!-- Action Buttons -->
                            <div class="border-top pt-4 mt-4">
                                <div class="d-grid gap-3">
                                    {#if videoID && isEditing}
                                        <!-- Show only Update Video button -->
                                        <button
                                            class="btn btn-outline-warning"
                                            disabled={isExporting ||
                                                isSaving ||
                                                settingsUnchanged()}
                                            onclick={() => handleExport()}>
                                            {#if isSaving}
                                                <span class="spinner-border spinner-border-sm me-2"
                                                ></span>
                                                Updating...
                                            {:else}
                                                <i class="bi bi-pencil-square me-2"></i>
                                                Update Video
                                            {/if}
                                        </button>
                                    {:else}
                                        <!-- Always show Export Video button -->
                                        <button
                                            class="btn btn-primary"
                                            onclick={() => handleExport()}
                                            disabled={isExporting ||
                                                isSaving ||
                                                settingsUnchanged()}>
                                            {#if isExporting}
                                                <span class="spinner-border spinner-border-sm me-2"
                                                ></span>
                                                Processing...
                                            {:else}
                                                <i class="bi bi-download me-2"></i>
                                                Export Video
                                            {/if}
                                        </button>
                                        <!-- Show normal Save / Sign In buttons -->
                                        {#if isLoggedIn}
                                            <button
                                                class="btn btn-outline-primary"
                                                disabled={isExporting || isSaving}
                                                onclick={() => handleSaveToCloud()}>
                                                {#if isSaving}
                                                    <span
                                                        class="spinner-border spinner-border-sm me-2"
                                                    ></span>
                                                    Saving...
                                                {:else}
                                                    <i class="bi bi-cloud me-2"></i>
                                                    Save to Cloud
                                                {/if}
                                            </button>
                                        {:else}
                                            <div class="text-center p-3 bg-light rounded">
                                                <p class="small text-muted mb-2">
                                                    Sign in to save your videos to the cloud
                                                </p>
                                                <a
                                                    href="/login"
                                                    class="btn btn-outline-secondary btn-sm"
                                                    >Sign In</a>
                                            </div>
                                        {/if}
                                    {/if}

                                    {#if isExporting || isSaving}
                                        <div class="progress mt-2" style="height: 6px;">
                                            <div
                                                class="progress-bar progress-bar-animated bg-success"
                                                role="progressbar"
                                                style="width: {$jobProgress}%;"
                                                aria-valuenow={$jobProgress}
                                                aria-valuemin="0"
                                                aria-valuemax="100">
                                            </div>
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            {/if}
        </div>
    </main>
</div>
