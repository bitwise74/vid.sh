<script lang="ts">
    import { type User } from '$lib/api/User'
    import ToastContainer from '$lib/components/toast/ToastContainer.svelte'
    import { dashboardView } from '$lib/stores/appControl'
    import { isLoggedIn, loadedVideosCount, user } from '$lib/stores/AppVars'
    import { videos } from '$lib/stores/VideoStore'
    import { onMount } from 'svelte'
    import '../app.css'

    const { data, children }: any = $props()

    $effect(() => {
        if (data.loggedIn) {
            user.set(data)
            videos.set(data.videos)
            isLoggedIn.set(true)
            loadedVideosCount.set(data.videos?.length ?? 0)
            return
        }

        isLoggedIn.set(false)
        loadedVideosCount.set(0)
        videos.set([])
    })

    onMount(() => {
        dashboardView.set(localStorage.getItem('view') || 'grid')
    })
</script>

{@render children?.()}

<ToastContainer />
