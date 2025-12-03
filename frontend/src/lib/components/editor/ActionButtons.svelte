<script lang="ts">
    import { goto, invalidateAll } from '$app/navigation'
    import { ExportVideo, UpdateFile, UploadFile, type Video } from '$lib/api/Files'
    import { isLoggedIn, jobStats, user } from '$lib/stores/AppVars'
    import { isExporting, isPreciseTrimming, isSaving, selectedFile, settingsUnchanged, videoName } from '$lib/stores/EditOptions'
    import { AutoCopy } from '$lib/utils/autoCopy'
    import { DownloadBlob } from '$lib/utils/downloadBlob'
    import { GetExportCropCoords } from '$lib/utils/getExportCropCoordinates'
    import { onMount } from 'svelte'
    import { toastStore } from '../../stores/ToastStore'
    import { losslessExport, targetSize, trimEnd, trimStart } from './../../stores/EditOptions'

    interface Props {
        videoID: string | null
        turnstileToken: string
    }

    let { videoID, turnstileToken }: Props = $props()

    let isBusy = $derived($isExporting || $isSaving)
    let isDisabled = $derived(isBusy || $settingsUnchanged)

    onMount(() => {
        if ('Notification' in window && Notification.permission !== 'granted' && Notification.permission !== 'denied') {
            Notification.requestPermission().then((permission) => {
                if (permission === 'granted') {
                    new Notification('Notifications Enabled', {
                        body: 'You will now receive desktop notifications if anything important happens',
                        icon: '/favicon.svg'
                    })
                }
            })
        }
    })

    function getEditOpts() {
        const crop = GetExportCropCoords()

        return {
            losslessExport: $losslessExport,
            targetSize: $targetSize,
            trimStart: $trimStart,
            trimEnd: $trimEnd,
            preciseTrimming: $isPreciseTrimming,
            cropH: crop.h,
            cropW: crop.w,
            cropX: crop.x,
            cropY: crop.y
        }
    }

    async function redirectToDashboard(video: Video) {
        await invalidateAll()
        await AutoCopy(video.file_key)
        await goto('/dashboard')
    }

    async function handleExport() {
        if (!$selectedFile) return // Typeguard

        isExporting.set(true)

        const processedBlob = await ExportVideo($selectedFile, getEditOpts(), turnstileToken, false)
            .catch((err) => {
                toastStore.error({
                    title: 'Failed to export video',
                    message: err.message,
                    duration: 10000
                })
                console.error(err)
                return
            })
            .finally(() => {
                setTimeout(() => {
                    isExporting.set(false)
                }, 3000)
            })

        if (processedBlob) {
            DownloadBlob(processedBlob, 'edited_' + $videoName)
        }
    }

    async function handleSave() {
        if (!$selectedFile) return

        isSaving.set(true)

        if ($settingsUnchanged) {
            const video = await UploadFile($selectedFile)
                .catch((err) => {
                    toastStore.error({
                        title: 'Failed to export video',
                        message: err.message,
                        duration: 10000
                    })
                    console.error(err)
                    return
                })
                .finally(() => {
                    isSaving.set(false)
                })

            if (!video) return

            await redirectToDashboard(video)
        }

        const video = await ExportVideo($selectedFile, getEditOpts(), turnstileToken, true)
            .catch((err) => {
                toastStore.error({
                    title: 'Failed to export video',
                    message: err.message,
                    duration: 10000
                })
                console.error(err)
                return
            })
            .finally(() => {
                isSaving.set(false)
            })

        if (!video) return // Typeguard

        await redirectToDashboard(video)
    }

    async function handleUpdate() {
        if (!videoID) return // Typeguard

        isSaving.set(true)

        const video = await UpdateFile(videoID, {
            processing_options: getEditOpts()
        })
            .catch((err) => {
                toastStore.error({
                    title: 'Updating video failed',
                    message: err.message,
                    duration: 10000
                })
                console.error('Failed to update video', err)
                return
            })
            .finally(() => {
                isSaving.set(false)
            })

        if (!video) return // Typeguard

        await redirectToDashboard(video)
    }
</script>
<!-- TODO: rewrite the player so it can handle very large files -->
<div class="border-top mt-4 pt-4">
    <div class="d-grid gap-3">
        <!-- Editing an existing video -->
        <!-- TODO: add save as copy option -->
        {#if videoID}
            <button class="btn btn-outline-warning" disabled={isDisabled} onclick={handleUpdate}>
                {#if $isSaving}
                    <span class="spinner-border spinner-border-sm me-2"></span>
                    Updating...
                {:else}
                    <i class="bi bi-pencil-square me-2"></i>
                    Update Video
                {/if}
            </button>
        {/if}

        <!-- Editing an uploaded video -->
        {#if !videoID}
            <button class="btn btn-primary" disabled={isDisabled} onclick={handleExport}>
                {#if $isExporting}
                    <span class="spinner-border spinner-border-sm me-2"></span>
                    Processing...
                {:else}
                    <i class="bi bi-download me-2"></i>
                    Export Video
                {/if}
            </button>
            {#if $isLoggedIn}
                <button class="btn btn-outline-primary" disabled={isDisabled} onclick={handleSave}>
                    {#if $isSaving}
                        <span class="spinner-border spinner-border-sm me-2"></span>
                        Saving...
                    {:else}
                        <i class="bi bi-cloud me-2"></i>
                        Save to Cloud
                    {/if}
                </button>
            {:else}
                <div class="bg-body-tertiary rounded p-3 text-center">
                    <p class="small text-muted mb-2">Sign in to save your videos to the cloud</p>
                    <a href="/login" class="btn btn-outline-secondary btn-sm">Sign In</a>
                </div>
            {/if}
        {/if}

        {#if $isExporting || $isSaving}
            <p class="small text-muted mb-0 mt-3 text-center upload-text">
                {#if $jobStats.state}
                    {#key $jobStats.state}
                        {#each $jobStats.state!.split('') as l, i}
                            <span style="animation-delay:{i * 0.01}s ">{l}</span>
                        {/each}
                    {/key}
                {/if}
            </p>
            <div class="progress" style="height: 5px;">
                <div
                    class="progress-bar progress-bar-animated bg-success"
                    role="progressbar"
                    style="width: {$jobStats.progress}%;"
                    aria-valuenow={$jobStats.progress}
                    aria-valuemin="0"
                    aria-valuemax="100">
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
    .upload-text span {
        display: inline-block;
        opacity: 0;
        transform: translateY(10px);
        filter: blur(4px);
        animation: fadeUpBlur 0.5s ease forwards;
    }

    .upload-text {
        font-variant-ligatures: none;
        word-spacing: normal;
        white-space: pre;
    }

    @keyframes fadeUpBlur {
        0% {
            opacity: 0;
            transform: translateY(10px);
            filter: blur(4px);
        }
        100% {
            opacity: 1;
            transform: translateY(0);
            filter: blur(0);
        }
    }
</style>
