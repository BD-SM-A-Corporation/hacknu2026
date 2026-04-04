<?php

namespace App\Console\Commands;

use App\Models\Locomotive;
use App\Models\TelemetryHistory;
use Carbon\Carbon;
use Illuminate\Console\Command;

class CheckLocomotiveStatus extends Command
{
    protected $signature = 'telemetry:check-status';

    protected $description = 'Checks the last ping time of locomotives and marks them Offline if > 5 mins';

    public function handle()
    {
        $locomotives = Locomotive::all();
        $fiveMinsAgo = Carbon::now()->subMinutes(5);

        foreach ($locomotives as $loco) {
            $lastPing = TelemetryHistory::where('locomotive_id', $loco->id)->max('timestamp');

            if (! $lastPing || Carbon::parse($lastPing)->isBefore($fiveMinsAgo)) {
                if ($loco->status !== 'Offline') {
                    $loco->update(['status' => 'Offline']);
                    $this->info("Locomotive {$loco->id} marked as Offline.");
                }
            } else {
                if ($loco->status === 'Offline') {
                    $loco->update(['status' => 'Active']);
                    $this->info("Locomotive {$loco->id} came back Online.");
                }
            }
        }
    }
}
