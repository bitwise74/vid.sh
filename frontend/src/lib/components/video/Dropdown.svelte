<script lang="ts">
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import { DeleteVideo, UpdateVideo } from '$lib/api/Files'
    import { currentVideoURL, stats, videos } from '$lib/stores/VideoStore'
    import { toastStore } from '../toast/toastStore'

    const { fileKey } = $props()

    let showRenameModal = $state(false)
    let currentVideoId: number
    let renameValue = $state('')
    let oldName = ''
    let ext = ''

    async function handleVideoAction(action: string) {
        const video = $videos.find((v) => v.file_key === fileKey)
        if (!video) return
        const videoURL = `${PUBLIC_CDN_URL}/${fileKey}`

        ext = '.' + fileKey.split('.')[fileKey.split('.').length - 1]

        switch (action) {
            case 'play':
                currentVideoURL.set(videoURL)
                break
            case 'edit':
                window.location.href = `/editor?id=${video.id}`
                break
            case 'rename':
                oldName = video.name
                renameValue = video.name
                currentVideoId = video.id
                showRenameModal = true
                break
            case 'assign-labels':
                break
            case 'download':
                window.location.href = videoURL
                break
            case 'share':
                if (navigator.share) {
                    await navigator.share({ url: videoURL })
                } else if (navigator.clipboard) {
                    await navigator.clipboard.writeText(videoURL)
                    toastStore.info({
                        title: 'Link copied to clipboard',
                        duration: 2000
                    })
                }
                break
            case 'copy-link':
                toastStore.info({
                    title: 'Link copied to clipboard',
                    duration: 2000
                })
                navigator.clipboard.writeText(videoURL)
                break
            case 'delete':
                const resp = await DeleteVideo(video.id)
                videos.set($videos.filter((v) => v.id != video.id))
                stats.set(resp)
                toastStore.success({
                    title: 'Video deleted',
                    duration: 2000
                })
                break
            default:
                toastStore.warning({
                    title: 'Not implemented yet'
                })
                console.error('unknown video action', action)
        }
    }

    async function confirmRename() {
        if (!renameValue || !currentVideoId) return
        try {
            const nameWithoutExt = renameValue.replace(/\?.*$/, '').replace(/\.[^/.]+$/, '')
            const fullName = nameWithoutExt + ext

            const newVid = await UpdateVideo(currentVideoId, { name: fullName })
            if (!newVid) return

            toastStore.success({
                title: 'Video renamed',
                duration: 3000
            })
            const vids = [...$videos]
            const idx = vids.findIndex((v) => v.name === oldName)
            if (idx !== -1) {
                vids.splice(idx, 1, newVid)
                videos.set(vids)
            }

            showRenameModal = false
        } catch (err) {
            toastStore.error({
                title: 'Failed to rename video',
                message: err.message,
                duration: 10000
            })
            console.error(err)
        }
    }
</script>

<div class="dropdown">
    <button class="btn btn-sm btn-outline-secondary" type="button" data-bs-toggle="dropdown" aria-label="Options" aria-expanded="false">
        <i class="bi bi-three-dots-vertical"></i>
    </button>
    <ul class="dropdown-menu">
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('play')}><i class="bi bi-play me-2"></i>Play</button>
        </li>
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('rename')}><i class="bi bi-pencil me-2"></i>Rename</button>
        </li>
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('edit')}><i class="bi bi-pencil me-2"></i>Edit</button>
        </li>
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('assign-labels')}><i class="bi bi-tags me-2"></i>Assign Labels</button>
        </li>
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('download')}><i class="bi bi-download me-2"></i>Download</button>
        </li>
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('share')}><i class="bi bi-share me-2"></i>Share</button>
        </li>
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('copy-link')}><i class="bi bi-link me-2"></i>Copy Link</button>
        </li>
        <li><hr class="dropdown-divider" /></li>
        <li>
            <button class="dropdown-item text-danger" onclick={() => handleVideoAction('delete')}><i class="bi bi-trash me-2"></i>Delete</button>
        </li>
    </ul>
</div>

{#if showRenameModal}
    <div class="modal show d-flex justify-content-center align-items-center" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
        <div class="modal-dialog modal-lg">
            <div class="modal-content">
                <div class="modal-body">
                    <input type="text" class="form-control" bind:value={renameValue} />
                </div>
                <div class="modal-footer">
                    <button class="btn btn-secondary" onclick={() => (showRenameModal = false)}>Cancel</button>
                    <button class="btn btn-primary" onclick={confirmRename}>Rename</button>
                </div>
            </div>
        </div>
    </div>
{/if}
