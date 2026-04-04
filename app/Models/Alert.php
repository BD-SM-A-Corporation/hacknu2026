<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Alert extends Model
{
    protected $fillable = [
        'locomotive_id',
        'type',
        'message',
        'resolved_at',
    ];
}
