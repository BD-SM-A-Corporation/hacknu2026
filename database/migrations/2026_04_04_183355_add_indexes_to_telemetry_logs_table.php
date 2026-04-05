<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration {
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::table('telemetry_logs', function (Blueprint $table) {
            // Check if index exists; if not, create it.
            if (!Schema::hasIndex('telemetry_logs', 'idx_loco_time')) {
                $table->index(['locomotive_id', 'timestamp'], 'idx_loco_time');
            }

            if (!Schema::hasIndex('telemetry_logs', 'idx_timestamp')) {
                $table->index('timestamp', 'idx_timestamp');
            }
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('telemetry_logs', function (Blueprint $table) {
            $table->dropIndex('idx_loco_time');
            $table->dropIndex('idx_timestamp');
        });
    }
};
