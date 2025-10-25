<script lang="ts">
    import { RequestPasswordReset } from '$lib/api-v2/Auth'
    import { loadedInitialData } from '$lib/stores/AppVars'
    import { toastStore } from '../../stores/ToastStore'

    type FormState = 'waiting' | 'submitting' | 'success' | 'error'

    let email = $state('')
    let formState = $state<FormState>('waiting')

    async function handleSubmit(e: Event) {
        e.preventDefault()
        formState = 'submitting'

        try {
            if (!email) {
                throw new Error('Please fill in all fields')
            }

            if (!email.includes('@')) {
                throw new Error('Please enter a valid email address')
            }

            await RequestPasswordReset(email)
            loadedInitialData.set(false)
            formState = 'success'

            toastStore.success({
                title: 'Password reset email sent',
                message: 'Please check your inbox for further instructions',
                duration: 30000
            })
        } catch (err) {
            console.error('Failed to reset password', err)
            toastStore.error({
                title: 'Failed to reset password',
                message: err,
                duration: 10000
            })

            formState = 'error'
        }
    }
</script>

<div class="card border-0 shadow-lg" style="max-width: 28rem; margin: 0 auto;">
    <div class="border-0 py-4 text-center">
        <h2 class="card-title h3 fw-bold mb-0">Reset Password</h2>
        <p class="text-muted-foreground mb-0">Enter your email to reset your password</p>
    </div>
    <div class="card-body p-4">
        <form>
            <div class="mb-3">
                <label for="email" class="form-label">Email</label>
                <div class="input-group">
                    <span class="input-group-text">
                        <i class="bi bi-envelope text-muted"></i>
                    </span>
                    <input type="email" class="form-control" id="email" placeholder="Enter your email" bind:value={email} />
                </div>
            </div>
        </form>

        <button type="submit" class="btn btn-gradient w-100 mb-4" disabled={formState !== 'waiting' || !email} onclick={handleSubmit}>
            {#if formState === 'submitting'}
                <span class="spinner-border spinner-border-sm me-2" role="status"></span>
                Processing...
            {:else}
                Send Reset Link
            {/if}
        </button>

        <div class="text-center">
            <span class="text-muted-foreground">Don't have an account? </span>
            <a href="/register" class="text-decoration-none fw-medium">Sign up</a>
        </div>
    </div>
</div>
