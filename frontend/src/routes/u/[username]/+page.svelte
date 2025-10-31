<script lang="ts">
    import { page } from '$app/state'
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import Player from '$lib/components/dashboard/Player.svelte'
    import Header from '$lib/components/Header.svelte'
    import Card from '$lib/components/video/Card.svelte'
    import { dashboardView } from '$lib/stores/appControl'
    import { currentVideoURL } from '$lib/stores/VideoStore'
    import { onDestroy } from 'svelte'
    import type { PageProps } from './$types'

    let { data }: { data: PageProps } = $props()

    const username = page.params.username

    onDestroy(() => {
        currentVideoURL.set('')
    })
</script>

<svelte:head>
    {#if !(data as any).found}
        <title>Profile not found - Vid.sh</title>
        <meta property="og:title" content="Profile not found - Vid.sh" />
        <meta property="og:description" content={`Either their profile is private, was deleted, or never exited`} />
        <meta property="theme-color" content="#5733E7" />
    {/if}

    {#if (data as any).found}
        <title>{username} on Vid.sh</title>
        <meta property="og:title" content={`@${username}`} />
        <meta property="og:image" content={`${PUBLIC_CDN_URL}/avatars/${data.avatarHash}`} />
        <meta property="og:description" content={`View ${data.videos?.length || 0} videos posted by them on vid.sh`} />
        <meta property="theme-color" content="#5733E7" />
    {/if}
</svelte:head>

{#if $currentVideoURL != ''}
    <Player />
{/if}

{#if !(data as any).found}
    <Header title="Profile Not Found" />
    <div class="container d-flex justify-content-center align-items-center" style="height: 80vh;">
        <div class="card text-center shadow-lg border-0 p-4 bg-body-tertiary" style="max-width: 420px;">
            <div class="card-body">
                <div class="mb-3">
                    <i class="bi bi-person-x display-3 text-danger"></i>
                </div>
                <h4 class="fw-semibold mb-2">Profile Not Found</h4>
                <p class="text-muted mb-4">The user profile you're looking for doesn't exist or may have been removed.</p>
                <a href="/" class="btn btn-primary px-4">Go Home</a>
            </div>
        </div>
    </div>
{:else}
    <Header title={`${username}'s profile`} />
    <div class="container">
        <div class="row position-relative justify-content-center">
            <div class="banner position-relative mt-3" style="width: 900px; height: 150px; z-index: 1; border-radius: 8px;">
                <!-- <img
                    alt="profile_banner"
                    class="h-100 w-auto d-block mx-auto rounded-3 shadow-sm"
                    style="object-fit: cover; object-position: center;"
                    src="https://cdn.discordapp.com/attachments/1311343659773329450/1420179841621164082/banner.png?ex=68fb594b&is=68fa07cb&hm=40f3d9e70d27cf2a58f19d0aa481410937e428669f1cbfeb0ba87d1d69497be3&" /> -->

                <div class="position-absolute start-50 translate-middle-x z-2">
                    <img
                        alt="profile_picture"
                        width="128"
                        height="128"
                        class="rounded-circle object-fit-cover border border-3 border-white shadow"
                        src={`${PUBLIC_CDN_URL}/avatars/${data.avatarHash}`}
                        style="background: #fff;" />
                </div>
            </div>

            <h1 class="text-center fs-4 mt-1">{username}</h1>
        </div>
        <!-- <div class="row mb-1">
            <Search tags={[]}></Search>
        </div> -->
        <div class="row">
            <p class="text-center small text-muted">Only showing up to 25 videos for now. Profiles are very experimental</p>
            {#if data.videos.length === 0}
                <div class="py-5 text-center w-100">
                    <div class="d-flex align-items-center justify-content-center rounded-circle bg-body-tertiary mx-auto mb-4" style="width: 96px; height: 96px;">
                        <i class="bi bi-file-earmark-play text-muted display-4"></i>
                    </div>
                    <h4 class="fw-semibold mb-3">No videos yet</h4>
                    <p class="text-muted mb-4">This user hasn't uploaded any videos yet.</p>
                </div>
            {/if}

            {#if $dashboardView === 'list'}
                {#each data.videos as vid, i (vid.file_key)}
                    <Card
                        {i}
                        video={{
                            id: 1,
                            name: vid.name,
                            format: 'video/mp4',
                            created_at: vid.created_at,
                            file_key: vid.file_key,
                            duration: vid.duration,
                            thumbnail_url: vid.thumbnail_url,
                            video_url: vid.video_url,
                            size: vid.size
                        } as any}
                        isProfile={true} />
                {/each}
            {:else}
                {#each data.videos as vid, i (vid.file_key)}
                    <Card
                        {i}
                        video={{
                            id: 1,
                            name: vid.name,
                            format: 'video/mp4',
                            created_at: vid.created_at,
                            file_key: vid.file_key,
                            duration: vid.duration,
                            thumbnail_url: vid.thumbnail_url,
                            video_url: vid.video_url,
                            size: vid.size
                        } as any}
                        isProfile={true} />
                {/each}
            {/if}
        </div>
    </div>
{/if}
