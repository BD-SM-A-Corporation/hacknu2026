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
    import { activeLocomotiveId } from '@/lib/telemetry';
    import {
        openLocomotiveSelector,
        openWebSocketConfig,
    } from '@/lib/uiStore';
    import { isRealtime, formattedRange } from '@/lib/dateRangeStore';
    import { wsClient } from '@/lib/websocketClient';

    onMount(() => {
        wsClient.connect();
    });
</script>

<AppHead title="Дэшборд - Телеметрия" />

<!-- Locomotive Status Bar -->
<div
    class="flex flex-col sm:flex-row justify-between items-center bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl p-4 shadow-sm mb-6"
>
    <div class="flex items-center gap-4">
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
        <div>
            <span class="text-sm opacity-60">Режим:</span>
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
{#if $activeLocomotiveId}
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-5">
        <!-- Main stats -->
        <HealthWidget />
        <SpeedWidget />
        <TemperatureWidget />
        <PressureWidget />
        <FuelWidget />
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
            Выберите локомотив из списка для мониторинга телеметрии в реальном
            времени
        </p>
        <button
            on:click={openLocomotiveSelector}
            class="px-6 py-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-medium transition-colors shadow-sm"
        >
            Выбрать локомотив
        </button>
    </div>
{/if}

<!-- Dynamic Modals Rendered At Layout Level -->
<LocomotiveSelectorModal />
<WebSocketConfigModal />
<UserSettingsModal />
