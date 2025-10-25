<script>
    import { page } from '$app/state'
    import Error from '$lib/components/Error.svelte'
    import Header from '$lib/components/Header.svelte'
    import NewPasswd from '$lib/components/password/NewPasswd.svelte'

    const token = page.url.searchParams.get('token')
</script>

<div class="min-vh-100 gradient-bg">
    <Header page="reset_passwd" title="Reset Password" />
    <main class="container py-5">
        <div class="row justify-content-center align-items-center">
            {#if !token}
                <Error title="Invalid password reset link" message="Please request a new one." buttons={[{ text: 'Request New Link', href: '/forgot-password', class: 'btn btn-gradient w-100' }]} />
            {:else}
                <NewPasswd {token} />
            {/if}
        </div>
    </main>
</div>
