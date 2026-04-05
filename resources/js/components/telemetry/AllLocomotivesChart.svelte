<script lang="ts">
    import { Chart, registerables } from 'chart.js';

    Chart.register(...registerables);

    let { startDate = null, endDate = null } = $props<{
        startDate?: Date | null;
        endDate?: Date | null;
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

    let stepSelection = $state('auto');

    $effect(() => {
        let effStart = startDate;
        let effEnd = endDate;
        if (!effStart || !effEnd) {
            // Default to last 24 hours if no date range is provided
            effEnd = new Date();
            effStart = new Date(effEnd.getTime() - 24 * 3600 * 1000);
        }
        
        if (effStart && effEnd) {
            fetchData(effStart, effEnd);
        }
    });

    async function fetchData(start: Date, end: Date) {
        loading = true;
        hasData = false;
        try {
            let stepParam = '';
            let periodSecs = (end.getTime() - start.getTime()) / 1000;

            if (stepSelection === 'auto') {
                if (periodSecs > 60) {
                    const calculatedStep = Math.round(periodSecs / 60);
                    stepParam = `&step=${calculatedStep}`;
                }
            } else if (stepSelection !== '1s') {
                const num = parseInt(stepSelection);
                const unit = stepSelection.slice(-1);
                let multiplier = 1;
                if (unit === 'm') multiplier = 60;
                if (unit === 'h') multiplier = 3600;
                if (unit === 'd') multiplier = 86400;
                stepParam = `&step=${num * multiplier}`;
            }

            const res = await fetch(
                `/api/telemetry/analytics/all?start_date=${start.toISOString()}&end_date=${end.toISOString()}${stepParam}`,
            );
            if (res.ok) {
                const data = await res.json();
                hasData = data.length > 0;
                currentData = data;
                setTimeout(() => {
                    renderChart(data, start, end);
                }, 50);
            }
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    function renderChart(data: any[], start: Date, end: Date) {
        if (!chartCanvas) return;
        if (chart) chart.destroy();
        if (!hasData) return;

        let periodSecs = (end.getTime() - start.getTime()) / 1000;
        const showDate = periodSecs >= 86400;

        const labels = data.map((d) => {
            const dt = new Date(d.timestamp);
            if (showDate) {
                return dt.toLocaleString('ru-RU', {
                    day: '2-digit',
                    month: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit',
                });
            }
            return dt.toLocaleTimeString([], { hour12: false });
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
                        label: 'Ср. Скорость',
                        data: speed,
                        borderColor: '#059669',
                        hidden: !showSpeed,
                        tension: 0.2,
                        borderWidth: 1.5,
                        yAxisID: 'y',
                    },
                    {
                        label: 'Ср. Температура',
                        data: temp,
                        borderColor: '#dc2626',
                        hidden: !showTemperature,
                        tension: 0.2,
                        borderWidth: 1.5,
                        yAxisID: 'yTemp',
                    },
                    {
                        label: 'Ср. Давление',
                        data: press,
                        borderColor: '#2563eb',
                        hidden: !showPressure,
                        tension: 0.2,
                        borderWidth: 1.5,
                        yAxisID: 'yPress',
                    },
                    {
                        label: 'Ср. Топливо',
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
                plugins: { legend: { display: false } },
                elements: { point: { radius: 0 } },
                interaction: { mode: 'index', intersect: false },
                scales: {
                    x: {
                        ticks: { maxTicksLimit: 15 },
                    },
                    y: {
                        type: 'linear',
                        display: true,
                        position: 'left',
                        title: {
                            display: true,
                            text: 'Ср. Скорость / Топливо',
                        },
                    },
                    yTemp: {
                        type: 'linear',
                        display: showTemperature,
                        position: 'right',
                        title: { display: true, text: 'Ср. Температура (°C)' },
                        grid: { drawOnChartArea: false },
                    },
                    yPress: {
                        type: 'linear',
                        display: showPressure,
                        position: 'right',
                        title: { display: true, text: 'Ср. Давление' },
                        grid: { drawOnChartArea: false },
                    },
                },
            },
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

            if (chart.options?.scales?.yTemp) {
                chart.options.scales.yTemp.display = temp;
            }
            if (chart.options?.scales?.yPress) {
                chart.options.scales.yPress.display = press;
            }

            chart.update();
        }
    });

</script>

<div class="relative bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl p-6 shadow-sm min-h-[400px] h-full w-full flex flex-col">
    <div class="flex justify-between items-center mb-4">
        <div>
            <h3 class="font-bold text-lg">Агрегированная аналитика (все локомотивы)</h3>
            <p class="text-sm text-zinc-500">Средние значения по всем локомотивам</p>
        </div>
        <div class="flex items-center gap-2 relative z-20">
            <span class="text-xs text-zinc-500">Шаг:</span>
            <select
                bind:value={stepSelection}
                onchange={() => {
                    let effStart = startDate;
                    let effEnd = endDate;
                    if (!effStart || !effEnd) {
                        effEnd = new Date();
                        effStart = new Date(effEnd.getTime() - 24 * 3600 * 1000);
                    }
                    if (effStart && effEnd) {
                        fetchData(effStart, effEnd);
                    }
                }}
                class="p-1 text-sm rounded bg-zinc-100 border-zinc-200 dark:bg-zinc-800 dark:border-zinc-700 focus:ring-emerald-500"
            >
                <option value="auto">Авто</option>
                <option value="1m">1 мин</option>
                <option value="5m">5 мин</option>
                <option value="10m">10 мин</option>
                <option value="30m">30 мин</option>
                <option value="1h">1 час</option>
            </select>
        </div>
    </div>
    
    <div class="flex-1 relative w-full pr-36">
        {#if loading}
            <div class="absolute inset-0 flex justify-center items-center z-10 bg-white/50 dark:bg-zinc-900/50">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-600"></div>
            </div>
        {/if}
        
        {#if !hasData && !loading}
            <div class="absolute inset-0 flex flex-col justify-center items-center text-zinc-500 z-10">
                <p>Нет данных для отображения</p>
            </div>
        {/if}

        <canvas bind:this={chartCanvas}></canvas>
        
        {#if !loading && hasData}
            <div
                class="absolute top-4 right-0 z-10 bg-white/90 dark:bg-zinc-800/90 backdrop-blur-md p-3 border border-zinc-200 dark:border-zinc-700 rounded-lg shadow-sm flex flex-col gap-1.5 text-sm min-w-[130px]"
            >
                <span class="text-xs uppercase text-zinc-500 font-medium tracking-wider mb-1 px-1">Метрики</span>

                <label class="flex items-center gap-2.5 cursor-pointer transition-colors px-2 py-1.5 rounded-md {showSpeed ? 'bg-emerald-500/15 text-emerald-700 dark:bg-emerald-500/25 dark:text-emerald-300' : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-700/50'}">
                    <input type="checkbox" bind:checked={showSpeed} class="rounded border-zinc-300 text-emerald-600 focus:ring-emerald-500 dark:bg-zinc-800 dark:border-zinc-600" />
                    <span class="font-medium">Скорость</span>
                </label>
                <label class="flex items-center gap-2.5 cursor-pointer transition-colors px-2 py-1.5 rounded-md {showTemperature ? 'bg-red-500/15 text-red-700 dark:bg-red-500/25 dark:text-red-300' : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-700/50'}">
                    <input type="checkbox" bind:checked={showTemperature} class="rounded border-zinc-300 text-red-600 focus:ring-red-500 dark:bg-zinc-800 dark:border-zinc-600" />
                    <span class="font-medium">Температура</span>
                </label>
                <label class="flex items-center gap-2.5 cursor-pointer transition-colors px-2 py-1.5 rounded-md {showPressure ? 'bg-blue-500/15 text-blue-700 dark:bg-blue-500/25 dark:text-blue-300' : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-700/50'}">
                    <input type="checkbox" bind:checked={showPressure} class="rounded border-zinc-300 text-blue-600 focus:ring-blue-500 dark:bg-zinc-800 dark:border-zinc-600" />
                    <span class="font-medium">Давление</span>
                </label>
                <label class="flex items-center gap-2.5 cursor-pointer transition-colors px-2 py-1.5 rounded-md {showFuel ? 'bg-amber-500/15 text-amber-700 dark:bg-amber-500/25 dark:text-amber-300' : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-700/50'}">
                    <input type="checkbox" bind:checked={showFuel} class="rounded border-zinc-300 text-amber-600 focus:ring-amber-500 dark:bg-zinc-800 dark:border-zinc-600" />
                    <span class="font-medium">Топливо</span>
                </label>
            </div>
        {/if}
    </div>
</div>
