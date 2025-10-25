<script lang="ts">
    import { currentVideoURL, videos } from '$lib/stores/VideoStore'
    import Card from '../video/Card.svelte'
    import Player from './Player.svelte'

    let { isProfile = false } = $props()
</script>

{#if $currentVideoURL !== ''}
    <Player />
{/if}

<div class="row">
    {#if $videos.length === 0}
        <div class="py-5 text-center">
            <div class="d-flex align-items-center justify-content-center rounded-circle bg-body-tertiary mx-auto mb-4" style="width: 96px; height: 96px;">
                <i class="bi bi-file-earmark-play text-muted display-4"></i>
            </div>
            <h4 class="fw-semibold mb-3">No videos yet</h4>
            <p class="text-muted mb-4">Upload your first video by drag and dropping it here or going to the editor</p>
            <a href="/editor" class="btn btn-dark">Open Editor</a>
        </div>
    {:else}
        {#each $videos as video, i (video.id)}
            <Card {video} {isProfile} {i} />
        {/each}
    {/if}
</div>
