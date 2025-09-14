<script lang="ts">
    import { page } from '$app/state'
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import { TokenValidationResult, ValidateToken } from '$lib/api/Auth'
    import { LoadVideo, OwnsVideo } from '$lib/api/Files'
    import ActionButtons from '$lib/components/editor/ActionButtons.svelte'
    import Compress from '$lib/components/editor/tabs/Compress.svelte'
    import Crop from '$lib/components/editor/tabs/Crop.svelte'
    import Trim from '$lib/components/editor/tabs/Trim.svelte'
    import Header from '$lib/components/Header.svelte'
    import { toastStore } from '$lib/components/toast/toastStore'
    import VideoPlayer from '$lib/components/video/Player.svelte'
    import VideoUpload from '$lib/components/video/Upload.svelte'
    import { isLoggedIn } from '$lib/stores/AppVars'
    import { exportFormat, exportFps, selectedFile, trimEnd, trimStart, videoDuration, videoSource } from '$lib/stores/EditOptions'
    import { currentTime } from '$lib/stores/VideoStore'
    import { getCookie } from '$lib/utils/Cookies'
    import { onDestroy, onMount } from 'svelte'
    import { Turnstile } from 'svelte-turnstile'

    const videoID = page.url.searchParams.get('id')

    let videoName = $state('')
    let videoSize = $state(0)
    let turnstileToken = $state('')

    // Mainly checks for when the video ID query is present
    onMount(async () => {
        if (getCookie('logged_in') == '1') {
            isLoggedIn.set((await ValidateToken()) == TokenValidationResult.VALID)
        }

        if (!videoID) return

        // Reject if not logged in but trying to edit existing video
        if (!$isLoggedIn) {
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
        videoDuration.set(videoData.duration)
        videoSource.set(`${PUBLIC_CDN_URL}/${videoData?.file_key}`)
        trimEnd.set(videoData.duration)
    })

    onDestroy(() => {
        if (videoSource) URL.revokeObjectURL($videoSource)
        handleVideoClear()
    })

    const handleVideoSelect = (f: File) => {
        if ($videoSource) URL.revokeObjectURL($videoSource)
        videoSource.set(URL.createObjectURL(f))

        // This assignment changes what the processing function does. That is
        // if a video is not selected the function will only send edit instructions
        // to the server for a file in the S3 server.
        selectedFile.set(f)

        if (f.type !== 'video/mp4') {
            toastStore.warning({
                title: 'Video will be remuxed',
                message: 'Files other than mp4 are not officially supported. Editing tools should work fine, but the video itself will be remuxed into an mp4 file',
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
                message: "Upload your video to the dashboard first to convert it into an mp4 file, then click on it's dropdown menu and choose 'edit'",
                duration: 30000
            })
        }

        videoSize = parseFloat((f.size / (1024 * 1024)).toFixed(2))
        videoName = f.name

        const video = document.createElement('video')

        const onMeta = () => {
            videoDuration.set(video.duration)
            console.log(video.duration)
            trimEnd.set(video.duration)

            video.removeEventListener('loadedmetadata', onMeta)
            video.src = ''
            video.load()
        }

        video.addEventListener('loadedmetadata', onMeta)
        video.src = $videoSource
    }

    const handleVideoClear = () => {
        videoSource.set('')
        selectedFile.set(null)
        videoDuration.set(0)
        currentTime.set(0)
        trimStart.set(0)
        trimEnd.set(0)
        videoSize = 0
        exportFormat.set('mp4')
        exportFps.set("Don't change")
        if ($videoSource) URL.revokeObjectURL($videoSource)
    }

    function onTimeUpdate(video: HTMLVideoElement) {
        if (video.currentTime >= $trimEnd) {
            video.pause()
            video.currentTime = $trimEnd
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
            {#if $videoSource}
                <div class="col-lg-8">
                    <div class="mb-4">
                        <VideoUpload videoSelected {videoName} {videoSize} onClear={handleVideoClear} />
                    </div>
                    <VideoPlayer {onTimeUpdate} />
                </div>
            {:else}
                <div class="mb-4">
                    <VideoUpload onVideoSelect={handleVideoSelect} videoSelected={false} onClear={handleVideoClear} />
                </div>
            {/if}

            {#if $videoSource}
                <div class="col-lg-4">
                    <div class="card shadow">
                        <div class="card-body">
                            <ul class="nav nav-tabs mb-4" role="tablist">
                                <li class="nav-item" role="presentation">
                                    <button class="nav-link active" data-bs-toggle="tab" data-bs-target="#trim-tab" type="button">
                                        <i class="bi-scissors me-1"></i> Trim
                                    </button>
                                </li>
                                <li class="nav-item" role="presentation">
                                    <button class="nav-link" data-bs-toggle="tab" data-bs-target="#crop-tab" type="button">
                                        <i class="bi-crop me-1"></i> Crop
                                    </button>
                                </li>
                                <li class="nav-item" role="presentation">
                                    <button class="nav-link" data-bs-toggle="tab" data-bs-target="#compress-tab" type="button">
                                        <i class="bi-file-earmark-arrow-down me-1"></i> Compress
                                    </button>
                                </li>
                            </ul>

                            <div class="tab-content">
                                <Trim />
                                <Compress {videoSize} />
                                <Crop />
                            </div>

                            <ActionButtons {videoID} {turnstileToken} />
                        </div>
                    </div>
                </div>
            {/if}
        </div>
    </main>
</div>
