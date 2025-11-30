<script lang="ts">
    import { user } from '$lib/stores/AppVars'

    function convertGB(b: number, point = 2): string {
        return (b / 1000000000).toFixed(point)
    }
</script>

<div class="row g-3 mb-4">
    <div class="col-md-6 col-lg-6">
        <div class="card border-2 rounded-3 shadow-sm stat-card-slim h-100 bg-body-tertiary">
            <div class="card-body py-3 px-4 d-flex align-items-center justify-content-between">
                <div class="d-flex flex-column">
                    <span class="text-muted small text-uppercase fw-semibold">Total Videos</span>
                    <span class="fs-4 fw-bold text-primary mb-0">
                        {$user.stats?.uploadedFiles || 0}
                    </span>
                </div>
                <div class="stat-icon-sm rounded-circle bg-primary bg-opacity-10 text-primary d-flex align-items-center justify-content-center">
                    <i class="bi bi-file-earmark-play-fill fs-5"></i>
                </div>
            </div>
        </div>
    </div>

    <div class="col-md-6 col-lg-6">
        <div class="card border-2 shadow-sm rounded-3 stat-card-slim h-100 bg-body-tertiary">
            <div class="card-body py-3 px-4">
                <div class="d-flex align-items-center justify-content-between mb-1">
                    <span class="text-muted small text-uppercase fw-semibold">Storage Used</span>
                    <div class="stat-icon-sm rounded-circle bg-info bg-opacity-10 text-info d-flex align-items-center justify-content-center">
                        <i class="bi bi-hdd-fill fs-5"></i>
                    </div>
                </div>
                <div class="d-flex align-items-center justify-content-between">
                    <span class="fw-semibold">
                        {convertGB($user.stats?.usedStorage || 0)}
                    </span>
                    <span class="text-muted small">
                        / {convertGB($user.stats?.maxStorage || 0, 0)} GB
                    </span>
                </div>
                <div class="progress mt-2" style="height: 4px;">
                    <div class="progress-bar bg-info" role="progressbar" style="width: {(($user.stats?.usedStorage || 0) / ($user.stats?.maxStorage || 1)) * 100}%"></div>
                </div>
            </div>
        </div>
    </div>
</div>
