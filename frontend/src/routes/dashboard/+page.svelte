<script lang="ts">
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import { FetchFiles, UploadFile, type Video } from '$lib/api/Files'
    import VideoList from '$lib/components/dashboard/List.svelte'
    import Search from '$lib/components/dashboard/Search.svelte'
    import StatBlocks from '$lib/components/dashboard/StatBlocks.svelte'
    import Header from '$lib/components/Header.svelte'
    import { videoSelectionToast } from '$lib/components/toast/SelectionToast'
    import { loadedVideosCount, user } from '$lib/stores/AppVars'
    import { toastStore } from '$lib/stores/ToastStore'
    import { videos } from '$lib/stores/VideoStore'
    import { onDestroy, onMount } from 'svelte'

    let page = 0
    let isLoading = false
    let perPage = $state('10')
    let sortBy = $state('newest')

    let dropOverlay: HTMLElement | null = null
    let sentinel: Element | null = null
    let observer: IntersectionObserver | null = null
    let statsCollapsed = $state(false)

    const allLoaded = () => {
        // TODO: find way to not spam requests
        if (!$user || !$user.stats) return false
        return $loadedVideosCount >= $user.stats.uploadedFiles
    }

    function showOverlay(e: DragEvent) {
        if (e.dataTransfer?.types.includes('Files') && dropOverlay) {
            dropOverlay.classList.remove('d-none')
            dropOverlay.classList.add('d-flex')
        }
    }

    function hideOverlay() {
        if (!dropOverlay) return
        dropOverlay.classList.add('d-none')
        dropOverlay.classList.remove('d-flex')
    }

    onMount(async () => {
        if (typeof window === 'undefined') return

        dropOverlay = document.getElementById('dropOverlay')
        window.addEventListener('dragenter', showOverlay)
        window.addEventListener('dragend', hideOverlay)

        if (sentinel) {
            observer = new IntersectionObserver((entries) => {
                // Clean up
                if (allLoaded()) {
                    observer?.disconnect()

                    if (sentinel?.parentNode) {
                        sentinel.parentNode.removeChild(sentinel)
                    }

                    return
                }

                if (entries[0].isIntersecting) {
                    loadModeContent()
                }
            })
            observer.observe(sentinel)
        }

        // Check for missing data
        if ($user.avatarHash == '' && localStorage.getItem('notifNoAvatar') !== 'true') {
            toastStore.info({
                title: 'You can set a profile picture now',
                message: "Go to settings for details (this won't show up again)",
                duration: 10000
            })
            localStorage.setItem('notifNoAvatar', 'true')
        }

        if ($user.username == '' && localStorage.getItem('notifNoUsername') !== 'true') {
            setTimeout(() => {
                toastStore.info({
                    title: "You don't have a nickname set yet",
                    message: "Go to settings for details (this won't show up again)",
                    duration: 10000
                })
            }, 3500)
            localStorage.setItem('notifNoUsername', 'true')
        }
    })

    onDestroy(() => {
        if (typeof window !== 'undefined') {
            window.removeEventListener('dragenter', showOverlay)
            window.removeEventListener('dragend', hideOverlay)
        }
        observer?.disconnect()
    })

    async function handleDrop(e: DragEvent) {
        hideOverlay()
        e.preventDefault()

        if (!e.dataTransfer || !videos) return

        const files = Array.from(e.dataTransfer.files)
        const videoFile = files.find((f) => ['video/mp4', 'video/quicktime', 'video/x-matroska'].includes(f.type))

        if (!videoFile) {
            toastStore.error({
                title: 'No valid files detected',
                message: 'Please use one of the supported formats (mp4, mov, mkv)',
                duration: 10000
            })
            return
        }

        if (videoFile.type === 'video/x-matroska') {
            toastStore.info({
                title: '.mkv file detected',
                message: 'These files usually take longer to process',
                duration: 10000
            })
        }

        videos.set([
            {
                name: videoFile.name,
                size: videoFile.size,
                format: 'video/mp4',
                created_at: Date.now() / 1000,
                state: 'processing'
            } as unknown as Video,
            ...$videos
        ])

        try {
            const newVid = await UploadFile(videoFile)
            if (!newVid) return

            videos.set([newVid, ...$videos.splice(1)])
        } catch (error) {
            // Remove processing vid if failed
            videos.set([...$videos.splice(1)])
            toastStore.error({
                title: 'Failed to save video to cloud',
                message: error.message,
                duration: 10000
            })
        }
    }

    async function loadModeContent() {
        if (isLoading) return
        isLoading = true

        page++

        try {
            const newVideos = await FetchFiles({
                limit: parseInt(perPage),
                page: page,
                sort: sortBy as any, // TODO: fix this ugly
                tags: ''
            })

            for (const vid of newVideos) {
                vid.thumbnail_url = `${PUBLIC_CDN_URL}/${vid.file_key.replace('.mp4', '.webp')}`
                vid.video_url = `${PUBLIC_CDN_URL}/${vid.file_key}`
            }

            loadedVideosCount.set($loadedVideosCount + newVideos.length)

            // TODO: should append based on sort order
            videos.set([...$videos, ...newVideos])
        } catch (error) {
            toastStore.error({
                title: 'Failed to load more content',
                message: error.message
            })
        } finally {
            isLoading = false
        }
    }

    $effect(() => videoSelectionToast.subscribe(() => {}))

    function handleStatsMode(event: CustomEvent<boolean>) {
        statsCollapsed = event.detail
    }
</script>

<svelte:head>
    <title>Dashboard - vid.sh</title>
</svelte:head>

<div class="dashboard-root min-vh-100 position-relative gradient-bg-dark">
    <div
        class="drop-overlay position-fixed w-100 h-100 d-none justify-content-center align-items-center z-3 start-0 top-0"
        id="dropOverlay"
        role="none"
        ondragenter={showOverlay}
        ondragover={(e) => e.preventDefault()}
        ondragend={hideOverlay}
        ondragleave={hideOverlay}
        ondrop={handleDrop}>
        <div class="drop-overlay__panel">
            <i class="bi bi-cloud-upload mb-3"></i>
            <h3>Drop files to upload</h3>
            <p>MP4, MOV, or MKV up to your plan limits</p>
        </div>
    </div>

    <Header title="Dashboard" page="dashboard" />
    <main class="dashboard-main container py-4 py-md-5">
        <div class={`dashboard-stack ${statsCollapsed ? 'stack-collapsed' : ''}`}>
            <div class="dashboard-stack__stats">
                <StatBlocks on:modeChange={handleStatsMode} />
            </div>
            <div class="dashboard-stack__search">
                <Search />
            </div>
        </div>

        <section class="video-panel glass-panel">
            <div class="video-panel__header">
                <div>
                    <p class="eyebrow">Your uploads</p>
                    <h3 class="section-title">Browse and manage</h3>
                </div>
                <span class="section-hint">{$loadedVideosCount} / {$user.stats?.uploadedFiles || 0} videos</span>
            </div>
            <VideoList />
        </section>

        <div bind:this={sentinel} class="infinite-sentinel"></div>
    </main>

    <div class="text-center mb-4">
        {#if $user && $user.stats && $user.stats.uploadedFiles != 0 && $user.stats.uploadedFiles > parseInt(perPage)}
            {#if !allLoaded()}
                <p class="text-muted small">Scroll down to load more</p>
            {/if}
        {/if}
    </div>
</div>

<style>
    .dashboard-root {
        padding-bottom: 4rem;
    }

    .dashboard-main {
        display: flex;
        flex-direction: column;
        gap: 2rem;
    }

    .dashboard-stack {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }

    .dashboard-stack__stats,
    .dashboard-stack__search {
        width: 100%;
    }

    .glass-panel {
        border: 1px solid var(--dashboard-surface-border);
        border-radius: 2rem;
        padding: 2rem;
        background: var(--dashboard-surface-bg);
        box-shadow: var(--dashboard-surface-shadow);
        color: var(--dashboard-text-primary);
    }

    .video-panel__header {
        display: flex;
        justify-content: space-between;
        gap: 1rem;
        align-items: baseline;
        margin-bottom: 1.25rem;
    }

    .eyebrow {
        text-transform: uppercase;
        letter-spacing: 0.25em;
        font-size: 0.7rem;
        color: var(--dashboard-text-muted);
        margin-bottom: 0.35rem;
    }

    .section-title {
        margin: 0;
        font-size: 1.5rem;
    }

    .section-hint {
        color: var(--dashboard-text-muted);
        font-size: 0.9rem;
    }

    .drop-overlay {
        background: var(--dashboard-overlay);
        backdrop-filter: blur(8px);
    }

    .drop-overlay__panel {
        border: 1.5px dashed var(--dashboard-overlay-panel-border);
        border-radius: 1.5rem;
        padding: 2.5rem 3rem;
        text-align: center;
        color: var(--dashboard-text-primary);
        background: var(--dashboard-overlay-panel-bg);
        box-shadow: var(--dashboard-surface-shadow);
    }

    .drop-overlay__panel i {
        font-size: 3rem;
        color: #7f5dff;
    }

    .drop-overlay__panel p {
        margin: 0;
        color: var(--dashboard-text-muted);
    }

    .infinite-sentinel {
        height: 1px;
    }

    @media (min-width: 992px) {
        .dashboard-stack.stack-collapsed {
            flex-direction: row;
            align-items: flex-start;
            gap: 1.5rem;
        }

        .dashboard-stack.stack-collapsed .dashboard-stack__stats {
            flex: 0 0 360px;
        }

        .dashboard-stack.stack-collapsed .dashboard-stack__search {
            flex: 1;
        }
    }

    @media (max-width: 768px) {
        .glass-panel {
            padding: 1.5rem;
            border-radius: 1.5rem;
        }

        .video-panel__header {
            flex-direction: column;
            align-items: flex-start;
        }
    }
</style>
