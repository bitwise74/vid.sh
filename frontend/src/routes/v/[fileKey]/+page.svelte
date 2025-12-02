<script lang="ts">
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import type { Video } from '$lib/api/Files'
    import type { PageProps } from './$types'

    type FetchedVideo = {
        file?: Video
        user: {
            avatarURL: string
            username: string
        }
    }

    let { data }: PageProps = $props()
    let d = data.video as FetchedVideo
</script>

<svelte:head>
    {#if d?.file === undefined}
        <meta property="og:title" content="No such file" />
        <meta property="og:description" content="Either the file doesn't exist, was deleted, or is private" />
        <meta property="og:url" content="https://bitwise0x.dev" />
        <meta property="theme-color" content="#5733E7" />
    {:else}
        <meta property="og:title" content={d.file.name.replace('.mp4', '')} />

        <meta property="og:type" content="video.other" />
        <meta property="og:site_name" content={`Video by @${d.user.username}`} />

        <meta property="og:video:type" content={d.file.format} />
        <meta property="og:video" content={`${PUBLIC_CDN_URL}/${d.file.file_key}`} />
        <meta property="og:video:secure_url" content={`${PUBLIC_CDN_URL}/${d.file.file_key}`} />
        <meta property="theme-color" content="#5733E7" />

        <!-- <link type="application/json+oembed" href={`${PUBLIC_BASE_URL}/api/oembed?username=${d.user.username}&size=${d.file.size}`} /> -->
    {/if}
</svelte:head>

{#if d?.file === undefined}
    <h1>No such file</h1>
    <p>Either the file doesn't exist, was deleted, or is private</p>
{:else}
    <video autoplay controls src={`${PUBLIC_CDN_URL}/${d.file.file_key}`}> <track kind="captions" /></video>
{/if}

<style>
    :global(body) {
        margin: 0;
        background: black;
        color: white;
    }

    video {
        position: absolute;
        inset: 0;
        max-width: 100%;
        max-height: 100%;
        margin: auto;
        user-select: none;
    }
</style>
