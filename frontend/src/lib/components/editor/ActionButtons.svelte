<script lang="ts">
    import { goto } from '$app/navigation'
    import { ExportVideo, UpdateVideo, UploadVideo } from '$lib/api/Files'
    import { isLoggedIn } from '$lib/stores/AppVars'
    import { exportFormat, exportFps, isExporting, isSaving, selectedFile, settingsUnchanged, videoName } from '$lib/stores/EditOptions'
    import { DownloadBlob } from '$lib/utils/DownloadBlob'
    import { GetExportCropCoords } from '$lib/utils/GetExportCropCoords'
    import { jobProgress } from '../../../routes/editor/Logic'
    import { toastStore } from '../toast/toastStore'
    import { losslessExport, targetSize, trimEnd, trimStart } from './../../stores/EditOptions'

    interface Props {
        videoID: string | null
        turnstileToken: string
    }

    let { videoID, turnstileToken }: Props = $props()
    const disableBtn = (): boolean => $isExporting || $isSaving || $settingsUnchanged

    async function handleExport() {
        isExporting.set(true)

        try {
            if ($selectedFile) {
                // Handle for attached files (send to server)

                const processedBlob = await ExportVideo(
                    {
                        file: $selectedFile,
                        processingOpts: {
                            format: $exportFormat,
                            fps: parseInt($exportFps),
                            losslessExport: $losslessExport,
                            targetSize: $targetSize,
                            trimStart: $trimStart,
                            trimEnd: $trimEnd,
                            saveToCloud: false,
                            crop: GetExportCropCoords()
                        }
                    },
                    turnstileToken
                )

                DownloadBlob(processedBlob, 'edited_' + $videoName)
            } else {
                // Handle send instructions

                const video = await UpdateVideo(videoID!, {
                    processing_options: {
                        format: $exportFormat,
                        fps: parseInt($exportFps),
                        losslessExport: $losslessExport,
                        targetSize: $targetSize,
                        trimStart: $trimStart,
                        trimEnd: $trimEnd,
                        saveToCloud: false,
                        crop: GetExportCropCoords()
                    }
                })
                if (!video) return
                window.location.href = '/dashboard'
            }
        } catch (err) {
            toastStore.error({
                title: 'Failed to export video',
                message: 'Check the console for more details',
                duration: 10000
            })
            console.error(err)
        } finally {
            isExporting.set(false)
        }
    }

    async function handleSave() {
        if (!$selectedFile) return

        isSaving.set(true)

        try {
            if ($settingsUnchanged) {
                // Uploads can only happen with attached video files
                await UploadVideo($selectedFile)
                goto('/dashboard')
                return
            }

            await ExportVideo(
                {
                    file: $selectedFile,
                    processingOpts: {
                        format: $exportFormat,
                        fps: parseInt($exportFps),
                        losslessExport: $losslessExport,
                        targetSize: $targetSize,
                        trimStart: $trimStart,
                        trimEnd: $trimEnd,
                        saveToCloud: true,
                        crop: GetExportCropCoords()
                    }
                },
                turnstileToken
            )
            goto('/dashboard')
        } catch (err) {
            toastStore.error({
                title: 'Failed to save video',
                message: 'Check the console for more details',
                duration: 10000
            })
            console.error(err)
        } finally {
            isSaving.set(false)
        }
    }
</script>

<div class="border-top pt-4 mt-4">
    <div class="d-grid gap-3">
        {#if videoID}
            <button class="btn btn-outline-warning" disabled={disableBtn()} onclick={handleExport}>
                {#if $isSaving}
                    <span class="spinner-border spinner-border-sm me-2"></span>
                    Updating...
                {:else}
                    <i class="bi bi-pencil-square me-2"></i>
                    Update Video
                {/if}
            </button>
        {:else}
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
                <div class="text-center p-3 bg-light rounded">
                    <p class="small text-muted mb-2">Sign in to save your videos to the cloud</p>
                    <a href="/login" class="btn btn-outline-secondary btn-sm">Sign In</a>
                </div>
            {/if}
        {/if}

        {#if $isExporting || $isSaving}
            <div class="progress mt-2" style="height: 6px;">
                <div class="progress-bar progress-bar-animated bg-success" role="progressbar" style="width: {$jobProgress}%;" aria-valuenow={$jobProgress} aria-valuemin="0" aria-valuemax="100"></div>
            </div>
        {/if}
    </div>
</div>
