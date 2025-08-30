<script lang="ts">
    import { goto } from '$app/navigation'
    import { type Video, UploadVideo } from '$lib/api/Files'
    import { videos } from '$lib/stores/VideoStore'
    import { getCookie } from '$lib/utils/Cookies'
    import { toastStore } from './toast/toastStore'

    let { page = '', subpage = '', title = '' } = $props()
    let loggedIn = $state(false)

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
                          class: 'btn-gradient'
                      }
                  ]
                : [
                      {
                          text: 'Log in',
                          action: handleLogin,
                          class: 'btn-gradient'
                      },
                      {
                          text: 'Register',
                          href: '/register',
                          class: 'btn-gradient'
                      }
                  ],
        dashboard: () => [
            {
                text: 'Upload Video',
                action: () => {
                    const input = document.createElement('input')
                    input.type = 'file'
                    input.accept = 'video/*'

                    input.onchange = async (e: Event) => {
                        const target = e.target as HTMLInputElement
                        if (!target.files || target.files.length === 0) return

                        const videoFile = Array.from(target.files).find((f) =>
                            ['video/mp4', 'video/quicktime', 'video/x-matroska'].includes(f.type)
                        )
                        if (!videoFile) {
                            toastStore.error({
                                title: 'No valid files detected',
                                message: 'Please use one of the supported formats (mp4, mov, mkv)',
                                duration: 10000
                            })
                            return
                        }

                        if (videoFile.type === 'video/x-matroska') {
                            toastStore.info({
                                title: '.mkv file detected',
                                message: 'These files usually take longer to process',
                                duration: 10000
                            })
                        }

                        videos.set([
                            {
                                name: videoFile.name,
                                size: videoFile.size,
                                format: videoFile.type,
                                created_at: Date.now() / 1000,
                                state: 'processing'
                            } as Video,
                            ...$videos
                        ])

                        try {
                            const newVid = await UploadVideo(videoFile)
                            if (!newVid) return

                            videos.set([newVid, ...$videos.slice(1)])
                        } catch (error) {
                            // Remove processing vid if failed
                            videos.set([...$videos.slice(1)])
                            toastStore.error({
                                title: 'Failed to save video to cloud',
                                message: error.message,
                                duration: 10000
                            })
                            console.error('POST /API/FILES', error)
                        }
                    }

                    input.click()
                },
                icon: 'bi-upload',
                class: 'btn-outline-dark'
            },
            {
                text: 'Editor',
                href: '/editor',
                icon: 'bi-pencil',
                class: 'btn-dark'
            }
        ],
        login: () => [
            { text: 'Go Back', href: '/', icon: 'bi-arrow-left', class: 'btn-outline-dark' }
        ],
        verify: () => [
            { text: 'Go Back', href: '/', icon: 'bi-arrow-left', class: 'btn-outline-dark' }
        ],
        register: () => [
            { text: 'Go Back', href: '/', icon: 'bi-arrow-left', class: 'btn-outline-dark' }
        ],
        editor: () => [
            loggedIn
                ? {
                      text: 'Dashboard',
                      href: '/dashboard',
                      icon: 'bi-view-stacked',
                      class: 'btn-gradient'
                  }
                : {
                      text: 'Go Back',
                      href: '/',
                      icon: 'bi-arrow-left',
                      class: 'btn-outline-dark'
                  }
        ]
    }

    const buttons = buttonConfig[subpage ? `${page}/${subpage}` : page]?.() || []
</script>

<header class="border-bottom bg-white-95 sticky-top-custom">
    <div class="container">
        <div class="d-flex align-items-center justify-content-between py-3">
            <div class="d-flex align-items-center">
                <a href="/" class="d-flex align-items-center text-decoration-none me-3">
                    <div
                        class="d-flex align-items-center justify-content-center rounded-2 gradient-primary me-2"
                        style="width: 32px; height: 32px;">
                        <i class="bi bi-play-fill text-white"></i>
                    </div>
                    <span class="fs-4 fw-bold text-dark">vid.sh</span>
                </a>

                {#if title}
                    <div class="vr me-3"></div>
                    <h1 class="fs-5 fw-semibold mb-0">{title}</h1>
                {/if}
            </div>

            <!-- Render buttons dynamically -->
            <div class="d-flex align-items-center gap-3">
                {#each buttons as btn}
                    {#if btn.href}
                        <a href={btn.href} class="btn btn-sm {btn.class}">
                            {#if btn.icon}<i class="bi {btn.icon} me-1"></i>{/if}
                            {btn.text}
                        </a>
                    {:else if btn.action}
                        <button onclick={btn.action} class="btn btn-sm {btn.class}">
                            {#if btn.icon}<i class="bi {btn.icon} me-1"></i>{/if}
                            {btn.text}
                        </button>
                    {/if}
                {/each}
            </div>
        </div>
    </div>
</header>
