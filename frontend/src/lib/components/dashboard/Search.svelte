<script lang="ts">
    import { SearchFiles } from '$lib/api/Files'
    import { dashboardView, perPage } from '$lib/stores/appControl'
    import { videos } from '$lib/stores/VideoStore'
    import { onDestroy } from 'svelte'

    let timeout: number | undefined

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

    function toggleView(val: 'list' | 'grid') {
        localStorage.setItem('view', val)
        dashboardView.set(val)
    }

    onDestroy(() => {
        clearTimeout(timeout)
    })
</script>

<div class="search-panel">
    <div class="search-panel__primary">
        <label class="search-field" for="dashboard-search">
            <i class="bi bi-search"></i>
            <input id="dashboard-search" type="text" placeholder="Search your library" oninput={handleInput} />
        </label>

        <button type="button" class="chip-button" disabled>
            <i class="bi bi-funnel"></i>
            Filters soon
        </button>
    </div>

    <div class="search-panel__secondary">
        <div class="control-block">
            <span class="control-label">Results</span>
            <select bind:value={$perPage} aria-label="Results per page">
                <option value="10">10 / page</option>
                <option value="20">20 / page</option>
                <option value="50">50 / page</option>
                <option value="100">100 / page</option>
            </select>
        </div>

        <div class="view-toggle" role="group" aria-label="Toggle layout">
            <button type="button" class:selected={$dashboardView === 'list'} onclick={() => toggleView('list')}>
                <i class="bi bi-list-ul"></i>
                <span>List</span>
            </button>
            <button type="button" class:selected={$dashboardView === 'grid'} onclick={() => toggleView('grid')}>
                <i class="bi bi-grid-3x3"></i>
                <span>Grid</span>
            </button>
        </div>
    </div>
</div>

<style>
    .search-panel {
        border: 1px solid var(--dashboard-surface-border);
        border-radius: 1.5rem;
        padding: 1.5rem;
        background: var(--dashboard-surface-bg);
        box-shadow: var(--dashboard-surface-shadow);
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
        color: var(--dashboard-text-primary);
    }

    .search-panel__primary,
    .search-panel__secondary {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
        align-items: center;
        justify-content: space-between;
    }

    .search-field {
        flex: 1;
        display: flex;
        align-items: center;
        gap: 0.75rem;
        border-radius: 999px;
        padding: 0.75rem 1.25rem;
        border: 1px solid var(--dashboard-chip-border);
        color: var(--dashboard-text-muted);
        cursor: text;
        background: var(--dashboard-card-bg);
    }

    .search-field input {
        flex: 1;
        border: none;
        background: transparent;
        color: var(--dashboard-text-primary);
    }

    .search-field input:focus {
        outline: none;
    }

    .chip-button {
        border-radius: 999px;
        border: 1px solid var(--dashboard-chip-border);
        padding: 0.65rem 1.2rem;
        background: var(--dashboard-chip-bg);
        color: var(--dashboard-text-muted);
        display: inline-flex;
        align-items: center;
        gap: 0.4rem;
        font-size: 0.85rem;
    }

    .control-block {
        display: flex;
        flex-direction: column;
        gap: 0.2rem;
        color: var(--dashboard-text-muted);
    }

    .control-block select {
        background: var(--dashboard-card-bg);
        color: var(--dashboard-text-primary);
        border: 1px solid var(--dashboard-chip-border);
        border-radius: 0.75rem;
        padding: 0.4rem 0.8rem;
        min-width: 140px;
    }

    .view-toggle {
        display: inline-flex;
        border: 1px solid var(--dashboard-chip-border);
        border-radius: 999px;
        overflow: hidden;
        background: var(--dashboard-chip-bg);
    }

    .view-toggle button {
        border: none;
        background: transparent;
        color: var(--dashboard-text-muted);
        padding: 0.5rem 1rem;
        font-size: 0.85rem;
        display: flex;
        align-items: center;
        gap: 0.4rem;
        cursor: pointer;
    }

    .view-toggle button.selected {
        background: linear-gradient(90deg, #7f00ff, #007fff);
        color: #fff;
    }

    :global(.stack-collapsed .dashboard-stack__search .search-panel) {
        margin-bottom: 0;
    }

    @media (max-width: 768px) {
        .search-panel__secondary {
            flex-direction: column;
            align-items: flex-start;
        }

        .view-toggle {
            width: 100%;
            justify-content: space-between;
        }

        .view-toggle button {
            flex: 1;
            justify-content: center;
        }
    }
</style>
