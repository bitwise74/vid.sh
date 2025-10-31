<script lang="ts">
    import { type User } from '$lib/api/User'
    import ToastContainer from '$lib/components/toast/ToastContainer.svelte'
    import { dashboardView } from '$lib/stores/appControl'
    import { isLoggedIn, loadedVideosCount, user } from '$lib/stores/AppVars'
    import { videos } from '$lib/stores/VideoStore'
    import { onMount } from 'svelte'
    import '../app.css'

    const { children, data }: { children: any; data: User | null } = $props()

    if (data) {
        user.set(data)
        videos.set(data.videos)

        isLoggedIn.set(true)
        loadedVideosCount.set(data.videos?.length ?? 0)
    }

    onMount(() => {
        dashboardView.set(localStorage.getItem('view') || 'grid')
    })
</script>

{@render children?.()}

<ToastContainer />
