import { writable, derived, get } from 'svelte/store';

export interface TelemetryData {
    locomotiveId: string;
    speed: number; // km/h
    temperature: number; // Celsius
    pressure: number; // atm or kPa
    fuelLevel: number; // percentage (0-100)
    healthScore: number; // 0-100 score
    alerts: string[]; // List of real-time anomaly alerts
    timestamp: string; // ISO string
}

// Global telemetry state per locomotive ID, or just currently selected locomotive
export const telemetryData = writable<TelemetryData>({
    locomotiveId: '',
    speed: 0,
    temperature: 0,
    pressure: 0,
    fuelLevel: 0,
    healthScore: 100,
    alerts: [],
    timestamp: new Date().toISOString(),
});

export const allTelemetry = writable<Record<string, TelemetryData>>({});

export const averageFleetTelemetry = derived(allTelemetry, ($all) => {
    const keys = Object.keys($all);
    const count = keys.length;
    
    if (count === 0) {
        return {
            locomotiveId: 'ALL',
            speed: 0,
            temperature: 0,
            pressure: 0,
            fuelLevel: 0,
            healthScore: 100,
            alerts: [],
            timestamp: new Date().toISOString()
        };
    }

    let speed = 0, temp = 0, press = 0, fuel = 0, health = 0;
    
    keys.forEach(k => {
        speed += Number($all[k].speed) || 0;
        temp += Number($all[k].temperature) || 0;
        press += Number($all[k].pressure) || 0;
        fuel += Number($all[k].fuelLevel) || 0;
        health += Number($all[k].healthScore) || 0;
    });

    return {
        locomotiveId: 'ALL',
        speed: speed / count,
        temperature: temp / count,
        pressure: press / count,
        fuelLevel: fuel / count,
        healthScore: Math.round(health / count),
        alerts: [],
        timestamp: new Date().toISOString()
    };
});

export const activeLocomotiveId = writable<string | null>(
    typeof window !== 'undefined'
        ? localStorage.getItem('activeLocomotiveId')
        : null,
);

if (typeof window !== 'undefined') {
    activeLocomotiveId.subscribe((value) => {
        if (value) {
            localStorage.setItem('activeLocomotiveId', value);
        }
    });
}

// Functions to update the telemetry safely
export function updateTelemetry(data: Partial<TelemetryData>) {
    const id = data.locomotiveId;
    if (!id) return;

    allTelemetry.update((map) => {
        const current = map[id] || {
            locomotiveId: id,
            speed: 0,
            temperature: 0,
            pressure: 0,
            fuelLevel: 0,
            healthScore: 100,
            alerts: [],
            timestamp: new Date().toISOString(),
        };
        map[id] = { ...current, ...data, timestamp: new Date().toISOString() };
        return map;
    });

    if (id !== get(activeLocomotiveId)) {
        // Skip updates for the main focused widget if it's not the active one
        if (get(activeLocomotiveId) === null) {
            activeLocomotiveId.set(id);
        } else {
            return;
        }
    }

    telemetryData.update((current) => ({
        ...current,
        ...data,
        timestamp: new Date().toISOString(),
    }));
}
