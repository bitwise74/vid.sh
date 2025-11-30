<script lang="ts">
    import { invalidateAll } from '$app/navigation'
    import { ResetPassword } from '$lib/api/Auth'
    import { toastStore } from '$lib/stores/ToastStore'

    type Props = {
        token: string
    }

    let newPassword = $state('')
    let confirmPassword = $state('')
    let isSubmitting = $state(false)

    let { token }: Props = $props()

    async function handleSubmit() {
        if (!newPassword || !confirmPassword) {
            toastStore.error({
                title: 'Please fill in all fields',
                duration: 10000
            })
            return
        }

        if (newPassword !== confirmPassword) {
            toastStore.error({
                title: 'Passwords do not match',
                duration: 10000
            })
            return
        }

        try {
            await ResetPassword(token, newPassword)
            await invalidateAll()

            toastStore.success({
                title: 'Password reset successful',
                message: 'You can now log in with your new password',
                duration: 10000
            })
        } catch (err) {
            console.error('Failed to reset password', err)
            toastStore.error({
                title: 'Failed to reset password',
                message: err.message,
                duration: 30000
            })
        }
    }
</script>

<div class="card border-0 shadow-lg" style="max-width: 28rem; margin: 0 auto;">
    <div class="border-0 py-4 text-center">
        <h2 class="card-title h3 fw-bold mb-0">Reset Password</h2>
    </div>
    <div class="card-body p-4">
        <form>
            <div class="mb-3">
                <label for="new-password" class="form-label">New Password</label>
                <div class="input-group">
                    <span class="input-group-text">
                        <i class="bi bi-lock text-muted"></i>
                    </span>
                    <input type="password" class="form-control" id="new-password" placeholder="Enter your new password" bind:value={newPassword} />
                </div>

                <label for="confirm-password" class="form-label mt-3">Confirm Password</label>
                <div class="input-group">
                    <span class="input-group-text">
                        <i class="bi bi-lock text-muted"></i>
                    </span>
                    <input type="password" class="form-control" id="confirm-password" placeholder="Confirm your new password" bind:value={confirmPassword} />
                </div>
            </div>
        </form>

        <button type="submit" class="btn btn-gradient w-100 mb-4" disabled={!newPassword || newPassword !== confirmPassword} onclick={handleSubmit}>
            {#if isSubmitting}
                <span class="spinner-border spinner-border-sm me-2" role="status"></span>
                Processing...
            {:else}
                Confirm
            {/if}
        </button>
    </div>
</div>
