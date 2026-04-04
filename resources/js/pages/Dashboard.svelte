<script module lang="ts">
    import { dashboard } from '@/routes';

    export const layout = {
        breadcrumbs: [
            {
                title: 'Центр управления',
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

    // UI Modals
    import LocomotiveSelectorModal from '@/components/telemetry/LocomotiveSelectorModal.svelte';
    import UserSettingsModal from '@/components/telemetry/UserSettingsModal.svelte';
    import WebSocketConfigModal from '@/components/telemetry/WebSocketConfigModal.svelte';

    // Stores
    import { activeLocomotiveId } from '@/lib/telemetry';
    import {
        openLocomotiveSelector,
        openWebSocketConfig,
        openUserSettings,
    } from '@/lib/uiStore';
    import { wsClient } from '@/lib/websocketClient';

    onMount(() => {
        wsClient.connect();
    });
</script>

<AppHead title="Dashboard - Телеметрия" />

<!-- Control Action Bar -->
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
        <button
            on:click={openUserSettings}
            aria-label="Настройки пользователя"
            class="p-2 bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 rounded-lg transition-colors border border-zinc-200 dark:border-zinc-700"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                ><circle cx="12" cy="12" r="3"></circle><path
                    d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
                ></path></svg
            >
        </button>
    </div>
</div>

<!-- Telemetry Dashboard Grid -->
{#if $activeLocomotiveId}
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <!-- Main stats -->
        <SpeedWidget />
        <TemperatureWidget />
        <PressureWidget />
        <FuelWidget />
    </div>

    <!-- Live Charts placeholder / extra wide widget area -->
    <div class="mt-6">
        <div
            class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl p-6 min-h-64 flex items-center justify-center shadow-sm"
        >
            <div class="text-center opacity-50">
                <svg
                    class="mx-auto h-12 w-12 mb-3"
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    ><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"
                    ></polyline></svg
                >
                <p>Место для графиков динамики (Chart.js / D3)</p>
            </div>
        </div>
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
