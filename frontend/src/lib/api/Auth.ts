import { goto } from '$app/navigation'
import { PUBLIC_BASE_URL } from '$env/static/public'

export enum TokenValidationResult {
        "VALID",
        "INVALID",
        "NOT_VERIFIED"
}

export interface AuthForm {
        password: string
        email: string
}

export async function Login(form: AuthForm) {
        console.log("attemping to")
        const resp = await fetch(`${PUBLIC_BASE_URL}/api/users/login`, {
                method: "POST",
                credentials: "include",
                body: JSON.stringify(form),
                headers: {
                        "Content-Type": "application/json"
                }
        })

        if (resp.ok) return

        const body = await resp.json()

        // fuck javascript
        if (body.verified === false) {
                goto("/verify")
                return
        }

        throw new Error(body.error)
}

export async function Register(form: AuthForm): Promise<string> {
        const resp = await fetch(`${PUBLIC_BASE_URL}/api/users`, {
                method: "POST",
                body: JSON.stringify(form),
                headers: {
                        "Content-Type": "application/json",
                }
        })


        const body = await resp.json()
        if (resp.status === 200) return body.userID

        console.error(`[${body.requestID}] Login request failed: ${body.error}`)
        throw new Error(body.error)
}

export async function ValidateToken(): Promise<TokenValidationResult> {
        const resp = await fetch(`${PUBLIC_BASE_URL}/api/validate`, {
                method: "GET",
                credentials: 'include'
        })

        if (resp.status === 200) return TokenValidationResult.VALID

        const body = await resp.json()

        switch (body.error) {
                case "account_not_verified": {
                        return TokenValidationResult.NOT_VERIFIED
                }

                case "token_invalid":
                case "token_expired": {
                        return TokenValidationResult.INVALID
                }

                default: throw new Error(body.error)
        }
}