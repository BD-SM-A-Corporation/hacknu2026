<script lang="ts">
    import { telemetryData } from '@/lib/telemetry';
    import WidgetCard from './WidgetCard.svelte';

    export let overrideHealth: number | undefined = undefined;
    $: score = overrideHealth !== undefined ? overrideHealth : $telemetryData.healthScore;

    // Normal health is > 85, warning > 60, critical < 60
    $: status = (score < 60 ? 'critical' : score < 85 ? 'warning' : 'normal') as 'critical' | 'warning' | 'normal';
</script>

<WidgetCard title="Индекс Здоровья" icon="❤️" {status}>
    <div class="text-center">
        <span class="text-5xl font-bold font-mono tracking-tight">{score.toFixed(0)}</span>
        <span class="text-lg opacity-70 ml-1">%</span>

        <div class="mt-4 w-full bg-zinc-200 dark:bg-zinc-700 rounded-full h-2 overflow-hidden flex">
            <div
                class={`h-full transition-all duration-300 ease-out
                    ${status === 'critical' ? 'bg-red-500' : status === 'warning' ? 'bg-amber-500' : 'bg-emerald-500'}`}
                style="width: {score}%"
            ></div>
        </div>
    </div>
</WidgetCard>
