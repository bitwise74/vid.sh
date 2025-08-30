<script lang="ts">
    import ToastContainer from '$lib/components/toast/ToastContainer.svelte'
    import '../app.css'
  interface Props {
    children?: import('svelte').Snippet;
  }

  let { children }: Props = $props();
</script>

{@render children?.()}

<ToastContainer />
