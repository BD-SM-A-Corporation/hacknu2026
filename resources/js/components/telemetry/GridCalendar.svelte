<script lang="ts">
    import { onMount } from 'svelte';

    let { onRangeSelected } = $props<{
        onRangeSelected: (start: Date, end: Date) => void;
    }>();

    let baseDate = new Date();
    baseDate.setHours(0, 0, 0, 0);
    // Get monday of current week
    const day = baseDate.getDay();
    const diff = baseDate.getDate() - day + (day === 0 ? -6 : 1);
    let startOfWeek = $state(new Date(baseDate.setDate(diff)));
    
    // Derived values
    let days = $derived(Array.from({ length: 7 }).map((_, i) => {
        const d = new Date(startOfWeek);
        d.setDate(d.getDate() + i);
        return d;
    }));
    let hours = Array.from({ length: 24 }).map((_, i) => i);

    let isDragging = $state(false);
    let selectionStart = $state<{dayIdx: number, hour: number} | null>(null);
    let selectionCurrent = $state<{dayIdx: number, hour: number} | null>(null);

    function handleMouseDown(dayIdx: number, hour: number) {
        isDragging = true;
        selectionStart = { dayIdx, hour };
        selectionCurrent = { dayIdx, hour };
    }

    function handleMouseEnter(dayIdx: number, hour: number) {
        if (isDragging) {
            selectionCurrent = { dayIdx, hour };
        }
    }

    function handleMouseUp() {
        if (isDragging && selectionStart && selectionCurrent) {
            isDragging = false;
            
            // Calculate actual date Range
            let d1 = selectionStart.dayIdx;
            let h1 = selectionStart.hour;
            let d2 = selectionCurrent.dayIdx;
            let h2 = selectionCurrent.hour;
            
            // Reorder if dragging backwards in time
            const t1 = d1 * 24 + h1;
            const t2 = d2 * 24 + h2;
            
            let finalStart = { dayIdx: d1, hour: h1 };
            let finalEnd = { dayIdx: d2, hour: h2 };
            
            if (t1 > t2) {
                finalStart = { dayIdx: d2, hour: h2 };
                finalEnd = { dayIdx: d1, hour: h1 };
            }
            
            const startDate = new Date(days[finalStart.dayIdx]);
            startDate.setHours(finalStart.hour, 0, 0, 0);
            
            const endDate = new Date(days[finalEnd.dayIdx]);
            // Make exclusive for API (up to next hour)
            endDate.setHours(finalEnd.hour + 1, 0, 0, 0);
            
            onRangeSelected(startDate, endDate);
        }
    }

    function isSelected(dayIdx: number, hour: number) {
        if (!selectionStart || !selectionCurrent) return false;
        
        const cellTime = dayIdx * 24 + hour;
        
        const t1 = selectionStart.dayIdx * 24 + selectionStart.hour;
        const t2 = selectionCurrent.dayIdx * 24 + selectionCurrent.hour;
        
        const minTime = Math.min(t1, t2);
        const maxTime = Math.max(t1, t2);
        
        return cellTime >= minTime && cellTime <= maxTime;
    }

    function nextWeek() {
        const next = new Date(startOfWeek);
        next.setDate(next.getDate() + 7);
        startOfWeek = next;
    }
    
    function prevWeek() {
        const prev = new Date(startOfWeek);
        prev.setDate(prev.getDate() - 7);
        startOfWeek = prev;
    }
</script>

<div class="border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden bg-white dark:bg-zinc-900 select-none">
    <div class="flex justify-between items-center p-4 border-b border-zinc-200 dark:border-zinc-800">
        <div>
            <h3 class="font-medium text-lg">Выбор исторического периода</h3>
            <p class="text-sm text-zinc-500">Зажмите мышь и проведите по таймлайну, чтобы выбрать период</p>
        </div>
        <div class="flex gap-2">
            <button onclick={prevWeek} class="p-2 border rounded hover:bg-zinc-100 dark:border-zinc-700 dark:hover:bg-zinc-800">&larr;</button>
            <button onclick={nextWeek} class="p-2 border rounded hover:bg-zinc-100 dark:border-zinc-700 dark:hover:bg-zinc-800">&rarr;</button>
        </div>
    </div>
    
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div class="flex overflow-x-auto" onmouseleave={() => { if(isDragging) handleMouseUp() }} onmouseup={handleMouseUp} role="presentation">
        <div class="flex-none w-16 border-r border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900">
            <div class="h-10"></div>
            {#each hours as hour}
                <div class="h-8 flex items-center justify-center text-xs text-zinc-400 border-t border-zinc-100 dark:border-zinc-800/50">
                    {hour}:00
                </div>
            {/each}
        </div>
        
        <div class="flex-1 flex min-w-[700px]">
            {#each days as day, dIdx}
                <div class="flex-1 flex flex-col border-r border-zinc-100 dark:border-zinc-800/50 last:border-0 relative group">
                    <div class="h-10 flex flex-col justify-center items-center font-medium bg-zinc-50 dark:bg-zinc-800/40 text-sm border-b border-zinc-200 dark:border-zinc-800">
                        <span>{day.toLocaleDateString('ru-RU', { weekday: 'short' })}</span>
                        <span class="text-xs text-zinc-500 font-normal">{day.getDate().toString().padStart(2,'0')}.{(day.getMonth()+1).toString().padStart(2,'0')}</span>
                    </div>
                    {#each hours as hour}
                        <!-- svelte-ignore a11y_interactive_supports_focus -->
                        <!-- svelte-ignore a11y_click_events_have_key_events -->
                        <div 
                            class="h-8 border-t border-zinc-100 dark:border-zinc-800/50 transition-colors duration-75 {isSelected(dIdx, hour) ? 'bg-emerald-200 dark:bg-emerald-800' : 'hover:bg-zinc-100 dark:hover:bg-zinc-800 cursor-pointer'}"
                            onmousedown={() => handleMouseDown(dIdx, hour)}
                            onmouseenter={() => handleMouseEnter(dIdx, hour)}
                            role="button"
                        ></div>
                    {/each}
                </div>
            {/each}
        </div>
    </div>
</div>
