<script lang="ts">
  import Toast from './Toast.svelte'
  import { toastStore } from './toastStore'
</script>

<div class="toast-container position-fixed bottom-0 end-0 p-3" style="z-index: 9999;">
  {#each $toastStore as toast (toast.id)}
    <Toast toast={toast}/>
  {/each}
</div>
