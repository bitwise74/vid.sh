<script lang="ts">
    import { toastStore } from '../../stores/ToastStore'
    import Toast from './Toast.svelte'
</script>

<div class="toast-container position-fixed bottom-0 end-0 p-3" style="z-index: 9999;">
    {#each $toastStore as toast (toast.id)}
        <Toast {toast} />
    {/each}
</div>
