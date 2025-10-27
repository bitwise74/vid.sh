<script lang="ts">
    import { toastStore, type Toast } from '../../stores/ToastStore'

    let { toast }: { toast: Toast } = $props()
    let isLeaving = $state(false)

    function getToastIcon(type: string) {
        switch (type) {
            case 'success':
                return 'bi-check-circle'
            case 'error':
                return 'bi-exclamation-circle'
            case 'warning':
                return 'bi-exclamation-triangle'
            case 'question':
                return 'bi-question-circle'
            default:
                return 'bi-info-circle'
        }
    }

    function getToastColor(type: string) {
        switch (type) {
            case 'success':
                return 'success'
            case 'error':
                return 'danger'
            case 'warning':
                return 'warning'
            case 'loading':
                return 'transparent'
            default:
                return 'black'
        }
    }

    if (toast.duration && toast.duration > 0 && toast.dismissible) {
        setTimeout(() => {
            isLeaving = true
            setTimeout(() => {
                toastStore.remove(toast.id!)
            }, 1500)
        }, toast.duration)
    }
</script>

<div class="toast bg-body-tertiary show align-items-center border-1 {isLeaving ? 'slidefade-out' : 'slidefade-in'}" role="alert" aria-live="assertive" aria-atomic="true">
    <div class="d-flex">
        <div class="toast-body d-flex align-items-center">
            {#if toast.type === 'loading'}
                <div class="spinner-border text-success me-3 spinner-border-sm" role="status"></div>
            {:else}
                <i class="{getToastIcon(toast.type)} fs-5 me-3"></i>
            {/if}

            <!-- text + buttons stacked -->
            <div class="flex-grow-1 d-flex flex-column">
                <div>
                    <div class="fw-semibold">{toast.title}</div>
                    {#if toast.message}
                        {#if toast.message[0] === '<'}
                            {@html toast.message}
                        {:else}
                            {toast.message}
                        {/if}
                    {/if}
                </div>

                {#if toast.buttons && toast.buttons.length > 0}
                    <div class="mt-2 d-flex flex-wrap gap-2">
                        {#each toast.buttons as button}
                            <button
                                type="button"
                                class="btn btn-sm {button.class || 'btn-outline-primary'}"
                                onclick={() => {
                                    button.action()
                                    if (toast.dismissible) toastStore.remove(toast.id)
                                }}>
                                {button.text}
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
        </div>

        {#if toast.dismissible}
            <button type="button" class="btn-close btn-close-dark m-auto me-2" aria-label="Close" onclick={() => toastStore.remove(toast.id)}> </button>
        {/if}
    </div>

    {#if toast.dismissible}
        <div style="--toast-duration:{toast.duration}ms;" class="toast-progress bg-{getToastColor(toast.type)}"></div>
    {/if}
</div>

<style>
    .toast {
        min-width: 300px;
        max-width: 400px;
        position: relative;
        overflow: hidden;
    }

    .slidefade-in {
        animation: slidefade-in 1.5s cubic-bezier(0.1, 1, 0, 1) forwards;
    }

    .slidefade-out {
        animation: slidefade-out 1.5s cubic-bezier(0.1, 1, 0, 1) forwards;
    }

    /* Progress bar */
    .toast-progress {
        position: absolute;
        bottom: 0;
        left: 0;
        height: 3px;
        background-color: rgba(255, 255, 255, 0.8);
        animation: progress-bar var(--toast-duration, 5000ms) linear forwards;
    }

    @keyframes progress-bar {
        from {
            width: 100%;
        }
        to {
            width: 0%;
        }
    }
</style>
