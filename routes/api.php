<?php

use App\Http\Controllers\Api\EndpointController;
use App\Http\Controllers\Api\LocomotiveController;
use App\Http\Controllers\Api\TelemetryAnalyticsController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

Route::middleware('auth:sanctum')->group(function () {
    Route::get('/user', function (Request $request) {
        return $request->user();
    });

    Route::get('locomotives/{locomotive}/stats', [LocomotiveController::class, 'stats']);
    Route::get('map/positions', [LocomotiveController::class, 'positions']);
    Route::apiResource('locomotives', LocomotiveController::class);
    Route::apiResource('endpoints', EndpointController::class);
    Route::get('telemetry/analytics', [TelemetryAnalyticsController::class, 'index']);
});
