<script lang="ts">
    import { telemetryData } from '@/lib/telemetry';
    import WidgetCard from './WidgetCard.svelte';

    $: speed = $telemetryData.speed;

    // Normal speed range check. E.g. 0-120 is normal, 120-140 warning, >140 critical
    $: status = (speed > 140 ? 'critical' : speed > 120 ? 'warning' : 'normal') as 'critical' | 'warning' | 'normal';
</script>

<WidgetCard title="Скорость" icon="🏎️" {status}>
    <div class="text-center">
        <span class="text-5xl font-bold font-mono tracking-tight">{speed.toFixed(1)}</span>
        <span class="text-lg opacity-70 ml-1">км/ч</span>

        <div class="mt-4 w-full bg-zinc-200 dark:bg-zinc-700 rounded-full h-2 overflow-hidden">
            <div
                class={`h-full transition-all duration-300 ease-out
                    ${status === 'critical' ? 'bg-red-500' : status === 'warning' ? 'bg-amber-500' : 'bg-emerald-500'}`}
                style="width: {Math.min((speed / 200) * 100, 100)}%"
            ></div>
        </div>
    </div>
</WidgetCard>
