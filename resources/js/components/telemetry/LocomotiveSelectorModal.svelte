<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { index } from '@/actions/App/Http/Controllers/Api/LocomotiveController';
    import { activeLocomotiveId } from '@/lib/telemetry';
    import {
        isLocomotiveSelectorOpen,
        closeLocomotiveSelector
    } from '@/lib/uiStore';
    import ModalContainer from './ModalContainer.svelte';

    let searchQuery = '';
    let allLocomotives: { id: string, model: string, depot: string, status?: string }[] = [];
    let isLoading = true;

    function handleLocomotivePing(e: Event) {
        const customEvent = e as CustomEvent<any>;
        const data = customEvent.detail;
        const id = data.locomotiveId;

        let newStatus = 'active';

        // Check anomaly, critical triggers. Exclude KZ8A purely from fuel check since electric locos report 0 fuel
        if (data.is_anomaly || data.temperature > 100 || data.pressure < 4.5 || (!data.locomotiveId.includes('KZ8A') && data.fuelLevel < 10)) {
            newStatus = 'warning';
        } else if (data.speed < 1) {
            newStatus = 'stopped';
        }

        const index = allLocomotives.findIndex(l => l.id === id);

        if (index === -1) {
            // Automatically push new locomotives to the visual list
            allLocomotives = [...allLocomotives, { id, model: 'Авто-обнаружение', depot: 'Сеть', status: newStatus }];
        } else {
            if (allLocomotives[index].status !== newStatus) {
                allLocomotives[index].status = newStatus;
                allLocomotives = [...allLocomotives];
            }
        }
    }

    onMount(async () => {
        try {
            const response = await fetch(index.url(), {
                credentials: 'same-origin',
                headers: {
                    'Accept': 'application/json',
                    'X-Requested-With': 'XMLHttpRequest'
                }
            });

            if (response.ok) {
                const res = await response.json();
                allLocomotives = res.data || [];
            }
        } catch (e) {
            console.error('Failed to load locomotives:', e);
        } finally {
            isLoading = false;
        }

        if (typeof window !== 'undefined') {
            window.addEventListener('locomotivePing', handleLocomotivePing);
        }
    });

    onDestroy(() => {
        if (typeof window !== 'undefined') {
            window.removeEventListener('locomotivePing', handleLocomotivePing);
        }
    });

    $: filteredLocomotives = allLocomotives.filter(l =>
        l.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
        l.model.toLowerCase().includes(searchQuery.toLowerCase())
    );

    function selectLocomotive(id: string) {
        $activeLocomotiveId = id;
        closeLocomotiveSelector();
    }
</script>

<ModalContainer
    isOpen={$isLocomotiveSelectorOpen}
    title="Выбор Локомотива"
    on:close={closeLocomotiveSelector}
>
    <div class="relative w-full mb-4">
        <svg class="absolute left-3 top-2.5 h-4 w-4 text-zinc-400" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
        <input
            type="text"
            bind:value={searchQuery}
            placeholder="Поиск по ID или типу..."
            class="w-full pl-9 pr-4 py-2 bg-zinc-100 dark:bg-zinc-800 border-none rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none"
        />
    </div>

    <div class="space-y-2 max-h-80 overflow-y-auto pr-2">
        {#if isLoading}
            <div class="text-center py-8 opacity-50">Загрузка локомотивов...</div>
        {:else}
            {#each filteredLocomotives as loco (loco.id)}
            <button
                class={`w-full text-left p-3 rounded-xl border transition-all ${
                    $activeLocomotiveId === loco.id
                    ? 'border-emerald-500 bg-emerald-50 dark:bg-emerald-500/10'
                    : 'border-zinc-200 dark:border-zinc-800 hover:border-emerald-300 dark:hover:border-emerald-500/50'
                }`}
                on:click={() => selectLocomotive(loco.id)}
            >
                <div class="flex justify-between items-center mb-1">
                    <span class="font-bold text-lg">{loco.id}</span>
                    <span class={`text-xs px-2 py-1 rounded-full ${
                        loco.status === 'active' || loco.status === 'Active' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400' :
                        loco.status === 'warning' ? 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400' :
                        loco.status === 'stopped' ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400' :
                        'bg-zinc-100 text-zinc-700 dark:bg-zinc-800 dark:text-zinc-400'
                    }`}>
                        {loco.status === 'active' || loco.status === 'Active' ? 'В пути' :
                         loco.status === 'warning' ? 'Внимание / Аномалия' :
                         loco.status === 'stopped' ? 'Остановка' :
                         loco.status === 'Offline' ? 'Не в сети' : 'В депо'}
                    </span>
                </div>
                <div class="flex justify-between text-sm opacity-70">
                    <span>Модель: {loco.model}</span>
                    <span>{loco.depot || 'Без депо'}</span>
                </div>
            </button>
            {:else}
                <div class="text-center py-8 opacity-50">
                    Локомотивы не найдены
                </div>
            {/each}
        {/if}
    </div>
</ModalContainer>
