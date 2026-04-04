<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Locomotive;
use App\Http\Resources\LocomotiveResource;
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
}
