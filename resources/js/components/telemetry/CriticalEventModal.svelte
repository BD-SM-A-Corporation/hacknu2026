<script lang="ts">
    import { Chart, registerables } from 'chart.js';

    Chart.register(...registerables);

    let {
        isOpen = $bindable(false),
        event = null,
    } = $props<{
        isOpen: boolean;
        event: any | null;
    }>();

    let loading = $state(false);
    let chartCanvas = $state<HTMLCanvasElement>();
    let chart: Chart | null = null;

    let showSpeed = $state(true);
    let showTemperature = $state(true);
    let showPressure = $state(true);
    let showFuel = $state(true);
    let hasData = $state(false);
    let currentData = $state<any[]>([]);
    
    let stepSelection = $state('10s');

    $effect(() => {
        if (isOpen && event) {
            fetchData();
        }
    });

    function close() {
        isOpen = false;
        if (chart) {
            chart.destroy();
            chart = null;
        }
    }

    async function fetchData() {
        if (!event) return;
        loading = true;
        hasData = false;
        
        try {
            // Window: 3 hours before, 3 hours after
            const eventTime = new Date(event.timestamp);
            const start = new Date(eventTime.getTime() - 3 * 3600 * 1000);
            const end = new Date(eventTime.getTime() + 3 * 3600 * 1000);

            let stepParam = '';
            if (stepSelection !== '1s') {
                const num = parseInt(stepSelection);
                const unit = stepSelection.slice(-1);
                let multiplier = 1;
                if (unit === 'm') multiplier = 60;
                if (unit === 'h') multiplier = 3600;
                stepParam = `&step=${num * multiplier}`;
            }

            const res = await fetch(
                `/api/telemetry/analytics?locomotive_id=${event.locomotive_id}&start_date=${start.toISOString()}&end_date=${end.toISOString()}${stepParam}`,
            );
            if (res.ok) {
                const data = await res.json();
                hasData = data.length > 0;
                currentData = data;
                setTimeout(() => {
                    renderChart(data);
                }, 50);
            }
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    function downloadChart() {
        if (!chartCanvas || !chart) return;
        const link = document.createElement('a');
        link.download = `anomaly_${event.locomotive_id}_${new Date(event.timestamp).getTime()}.png`;
        link.href = chart.toBase64Image();
        link.click();
    }

    function downloadCSV() {
        if (!currentData.length) return;

        let csvContent =
            'Дата и Время;Скорость (км/ч);Температура (°C);Давление;Уровень топлива\n';
        currentData.forEach((item) => {
            const dt = new Date(item.timestamp).toLocaleString('ru-RU');
            const speed = Number(item.speed).toFixed(2);
            const temp = Number(item.temperature).toFixed(2);
            const press = Number(item.pressure).toFixed(2);
            const fuel = Number(item.fuel_level).toFixed(2);
            csvContent += `"${dt}";${speed};${temp};${press};${fuel}\n`;
        });

        const blob = new Blob(['\uFEFF' + csvContent], {
            type: 'text/csv;charset=utf-8;',
        });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = `anomaly_${event.locomotive_id}_${new Date(event.timestamp).getTime()}.csv`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    function renderChart(data: any[]) {
        if (!chartCanvas) return;
        if (chart) chart.destroy();
        if (!hasData) return;

        const labels = data.map((d) => {
            return new Date(d.timestamp).toLocaleTimeString([], { hour12: false });
        });

        const speed = data.map((d) => Number(d.speed));
        const temp = data.map((d) => Number(d.temperature));
        const press = data.map((d) => Number(d.pressure));
        const fuel = data.map((d) => Number(d.fuel_level));

        chart = new Chart(chartCanvas, {
            type: 'line',
            data: {
                labels,
                datasets: [
                    {
                        label: 'Скорость',
                        data: speed,
                        borderColor: '#059669',
                        hidden: !showSpeed,
                        tension: 0.2,
                        borderWidth: 1.5,
                        yAxisID: 'y',
                    },
                    {
                        label: 'Температура',
                        data: temp,
                        borderColor: '#dc2626',
                        hidden: !showTemperature,
                        tension: 0.2,
                        borderWidth: 1.5,
                        yAxisID: 'yTemp',
                    },
                    {
                        label: 'Давление',
                        data: press,
                        borderColor: '#2563eb',
                        hidden: !showPressure,
                        tension: 0.2,
                        borderWidth: 1.5,
                        yAxisID: 'yPress',
                    },
                    {
                        label: 'Топливо',
                        data: fuel,
                        borderColor: '#d97706',
                        hidden: !showFuel,
                        tension: 0.2,
                        borderWidth: 1.5,
                        yAxisID: 'y',
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { 
                    legend: { display: false }
                },
                elements: { point: { radius: 0 } },
                interaction: { mode: 'index', intersect: false },
                scales: {
                    x: { ticks: { maxTicksLimit: 15 } },
                    y: {
                        type: 'linear',
                        display: true,
                        position: 'left',
                        title: { display: true, text: 'Скорость / Топливо' },
                    },
                    yTemp: {
                        type: 'linear',
                        display: showTemperature,
                        position: 'right',
                        title: { display: true, text: 'Температура (°C)' },
                        grid: { drawOnChartArea: false },
                    },
                    yPress: {
                        type: 'linear',
                        display: showPressure,
                        position: 'right',
                        title: { display: true, text: 'Давление' },
                        grid: { drawOnChartArea: false },
                    },
                },
            },
            plugins: [
                {
                    id: 'custom_canvas_background_color_anomaly',
                    beforeDraw: (chartInstance) => {
                        const ctx = chartInstance.canvas.getContext('2d');
                        if (!ctx) return;
                        ctx.save();
                        ctx.globalCompositeOperation = 'destination-over';
                        const isDark = document.documentElement.classList.contains('dark');
                        ctx.fillStyle = isDark ? '#18181b' : '#ffffff';
                        ctx.fillRect(0, 0, chartInstance.width, chartInstance.height);
                        ctx.restore();
                    },
                },
            ],
        });
    }

    $effect(() => {
        const speed = showSpeed;
        const temp = showTemperature;
        const press = showPressure;
        const fuel = showFuel;

        if (chart) {
            chart.data.datasets[0].hidden = !speed;
            chart.data.datasets[1].hidden = !temp;
            chart.data.datasets[2].hidden = !press;
            chart.data.datasets[3].hidden = !fuel;

            if (chart.options?.scales?.yTemp) chart.options.scales.yTemp.display = temp;
            if (chart.options?.scales?.yPress) chart.options.scales.yPress.display = press;

            chart.update();
        }
    });
</script>

<svelte:window onkeydown={(e) => e.key === 'Escape' && isOpen && close()} />

{#if isOpen && event}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
        onclick={close}
    >
        <div
            class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl shadow-2xl w-full max-w-5xl flex flex-col h-[85vh]"
            onclick={(e) => e.stopPropagation()}
        >
            <div class="p-4 border-b border-zinc-200 dark:border-zinc-800 flex justify-between items-center bg-zinc-50 dark:bg-zinc-800/30 rounded-t-xl overflow-x-auto gap-4">
                <div class="flex-shrink-0">
                    <h2 class="text-lg font-bold flex items-center gap-2">
                        <span class="text-red-500">Детали аномалии</span>
                        <span class="text-emerald-600 font-mono ml-2 block sm:inline">{event.locomotive_id}</span>
                    </h2>
                    <p class="text-sm text-zinc-500 mt-1">Время события: {new Date(event.timestamp).toLocaleString('ru-RU')}</p>
                </div>
                
                <div class="flex items-center gap-4 text-sm text-zinc-500 flex-shrink-0">
                    <div class="flex items-center gap-2 relative z-20">
                        <span class="text-xs">Масштаб (Шаг):</span>
                        <select
                            bind:value={stepSelection}
                            onchange={() => fetchData()}
                            class="p-1 rounded bg-zinc-100 border-zinc-200 dark:bg-zinc-800 dark:border-zinc-700 focus:ring-emerald-500"
                        >
                            <option value="1s">1 сек</option>
                            <option value="5s">5 сек</option>
                            <option value="10s">10 сек</option>
                            <option value="30s">30 сек</option>
                            <option value="1m">1 мин</option>
                            <option value="5m">5 мин</option>
                        </select>
                    </div>

                    <div class="flex items-center gap-2">
                        <button
                            onclick={downloadCSV}
                            class="px-3 py-1.5 flex items-center gap-1 bg-zinc-100 hover:bg-zinc-200 text-zinc-700 dark:bg-zinc-800 dark:hover:bg-zinc-700 dark:text-zinc-300 text-sm rounded shadow-sm transition"
                            title="Скачать данные как CSV (Excel)"
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8" /><path d="M8 13h2"/><path d="M8 17h2"/><path d="M14 13h2"/><path d="M14 17h2"/></svg>
                            CSV
                        </button>
                        <button
                            onclick={downloadChart}
                            class="px-3 py-1.5 flex items-center gap-1 bg-emerald-600/10 text-emerald-600 hover:bg-emerald-600 hover:text-white dark:bg-emerald-500/20 dark:text-emerald-400 dark:hover:bg-emerald-600 dark:hover:text-white text-sm rounded shadow-sm transition"
                            title="Скачать график как PNG"
                        >
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
                            Скачать
                        </button>
                        <button onclick={close} class="p-2 ml-2 hover:bg-zinc-200 dark:hover:bg-zinc-700 rounded transition">&times;</button>
                    </div>
                </div>
            </div>

            <!-- Detail Cards -->
            <div class="grid grid-cols-4 gap-4 p-4 border-b border-zinc-100 dark:border-zinc-800/50">
                <div class="p-3 bg-red-50 dark:bg-red-900/10 rounded-lg border border-red-100 dark:border-red-900/30">
                    <p class="text-xs text-red-600 dark:text-red-400 font-medium">Health Score</p>
                    <p class="text-2xl font-bold text-red-700 dark:text-red-300">{event.health_score || 0}</p>
                </div>
                <div class="p-3 bg-zinc-50 dark:bg-zinc-800/50 rounded-lg border border-zinc-200 dark:border-zinc-700">
                    <p class="text-xs text-zinc-500 font-medium">Скорость</p>
                    <p class="text-xl font-bold">{Number(event.speed).toFixed(1)} км/ч</p>
                </div>
                <div class="p-3 bg-zinc-50 dark:bg-zinc-800/50 rounded-lg border border-zinc-200 dark:border-zinc-700">
                    <p class="text-xs text-zinc-500 font-medium">Температура</p>
                    <p class="text-xl font-bold">{Number(event.temperature).toFixed(1)} &deg;C</p>
                </div>
                <div class="p-3 bg-zinc-50 dark:bg-zinc-800/50 rounded-lg border border-zinc-200 dark:border-zinc-700">
                    <p class="text-xs text-zinc-500 font-medium">Давление</p>
                    <p class="text-xl font-bold">{Number(event.pressure).toFixed(1)} атм</p>
                </div>
            </div>

            <div class="flex-1 p-6 flex flex-col relative bg-zinc-50/50 dark:bg-zinc-900/50 overflow-hidden">
                {#if loading}
                    <div class="flex-1 flex justify-center items-center">
                        <div class="flex flex-col items-center gap-4 text-emerald-600">
                            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600"></div>
                            <p class="font-medium">Загрузка данных +- 3 часа...</p>
                        </div>
                    </div>
                {:else if !hasData}
                    <div class="flex-1 flex flex-col justify-center items-center text-zinc-500 gap-4">
                        <p>Нет подробных данных для этого периода.</p>
                    </div>
                {/if}

                {#if !loading && hasData}
                    <div class="absolute top-6 right-6 z-10 bg-white/90 dark:bg-zinc-800/90 backdrop-blur-md p-3 border border-zinc-200 dark:border-zinc-700 rounded-lg shadow-sm flex flex-col gap-1.5 text-sm min-w-[140px]">
                        <span class="text-xs uppercase text-zinc-500 font-medium tracking-wider mb-1 px-1">Метрики</span>
                        <label class="flex items-center gap-2.5 cursor-pointer transition-colors px-2 py-1.5 rounded-md {showSpeed ? 'bg-emerald-500/15 text-emerald-700 dark:bg-emerald-500/25 dark:text-emerald-300' : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-700/50'}">
                            <input type="checkbox" bind:checked={showSpeed} class="rounded border-zinc-300 text-emerald-600" />
                            <span class="font-medium">Скорость</span>
                        </label>
                        <label class="flex items-center gap-2.5 cursor-pointer transition-colors px-2 py-1.5 rounded-md {showTemperature ? 'bg-red-500/15 text-red-700 dark:bg-red-500/25 dark:text-red-300' : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-700/50'}">
                            <input type="checkbox" bind:checked={showTemperature} class="rounded border-zinc-300 text-red-600" />
                            <span class="font-medium">Температура</span>
                        </label>
                        <label class="flex items-center gap-2.5 cursor-pointer transition-colors px-2 py-1.5 rounded-md {showPressure ? 'bg-blue-500/15 text-blue-700 dark:bg-blue-500/25 dark:text-blue-300' : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-700/50'}">
                            <input type="checkbox" bind:checked={showPressure} class="rounded border-zinc-300 text-blue-600" />
                            <span class="font-medium">Давление</span>
                        </label>
                        <label class="flex items-center gap-2.5 cursor-pointer transition-colors px-2 py-1.5 rounded-md {showFuel ? 'bg-amber-500/15 text-amber-700 dark:bg-amber-500/25 dark:text-amber-300' : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-700/50'}">
                            <input type="checkbox" bind:checked={showFuel} class="rounded border-zinc-300 text-amber-600" />
                            <span class="font-medium">Топливо</span>
                        </label>
                    </div>
                    <div class="flex-1 w-full relative">
                        <canvas bind:this={chartCanvas}></canvas>
                    </div>
                {/if}
            </div>
        </div>
    </div>
{/if}
