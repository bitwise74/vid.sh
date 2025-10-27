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
</script>

<svelte:head>
    <title>Dashboard - vid.sh</title>
</svelte:head>

<div class="min-vh-100 position-relative">
    <div
        class="position-fixed w-100 h-100 bg-dark d-none justify-content-center align-items-center z-3
           start-0 top-0 bg-opacity-50"
        id="dropOverlay"
        role="none"
        ondragenter={showOverlay}
        ondragover={(e) => e.preventDefault()}
        ondragend={hideOverlay}
        ondragleave={hideOverlay}
        ondrop={handleDrop}>
        <div class="border-3 border-light rounded-3 border border-dashed p-5 text-center">
            <i class="bi bi-cloud-upload text-light display-1 mb-3"></i>
            <h3 class="text-light fw-semibold">Drop files to upload</h3>
        </div>
    </div>

    <Header title="Dashboard" page="dashboard" />
    <main class="container py-4">
        <StatBlocks />
        <Search tags={[]} />
        <VideoList />

        <div bind:this={sentinel}></div>
    </main>

    <div class="text-center mb-4">
        {#if $user && $user.stats && $user.stats.uploadedFiles != 0 && $user.stats.uploadedFiles > parseInt(perPage)}
            {#if !allLoaded()}
                <p class="text-muted small">Scroll down to load more</p>
            {:else}
                <p class="text-muted small">All videos loaded</p>
            {/if}
        {/if}
    </div>
</div>
