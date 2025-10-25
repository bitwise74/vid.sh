<script lang="ts">
    import { goto } from '$app/navigation'
    import { getCookie } from '$lib/utils/Cookies'
    import { UploadFileButton } from '$lib/utils/Upload'
    import { onMount } from 'svelte'
    import Avatar from './user/avatar.svelte'

    let { page = '', subpage = '', title = '' } = $props()
    let loggedIn = $state(false)
    let href = $state('')

    onMount(() => {
        href = localStorage.getItem('optDisableRoot') === 'true' ? '/dashboard' : '/'
    })

    if (getCookie('logged_in') === '1') {
        loggedIn = true
    }

    const handleLogin = async () => {
        goto(loggedIn ? '/dashboard' : '/login')
    }

    const buttonConfig = {
        main: () =>
            loggedIn
                ? [
                      {
                          text: 'Dashboard',
                          href: '/dashboard',
                          icon: 'bi-view-stacked',
                          class: 'btn-dark'
                      }
                  ]
                : [
                      {
                          text: 'Log in',
                          action: handleLogin,
                          class: 'btn-gradient',
                          icon: 'bi-box-arrow-in-left',
                          important: true
                      },
                      {
                          text: 'Register',
                          href: '/register',
                          class: 'btn-gradient',
                          icon: 'bi-plus',
                          important: true
                      }
                  ],
        dashboard: () => [
            {
                text: 'Upload',
                action: async () => await UploadFileButton(),
                icon: 'bi-upload',
                class: 'btn-dark'
            },
            {
                text: 'Editor',
                href: '/editor',
                icon: 'bi-pencil',
                class: 'btn-dark'
            }
        ],
        login: () => [{ text: 'Go Back', href: '/', icon: 'bi-arrow-left', class: 'btn-dark', important: true }],
        verify: () => [{ text: 'Go Back', href: '/', icon: 'bi-arrow-left', class: 'btn-dark', important: true }],
        register: () => [{ text: 'Go Back', href: '/', icon: 'bi-arrow-left', class: 'btn-dark', important: true }],
        editor: () => [
            loggedIn
                ? {
                      text: 'Dashboard',
                      href: '/dashboard',
                      icon: 'bi-view-stacked',
                      class: 'btn-dark'
                  }
                : {
                      text: 'Go Back',
                      href: '/',
                      icon: 'bi-arrow-left',
                      class: 'btn-dark'
                  }
        ],
        profile: () => [
            {
                text: 'Go Back',
                href: '/dashboard',
                icon: 'bi-arrow-left',
                class: 'btn-dark'
            }
        ],
        settings: () => [
            {
                text: 'Go Back',
                href: '/dashboard',
                icon: 'bi-arrow-left',
                class: 'btn-dark'
            }
        ],
        forgot_password: () => [
            {
                text: 'Go Back',
                href: '/login',
                icon: 'bi-arrow-left',
                class: 'btn-dark',
                important: true
            }
        ]
    }

    const buttons = buttonConfig[subpage ? `${page}/${subpage}` : page]?.() || []
</script>

<header class="bg-body-tertiary sticky-top-custom">
    <div class="container">
        <div class="d-flex align-items-center justify-content-between py-3">
            <div class="d-flex align-items-center">
                <a {href} class="d-flex align-items-center text-decoration-none me-3">
                    <img src="/favicon.svg" width="38" height="38" class="navbar-brand" alt="logo" />
                    <span class="fs-4 ps-2 fw-bold text-gradient">vid.sh</span>
                </a>

                <div class="d-none d-sm-flex">
                    {#if title}
                        <div class="vr me-3"></div>
                        <h1 class="fs-5 fw-semibold mb-0">{title}</h1>
                    {/if}
                </div>
            </div>

            <div class="d-flex align-items-center gap-3 flex-wrap">
                {#each buttons as btn}
                    {#if btn.href}
                        <a href={btn.href} class="btn btn-sm shadow-bottom rounded-2 btn-lg p-2 px-3 {btn.class} {!btn.important ? 'd-none d-md-inline-flex' : ''}">
                            {#if btn.icon}<i class="{btn.icon} me-1"></i>{/if}
                            {btn.text}
                        </a>
                    {:else if btn.action}
                        <button onclick={btn.action} class="rounded-2 btn btn-sm shadow-bottom btn-lg p-2 px-3 {btn.class} {!btn.important ? 'd-none d-md-inline-flex' : ''}">
                            {#if btn.icon}<i class="{btn.icon} me-1"></i>{/if}
                            {btn.text}
                        </button>
                    {/if}
                {/each}
                {#if loggedIn}
                    <Avatar />
                {/if}
            </div>
        </div>
    </div>
</header>
