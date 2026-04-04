<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { fade, scale } from 'svelte/transition';

    export let isOpen = false;
    export let title = '';

    const dispatch = createEventDispatcher();

    function close() {
        dispatch('close');
    }
</script>

<svelte:window on:keydown={(e) => { if (e.key === 'Escape' && isOpen) close(); }} />

{#if isOpen}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-zinc-900/50 backdrop-blur-sm"
        transition:fade={{ duration: 200 }}
        on:click={close}
    >
        <div
            class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl shadow-2xl w-full max-w-lg max-h-[90vh] flex flex-col overflow-hidden"
            transition:scale={{ start: 0.95, duration: 200 }}
            on:click|stopPropagation
        >
            <div class="px-6 py-4 border-b border-zinc-200 dark:border-zinc-800 flex justify-between items-center">
                <h2 class="text-lg font-semibold">{title}</h2>
                <button class="p-1 rounded-md hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors" on:click={close}>
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                </button>
            </div>

            <div class="p-6 overflow-y-auto flex-grow">
                <slot />
            </div>
        </div>
    </div>
{/if}
