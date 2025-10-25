<script lang="ts">
    import { SearchFiles } from '$lib/api-v2/Files'
    import { perPage } from '$lib/stores/FilterOpts'
    import { view } from '$lib/stores/UserPreferences'
    import { videos } from '$lib/stores/VideoStore'
    import { onDestroy } from 'svelte'

    let timeout: number | undefined

    let { tags }: { tags: any[] } = $props()

    // Debounce requests
    function handleInput(e: any) {
        clearTimeout(timeout)
        const search = e.target.value

        timeout = setTimeout(async () => {
            videos.set(
                await SearchFiles({
                    limit: parseInt($perPage),
                    page: 0, // TODO: implement pagination for searching,
                    query: search
                })
            )
        }, 300)
    }

    function toggleView(val: string) {
        localStorage.setItem('view', val)
        view.set(val)
    }

    onDestroy(() => {
        clearTimeout(timeout)
    })
</script>
<div class="row g-3 mb-4">
    <div class="d-flex justify-content-center">
        <div class="input-group shadow-sm" style="max-width: 500px;">
            <span class="input-group-text bg-opacity-75">
                <i class="bi-search"></i>
            </span>
            <input type="text" class="form-control" placeholder="Search videos..." oninput={handleInput} />
            <button disabled class="input-group-text dropdown-toggle bg-opacity-75" aria-label="filters" data-bs-toggle="dropdown" aria-expanded="false" data-bs-auto-close="outside">
                <i class="bi-funnel"></i>
            </button>
            <div class="dropdown-center">
                <ul class="dropdown-menu p-3 animate" style="min-width: 600px; max-width: 800px;">
                    <li>
                        <div class="row g-3">
                            <!-- Filter By -->
                            <div class="col-md-4">
                                <label class="form-label small text-muted" for="filter-by">Filter by</label>
                                <select class="form-select" id="filter-by"> </select>
                            </div>

                            <!-- Sort By -->
                            <div class="col-md-4">
                                <label class="form-label small text-muted" for="sort-by">Sort by</label>
                                <select class="form-select" id="sort-by">
                                    <option value="latest">Latest</option>
                                    <option value="popular">Most Popular</option>
                                    <option value="rating">Highest Rated</option>
                                </select>
                            </div>

                            <div class="col-md-4">
                                <label class="form-label small text-muted" for="results-per-page">Results per page</label>
                                <select bind:value={$perPage} class="form-select" id="results-per-page">
                                    <option value="10">10</option>
                                    <option value="20" selected>20</option>
                                    <option value="50">50</option>
                                    <option value="100">100</option>
                                </select>
                            </div>

                            <div class="col-md-12">
                                <p class="form-label small text-muted user-select-none mb-2">Tags</p>
                                <div class="p-2 border rounded" style="max-height: 120px; overflow-y: auto;">
                                    {#each tags as tag, i}
                                        <input class="form-check-input" type="checkbox" value="1" id={`tag${i}`} />
                                        <label class="form-check-input" for={`tag${i}`}>
                                            {tag.name} <span class="text-muted">({tag.assignedCount})</span>
                                        </label>
                                    {/each}
                                </div>
                            </div>
                        </div>
                    </li>

                    <li class="mt-3 pt-2 border-top">
                        <div class="d-flex justify-content-end gap-2">
                            <button type="button" class="btn btn-sm btn-outline-secondary">Clear</button>
                            <button type="button" class="btn btn-sm btn-primary">Apply</button>
                        </div>
                    </li>
                </ul>
            </div>
        </div>

        <div class="input-group-text ms-3 bg-opacity-75">
            <button
                type="button"
                class="btn bg-transparent border-0 m-0 p-1"
                data-bs-toggle="tooltip"
                aria-label="Toggle list view"
                data-placement="top"
                title="Toggle list view"
                onclick={() => toggleView('list')}>
                <i class="bi-list-ul"></i>
            </button>
            <div class="vr m-2"></div>
            <button
                type="button"
                class="btn bg-transparent border-0 m-0 p-1"
                data-bs-toggle="tooltip"
                aria-label="Toggle grid view"
                data-placement="top"
                title="Toggle list view"
                onclick={() => toggleView('grid')}>
                <i class="bi-grid-3x3"></i>
            </button>
        </div>
    </div>
</div>

<style>
    .animate {
        animation: slidefade-in 0.35s cubic-bezier(0.16, 1, 0.3, 1);
    }
</style>
