<script lang="ts">
    import WidgetCard from './WidgetCard.svelte';
    import { telemetryData } from '@/lib/telemetry';

    export let overridePressure: number | undefined = undefined;
    $: pressure = overridePressure !== undefined ? overridePressure : $telemetryData.pressure;

    // Pressure logic: normal 4-8. warning 3-4 or 8-10. critical < 3 or > 10.
    $: status = (pressure < 3 || pressure > 10) ? 'critical' :
                (pressure < 4 || pressure > 8) ? 'warning' : 'normal';

    // Convert to percentage for visual bar (assuming 15 atm max)
    $: percent = Math.min((pressure / 15) * 100, 100);
</script>

<WidgetCard title="Давление масла" icon="🗜️" {status}>
    <div class="whitespace-nowrap flex flex-col items-center w-full">
        <div class="text-center">
            <span class="text-5xl font-bold font-mono tracking-tight">{pressure.toFixed(2)}</span>
            <span class="text-lg opacity-70 ml-1">атм</span>
        </div>

        <div class="mt-4 flex w-full justify-between items-end gap-1 h-8 opacity-80 px-2">
            {#each Array(15) as _, i}
                <div class={`w-full rounded-sm transition-all duration-300
                    ${(i+1) <= (pressure)
                        ? (status === 'critical' ? 'bg-red-500' : status === 'warning' ? 'bg-amber-500' : 'bg-emerald-500')
                        : 'bg-zinc-200 dark:bg-zinc-700 h-2'
                    }`}
                    style="height: {(i+1) <= pressure ? '100%' : '20%'}"
                ></div>
            {/each}
        </div>

        {#if status === 'critical'}
            <div class="mt-2 text-xs font-bold uppercase animate-pulse">Критическое давление!</div>
        {/if}
    </div>
</WidgetCard>
