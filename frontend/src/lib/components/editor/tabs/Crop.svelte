<script lang="ts">
    import { isCroppingEnabled } from '$lib/stores/EditOptions'
    import { onMount } from 'svelte'

    let firstVisit = true
    const cropTabBtn = document.querySelector('[data-bs-target="#crop-tab"')

    onMount(() => {
        cropTabBtn?.addEventListener('show.bs.tab', () => {
            if (firstVisit) {
                isCroppingEnabled.set(true)
            }
        })
    })
</script>

<div class="tab-pane user-select-none" id="crop-tab">
    <input class="form-check-input" type="checkbox" id="toggle-cropping" bind:checked={$isCroppingEnabled} />
    <label for="toggle-cropping" class="small">Enable cropping</label>
</div>
