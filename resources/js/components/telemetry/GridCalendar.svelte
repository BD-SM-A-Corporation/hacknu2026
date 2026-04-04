<script lang="ts">
    import { onMount } from 'svelte';

    let { onRangeSelected } = $props<{
        onRangeSelected: (start: Date, end: Date) => void;
    }>();

    type ViewMode = 'month' | 'week' | 'day';
    let viewMode = $state<ViewMode>('week');
    
    let baseDate = $state(new Date());
    
    // Derived date helpers
    let monthDays = $derived((() => {
        const year = baseDate.getFullYear();
        const month = baseDate.getMonth();
        const firstDay = new Date(year, month, 1);
        const lastDay = new Date(year, month + 1, 0);
        
        const startOffset = (firstDay.getDay() + 6) % 7; // monday=0
        
        let days: Date[] = [];
        for(let i=startOffset; i>0; i--) {
            days.push(new Date(year, month, 1 - i));
        }
        for(let i=1; i<=lastDay.getDate(); i++) {
            days.push(new Date(year, month, i));
        }
        const endOffset = (7 - days.length % 7) % 7;
        for(let i=1; i<=endOffset; i++) {
            days.push(new Date(year, month + 1, i));
        }
        return days;
    })());

    let weekDays = $derived((() => {
        const day = baseDate.getDay();
        const diff = baseDate.getDate() - day + (day === 0 ? -6 : 1);
        let days: Date[] = [];
        for (let i = 0; i < 7; i++) {
            const d = new Date(baseDate);
            d.setHours(0, 0, 0, 0);
            d.setDate(diff + i);
            days.push(d);
        }
        return days;
    })());

    let headerText = $derived((() => {
        if (viewMode === 'month') {
            return baseDate.toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' });
        } else if (viewMode === 'week') {
            const start = weekDays[0];
            const end = weekDays[6];
            if (start.getMonth() === end.getMonth()) {
                return `${start.getDate()} - ${end.getDate()} ${start.toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' })}`;
            } else {
                return `${start.getDate()} ${start.toLocaleDateString('ru-RU', { month: 'short' })} - ${end.getDate()} ${end.toLocaleDateString('ru-RU', { month: 'short', year: 'numeric' })}`;
            }
        } else {
            return baseDate.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' });
        }
    })());

    // Selection mechanics
    let startMark = $state<number | null>(null);
    let currentMark = $state<number | null>(null);
    let isDragging = $state(false);

    let selectedRange = $derived((() => {
        if (startMark === null || currentMark === null) return null;
        return {
            min: Math.min(startMark, currentMark),
            max: Math.max(startMark, currentMark)
        };
    })());

    function isSelected(ts: number, te: number) {
        if (!selectedRange) return false;
        return ts < selectedRange.max && te > selectedRange.min;
    }

    function handleMouseDown(ts: number, te: number) {
        isDragging = true;
        startMark = ts;
        currentMark = te;
    }

    function handleMouseEnter(ts: number, te: number) {
        if (!isDragging || startMark === null) return;
        if (ts < startMark) {
            currentMark = ts;
        } else {
            currentMark = te;
        }
    }

    function handleMouseUp() {
        if (isDragging && startMark !== null && currentMark !== null) {
            isDragging = false;
            const min = Math.min(startMark, currentMark);
            const max = Math.max(startMark, currentMark);
            onRangeSelected(new Date(min), new Date(max));
        }
    }

    function nextPeriod() {
        const next = new Date(baseDate);
        if (viewMode === 'month') {
            next.setMonth(next.getMonth() + 1);
        } else if (viewMode === 'week') {
            next.setDate(next.getDate() + 7);
        } else {
            next.setDate(next.getDate() + 1);
        }
        baseDate = next;
    }
    
    function prevPeriod() {
        const prev = new Date(baseDate);
        if (viewMode === 'month') {
            prev.setMonth(prev.getMonth() - 1);
        } else if (viewMode === 'week') {
            prev.setDate(prev.getDate() - 7);
        } else {
            prev.setDate(prev.getDate() - 1);
        }
        baseDate = prev;
    }
</script>

<div class="border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden bg-white dark:bg-zinc-900 select-none">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center p-4 border-b border-zinc-200 dark:border-zinc-800 gap-4">
        <div>
            <h3 class="font-medium text-lg capitalize">{headerText}</h3>
            <p class="text-sm text-zinc-500">Зажмите мышь и тяните, чтобы выделить период времени</p>
        </div>
        <div class="flex items-center gap-4">
            <!-- View Mode Switcher -->
            <div class="flex p-1 bg-zinc-100 dark:bg-zinc-800 rounded-lg">
                <button 
                    onclick={() => viewMode = 'month'} 
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors {viewMode === 'month' ? 'bg-white dark:bg-zinc-700 shadow-sm text-zinc-900 dark:text-zinc-100' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'}">
                    Месяц
                </button>
                <button 
                    onclick={() => viewMode = 'week'} 
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors {viewMode === 'week' ? 'bg-white dark:bg-zinc-700 shadow-sm text-zinc-900 dark:text-zinc-100' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'}">
                    Неделя
                </button>
                <button 
                    onclick={() => viewMode = 'day'} 
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors {viewMode === 'day' ? 'bg-white dark:bg-zinc-700 shadow-sm text-zinc-900 dark:text-zinc-100' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200'}">
                    День
                </button>
            </div>
            
            <!-- Navigation -->
            <div class="flex gap-1.5">
                <button onclick={prevPeriod} aria-label="Предыдущий период" class="p-1.5 border border-zinc-200 dark:border-zinc-700 rounded-lg hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors">
                    <svg class="w-5 h-5 text-zinc-600 dark:text-zinc-400" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
                </button>
                <button onclick={nextPeriod} aria-label="Следующий период" class="p-1.5 border border-zinc-200 dark:border-zinc-700 rounded-lg hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors">
                    <svg class="w-5 h-5 text-zinc-600 dark:text-zinc-400" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
                </button>
            </div>
        </div>
    </div>
    
    <!-- Calendar Body -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div class="relative w-full" onmouseleave={() => { if(isDragging) handleMouseUp(); }} onmouseup={handleMouseUp} role="presentation">
        
        {#if viewMode === 'month'}
            <!-- MONTH VIEW -->
            <div class="grid grid-cols-7 border-b border-zinc-200 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-900">
                {#each ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс'] as d}
                    <div class="py-2.5 text-center text-xs font-medium text-zinc-500 uppercase tracking-wider border-r last:border-0 border-zinc-200 dark:border-zinc-800">{d}</div>
                {/each}
            </div>
            <div class="grid grid-cols-7 auto-rows-fr">
                {#each monthDays as day}
                    {@const ts = new Date(day).setHours(0,0,0,0)}
                    {@const te = new Date(day).setHours(24,0,0,0)}
                    <!-- svelte-ignore a11y_interactive_supports_focus -->
                    <!-- svelte-ignore a11y_click_events_have_key_events -->
                    <div 
                        class="min-h-24 p-2 border-b border-r border-zinc-100 dark:border-zinc-800/50 transition-colors duration-75 select-none
                        {day.getMonth() !== baseDate.getMonth() ? 'opacity-40 bg-zinc-50/50 dark:bg-zinc-900/50' : 'bg-white dark:bg-zinc-900'}
                        {isSelected(ts, te) ? 'bg-emerald-100/60 dark:bg-emerald-900/40' : 'hover:bg-zinc-50 dark:hover:bg-zinc-800/60 cursor-pointer'}"
                        onmousedown={() => handleMouseDown(ts, te)}
                        onmouseenter={() => handleMouseEnter(ts, te)}
                        role="button" tabindex="0"
                    >
                        <div class="font-medium text-sm flex items-center justify-center {day.getDate() === new Date().getDate() && day.getMonth() === new Date().getMonth() ? 'bg-emerald-500 text-white rounded-full w-7 h-7 shadow-sm' : 'text-zinc-700 dark:text-zinc-300 w-7 h-7'}">
                            {day.getDate()}
                        </div>
                    </div>
                {/each}
            </div>
            
        {:else if viewMode === 'week'}
            <!-- WEEK VIEW -->
            <div class="flex overflow-x-auto min-h-[400px]">
                <div class="flex-none w-16 border-r border-zinc-200 dark:border-zinc-800 bg-zinc-50/80 dark:bg-zinc-900/80 sticky left-0 z-10 backdrop-blur-md">
                    <div class="h-12 border-b border-zinc-200 dark:border-zinc-800"></div>
                    {#each Array.from({length: 24}) as _, hour}
                        <div class="h-10 flex items-center justify-center text-xs text-zinc-400 font-medium border-t border-zinc-100 dark:border-zinc-800/50">
                            {hour.toString().padStart(2, '0')}:00
                        </div>
                    {/each}
                </div>
                
                <div class="flex-1 flex min-w-[700px]">
                    {#each weekDays as day}
                        <div class="flex-1 flex flex-col border-r border-zinc-100 dark:border-zinc-800/50 last:border-0 relative">
                            <div class="h-12 flex flex-col justify-center items-center bg-zinc-50 dark:bg-zinc-800/40 border-b border-zinc-200 dark:border-zinc-800 sticky top-0 z-10 backdrop-blur-md">
                                <span class="font-medium text-zinc-700 dark:text-zinc-300 text-sm capitalize">{day.toLocaleDateString('ru-RU', { weekday: 'short' })}</span>
                                <span class="text-xs text-zinc-500">{day.getDate().toString().padStart(2,'0')}.{(day.getMonth()+1).toString().padStart(2,'0')}</span>
                            </div>
                            {#each Array.from({length: 24}) as _, hour}
                                {@const ts = new Date(day).setHours(hour, 0, 0, 0)}
                                {@const te = new Date(day).setHours(hour + 1, 0, 0, 0)}
                                <!-- svelte-ignore a11y_interactive_supports_focus -->
                                <!-- svelte-ignore a11y_click_events_have_key_events -->
                                <div 
                                    class="h-10 border-t border-zinc-100 dark:border-zinc-800/50 transition-colors duration-75 
                                    {isSelected(ts, te) ? 'bg-emerald-200/60 dark:bg-emerald-800/60 shadow-[inset_0_0_0_1px_rgba(16,185,129,0.2)]' : 'hover:bg-zinc-100 dark:hover:bg-zinc-800 cursor-pointer'}"
                                    onmousedown={() => handleMouseDown(ts, te)}
                                    onmouseenter={() => handleMouseEnter(ts, te)}
                                    role="button" tabindex="0"
                                ></div>
                            {/each}
                        </div>
                    {/each}
                </div>
            </div>
            
        {:else}
            <!-- DAY VIEW -->
            <div class="flex overflow-x-auto min-h-[400px]">
                <div class="flex-none w-16 border-r border-zinc-200 dark:border-zinc-800 bg-zinc-50/80 dark:bg-zinc-900/80 sticky left-0 z-10">
                    <div class="h-14 border-b border-zinc-200 dark:border-zinc-800"></div>
                    {#each Array.from({length: 24}) as _, hour}
                        <div class="h-14 flex items-center justify-center text-xs text-zinc-500 font-medium border-t border-zinc-100 dark:border-zinc-800/50">
                            {hour.toString().padStart(2, '0')}:00
                        </div>
                    {/each}
                </div>
                
                <div class="flex-1 flex flex-col bg-white dark:bg-zinc-900 border-r border-zinc-100 dark:border-zinc-800/50">
                    <div class="h-14 flex flex-col justify-center items-center bg-zinc-50 dark:bg-zinc-800/40 border-b border-zinc-200 dark:border-zinc-800 sticky top-0 z-10 flex-shrink-0">
                        <span class="text-zinc-700 dark:text-zinc-200 font-medium capitalize">{baseDate.toLocaleDateString('ru-RU', { weekday: 'long' })}</span>
                        <span class="text-sm font-bold text-emerald-600 dark:text-emerald-500">{baseDate.toLocaleDateString('ru-RU')}</span>
                    </div>
                    <div class="flex-1 flex flex-col pt-px">
                        {#each Array.from({length: 24}) as _, hour}
                            {@const ts = new Date(baseDate).setHours(hour, 0, 0, 0)}
                            {@const te = new Date(baseDate).setHours(hour + 1, 0, 0, 0)}
                            <!-- svelte-ignore a11y_interactive_supports_focus -->
                            <!-- svelte-ignore a11y_click_events_have_key_events -->
                            <div 
                                class="h-14 border-t border-zinc-100 dark:border-zinc-800/50 transition-colors duration-75 flex items-center px-4 group
                                {isSelected(ts, te) ? 'bg-emerald-100/80 dark:bg-emerald-900/50 border-l-4 border-l-emerald-500' : 'hover:bg-zinc-50 dark:hover:bg-zinc-800/60 cursor-pointer border-l-4 border-l-transparent'}"
                                onmousedown={() => handleMouseDown(ts, te)}
                                onmouseenter={() => handleMouseEnter(ts, te)}
                                role="button" tabindex="0"
                            >
                                <span class="text-xs opacity-0 group-hover:opacity-100 transition-opacity {isSelected(ts,te) ? 'text-emerald-700 dark:text-emerald-400 font-medium' : 'text-zinc-400'}">
                                    Нажми и тяни
                                </span>
                            </div>
                        {/each}
                    </div>
                </div>
            </div>
        {/if}
        
    </div>
</div>
