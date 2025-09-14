<script lang="ts">
    let {
        onVideoSelect,
        videoSelected = false,
        videoName,
        videoSize,
        onClear
    }: {
        onVideoSelect?: (f: File) => void
        onClear: () => void
        videoSelected: boolean
        videoName?: string
        videoSize?: number
    } = $props()
    let isDragOver = $state(false)

    function handleDrop(e: DragEvent) {
        e.preventDefault()
        isDragOver = false

        const files = Array.from(e.dataTransfer!.files)
        const videoFile = files.find((file) => file.type.startsWith('video/'))

        if (videoFile) {
            onVideoSelect?.(videoFile)
        }
    }

    function handleDragOver(e: Event) {
        e.preventDefault()
        isDragOver = true
    }

    function handleDragLeave(e: Event) {
        e.preventDefault()
        isDragOver = false
    }

    function handleFileSelect(e: Event) {
        const file = e.target!.files?.[0]
        if (file && file.type.startsWith('video/')) {
            onVideoSelect?.(file)
        }
    }
</script>

{#if videoSelected}
    <div class="card border-success bg-success bg-opacity-10">
        <div class="card-body p-4">
            <div class="d-flex align-items-center justify-content-between">
                <div class="d-flex align-items-center">
                    <i class="bi bi-file-earmark-play text-success fs-2 me-3"></i>
                    <div>
                        <p class="fw-medium text-success mb-1">{videoName}</p>
                        <p class="small text-success mb-0">
                            {videoSize} MB
                        </p>
                    </div>
                </div>
                <button class="btn btn-outline-success btn-sm" aria-label="Clear video" onclick={onClear}>
                    <i class="bi bi-x"></i>
                </button>
            </div>
        </div>
    </div>
{:else}
    <label
        for="video-upload"
        class="card border-2 border-dashed {isDragOver ? 'border-primary bg-primary bg-opacity-10' : 'border-secondary'}"
        style="cursor: pointer; transition: all 0.3s ease;"
        ondrop={handleDrop}
        ondragover={handleDragOver}
        ondragleave={handleDragLeave}
        aria-label="Upload your video">
        <div class="card-body p-5 text-center">
            <i class="bi bi-cloud-upload text-secondary display-4 mb-4"></i>
            <h5 class="fw-semibold mb-3">Upload your video</h5>
            <p class="text-muted-foreground mb-4">Drag and drop your video file here, or click to browse</p>

            <!-- Hidden input triggered by the wrapping label -->
            <input type="file" accept="video/*" onchange={handleFileSelect} class="d-none" id="video-upload" />

            <div class="btn btn-gradient mt-2">Choose File</div>

            <p class="small text-muted mt-3">MP4 up to 2GB</p>
        </div>
    </label>
{/if}
