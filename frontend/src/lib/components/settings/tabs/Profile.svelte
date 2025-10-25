<script lang="ts">
    import { PUBLIC_BASE_URL, PUBLIC_CDN_URL } from '$env/static/public'
    import { UpdateUser } from '$lib/api-v2/User'
    import { user } from '$lib/stores/AppVars'
    import { toastStore } from '$lib/stores/ToastStore'

    let usernameFieldState = $state<'idling' | 'taken' | 'error' | 'success'>('idling')
    let timeout

    function handleUploadAvatar() {
        const input = document.createElement('input')
        input.type = 'file'
        input.accept = 'image/*'
        input.click()

        input.onchange = async (e) => {
            var file = (e.target as HTMLInputElement).files![0]

            // 10MB hard limit
            if (file.size > 10000000) {
                toastStore.error({
                    title: 'File is too big',
                    message: 'Max upload size for profile pictures is 10MB'
                })
                return
            }

            // TODO: add cropping

            const profile = await UpdateUser({
                avatar: file
            })

            // We know for sure its safe
            user.update((u) => ({ ...(u as any), avatarHash: profile.avatarHash! }))
        }

        input.onclose = () => {
            document.removeChild(input)
        }
    }

    function handleUpdateUsername(e: Event) {
        clearTimeout(timeout)

        timeout = setTimeout(() => {
            const newUsername = (e.target as HTMLInputElement).value

            UpdateUser({
                username: newUsername
            })
                .then(() => {
                    usernameFieldState = 'success'
                    user.update((u) => ({ ...(u as any), username: newUsername }))
                })
                .catch((err) => {
                    if (err.status === 409) {
                        usernameFieldState = 'taken'
                        return
                    }

                    usernameFieldState = 'error'

                    toastStore.error({
                        title: 'Failed to update username',
                        message: err.message || 'Check the console for details'
                    })
                })
        }, 1000)
    }

    function toggleProfile(e: Event) {
        const enabled = (e.target as HTMLInputElement).checked

        UpdateUser({
            publicProfileEnabled: enabled
        })
            .then(() => {
                user.update((u) => ({ ...(u as any), publicProfileEnabled: enabled }))

                if (!enabled) {
                    toastStore.success({
                        title: 'Public profile disabled'
                    })
                    return
                }

                toastStore.success({
                    title: 'Done!',
                    message: `<p>Your public profile is available <a href=${PUBLIC_BASE_URL}/u/${$user.username}>here</a></p>`
                })
            })
            .catch((err) => {
                toastStore.error({
                    title: 'Failed to update public profile setting',
                    message: err.message || 'Check the console for details'
                })
            })
    }
</script>

<div class="tab-pane show active" id="profile-tab">
    <h4>Public profile</h4>
    <hr />
    <div class="row">
        <div class="col-sm-11 col-md-6 col-lg-4 order-1 order-md-2 d-flex justify-content-md-end pe-4 mb-4">
            <div class="position-relative" style="width: 200px; height: 200px">
                <img
                    src={$user.avatarHash ? `${PUBLIC_CDN_URL}/avatars/${$user.avatarHash}` : 'placeholder.svg'}
                    alt="profile_picture"
                    class="rounded-circle object-fit-cover"
                    width="200"
                    height="200" />
                <button class="btn btn-sm btn-dark position-absolute bottom-0 start-0 m-2" onclick={() => handleUploadAvatar()}>
                    <i class="bi-pencil me-1"></i> Edit
                </button>
            </div>
        </div>

        <div class="col-md-6 col-lg-8 order-2 order-md-1">
            <label for="username" class="form-label fw-semibold">Username</label>

            <div class="input-group">
                <span
                    class="input-group-text border-1
      {usernameFieldState === 'taken' ? 'border-danger' : ''}
      {usernameFieldState === 'error' ? 'border-danger' : ''}
      {usernameFieldState === 'success' ? 'border-success' : ''}"
                    style="border-right: 0;">@</span>

                <input
                    type="text"
                    id="username"
                    class="form-control border-1
      {usernameFieldState === 'taken' ? 'border-danger' : ''}
      {usernameFieldState === 'error' ? 'border-danger' : ''}
      {usernameFieldState === 'success' ? 'border-success' : ''}"
                    value={$user.username}
                    onchange={handleUpdateUsername}
                    required />
            </div>

            {#if usernameFieldState === 'taken'}
                <div class="small text-danger mt-2">Username is taken</div>
            {:else if usernameFieldState === 'error'}
                <div class="small text-danger mt-2">Failed to update username</div>
            {:else if usernameFieldState === 'success'}
                <div class="small text-success mt-2">Username updated!</div>
            {/if}

            <p class="small text-muted mt-2">
                Your username will be shown in rich embeds and will be used for your @handle. You only need to set this if you want rich embeds or your public profile to work
            </p>
            <!-- <div class="col-md-6 col-lg-8 order-3 order-md-2"> -->
            <input id="lossless_export" class="form-check-input" type="checkbox" checked={$user.publicProfileEnabled} onclick={(e) => toggleProfile(e)} />
            <label for="lossless_export" class="form-label fw-semibold">Enable public profile</label>
            <p class="small text-muted">If toggled, a public profile will be created on which others can view videos you uploaded that are public. To enable this you need to set a username</p>
            <!-- </div> -->
        </div>
    </div>
</div>

<style>
    .form-control:focus {
        box-shadow: none !important;
        border-color: inherit !important;
    }
</style>
