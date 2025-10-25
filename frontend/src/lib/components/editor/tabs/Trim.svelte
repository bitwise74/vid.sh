<script lang="ts">
    import { trimEnd, trimStart, videoDuration } from '$lib/stores/EditOptions'
    import { currentTime } from '$lib/stores/VideoStore'
    import RangeSlider from 'svelte-range-slider-pips'
    import { get } from 'svelte/store'

    function handleTrimChange(
        e: CustomEvent<{
            values: number[]
        }>
    ) {
        const [start, end] = e.detail.values
        trimStart.set(start)
        currentTime.set(start)
        trimEnd.set(end)
    }

    function formatTime(time: number) {
        const minutes = Math.floor(time / 60)
        const seconds = Math.floor(time % 60)
        const milliseconds = Math.floor((time % 1) * 100)
        return `${minutes}:${seconds.toString().padStart(2, '0')}.${milliseconds.toString().padStart(2, '0')}`
    }

    function parseTime(str: string): number | null {
        const parts = str.split(':').map(Number)
        if (parts.some(isNaN)) return null
        if (parts.length === 2) return parts[0] * 60 + parts[1]
        return parts[0]
    }

    function clampTrimValues(start: number, end: number): [number, number] {
        const clampedStart = Math.max(0, Math.min(start, $videoDuration))
        const clampedEnd = Math.max(clampedStart, Math.min(end, $videoDuration))
        return [clampedStart, clampedEnd]
    }

    function handleTrimStartChange(e: FocusEvent) {
        const target = e.target as HTMLSpanElement
        const val = parseTime(target.innerText)
        if (val !== null) {
            const [newTrimStart, newTrimEnd] = clampTrimValues(val, $trimEnd)

            trimStart.set(newTrimStart)
            trimEnd.set(newTrimEnd)
        }
        target.innerText = formatTime($trimStart)
    }

    function handleTrimEndChange(e: FocusEvent) {
        const target = e.target as HTMLSpanElement
        const val = parseTime(target.innerText)
        if (val !== null) {
            const [newTrimStart, newTrimEnd] = clampTrimValues($trimStart, val)

            trimStart.set(newTrimStart)
            trimEnd.set(newTrimEnd)
        }
        target.innerText = formatTime($trimEnd)
    }
</script>

<div class="tab-pane show active" id="trim-tab">
    <p>
        <span class="text-muted-foreground user-select-none">From:</span>
        <span
            class="trim-time editable"
            contenteditable="true"
            role="textbox"
            tabindex="0"
            onblur={handleTrimStartChange}
            onkeydown={(e: KeyboardEvent) => {
                if (e.key === 'Enter') {
                    ;(e.target as HTMLSpanElement).blur()
                    e.preventDefault()
                }
            }}>
            {formatTime($trimStart)}
        </span>
        <span class="text-muted-foreground user-select-none">to:</span>
        <span
            class="trim-time editable"
            contenteditable="true"
            role="textbox"
            tabindex="0"
            onblur={handleTrimEndChange}
            onkeydown={(e: KeyboardEvent) => {
                if (e.key === 'Enter') {
                    ;(e.target as HTMLSpanElement).blur()
                    e.preventDefault()
                }
            }}>
            {formatTime($trimEnd)}
        </span>
    </p>

    <div class="mb-3">
        <div style="cursor: pointer;">
            <RangeSlider
                darkmode="auto"
                spring={false}
                rangeGapMin={$videoDuration * 0.025}
                values={[$trimStart, $trimEnd]}
                min={0}
                max={$videoDuration}
                step={0.01}
                float
                range
                draggy={true}
                on:stop={handleTrimChange}
                formatter={formatTime}
                style="--range-range: #1463DA">
            </RangeSlider>
        </div>

        <div class="d-flex justify-content-between">
            <button aria-label="Set Trim Start to Current Time" class="btn btn-outline-primary" title="Set Trim Start to Current Time" onclick={() => trimStart.set(get(currentTime))}>
                <i class="bi bi-skip-start-fill"></i>
            </button>
            <button aria-label="Set Trim End to Current Time" class="btn btn-outline-primary" title="Set Trim End to Current Time" onclick={() => trimEnd.set(get(currentTime))}>
                <i class="bi bi-skip-end-fill"></i>
            </button>
        </div>
    </div>

    <div class="small text-muted user-select-none">
        New duration: {formatTime($trimEnd - $trimStart)}<br />
        Tip: Click on the timestamps to edit
    </div>
</div>
