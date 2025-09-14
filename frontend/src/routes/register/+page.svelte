<script>
    import Header from '$lib/components/Header.svelte'
    import RegisterForm from '$lib/components/RegisterForm.svelte'
</script>

<svelte:head>
    <title>Create Account - vid.sh</title>
</svelte:head>

<div class="min-vh-100 gradient-bg">
    <Header page="register" />

    <main class="container py-5">
        <div class="row justify-content-center align-items-center" style="min-height: calc(100vh - 200px);">
            <div class="col-lg-10">
                <div class="row align-items-center g-5">
                    <RegisterForm />
                </div>
            </div>
        </div>
    </main>
</div>
