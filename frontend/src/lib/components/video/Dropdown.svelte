<script lang="ts">
    import { PUBLIC_BASE_URL, PUBLIC_CDN_URL } from '$env/static/public'
    import { DeleteFiles, UpdateFile, type Video } from '$lib/api/Files'
    import { selectedVideos } from '$lib/stores/appControl'
    import { user } from '$lib/stores/AppVars'
    import { currentVideoURL, videos } from '$lib/stores/VideoStore'
    import { toastStore } from '../../stores/ToastStore'

    type Props = {
        video: Video
        isProfile?: boolean
        isListView?: boolean
        expanded?: boolean
        position?: { x: number; y: number } | null
    }

    let { video, isProfile = false, expanded = false, position = null }: Props = $props()

    const manualMenuStyle = $derived(expanded && position ? `--dropdown-x:${position.x}px;--dropdown-y:${position.y}px` : null)

    let showRenameModal = $state(false)
    let currentVideoId: string
    let renameValue = $state('')
    let oldName = ''
    let ext = ''

    async function handleVideoAction(action: string) {
        let videoURL = `${PUBLIC_CDN_URL}/${video.file_key}`

        ext = '.' + video.file_key.split('.')[video.file_key.split('.').length - 1]

        switch (action) {
            case 'play':
                currentVideoURL.set(videoURL)
                break
            case 'edit':
                window.location.href = `/editor?id=${video.file_key}`
                break
            case 'rename':
                // Quick and dirty, TODO: improve later
                oldName = video.name
                renameValue = video.name.replace('.mp4', '')
                currentVideoId = video.id
                showRenameModal = true
                break
            case 'download':
                window.location.href = videoURL
                break
            case 'share':
                if (localStorage.getItem('optRichEmbeds') === 'true') {
                    videoURL = `${PUBLIC_BASE_URL}/v/${video.file_key}`
                }

                if (navigator.share) {
                    await navigator.share({ url: videoURL })
                } else if (navigator.clipboard) {
                    await navigator.clipboard.writeText(videoURL)
                    toastStore.info({
                        title: 'Link copied to clipboard'
                    })
                }
                break
            case 'delete':
                const resp = await DeleteFiles([video.id])
                videos.delete('id', video.id)
                user.update((u) => {
                    u.stats = resp
                    return u
                })
                toastStore.success({
                    title: 'Video deleted'
                })

                // In case the video was bulk selected
                if ($selectedVideos.includes(video.id.toString())) {
                    selectedVideos.update((ids) => ids.filter((id) => id !== video.id.toString()))
                }

                break
            case 'toggle_private': {
                const newState = !video.private

                const newVid = await UpdateFile(video.id, { private: newState }).catch((err) => {
                    toastStore.error({
                        title: 'Failed to update video',
                        message: err.message
                    })
                    return
                })

                if (!newVid) return

                toastStore.success({
                    title: `File marked as ${!newState ? 'public' : 'private'}`
                })

                video.private = newState
                videos.replace(video.file_key, video)
                return
            }
            default:
                toastStore.warning({
                    title: 'Not implemented yet'
                })
                console.error('unknown video action', action)
        }
    }

    async function confirmRename() {
        try {
            const nameWithoutExt = renameValue.replace(/\?.*$/, '').replace(/\.[^/.]+$/, '')
            const fullName = nameWithoutExt + ext

            const newVid = await UpdateFile(currentVideoId, { name: fullName }).catch((err) => {
                toastStore.error({
                    title: 'Failed to rename video',
                    message: err.message,
                    duration: 10000
                })
            })
            if (!newVid) return

            newVid.thumbnail_url = `${PUBLIC_CDN_URL}/${newVid.file_key.replace('.mp4', '.webp')}`
            newVid.video_url = `${PUBLIC_CDN_URL}/${newVid.file_key}`

            toastStore.success({
                title: 'Video renamed'
            })

            videos.replace(video.file_key, newVid)
        } catch (err) {
            toastStore.error({
                title: 'Failed to rename video',
                message: err.message,
                duration: 10000
            })
            console.error(err)
        } finally {
            showRenameModal = false
        }
    }
</script>

<div class="dropdown-center">
    <button class="btn btn-sm btn-outline-secondary rounded-3" type="button" data-bs-toggle="dropdown" aria-label="Options" aria-expanded={expanded}>
        <i class="bi-three-dots fs-6 fw-bold"></i>
    </button>
    <ul class={`dropdown-menu animate${expanded && position ? ' show-manual' : ''}`} style={manualMenuStyle ?? undefined}>
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('play')}><i class="bi bi-play me-2"></i>Play</button>
        </li>
        <li><hr class="dropdown-divider" /></li>
        {#if !isProfile}
            <li>
                <button class="dropdown-item" onclick={() => handleVideoAction('rename')}><i class="bi bi-pencil me-2"></i>Rename</button>
            </li>
            <li>
                <button class="dropdown-item" disabled onclick={() => handleVideoAction('assign-labels')}><i class="bi bi-tags me-2"></i>Tags</button>
            </li>
        {/if}
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('share')}><i class="bi bi-share me-2"></i>Share</button>
        </li>
        <li>
            <button class="dropdown-item" onclick={() => handleVideoAction('download')}><i class="bi bi-download me-2"></i>Download</button>
        </li>
        {#if !isProfile}
            <li><hr class="dropdown-divider" /></li>
            <li>
                <button class="dropdown-item" onclick={() => handleVideoAction('edit')}><i class="bi bi-pencil me-2"></i>Edit</button>
            </li>
            <li>
                <button class="dropdown-item text-{video.private ? 'danger' : 'body'}" onclick={() => handleVideoAction('toggle_private')}>
                    <i class="bi-{video.private ? `eye` : `eye-slash`} me-2"></i>Make {video.private ? 'public' : 'private'}
                </button>
            </li>
            <li>
                <button class="dropdown-item text-danger" onclick={() => handleVideoAction('delete')}><i class="bi bi-trash me-2"></i>Delete</button>
            </li>
        {/if}
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
                    <button class="btn btn-primary" disabled={renameValue === video.name.replace('.mp4', '') || renameValue.trim() === ''} onclick={confirmRename}>Rename</button>
                </div>
            </div>
        </div>
    </div>
{/if}

<style>
    .animate {
        animation: slidefade-in 0.35s cubic-bezier(0.16, 1, 0.3, 1);
    }

    .dropdown-menu.show-manual {
        display: block;
        position: fixed;
        inset: auto;
        top: var(--dropdown-y, 0px);
        left: var(--dropdown-x, 0px);
        transform: translate(-50%, calc(-100% - 8px));
        z-index: 1080;
    }
</style>
