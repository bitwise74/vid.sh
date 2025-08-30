<script lang="ts">
    import { goto } from '$app/navigation'
    import { Login } from '$lib/api/Auth'
    import { writable } from 'svelte/store'
    import { toastStore } from './toast/toastStore'

    const showPassword = writable(false)
    const isLoading = writable(false)
    const formData = writable({ email: '', password: '' } as { [key: string]: any })

    async function handleSubmit(e: Event) {
        e.preventDefault()
        isLoading.set(true)

        const { email, password } = $formData

        try {
            if (!email || !password) {
                throw new Error('Please fill in all fields')
            }

            if (!email.includes('@')) {
                throw new Error('Please enter a valid email address')
            }

            await Login({ email, password })
            goto('/dashboard')
        } catch (err) {
                console.error("Failed to login", err)
            toastStore.error({
                title: 'Failed to login',
                message: err,
                duration: 10000
            })
        } finally {
            isLoading.set(false)
        }
    }

    function handleInputChange(field: string, e: Event) {
        const value = (e.target as HTMLInputElement).value
        formData.update((data) => ({ ...data, [field]: value }))
    }
</script>

<div class="card shadow-lg border-0" style="max-width: 28rem; margin: 0 auto;">
    <div class="card-header bg-white text-center border-0 py-4">
        <h2 class="card-title h3 fw-bold mb-2">Welcome back</h2>
        <p class="text-muted-foreground mb-0">Sign in to your vid.sh account to continue editing</p>
    </div>
    <div class="card-body p-4">
        <form onsubmit={handleSubmit}>
            <div class="mb-3">
                <label for="email" class="form-label">Email</label>
                <div class="input-group">
                    <span class="input-group-text">
                        <i class="bi bi-envelope text-muted"></i>
                    </span>
                    <input
                        type="email"
                        class="form-control"
                        id="email"
                        placeholder="Enter your email"
                        bind:value={$formData.email}
                        oninput={(e) => handleInputChange('email', e)} />
                </div>
            </div>

            <div class="mb-3">
                <label for="password" class="form-label">Password</label>
                <div class="input-group">
                    <span class="input-group-text">
                        <i class="bi bi-lock text-muted"></i>
                    </span>
                    <input
                        type={$showPassword ? 'text' : 'password'}
                        class="form-control"
                        id="password"
                        placeholder="Enter your password"
                        bind:value={$formData.password}
                        oninput={(e) => handleInputChange('password', e)}
                        required />
                    <button
                        type="button"
                        class="btn btn-outline-secondary"
                        onclick={() => showPassword.set(!$showPassword)}
                        aria-label="Toggle password visibility">
                        <i class="bi bi-eye{$showPassword ? '-slash' : ''}"></i>
                    </button>
                </div>
            </div>

            <div class="d-flex justify-content-between align-items-center mb-4">
                <div class="form-check">
                    <input class="form-check-input" type="checkbox" id="remember" />
                    <label class="form-check-label small" for="remember"> Remember me </label>
                </div>
                <a href="/forgot-password" class="small text-decoration-none">Forgot password?</a>
            </div>

            <button type="submit" class="btn btn-gradient w-100 mb-4" disabled={$isLoading}>
                {#if $isLoading}
                    <span class="spinner-border spinner-border-sm me-2" role="status"></span>
                    Signing in...
                {:else}
                    Sign In
                {/if}
            </button>

            <div class="text-center mb-4">
                <div class="position-relative">
                    <hr />
                    <span
                        class="position-absolute top-50 start-50 translate-middle bg-white px-3 small text-muted">
                        Or continue with (coming soon)
                    </span>
                </div>
            </div>

            <div class="row g-2 mb-4">
                <div class="col-6">
                    <button
                        type="button"
                        class="btn btn-outline-secondary w-100"
                        aria-label="Sign in with Google"
                        disabled>
                        <i class="bi bi-google me-1"></i>
                        Google
                    </button>
                </div>
                <div class="col-6">
                    <button
                        type="button"
                        class="btn btn-outline-secondary w-100"
                        aria-label="Sign in with GitHub"
                        disabled>
                        <i class="bi bi-github me-1"></i>
                        GitHub
                    </button>
                </div>
            </div>
        </form>

        <div class="text-center">
            <span class="text-muted-foreground">Don't have an account? </span>
            <a href="/register" class="text-decoration-none fw-medium">Sign up</a>
        </div>
    </div>
</div>
