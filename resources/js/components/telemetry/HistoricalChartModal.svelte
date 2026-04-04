<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { activeLocomotiveId } from '@/lib/telemetry';
    import { Chart, registerables } from 'chart.js';

    Chart.register(...registerables);

    let { isOpen = $bindable(false), startDate, endDate } = $props<{
        isOpen: boolean;
        startDate: Date | null;
        endDate: Date | null;
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
    
    let stepSelection = $state("auto");

    $effect(() => {
        if (isOpen && startDate && endDate && $activeLocomotiveId) {
            fetchData();
        }
    });

    async function fetchData() {
        loading = true;
        hasData = false;
        try {
            let stepParam = "";
            let periodSecs = (endDate!.getTime() - startDate!.getTime()) / 1000;
            
            if (stepSelection === "auto") {
                // Auto: try to get roughly ~60 points
                if (periodSecs > 60) {
                    const calculatedStep = Math.round(periodSecs / 60);
                    stepParam = `&step=${calculatedStep}`;
                }
            } else if (stepSelection !== "1s") {
                // E.g. "5s" => 5, "1m" => 60
                const num = parseInt(stepSelection);
                const unit = stepSelection.slice(-1);
                let multiplier = 1;
                if (unit === "m") multiplier = 60;
                if (unit === "h") multiplier = 3600;
                if (unit === "d") multiplier = 86400;
                stepParam = `&step=${num * multiplier}`;
            }

            const res = await fetch(`/api/telemetry/analytics?locomotive_id=${$activeLocomotiveId}&start_date=${startDate!.toISOString()}&end_date=${endDate!.toISOString()}${stepParam}`);
            if (res.ok) {
                const data = await res.json();
                hasData = data.length > 0;
                currentData = data;
                // Add a small delay for DOM to render the canvas if it was hidden
                setTimeout(() => {
                    renderChart(data);
                }, 50);
            }
        } catch(e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    function downloadChart() {
        if (!chartCanvas || !chart) return;
        const link = document.createElement('a');
        link.download = `analytics_KZ8A_${$activeLocomotiveId}.png`;
        link.href = chart.toBase64Image();
        link.click();
    }

    function downloadCSV() {
        if (!currentData.length) return;
        
        let csvContent = "Дата и Время;Скорость (км/ч);Температура (°C);Давление;Уровень топлива\n";
        currentData.forEach(item => {
            const dt = new Date(item.timestamp).toLocaleString('ru-RU');
            const speed = Number(item.speed).toFixed(2);
            const temp = Number(item.temperature).toFixed(2);
            const press = Number(item.pressure).toFixed(2);
            const fuel = Number(item.fuel_level).toFixed(2);
            csvContent += `"${dt}";${speed};${temp};${press};${fuel}\n`;
        });
        
        const blob = new Blob(["\uFEFF" + csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = `analytics_KZ8A_${$activeLocomotiveId}.csv`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    function renderChart(data: any[]) {
        if (!chartCanvas) return;
        if (chart) chart.destroy();
        if (!hasData) return;

        // Отключаем фронтенд-прореживание, так как сервер уже агрегировал данные
        let periodSecs = (endDate!.getTime() - startDate!.getTime()) / 1000;
        const showDate = periodSecs >= 86400; // Если период больше 1 дня, то показываем также и дату
        
        const labels = data.map(d => {
            const dt = new Date(d.timestamp);
            if (showDate) {
                return dt.toLocaleString('ru-RU', { day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit' });
            }
            return dt.toLocaleTimeString([], { hour12: false });
        });
        
        // Убеждаемся, что парсим в Number т.к. AVG() в postgres возвращает строки для float/decimal типов
        const speed = data.map(d => Number(d.speed));
        const temp = data.map(d => Number(d.temperature));
        const press = data.map(d => Number(d.pressure));
        const fuel = data.map(d => Number(d.fuel_level));

        chart = new Chart(chartCanvas, {
            type: 'line',
            data: {
                labels,
                datasets: [
                    { label: 'Скорость', data: speed, borderColor: '#059669', hidden: !showSpeed, tension: 0.2, borderWidth: 1.5, yAxisID: 'y' },
                    { label: 'Температура', data: temp, borderColor: '#dc2626', hidden: !showTemperature, tension: 0.2, borderWidth: 1.5, yAxisID: 'yTemp' },
                    { label: 'Давление', data: press, borderColor: '#2563eb', hidden: !showPressure, tension: 0.2, borderWidth: 1.5, yAxisID: 'yPress' },
                    { label: 'Топливо', data: fuel, borderColor: '#d97706', hidden: !showFuel, tension: 0.2, borderWidth: 1.5, yAxisID: 'y' },
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { display: false } },
                elements: { point: { radius: 0 } },
                interaction: { mode: 'index', intersect: false },
                scales: {
                    x: {
                        ticks: { maxTicksLimit: 15 } // Не перегружаем ось X текстом
                    },
                    y: {
                        type: 'linear',
                        display: true,
                        position: 'left',
                        title: { display: true, text: 'Шкала: Скорость / Топливо' }
                    },
                    yTemp: {
                        type: 'linear',
                        display: showTemperature,
                        position: 'right',
                        title: { display: true, text: 'Температура (°C)' },
                        grid: { drawOnChartArea: false } // Оставляем сетку только для главной оси
                    },
                    yPress: {
                        type: 'linear',
                        display: showPressure,
                        position: 'right',
                        title: { display: true, text: 'Давление' },
                        grid: { drawOnChartArea: false }
                    }
                }
            },
            plugins: [{
                id: 'custom_canvas_background_color',
                beforeDraw: (chartInstance) => {
                    const ctx = chartInstance.canvas.getContext('2d');
                    if (!ctx) return; // Fix TS lint
                    ctx.save();
                    ctx.globalCompositeOperation = 'destination-over';
                    const isDark = document.documentElement.classList.contains('dark');
                    ctx.fillStyle = isDark ? '#18181b' : '#ffffff'; 
                    ctx.fillRect(0, 0, chartInstance.width, chartInstance.height);
                    ctx.restore();
                }
            }]
        });
    }

    $effect(() => {
        if (chart) {
            chart.data.datasets[0].hidden = !showSpeed;
            chart.data.datasets[1].hidden = !showTemperature;
            chart.data.datasets[2].hidden = !showPressure;
            chart.data.datasets[3].hidden = !showFuel;
            
            // Динамически скрываем оси, если метрика отключена
            if (chart.options?.scales?.yTemp) {
                chart.options.scales.yTemp.display = showTemperature;
            }
            if (chart.options?.scales?.yPress) {
                chart.options.scales.yPress.display = showPressure;
            }
            
            chart.update();
        }
    });

    function close() {
        isOpen = false;
        if (chart) {
            chart.destroy();
            chart = null;
        }
    }
</script>

<svelte:window onkeydown={(e) => { if (e.key === 'Escape' && isOpen) close(); }} />

{#if isOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" onclick={close}>
        <div class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl shadow-2xl w-full max-w-6xl flex flex-col h-[85vh]" onclick={(e) => e.stopPropagation()}>
            <div class="p-4 border-b border-zinc-200 dark:border-zinc-800 flex justify-between items-center bg-zinc-50 dark:bg-zinc-800/30 rounded-t-xl">
                <h2 class="text-lg font-bold">
                    Архив телеметрии 
                    <span class="text-emerald-600 font-mono ml-2 block sm:inline">{$activeLocomotiveId}</span>
                </h2>
                <div class="flex items-center gap-4 text-sm text-zinc-500">
                    <div>{startDate?.toLocaleString('ru-RU')} &mdash; {endDate?.toLocaleString('ru-RU')}</div>
                    
                    <div class="flex items-center gap-2 relative z-20">
                        <span class="text-xs">Шаг:</span>
                        <select bind:value={stepSelection} onchange={() => fetchData()} class="p-1 rounded bg-zinc-100 border-zinc-200 dark:bg-zinc-800 dark:border-zinc-700 focus:ring-emerald-500">
                            <option value="auto">Авто</option>
                            <option value="1s">1 сек</option>
                            <option value="5s">5 сек</option>
                            <option value="10s">10 сек</option>
                            <option value="30s">30 сек</option>
                            <option value="1m">1 мин</option>
                            <option value="5m">5 мин</option>
                            <option value="10m">10 мин</option>
                            <option value="30m">30 мин</option>
                            <option value="1h">1 час</option>
                            <option value="3h">3 часа</option>
                            <option value="6h">6 часов</option>
                            <option value="12h">12 часов</option>
                            <option value="1d">1 день</option>
                        </select>
                    </div>
                </div>
                
                <div class="flex items-center gap-2">
                    <button 
                        onclick={downloadCSV} 
                        class="px-3 py-1.5 flex items-center gap-1 bg-zinc-100 hover:bg-zinc-200 text-zinc-700 dark:bg-zinc-800 dark:hover:bg-zinc-700 dark:text-zinc-300 text-sm rounded shadow-sm transition"
                        title="Скачать данные как CSV (Excel)"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><path d="M8 13h2"/><path d="M8 17h2"/><path d="M14 13h2"/><path d="M14 17h2"/></svg>
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
            
            <div class="flex-1 p-6 flex flex-col relative bg-zinc-50/50 dark:bg-zinc-900/50 overflow-hidden">
                {#if loading}
                    <div class="flex-1 flex justify-center items-center">
                        <div class="flex flex-col items-center gap-4 text-emerald-600">
                            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600"></div>
                            <p class="font-medium">Загрузка данных...</p>
                        </div>
                    </div>
                {:else if !hasData}
                    <div class="flex-1 flex flex-col justify-center items-center text-zinc-500 gap-4">
                        <svg class="w-16 h-16 opacity-30" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                           <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        Аналитические данные для выбранного локомотива в данном периоде отсутствуют.
                    </div>
                {/if}
                
                {#if !loading && hasData}
                    <div class="absolute top-6 right-6 z-10 bg-white/90 dark:bg-zinc-800/90 backdrop-blur-md p-3 border border-zinc-200 dark:border-zinc-700 rounded-lg shadow-sm flex flex-col gap-2 text-sm font-medium">
                        <span class="text-xs uppercase text-zinc-500 mb-1">Метрики</span>
                        <label class="flex items-center gap-2 cursor-pointer hover:opacity-80">
                            <input type="checkbox" bind:checked={showSpeed} class="rounded border-zinc-300 text-emerald-600 focus:ring-emerald-500">
                            <span class="text-zinc-700 dark:text-zinc-300">Скорость</span>
                        </label>
                        <label class="flex items-center gap-2 cursor-pointer hover:opacity-80">
                            <input type="checkbox" bind:checked={showTemperature} class="rounded border-zinc-300 text-red-600 focus:ring-red-500">
                            <span class="text-zinc-700 dark:text-zinc-300">Температура</span>
                        </label>
                        <label class="flex items-center gap-2 cursor-pointer hover:opacity-80">
                            <input type="checkbox" bind:checked={showPressure} class="rounded border-zinc-300 text-blue-600 focus:ring-blue-500">
                            <span class="text-zinc-700 dark:text-zinc-300">Давление</span>
                        </label>
                        <label class="flex items-center gap-2 cursor-pointer hover:opacity-80">
                            <input type="checkbox" bind:checked={showFuel} class="rounded border-zinc-300 text-amber-600 focus:ring-amber-500">
                            <span class="text-zinc-700 dark:text-zinc-300">Топливо</span>
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
