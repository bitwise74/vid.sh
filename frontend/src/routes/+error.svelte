<script lang="ts">
    import { page } from '$app/state'
    import Header from '$lib/components/Header.svelte'
</script>

<Header title={`Error ${page.status} - ${page.error?.message || 'Something went wrong'}`} />

<div class="container mt-5 text-center">
    {#if page.status === 404}
        <i class="bi-exclamation-triangle display-1 text-danger mb-3"></i>
        <h1 class="display-3 text-danger mb-3">404</h1>
        <h2 class="mb-3">Page Not Found</h2>
        <p class="text-muted mb-4">Sorry, we couldn't find the page you're looking for. It might have been moved, renamed, or never existed.</p>
        <a href="/" class="btn btn-primary btn-lg shadow-bottom">
            <i class="bi bi-house-door-fill me-2"></i>Go Back Home
        </a>
    {:else}
        <i class="bi-bug display-1 text-warning mb-3"></i>
        <h1 class="display-3 text-warning mb-3">{page.status}</h1>
        <h2 class="mb-3">Something went really wrong</h2>
        <p class="text-muted mb-4">
            {page.error?.message || 'An unexpected error occurred. Please try again later.'}
        </p>
        <a href="/" class="btn btn-primary btn-lg shadow-bottom">
            <i class="bi bi-house-door-fill me-2"></i>Return Home
        </a>
    {/if}
</div>
