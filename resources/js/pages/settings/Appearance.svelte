<script module lang="ts">
    import { edit as editAppearance } from '@/routes/appearance';

    export const layout = {
        breadcrumbs: [
            {
                title: 'Настройки внешнего вида',
                href: editAppearance(),
            },
        ],
    };
</script>

<script lang="ts">
    import AppearanceTabs from '@/components/AppearanceTabs.svelte';
    import AppHead from '@/components/AppHead.svelte';
    import Heading from '@/components/Heading.svelte';
</script>

<AppHead title="Настройки внешнего вида" />

<h1 class="sr-only">Настройки внешнего вида</h1>

<div class="space-y-6">
    <Heading
        variant="small"
        title="Внешний вид"
        description="Настройте тему и оформление вашего аккаунта"
    />
    <AppearanceTabs />
</div>
