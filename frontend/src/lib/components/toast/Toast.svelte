<script lang="ts">
    import { toastStore, type Toast } from './toastStore'

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
            case 'info':
                return 'bi-info-circle'
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
            case 'info':
                return 'black'
            default:
                return 'black'
        }
    }
</script>

<div class="toast show align-items-center text-bg-light border-1 {isLeaving ? 'slidefade-out' : 'slidefade-in'}" role="alert" aria-live="assertive" aria-atomic="true">
    <div class="d-flex">
        <div class="toast-body d-flex align-items-center">
            <i class="{getToastIcon(toast.type)} me-3 fs-5"></i>
            <div class="flex-grow-1">
                <div class="fw-semibold">{toast.title}</div>
                {#if toast.message}
                    <div class="small opacity-75">{toast.message}</div>
                {/if}
            </div>
        </div>
        {#if toast.dismissible}
            <button type="button" class="btn-close btn-close-dark me-2 m-auto" aria-label="Close" onclick={() => toastStore.remove(toast.id)}></button>
        {/if}
    </div>
    <div style="--toast-duration:{toast.duration}ms;" class="toast-progress bg-{getToastColor(toast.type)}"></div>
</div>

<style>
    .toast {
        min-width: 300px;
        max-width: 400px;
        position: relative;
        overflow: hidden;
    }

    /* Entrance animation */
    @keyframes slidefade-in {
        from {
            transform: translateY(30%);
            opacity: 0;
        }
        to {
            transform: translateY(0);
            opacity: 1;
        }
    }

    /* Exit animation */
    @keyframes slidefade-out {
        from {
            transform: translateY(0);
            opacity: 1;
        }
        to {
            transform: translateY(30%);
            opacity: 0;
        }
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
