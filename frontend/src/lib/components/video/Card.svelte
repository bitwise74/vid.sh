<script lang="ts">
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import type { Video } from '$lib/api/Files'
    import { currentVideoURL } from '$lib/stores/VideoStore'
    import VideoDropdown from './Dropdown.svelte'

    let { video }: { video: Video } = $props()

    function formatDuration(seconds: number): string {
        const h = Math.floor(seconds / 3600)
        const m = Math.floor((seconds % 3600) / 60)
        const s = Math.floor(seconds % 60)

        if (h > 0) {
            return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
        } else {
            return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
        }
    }

    function formatSize(bytes: number): string {
        if (bytes < 1024) return `${bytes} B`
        const kb = bytes / 1024
        if (kb < 1024) return `${kb.toFixed(2)} KB`
        const mb = kb / 1024
        if (mb < 1024) return `${mb.toFixed(2)} MB`
        const gb = mb / 1024
        return `${gb.toFixed(2)} GB`
    }

    function formatDate(unix: number): string {
        const date = new Date(unix * 1000)
        return date.toLocaleDateString(undefined, {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        })
    }

//     const tags = {
//         tf2: '#495544',
//         games: '#954951'
//     }

//     const vidTags = ['tf2', 'games']
</script>

<div class="col-sm-6 col-lg-4 col-xl-3 mb-3">
    <div class="card d-flex flex-column h-100 border-0 shadow-sm card-hover">
        <div class="position-relative aspect-video bg-light overflow-hidden rounded-2">
            {#if video.state == 'processing'}
                <img
                    src={'placeholder.svg'}
                    class="w-100 h-100 object-fit-cover"
                    alt="placeholder" />
                <div
                    class="position-absolute top-0 start-0 w-100 h-100 d-flex align-items-center justify-content-center flex-column opacity-1 gap-1"
                    style="background: rgba(0,0,0,0.5);">
                    <i class="spinner-border text-white" role="status"></i>
                </div>
                <div class="position-absolute bottom-0 end-0 m-2">
                    <span class="badge bg-dark">Processing...</span>
                </div>
            {:else}
                <img
                    src={`${PUBLIC_CDN_URL}/${video.thumb_key}`}
                    alt={video.name}
                    class="w-100 h-100"
                    style="object-fit: contain;" />
                <div
                    class="position-absolute top-0 start-0 w-100 h-100 d-flex align-items-center justify-content-center opacity-0 hover-overlay"
                    style="background: rgba(0,0,0,0.3);">
                    <button
                        class="btn bg-black btn-sm text-white"
                        onclick={() => currentVideoURL.set(`${PUBLIC_CDN_URL}/${video.file_key}`)}
                        aria-label="Play video">
                        <i class="bi bi-play-fill me-1"></i>Play
                    </button>
                </div>
                <div class="position-absolute bottom-0 end-0 m-2">
                    <span class="badge bg-black bg-opacity-75"
                        >{formatDuration(video.duration)}</span>
                </div>
            {/if}
        </div>

        <div class="card-body p-3">
            <div class="d-flex align-items-start justify-content-between mb-2">
                <h6 class="card-title flex-grow-1 overflow-hidden pe-3">
                    {video.name.substr(0, video.name.length-4)}
                </h6>
                <VideoDropdown fileKey={video.file_key} />
            </div>

            <div class="d-flex align-items-center justify-content-between small text-muted">
                <span>{formatDate(video.created_at)}</span>
                <span class="badge bg-dark p-2">{video.format.toUpperCase().slice(6)}</span>
            </div>

            <div class="d-flex align-items-center justify-content-between small text-muted">
                <span>{formatSize(video.size)}</span>
                <span>{video.views} views</span>
            </div>

            <!-- {#if vidTags.length > 0}
                <div class="mt-3 d-flex flex-wrap gap-1" id="tags">
                    {#each vidTags as t}
                        <span
                            class="badge rounded-pill d-flex align-items-center gap-1"
                            style="background-color: {tags[
                                t
                            ]}; color: #fff; font-size: 0.75rem; padding: 0.25em 0.5em;">
                            <i class="bi bi-circle-fill" style="font-size: 0.5rem;"></i>
                            {t}
                        </span>
                    {/each}
                </div>
            {/if} -->
        </div>
    </div>
</div>

<style>
    .hover-overlay {
        transition: opacity 0.1s ease;
    }

    .card-hover:hover .hover-overlay {
        opacity: 1 !important;
    }
</style>
