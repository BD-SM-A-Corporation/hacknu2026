import { updateTelemetry } from './telemetry';

export class WebSocketClient {
    private ws: WebSocket | null = null;
    private url: string;
    private reconnectAttempts: number = 0;
    private maxReconnectAttempts: number = 5;
    private reconnectTimeout: number = 2000;

    // Callbacks
    public onConnect: () => void = () => {};
    public onDisconnect: () => void = () => {};
    public onError: (err: Event) => void = () => {};

    constructor(url: string) {
        this.url = url;
    }

    public connect() {
        if (
            this.ws &&
            (this.ws.readyState === WebSocket.CONNECTING ||
                this.ws.readyState === WebSocket.OPEN)
        ) {
            return;
        }

        try {
            this.ws = new WebSocket(this.url);

            this.ws.onopen = () => {
                console.log('WebSocket connected to', this.url);
                this.reconnectAttempts = 0;
                this.onConnect();
            };

            this.ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);

                    // Expected format: { type: 'telemetry', data: { locomotiveId, speed, temperature, pressure, fuelLevel } }
                    if (data.type === 'telemetry' && data.data) {
                        updateTelemetry(data.data);
                    }
                } catch (e) {
                    console.error('Failed to parse WebSocket message', e);
                }
            };

            this.ws.onclose = () => {
                console.log('WebSocket disconnected');
                this.onDisconnect();
                this.ws = null;
                this.attemptReconnect();
            };

            this.ws.onerror = (error) => {
                this.onError(error);
            };
        } catch (e) {
            console.error('Failed to initialize WebSocket', e);
            this.attemptReconnect();
        }
    }

    public disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    public setUrl(url: string) {
        this.url = url;
        this.disconnect();
        this.connect();
    }

    private attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(
                `Reconnecting... Attempt ${this.reconnectAttempts} of ${this.maxReconnectAttempts}`,
            );
            setTimeout(
                () => this.connect(),
                this.reconnectTimeout * this.reconnectAttempts,
            );
        } else {
            console.error('Max reconnect attempts reached.');
        }
    }
}

// Export a singleton instance with a default/fallback endpoint for now
export const wsClient = new WebSocketClient('ws://localhost/telemetry/');
