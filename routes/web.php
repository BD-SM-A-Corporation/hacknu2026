<?php

use Illuminate\Support\Facades\Route;
use Laravel\Fortify\Features;

Route::redirect('/', '/dashboard')->name('home');

Route::middleware(['auth', 'verified'])->group(function () {
    Route::inertia('dashboard', 'Dashboard', [
        'canRegister' => Features::enabled(Features::registration()),
    ])->name('dashboard');
});

require __DIR__ . '/settings.php';
