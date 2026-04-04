<script lang="ts">
    import { isUserSettingsOpen, closeUserSettings } from '@/lib/uiStore';
    import ModalContainer from './ModalContainer.svelte';

    // Simple state for ui presets
    let darkMode = true;
    let highContrast = false;
    let alertsEnabled = true;

    // Simulated password logic for SPA look-and-feel
    let currentPassword = '';
    let newPassword = '';
    let confirmPassword = '';
    let passwordMessage = '';

    function toggleDarkMode() {
        darkMode = !darkMode;

        if (darkMode) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }

    function savePassword() {
        if (!currentPassword || !newPassword || !confirmPassword) {
            passwordMessage = 'Заполните все поля';

            return;
        }

        if (newPassword !== confirmPassword) {
            passwordMessage = 'Новые пароли не совпадают';

            return;
        }

        // Simulated API call success
        passwordMessage = 'Пароль успешно изменен!';
        setTimeout(() => {
            currentPassword = '';
            newPassword = '';
            confirmPassword = '';
            passwordMessage = '';
        }, 3000);
    }
</script>

<ModalContainer
    isOpen={$isUserSettingsOpen}
    title="Настройки пользователя"
    on:close={closeUserSettings}
>
    <div class="space-y-6">
        <div>
            <h3 class="font-medium mb-3 opacity-80">Внешний вид (Пресеты)</h3>
            <div class="space-y-3">
                <label class="flex items-center justify-between p-3 rounded-lg border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/50 cursor-pointer">
                    <span class="text-sm">Тёмная тема</span>
                    <input type="checkbox" checked={darkMode} on:change={toggleDarkMode} class="rounded text-emerald-500 focus:ring-emerald-500" />
                </label>
                <label class="flex items-center justify-between p-3 rounded-lg border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/50 cursor-pointer">
                    <span class="text-sm">Высокая контрастность</span>
                    <input type="checkbox" bind:checked={highContrast} class="rounded text-emerald-500 focus:ring-emerald-500" />
                </label>
                <label class="flex items-center justify-between p-3 rounded-lg border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/50 cursor-pointer">
                    <span class="text-sm">Звуковые оповещения (критические)</span>
                    <input type="checkbox" bind:checked={alertsEnabled} class="rounded text-emerald-500 focus:ring-emerald-500" />
                </label>
            </div>
        </div>

        <hr class="border-zinc-200 dark:border-zinc-800" />

        <div>
            <h3 class="font-medium mb-3 opacity-80">Безопасность</h3>
            <div class="space-y-3">
                <input
                    type="password"
                    bind:value={currentPassword}
                    placeholder="Текущий пароль"
                    class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800 border-none rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none text-sm"
                />
                <input
                    type="password"
                    bind:value={newPassword}
                    placeholder="Новый пароль"
                    class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800 border-none rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none text-sm"
                />
                <input
                    type="password"
                    bind:value={confirmPassword}
                    placeholder="Подтверждение нового пароля"
                    class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800 border-none rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none text-sm"
                />
                <button
                    on:click={savePassword}
                    class="w-full py-2 bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 rounded-lg text-sm font-medium transition-colors hover:opacity-90"
                >
                    Сменить пароль
                </button>

                {#if passwordMessage}
                    <div class={`text-sm text-center ${passwordMessage.includes('успешно') ? 'text-emerald-500' : 'text-red-500'}`}>
                        {passwordMessage}
                    </div>
                {/if}
            </div>
        </div>
    </div>
</ModalContainer>
