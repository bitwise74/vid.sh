<script lang="ts">
    import { currentVideoURL, videos } from '$lib/stores/VideoStore'
    import VideoCard from './Card.svelte'
</script>
<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<!-- svelte-ignore a11y_click_events_have_key_events -->
{#if $currentVideoURL !== ''}
    <div class="modal-backdrop show pt-0" onclick={() => currentVideoURL.set('')}></div>
    <div class="modal show d-block pt-0 d-flex" tabindex="-1" onclick={() => currentVideoURL.set('')}>
        <div class="modal-dialog modal-xl py-0 align-content-center justify-content-center" onclick={(e) => e.stopPropagation()}>
            <div class="modal-content">
                <!-- TODO: fix video player -->
                <!-- <VideoPlayerDashhboard video={currentVideoURL} isOpen={true} onClose={() => {console.log("closing")}}></VideoPlayerDashhboard> -->
                <video controls autoplay class="w-100" volume="0.75" src={$currentVideoURL}>
                    <track kind="captions" />
                </video>
            </div>
        </div>
    </div>
{/if}

{#if $videos.length === 0}
    <div class="text-center py-5">
        <div class="d-flex align-items-center justify-content-center rounded-circle bg-light mx-auto mb-4" style="width: 96px; height: 96px;">
            <i class="bi bi-file-earmark-play text-muted display-4"></i>
        </div>
        <h4 class="fw-semibold mb-3">No videos yet</h4>
        <p class="text-muted mb-4">Upload your first video by drag and dropping it here or going to the editor</p>
        <a href="/editor" class="btn btn-dark">Open Editor</a>
    </div>
{:else}
    {#each $videos as video}
        <VideoCard {video} />
    {/each}
{/if}
