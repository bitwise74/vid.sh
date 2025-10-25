<script lang="ts">
    import { currentVideoURL } from '$lib/stores/VideoStore'
</script>

<div role="none" class="modal-backdrop show pt-0" onclick={() => currentVideoURL.set('')}></div>

<div role="none" class="modal show d-block d-flex pt-0" tabindex="-1" onclick={() => currentVideoURL.set('')}>
    <div role="none" class="modal-dialog modal-xl align-content-center" onclick={(e) => e.stopPropagation()}>
        <div class="modal-content">
            <video controls autoplay class="w-100" volume="0.75" src={$currentVideoURL}>
                <track kind="captions" />
            </video>
        </div>
    </div>
</div>
