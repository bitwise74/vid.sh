<script lang="ts">
    import type { Video } from '$lib/api/Files'
    import { dashboardView, selectedVideos as selected } from '$lib/stores/appControl'
    import { currentVideoURL } from '$lib/stores/VideoStore'
    import { FormatDate, FormatDuration, FormatSize } from '$lib/utils/format'
    import { onMount } from 'svelte'
    import VideoDropdown from './Dropdown.svelte'

    type Props = {
        video: Video
        isProfile?: boolean
        i: number
    }

    let animationSpeed: number = $state(1)

    onMount(() => {
        animationSpeed = parseFloat(localStorage.getItem('optAnimationSpeed') || '1')
    })

    let { video, isProfile = false, i }: Props = $props()

    function handleOnTick(e: Event) {
        const input = e.target as HTMLInputElement
        if (input.checked) {
            selected.update((ids) => {
                if (!ids.includes(video.id)) {
                    ids.push(video.id)
                }
                return ids
            })
        } else {
            selected.update((ids) => ids.filter((id) => id !== video.id))
        }
    }
</script>

<!--
 Very hacky fix for the dropdown rendering in the wrong spot:
 remove the animation class after animation end to reset the stacking context.
 You do what u gotta do I guess.
-->

{#if $dashboardView === 'list'}
    <div class="row justify-content-center card-animate" style="animation-delay:{i * 0.025}s" onanimationend={(e) => e.currentTarget.classList.remove('card-animate')}>
        <div class="col-xl-10">
            <div class="card d-flex flex-row align-items-center border-0 rounded-3 bg-body-tertiary shadow-sm mb-1">
                <input type="checkbox" class="form-check-input m-3" aria-label="Select video" checked={$selected.includes(video.id)} onchange={handleOnTick} />

                <div class="thumb position-relative rounded overflow-hidden" style="width:150px;height:92px;">
                    {#if video.state === 'processing'}
                        <img src="placeholder.svg" alt="processing" class="w-100 h-100 object-fit-cover bg-body" />
                        <div class="overlay d-flex align-items-center justify-content-center">
                            <div class="spinner-border text-white spinner-border" role="status"></div>
                        </div>
                        <span class="badge position-absolute bottom-0 end-0 m-1 bg-dark small">Processing...</span>
                    {:else}
                        <img src={video.thumbnail_url} alt={video.name} class="w-100 h-100 object-fit-cover bg-body" />
                        <div class="hover-overlay d-flex align-items-center justify-content-center">
                            <button class="btn btn-sm bg-black text-white" onclick={() => currentVideoURL.set(video.video_url!)} aria-label="Play video">
                                <i class="bi-play-fill me-1"></i>Play
                            </button>
                        </div>
                        <span class="badge position-absolute bottom-0 end-0 m-1 bg-black bg-opacity-75 small">
                            {FormatDuration(video.duration)}
                        </span>
                    {/if}
                </div>

                <div class="card-body flex-grow-1 p-3 d-flex flex-column">
                    <div class="d-flex align-items-center mb-2">
                        <h6 class="flex-grow-1 text-truncate mb-0 pe-3">{video.name.slice(0, -4)}</h6>
                        <VideoDropdown {video} {isProfile} />
                    </div>
                    <div class="d-flex align-items-center gap-3 small text-muted">
                        <span>{FormatSize(video.size)}</span>
                        <span class="d-sm-flex d-none">{FormatDate(video.created_at)}</span>
                        <span class="badge bg-dark d-sm-flex d-none">{video.format.toUpperCase().slice(6)}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
{:else}
    <div class="col-sm-6 col-12 col-lg-4 col-xl-3 mb-3 card-animate" style="animation-delay:{i * 0.02}s" onanimationend={(e) => e.currentTarget.classList.remove('card-animate')}>
        <div class="card h-100 border-0 rounded-3 bg-body-tertiary shadow-sm">
            <div class="thumb position-relative bg-black rounded-3 aspect-video overflow-hidden">
                {#if video.state == 'processing'}
                    <img src="placeholder.svg" alt="processing" class="w-100 h-100 object-fit-cover" />
                    <div class="overlay d-flex align-items-center justify-content-center flex-column gap-1">
                        <i class="spinner-border-sm text-white" role="status"></i>
                        <span class="badge bg-dark">Processing...</span>
                    </div>
                {:else}
                    <img src={video.thumbnail_url} alt={video.name} class="w-100 h-100 object-fit-cover" />
                    <div class="hover-overlay d-flex align-items-center justify-content-center">
                        <button class="btn btn-sm bg-black text-white" onclick={() => currentVideoURL.set(video.video_url!)} aria-label="Play video">
                            <i class="bi bi-play-fill me-1"></i>Play
                        </button>
                    </div>
                    <span class="badge position-absolute bottom-0 end-0 m-2 bg-black bg-opacity-75">
                        {FormatDuration(video.duration)}
                    </span>
                {/if}
            </div>

            <div class="card-body p-3">
                <div class="d-flex align-items-start justify-content-between mb-2">
                    <h6 class="flex-grow-1 text-truncate pe-3 mb-0">{video.name.slice(0, -4)}</h6>
                    <VideoDropdown {video} {isProfile} />
                </div>

                {#if !isProfile}
                    <div class="d-flex justify-content-between small text-muted">
                        <span>{FormatDate(video.created_at)}</span>
                        <span class="badge bg-dark p-2">{video.format.toUpperCase().slice(6)}</span>
                    </div>
                    <div class="d-flex justify-content-between small text-muted">
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

    .card-animate {
        opacity: 0;
        animation: slidefade-in 0.3s forwards;
    }

    .thumb:hover .hover-overlay {
        opacity: 1;
    }
    .hover-overlay {
        position: absolute;
        inset: 0;
        opacity: 0;
        background: rgba(0, 0, 0, 0.3);
        transition: opacity 0.2s;
    }
    .overlay {
        position: absolute;
        inset: 0;
        background: rgba(0, 0, 0, 0.5);
    }
</style>
