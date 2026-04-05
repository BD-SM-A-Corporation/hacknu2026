<script lang="ts">
    import { telemetryData } from '@/lib/telemetry';
    import WidgetCard from './WidgetCard.svelte';

    let { overrideFuel = undefined }: { overrideFuel?: number } = $props();
    const fuelLevel = $derived(overrideFuel !== undefined ? overrideFuel : $telemetryData.fuelLevel);
    const status = $derived(
        fuelLevel <= 10 ? 'critical' : fuelLevel <= 20 ? 'warning' : 'normal',
    );
</script>

<WidgetCard title="Уровень топлива" icon="⛽" {status}>
    <div class="flex items-center justify-center w-full relative">
        <svg viewBox="0 0 100 100" class="w-28 h-28 transform -rotate-90">
            <circle
                cx="50"
                cy="50"
                r="40"
                class="stroke-zinc-200 dark:stroke-zinc-700 fill-none"
                stroke-width="10"
            />
            <circle
                cx="50"
                cy="50"
                r="40"
                class={`fill-none transition-all duration-1000 ease-out
                    ${status === 'critical' ? 'stroke-red-500' : status === 'warning' ? 'stroke-amber-500' : 'stroke-emerald-500'}`}
                stroke-width="10"
                stroke-dasharray="251.2"
                stroke-dashoffset={251.2 - (251.2 * fuelLevel) / 100}
                stroke-linecap="round"
            />
        </svg>
        <div class="absolute flex flex-col items-center justify-center">
            <span class="text-2xl font-bold font-mono"
                >{Math.round(fuelLevel)}%</span
            >
        </div>
    </div>
</WidgetCard>
