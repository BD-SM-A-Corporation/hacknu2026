<?php

use Illuminate\Support\Facades\Route;
use Laravel\Fortify\Features;

Route::redirect('/', '/dashboard')->name('home');

Route::middleware(['auth', 'verified'])->group(function () {
    Route::inertia('dashboard', 'Dashboard', [
        'canRegister' => Features::enabled(Features::registration()),
    ])->name('dashboard');

    Route::inertia('analytics', 'Analytics')->name('analytics');

    Route::get('docs/{version?}/{page?}', [\App\Http\Controllers\DocumentationController::class, 'show'])
        ->where('page', '.*')
        ->name('docs.show');
});

require __DIR__.'/settings.php';
