<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { get } from 'svelte/store';
    import { Chart, registerables } from 'chart.js';
    import { telemetryData, activeLocomotiveId } from '@/lib/telemetry';

    Chart.register(...registerables);

    let canvas: HTMLCanvasElement;
    let chart: Chart;

    const maxDataPoints = 60;

    // Data buffers
    let labels: string[] = [];
    let speedData: number[] = [];
    let temperatureData: number[] = [];
    let pressureData: number[] = [];
    let fuelData: number[] = [];

    // Toggles
    let showSpeed = $state(true);
    let showTemperature = $state(true);
    let showPressure = $state(true);
    let showFuel = $state(true);

    let unsubscribe: () => void;

    onMount(async () => {
        chart = new Chart(canvas, {
            // ... existing config ...
            type: 'line',
            data: {
                labels,
                datasets: [
                    {
                        label: 'Скорость (км/ч)',
                        data: speedData,
                        borderColor: '#059669', // Emerald 600
                        tension: 0.4,
                        hidden: !showSpeed
                    },
                    {
                        label: 'Температура (°C)',
                        data: temperatureData,
                        borderColor: '#dc2626', // Red 600
                        tension: 0.4,
                        hidden: !showTemperature
                    },
                    {
                        label: 'Давление (атм)',
                        data: pressureData,
                        borderColor: '#2563eb', // Blue 600
                        tension: 0.4,
                        hidden: !showPressure
                    },
                    {
                        label: 'Топливо (%)',
                        data: fuelData,
                        borderColor: '#d97706', // Amber 600
                        tension: 0.4,
                        hidden: !showFuel
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: false, // For better realtime performance without jumping
                interaction: {
                    mode: 'index',
                    intersect: false,
                },
                scales: {
                    x: {
                        display: true,
                    },
                    y: {
                        beginAtZero: true
                    }
                },
                plugins: {
                    legend: {
                        display: false // We use our own toggles
                    }
                }
            }
        });

        const id = get(activeLocomotiveId);
        if (id) {
            // 1. Try to load from localStorage first
            const cached = localStorage.getItem(`telemetry_chart_${id}`);
            if (cached) {
                try {
                    const data = JSON.parse(cached);
                    labels = data.labels || [];
                    speedData = data.speedData || [];
                    temperatureData = data.temperatureData || [];
                    pressureData = data.pressureData || [];
                    fuelData = data.fuelData || [];

                    // Update chart with cached data
                    chart.data.labels = labels;
                    chart.data.datasets[0].data = speedData;
                    chart.data.datasets[1].data = temperatureData;
                    chart.data.datasets[2].data = pressureData;
                    chart.data.datasets[3].data = fuelData;
                    chart.update();
                } catch (e) {
                    console.error("Failed to parse cached chart data", e);
                }
            }

            // 2. Refresh from API to get any missed data
            try {
                const res = await fetch(`/api/locomotives/${id}/stats`);
                if (res.ok) {
                    const history = await res.json();

                    if (history.length > 0) {
                        const sortedHistory = history.reverse();
                        const newLabels: any[] = [];
                        const newSpeed: any[] = [];
                        const newTemp: any[] = [];
                        const newPress: any[] = [];
                        const newFuel: any[] = [];

                        sortedHistory.forEach((record: any) => {
                            const time = new Date(record.timestamp).toLocaleTimeString();
                            newLabels.push(time);
                            newSpeed.push(record.speed);
                            newTemp.push(record.temperature);
                            newPress.push(record.pressure);
                            newFuel.push(record.fuel_level || record.fuelLevel);
                        });

                        labels = newLabels;
                        speedData = newSpeed;
                        temperatureData = newTemp;
                        pressureData = newPress;
                        fuelData = newFuel;

                        chart.data.labels = labels;
                        chart.data.datasets[0].data = speedData;
                        chart.data.datasets[1].data = temperatureData;
                        chart.data.datasets[2].data = pressureData;
                        chart.data.datasets[3].data = fuelData;
                        chart.update();

                        // Update cache
                        localStorage.setItem(`telemetry_chart_${id}`, JSON.stringify({
                            labels, speedData, temperatureData, pressureData, fuelData
                        }));
                    }
                }
            } catch (err) {
                console.error("Failed to fetch historical stats", err);
            }
        }

        unsubscribe = telemetryData.subscribe((data) => {
            if (!data.locomotiveId || (data.locomotiveId !== id && id !== '')) return; // Ensure we stick to the initial ID for this modal instance

            const timeString = new Date(data.timestamp).toLocaleTimeString();

            // Check if last added was same timestamp (avoid duplicates from fast pings)
            if (labels.length > 0 && labels[labels.length - 1] === timeString) return;

            labels.push(timeString);
            speedData.push(data.speed);
            temperatureData.push(data.temperature);
            pressureData.push(data.pressure);
            fuelData.push(data.fuelLevel);

            if (labels.length > maxDataPoints) {
                labels.shift();
                speedData.shift();
                temperatureData.shift();
                pressureData.shift();
                fuelData.shift();
            }

            if (chart) {
                chart.update();
                // Save to localStorage on every update
                if (id) {
                    localStorage.setItem(`telemetry_chart_${id}`, JSON.stringify({
                        labels, speedData, temperatureData, pressureData, fuelData
                    }));
                }
            }
        });
    });

    onDestroy(() => {
        if (unsubscribe) unsubscribe();
        if (chart) chart.destroy();
    });

    // Reactively update chart visibility when toggles change
    $effect(() => {
        if (chart) {
            chart.data.datasets[0].hidden = !showSpeed;
            chart.data.datasets[1].hidden = !showTemperature;
            chart.data.datasets[2].hidden = !showPressure;
            chart.data.datasets[3].hidden = !showFuel;
            chart.update();
        }
    });
</script>

<div class="relative bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl p-6 shadow-sm min-h-[300px] h-[400px]">
    <!-- Custom Controls over the chart -->
    <div class="absolute top-4 right-4 z-10 bg-white/80 dark:bg-zinc-800/80 backdrop-blur-sm p-3 border border-zinc-200 dark:border-zinc-700 rounded-lg shadow-sm flex flex-col gap-2 text-sm">
        <span class="font-medium text-xs text-zinc-500 uppercase tracking-wider mb-1">Отображение</span>
        <label class="flex items-center gap-2 cursor-pointer transition-opacity hover:opacity-80">
            <input type="checkbox" bind:checked={showSpeed} class="rounded border-zinc-300 text-emerald-600 focus:ring-emerald-500 dark:bg-zinc-900" />
            <span class="text-zinc-700 dark:text-zinc-300">Скорость</span>
        </label>
        <label class="flex items-center gap-2 cursor-pointer transition-opacity hover:opacity-80">
            <input type="checkbox" bind:checked={showTemperature} class="rounded border-zinc-300 text-red-600 focus:ring-red-500 dark:bg-zinc-900" />
            <span class="text-zinc-700 dark:text-zinc-300">Температура</span>
        </label>
        <label class="flex items-center gap-2 cursor-pointer transition-opacity hover:opacity-80">
            <input type="checkbox" bind:checked={showPressure} class="rounded border-zinc-300 text-blue-600 focus:ring-blue-500 dark:bg-zinc-900" />
            <span class="text-zinc-700 dark:text-zinc-300">Давление</span>
        </label>
        <label class="flex items-center gap-2 cursor-pointer transition-opacity hover:opacity-80">
            <input type="checkbox" bind:checked={showFuel} class="rounded border-zinc-300 text-amber-600 focus:ring-amber-500 dark:bg-zinc-900" />
            <span class="text-zinc-700 dark:text-zinc-300">Топливо</span>
        </label>
    </div>

    <!-- Chart container -->
    <div class="w-full h-full pt-4 pr-32">
        <canvas bind:this={canvas}></canvas>
    </div>
</div>
