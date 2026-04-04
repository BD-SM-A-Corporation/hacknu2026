<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class LocomotiveResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'model' => $this->model,
            'depot' => $this->depot,
            'status' => 'active', // mock or compute from telemetry later
            'created_at' => $this->created_at,
        ];
    }
}
