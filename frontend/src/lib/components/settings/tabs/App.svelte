<script lang="ts">
    import RangeSlider from 'svelte-range-slider-pips'

    const set = (key: string, e: Event, type = 'text') => {
        let value: any

        if (type === 'bool') {
            value = (e.target as HTMLInputElement).checked
        } else if (type === 'number') {
            value = parseInt((e.target as HTMLInputElement).value)
        } else {
            value = (e.target as HTMLInputElement).value
        }

        if (type === 'number' && isNaN(parseInt(value as any))) {
            return
        }

        localStorage.setItem(key, String(value))
    }

    const getSetting = (key: string) => {
        return localStorage.getItem(key)
    }
</script>

<div class="tab-pane" id="appearance-tab">
    <h4>App</h4>
    <hr />
    <div class="row mb-2">
        <div class="col-6">
            <label for="optPreferredTheme" class="form-label fw-semibold">Preferred theme</label>
            <select class="form-select mb-1" aria-label="Preferred theme" onchange={(e) => set('optPreferredTheme', e)} value={getSetting('optPreferredTheme')}>
                <option value="light">Light</option>
                <option value="dark">Dark</option>
                <option value="auto">System</option>
            </select>
            <p class="small text-muted">If specified the automatic theme will be overwritten</p>
        </div>

        <div class="col-6">
            <label for="optDefaultTargetSize" class="form-label fw-semibold">Default Target Size</label>
            <div class="input-group">
                <input
                    type="number"
                    max={2000}
                    min={0}
                    onchange={(e) => set('optDefaultTargetSize', e, 'number')}
                    id="optDefaultTargetSize"
                    class="form-control"
                    value={getSetting('optDefaultTargetSize')} />
                <span class="input-group-text text-muted">MB</span>
            </div>
            <p class="small text-muted">If set the editor tab will automatically set the target file size if the file is big enough. <b>To disable simply set it to 0MB</b></p>
        </div>
    </div>

    <div class="row mb-2 user-select-none">
        <div>
            <input id="lossless_export" class="form-check-input" type="checkbox" onclick={(e) => set('optLosslessExport', e, 'bool')} checked={getSetting('optLosslessExport') === 'true'} />
            <label for="lossless_export" class="form-label fw-semibold">Default Lossless Exports</label>
            <p class="small text-muted">If set the editor will automatically use lossless exports. Files may end up being bigger than the target size</p>
        </div>
        <div>
            <input type="checkbox" id="auto_copy_url" class="form-check-input" onclick={(e) => set('optAutoCopy', e, 'bool')} checked={getSetting('optAutoCopy') === 'true'} />
            <label for="auto_copy_url" class="form-label fw-semibold">Auto Copy URLs</label>
            <p class="small text-muted">Automatically copy the video URL when saving to cloud</p>
        </div>
        <div>
            <input type="checkbox" id="disable_main_page" class="form-check-input" onclick={(e) => set('optDisableRoot', e, 'bool')} checked={getSetting('optDisableRoot') === 'true'} />
            <label for="disable_main_page" class="form-label fw-semibold">Disable Main Page</label>
            <p class="small text-muted">Automatically redirect from the main page to the dashboard</p>
        </div>
        <!-- <div>
            <input type="checkbox" id="auto_merge_audio" class="form-check-input" onclick={(e) => handleChangeSetting('optAutoMergeAudio', e)} />
            <label for="auto_merge_audio" class="form-label fw-semibold">Auto Merge Audio Tracks</label>
            <p class="small text-muted">If multiple audio tracks are detected automatically merge them (or you won't be able to hear all audio tracks)</p>
        </div> -->
        <div>
            <input type="checkbox" id="use_rich_embeds" class="form-check-input" onclick={(e) => set('optRichEmbeds', e, 'bool')} checked={getSetting('optRichEmbeds') === 'true'} />
            <label for="use_rich_embeds" class="form-label fw-semibold">Use Rich Embeds</label>
            <p class="small text-muted">
                If checked and both the username and profile picture are set videos shared via URLs will have a pretty rich embed around them (check image if you don't know what that means)
            </p>
        </div>
    </div>

    <hr />
    <h4>Visuals</h4>
    <div class="row mb-2 user-select-none">
        <div class="col-6">
            <label for="animation_speed" class="form-label fw-semibold">Animation Speed</label>
            <RangeSlider
                min={0}
                max={2}
                step={0.5}
                pips
                darkmode={false}
                all={'label'}
                value={parseInt(getSetting('optAnimationSpeed') || '1')}
                springValues={{ stiffness: 0.25, damping: 0.75 }}
                on:stop={(e) => {
                    const sliderValue = e.detail.value as number
                    localStorage.setItem('optAnimationSpeed', String(sliderValue))
                }} />
            <p class="small text-muted">Adjust the speed of animations throughout the app (1.0 is normal speed, 0 disables animations)</p>
        </div>
    </div>
</div>
