<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    export let title: string;
    export let icon: string = '';
    // Status can be normal, warning (amber), critical (red)
    export let status: 'normal' | 'warning' | 'critical' = 'normal';

    const dispatch = createEventDispatcher();
</script>

<div
    class={`flex flex-col rounded-xl border p-4 shadow-sm transition-colors duration-300
    ${
        status === 'critical'
            ? 'bg-red-500/10 border-red-500/50 text-red-700 dark:text-red-400 dark:bg-red-900/20 dark:border-red-500/30'
            : status === 'warning'
              ? 'bg-amber-500/10 border-amber-500/50 text-amber-700 dark:text-amber-400 dark:bg-amber-900/20 dark:border-amber-500/30'
              : 'bg-white border-zinc-200 text-zinc-900 dark:bg-zinc-800/80 dark:border-zinc-700/50 dark:text-zinc-100'
    }
`}
>
    <div class="flex items-center justify-between mb-2 opacity-80">
        <h3 class="text-sm font-medium uppercase tracking-wider">{title}</h3>
        {#if icon}
            <span class="text-xl">{icon}</span>
        {/if}
    </div>

    <div class="flex-grow flex items-center justify-center">
        <slot />
    </div>
</div>
