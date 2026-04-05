<script module lang="ts">
    import { edit } from '@/routes/security';

    export const layout = {
        breadcrumbs: [
            {
                title: 'Настройки безопасности',
                href: edit(),
            },
        ],
    };
</script>

<script lang="ts">
    import { Form } from '@inertiajs/svelte';
    import ShieldCheck from 'lucide-svelte/icons/shield-check';
    import { onDestroy } from 'svelte';
    import SecurityController from '@/actions/App/Http/Controllers/Settings/SecurityController';
    import AppHead from '@/components/AppHead.svelte';
    import Heading from '@/components/Heading.svelte';
    import InputError from '@/components/InputError.svelte';
    import PasswordInput from '@/components/PasswordInput.svelte';
    import TwoFactorRecoveryCodes from '@/components/TwoFactorRecoveryCodes.svelte';
    import TwoFactorSetupModal from '@/components/TwoFactorSetupModal.svelte';
    import { Button } from '@/components/ui/button';
    import { Label } from '@/components/ui/label';
    import { twoFactorAuthState } from '@/lib/twoFactorAuth.svelte';
    import { disable, enable } from '@/routes/two-factor';

    let {
        canManageTwoFactor = false,
        requiresConfirmation = false,
        twoFactorEnabled = false,
    }: {
        canManageTwoFactor?: boolean;
        requiresConfirmation?: boolean;
        twoFactorEnabled?: boolean;
    } = $props();

    const twoFactorAuth = twoFactorAuthState();
    let showSetupModal = $state(false);

    onDestroy(() => twoFactorAuth.clearTwoFactorAuthData());
</script>

<AppHead title="Настройки безопасности" />

<h1 class="sr-only">Настройки безопасности</h1>

<div class="space-y-6">
    <Heading
        variant="small"
        title="Обновление пароля"
        description="Убедитесь, что ваш аккаунт использует длинный и сложный пароль для безопасности"
    />

    <Form
        {...SecurityController.update.form()}
        class="space-y-6"
        options={{ preserveScroll: true }}
        resetOnSuccess
        resetOnError={['password', 'password_confirmation', 'current_password']}
    >
        {#snippet children({ errors, processing, recentlySuccessful })}
            <div class="grid gap-2">
                <Label for="current_password">Текущий пароль</Label>
                <PasswordInput
                    id="current_password"
                    name="current_password"
                    class="mt-1 block w-full"
                    autocomplete="current-password"
                    placeholder="Текущий пароль"
                />
                <InputError message={errors.current_password} />
            </div>

            <div class="grid gap-2">
                <Label for="password">Новый пароль</Label>
                <PasswordInput
                    id="password"
                    name="password"
                    class="mt-1 block w-full"
                    autocomplete="new-password"
                    placeholder="Новый пароль"
                />
                <InputError message={errors.password} />
            </div>

            <div class="grid gap-2">
                <Label for="password_confirmation">Подтвердите пароль</Label>
                <PasswordInput
                    id="password_confirmation"
                    name="password_confirmation"
                    class="mt-1 block w-full"
                    autocomplete="new-password"
                    placeholder="Подтвердите пароль"
                />
                <InputError message={errors.password_confirmation} />
            </div>

            <div class="flex items-center gap-4">
                <Button
                    type="submit"
                    disabled={processing}
                    data-test="update-password-button"
                >
                    Сохранить пароль
                </Button>

                {#if recentlySuccessful}
                    <p class="text-sm text-neutral-600">Сохранено.</p>
                {/if}
            </div>
        {/snippet}
    </Form>
</div>

{#if canManageTwoFactor}
    <div class="space-y-6">
        <Heading
            variant="small"
            title="Двухфакторная аутентификация"
            description="Управление настройками двухфакторной аутентификации"
        />

        {#if !twoFactorEnabled}
            <div class="flex flex-col items-start justify-start space-y-4">
                <p class="text-muted-foreground text-sm">
                    Когда вы включите двухфакторную аутентификацию, во время входа у вас будет запрашиваться пин-код безопасности. Этот код можно получить в приложении-аутентификаторе (TOTP) на вашем мобильном телефоне.
                </p>

                <div>
                    {#if twoFactorAuth.hasSetupData()}
                        <Button onclick={() => (showSetupModal = true)}>
                            <ShieldCheck class="size-4" />Продолжить настройку
                        </Button>
                    {:else}
                        <Form
                            {...enable.form()}
                            onSuccess={() => (showSetupModal = true)}
                        >
                            {#snippet children({ processing })}
                                <Button type="submit" disabled={processing}>
                                    Включить 2FA
                                </Button>
                            {/snippet}
                        </Form>
                    {/if}
                </div>
            </div>
        {:else}
            <div class="flex flex-col items-start justify-start space-y-4">
                <p class="text-muted-foreground text-sm">
                    При входе в систему вы будете вводить случайный пин-код безопасности, который сможете получить в приложении-аутентификаторе (TOTP) на вашем телефоне.
                </p>

                <div class="relative inline">
                    <Form {...disable.form()}>
                        {#snippet children({ processing })}
                            <Button
                                variant="destructive"
                                type="submit"
                                disabled={processing}
                            >
                                Отключить 2FA
                            </Button>
                        {/snippet}
                    </Form>
                </div>

                <TwoFactorRecoveryCodes />
            </div>
        {/if}

        <TwoFactorSetupModal
            bind:isOpen={showSetupModal}
            {requiresConfirmation}
            {twoFactorEnabled}
        />
    </div>
{/if}
