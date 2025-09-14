<script lang="ts">
    import { goto } from '$app/navigation'
    import { Register } from '$lib/api/Auth'
    import '../../app.css'
    import { toastStore } from './toast/toastStore'

    let showPassword = $state(false)
    let showConfirmPassword = $state(false)
    let isLoading = $state(false)
    let success = $state(false)
    let formData = $state({
        email: '',
        password: '',
        confirmPassword: '',
        agreeToTerms: false
    } as { [key: string]: any })

    function validatePassword(password: string) {
        const minLength = password.length >= 8

        return {
            minLength,
            isValid: minLength
        }
    }

    let passwordValidation = $state(validatePassword(formData.password))

    async function handleSubmit(e: Event) {
        e.preventDefault()
        isLoading = true

        try {
            if (!formData.email || !formData.password || !formData.confirmPassword) {
                throw new Error('Please fill in all required fields')
            }

            if (!formData.email.includes('@')) {
                throw new Error('Please enter a valid email address')
            }

            if (!passwordValidation.isValid) {
                throw new Error('Password does not meet requirements')
            }

            if (formData.password !== formData.confirmPassword) {
                throw new Error('Passwords do not match')
            }

            if (!formData.agreeToTerms) {
                throw new Error('You must agree to the Terms of Service and Privacy Policy')
            }

            const userID = await Register({ email: formData.email, password: formData.password })
            goto('/verify?userID=' + userID)
        } catch (err) {
            toastStore.error({
                title: 'Failed to register account',
                message: err.message,
                duration: 10000
            })
        } finally {
            isLoading = false
        }
    }

    function handleInputChange(field: string, value: any) {
        formData[field] = value
        passwordValidation = validatePassword(formData.password)
    }

    function handleCheckboxChange(field: string, checked: any) {
        formData[field] = checked
    }
</script>

<div class="card shadow-lg border-0" style="max-width: 28rem; margin: 0 auto;">
    <div class="card-header bg-white text-center border-0 py-4">
        <h2 class="card-title h3 fw-bold mb-2">Create your account</h2>
        <p class="text-muted-foreground mb-0">Provide your details in the form</p>
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
                        placeholder="john@example.com"
                        bind:value={formData.email}
                        oninput={(e) => handleInputChange('email', (e.target as HTMLInputElement).value)}
                        required />
                </div>
            </div>

            <div class="mb-3">
                <label for="password" class="form-label">Password</label>
                <div class="input-group">
                    <span class="input-group-text">
                        <i class="bi bi-lock text-muted"></i>
                    </span>
                    <input
                        type={showPassword ? 'text' : 'password'}
                        class="form-control"
                        id="password"
                        placeholder="Create a strong password"
                        bind:value={formData.password}
                        oninput={(e) => handleInputChange('password', (e.target as HTMLInputElement).value)}
                        required />
                    <button type="button" class="btn btn-outline-secondary" onclick={() => (showPassword = !showPassword)} aria-label="Toggle password visibility">
                        <i class="bi bi-eye{showPassword ? '-slash' : ''}"></i>
                    </button>
                </div>

                {#if formData.password}
                    <div class="mt-2">
                        <div class="small d-flex align-items-center mb-1 {passwordValidation.minLength ? 'text-success' : 'text-muted'}">
                            <div class="rounded-circle me-2 {passwordValidation.minLength ? 'bg-success' : 'bg-secondary'}" style="width: 4px; height: 4px;"></div>
                            <span>At least 8 characters</span>
                        </div>
                    </div>
                {/if}
            </div>

            <div class="mb-3">
                <label for="confirmPassword" class="form-label">Confirm Password</label>
                <div class="input-group">
                    <span class="input-group-text">
                        <i class="bi bi-lock text-muted"></i>
                    </span>
                    <input
                        type={showConfirmPassword ? 'text' : 'password'}
                        class="form-control"
                        id="confirmPassword"
                        placeholder="Confirm your password"
                        bind:value={formData.confirmPassword}
                        oninput={(e) => handleInputChange('confirmPassword', (e.target as HTMLInputElement).value)}
                        required />
                    <button type="button" class="btn btn-outline-secondary" onclick={() => (showConfirmPassword = !showConfirmPassword)} aria-label="Toggle confirm password visibility">
                        <i class="bi bi-eye{showConfirmPassword ? '-slash' : ''}"></i>
                    </button>
                </div>
                {#if formData.confirmPassword && formData.password !== formData.confirmPassword}
                    <div class="small text-danger mt-1">Passwords do not match</div>
                {/if}
            </div>

            <div class="mb-4">
                <div class="form-check mb-3">
                    <input
                        class="form-check-input"
                        type="checkbox"
                        id="agreeToTerms"
                        bind:checked={formData.agreeToTerms}
                        onchange={(e) => handleCheckboxChange('agreeToTerms', (e.target as HTMLInputElement).checked)} />
                    <label class="form-check-label small" for="agreeToTerms">
                        I agree to the
                        <a href="/terms" class="text-decoration-none">Terms of Service</a> and
                        <a href="/privacy" class="text-decoration-none">Privacy Policy</a>
                    </label>
                </div>
            </div>

            <button type="submit" class="btn btn-gradient w-100 mb-4" disabled={isLoading || !passwordValidation.isValid}>
                {#if isLoading}
                    <span class="spinner-border spinner-border-sm me-2" role="status"></span>
                    Creating account...
                {:else}
                    Create Account
                {/if}
            </button>

            <div class="text-center mb-4">
                <div class="position-relative">
                    <hr />
                    <span class="position-absolute top-50 start-50 translate-middle bg-white px-3 small text-muted"> Or continue with </span>
                </div>
            </div>

            <div class="row g-2 mb-4">
                <div class="col-6">
                    <button type="button" class="btn btn-outline-secondary w-100" aria-label="Sign in with Google">
                        <i class="bi bi-google me-1"></i>
                        Google
                    </button>
                </div>
                <div class="col-6">
                    <button type="button" class="btn btn-outline-secondary w-100" aria-label="Sign in with GitHub">
                        <i class="bi bi-github me-1"></i>
                        GitHub
                    </button>
                </div>
            </div>
        </form>

        <div class="text-center">
            <span class="text-muted-foreground">Already have an account? </span>
            <a href="/login" class="text-decoration-none fw-medium" aria-label="Sign in">Sign in</a>
        </div>
    </div>
</div>
