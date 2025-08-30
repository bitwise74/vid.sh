<script lang="ts">
    import { currentTime } from '$lib/stores/VideoStore'
    import RangeSlider from 'svelte-range-slider-pips'

    let {
        src,
        duration,
        trimStart,
        trimEnd,
        onTimeUpdate
    }: {
        src: string
        duration: number
        trimStart: number
        trimEnd: number
        onTimeUpdate: (video: HTMLVideoElement) => void
    } = $props()

    let video: HTMLVideoElement

    let isPlaying = $state(false)
    let isMuted = $state(false)
    let volume = $state(0.75)

    let wasPaused = false

    function togglePlay() {
        if (video.currentTime >= trimEnd && trimEnd > trimStart) {
            video.currentTime = trimStart
            currentTime.set(trimStart)
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

    function formatTime(time: number) {
        const minutes = Math.floor(time / 60)
        const seconds = Math.floor(time % 60)
        const milliseconds = Math.floor((time % 1) * 100)
        return `${minutes}:${seconds.toString().padStart(2, '0')}.${milliseconds.toString().padStart(2, '0')}`
    }
</script>

<div>
    <div class="position-relative aspect-video bg-dark rounded overflow-hidden mb-3">
        <video
            bind:this={video}
            bind:currentTime={$currentTime}
            {src}
            class="w-100 h-100"
            style="object-fit: contain;"
            onclick={togglePlay}
            onplay={() => (isPlaying = true)}
            onpause={() => (isPlaying = false)}
            ontimeupdate={() => onTimeUpdate(video)}
            >
            <track kind="captions" />
        </video>
    </div>

    <div>
        <div class="mb-3">
            <RangeSlider
                spring={false}
                limits={[trimStart, trimEnd]}
                max={duration}
                step={0.05}
                value={$currentTime}
                formatter={formatTime}
                on:start={() => handleSeekStart()}
                on:stop={(e) => handleSeekStop(e)} />
            <div class="d-flex justify-content-between small text-muted">
                <span>{formatTime($currentTime)}</span>
                <span>{formatTime(duration)}</span>
            </div>
        </div>

        <div class="d-flex align-items-center justify-content-between">
            <div class="d-flex align-items-center gap-2">
                <button
                    class="btn btn-outline-dark btn-sm"
                    onclick={togglePlay}
                    aria-label="Play/Pause">
                    <i class="bi bi-{isPlaying ? 'pause' : 'play'}-fill"></i>
                </button>

                <div class="d-flex align-items-center gap-2">
                    <button
                        class="btn btn-outline-dark btn-sm"
                        onclick={toggleMute}
                        aria-label="Mute/Unmute">
                        <i class="bi bi-volume-{isMuted ? 'mute' : 'up'}-fill"></i>
                    </button>
                    <input
                        type="range"
                        class="form-range"
                        style="width: 80px;"
                        min="0"
                        max="1"
                        step="0.01"
                        value={volume}
                        oninput={handleVolumeChange} />
                </div>
            </div>

            <button
                class="btn btn-outline-dark btn-sm"
                aria-label="Fullscreen"
                onclick={() => video.requestFullscreen()}>
                <i class="bi bi-arrows-fullscreen"></i>
            </button>
        </div>
    </div>
</div>
