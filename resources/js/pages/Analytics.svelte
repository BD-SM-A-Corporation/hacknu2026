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
    import AllLocomotivesChart from '@/components/telemetry/AllLocomotivesChart.svelte';
    import CriticalEventModal from '@/components/telemetry/CriticalEventModal.svelte';
    import { openLocomotiveSelector } from '@/lib/uiStore';

    let historyModalOpen = $state(false);
    let historyRangeStart = $state<Date | null>(null);
    let historyRangeEnd = $state<Date | null>(null);

    let viewMode = $state<'single' | 'all'>('single');
    let anomalies = $state<any[]>([]);
    let loadingAnomalies = $state(false);
    let criticalModalOpen = $state(false);
    let selectedEvent = $state<any>(null);

    async function fetchAnomalies() {
        loadingAnomalies = true;
        try {
            const res = await fetch('/api/telemetry/anomalies?limit=25');
            if (res.ok) anomalies = await res.json();
        } catch (e) {
            console.error(e);
        } finally {
            loadingAnomalies = false;
        }
    }

    $effect(() => {
        if (viewMode === 'all') fetchAnomalies();
    });

    function openCriticalModal(event: any) {
        selectedEvent = event;
        criticalModalOpen = true;
    }

    function downloadAnomaliesCSV() {
        if (!anomalies.length) return;

        let csvContent = 'ID локомотива;Дата и Время;Health Score;Скорость (км/ч);Температура (°C);Давление;Уровень топлива\n';
        anomalies.forEach((item) => {
            const dt = new Date(item.timestamp).toLocaleString('ru-RU');
            const hs = item.health_score || 0;
            const speed = Number(item.speed).toFixed(2);
            const temp = Number(item.temperature).toFixed(2);
            const press = Number(item.pressure).toFixed(2);
            const fuel = Number(item.fuel_level).toFixed(2);
            csvContent += `"${item.locomotive_id}";"${dt}";${hs};${speed};${temp};${press};${fuel}\n`;
        });

        const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = `critical_events_${new Date().toISOString().split('T')[0]}.csv`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    onMount(() => {
        wsClient.connect();
    });

    function handleRangeSelected(start: Date, end: Date) {
        historyRangeStart = start;
        historyRangeEnd = end;
        if (viewMode === 'single') {
            historyModalOpen = true;
        }
    }
</script>

<AppHead title="Аналитика - Телеметрия" />

<div class="mb-4">
    <h1 class="text-2xl font-bold tracking-tight">Отчеты и Аналитика</h1>
    <p class="text-zinc-500 dark:text-zinc-400">Графики в реальном времени и исторически собранные данные</p>
</div>

<div class="mb-6 flex gap-4">
    <div class="inline-flex bg-zinc-100 dark:bg-zinc-800 p-1 rounded-lg">
        <button 
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors {viewMode === 'single' ? 'bg-white dark:bg-zinc-700 shadow flex items-center gap-2' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-300'}"
            onclick={() => viewMode = 'single'}
        >
            <span class="w-2 h-2 rounded-full {viewMode === 'single' && $activeLocomotiveId ? 'bg-emerald-500 animate-pulse' : 'bg-transparent'}"></span>
            Один локомотив
        </button>
        <button 
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors {viewMode === 'all' ? 'bg-white dark:bg-zinc-700 shadow' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-300'}"
            onclick={() => viewMode = 'all'}
        >
            Все локомотивы
        </button>
    </div>
</div>

{#if viewMode === 'all'}
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
        <!-- Main Chart Area -->
        <div class="xl:col-span-2 flex flex-col gap-6">
            <AllLocomotivesChart startDate={historyRangeStart} endDate={historyRangeEnd} />
            
            <section class="mt-2 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 p-6 rounded-xl shadow-sm">
                <div class="mb-4">
                    <h2 class="text-lg font-semibold">История Данных</h2>
                    <p class="text-sm text-zinc-500">Выберите необходимый диапазон времени в календаре ниже для перерендера графика.</p>
                </div>
                <GridCalendar onRangeSelected={handleRangeSelected} />
            </section>
        </div>
        
        <!-- Anomalies Column -->
        <div class="flex flex-col gap-4">
            <h2 class="text-lg font-semibold flex items-center justify-between">
                Критические события
                <div class="flex items-center gap-1">
                    <button onclick={downloadAnomaliesCSV} aria-label="Скачать список" class="p-1 hover:bg-zinc-200 dark:hover:bg-zinc-800 rounded transition text-zinc-500" title="Экспорт в CSV">
                        <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8" /><path d="M8 13h2"/><path d="M8 17h2"/><path d="M14 13h2"/><path d="M14 17h2"/></svg>
                    </button>
                    <button aria-label="Обновить" onclick={fetchAnomalies} class="p-1 hover:bg-zinc-200 dark:hover:bg-zinc-800 rounded transition text-zinc-500" title="Обновить">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
                    </button>
                </div>
            </h2>
            
            <div class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden flex flex-col min-h-[500px] h-full shadow-sm">
                {#if loadingAnomalies}
                    <div class="p-8 flex justify-center mt-20"><div class="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-600"></div></div>
                {:else if anomalies.length === 0}
                    <div class="p-8 flex flex-col items-center justify-center text-zinc-500 h-full mt-10">
                        <svg class="w-12 h-12 mb-3 text-zinc-300 dark:text-zinc-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                        <p>Нет критических событий</p>
                    </div>
                {:else}
                    <div class="overflow-y-auto flex-1 p-3">
                        {#each anomalies as anomaly}
                            <button 
                                class="w-full text-left p-3 mb-2 rounded-lg border border-red-200 dark:border-red-900/40 bg-red-50/50 dark:bg-red-900/10 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors shadow-sm"
                                onclick={() => openCriticalModal(anomaly)}
                            >
                                <div class="flex justify-between items-start mb-1">
                                    <span class="font-mono text-sm font-bold text-red-700 dark:text-red-400">{anomaly.locomotive_id}</span>
                                    <span class="text-xs font-medium text-zinc-500">{new Date(anomaly.timestamp).toLocaleTimeString('ru-RU')}</span>
                                </div>
                                <div class="text-xs text-zinc-600 dark:text-zinc-400 flex flex-wrap items-center gap-2 mt-2">
                                    <span class="bg-white dark:bg-zinc-800 px-1.5 py-0.5 rounded border border-red-100 dark:border-red-900/30 font-mono text-red-600 dark:text-red-400 font-bold">HS: {anomaly.health_score || 0}</span>
                                    <span class="bg-zinc-100 dark:bg-zinc-800 px-1.5 py-0.5 rounded border border-zinc-200 dark:border-zinc-700">{Number(anomaly.speed).toFixed(1)} км/ч</span>
                                </div>
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
        </div>
    </div>
    <CriticalEventModal bind:isOpen={criticalModalOpen} event={selectedEvent} />
{:else if $activeLocomotiveId}
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
