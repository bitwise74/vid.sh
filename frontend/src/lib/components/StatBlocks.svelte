<script lang="ts">
    import { stats } from '$lib/stores/VideoStore'


    function convertGB(b: number, point = 2): string {
        return (b/1000000000).toFixed(point)
    }
</script>

<div class="row g-4 mb-5">
    <div class="col-md-6 col-lg-3 d-flex">
        <div class="card border-1 shadow-sm h-100 w-100">
            <div class="card-body d-flex flex-column justify-content-between">
                <div class="d-flex align-items-center justify-content-between">
                    <div>
                        <h6 class="text-muted mb-2">Total Videos</h6>
                        <h3 class="mb-0">{$stats.uploadedFiles}</h3>
                    </div>
                    <i class="bi bi-file-earmark-play text-muted fs-2"></i>
                </div>
            </div>
        </div>
    </div>

    <div class="col-md-6 col-lg-3 d-flex">
        <div class="card border-1 shadow-sm h-100 w-100">
            <div class="card-body d-flex flex-column justify-content-between">
                <div class="d-flex align-items-center justify-content-between">
                    <div class="w-100">
                        <h6 class="text-muted mb-2">Storage Used</h6>
                        <h3 class="mb-2">
                            {convertGB($stats.usedStorage)}/{convertGB($stats.maxStorage, 0) + ' GB'}
                        </h3>
                        <div class="progress mb-2" style="height: 8px;">
                            <div
                                class="progress-bar bg-dark"
                                style="width: {($stats.usedStorage / $stats.maxStorage) * 100}%">
                            </div>
                        </div>
                    </div>
                    <i class="bi bi-hdd text-muted fs-2"></i>
                </div>
            </div>
        </div>
    </div>

    <div class="col-md-6 col-lg-3 d-flex">
        <div class="card border-1 shadow-sm h-100 w-100">
            <div class="card-body d-flex flex-column justify-content-between">
                <div class="d-flex align-items-center justify-content-between">
                    <div>
                        <h6 class="text-muted mb-2">Total Views</h6>
                        <h3 class="mb-0">{$stats.totalViews}</h3>
                    </div>
                    <i class="bi bi-eye text-muted fs-2"></i>
                </div>
            </div>
        </div>
    </div>

    <div class="col-md-6 col-lg-3 d-flex">
        <div class="card border-1 shadow-sm h-100 w-100">
            <div class="card-body d-flex flex-column justify-content-between">
                <div class="d-flex align-items-center justify-content-between">
                    <div>
                        <h6 class="text-muted mb-2">Watch Time</h6>
                        <h3 class="mb-0">{$stats.totalWatchTime}</h3>
                    </div>
                    <i class="bi bi-clock text-muted fs-2"></i>
                </div>
            </div>
        </div>
    </div>
</div>