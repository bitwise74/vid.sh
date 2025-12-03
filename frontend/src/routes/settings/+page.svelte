<script lang="ts">
    import Header from '$lib/components/Header.svelte'
    import Appearance from '$lib/components/settings/tabs/App.svelte'
    import Danger from '$lib/components/settings/tabs/Danger.svelte'
    import Profile from '$lib/components/settings/tabs/Profile.svelte'
    import Security from '$lib/components/settings/tabs/Security.svelte'
    import { user } from '$lib/stores/AppVars'

    const navTabs = [
        {
            id: 'profile-tab',
            icon: 'bi-person',
            title: 'Profile',
            description: 'Avatar, username, social presets',
            active: true
        },
        {
            id: 'appearance-tab',
            icon: 'bi-window',
            title: 'App',
            description: 'Editor defaults and themes'
        },
        {
            id: 'security-tab',
            icon: 'bi-shield-check',
            title: 'Security',
            description: 'Audit log & safety checks'
        },
        {
            id: 'danger-tab',
            icon: 'bi-exclamation-triangle',
            title: 'Danger zone',
            description: 'Deletion tools'
        }
    ]

    const supportLinks = [
        {
            title: 'Need a hand?',
            description: 'Browse through the commonly asked questions.',
            icon: 'bi-life-preserver',
            href: '/help',
            cta: 'Open help center'
        },
        {
            title: 'Talk to us',
            description: 'Reach out for help or to report a security issue.',
            icon: 'bi-chat-dots',
            href: 'mailto:eryk.sarelo@proton.me',
            cta: 'Email support'
        }
    ]

    const formatBytes = (bytes = 0) => {
        if (!bytes) return '0 B'

        const units = ['B', 'KB', 'MB', 'GB', 'TB']
        let value = bytes
        let unitIndex = 0

        while (value >= 1024 && unitIndex < units.length - 1) {
            value /= 1024
            unitIndex++
        }

        const precision = value >= 10 || unitIndex === 0 ? 0 : 1
        return `${value.toFixed(precision)} ${units[unitIndex]}`
    }

    const getUsagePercent = () => {
        const stats = $user?.stats
        if (!stats?.maxStorage) return 0
        return Math.round((stats.usedStorage / stats.maxStorage) * 100)
    }
</script>

<div class="settings-shell gradient-bg min-vh-100">
    <Header title="Settings" page="settings" />

    <main class="container py-5">
        <section class="hero-panel rounded-4 bg-white-95 shadow-sm p-4 p-lg-5 mb-5">
            <div class="row g-4 align-items-center">
                <div class="col-lg-7">
                    <p class="text-uppercase text-primary fw-semibold small mb-2">Control center</p>
                    <h1 class="fw-semibold mb-3">Tune your preferences</h1>
                    <!-- <p class="text-muted mb-4">Keep your profile, editor defaults, and safety controls in sync.</p> -->

                    <div class="d-flex flex-wrap gap-2">
                        <div class="stat-pill rounded-4">
                            <small class="text-muted">Uploads</small>
                            <strong>{($user?.stats?.uploadedFiles ?? 0).toLocaleString()}</strong>
                        </div>
                        <div class="stat-pill rounded-4">
                            <small class="text-muted">Storage</small>
                            <strong>{formatBytes($user?.stats?.usedStorage ?? 0)} / {formatBytes($user?.stats?.maxStorage ?? 0)}</strong>
                        </div>
                        <div class="stat-pill rounded-4">
                            <small class="text-muted">Visibility</small>
                            <strong>{$user.publicProfileEnabled ? 'Public' : 'Private'}</strong>
                        </div>
                    </div>
                </div>

                <div class="col-lg-5">
                    <div class="insight-card rounded-4 p-4 h-100">
                        <div class="d-flex align-items-center gap-3 mb-3">
                            <span class="badge-soft">
                                <i class="bi-shield-lock"></i>
                            </span>
                            <div>
                                <p class="fw-semibold mb-0">Account health</p>
                                <small class="text-muted">Security & visibility overview</small>
                            </div>
                        </div>

                        <div class="d-flex justify-content-between align-items-center mb-2">
                            <small class="text-muted">Storage used</small>
                            <span class="fw-semibold">{getUsagePercent()}%</span>
                        </div>
                        <div class="progress bg-body-secondary mb-3" style="height: 6px;">
                            <div class="progress-bar bg-primary" role="progressbar" style={`width: ${getUsagePercent()}%`} aria-valuenow={getUsagePercent()} aria-valuemin="0" aria-valuemax="100">
                            </div>
                        </div>
                        <p class="small text-muted mb-0">
                            {getUsagePercent() < 80 ? 'Plenty of room left for uploads.' : 'You are nearing your limit. Consider freeing up space.'}
                        </p>
                    </div>
                </div>
            </div>
        </section>

        <div class="row g-4 align-items-start">
            <div class="col-lg-4 col-xl-3">
                <div class="nav-panel rounded-4 bg-white-95 shadow-sm p-4 mb-4">
                    <p class="text-uppercase text-muted small fw-semibold mb-3">Navigation</p>
                    <div class="d-flex flex-column gap-2" role="tablist">
                        {#each navTabs as tab}
                            <button class={`tab-pill ${tab.active ? 'active' : ''}`} type="button" data-bs-toggle="tab" data-bs-target={`#${tab.id}`} role="tab" aria-controls={tab.id}>
                                <span class="icon">
                                    <i class={`bi ${tab.icon}`}></i>
                                </span>
                                <div class="text-start">
                                    <span class="fw-semibold d-block">{tab.title}</span>
                                    <small class="text-muted">{tab.description}</small>
                                </div>
                            </button>
                        {/each}
                    </div>
                </div>
            </div>

            <div class="col-lg-8 col-xl-9">
                <div class="content-panel rounded-4 bg-white-95 shadow-sm p-4">
                    <div class="tab-content">
                        <Profile />
                        <Appearance />
                        <Security />
                        <Danger />
                    </div>
                </div>
            </div>
        </div>
    </main>
</div>

<style>
    .settings-shell {
        background: linear-gradient(135deg, #f5f3ff 0%, #eef2ff 45%, #e0f2fe 100%);
    }

    .hero-panel {
        border: 1px solid rgba(255, 255, 255, 0.5);
    }

    .stat-pill {
        min-width: 140px;
        padding: 0.85rem 1.25rem;
        border: 1px solid rgba(124, 58, 237, 0.15);
        background: rgba(255, 255, 255, 0.8);
        display: flex;
        flex-direction: column;
        gap: 0.15rem;
    }

    .insight-card {
        background: radial-gradient(circle at top left, rgba(124, 58, 237, 0.08), rgba(59, 130, 246, 0.08));
        border: 1px solid rgba(124, 58, 237, 0.15);
    }

    .badge-soft {
        width: 48px;
        height: 48px;
        border-radius: 16px;
        background: rgba(124, 58, 237, 0.12);
        color: #7c3aed;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        font-size: 1.25rem;
    }

    .nav-panel {
        border: 1px solid rgba(255, 255, 255, 0.5);
    }

    .tab-pill {
        border: 1px solid transparent;
        border-radius: 16px;
        padding: 0.85rem 1rem;
        display: flex;
        align-items: center;
        gap: 0.85rem;
        background: rgba(255, 255, 255, 0.7);
        color: inherit;
        text-align: left;
        transition:
            border-color 0.2s ease,
            transform 0.2s ease;
    }

    .tab-pill:hover {
        border-color: rgba(124, 58, 237, 0.4);
        transform: translateY(-1px);
    }

    .tab-pill.active {
        border-color: rgba(124, 58, 237, 0.6);
        box-shadow: 0 8px 24px rgba(124, 58, 237, 0.1);
        background: #fff;
    }

    .tab-pill .icon {
        width: 40px;
        height: 40px;
        border-radius: 12px;
        background: rgba(124, 58, 237, 0.08);
        color: #7c3aed;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        font-size: 1.1rem;
    }

    .status-item {
        background: rgba(255, 255, 255, 0.8);
    }

    .support-card .support-link {
        color: inherit;
        transition: background 0.2s ease;
    }

    .support-card .support-link:hover {
        background: rgba(124, 58, 237, 0.06);
    }

    .icon-pill {
        width: 42px;
        height: 42px;
        border-radius: 12px;
        background: rgba(59, 130, 246, 0.1);
        color: #2563eb;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        font-size: 1.1rem;
    }

    .content-panel {
        border: 1px solid rgba(255, 255, 255, 0.6);
    }

    @media (max-width: 991.98px) {
        .tab-pill {
            border-radius: 12px;
        }
    }

    @media (max-width: 575.98px) {
        .stat-pill {
            width: 100%;
        }
    }
</style>
