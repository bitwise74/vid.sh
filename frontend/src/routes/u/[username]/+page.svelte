<script lang="ts">
    import { page } from '$app/state'
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import Search from '$lib/components/dashboard/Search.svelte'
    import Header from '$lib/components/Header.svelte'
    import Card from '$lib/components/video/Card.svelte'
    import { view } from '$lib/stores/UserPreferences'
    import type { ProfileData } from './+page.server'

    let {
        data
    }: {
        data: ProfileData
    } = $props()

    const username = page.params.username
</script>

<Header title={`${username}'s profile`} />

{#if !data.public}
    <div class="d-flex justify-content-center align-items-center" style="height: 60vh;">
        <div class="text-center p-5 rounded-4 shadow-sm border bg-body-tertiary">
            <div class="d-flex justify-content-center align-items-center mb-4 bg-dark" style="width: 100px; height: 100px; border-radius: 50%;">
                <i class="bi-lock fs-1 text-white"></i>
            </div>
            <h4 class="fw-semibold mb-0">This profile is private</h4>
        </div>
    </div>
{:else}
    <div class="container">
        <div class="row position-relative justify-content-center">
            <div class="banner position-relative mt-3" style="width: 900px; height: 150px; z-index: 1; border-radius: 8px;">
                <img
                    alt="profile_banner"
                    class="h-100 w-auto d-block mx-auto rounded-3 shadow-sm"
                    style="object-fit: cover; object-position: center;"
                    src="https://cdn.discordapp.com/attachments/1311343659773329450/1420179841621164082/banner.png?ex=68fb594b&is=68fa07cb&hm=40f3d9e70d27cf2a58f19d0aa481410937e428669f1cbfeb0ba87d1d69497be3&" />

                <div class="position-absolute start-50 translate-middle-x" style="bottom: -32px; z-index: 2;">
                    <img
                        alt="profile_picture"
                        width="128"
                        height="128"
                        class="rounded-circle object-fit-cover border border-3 border-white shadow"
                        src={`${PUBLIC_CDN_URL}/avatars/${data.avatarHash}`}
                        style="background: #fff;" />
                </div>
            </div>

            <h1 class="text-center fs-4 mt-5">{username}</h1>
        </div>
        <div class="row mb-1">
            <Search tags={[]}></Search>
        </div>
        <div class="row">
            {#if $view === 'list'}
                {#each data.videos as vid}
                    <Card
                        video={{
                            id: 1,
                            name: vid.name,
                            format: 'video/mp4',
                            created_at: vid.created_at,
                            file_key: vid.file_key,
                            duration: vid.duration,
                            thumbnail_url: vid.thumbnail_url,
                            video_url: vid.video_url
                        } as any}
                        isProfile={true} />
                {/each}
            {:else}
                {#each data.videos as vid}
                    <Card
                        video={{
                            id: 1,
                            name: vid.name,
                            format: 'video/mp4',
                            created_at: vid.created_at,
                            file_key: vid.file_key,
                            duration: vid.duration,
                            thumbnail_url: vid.thumbnail_url,
                            video_url: vid.video_url
                        } as any}
                        isProfile={true} />
                {/each}
            {/if}
        </div>
    </div>
{/if}
