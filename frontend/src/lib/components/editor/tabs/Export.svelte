<script lang="ts">
    import { exportFormat, exportFps } from '$lib/stores/EditOptions'
</script>

<div class="tab-pane user-select-none" id="export-tab">
    <div class="mb-3">
        <label class="form-label" for="export-format">Format</label>
        <p class="small text-muted">Some formats may take longer to process</p>
        <select class="form-select" id="export-format" disabled bind:value={$exportFormat}>
            <option value="mp4">mp4</option>
            <option value="mov">mov</option>
            <option value="avi">avi</option>
            <option value="gif">gif</option>
        </select>
    </div>
    <div class="mb-3">
        <label class="form-label" for="export-framerate"> Framerate </label>
        <p class="small text-muted">Reducing the framerate will also reduce the file size</p>

        <select class="form-select" id="export-framerate" disabled bind:value={$exportFps}>
            <option value="Don't change">Don't change</option>
            <option value="60">60</option>
            <option value="30">30</option>
            <option value="24">24</option>
        </select>
    </div>
</div>
