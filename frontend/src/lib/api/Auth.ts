import { goto } from '$app/navigation'
import { PUBLIC_BASE_URL } from '$env/static/public'
import { toastStore } from '$lib/stores/ToastStore'

export type AuthForm = {
    password: string
    email: string
    remember: boolean
}

/**
 * Authorizes a user from credentials
 * @param form Auth form with login credentials
 * @returns User ID and verification status
 */
export async function Login(form: AuthForm) {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/users/login`, {
        method: 'POST',
        body: JSON.stringify(form),
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json'
        }
    })

    if (!req.ok) {
        const body = await req.json()
        console.error(`[Auth/Login]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }
}

/**
 * Registers a user
 * @param form Credentials
 */
export async function Register(form: Omit<AuthForm, 'remember'>) {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/users/register`, {
        method: 'POST',
        body: JSON.stringify(form)
    })

    const body = await req.json()

    if (!req.ok) {
        console.error(`[Auth/Register]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    goto('/verify')
}

/**
 * Requests a password reset link
 * @param email User email
 */
export async function RequestPasswordReset(email: string) {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/users/password-reset`, {
        method: 'POST',
        body: JSON.stringify({ email })
    })

    const body = await req.json()

    if (!req.ok) {
        console.error(`[Auth/RequestPasswordReset]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }

    return body
}

/**
 * Resets a user's password
 * @param token Password reset token
 * @param newPassword New password
 */
export async function ResetPassword(token: string, newPassword: string) {
    const req = await fetch(`${PUBLIC_BASE_URL}/api/users/password-reset/confirm`, {
        method: 'POST',
        body: JSON.stringify({ token, newPassword })
    })

    const body = await req.json()

    if (!req.ok) {
        console.error(`[Auth/ResetPassword]: Request failed, requestID: ${body.requestID}`, body.error)
        throw new Error(body.error, { cause: req })
    }
}

/**
 * Logs out the current user
 */
export async function Logout() {
    await fetch(`${PUBLIC_BASE_URL}/api/users/logout`, {
        method: 'POST',
        credentials: 'include'
    })

    if (window.location.pathname === '/') {
        window.location.reload()
    } else {
        goto('/')
    }

    localStorage.setItem("optDisableRoot", "false")

    toastStore.success({ title: 'Logged out successfully' })
}
