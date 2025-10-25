<script lang="ts">
    import { onMount } from 'svelte'

    type EventType = 'login' | 'password_change' | 'profile_update' | 'other'

    type AuditEvent = {
        id: string
        userId: string
        type: EventType
        timestamp: string
        ip?: string
        success?: boolean
        details?: Record<string, unknown>
    }

    let page = $state(0)
    let results = $state(20)

    const ExampleEvents: AuditEvent[] = []

    const typeToColor = {
        login: (success: boolean) => (success ? 'success' : 'danger'),
        password_change: (success: boolean) => (success ? 'warning' : 'danger'),
        profile_update: () => 'info',
        other: () => 'secondary'
    }

    onMount(async () => {})
</script>

<!-- Not really possible since the initla login would be there but you never know -->
{#if ExampleEvents.length === 0}
    <p class="text-muted">No audit log events to display.</p>
{:else}
    <label for="eventType" class="small text-muted">Filter by:</label>
    <div class="select">
        <select class="form-control" id="eventType">
            <option value="all">All Events</option>
            <option value="login">Login Attempts</option>
            <option value="password_change">Password Changes</option>
            <option value="profile_update">Profile Updates</option>
            <option value="other">Other Events</option>
        </select>
    </div>
    <table class="table table-striped mt-3">
        <thead>
            <tr>
                <th scope="col">Timestamp</th>
                <th scope="col">Type</th>
                <th scope="col">IP</th>
                <th scope="col">Details</th>
            </tr>
        </thead>
        <tbody>
            {#each ExampleEvents as event}
                <tr>
                    <td>{new Date(event.timestamp).toLocaleString()}</td>
                    <td class="text-{typeToColor[event.type](event.success || false)}">{event.type.replace('_', ' ')} {event.success == false ? '(failed)' : ''}</td>
                    <td>{event.ip || 'Unknown'}</td>
                    <td>
                        {#if event.details}
                            <pre class="small">{JSON.stringify(event.details, null, 2)}</pre>
                        {:else}
                            Not available
                        {/if}
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
{/if}
