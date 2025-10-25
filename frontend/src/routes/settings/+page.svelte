<script lang="ts">
    import Header from '$lib/components/Header.svelte'
    import Appearance from '$lib/components/settings/tabs/App.svelte'
    import Danger from '$lib/components/settings/tabs/Danger.svelte'
    import Profile from '$lib/components/settings/tabs/Profile.svelte'
    import Security from '$lib/components/settings/tabs/Security.svelte'

    const buttonClass = 'rounded-3 border-0 list-group-item list-group-item-action flex-fill'
</script>

<Header title="Settings" page="settings" />
<div class="p-5 container-lg">
    <div class="row">
        <div class="col-md-3 mb-4">
            <div class="list-group d-flex" role="tablist">
                <button class="{buttonClass} active" data-bs-toggle="tab" data-bs-target="#profile-tab" type="button" role="tab"><i class="bi-person me-1"></i> Profile </button>
                <button class={buttonClass} data-bs-toggle="tab" data-bs-target="#appearance-tab" type="button" role="tab"><i class="bi-window me-1"></i> App </button>
                <button class={buttonClass} data-bs-toggle="tab" data-bs-target="#security-tab" type="button" role="tab"><i class="bi-shield-check me-1"></i> Security </button>
                <button class="{buttonClass} text-danger" data-bs-toggle="tab" data-bs-target="#danger-tab" type="button" role="tab"><i class="bi-exclamation-triangle me-1"></i> Danger Zone </button>
            </div>
        </div>

        <div class="col tab-content">
            <Profile />
            <Appearance />
            <Security />
            <Danger />
        </div>
    </div>
</div>
