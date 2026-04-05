<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\TelemetryHistory;
use Illuminate\Http\Request;

class TelemetryAnalyticsController extends Controller
{
    /**
     * Display a listing of the resource.
     */
    public function index(Request $request)
    {
        $request->validate([
            'locomotive_id' => 'required|string',
            'start_date' => 'required|date',
            'end_date' => 'required|date|after_or_equal:start_date',
            'step' => 'nullable|integer|min:1',
        ]);

        $query = TelemetryHistory::where('locomotive_id', $request->locomotive_id)
            ->whereBetween('timestamp', [$request->start_date, $request->end_date]);

        if ($request->filled('step')) {
            $step = (int) $request->step;
            $data = $query->selectRaw("
                date_bin('$step seconds', timestamp, TIMESTAMP '2026-01-01') as timestamp,
                AVG(speed) as speed,
                AVG(temperature) as temperature,
                AVG(pressure) as pressure,
                AVG(fuel_level) as fuel_level
            ")
            ->groupByRaw("date_bin('$step seconds', timestamp, TIMESTAMP '2026-01-01')")
            ->orderBy('timestamp', 'asc')
            ->get();
        } else {
            $data = $query->select('timestamp', 'speed', 'temperature', 'pressure', 'fuel_level')
                ->orderBy('timestamp', 'asc')
                ->get();
        }

        return response()->json($data);
    }

    public function allLocomotives(Request $request)
    {
        $request->validate([
            'start_date' => 'required|date',
            'end_date' => 'required|date|after_or_equal:start_date',
            'step' => 'nullable|integer|min:1',
        ]);

        $query = TelemetryHistory::whereBetween('timestamp', [$request->start_date, $request->end_date]);

        $step = $request->filled('step') ? (int) $request->step : 300; // default 5 minutes

        $data = $query->selectRaw("
            date_bin('$step seconds', timestamp, TIMESTAMP '2026-01-01') as timestamp,
            AVG(speed) as speed,
            AVG(temperature) as temperature,
            AVG(pressure) as pressure,
            AVG(fuel_level) as fuel_level
        ")
        ->groupByRaw("date_bin('$step seconds', timestamp, TIMESTAMP '2026-01-01')")
        ->orderBy('timestamp', 'asc')
        ->get();

        return response()->json($data);
    }

    public function anomaliesHistory(Request $request)
    {
        $limit = $request->query('limit', 50);

        $data = TelemetryHistory::where('is_anomaly', true)
            ->select('locomotive_id', 'timestamp', 'speed', 'temperature', 'pressure', 'fuel_level', 'health_score')
            ->orderBy('timestamp', 'desc')
            ->limit($limit)
            ->get();
            
        return response()->json($data);
    }

    /**
     * Store a newly created resource in storage.
     */
    public function store(Request $request)
    {
        //
    }

    /**
     * Display the specified resource.
     */
    public function show(string $id)
    {
        //
    }

    /**
     * Update the specified resource in storage.
     */
    public function update(Request $request, string $id)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     */
    public function destroy(string $id)
    {
        //
    }
}
