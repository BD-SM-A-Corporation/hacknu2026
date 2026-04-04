<script lang="ts">
    import { onMount } from 'svelte';
    import { index, store as storeEndpoint } from '@/actions/App/Http/Controllers/Api/EndpointController';

    import { isWebSocketConfigOpen, closeWebSocketConfig } from '@/lib/uiStore';
    import { wsClient } from '@/lib/websocketClient';
    import ModalContainer from './ModalContainer.svelte';

    let endpointUrl = 'ws://localhost:8080/telemetry';
    let savedEndpoints: { id: number, url: string, is_active: boolean }[] = [];
    let isConnected = false;
    let importError = '';

    onMount(async () => {
        try {
            const res = await fetch(index.url());

            if (res.ok) {
                const data = await res.json();
                savedEndpoints = data.data || [];
            }
        } catch (e) {
            console.error('Failed to load endpoints');
        }
    });

    wsClient.onConnect = () => isConnected = true;
    wsClient.onDisconnect = () => isConnected = false;

    async function applyEndpoint(urlToApply = endpointUrl) {
        endpointUrl = urlToApply;
        wsClient.setUrl(endpointUrl);

        // Also save to backend
        try {
            await fetch(storeEndpoint.url(), {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Accept': 'application/json' },
                body: JSON.stringify({ url: endpointUrl, is_active: true })
            });
            // Update list silently
            const res = await fetch(index.url());

            if (res.ok) {
                const data = await res.json();
                savedEndpoints = data.data || [];
            }
        } catch(e) {
            console.error('Failed to save endpoint', e);
        }
    }

    function handleFileUpload(event: Event) {
        const input = event.target as HTMLInputElement;

        if (!input.files || input.files.length === 0) return;

        const file = input.files[0];
        const reader = new FileReader();

        reader.onload = (e) => {
            try {
                const text = e.target?.result as string;

                if (file.name.endsWith('.json')) {
                    const data = JSON.parse(text);

                    if (data.endpoint) endpointUrl = data.endpoint;
                    else importError = 'В JSON не найден ключ "endpoint"';
                } else if (file.name.endsWith('.csv')) {
                    // Primitive CSV logic: assuming first line has url
                    const firstLine = text.split('\n')[0];

                    if (firstLine.startsWith('ws://') || firstLine.startsWith('wss://')) {
                        endpointUrl = firstLine.trim();
                    } else {
                        importError = 'Неверный формат CSV эндпоинта';
                    }
                }

                if (!importError) {
                    applyEndpoint(endpointUrl);
                }
            } catch (err) {
                importError = 'Ошибка парсинга файла';
            }
        };
        reader.readAsText(file);
    }
</script>

<ModalContainer
    isOpen={$isWebSocketConfigOpen}
    title="Конфигурация WebSocket"
    on:close={closeWebSocketConfig}
>
    <div class="space-y-6">
        <div>
            <label class="block text-sm font-medium mb-2 opacity-80">Текущий Endpoint</label>
            <div class="flex flex-col gap-2">
                <div class="flex gap-2">
                    <input
                        type="text"
                        bind:value={endpointUrl}
                        class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800 border-none rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none font-mono text-sm"
                    />
                    <button
                        on:click={() => applyEndpoint()}
                        class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg transition-colors"
                    >
                        Подключить
                    </button>
                </div>

                {#if savedEndpoints.length > 0}
                <div class="mt-2 text-sm opacity-80 mb-1">Сохраненные эндпоинты:</div>
                <div class="flex flex-wrap gap-2">
                    {#each savedEndpoints as endp}
                        <button
                            on:click={() => applyEndpoint(endp.url)}
                            class="px-3 py-1 bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 rounded-full text-xs font-mono transition-colors"
                        >
                            {endp.url}
                        </button>
                    {/each}
                </div>
                {/if}
            </div>

            <div class="mt-3 flex items-center gap-2 text-sm">
                <span class={`w-2.5 h-2.5 rounded-full ${isConnected ? 'bg-emerald-500' : 'bg-red-500'}`}></span>
                <span class="opacity-80">{isConnected ? 'Подключено' : 'Отключено'}</span>
            </div>
        </div>

        <hr class="border-zinc-200 dark:border-zinc-800" />

        <div>
            <h3 class="text-sm font-medium mb-2 opacity-80">Импорт конфигурации (JSON/CSV)</h3>
            <div class="border-2 border-dashed border-zinc-300 dark:border-zinc-700 rounded-xl p-8 text-center hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors relative cursor-pointer">
                <input
                    type="file"
                    accept=".json,.csv"
                    on:change={handleFileUpload}
                    class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                />
                <svg class="mx-auto h-8 w-8 text-zinc-400 mb-2" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="17 8 12 3 7 8"></polyline><line x1="12" y1="3" x2="12" y2="15"></line></svg>
                <div class="text-sm opacity-80">
                    Нажмите или перетащите файл с настройками
                </div>
            </div>
            {#if importError}
                <div class="mt-2 text-red-500 text-sm flex items-center gap-1">
                    <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
                    {importError}
                </div>
            {/if}
        </div>
    </div>
</ModalContainer>
