<script lang="ts">
    import { PUBLIC_BASE_URL } from '$env/static/public'
    import Header from '$lib/components/Header.svelte'
    import { onMount } from 'svelte'
    import '../../app.css'

    let status: 'loading' | 'success' | 'error' | 'info' = $state('info')
    let message = $state('Please wait...')

    let userID: string | null = $state('')

    let resendCooldown = 0
    let countdown = $state(resendCooldown)
    let disabled = $state(true)
    let interval: number

    async function handleResendCode() {
        alert("Feature isn't finished yet! Please contact me directly if you have issues with verification (bit.wise. on discord)")
    }

    onMount(async () => {
        const params = new URLSearchParams(window.location.search)
        const token = params.get('token')
        userID = params.get('user_id')

        // Page after register/login while not verified
        if (!token) {
            interval = setInterval(() => {
                countdown--
                if (countdown <= 0) {
                    disabled = false
                    clearInterval(interval)
                }
            }, 1000)
            return
        }

        try {
            const res = await fetch(`${PUBLIC_BASE_URL}/api/users/verify?token=${token}&user_id=${userID}`, {
                method: 'POST',
                credentials: 'include'
            })

            if (res.ok) {
                status = 'success'
                message = 'Your email has been verified successfully! 🎉'
            } else {
                const body = await res.json().catch(() => ({}))
                status = 'error'
                message = body.error || 'Verification failed. The link may have expired.'
            }
        } catch (err) {
            status = 'error'
            message = 'Unable to contact server. Please try again later.'
        }
    })

    function formatTime(seconds: number): string {
        const m = Math.floor(seconds / 60)
            .toString()
            .padStart(2, '0')
        const s = (seconds % 60).toString().padStart(2, '0')
        return `${m}:${s}`
    }
</script>

<svelte:head>
    <title>vid.sh - Email Verification</title>
</svelte:head>

<div class="d-flex flex-column min-vh-100">
    <Header title="Verify email" page="verify" />

    <main class="flex-grow-1">
        <section class="position-relative py-5 overflow-hidden" style="padding-top: 6rem; padding-bottom: 6rem;">
            <div class="position-absolute top-0 start-0 w-100 h-100 gradient-bg"></div>
            <div class="container position-relative">
                <div class="row justify-content-center">
                    <div class="col-lg-8 text-center">
                        {#if status === 'info'}
                            <div class="mb-2">
                                <i class="bi bi-envelope-check text-gradient" style="font-size: 5rem" role="status"></i>
                            </div>
                            <h2 class="fw-bold mb-3">One more thing...</h2>
                            <p class="text-muted">
                                We've sent you a verification email. You must verify your account before using vid.sh. This is to prevent bots from abusing the website for infinite storage and free
                                video processing. <span class="fw-bold"> Click one of the buttons to go to your preferred email website</span>
                            </p>

                            <div class="d-flex gap-2 my-4 flex-wrap align-items-center justify-content-center">
                                <a href="https://mail.google.com" target="_blank" class="btn btn-outline-primary d-flex align-items-center gap-1">
                                    <i class="bi bi-google"></i> Gmail
                                </a>

                                <a href="https://mail.proton.me" target="_blank" class="btn btn-outline-dark d-flex align-items-center gap-1">
                                    <i class="bi bi-lock-fill"></i> Proton Mail
                                </a>

                                <a href="https://outlook.live.com" target="_blank" class="btn btn-outline-info d-flex align-items-center gap-1">
                                    <i class="bi bi-microsoft"></i> Outlook
                                </a>
                            </div>

                            {#if userID}
                                <p class="text-muted-foreground">Didn't receive the message? Check your spam folder or click send again</p>

                                <button {disabled} onclick={handleResendCode} class="btn btn-gradient">
                                    <i class="bi bi-arrow-clockwise"></i>
                                    Resend Mail
                                    {#if resendCooldown > 0}
                                        ({formatTime(countdown)})
                                    {/if}
                                </button>
                            {/if}
                        {:else if status === 'success'}
                            <div class="mb-4">
                                <i class="bi bi-check-circle text-success display-1"></i>
                            </div>
                            <h2 class="fw-bold mb-3">Email Verified</h2>
                            <p class="fs-5 text-muted-foreground mb-4">{message}</p>
                            <a href="/dashboard" class="btn btn-gradient btn-lg px-4">
                                <i class="bi bi-house me-2"></i>
                                Go to Dashboard
                            </a>
                        {:else if status === 'error'}
                            <div class="mb-4">
                                <i class="bi bi-x-circle text-danger display-1"></i>
                            </div>
                            <h2 class="fw-bold mb-3">Verification Failed</h2>
                            <p class="fs-5 text-muted-foreground mb-4">{message}</p>
                            <button onclick={handleResendCode} class="btn btn-gradient btn-lg px-4">
                                <i class="bi bi-arrow-clockwise me-2"></i>
                                Resend Verification
                            </button>
                        {/if}
                    </div>
                </div>
            </div>
        </section>
    </main>

    <footer class="border-top bg-light py-4">
        <div class="container text-center">
            <p class="small text-muted-foreground mb-0">© 2025 vid.sh. All rights reserved.</p>
        </div>
    </footer>
</div>
