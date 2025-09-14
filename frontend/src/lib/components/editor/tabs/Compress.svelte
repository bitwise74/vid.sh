<script lang="ts">
    import { losslessExport, targetSize } from '$lib/stores/EditOptions'
    import RangeSlider from 'svelte-range-slider-pips'

    let { videoSize }: { videoSize: number } = $props()
    let isCompressingDisabled = videoSize < 5

    const handleCompressChange = (e: CustomEvent<{ value: number }>) => {
        targetSize.set(e.detail.value)
    }
</script>

<div class="tab-pane" id="compress-tab">
    <div class="mb-3">
        <label for="target-size" class={isCompressingDisabled ? 'text-muted' : ''}>
            Target Size:
            {#if $targetSize <= 0 || $targetSize >= videoSize}
                Not set
            {:else}
                {$targetSize} MB
            {/if}
        </label>

        <div style="cursor: pointer;">
            <RangeSlider
                spring={false}
                value={0}
                disabled={isCompressingDisabled}
                min={0}
                max={videoSize}
                step={1}
                float
                on:change={handleCompressChange}
                formatter={(val) => {
                    if (val == 0 || val == videoSize) {
                        targetSize.set(-1)
                        return 'Not set'
                    }

                    return val.toString() + ' MB'
                }} />
        </div>
    </div>
    <div class="mb-3">
        <label class="form-label" for="processing-speed">Other processing options</label>
        <p class="small text-muted">Lossless exports produce a better video quality, but take longer to finish and don't optimize the video size.</p>
        <input class="form-check-input" type="checkbox" id="losslessExport" bind:checked={$losslessExport} />
        <label class="form-check-label small" for="losslessExport">Lossless export</label>
    </div>
</div>
