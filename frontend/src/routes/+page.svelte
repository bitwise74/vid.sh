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
    <title>vid.sh - Simple video editor</title>
    <meta property="og:title" content="Vid.sh - Suckless video editor" />
    <meta
        property="og:description"
        content="Edit your videos right in your browser, store them securely in the cloud, and share them anywhere. All for free and open source. No downloads, no ads, no spying and no bullshit." />
    <meta property="og:url" content="https://bitwise0x.dev" />
    <meta property="theme-color" content="#5733E7" />
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
                            Edit, Host & Share
                            <span class="d-block text-gradient">Videos Easily</span>
                        </h1>

                        <p class="fs-5 text-muted-foreground mb-5 mx-auto" style="max-width: 42rem;">
                            Edit your videos right in your browser, store them securely in the cloud, and share them anywhere. All for free and open source. No downloads, no ads, no spying and no
                            bullshit.
                        </p>

                        <div class="d-flex flex-column flex-sm-row gap-3 justify-content-center mb-5">
                            <a href="/editor" class="btn btn-dark btn-lg px-4">
                                <i class="bi-pencil me-2"></i>
                                Open Editor
                            </a>
                        </div>

                        <div class="row g-4 text-center">
                            <div class="col-6 col-md-3">
                                <div class="fs-2 fw-bold text-primary">2GB</div>
                                <div class="small text-muted-foreground">File Size Limit</div>
                            </div>
                            <div class="col-6 col-md-3">
                                <div class="fs-2 fw-bold text-primary">10GB</div>
                                <div class="small text-muted-foreground">Free Storage</div>
                            </div>
                            <div class="col-6 col-md-3">
                                <div class="fs-2 fw-bold text-primary">Free</div>
                                <div class="small text-muted-foreground">Video Editor</div>
                            </div>
                            <div class="col-6 col-md-3">
                                <div class="fs-2 fw-bold text-primary">Free</div>
                                <div class="small text-muted-foreground">Video Sharing anywhere</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <section id="features" class="py-5" style="padding-top: 5rem; padding-bottom: 8rem;">
            <div class="container">
                <div class="row justify-content-center text-center mb-5">
                    <div class="col-lg-8">
                        <h2 class="display-4 fw-bold mb-4">Everything You Need to Quickly Edit Videos</h2>
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
                                <p class="card-text text-muted-foreground">Simple stuff like cropping, trimming, compressing and converting your videos is done completely on-device.</p>
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
                                <p class="card-text text-muted-foreground">Registered users can host their videos in the cloud to share them easily on platforms with file size limits.</p>
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
                                <p class="card-text text-muted-foreground">Sharing your cloud-hosted videos is as easy as clicking a button and embedding the video anywhere you want.</p>
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
                </div>
            </div>
        </section>
    </main>

    <footer class="border-top py-5">
        <div class="container">
            <div class="row g-4">
                <div class="col-md-3 col-sm-4">
                    <h6 class="fw-semibold mb-3">Product</h6>
                    <ul class="list-unstyled">
                        <li class="mb-2">
                            <a href="/video-editor" class="text-decoration-none text-muted-foreground">Video Editor</a>
                        </li>
                        <li class="mb-2">
                            <a href="/dashboard" class="text-decoration-none text-muted-foreground">Dashboard</a>
                        </li>
                        <li class="mb-2">
                            <a href="/analytics" class="text-decoration-none text-muted-foreground">Analytics</a>
                        </li>
                        <li class="mb-2">
                            <a href="/api" class="text-decoration-none text-muted-foreground">API</a>
                        </li>
                    </ul>
                </div>

                <div class="col-md-3 col-sm-4">
                    <h6 class="fw-semibold mb-3">Company</h6>
                    <ul class="list-unstyled">
                        <li class="mb-2">
                            <a href="/about" class="text-decoration-none text-muted-foreground">About</a>
                        </li>
                        <li class="mb-2">
                            <a href="/blog" class="text-decoration-none text-muted-foreground">Blog</a>
                        </li>
                        <li class="mb-2">
                            <a href="/careers" class="text-decoration-none text-muted-foreground">Careers</a>
                        </li>
                        <li class="mb-2">
                            <a href="/contact" class="text-decoration-none text-muted-foreground">Contact</a>
                        </li>
                    </ul>
                </div>

                <div class="col-md-3 col-sm-4">
                    <h6 class="fw-semibold mb-3">Support</h6>
                    <ul class="list-unstyled">
                        <li class="mb-2">
                            <a href="/help-center" class="text-decoration-none text-muted-foreground">Help Center</a>
                        </li>
                        <li class="mb-2">
                            <a href="/documentation" class="text-decoration-none text-muted-foreground">Documentation</a>
                        </li>
                        <li class="mb-2">
                            <a href="/community" class="text-decoration-none text-muted-foreground">Community</a>
                        </li>
                        <li class="mb-2">
                            <a href="/status" class="text-decoration-none text-muted-foreground">Status</a>
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
