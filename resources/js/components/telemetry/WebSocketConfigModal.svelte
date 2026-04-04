<script lang="ts">
    import ModalContainer from './ModalContainer.svelte';
    import { isWebSocketConfigOpen, closeWebSocketConfig } from '@/lib/uiStore';
    import { wsClient } from '@/lib/websocketClient';

    let endpointUrl = 'ws://localhost:8080/telemetry';
    let isConnected = false;
    let importError = '';

    wsClient.onConnect = () => isConnected = true;
    wsClient.onDisconnect = () => isConnected = false;

    function applyEndpoint() {
        wsClient.setUrl(endpointUrl);
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
                    applyEndpoint();
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
            <div class="flex gap-2">
                <input
                    type="text"
                    bind:value={endpointUrl}
                    class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800 border-none rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none font-mono text-sm"
                />
                <button
                    on:click={applyEndpoint}
                    class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg transition-colors"
                >
                    Подключить
                </button>
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
