<script lang="ts">
    import { loadedVideosCount, user } from '$lib/stores/AppVars'

    function convertGB(b: number, point = 2): string {
        return (b / 1000000000).toFixed(point)
    }
</script>
<div class="card bg-body-tertiary text-white border-0 rounded-4 p-4 mb-4">
    <div class="row g-4">
        <div class="col-lg-8 d-flex flex-column justify-content-between">
            <div class="mb-4">
                <span class="badge rounded-pill bg-primary bg-opacity-25 text-primary mb-3 px-3 py-2">WELCOME BACK</span>
                <h2 class="fw-bold mb-2">Hey Bitwise, ready to edit some videos?</h2>
                <p class="text-very-muted" style="max-width: 600px;">Processing server is at low capacity today. Editing videos should be faster than usual!</p>

                <div class="d-flex gap-3 mt-4">
                    <button class="btn btn-primary border-0 fw-semibold px-4 py-2 d-flex align-items-center gap-2 button-gradient">
                        <i class="bi bi-scissors"></i> Start new edit
                    </button>
                    <button class="btn btn-outline-secondary px-4 py-2"> Manage library </button>
                </div>
            </div>

            <div class="row g-3 mt-auto">
                <div class="col-md-6">
                    <div class="card bg-black bg-opacity-25 border border-dark-gray border-opacity-25 rounded-4 h-100">
                        <div class="card-body">
                            <div class="text-muted small text-uppercase fw-bold mb-1">Videos This Session</div>
                            <div class="d-flex align-items-end justify-content-between">
                                <span class="display-6 fw-bold"> {$loadedVideosCount} </span>
                                <span class="text-primary small mb-2">Loaded in view</span>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="col-md-6">
                    <div class="card bg-black bg-opacity-25 border border-dark-gray border-opacity-25 rounded-4 h-100">
                        <div class="card-body">
                            <div class="text-primary small text-uppercase fw-bold mb-1">Total Uploads</div>
                            <div class="d-flex align-items-end justify-content-between">
                                <span class="display-6 fw-bold text-white">
                                    {$user.stats?.uploadedFiles || 0}
                                </span>
                                <span class="text-primary small mb-2">Across account</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="col-lg-4">
            <div class="card bg-black bg-opacity-25 order border-dark-gray border-opacity-25 rounded-4 h-100">
                <div class="card-body p-4 d-flex flex-column justify-content-center">
                    <div class="d-flex justify-content-between align-items-center mb-2">
                        <span class="text-primary small text-uppercase fw-bold">Storage Usage</span>
                        <span class="badge bg-primary bg-opacity-25 text-primary rounded-pill">
                            {Math.round((($user.stats?.usedStorage || 0) / ($user.stats?.maxStorage || 1)) * 100)}%
                        </span>
                    </div>

                    <h2 class="fw-bold mb-3">{convertGB($user.stats?.usedStorage || 0)}</h2>

                    <div class="progress bg-secondary bg-opacity-25 mb-3" style="height: 8px;">
                        <div class="progress-bar bg-primary" role="progressbar" style="width: {(($user.stats?.usedStorage || 0) / ($user.stats?.maxStorage || 1)) * 100}%"></div>
                    </div>

                    <div class="d-flex justify-content-between text-primary small mb-4 border-bottom border-secondary border-opacity-25 pb-4">
                        <span>{convertGB($user.stats?.maxStorage || 0, 0)} GB total</span>
                        <span>{$user.stats?.uploadedFiles || 0} files</span>
                    </div>

                    <div class="d-flex justify-content-between text-primary small mb-2">
                        <span>Default video visibility</span>
                        <span class="text-white">{$user.defaultPrivateVideos ? 'Private' : 'Public'}</span>
                    </div>
                    <div class="d-flex justify-content-between text-primary small">
                        <span>Profile visibility</span>
                        <span class="text-white">{$user.publicProfileEnabled ? 'Private' : 'Public'}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .text-very-muted {
        color: rgba(255, 255, 255, 0.4);
    }

    .border-dark-gray {
        border-color: rgba(255, 255, 255, 0.1) !important;
    }

    .button-gradient {
        background: linear-gradient(90deg, #7f00ff 0%, #007fff 100%);
    }
</style>
