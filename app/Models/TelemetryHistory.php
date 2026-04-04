<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class TelemetryHistory extends Model
{
    // The table is managed by the Golang microservice / GORM. Please don't touch
    protected $table = 'telemetry_logs';

    // Prevent Laravel from touching timestamps if they are managed by GORM
    public $timestamps = false; // GORM typically uses created_at, updated_at directly, but we only have `timestamp` here based on Go struct. Wait! If GORM adds them, we can let it be. But we won't write from Laravel anyway.

    protected $guarded = []; // Make read-only by convention

    protected $casts = [
        'timestamp' => 'datetime',
    ];
}
