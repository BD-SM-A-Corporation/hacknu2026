<script module lang="ts">
    import { dashboard } from '@/routes';

    export const layout = {
        breadcrumbs: [
            {
                title: 'Дэшборд',
                href: dashboard(),
            },
        ],
    };
</script>

<script lang="ts">
    import { onMount } from 'svelte';
    import AppHead from '@/components/AppHead.svelte';

    // Widgets
    import SpeedWidget from '@/components/telemetry/SpeedWidget.svelte';
    import TemperatureWidget from '@/components/telemetry/TemperatureWidget.svelte';
    import PressureWidget from '@/components/telemetry/PressureWidget.svelte';
    import FuelWidget from '@/components/telemetry/FuelWidget.svelte';
    import HealthWidget from '@/components/telemetry/HealthWidget.svelte';
    import LocomotiveMap from '@/components/telemetry/LocomotiveMap.svelte';

    // UI Modals
    import LocomotiveSelectorModal from '@/components/telemetry/LocomotiveSelectorModal.svelte';
    import UserSettingsModal from '@/components/telemetry/UserSettingsModal.svelte';
    import WebSocketConfigModal from '@/components/telemetry/WebSocketConfigModal.svelte';

    // Stores
    import {
        activeLocomotiveId,
        telemetryData,
        allTelemetry,
        averageFleetTelemetry,
    } from '@/lib/telemetry';
    import CriticalEventModal from '@/components/telemetry/CriticalEventModal.svelte';
    import {
        openLocomotiveSelector,
        openWebSocketConfig,
    } from '@/lib/uiStore';
    import { isRealtime, formattedRange } from '@/lib/dateRangeStore';
    import { wsClient } from '@/lib/websocketClient';

    let now = Date.now();
    let combinedAlerts: string[] = [];

    // Fleet state
    let viewMode = 'single'; // single | fleet
    let anomalies: any[] = [];
    let criticalModalOpen = false;
    let selectedEvent: any = null;

    async function fetchAnomalies() {
        try {
            const res = await fetch('/api/telemetry/anomalies?limit=25');
            if (res.ok) anomalies = await res.json();
        } catch (e) {
            console.error(e);
        }
    }

    $: if (viewMode === 'fleet') {
        fetchAnomalies();
    }

    $: displayAnomalies = (() => {
        let active: any[] = [];
        
        Object.values($allTelemetry).forEach(loco => {
            if (!loco.locomotiveId) return;

            const locoTime = new Date(loco.timestamp).getTime();
            const isDropped = $isRealtime && (now - locoTime > 10000);

            if (isDropped) {
                active.push({
                    locomotive_id: loco.locomotiveId,
                    timestamp: new Date().toISOString(),
                    health_score: loco.healthScore || 0,
                    speed: loco.speed || 0,
                    temperature: loco.temperature || 0,
                    pressure: loco.pressure || 0,
                    fuel_level: loco.fuelLevel || 0,
                    _label: "ОБРЫВ СВЯЗИ"
                });
            } else if (loco.alerts && loco.alerts.length > 0) {
                loco.alerts.forEach(a => {
                    active.push({
                        locomotive_id: loco.locomotiveId,
                        timestamp: new Date(loco.timestamp).toISOString(),
                        health_score: loco.healthScore,
                        speed: loco.speed,
                        temperature: loco.temperature,
                        pressure: loco.pressure,
                        fuel_level: loco.fuelLevel,
                        _label: a
                    });
                });
            }
        });

        const combined = [...active, ...anomalies];
        // optional filter unique? But active changes often, so history + active is fine
        return combined.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
    })();

    function openCriticalModal(event: any) {
        selectedEvent = event;
        criticalModalOpen = true;
    }

    onMount(() => {
        wsClient.connect();

        const interval = setInterval(() => {
            now = Date.now();
        }, 2000);
        return () => clearInterval(interval);
    });

    $: {
        const activeAlerts: string[] = [];

        Object.values($allTelemetry).forEach(loco => {
            if (!loco.locomotiveId) return;

            const locoTime = new Date(loco.timestamp).getTime();
            const isDropped = $isRealtime && (now - locoTime > 10000);

            if (isDropped) {
                activeAlerts.push(`ОБРЫВ СВЯЗИ (${loco.locomotiveId})`);
            }
            if (loco.alerts && loco.alerts.length > 0) {
                loco.alerts.forEach(a => activeAlerts.push(`${a} (${loco.locomotiveId})`));
            }
        });

        combinedAlerts = activeAlerts;
    }
</script>

<AppHead title="Дэшборд - Телеметрия" />

<!-- Global Alerts -->
{#if combinedAlerts.length > 0}
    <div class="fixed top-20 right-4 z-50 flex flex-col gap-2 max-w-sm pointer-events-none">
        {#each combinedAlerts as alert}
            <div class="bg-red-600 text-white px-4 py-3 rounded-xl shadow-lg border border-red-500 font-semibold flex items-center gap-3 animate-in slide-in-from-right-10 pointer-events-auto">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
                {alert}
            </div>
        {/each}
    </div>
{/if}

<!-- Locomotive Status Bar -->
<div
    class="flex flex-col sm:flex-row justify-between items-center bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl p-4 shadow-sm mb-6"
>
    <div class="flex items-center gap-4">
        <div>
            <span class="text-sm opacity-60">Режим обзора:</span>
            <div class="flex items-center gap-2 mt-1 bg-zinc-100 dark:bg-zinc-800 p-1 rounded-lg w-max">
                <button 
                    class="px-3 py-1 rounded-md text-sm font-medium transition-colors {viewMode === 'single' ? 'bg-white dark:bg-zinc-700 shadow' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-300'}"
                    on:click={() => viewMode = 'single'}
                >Один</button>
                <button 
                    class="px-3 py-1 rounded-md text-sm font-medium transition-colors {viewMode === 'fleet' ? 'bg-white dark:bg-zinc-700 shadow' : 'text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-300'}"
                    on:click={() => viewMode = 'fleet'}
                >Весь парк</button>
            </div>
        </div>
        <div class="hidden sm:block w-px h-10 bg-zinc-200 dark:bg-zinc-700 mx-2"></div>
        
        {#if viewMode === 'single'}
            <div>
                <span class="text-sm opacity-60">Текущий локомотив:</span>
                <div class="text-xl font-bold flex items-center gap-2">
                    {#if $activeLocomotiveId}
                        <span class="text-emerald-600 dark:text-emerald-400"
                            >{$activeLocomotiveId}</span
                        >
                    {:else}
                        <span class="opacity-50">Не выбран</span>
                    {/if}
                </div>
            </div>
            <div class="hidden sm:block w-px h-10 bg-zinc-200 dark:bg-zinc-700 mx-2"></div>
        {/if}

        <div>
            <span class="text-sm opacity-60">Соединение:</span>
            <div class="text-md font-medium mt-0.5">
                {#if $isRealtime}
                    <span class="inline-flex items-center gap-1.5 text-emerald-600 dark:text-emerald-400 font-semibold">
                        <span class="relative flex h-2 w-2">
                            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                        </span>
                        В реальном времени
                    </span>
                {:else}
                    <span class="text-zinc-700 dark:text-zinc-300">
                        Исторические данные: {$formattedRange}
                    </span>
                {/if}
            </div>
        </div>
    </div>

    <div class="flex gap-2 mt-4 sm:mt-0">
        <button
            on:click={openLocomotiveSelector}
            class="px-4 py-2 bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 rounded-lg text-sm font-medium transition-colors border border-zinc-200 dark:border-zinc-700"
        >
            Сменить локомотив
        </button>
        <button
            on:click={openWebSocketConfig}
            class="px-4 py-2 bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 rounded-lg text-sm font-medium transition-colors border border-zinc-200 dark:border-zinc-700"
        >
            Настройки WS
        </button>
    </div>
</div>

<!-- Map Component ALWAYS shown -->
<LocomotiveMap />

<!-- Telemetry Dashboard Grid -->
{#if viewMode === 'fleet' || $activeLocomotiveId}
    <div class="grid gap-6 {viewMode === 'fleet' ? 'xl:grid-cols-4' : 'lg:grid-cols-5'} mb-6 mt-6">
        {#if viewMode === 'fleet'}
            <div class="xl:col-span-3 flex flex-col gap-6">
                {#each Object.values($allTelemetry) as loco}
                    {#if loco.locomotiveId}
                        <div class="bg-zinc-50 dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 rounded-xl p-4 shadow-sm flex flex-col gap-4">
                            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                                <h3 class="font-bold text-lg text-zinc-900 dark:text-white flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-emerald-600 dark:text-emerald-500"><rect x="4" y="4" width="16" height="16" rx="2" ry="2"></rect><rect x="9" y="9" width="6" height="6"></rect><line x1="9" y1="1" x2="9" y2="4"></line><line x1="15" y1="1" x2="15" y2="4"></line><line x1="9" y1="20" x2="9" y2="23"></line><line x1="15" y1="20" x2="15" y2="23"></line><line x1="20" y1="9" x2="23" y2="9"></line><line x1="20" y1="14" x2="23" y2="14"></line><line x1="1" y1="9" x2="4" y2="9"></line><line x1="1" y1="14" x2="4" y2="14"></line></svg>
                                    Локомотив {loco.locomotiveId}
                                </h3>
                                <div class="flex items-center gap-2">
                                    <button 
                                        on:click={() => { $activeLocomotiveId = loco.locomotiveId; viewMode = 'single'; }}
                                        class="text-xs font-medium bg-emerald-100/50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 hover:bg-emerald-100 dark:hover:bg-emerald-900/50 px-2 py-1 rounded transition-colors"
                                    >
                                        К анализу →
                                    </button>
                                    <span class="text-xs font-medium bg-white dark:bg-zinc-800 px-2 py-1 rounded border border-zinc-200 dark:border-zinc-700 shadow-sm text-zinc-500">
                                        Обновлено: {new Date(loco.timestamp).toLocaleTimeString('ru-RU')}
                                    </span>
                                </div>
                            </div>
                            <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
                                <HealthWidget overrideHealth={loco.healthScore} />
                                <SpeedWidget overrideSpeed={loco.speed} />
                                <TemperatureWidget overrideTemperature={loco.temperature} />
                                <PressureWidget overridePressure={loco.pressure} />
                                <FuelWidget overrideFuel={loco.fuelLevel} />
                            </div>
                        </div>
                    {/if}
                {/each}
                {#if Object.values($allTelemetry).length === 0}
                    <div class="bg-zinc-50 dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 rounded-xl py-12 flex items-center justify-center text-zinc-500 shadow-sm">
                        Ожидание данных телеметрии...
                    </div>
                {/if}
            </div>

            <div class="flex flex-col gap-4">
                <h2 class="text-lg font-semibold flex items-center justify-between">
                    Критические события
                    <button on:click={fetchAnomalies} class="p-1 hover:bg-zinc-200 dark:hover:bg-zinc-800 rounded transition text-zinc-500">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
                    </button>
                </h2>
                
                <div class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden flex flex-col h-[500px] shadow-sm">
                    {#if displayAnomalies.length === 0}
                        <div class="p-8 flex flex-col items-center justify-center text-zinc-500 h-full">
                            <p>Нет критических событий</p>
                        </div>
                    {:else}
                        <div class="overflow-y-auto flex-1 p-3">
                            {#each displayAnomalies as anomaly}
                                <button 
                                    class="w-full text-left p-3 mb-2 rounded-lg border border-red-200 dark:border-red-900/40 bg-red-50/50 dark:bg-red-900/10 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors shadow-sm"
                                    on:click={() => openCriticalModal(anomaly)}
                                >
                                    <div class="flex justify-between items-start mb-1">
                                        <div class="flex flex-col items-start gap-1">
                                            <span class="font-mono text-sm font-bold text-red-700 dark:text-red-400">{anomaly.locomotive_id}</span>
                                            {#if anomaly._label}
                                                <span class="bg-red-600 text-white text-[10px] uppercase font-bold px-1.5 py-0.5 rounded shadow-sm">{anomaly._label}</span>
                                            {/if}
                                        </div>
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
        {:else}
            <!-- Main stats -->
            <HealthWidget />
            <SpeedWidget />
            <TemperatureWidget />
            <PressureWidget />
            <FuelWidget />
        {/if}
    </div>
{:else}
    <div
        class="flex-1 flex flex-col items-center justify-center p-12 border-2 border-dashed border-zinc-300 dark:border-zinc-700 rounded-xl bg-zinc-50 dark:bg-zinc-800/30"
    >
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
            ><path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"
            ></path><line x1="4" y1="22" x2="4" y2="15"></line></svg
        >
        <h2 class="text-xl font-bold mb-2">Ожидание выбора</h2>
        <p class="text-zinc-500 mb-6 text-center max-w-sm">
            Выберите локомотив из списка для мониторинга в реальном времени, 
            либо переключитесь в режим обзора «Весь парк».
        </p>
        <button
            on:click={openLocomotiveSelector}
            class="px-6 py-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-medium transition-colors shadow-sm"
        >
            Выбрать локомотив
        </button>
    </div>
{/if}

<CriticalEventModal bind:isOpen={criticalModalOpen} event={selectedEvent} />

<!-- Dynamic Modals Rendered At Layout Level -->
<LocomotiveSelectorModal />
<WebSocketConfigModal />
<UserSettingsModal />
