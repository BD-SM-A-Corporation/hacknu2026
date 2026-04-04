<script lang="ts">
    import Calendar from 'lucide-svelte/icons/calendar';
    import ChevronDown from 'lucide-svelte/icons/chevron-down';
    import {
        dateRange,
        activePreset,
        presets,
        setPreset,
        setCustomRange,
        formattedRange,
        formatDateInput,
        isRealtime,
        activateRealtime,
        type PresetKey,
    } from '@/lib/dateRangeStore';

    let showDropdown = $state(false);
    let customStart = $state(formatDateInput($dateRange.start));
    let customEnd = $state(formatDateInput($dateRange.end));

    function handlePresetClick(key: PresetKey) {
        setPreset(key);
        customStart = formatDateInput($dateRange.start);
        customEnd = formatDateInput($dateRange.end);
        showDropdown = false;
    }

    function handleCustomApply() {
        if (customStart && customEnd) {
            setCustomRange(new Date(customStart), new Date(customEnd));
            showDropdown = false;
        }
    }

    function handleToggle() {
        showDropdown = !showDropdown;
        if (showDropdown) {
            customStart = formatDateInput($dateRange.start);
            customEnd = formatDateInput($dateRange.end);
        }
    }

    function handleClickOutside(e: MouseEvent) {
        const target = e.target as HTMLElement;
        if (!target.closest('.date-range-picker')) {
            showDropdown = false;
        }
    }
</script>

<svelte:window on:click={handleClickOutside} />

<div class="date-range-picker relative">
    <!-- Trigger Button -->
    <button
        onclick={handleToggle}
        class="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-1.5 text-sm font-medium shadow-sm transition-all hover:border-zinc-300 dark:border-zinc-700 dark:hover:border-zinc-600 {$isRealtime ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400' : 'bg-white text-zinc-700 hover:bg-zinc-50 dark:bg-zinc-800 dark:text-zinc-200 dark:hover:bg-zinc-750'}"
    >
        {#if $isRealtime}
            <div class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </div>
            <span class="hidden sm:inline">В реальном времени</span>
        {:else}
            <Calendar class="h-4 w-4 opacity-60" />
            <span class="hidden sm:inline">{$formattedRange}</span>
        {/if}
        <ChevronDown
            class="h-3.5 w-3.5 opacity-50 transition-transform {showDropdown
                ? 'rotate-180'
                : ''}"
        />
    </button>

    <!-- Dropdown -->
    {#if showDropdown}
        <div
            class="absolute left-0 top-full z-50 mt-2 w-72 rounded-xl border border-zinc-200 bg-white p-4 shadow-xl dark:border-zinc-700 dark:bg-zinc-900"
        >
            <div class="mb-4">
                <button
                    onclick={() => { activateRealtime(); showDropdown = false; }}
                    class="flex w-full items-center justify-center gap-2 rounded-lg py-2 text-sm font-semibold transition-all {$isRealtime ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/50 dark:text-emerald-300' : 'bg-zinc-100 text-zinc-600 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-400 dark:hover:bg-zinc-700'}"
                >
                    <div class="relative flex h-2 w-2">
                        {#if $isRealtime}
                            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                        {/if}
                        <span class="relative inline-flex rounded-full h-2 w-2 {$isRealtime ? 'bg-emerald-500' : 'bg-zinc-400'}"></span>
                    </div>
                    Реальное время
                </button>
            </div>

            <!-- Presets -->
            <div class="mb-2 text-xs font-semibold uppercase tracking-wider text-zinc-400">
                Предустановки
            </div>
            <div class="mb-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
                {#each presets as preset (preset.key)}
                    <button
                        onclick={() => handlePresetClick(preset.key)}
                        class="rounded-md px-2 py-1.5 text-xs font-medium transition-all {$activePreset ===
                        preset.key && !$isRealtime
                            ? 'bg-emerald-600 text-white'
                            : 'bg-zinc-100 text-zinc-600 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-400 dark:hover:bg-zinc-700'}"
                    >
                        {preset.label}
                    </button>
                {/each}
            </div>

            <div class="mb-2 text-xs font-semibold uppercase tracking-wider text-zinc-400">
                Произвольный период
            </div>

            <div class="flex flex-col gap-2">
                <div class="flex items-center gap-2">
                    <label for="date-range-start" class="text-xs text-zinc-500 dark:text-zinc-400 w-6">С</label>
                    <input
                        id="date-range-start"
                        type="date"
                        bind:value={customStart}
                        class="flex-1 rounded-lg border border-zinc-200 bg-zinc-50 px-3 py-1.5 text-sm dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-200"
                    />
                </div>
                <div class="flex items-center gap-2">
                    <label for="date-range-end" class="text-xs text-zinc-500 dark:text-zinc-400 w-6">По</label>
                    <input
                        id="date-range-end"
                        type="date"
                        bind:value={customEnd}
                        class="flex-1 rounded-lg border border-zinc-200 bg-zinc-50 px-3 py-1.5 text-sm dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-200"
                    />
                </div>
                <button
                    onclick={handleCustomApply}
                    class="mt-1 w-full rounded-lg bg-emerald-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-emerald-700"
                >
                    Применить
                </button>
            </div>
        </div>
    {/if}
</div>
