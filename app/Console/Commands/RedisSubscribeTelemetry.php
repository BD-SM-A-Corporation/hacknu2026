<?php

namespace App\Console\Commands;

use App\Events\LocomotiveAlertBroadcast;
use App\Models\Alert;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Redis;

class RedisSubscribeTelemetry extends Command
{
    protected $signature = 'telemetry:listen-redis';

    protected $description = 'Subscribes to locomotive-updates Redis channel and processes anomalies and alerts';

    public function handle()
    {
        ini_set('default_socket_timeout', -1);
        $this->info("Listening to Redis channel 'locomotive-updates'...");

        Redis::subscribe(['locomotive-updates'], function ($message) {
            $data = json_decode($message, true);
            if (!$data) {
                return;
            }

            $isAnomaly = $data['is_anomaly'] ?? false;
            $temp = $data['temperature'] ?? 0;
            $pressure = $data['pressure'] ?? 0;
            $locoId = $data['locomotiveId'] ?? null;
            $lat = $data['lat'] ?? null;
            $lng = $data['lng'] ?? null;
            $speed = $data['speed'] ?? null;
            $gpsCorrupted = $data['gps_corrupted'] ?? false;

            if (!$locoId) {
                return;
            }

            // Auto-register missing locomotives
            $loco = \App\Models\Locomotive::firstOrCreate(
                ['id' => $locoId],
                ['model' => 'Авто-обнаружение', 'status' => 'Active']
            );

            if (!$gpsCorrupted && $lat !== null && $lng !== null) {
                $loco->update([
                    'lat' => $lat,
                    'lng' => $lng,
                    'speed' => $speed,
                ]);
            }

            $alertType = null;
            $alertMsg = null;

            if ($isAnomaly) {
                $alertType = 'anomaly';
                $alertMsg = "Anomaly detected (speed/fuel spike) for locomotive {$locoId}";
            } elseif ($temp > 100) {
                $alertType = 'temperature';
                $alertMsg = "Critical Temperature ({$temp}°C) for locomotive {$locoId}";
            } elseif ($pressure < 4.5) {
                $alertType = 'pressure';
                $alertMsg = "Low Oil Pressure ({$pressure} atm) for locomotive {$locoId}";
            }

            if ($alertType) {
                Alert::create([
                    'locomotive_id' => $locoId,
                    'type' => $alertType,
                    'message' => $alertMsg,
                ]);

                event(new LocomotiveAlertBroadcast($locoId, $alertType, $alertMsg));
                $this->warn("Alert Dispatched: $alertMsg");
            }
        });
    }
}
