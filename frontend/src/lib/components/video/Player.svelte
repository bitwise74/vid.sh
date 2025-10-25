<script lang="ts">
    import { trimEnd, trimStart, videoDuration, videoSource } from '$lib/stores/EditOptions'
    import { currentTime } from '$lib/stores/VideoStore'
    import { FormatDuration } from '$lib/utils/Format'
    import RangeSlider from 'svelte-range-slider-pips'
    import CropBox from '../editor/CropBox.svelte'

    let {
        onTimeUpdate
    }: {
        onTimeUpdate: (video: HTMLVideoElement) => void
    } = $props()

    let video: HTMLVideoElement = $state(null as any)

    let isPlaying = $state(false)
    let isMuted = $state(false)
    let volume = $state(0.75)
    let wasPaused = false

    function togglePlay() {
        if (video.currentTime >= $trimEnd && $trimEnd > $trimStart) {
            video.currentTime = $trimStart
            currentTime.set($trimStart)
        }

        video.paused ? video.play() : video.pause()
    }

    function handleSeekStart() {
        if (!video.paused) {
            wasPaused = true
            video.pause()
        }
    }

    function handleSeekStop(e: CustomEvent<{ value: number }>) {
        video.currentTime = e.detail.value
        currentTime.set(e.detail.value)

        if (wasPaused) {
            video.play()
            wasPaused = false
        }
    }

    function handleVolumeChange(e: Event) {
        volume = parseFloat((e.target as HTMLInputElement).value)

        video.volume = volume
        isMuted = volume === 0
    }

    function toggleMute() {
        if (isMuted) {
            volume = 0.5
            isMuted = false
        } else {
            volume = 0
            isMuted = true
        }
    }
</script>

<div>
    <div class="position-relative d-inline-block mb-3 aspect-video overflow-hidden rounded">
        <video
            bind:this={video}
            bind:currentTime={$currentTime}
            src={$videoSource}
            class="w-100 h-100"
            style="object-fit: contain;"
            onclick={togglePlay}
            onplay={() => (isPlaying = true)}
            onpause={() => (isPlaying = false)}
            ontimeupdate={() => onTimeUpdate(video!)}>
            <track kind="captions" />
        </video>

        <CropBox {video} />
    </div>

    <div>
        <div style="cursor: pointer;">
            <RangeSlider
                darkmode="auto"
                spring={false}
                limits={[$trimStart, $trimEnd]}
                max={$videoDuration}
                step={0.05}
                value={$currentTime}
                formatter={FormatDuration}
                on:start={() => handleSeekStart()}
                on:stop={(e) => handleSeekStop(e)} />
        </div>
        <div class="d-flex justify-content-between small text-muted mb-3">
            <span>{FormatDuration($currentTime)}</span>
            <span>{FormatDuration($videoDuration)}</span>
        </div>

        <div class="d-flex align-items-center justify-content-between">
            <div class="d-flex align-items-center gap-2">
                <button class="btn btn-outline-primary btn-sm" onclick={togglePlay} aria-label="Play/Pause">
                    <i class="bi bi-{isPlaying ? 'pause' : 'play'}-fill"></i>
                </button>

                <div class="d-flex align-items-center gap-2">
                    <button class="btn btn-outline-primary btn-sm" onclick={toggleMute} aria-label="Mute/Unmute">
                        <i class="bi bi-volume-{isMuted ? 'mute' : 'up'}-fill"></i>
                    </button>
                    <input type="range" class="form-range" style="width: 80px;" min="0" max="1" step="0.01" value={volume} oninput={handleVolumeChange} />
                </div>
            </div>

            <button class="btn btn-outline-primary btn-sm" aria-label="Fullscreen" onclick={() => video.requestFullscreen()}>
                <i class="bi bi-arrows-fullscreen"></i>
            </button>
        </div>
    </div>
</div>
