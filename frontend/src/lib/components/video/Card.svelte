<script lang="ts">
    import type { Video } from '$lib/api-v2/Files'
    import { selectedVideos } from '$lib/stores/Dashboard'
    import { view } from '$lib/stores/UserPreferences'
    import { currentVideoURL } from '$lib/stores/VideoStore'
    import { FormatDate, FormatDuration, FormatSize } from '$lib/utils/Format'
    import VideoDropdown from './Dropdown.svelte'

    type Props = {
        video: Video
        isProfile?: boolean
        i: number
    }

    let { video, isProfile = false, i }: Props = $props()

    function handleOnTick(e: Event) {
        const input = e.target as HTMLInputElement
        if (input.checked) {
            selectedVideos.update((ids) => {
                if (!ids.includes(video.id)) {
                    ids.push(video.id)
                }
                return ids
            })
        } else {
            selectedVideos.update((ids) => ids.filter((id) => id !== video.id))
        }
    }
</script>

<!--
 Very hacky fix for the dropdown rendering in the wrong spot:
 remove the animation class after animation end to reset the stacking context.
 You do what u gotta do I guess.
-->

{#if $view === 'list'}
    <div class="row justify-content-center card-animate" style="animation-delay: {i * 0.025}s;" onanimationend={(e) => e.currentTarget.classList.remove('card-animate')}>
        <div class="col-xl-10">
            <div class="card card-hover mb-1 flex-row align-items-center shadow-bottom bg-body-tertiary border-0 rounded-3">
                <input type="checkbox" id="video-select" class="form-check-input m-3" aria-label="Select video" onchange={(e) => handleOnTick(e)} checked={$selectedVideos.includes(video.id)} />
                <div class="position-relative overflow-hidden">
                    {#if video.state === 'processing'}
                        <div class="position-relative rounded overflow-hidden" style="width:150px; height:92px;">
                            <img src="placeholder.svg" alt="processing" class="w-100 h-100 object-fit-cover bg-body" />
                            <div class="position-absolute top-0 start-0 w-100 h-100 d-flex align-items-center justify-content-center bg-dark bg-opacity-50">
                                <div class="spinner-border text-white spinner-border-sm" role="status"></div>
                            </div>
                            <span class="position-absolute bottom-0 end-0 m-1 badge bg-dark small">Processing...</span>
                        </div>
                    {:else}
                        <div class="position-relative rounded overflow-hidden thumbnail-wrap" style="width:150px; height:92px;">
                            <img src={video.thumbnail_url} alt={video.name} class="w-100 h-100 object-fit-cover bg-body" />
                            <div class="hover-overlay position-absolute opacity-0 top-0 start-0 w-100 h-100 d-flex align-items-center justify-content-center" style="background: rgba(0,0,0,0.3);">
                                <button class="btn btn-sm bg-black text-white" onclick={() => currentVideoURL.set(video.video_url!)} aria-label="Play video">
                                    <i class="bi-play-fill me-1"></i>Play
                                </button>
                            </div>
                            <span class="position-absolute bottom-0 end-0 m-1 badge bg-black bg-opacity-75 small">
                                {FormatDuration(video.duration)}
                            </span>
                        </div>
                    {/if}
                </div>

                <div class="card-body d-flex flex-column flex-grow-1 p-3" style="min-width: 0; flex: 1;">
                    <div class="d-flex align-items-center mb-2">
                        <h6 class="card-title flex-grow-1 text-truncate pe-3 mb-0">
                            {video.name.slice(0, -4)}
                        </h6>
                        <VideoDropdown {video} {isProfile} />
                    </div>

                    {#if !isProfile}
                        <div class="d-flex align-items-center gap-3 small text-muted mb-1">
                            <span>{FormatSize(video.size)}</span>
                            <span class="d-sm-flex d-none">{FormatDate(video.created_at)}</span>
                            <span class="badge bg-dark d-sm-flex d-none">{video.format.toUpperCase().slice(6)}</span>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
{:else}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <div class="col-sm-6 col-12 col-lg-4 col-xl-3 mb-3 card-animate" style="animation-delay: {i * 0.02}s;" onanimationend={(e) => e.currentTarget.classList.remove('card-animate')}>
        <div class="card d-flex flex-column h-100 card-hover border-0 shadow-sm bg-body-tertiary">
            <div class="position-relative bg-black rounded-2 aspect-video overflow-hidden">
                {#if video.state == 'processing'}
                    <img src={'placeholder.svg'} class="w-100 h-100 object-fit-cover" alt="placeholder" />
                    <div class="position-absolute w-100 h-100 d-flex align-items-center justify-content-center flex-column opacity-1 start-0 top-0 gap-1" style="background: rgba(0,0,0,0.5);">
                        <i class="spinner-border text-white" role="status"></i>
                    </div>
                    <div class="position-absolute bottom-0 end-0 m-2">
                        <span class="badge bg-dark">Processing...</span>
                    </div>
                {:else}
                    <img src={video.thumbnail_url} alt={video.name} class="w-100 h-100" style="object-fit: fill;" />
                    <div class="position-absolute w-100 h-100 d-flex align-items-center justify-content-center hover-overlay start-0 top-0 opacity-0" style="background: rgba(0,0,0,0.3);">
                        <button class="btn btn-sm bg-black text-white" onclick={() => currentVideoURL.set(video.video_url!)} aria-label="Play video">
                            <i class="bi bi-play-fill me-1"></i>Play
                        </button>
                    </div>
                    <div class="position-absolute bottom-0 end-0 m-2">
                        <span class="badge bg-black bg-opacity-75">{FormatDuration(video.duration)}</span>
                    </div>
                {/if}
            </div>

            <div class="card-body p-3 pb-3">
                <div class="d-flex align-items-start justify-content-between mb-2">
                    <h6 class="card-title flex-grow-1 overflow-hidden pe-3">
                        {video.name.substr(0, video.name.length - 4)}
                    </h6>
                    <VideoDropdown {video} {isProfile} />
                </div>
                {#if !isProfile}
                    <div class="d-flex align-items-center justify-content-between small text-muted">
                        <span>{FormatDate(video.created_at)}</span>
                        <span class="badge bg-dark p-2">{video.format.toUpperCase().slice(6)}</span>
                    </div>

                    <div class="d-flex align-items-center justify-content-between small text-muted">
                        <span>{FormatSize(video.size)}</span>
                    </div>
                {/if}
            </div>
        </div>
    </div>
{/if}

<style>
    .hover-overlay {
        transition: opacity 0.2s ease;
    }

    .card-hover:hover .hover-overlay {
        opacity: 1 !important;
    }

    #video-select {
        transform: scale(1.35);
        cursor: pointer;
    }

    .card-animate {
        opacity: 0;
        animation: slidefade-in 0.3s forwards;
    }
</style>
