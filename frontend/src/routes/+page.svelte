<script lang="ts">
    import { goto } from '$app/navigation'
    import { onMount } from 'svelte'
    import '../app.css'
    import Header from '../lib/components/Header.svelte'

    onMount(() => {
        if (localStorage.getItem('optDisableRoot') === 'true') {
            goto('/dashboard')
        }
    })
</script>

<svelte:head>
    <title>Vid.sh - Online Video Editor</title>
    <meta property="og:title" content="Vid.sh - Online Video Editor" />
    <meta name="description" content="Vid.sh is a free online video editor. Crop, trim, compress, convert, and share videos securely in the cloud, with no downloads or ads." />
    <meta name="keywords" content="online video editor, video cropping, video trimming, video conversion, free video editor, cloud video editor" />
    <meta property="og:url" content="https://bitwise0x.dev" />
    <meta property="og:type" content="website" />
    <meta name="twitter:title" content="Vid.sh - Online Video Editor" />
    <meta name="twitter:description" content="Edit your videos directly in the browser, store them securely in the cloud, and share anywhere." />
    <meta property="theme-color" content="#5733E7" />
    <link rel="canonical" href="https://bitwise0x.dev/" />

    <script type="application/ld+json">
        {
            "@context": "https://schema.org",
            "@type": "SoftwareApplication",
            "name": "Vid.sh",
            "url": "https://bitwise0x.dev",
            "applicationCategory": "Multimedia",
            "operatingSystem": "Web",
            "offers": {
                "@type": "Offer",
                "price": "0",
                "priceCurrency": "USD"
            },
            "description": "Vid.sh is a free online video editor. Crop, trim, compress, convert, and share videos in the cloud with no ads."
        }
    </script>
</svelte:head>

<div class="d-flex flex-column min-vh-100">
    <Header page="main" />

    <main class="flex-grow-1">
        <section class="position-relative py-5 overflow-hidden" style="padding-top: 5rem; padding-bottom: 8rem;">
            <div class="position-absolute top-0 start-0 w-100 h-100 gradient-bg"></div>
            <div class="container position-relative">
                <div class="row justify-content-center">
                    <div class="col-lg-10 text-center">
                        <h1 class="display-2 fw-bold mb-4">
                            Editing Clips Shouldn't Be
                            <span class="d-block text-gradient">Complicated</span>
                        </h1>

                        <p class="fs-5 text-muted-foreground mb-5 mx-auto" style="max-width: 42rem;">
                            Imagine a world where you can simply open your browser, edit your videos, export or store them in the cloud and be done. No downloads, no ads, no spying, no bullshit.
                        </p>

                        <div class="d-flex flex-column flex-sm-row gap-3 justify-content-center mb-5">
                            <a href="/editor" class="btn btn-gradient btn-lg px-4">
                                <i class="bi-pencil me-2"></i>
                                Editor
                            </a>
                            <a href="/dashboard" class="btn btn-gradient btn-lg px-4">
                                <i class="bi bi-grid me-2"></i>
                                Dashboard
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <section id="features" class="py-5" style="padding-top: 5rem; padding-bottom: 8rem;">
            <div class="container">
                <div class="row justify-content-center text-center mb-5">
                    <div class="col-lg-8">
                        <h2 class="display-4 fw-bold mb-4 fs-1">All The Features You Need <br /><span class="text-gradient">In One Place</span></h2>
                    </div>
                </div>

                <div class="row g-4">
                    <div class="col-md-6 col-lg-4">
                        <div class="card h-100 border-0 shadow card-hover bg-body-tertiary">
                            <div class="card-body p-4">
                                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" style="width: 48px; height: 48px; background-color: #f3e8ff;">
                                    <i class="bi bi-scissors text-primary fs-4"></i>
                                </div>
                                <h5 class="card-title">Basic Editor</h5>
                                <p class="card-text text-muted-foreground">Simple stuff like cropping, trimming, compressing and converting your videos.</p>
                            </div>
                        </div>
                    </div>

                    <div class="col-md-6 col-lg-4">
                        <div class="card h-100 border-0 shadow card-hover bg-body-tertiary">
                            <div class="card-body p-4">
                                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" style="width: 48px; height: 48px; background-color: #dcfce7;">
                                    <i class="bi bi-cloud text-success fs-4"></i>
                                </div>
                                <h5 class="card-title">Cloud Hosting</h5>
                                <p class="card-text text-muted-foreground">
                                    Registered users can host their videos in the cloud to share them easily on other platforms to bypass their file size limits.
                                </p>
                            </div>
                        </div>
                    </div>

                    <div class="col-md-6 col-lg-4">
                        <div class="card h-100 border-0 shadow card-hover bg-body-tertiary">
                            <div class="card-body p-4">
                                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" style="width: 48px; height: 48px; background-color: #fed7aa;">
                                    <i class="bi bi-share text-warning fs-4"></i>
                                </div>
                                <h5 class="card-title">Easy Sharing</h5>
                                <p class="card-text text-muted-foreground">Sharing your cloud-hosted videos is as easy as clicking a button and sharing it anywhere you want.</p>
                            </div>
                        </div>
                    </div>

                    <div class="col-md-6 col-lg-4">
                        <div class="card h-100 border-0 shadow card-hover bg-body-tertiary">
                            <div class="card-body p-4">
                                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" style="width: 48px; height: 48px; background-color: #fecaca;">
                                    <i class="bi bi-shield-check text-danger fs-4"></i>
                                </div>
                                <h5 class="card-title">Open Source</h5>
                                <p class="card-text text-muted-foreground">
                                    Everything from the backend to the frontend is completely open-sourced. This means you know exactly what happens with your data.
                                </p>
                            </div>
                        </div>
                    </div>

                    <div class="col-md-6 col-lg-4">
                        <div class="card h-100 border-0 shadow card-hover bg-body-tertiary">
                            <div class="card-body p-4">
                                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" style="width: 48px; height: 48px; background-color: #d8b4fe;">
                                    <i class="bi-shield text-purple fs-4"></i>
                                </div>
                                <h5 class="card-title">Privacy Focused</h5>
                                <p class="card-text text-muted-foreground">
                                    We take your privacy seriously. After processing your videos nothing is stored on the server permanently. There's absolutely no telemetry or other spying involved.
                                </p>
                            </div>
                        </div>
                    </div>

                    <div class="col-md-6 col-lg-4">
                        <div class="card h-100 border-0 shadow card-hover bg-body-tertiary">
                            <div class="card-body p-4">
                                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" style="width: 48px; height: 48px; background-color: #bbf7d0;">
                                    <i class="bi bi-wallet2 text-success fs-4"></i>
                                </div>
                                <h5 class="card-title">Free to Use</h5>
                                <p class="card-text text-muted-foreground">
                                    Vid.sh is completely free to use. There are no hidden charges and no premium plans. Just open the editor and start editing your videos right away. Although
                                    donations are appreciated to keep the project running.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    </main>

    <footer class="border-top py-5">
        <div class="container">
            <div class="row g-4">
                <div class="col-md-3 col-sm-4">
                    <h6 class="fw-semibold mb-3">Vid.sh</h6>
                    <ul class="list-unstyled">
                        <li class="mb-2">
                            <a href="/editor" class="text-decoration-none text-muted-foreground">Editor</a>
                        </li>
                        <li class="mb-2">
                            <a href="/dashboard" class="text-decoration-none text-muted-foreground">Dashboard</a>
                        </li>
                        <li class="mb-2">
                            <a href="/analytics" class="text-decoration-none text-muted-foreground">Analytics</a>
                        </li>
                        <li class="mb-2">
                            <a href="https://github.com/bitwise74/vid.sh" target="_blank" class="text-decoration-none text-muted-foreground">Source Code</a>
                        </li>
                    </ul>
                </div>

                <div class="col-md-3 col-sm-4">
                    <h6 class="fw-semibold mb-3">About The Creator</h6>
                    <ul class="list-unstyled">
                        <li class="mb-2">
                            <a href="https://bitwise74.github.io/bitwise74/" target="_blank" class="text-decoration-none text-muted-foreground">Portfolio</a>
                        </li>
                        <li class="mb-2">
                            <a href="https://github.com/bitwise74" target="_blank" class="text-decoration-none text-muted-foreground">My GitHub</a>
                        </li>
                        <li class="mb-2">
                            <a href="https://github.com/sponsors/bitwise74" target="_blank" class="text-decoration-none text-muted-foreground">Support Me</a>
                        </li>
                    </ul>
                </div>

                <div class="col-md-3 col-sm-4">
                    <h6 class="fw-semibold mb-3">Support</h6>
                    <ul class="list-unstyled">
                        <li class="mb-2">
                            <a href="/help" class="text-decoration-none text-muted-foreground">Help</a>
                        </li>
                        <li class="mb-2">
                            <a href="https://discord.gg/9VyyCkMSFP" target="_blank" class="text-decoration-none text-muted-foreground">Discord</a>
                        </li>
                    </ul>
                </div>
            </div>

            <div class="border-top pt-4 mt-4 d-flex flex-column flex-md-row justify-content-between align-items-center">
                <p class="small text-muted-foreground mb-0">© 2025 vid.sh. All rights reserved.</p>
                <div class="d-flex gap-4 mt-3 mt-md-0">
                    <a href="/privacy-policy" class="small text-decoration-none text-muted-foreground">Privacy Policy</a>
                    <a href="/terms-of-service" class="small text-decoration-none text-muted-foreground">Terms of Service</a>
                </div>
            </div>
        </div>
    </footer>
</div>

<style>
    .text-purple {
        color: #a855f7;
    }
</style>
