<script lang="ts">
    import { goto } from '$app/navigation'
    import { loadedVideosCount, user } from '$lib/stores/AppVars'
    import { createEventDispatcher, onMount } from 'svelte'

    function convertGB(b: number, point = 2): string {
        return (b / 1000000000).toFixed(point)
    }

    function formatNumber(value: number): string {
        return value.toLocaleString('en-US')
    }

    function storagePercent(): number {
        const used = $user.stats?.usedStorage || 0
        const max = $user.stats?.maxStorage || 1
        if (max === 0) {
            return 0
        }

        return Math.min(100, Math.round((used / max) * 100))
    }

    const dispatch = createEventDispatcher<{ modeChange: boolean }>()

    let collapsed = $state(false)

    function announceCollapse() {
        dispatch('modeChange', collapsed)
    }

    onMount(() => {
        collapsed = localStorage.getItem('statsCollapsed') === 'true'
        announceCollapse()
    })

    function toggleCollapsed() {
        collapsed = !collapsed
        if (typeof localStorage !== 'undefined') {
            localStorage.setItem('statsCollapsed', collapsed ? 'true' : 'false')
        }
        announceCollapse()
    }
</script>

<div class="stats-card rounded-4 p-4 shadow-vlg">
    <div class="d-flex justify-content-end align-items-center gap-2 mb-4 flex-wrap stats-card__header">
        <div class="text-very-muted small d-flex align-items-center gap-2 me-auto">
            <i class="bi bi-speedometer2 text-primary"></i>
            <span>Dashboard stats</span>
        </div>
        <button class="btn btn-sm btn-outline-secondary toggle-compact-btn" onclick={toggleCollapsed}>
            <i class="bi {collapsed ? 'bi-arrows-angle-expand' : 'bi-arrows-angle-contract'} me-1"></i>
            {collapsed ? 'Detailed view' : 'Compact mode'}
        </button>
    </div>

    {#if collapsed}
        <div class="compact-shell">
            <div class="compact-item">
                <span class="compact-label">Storage used</span>
                <span class="compact-value">{convertGB($user.stats?.usedStorage || 0)} GB</span>
                <span class="compact-subtext">{storagePercent()}% of {convertGB($user.stats?.maxStorage || 0, 0)} GB</span>
            </div>
            <div class="compact-divider"></div>
            <div class="compact-item">
                <span class="compact-label">Uploaded videos</span>
                <span class="compact-value">{formatNumber($user.stats?.uploadedFiles || 0)}</span>
                <span class="compact-subtext">All-time uploads</span>
            </div>
            <div class="compact-divider"></div>
            <div class="compact-item">
                <span class="compact-label">Public videos</span>
                <span class="compact-value">{formatNumber($loadedVideosCount || 0)}</span>
                <span class="compact-subtext">Currently published</span>
            </div>
        </div>
    {:else}
        <div class="row g-4">
            <div class="col-lg-8 d-flex flex-column justify-content-between">
                <div class="mb-4">
                    <span class="badge rounded-pill text-primary mb-3 px-3 py-2 stats-badge">WELCOME BACK</span>
                    <h2 class="fw-bold mb-2">Hey Bitwise, ready to edit some videos?</h2>
                    <p class="text-very-muted" style="max-width: 600px;">
                        Processing server is at <span class="text-success">low capacity</span> today. Editing videos should be faster than usual!
                    </p>

                    <div class="d-flex gap-3 mt-4 flex-wrap">
                        <button class="btn btn-primary border-0 fw-semibold px-4 py-2 d-flex align-items-center gap-2 button-gradient" onclick={() => goto('/editor')}>
                            <i class="bi bi-scissors"></i> Start new edit
                        </button>
                        <button
                            class="btn btn-outline-secondary px-4 py-2"
                            onclick={() => {
                                goto('/settings')
                            }}>
                            Change Settings
                        </button>
                    </div>
                </div>

                <div class="row g-3 mt-auto">
                    <div class="col-md-6">
                        <div class="stat-mini-card rounded-4 h-100">
                            <div class="text-muted small text-uppercase fw-bold mb-1">Public videos</div>
                            <div class="d-flex align-items-end justify-content-between">
                                <span class="display-6 fw-bold"> {$loadedVideosCount} </span>
                            </div>
                        </div>
                    </div>

                    <div class="col-md-6">
                        <div class="stat-mini-card rounded-4 h-100">
                            <div class="text-muted small text-uppercase fw-bold mb-1">Total Uploads</div>
                            <div class="d-flex align-items-end justify-content-between">
                                <span class="display-6 fw-bold">
                                    {$user.stats?.uploadedFiles || 0}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="col-lg-4">
                <div class="storage-card rounded-4 h-100">
                    <div class="card-body p-0 d-flex flex-column justify-content-center">
                        <div class="d-flex justify-content-between align-items-center mb-2">
                            <span class="text-primary small text-uppercase fw-bold">Storage Usage</span>
                            <span class="badge bg-primary bg-opacity-25 text-primary rounded-pill">
                                {storagePercent()}%
                            </span>
                        </div>

                        <h2 class="fw-bold mb-3">{convertGB($user.stats?.usedStorage || 0)}</h2>

                        <div class="progress mb-3 storage-progress">
                            <div class="progress-bar" role="progressbar" style={`width: ${storagePercent()}%`}></div>
                        </div>

                        <div class="d-flex justify-content-between text-muted small mb-4 border-bottom pb-4">
                            <span>{convertGB($user.stats?.maxStorage || 0, 0)} GB total</span>
                            <span>{$user.stats?.uploadedFiles || 0} files</span>
                        </div>

                        <div class="d-flex justify-content-between text-muted small mb-1">
                            <span>Default video visibility</span>
                            <span>{$user.defaultPrivateVideos ? 'Private' : 'Public'}</span>
                        </div>
                        <div class="d-flex justify-content-between text-muted small mb-1">
                            <span>Profile visibility</span>
                            <span>{$user.publicProfileEnabled ? 'Private' : 'Public'}</span>
                        </div>
                        <div class="d-flex justify-content-between text-muted small mb-1">
                            <span>Using rich embeds?</span>
                            <span>{localStorage.getItem('richEmbeds') === 'true' ? 'Yes' : 'No'}</span>
                        </div>
                        <div class="d-flex justify-content-between text-muted small">
                            <span>Default video target size</span>
                            <span>{localStorage.getItem('videoTargetSize') || 'Not set'}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    .stats-card {
        background: var(--dashboard-surface-bg);
        border: 1px solid var(--dashboard-surface-border);
        box-shadow: var(--dashboard-surface-shadow);
        color: var(--dashboard-text-primary);
    }

    .stats-badge {
        background: var(--dashboard-chip-bg);
        border: 1px solid var(--dashboard-chip-border);
    }

    .stats-card__header,
    .text-very-muted {
        color: var(--dashboard-text-muted);
    }

    .button-gradient {
        background: linear-gradient(90deg, #7f00ff 0%, #007fff 100%);
        color: #fff;
    }

    .shadow-vlg {
        box-shadow: var(--dashboard-surface-shadow);
    }

    .toggle-compact-btn {
        border-radius: 999px;
        padding-inline: 1rem;
        border-color: var(--dashboard-chip-border);
        color: var(--dashboard-text-primary);
    }

    .compact-shell {
        width: clamp(240px, 60vw, 820px);
        border: 1px solid var(--dashboard-chip-border);
        border-radius: 1.5rem;
        padding: 1rem 1.5rem;
        background: var(--dashboard-card-bg);
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .compact-item {
        display: flex;
        flex-direction: column;
        gap: 0.2rem;
        text-align: center;
        flex: 1;
    }

    .compact-label {
        text-transform: uppercase;
        font-size: 0.7rem;
        letter-spacing: 0.08em;
        color: var(--dashboard-text-muted);
    }

    .compact-value {
        font-size: 1.6rem;
        font-weight: 700;
        line-height: 1.2;
    }

    .compact-subtext {
        font-size: 0.8rem;
        color: var(--dashboard-text-muted);
    }

    .compact-divider {
        width: 1px;
        align-self: stretch;
        background: var(--dashboard-chip-border);
    }

    .stat-mini-card,
    .storage-card {
        background: var(--dashboard-card-bg);
        border: 1px solid var(--dashboard-card-border);
        padding: 1.5rem;
        height: 100%;
    }

    .storage-card {
        padding: 1.5rem;
    }

    .storage-progress {
        height: 8px;
        background: var(--dashboard-progress-bg);
    }

    .storage-progress .progress-bar {
        background: var(--dashboard-progress-fill);
    }

    @media (max-width: 768px) {
        .compact-shell {
            flex-direction: column;
            width: 100%;
        }
    }
</style>
