<script lang="ts">
    import { goto } from '$app/navigation'
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import { GetUser } from '$lib/api-v2/User'
    import ToastContainer from '$lib/components/toast/ToastContainer.svelte'
    import { isLoggedIn, loadedVideosCount, shouldRefetch, user } from '$lib/stores/AppVars'
    import { toastStore } from '$lib/stores/ToastStore'
    import { view } from '$lib/stores/UserPreferences'
    import { videos } from '$lib/stores/VideoStore'
    import { getCookie } from '$lib/utils/Cookies'
    import { onMount } from 'svelte'
    import '../app.css'

    const { children } = $props()
    let isLoading = $state(true)

    async function loadData() {
        try {
            if (getCookie('logged_in') !== '1') return

            const data = await GetUser(fetch)
            if (data) {
                const vids = data.videos || []

                for (let i = 0; i < vids.length; i++) {
                    const v = vids[i]

                    vids[i].thumbnail_url = `${PUBLIC_CDN_URL}/${v.file_key.split('.')[0]}.webp`
                    vids[i].video_url = `${PUBLIC_CDN_URL}/${v.file_key}${v.version > 1 ? `?v=${v.version}` : ''}`
                }

                user.set(data)
                videos.set(vids)
                isLoggedIn.set(true)
                loadedVideosCount.set(data.videos?.length ?? 0)
                view.set(localStorage.getItem('view') || 'list')
            } else {
                toastStore.error({
                    title: 'Session expired',
                    message: 'Please log in again',
                    duration: 10000
                })
                goto('/login')
            }
        } catch (err) {
            console.error(err)
        } finally {
            isLoading = false
        }
    }

    $effect(() => {
        if ($shouldRefetch) {
            loadData().finally(() => {
                shouldRefetch.set(false)
            })
        }
    })
</script>

{#if isLoading}
    <div></div>
{:else}
    {@render children?.()}
{/if}

<ToastContainer />
