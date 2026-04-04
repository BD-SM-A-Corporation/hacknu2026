<script module lang="ts">
</script>

<script lang="ts">
    import AppHead from '@/components/AppHead.svelte';
    import { setLayoutProps } from '@inertiajs/svelte';

    let { htmlContent, title } = $props<{
        htmlContent: string;
        title: string;
    }>();
    
    let breadcrumbs = $derived([
        { title: 'Документация', href: '/docs' },
        ...(title !== 'Документация' ? [{ title, href: '' }] : [])
    ]);

    $effect(() => {
        setLayoutProps({ breadcrumbs });
    });
</script>

<AppHead {title} />

<div class="px-4 py-8 mx-auto w-full max-w-6xl flex justify-center">
    <div class="w-full max-w-4xl bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-sm p-6 sm:p-10">
        <article class="prose prose-zinc dark:prose-invert max-w-none">
            {@html htmlContent}
        </article>
    </div>
</div>
