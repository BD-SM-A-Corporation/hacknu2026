import { writable, derived } from 'svelte/store';

export interface DateRange {
    start: Date;
    end: Date;
}

function startOfDay(date: Date): Date {
    const d = new Date(date);
    d.setHours(0, 0, 0, 0);
    return d;
}

function endOfDay(date: Date): Date {
    const d = new Date(date);
    d.setHours(23, 59, 59, 999);
    return d;
}

function subtractDays(date: Date, days: number): Date {
    const d = new Date(date);
    d.setDate(d.getDate() - days);
    return d;
}

function subtractMonths(date: Date, months: number): Date {
    const d = new Date(date);
    d.setMonth(d.getMonth() - months);
    return d;
}

export type PresetKey = 'week' | 'month' | '3months' | '6months' | 'custom';

export const presets: { key: PresetKey; label: string }[] = [
    { key: 'week', label: 'Неделя' },
    { key: 'month', label: 'Месяц' },
    { key: '3months', label: '3 мес.' },
    { key: '6months', label: '6 мес.' },
];

function getPresetRange(key: PresetKey): DateRange {
    const now = new Date();
    const end = endOfDay(now);
    switch (key) {
        case 'week':
            return { start: startOfDay(subtractDays(now, 7)), end };
        case 'month':
            return { start: startOfDay(subtractMonths(now, 1)), end };
        case '3months':
            return { start: startOfDay(subtractMonths(now, 3)), end };
        case '6months':
            return { start: startOfDay(subtractMonths(now, 6)), end };
        default:
            return { start: startOfDay(subtractMonths(now, 1)), end };
    }
}

// Default to realtime
const defaultRange = getPresetRange('week');

export const activePreset = writable<PresetKey>('week');
export const dateRange = writable<DateRange>(defaultRange);
export const isRealtime = writable<boolean>(true);

export const formattedRange = derived(dateRange, ($range) => {
    const opts: Intl.DateTimeFormatOptions = {
        day: 'numeric',
        month: 'short',
    };
    const startStr = $range.start.toLocaleDateString('ru-RU', opts);
    const endStr = $range.end.toLocaleDateString('ru-RU', opts);
    return `${startStr} — ${endStr}`;
});

export function setPreset(key: PresetKey) {
    activePreset.set(key);
    dateRange.set(getPresetRange(key));
    isRealtime.set(false);
}

export function setCustomRange(start: Date, end: Date) {
    activePreset.set('custom');
    dateRange.set({ start: startOfDay(start), end: endOfDay(end) });
    isRealtime.set(false);
}

export function activateRealtime() {
    isRealtime.set(true);
}

export function formatDateInput(date: Date): string {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, '0');
    const d = String(date.getDate()).padStart(2, '0');
    return `${y}-${m}-${d}`;
}
