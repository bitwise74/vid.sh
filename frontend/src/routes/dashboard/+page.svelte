<script lang="ts">
    import { goto } from '$app/navigation'
    import { LoadVideos, SearchVideos, UploadVideo, type Video } from '$lib/api/Files'
    import { LoadInitialData } from '$lib/api/User'
    import Header from '$lib/components/Header.svelte'
    import StatBlocks from '$lib/components/StatBlocks.svelte'
    import { toastStore } from '$lib/components/toast/toastStore'
    import VideoList from '$lib/components/video/List.svelte'
    import { stats, videos } from '$lib/stores/VideoStore'
    import { onDestroy, onMount } from 'svelte'

    let page = 0
    // Used to stop requests when everything has been loaded
    let loadedVideos = 0
    let loading = false
    let perPage = $state('10')
    let sortBy = $state('newest')

    let dropOverlay: HTMLElement | null
    let timeout: number | undefined
    let sentinel: Element
    let observer = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting && loadedVideos < $stats.uploadedFiles) {
            loadModeContent()
        }
    })

    onMount(async () => {
        dropOverlay = document.getElementById('dropOverlay')
        window.addEventListener('dragenter', showOverlay)
        window.addEventListener('dragend', hideOverlay)

        try {
            const data = await LoadInitialData()

            videos.set(data.videos)
            stats.set(data.stats)
            loadedVideos = data.videos.length
        } catch (error) {
            toastStore.error({
                title: 'Failed to load dashboard',
                message: error.message,
                duration: 10000
            })
            console.error('GET API/USERS: Failed to load initial data', error)
            goto('/')
        }

        if (sentinel) {
            observer.observe(sentinel)
        }
    })

    onDestroy(() => {
        clearTimeout(timeout)
        observer?.disconnect()
        window.removeEventListener('dragenter', showOverlay)
        window.removeEventListener('dragend', hideOverlay)
    })

    function handleInput(e: any) {
        clearTimeout(timeout)
        const search = e.target.value

        timeout = setTimeout(async () => {
            videos.set(await SearchVideos(search, parseInt(perPage)))
        }, 300)
    }

    function showOverlay(e: any) {
        if (e.dataTransfer.types.includes('Files') && dropOverlay) {
            dropOverlay.classList.remove('d-none')
            dropOverlay.classList.add('d-flex')
        }
    }

    function hideOverlay() {
        if (!dropOverlay) return
        dropOverlay.classList.add('d-none')
        dropOverlay.classList.remove('d-flex')
    }

    async function handleDrop(e: DragEvent) {
        hideOverlay()
        e.preventDefault()

        if (!e.dataTransfer) return

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
            } as Video,
            ...$videos
        ])

        try {
            const newVid = await UploadVideo(videoFile)
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
        if (loading) return
        loading = true

        page++

        try {
            const newVideos = await LoadVideos(page, parseInt(perPage))
            videos.set([...$videos, ...newVideos])
        } catch (error) {
            toastStore.error({
                title: 'Failed to load more content',
                message: error.message
            })
        } finally {
            loading = false
        }
    }
</script>

<svelte:head>
    <title>Dashboard - vid.sh</title>
</svelte:head>

<div class="min-vh-100 bg-light position-relative">
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <div
        class="position-fixed top-0 start-0 w-100 h-100 bg-dark bg-opacity-50 d-none
           justify-content-center align-items-center z-3"
        id="dropOverlay"
        ondragenter={showOverlay}
        ondragover={(e) => e.preventDefault()}
        ondragend={hideOverlay}
        ondragleave={hideOverlay}
        ondrop={handleDrop}>
        <div class="border border-3 border-dashed border-light rounded-4 p-5 text-center">
            <i class="bi bi-cloud-upload text-light display-1 mb-3"></i>
            <h3 class="text-light fw-semibold">Drop files to upload</h3>
        </div>
    </div>

    <Header title="Dashboard" page="dashboard" />
    <main class="container py-4">
        <StatBlocks />

        <div class="row g-3 mb-4">
            <div class="col-lg-8 col-md-7">
                <div class="input-group">
                    <span class="input-group-text">
                        <i class="bi bi-search"></i>
                    </span>
                    <input type="text" class="form-control" placeholder="Search videos..." oninput={handleInput} />
                </div>
            </div>
            <div class="col-lg-3 col-md-3 col-sm-6">
                <select class="form-select" bind:value={sortBy} disabled>
                    <option value="newest">Newest First</option>
                    <option value="oldest">Oldest First</option>
                    <option value="az">Name A-Z</option>
                    <option value="size">Size</option>
                    <option value="views">Most Views</option>
                </select>
            </div>
            <div class="col-lg-1 col-md-2 col-sm-6">
                <select class="form-select" bind:value={perPage} placeholder="Results per page">
                    <option value="10">10</option>
                    <option value="20">20</option>
                    <option value="50">50</option>
                    <option value="100">100</option>
                    <option value="250">250</option>
                </select>
            </div>
        </div>

        <div class="row">
            <VideoList />
        </div>
        <div bind:this={sentinel}></div>
    </main>
</div>
