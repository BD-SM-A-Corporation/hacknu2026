<script lang="ts">
    import { telemetryData } from '@/lib/telemetry';
    import WidgetCard from './WidgetCard.svelte';

    export let overrideTemperature: number | undefined = undefined;
    $: temperature = overrideTemperature !== undefined ? overrideTemperature : $telemetryData.temperature;

    // Engine temp logic: < 90 normal, 90-105 warning, > 105 critical
    $: status = (temperature > 105 ? 'critical' : temperature > 90 ? 'warning' : 'normal') as 'critical' | 'warning' | 'normal';
</script>

<WidgetCard title="Температура двигателя" icon="🌡️" {status}>
    <div class="text-center">
        <span class="text-5xl font-bold font-mono tracking-tight">{temperature.toFixed(1)}</span>
        <span class="text-lg opacity-70 ml-1">°C</span>

        {#if status === 'critical'}
            <div class="mt-2 text-xs font-bold uppercase animate-pulse">Опасно высокая температура</div>
        {:else if status === 'warning'}
            <div class="mt-2 text-xs font-medium uppercase">Температура выше нормы</div>
        {:else}
            <div class="mt-2 text-xs font-medium uppercase opacity-50">В пределах нормы</div>
        {/if}
    </div>
</WidgetCard>
