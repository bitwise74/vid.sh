<script lang="ts">
    import type { Video } from '$lib/api/Files'
    import { dashboardView, selectedVideos as selected } from '$lib/stores/appControl'
    import { currentVideoURL } from '$lib/stores/VideoStore'
    import { FormatDate, FormatDuration, FormatSize } from '$lib/utils/format'
    import { onMount } from 'svelte'
    import VideoDropdown from './Dropdown.svelte'
    import { PUBLIC_BASE_URL } from '$env/static/public'

    type Props = {
        video: Video
        isProfile?: boolean
        animDelay: number
        i: number
    }

    let animationSpeed: number = $state(1)
    let dropdownExpanded = $state(false)
    let dropdownPosition: { x: number; y: number } | null = $state(null)

    let { video, isProfile = false, animDelay }: Props = $props()

    onMount(() => {
        animationSpeed = parseFloat(localStorage.getItem('optAnimationSpeed') || '1')

        const handleWindowClick = () => {
            if (dropdownExpanded) {
                dropdownExpanded = false
            }
        }

        window.addEventListener('click', handleWindowClick)

        return () => {
            window.removeEventListener('click', handleWindowClick)
        }
    })

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

    function updateURL() {
        currentVideoURL.set(video.video_url!)

        if (localStorage.getItem('optRichEmbeds') === 'true') {
            // Overwrite for local testing. Normally it would point to the wrong port
            if (PUBLIC_BASE_URL.includes('localhost')) {
                let base = 'http://localhost:5173'
                window.history.replaceState({}, '', `${base}/v/${video.file_key}`)
            } else {
                window.history.replaceState({}, '', `${PUBLIC_BASE_URL}/v/${video.file_key}`)
            }
        }
    }

    function showContextMenu(event: MouseEvent) {
        event.preventDefault()
        dropdownPosition = { x: event.clientX, y: event.clientY }
        dropdownExpanded = true
    }
</script>

<!--
 Very hacky fix for the dropdown rendering in the wrong spot:
 remove the animation class after animation end to reset the stacking context.
 You do what u gotta do I guess.
-->

{#if $dashboardView === 'list'}
    <div class="row justify-content-center card-animate" style="animation-delay:{animDelay}s" onanimationend={(e) => e.currentTarget.classList.remove('card-animate')}>
        <div class="col-xl-10">
            <div class="video-card video-card--list" role="group" oncontextmenu={showContextMenu}>
                <label class="video-card__checkbox">
                    <input type="checkbox" aria-label="Select video" checked={$selected.includes(video.id)} onchange={handleOnTick} />
                </label>

                <div class="thumb position-relative rounded overflow-hidden">
                    {#if video.state === 'processing'}
                        <img src="placeholder.svg" alt="processing" class="w-100 h-100 object-fit-cover bg-body" />
                        <div class="overlay d-flex align-items-center justify-content-center">
                            <div class="spinner-border text-white spinner-border" role="status"></div>
                        </div>
                        <span class="badge position-absolute bottom-0 end-0 m-1 bg-dark small">Processing...</span>
                    {:else}
                        <img src={video.thumbnail_url} alt={video.name} class="w-100 h-100 object-fit-cover bg-body" />
                        <div class="hover-overlay d-flex align-items-center justify-content-center">
                            <button class="btn btn-sm bg-black text-white" onclick={() => updateURL()} aria-label="Play video">
                                <i class="bi-play-fill me-1"></i>Play
                            </button>
                        </div>
                        <span class="badge position-absolute bottom-0 end-0 m-1 bg-black bg-opacity-75 small">
                            {FormatDuration(video.duration)}
                        </span>
                    {/if}
                </div>

                <div class="video-card__body">
                    <div class="video-card__title-row">
                        <h6 class="video-card__title">{video.name.slice(0, -4)}</h6>
                        <VideoDropdown {video} {isProfile} expanded={dropdownExpanded} position={dropdownPosition} />
                    </div>
                    <div class="video-card__meta">
                        <span>{FormatSize(video.size)}</span>
                        <span class="d-none d-md-inline">{FormatDate(video.created_at)}</span>
                        {#if video.private}
                            <span class="status-pill">Private</span>
                        {/if}
                    </div>
                </div>
            </div>
        </div>
    </div>
    <!-- svelte-ignore a11y_no_static_element_interactions -->
{:else}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="col-sm-6 col-12 col-lg-4 col-xl-3 mb-2 card-animate" style="animation-delay:{animDelay}s" onanimationend={(e) => e.currentTarget.classList.remove('card-animate')}>
        <div class="video-card video-card--grid h-100" role="group" oncontextmenu={showContextMenu}>
            <div class="thumb position-relative rounded-3 aspect-video overflow-hidden">
                {#if video.state == 'processing'}
                    <img src="placeholder.svg" alt="processing" class="w-100 h-100 object-fit-cover" />
                    <div class="overlay d-flex align-items-center justify-content-center flex-column gap-1">
                        <i class="spinner-border-sm text-white" role="status"></i>
                        <span class="badge bg-dark">Processing...</span>
                    </div>
                {:else}
                    <img src={video.thumbnail_url} alt={video.name} class="w-100 h-100 object-fit-cover" />
                    <div class="hover-overlay d-flex align-items-center justify-content-center">
                        <button class="btn btn-sm bg-black text-white" onclick={() => updateURL()} aria-label="Play video">
                            <i class="bi bi-play-fill me-1"></i>Play
                        </button>
                    </div>
                    <span class="badge position-absolute bottom-0 end-0 m-2 bg-black bg-opacity-75">
                        {FormatDuration(video.duration)}
                    </span>
                {/if}
            </div>

            <div class="video-card__body">
                <div class="video-card__title-row">
                    <h6 class="video-card__title">{video.name.slice(0, -4)}</h6>
                    <VideoDropdown {video} {isProfile} expanded={dropdownExpanded} position={dropdownPosition} />
                </div>

                <div class="video-card__meta">
                    <span>{FormatDate(video.created_at)}</span>
                    <span>{FormatSize(video.size)}</span>
                    {#if video.private}
                        <span class="status-pill">Private</span>
                    {/if}
                </div>
            </div>
        </div>
    </div>
{/if}

<style>
    .video-card {
        background: var(--dashboard-card-bg);
        border: 1px solid var(--dashboard-card-border);
        border-radius: 1.25rem;
        padding: 1rem;
        display: flex;
        flex-direction: column;
        gap: 0.9rem;
        box-shadow: var(--dashboard-surface-shadow);
        transition:
            border-color 0.2s ease,
            transform 0.2s ease,
            box-shadow 0.2s ease;
        color: var(--dashboard-text-primary);
    }

    .video-card--list {
        flex-direction: row;
        align-items: stretch;
        gap: 1.25rem;
        padding: 1rem 1.5rem;
    }

    .video-card--grid {
        height: 100%;
    }

    .video-card:hover {
        border-color: var(--dashboard-surface-border);
        box-shadow: 0 30px 80px rgba(0, 0, 0, 0.15);
    }

    .video-card__checkbox {
        display: flex;
        align-items: center;
        margin-right: 0.5rem;
    }

    .video-card__checkbox input {
        appearance: none;
        width: 20px;
        height: 20px;
        border-radius: 6px;
        border: 1px solid var(--dashboard-chip-border);
        background: var(--dashboard-card-bg);
        margin: 0;
        position: relative;
        cursor: pointer;
    }

    .video-card__checkbox input:checked {
        background: linear-gradient(135deg, #7f00ff, #007fff);
        border-color: transparent;
    }

    .video-card__checkbox input:checked::after {
        content: '\2713';
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -58%);
        font-size: 0.85rem;
        color: #fff;
    }

    .video-card--list .thumb {
        width: 150px;
        height: 92px;
        flex-shrink: 0;
    }

    .video-card__body {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .video-card--grid .video-card__body {
        flex-grow: 1;
    }

    .video-card__title-row {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .video-card__title {
        flex: 1;
        margin: 0;
        font-size: 1rem;
        font-weight: 600;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .video-card__meta {
        display: flex;
        flex-wrap: wrap;
        gap: 0.75rem;
        font-size: 0.85rem;
        color: var(--dashboard-text-muted);
        align-items: center;
    }

    .status-pill {
        padding: 0.1rem 0.6rem;
        border-radius: 999px;
        background: rgba(255, 82, 82, 0.15);
        color: #d14343;
        font-weight: 600;
        font-size: 0.75rem;
        text-transform: uppercase;
        letter-spacing: 0.04em;
    }

    .hover-overlay {
        transition: opacity 0.2s ease;
    }

    .thumb {
        border-radius: 1rem;
        overflow: hidden;
        background: var(--dashboard-chip-bg);
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
        background: var(--dashboard-thumb-overlay);
        transition: opacity 0.2s;
    }
    .overlay {
        position: absolute;
        inset: 0;
        background: var(--dashboard-thumb-overlay);
    }

    @media (max-width: 768px) {
        .video-card--list {
            flex-direction: column;
            align-items: flex-start;
        }

        .video-card__checkbox {
            align-self: flex-start;
        }

        .video-card--list .thumb {
            width: 100%;
            height: auto;
            aspect-ratio: 16 / 9;
        }
    }
</style>
