<script lang="ts">
    import { PUBLIC_CDN_URL } from '$env/static/public'
    import { Logout } from '$lib/api-v2/Auth'
    import { user } from '$lib/stores/AppVars'
    import { UploadFileButton } from '$lib/utils/Upload'
</script>
<div class="dropdown">
    <button class="btn d-flex align-items-center text-decoration-none border-0" type="button" data-bs-toggle="dropdown" aria-expanded="false">
        <img src={$user.avatarHash ? `${PUBLIC_CDN_URL}/avatars/${$user.avatarHash}` : 'placeholder.svg'} alt="User Avatar" class="rounded-circle object-fit-cover me-1" width="40" height="40" />
        <i class="bi-caret-down-fill"></i>
    </button>

    <ul class="dropdown-menu dropdown-menu-end animate">
        {#if $user.publicProfileEnabled}
            <li>
                <a class="dropdown-item" href={`/u/${$user.username}`}>
                    <i class="bi-person me-2"></i>My Profile
                </a>
            </li>
        {/if}
        <li>
            <a class="dropdown-item" href="/settings">
                <i class="bi-gear me-2"></i>Settings
            </a>
        </li>
        <li><hr class="dropdown-divider" /></li>
        <li>
            <button class="dropdown-item" onclick={() => UploadFileButton()}>
                <i class="bi-file-earmark-arrow-up me-2"></i>Upload Video
            </button>
        </li>
        <li>
            <a class="dropdown-item" href="/editor">
                <i class="bi-pencil-square me-2"></i>Editor
            </a>
        </li>
        <li>
            <a class="dropdown-item" href="/dashboard">
                <i class="bi-view-stacked me-2"></i>Dashboard
            </a>
        </li>
        <li><hr class="dropdown-divider" /></li>
        <li style="pointer-events: none;">
            <span class="dropdown-item text-muted small">Signed in as <b>{$user.username}</b></span>
        </li>
        <li>
            <button class="dropdown-item text-danger" onclick={() => Logout()}>
                <i class="bi-box-arrow-right me-2"></i>Log out
            </button>
        </li>
    </ul>
</div>

<style>
    .animate {
        animation: slidefade-in 0.35s cubic-bezier(0.16, 1, 0.3, 1);
    }

    @keyframes slidefade-in {
        from {
            scale: 0.98;
            margin-top: -4px;
            opacity: 0;
        }
        to {
            scale: 1;
            margin-top: 0px;
            opacity: 1;
        }
    }
</style>
