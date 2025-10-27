<script lang="ts">
    import { goto } from '$app/navigation'
    import { ExportVideo, UpdateFile, UploadFile } from '$lib/api/Files'
    import { shouldRefetch } from '$lib/stores/appControl'
    import { isLoggedIn, jobStats } from '$lib/stores/AppVars'
    import { isExporting, isSaving, selectedFile, settingsUnchanged, videoName } from '$lib/stores/EditOptions'
    import { AutoCopy } from '$lib/utils/audoCopy'
    import { DownloadBlob } from '$lib/utils/downloadBlob'
    import { GetExportCropCoords } from '$lib/utils/getExportCropCoordinates'
    import { toastStore } from '../../stores/ToastStore'
    import { losslessExport, targetSize, trimEnd, trimStart } from './../../stores/EditOptions'

    interface Props {
        videoID: string | null
        turnstileToken: string
    }

    let { videoID, turnstileToken }: Props = $props()
    const disableBtn = (): boolean => $isExporting || $isSaving || $settingsUnchanged

    async function redirectToDashboard() {
        shouldRefetch.set(true)
        await goto('/dashboard')
    }

    async function handleExport() {
        if (!$selectedFile) return // Typeguard

        isExporting.set(true)
        const crop = GetExportCropCoords()

        const processedBlob = await ExportVideo(
            $selectedFile,
            {
                losslessExport: $losslessExport,
                targetSize: $targetSize,
                trimStart: $trimStart,
                trimEnd: $trimEnd,
                cropH: crop.h,
                cropW: crop.w,
                cropX: crop.x,
                cropY: crop.y
            },
            turnstileToken,
            false
        )
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
        const crop = GetExportCropCoords()

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

            await redirectToDashboard()
            return await AutoCopy(video.id, video.file_key)
        }

        const video = await ExportVideo(
            $selectedFile,
            {
                losslessExport: $losslessExport,
                targetSize: $targetSize,
                trimStart: $trimStart,
                trimEnd: $trimEnd,
                cropH: crop.h,
                cropW: crop.w,
                cropX: crop.x,
                cropY: crop.y
            },
            turnstileToken,
            true
        )
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

        await redirectToDashboard()
        await AutoCopy(video.id, video.file_key)
    }

    async function handleUpdate() {
        if (!videoID) return // Typeguard

        isSaving.set(true)
        const crop = GetExportCropCoords()

        const video = await UpdateFile(videoID, {
            processing_options: {
                losslessExport: $losslessExport,
                targetSize: $targetSize,
                trimStart: $trimStart,
                trimEnd: $trimEnd,
                cropH: crop.h,
                cropW: crop.w,
                cropX: crop.x,
                cropY: crop.y
            }
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

        await redirectToDashboard()
        await AutoCopy(video.id, video.file_key)
    }
</script>
<!-- TODO: rewrite the player so it can handle very large files -->
<div class="border-top mt-4 pt-4">
    <div class="d-grid gap-3">
        <!-- Editing an existing video -->
        {#if videoID}
            <button class="btn btn-outline-warning" disabled={disableBtn()} onclick={handleUpdate}>
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
            <button class="btn btn-primary" disabled={disableBtn()} onclick={handleExport}>
                {#if $isExporting}
                    <span class="spinner-border spinner-border-sm me-2"></span>
                    Processing...
                {:else}
                    <i class="bi bi-download me-2"></i>
                    Export Video
                {/if}
            </button>
            {#if $isLoggedIn}
                <button class="btn btn-outline-primary" disabled={disableBtn()} onclick={handleSave}>
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
