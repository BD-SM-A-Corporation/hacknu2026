<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Str;
use Illuminate\Support\Facades\File;
use Inertia\Inertia;

class DocumentationController extends Controller
{
    public function show($version = null, $page = null)
    {
        $path = 'index';
        if ($version && $page) {
            $path = $version . '/' . $page;
        } elseif ($version) {
            $path = $version;
            if (is_dir(base_path('docs/' . $path))) {
                $path .= '/index';
            }
        }

        $file = base_path('docs/' . $path . '.md');

        if (!File::exists($file)) {
            abort(404, 'Раздел документации не найден.');
        }

        $content = File::get($file);
        $htmlContent = Str::markdown($content);

        $title = 'Документация';
        if (preg_match('/^#\s+(.+)$/m', $content, $matches)) {
            $title = tap($matches[1], fn($t) => trim($t));
        }

        return Inertia::render('Docs/Show', [
            'htmlContent' => $htmlContent,
            'title' => $title,
        ]);
    }
}
