<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { writable } from 'svelte/store';
    import { activeLocomotiveId } from '@/lib/telemetry';
    import type { TelemetryData } from '@/lib/telemetry';
    import { openLocomotiveSelector } from '@/lib/uiStore';

    let mapElement: HTMLElement;
    let map: any;
    // Store markers
    let markers: Record<string, any> = {};
    let polylines: Record<string, any> = {};
    let modelsData: Record<string, string> = {};
    let L: any;

    const isFullscreen = writable(false);

    interface MapLocomotive {
        id: string;
        series: string;
        lat: number;
        lng: number;
        speed: number;
        status: string;
    }

    function toggleFullscreen() {
        isFullscreen.update(v => !v);
        setTimeout(() => {
            if (map) map.invalidateSize();
        }, 300); // 300ms accounts for transition delay
    }

    function getStatusColor(status: string) {
        if (status === 'warning' || status === 'Critical' || status === 'anomaly') return '#ef4444'; // Red
        if (status === 'stopped' || status === 'Offline' || status === 'depot') return '#71717a'; // Gray
        return '#10b981'; // Green
    }

    function createTrainIcon(color: string) {
        const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="${color}" width="24" height="24" stroke="white" stroke-width="1.5"><path d="M4 15s1-1 4-1 5 2 8 2 4-1 4-1V3s-1 1-4 1-5-2-8-2-4 1-4 1z"/><line x1="4" y1="22" x2="4" y2="15"/></svg>`;

        return L.divIcon({
            html: `<div class="relative w-10 h-10 rounded-full flex items-center justify-center bg-white dark:bg-zinc-800 shadow-lg border-2 z-[400]" style="border-color: ${color};">${svg}</div>`,
            className: 'custom-train-marker',
            iconSize: [40, 40],
            iconAnchor: [20, 20],
            popupAnchor: [0, -20]
        });
    }

    function bindPopup(marker: any, id: string, series: string, speed: number) {
        marker.bindPopup(`
            <div class="p-1 min-w-[200px]">
                <div class="font-bold text-base mb-1">${id} <span class="text-xs text-zinc-500">(${series})</span></div>
                <div class="text-sm mb-3 text-zinc-600">Скорость: <b class="text-zinc-900">${Math.round(speed)} км/ч</b></div>
                <div class="flex flex-col gap-2">
                    <button class="w-full text-xs font-semibold bg-blue-600 text-white px-3 py-1.5 rounded shadow-sm hover:bg-blue-700 transition" onclick="window.dispatchZoomRoute('${id}')">🔍 Показать маршрут</button>
                    <button class="w-full text-xs font-semibold bg-emerald-600 text-white px-3 py-1.5 rounded shadow-sm hover:bg-emerald-700 transition" onclick="window.dispatchOpenModal('${id}')">📈 Перейти к деталям</button>
                </div>
            </div>
        `);
    }

    onMount(async () => {
        // Safe dynamic import preventing SSR breakage
        const leaflet = await import('leaflet');
        L = leaflet.default || leaflet;
        await import('leaflet/dist/leaflet.css');

        map = L.map(mapElement, {
            attributionControl: false
        }).setView([48.0196, 66.9237], 5);

        L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
            attribution: '© OpenStreetMap contributors © CARTO',
            subdomains: 'abcd',
            maxZoom: 20
        }).addTo(map);

        const res = await fetch('/api/map/positions');
        const locos: MapLocomotive[] = await res.json();

        locos.forEach(async loco => {
            modelsData[loco.id] = loco.series || 'Тепловоз';

            // Draw historical path
            try {
                const histRes = await fetch(`/api/locomotives/${loco.id}/stats`);

                if (histRes.ok) {
                    const stats = await histRes.json();
                    const color = getStatusColor(loco.status);

                    const latlngs = stats
                        .filter((s: any) => s.lat && s.lng)
                        .reverse() // reverse to draw from oldest to newest
                        .map((s: any) => [parseFloat(s.lat), parseFloat(s.lng)]);

                    if (latlngs.length > 0) {
                        const path = L.polyline(latlngs, {
                            color: color,
                            weight: 3,
                            opacity: 0.6,
                            dashArray: '5, 8'
                        }).addTo(map);
                        polylines[loco.id] = path;

                        // Set map view to show all paths optionally here if wanted
                    }
                }
            } catch (e) {
                console.error("Failed to map path", e);
            }

            if (loco.lat && loco.lng) {
                const color = getStatusColor(loco.status);
                const marker = L.marker([loco.lat, loco.lng], { icon: createTrainIcon(color) }).addTo(map);
                bindPopup(marker, loco.id, modelsData[loco.id], loco.speed);
                markers[loco.id] = marker;
            }
        });

        window.addEventListener('locomotivePing', handlePing);
    });

    onDestroy(() => {
        if (typeof window !== 'undefined') {
            window.removeEventListener('locomotivePing', handlePing);
        }
        if (map) {
            map.remove();
        }
    });

    if (typeof window !== 'undefined') {
        (window as any).dispatchOpenModal = (id: string) => {
            activeLocomotiveId.set(id);
            openLocomotiveSelector();
        };
        (window as any).dispatchZoomRoute = (id: string) => {
            if (polylines[id] && map) {
                const bounds = polylines[id].getBounds();
                if (bounds.isValid()) {
                    map.fitBounds(bounds, { padding: [50, 50], maxZoom: 14 });
                } else if (markers[id]) {
                    map.setView(markers[id].getLatLng(), 14);
                }
            } else if (markers[id] && map) {
                map.setView(markers[id].getLatLng(), 14);
            }
        };
    }

    function handlePing(e: Event) {
        const customEvent = e as CustomEvent<TelemetryData>;
        const data = customEvent.detail;

        let id = data.locomotiveId || (data as any).locomotive_id;
        if (!id || (data as any).gps_corrupted || !L) return;

        let lat = (data as any).lat;
        let lng = (data as any).lng;

        if (lat && lng) {
            let newStatus = 'active';
            if ((data as any).is_anomaly || data.temperature > 100 || data.pressure < 4.5 || (!id.includes('KZ8A') && data.fuelLevel < 10)) {
                newStatus = 'warning';
            } else if (data.speed < 1) {
                newStatus = 'stopped';
            }

            const color = getStatusColor(newStatus);
            const series = modelsData[id] || 'Поезд';

            // Build polyline dynamically
            if (polylines[id]) {
                polylines[id].addLatLng([lat, lng]);
                polylines[id].setStyle({ color: color });
            } else if (map) {
                const path = L.polyline([[lat, lng]], {
                    color: color,
                    weight: 3,
                    opacity: 0.6,
                    dashArray: '5, 8'
                }).addTo(map);
                polylines[id] = path;
            }

            if (markers[id]) {
                const marker = markers[id];
                marker.setLatLng([lat, lng]);
                marker.setIcon(createTrainIcon(color));

                const isOpen = marker.getPopup()?.isOpen();
                bindPopup(marker, id, series, data.speed);
                if (isOpen) marker.openPopup();
            } else {
                const marker = L.marker([lat, lng], { icon: createTrainIcon(color) }).addTo(map);
                bindPopup(marker, id, series, data.speed);
                markers[id] = marker;
            }
        }
    }
</script>

<div class="{$isFullscreen ? 'fixed inset-0 z-[9999] bg-zinc-100 dark:bg-zinc-950 p-6' : 'mb-6'} w-full transition-all duration-300 relative group">
    <div class="{$isFullscreen ? 'h-full rounded-2xl' : 'h-[450px] rounded-xl'} w-full overflow-hidden border border-zinc-200 dark:border-zinc-800 shadow-sm relative z-0">
        <div bind:this={mapElement} class="w-full h-full"></div>
        <div class="absolute top-4 right-4 bg-white/95 dark:bg-zinc-900/95 backdrop-blur px-3 py-2 rounded-lg shadow-sm border border-zinc-200 dark:border-zinc-800 pointer-events-none z-[400] text-sm flex items-center shadow-lg font-medium">
            <span class="inline-block w-2.5 h-2.5 rounded-full bg-emerald-500 mr-2 animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.8)]"></span>
            GPS CARTO Active
        </div>

        <!-- Fullscreen Button -->
        <button onclick={toggleFullscreen} class="absolute bottom-4 right-4 z-[400] bg-white dark:bg-zinc-800 p-2.5 rounded-lg shadow-lg border border-zinc-200 dark:border-zinc-700 hover:bg-zinc-50 dark:hover:bg-zinc-700 transition" title="Полный экран">
            {#if $isFullscreen}
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8 3v3a2 2 0 0 1-2 2H3"/><path d="M21 8h-3a2 2 0 0 1-2-2V3"/><path d="M3 16h3a2 2 0 0 1 2 2v3"/><path d="M16 21v-3a2 2 0 0 1 2-2h3"/></svg>
            {:else}
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8 3H5a2 2 0 0 0-2 2v3"/><path d="M21 8V5a2 2 0 0 0-2-2h-3"/><path d="M3 16v3a2 2 0 0 0 2 2h3"/><path d="M16 21h3a2 2 0 0 0 2-2v-3"/></svg>
            {/if}
        </button>
    </div>
</div>

<style>
    :global(.leaflet-container) {
        z-index: 1 !important;
        font-family: inherit;
    }
    :global(.custom-train-marker) {
        background: transparent;
        border: none;
    }
</style>
