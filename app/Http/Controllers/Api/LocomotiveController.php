<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Http\Resources\LocomotiveResource;
use App\Models\Locomotive;
use App\Models\TelemetryHistory;
use Illuminate\Http\Request;

class LocomotiveController extends Controller
{
    /**
     * Display a listing of the resource.
     */
    public function index()
    {
        return LocomotiveResource::collection(Locomotive::all());
    }

    /**
     * Store a newly created resource in storage.
     */
    public function store(Request $request)
    {
        $validated = $request->validate([
            'id' => 'required|string|unique:locomotives,id',
            'model' => 'required|string|max:255',
            'depot' => 'nullable|string|max:255',
        ]);

        $locomotive = Locomotive::create($validated);

        return new LocomotiveResource($locomotive);
    }

    /**
     * Display the specified resource.
     */
    public function show(string $id)
    {
        $locomotive = Locomotive::findOrFail($id);

        return new LocomotiveResource($locomotive);
    }

    /**
     * Update the specified resource in storage.
     */
    public function update(Request $request, string $id)
    {
        $locomotive = Locomotive::findOrFail($id);
        $validated = $request->validate([
            'model' => 'sometimes|string|max:255',
            'depot' => 'nullable|string|max:255',
        ]);

        $locomotive->update($validated);

        return new LocomotiveResource($locomotive);
    }

    /**
     * Remove the specified resource from storage.
     */
    public function destroy(string $id)
    {
        $locomotive = Locomotive::findOrFail($id);
        $locomotive->delete();

        return response()->noContent();
    }

    /**
     * Display the latest 50 telemetry stats.
     */
    public function stats(string $id)
    {
        $stats = TelemetryHistory::where('locomotive_id', $id)
            ->orderBy('timestamp', 'desc')
            ->take(50)
            ->get();

        return response()->json($stats);
    }

    /**
     * Get real-time positions
     */
    public function positions()
    {
        $locomotives = Locomotive::select('id', 'model as series', 'lat', 'lng', 'speed', 'status')->get();
        return response()->json($locomotives);
    }
}
