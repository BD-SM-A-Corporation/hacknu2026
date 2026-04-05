<script module lang="ts">
    import { edit } from '@/routes/profile';

    export const layout = {
        breadcrumbs: [
            {
                title: 'Настройки профиля',
                href: edit(),
            },
        ],
    };
</script>

<script lang="ts">
    import { Form, page } from '@inertiajs/svelte';
    import ProfileController from '@/actions/App/Http/Controllers/Settings/ProfileController';
    import AppHead from '@/components/AppHead.svelte';
    import DeleteUser from '@/components/DeleteUser.svelte';
    import Heading from '@/components/Heading.svelte';
    import InputError from '@/components/InputError.svelte';
    import TextLink from '@/components/TextLink.svelte';
    import { Button } from '@/components/ui/button';
    import { Input } from '@/components/ui/input';
    import { Label } from '@/components/ui/label';
    import { send } from '@/routes/verification';

    let {
        mustVerifyEmail,
        status = '',
    }: {
        mustVerifyEmail: boolean;
        status?: string;
    } = $props();

    const user = $derived(page.props.auth.user);
</script>

<AppHead title="Настройки профиля" />

<h1 class="sr-only">Настройки профиля</h1>

<div class="flex flex-col space-y-6">
    <Heading
        variant="small"
        title="Информация профиля"
        description="Обновите ваше имя и адрес электронной почты"
    />

    <Form
        {...ProfileController.update.form()}
        class="space-y-6"
        options={{ preserveScroll: true }}
    >
        {#snippet children({ errors, processing, recentlySuccessful })}
            <div class="grid gap-2">
                <Label for="name">Имя</Label>
                <Input
                    id="name"
                    name="name"
                    class="mt-1 block w-full"
                    value={user.name}
                    required
                    autocomplete="name"
                    placeholder="Имя"
                />
                <InputError class="mt-2" message={errors.name} />
            </div>

            <div class="grid gap-2">
                <Label for="email">Электронная почта</Label>
                <Input
                    id="email"
                    type="email"
                    name="email"
                    class="mt-1 block w-full"
                    value={user.email}
                    required
                    autocomplete="username"
                    placeholder="Ваш email"
                />
                <InputError class="mt-2" message={errors.email} />
            </div>

            {#if mustVerifyEmail && !user.email_verified_at}
                <div>
                    <p class="-mt-4 text-sm text-muted-foreground">
                        Ваш адрес электронной почты не подтвержден.
                        <TextLink href={send()} as="button">
                            Нажмите здесь, чтобы отправить письмо с подтверждением еще раз.
                        </TextLink>
                    </p>

                    {#if status === 'verification-link-sent'}
                        <div class="mt-2 text-sm font-medium text-green-600">
                            Новая ссылка была отправлена на ваш адрес электронной почты.
                        </div>
                    {/if}
                </div>
            {/if}

            <div class="flex items-center gap-4">
                <Button
                    type="submit"
                    disabled={processing}
                    data-test="update-profile-button">Сохранить</Button
                >

                {#if recentlySuccessful}
                    <p class="text-sm text-neutral-600">Сохранено.</p>
                {/if}
            </div>
        {/snippet}
    </Form>
</div>

<DeleteUser />
