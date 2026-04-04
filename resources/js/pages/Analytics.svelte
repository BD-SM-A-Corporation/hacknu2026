<script module lang="ts">
    export const layout = {
        breadcrumbs: [
            {
                title: 'Аналитика',
                href: '/analytics',
            },
        ],
    };
</script>

<script lang="ts">
    import AppHead from '@/components/AppHead.svelte';
    import { activeLocomotiveId } from '@/lib/telemetry';
    import { wsClient } from '@/lib/websocketClient';
    import { onMount } from 'svelte';
    import TelemetryChart from '@/components/telemetry/TelemetryChart.svelte';
    import GridCalendar from '@/components/telemetry/GridCalendar.svelte';
    import HistoricalChartModal from '@/components/telemetry/HistoricalChartModal.svelte';
    import LocomotiveSelectorModal from '@/components/telemetry/LocomotiveSelectorModal.svelte';
    import WebSocketConfigModal from '@/components/telemetry/WebSocketConfigModal.svelte';
    import UserSettingsModal from '@/components/telemetry/UserSettingsModal.svelte';
    import { openLocomotiveSelector } from '@/lib/uiStore';

    let historyModalOpen = $state(false);
    let historyRangeStart = $state<Date | null>(null);
    let historyRangeEnd = $state<Date | null>(null);

    onMount(() => {
        wsClient.connect();
    });

    function handleRangeSelected(start: Date, end: Date) {
        historyRangeStart = start;
        historyRangeEnd = end;
        historyModalOpen = true;
    }
</script>

<AppHead title="Аналитика - Телеметрия" />

<div class="mb-4">
    <h1 class="text-2xl font-bold tracking-tight">Отчеты и Аналитика</h1>
    <p class="text-zinc-500 dark:text-zinc-400">Графики в реальном времени и исторически собранные данные</p>
</div>

{#if $activeLocomotiveId}
    <div class="grid gap-6">
        <!-- Live charts directly migrated from Dashboard -->
        <section>
            <div class="flex items-center justify-between mb-4 mt-6">
                <h2 class="text-lg font-semibold flex items-center gap-2">
                    <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                    Live Телеметрия
                </h2>
                <span class="px-3 py-1 bg-zinc-100 dark:bg-zinc-800 rounded-full text-sm font-mono border border-zinc-200 dark:border-zinc-700">
                    {$activeLocomotiveId}
                </span>
            </div>
            <TelemetryChart />
        </section>

        <!-- Historical selection Timeline migrated from instructions -->
        <section class="mt-8">
            <div class="mb-4">
                <h2 class="text-lg font-semibold">История Данных</h2>
                <p class="text-sm text-zinc-500">Выберите необходимый диапазон времени в календаре ниже, чтобы построить график и проанализировать статистику.</p>
            </div>
            
            <GridCalendar onRangeSelected={handleRangeSelected} />
        </section>
    </div>
    
    <HistoricalChartModal 
        bind:isOpen={historyModalOpen} 
        startDate={historyRangeStart} 
        endDate={historyRangeEnd} 
    />
{:else}
    <div class="flex-1 flex flex-col items-center justify-center p-12 border-2 border-dashed border-zinc-300 dark:border-zinc-700 rounded-xl bg-zinc-50 dark:bg-zinc-800/30 min-h-[500px]">
        <svg
            class="h-16 w-16 text-zinc-400 mb-4"
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
        ><path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"></path><line x1="4" y1="22" x2="4" y2="15"></line></svg>
        <h2 class="text-xl font-bold mb-2">Ожидание выбора</h2>
        <p class="text-zinc-500 mb-6 text-center max-w-sm">
            Выберите локомотив из списка для просмотра аналитики
        </p>
        <button
            onclick={openLocomotiveSelector}
            class="px-6 py-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-medium transition-colors shadow-sm"
        >
            Выбрать локомотив
        </button>
    </div>
{/if}

<!-- Modals -->
<LocomotiveSelectorModal />
<WebSocketConfigModal />
<UserSettingsModal />
