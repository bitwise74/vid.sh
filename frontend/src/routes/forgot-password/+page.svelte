<script>
    import Header from '$lib/components/Header.svelte'
    import RequestPasswdReset from '$lib/components/password/RequestPasswdReset.svelte'
</script>

<svelte:head>
    <title>Forgot Password - vid.sh</title>
</svelte:head>

<div class="min-vh-100 gradient-bg">
    <Header page="forgot_password" title="Forgot Password" />

    <main class="container py-5">
        <div class="row justify-content-center align-items-center" style="min-height: calc(100vh - 200px);">
            <div class="col-lg-10">
                <div class="row align-items-center g-5">
                    <RequestPasswdReset />
                </div>
            </div>
        </div>
    </main>
</div>
