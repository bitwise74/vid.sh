<script lang="ts">
    import { toastStore } from '$lib/components/toast/toastStore'
    import "../app.css"
    import Header from '../lib/components/Header.svelte'
</script>

<svelte:head>
  <title>vid.sh - Simple video editor</title>
          <meta property="og:image" content="https://private-user-images.githubusercontent.com/92938606/483901430-92dfbab2-50d9-48f6-b7d0-b15fa62597d2.png?jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3NTY1NzkwNDksIm5iZiI6MTc1NjU3ODc0OSwicGF0aCI6Ii85MjkzODYwNi80ODM5MDE0MzAtOTJkZmJhYjItNTBkOS00OGY2LWI3ZDAtYjE1ZmE2MjU5N2QyLnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFWQ09EWUxTQTUzUFFLNFpBJTJGMjAyNTA4MzAlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjUwODMwVDE4MzIyOVomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPThhOGRmN2RkNWU2M2E2ZTdhY2NhNmM4YjJlY2UxNTkyMWFhOWFhNzEyZDUwNzEyMzUzMjJkZWJhNjQxYWUyMGUmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0In0.H9fbsLwI5hGzoVYHDT3GVOXPDtSyF_WnFZOQFUnFGO8" />
</svelte:head>

<div class="d-flex flex-column min-vh-100">
  <Header page="main"/>

  <main class="flex-grow-1">
    <section class="position-relative py-5 overflow-hidden" style="padding-top: 5rem; padding-bottom: 8rem;">
      <div class="position-absolute top-0 start-0 w-100 h-100 gradient-bg"></div>
      <div class="container position-relative">
        <div class="row justify-content-center">
          <div class="col-lg-10 text-center">
            <h1 class="display-2 fw-bold mb-4">
              Edit, Host and Share Videos
              <span class="d-block text-gradient">Instantly</span>
            </h1>

            <p class="fs-5 text-muted-foreground mb-5 mx-auto" style="max-width: 42rem;">
Edit videos directly in your browser. Store them securely in the cloud. Share anywhere, instantly. 100% free, open source and zero tracking. No downloads, no ads, no bullshit.
            </p>
            
            <hr class="mb-5">

            <div class="row g-4 text-center">
              <div class="col-6 col-md-3">
                <div class="fs-2 fw-bold text-primary">Upload videos</div>
                <div class="small text-muted-foreground">Up to 2GB Each</div>
              </div>
              <div class="col-6 col-md-3">
                <div class="fs-2 fw-bold text-primary">10GB</div>
                <div class="small text-muted-foreground">of Free Cloud Storage</div>
              </div>
              <div class="col-6 col-md-3">
                <div class="fs-2 fw-bold text-primary">Edit videos</div>
                <div class="small text-muted-foreground">Completely Free</div>
              </div>
              <div class="col-6 col-md-3">
                <div class="fs-2 fw-bold text-primary">Share Videos</div>
                <div class="small text-muted-foreground">Instantly, Anywhere</div>
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
            <div class="card h-100 border-0 shadow card-hover">
              <div class="card-body p-4">
                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" 
                     style="width: 48px; height: 48px; background-color: #f3e8ff;">
                  <i class="bi bi-scissors text-primary fs-4"></i>
                </div>
                <h5 class="card-title">Basic Editor</h5>
                <p class="card-text text-muted-foreground">
                  Crop, trim, compress and convert videos directly in your browser. Fast and hassle-free.
                </p>
              </div>
            </div>
          </div>

          <div class="col-md-6 col-lg-4">
            <div class="card h-100 border-0 shadow card-hover">
              <div class="card-body p-4">
                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" 
                     style="width: 48px; height: 48px; background-color: #dcfce7;">
                  <i class="bi bi-cloud text-success fs-4"></i>
                </div>
                <h5 class="card-title">Cloud Hosting</h5>
                <p class="card-text text-muted-foreground">
                  Securely host videos in the cloud for easy sharing without worrying about size limits.
                </p>
              </div>
            </div>
          </div>

          <div class="col-md-6 col-lg-4">
            <div class="card h-100 border-0 shadow card-hover">
              <div class="card-body p-4">
                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" 
                     style="width: 48px; height: 48px; background-color: #fed7aa;">
                  <i class="bi bi-share text-warning fs-4"></i>
                </div>
                <h5 class="card-title">Easy Sharing</h5>
                <p class="card-text text-muted-foreground">
                  Share your videos anywhere with a single click
                </p>
              </div>
            </div>
          </div>

          <div class="col-md-6 col-lg-4">
            <div class="card h-100 border-0 shadow card-hover">
              <div class="card-body p-4">
                <div class="d-flex align-items-center justify-content-center rounded-3 mb-4" 
                     style="width: 48px; height: 48px; background-color: #fecaca;">
                  <i class="bi bi-shield-check text-danger fs-4"></i>
                </div>
                <h5 class="card-title">Open Source</h5>
                <p class="card-text text-muted-foreground">
                        Complete transparency: our backend and frontend are fully open source, so you know exactly what happens behind the scenes.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>

  <footer class="border-top bg-light py-5">
    <div class="container">
      <div class="row g-4">

        <div class="col-md-3 col-sm-4">
          <h6 class="fw-semibold mb-3">Product</h6>
          <ul class="list-unstyled">
            <li class="mb-2"><a href="/editor" class="text-decoration-none text-muted-foreground">Editor</a></li>
            <li class="mb-2"><a href="/dashboard" class="text-decoration-none text-muted-foreground">Dashboard</a></li>
            <li class="mb-2"><a href="/uptime" class="text-decoration-none text-muted-foreground">Uptime</a></li>
          </ul>
        </div>

        <div class="col-md-3 col-sm-4">
          <h6 class="fw-semibold mb-3">About Me</h6>
          <ul class="list-unstyled">
            <li class="mb-2"><a href="/github" class="text-decoration-none text-muted-foreground">My GitHub</a></li>
                        <!-- svelte-ignore a11y_click_events_have_key_events -->
            <li class="mb-2"><p onclick={() => {
                navigator.clipboard.writeText("eryk.sarelo@proton.me")
                toastStore.success({ title: "Email address copied to clipboard!" })
            }} class="text-decoration-none text-muted-foreground" style="cursor: pointer;">Send me an E-mail</p></li>
          </ul>
        </div>

        <div class="col-md-3 col-sm-4">
          <h6 class="fw-semibold mb-3">Support</h6>
          <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
          <ul class="list-unstyled">
            <li class="mb-2"><a href="/help" class="text-decoration-none text-muted-foreground">Commonly Asked Questions</a></li>
            <li class="mb-2"><a href="https://discord.gg/9VyyCkMSFP" class="text-decoration-none text-muted-foreground">Discord server</a></li>

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
