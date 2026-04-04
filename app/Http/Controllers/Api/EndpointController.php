<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\WebsocketEndpoint;
use App\Http\Resources\WebsocketEndpointResource;
use Illuminate\Http\Request;

class EndpointController extends Controller
{
    public function index()
    {
        return WebsocketEndpointResource::collection(WebsocketEndpoint::all());
    }

    public function store(Request $request)
    {
        $validated = $request->validate([
            'url' => 'required|url|unique:websocket_endpoints,url',
            'is_active' => 'boolean',
            'description' => 'nullable|string'
        ]);

        $endpoint = WebsocketEndpoint::create($validated);
        return new WebsocketEndpointResource($endpoint);
    }

    public function show(string $id)
    {
        $endpoint = WebsocketEndpoint::findOrFail($id);
        return new WebsocketEndpointResource($endpoint);
    }

    public function update(Request $request, string $id)
    {
        $endpoint = WebsocketEndpoint::findOrFail($id);
        $validated = $request->validate([
            'url' => 'sometimes|url|unique:websocket_endpoints,url,' . $id,
            'is_active' => 'boolean',
            'description' => 'nullable|string'
        ]);

        $endpoint->update($validated);
        return new WebsocketEndpointResource($endpoint);
    }

    public function destroy(string $id)
    {
        $endpoint = WebsocketEndpoint::findOrFail($id);
        $endpoint->delete();
        return response()->noContent();
    }
}
