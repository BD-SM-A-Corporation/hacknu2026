import { writable } from 'svelte/store';

// UI State management for single-page dynamic modals/tabs Without Routing

export const isLocomotiveSelectorOpen = writable<boolean>(false);
export const isWebSocketConfigOpen = writable<boolean>(false);
export const isUserSettingsOpen = writable<boolean>(false);

export function openLocomotiveSelector() {
    isLocomotiveSelectorOpen.set(true);
}

export function closeLocomotiveSelector() {
    isLocomotiveSelectorOpen.set(false);
}

export function openWebSocketConfig() {
    isWebSocketConfigOpen.set(true);
}

export function closeWebSocketConfig() {
    isWebSocketConfigOpen.set(false);
}

export function openUserSettings() {
    isUserSettingsOpen.set(true);
}

export function closeUserSettings() {
    isUserSettingsOpen.set(false);
}

export function closeAllModals() {
    isLocomotiveSelectorOpen.set(false);
    isWebSocketConfigOpen.set(false);
    isUserSettingsOpen.set(false);
}
