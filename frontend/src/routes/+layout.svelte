<script lang="ts">
    import ToastContainer from '$lib/components/toast/ToastContainer.svelte'
    import '../app.css'
    interface Props {
        children?: import('svelte').Snippet
    }

    let { children }: Props = $props()
</script>

<svelte:head>
    <link rel="icon" href="/favicon.svg" />
</svelte:head>
{@render children?.()}

<ToastContainer />
