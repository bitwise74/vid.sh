<script lang="ts">
    import { PUBLIC_DISCORD_PERMALINK } from '$env/static/public'

    type Props = {
        title: string
        message?: string
        code?: number | undefined
        buttons?: Array<{
            text: string
            href?: string
            icon?: string
            class?: string
        }>
    }

    let { title = 'An error occurred', message = 'Something went wrong. Please try again later.', code, buttons = [] }: Props = $props()
</script>

<div class="card border-0 shadow-lg" style="max-width: 28rem; margin: 0 auto;">
    <div class="border-0 py-4 text-center">
        <h2 class="card-title h3 fw-bold mb-2">{code || ''} {title}</h2>
        <p class="text-muted-foreground mb-0 text-center">{message}</p>
    </div>
    <div class="card-body p-4">
        {#each buttons as btn}
            {#if btn.href}
                <a href={btn.href} class="btn btn-sm shadow-bottom rounded-2 btn-lg p-2 px-3 {btn.class}">
                    {#if btn.icon}
                        <i class="{btn.icon} me-2"></i>
                    {/if}
                    {btn.text}
                </a>
            {:else}
                <button class="btn btn-sm shadow-bottom rounded-2 btn-lg p-2 px-3 {btn.class}" onclick={() => history.back()}>
                    {#if btn.icon}
                        <i class="{btn.icon} me-2"></i>
                    {/if}
                    {btn.text}
                </button>
            {/if}
        {/each}

        {#if PUBLIC_DISCORD_PERMALINK}
            <div class="text-center mt-4">
                <span class="text-muted-foreground">Need help? </span>
                <a href={PUBLIC_DISCORD_PERMALINK} class="text-decoration-none fw-medium">Join our Discord</a>
            </div>
        {/if}
    </div>
</div>
