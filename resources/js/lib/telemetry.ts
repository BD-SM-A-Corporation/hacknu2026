import { writable, get } from 'svelte/store';

export interface TelemetryData {
    locomotiveId: string;
    speed: number; // km/h
    temperature: number; // Celsius
    pressure: number; // atm or kPa
    fuelLevel: number; // percentage (0-100)
    timestamp: string; // ISO string
}

// Global telemetry state per locomotive ID, or just currently selected locomotive
export const telemetryData = writable<TelemetryData>({
    locomotiveId: '',
    speed: 0,
    temperature: 0,
    pressure: 0,
    fuelLevel: 0,
    timestamp: new Date().toISOString(),
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
    if (data.locomotiveId && data.locomotiveId !== get(activeLocomotiveId)) {
        // Skip updates for non-active locomotives if we are only focusing on one,
        // or we could maintain a dictionary of { [id]: data }. For now, we update the active one.
        if (get(activeLocomotiveId) === null) {
            // First time receiving data or no active selected
            activeLocomotiveId.set(data.locomotiveId);
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
