<script lang="ts">
    import { Link, page } from '@inertiajs/svelte';
    import BookOpen from 'lucide-svelte/icons/book-open';
    import Folder from 'lucide-svelte/icons/folder';
    import LayoutGrid from 'lucide-svelte/icons/layout-grid';
    import Menu from 'lucide-svelte/icons/menu';
    import Settings from 'lucide-svelte/icons/settings';
    import Activity from 'lucide-svelte/icons/activity';
    import {
        Avatar,
        AvatarFallback,
        AvatarImage,
    } from '@/components/ui/avatar';
    import { Button } from '@/components/ui/button';
    import {
        DropdownMenu,
        DropdownMenuContent,
        DropdownMenuTrigger,
    } from '@/components/ui/dropdown-menu';
    import {
        Sheet,
        SheetContent,
        SheetHeader,
        SheetTitle,
        SheetTrigger,
    } from '@/components/ui/sheet';
    import {
        Tooltip,
        TooltipContent,
        TooltipProvider,
        TooltipTrigger,
    } from '@/components/ui/tooltip';
    import UserMenuContent from '@/components/UserMenuContent.svelte';
    import { currentUrlState } from '@/lib/currentUrl.svelte';
    import { getInitials } from '@/lib/initials';
    import { toUrl } from '@/lib/utils';
    import { dashboard } from '@/routes';
    import { edit } from '@/routes/profile';
    import type { BreadcrumbItem } from '@/types';

    let {
        breadcrumbs = [],
    }: {
        breadcrumbs?: BreadcrumbItem[];
    } = $props();

    const auth = $derived(page.props.auth);
    const url = currentUrlState();

    const activeItemStyles = 'text-[#00AAD9] dark:text-[#00AAD9]';
    let mobileMenuOpen = $state(false);
</script>

<header
    class="sticky top-0 z-40 w-full border-b border-zinc-200/80 bg-white/80 backdrop-blur-xl dark:border-zinc-800/80 dark:bg-zinc-950/80"
>
    <div class="flex h-14 items-center px-4 lg:px-6">
        <!-- ========== LEFT SECTION ========== -->
        <div class="flex items-center gap-1 lg:gap-3">
            <!-- Mobile Menu -->
            <div class="lg:hidden">
                <Sheet bind:open={mobileMenuOpen}>
                    <SheetTrigger asChild>
                        {#snippet children(props)}
                            <Button
                                variant="ghost"
                                size="icon"
                                class="mr-1 h-9 w-9"
                                onclick={props.onclick}
                                aria-expanded={props['aria-expanded']}
                            >
                                <Menu class="h-5 w-5" />
                            </Button>
                        {/snippet}
                    </SheetTrigger>
                    <SheetContent side="left" class="w-[300px] p-6 z-[100]">
                        <SheetTitle class="sr-only">Навигация</SheetTitle>
                        <SheetHeader class="flex justify-start text-left">
                            <img
                                src="/KTZH-32x32.svg"
                                alt="КТЖ"
                                class="size-7 shrink-0"
                            />
                        </SheetHeader>
                        <div
                            class="flex h-full flex-1 flex-col space-y-6 pt-6 pb-10"
                        >
                            <nav class="-mx-3 space-y-1">
                                <Link
                                    href={toUrl(dashboard())}
                                    onclick={() => mobileMenuOpen = false}
                                    class="flex items-center gap-x-3 rounded-lg px-3 py-2 text-sm font-medium hover:bg-accent {url.isCurrentUrl(
                                        dashboard(),
                                        url.currentUrl,
                                    )
                                        ? activeItemStyles
                                        : ''}"
                                >
                                    <LayoutGrid class="h-5 w-5" />
                                    Дэшборд
                                </Link>
                                <Link
                                    href="/analytics"
                                    onclick={() => mobileMenuOpen = false}
                                    class="flex items-center gap-x-3 rounded-lg px-3 py-2 text-sm font-medium hover:bg-accent {url.isCurrentUrl(
                                        '/analytics',
                                        url.currentUrl,
                                    )
                                        ? activeItemStyles
                                        : ''}"
                                >
                                    <Activity class="h-5 w-5" />
                                    Аналитика
                                </Link>
                            </nav>
                            <div
                                class="flex flex-col space-y-3 border-t border-zinc-200 dark:border-zinc-800 pt-4"
                            >
                                <Link
                                    href={toUrl(edit())}
                                    onclick={() => mobileMenuOpen = false}
                                    class="flex items-center space-x-2 text-sm font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100"
                                >
                                    <Settings class="h-5 w-5" />
                                    <span>Настройки</span>
                                </Link>
                                <a
                                    href="https://github.com/BD-SM-A-Corporation/hacknu2026.git"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    onclick={() => mobileMenuOpen = false}
                                    class="flex items-center space-x-2 text-sm font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100"
                                >
                                    <Folder class="h-5 w-5" />
                                    <span>Репозиторий</span>
                                </a>
                                <Link
                                    href="/docs"
                                    onclick={() => mobileMenuOpen = false}
                                    class="flex items-center space-x-2 text-sm font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100"
                                >
                                    <BookOpen class="h-5 w-5" />
                                    <span>Документация</span>
                                </Link>
                            </div>
                        </div>
                    </SheetContent>
                </Sheet>
            </div>

            <!-- Logo -->
            <Link
                href={toUrl(dashboard())}
                class="flex items-center gap-2 transition-opacity hover:opacity-80"
            >
                <img src="/KTZH-32x32.svg" alt="КТЖ" />
                <span
                    class="hidden text-sm font-bold tracking-tight text-zinc-900 dark:text-zinc-100 md:inline lg:hidden"
                >
                    КТЖ
                </span>
                <span
                    class="hidden text-sm font-bold tracking-tight text-zinc-900 dark:text-zinc-100 lg:inline"
                >
                    КТЖ Монитор
                </span>
            </Link>

            <!-- Separator -->
            <div
                class="mx-2 hidden h-5 w-px bg-zinc-200 dark:bg-zinc-700 lg:block"
            ></div>

            <!-- Dashboard Button -->
            <Link
                href={toUrl(dashboard())}
                class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium transition-colors lg:flex
                    {url.isCurrentUrl(dashboard(), url.currentUrl)
                    ? 'bg-[#00AAD9]/10 text-[#00AAD9] dark:bg-[#00AAD9]/20 dark:text-[#00AAD9]'
                    : 'text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100'}"
            >
                <LayoutGrid class="h-4 w-4" />
                Дэшборд
            </Link>

            <!-- Separator -->
            <div
                class="mx-1 hidden h-5 w-px bg-zinc-200 dark:bg-zinc-700 lg:block"
            ></div>

            <!-- Analytics Button -->
            <Link
                href="/analytics"
                class="hidden items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium transition-colors lg:flex
                    {url.currentUrl.startsWith('/analytics')
                    ? 'bg-[#00AAD9]/10 text-[#00AAD9] dark:bg-[#00AAD9]/20 dark:text-[#00AAD9]'
                    : 'text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100'}"
            >
                <Activity class="h-4 w-4" />
                Аналитика
            </Link>
        </div>

        <!-- ========== RIGHT SECTION ========== -->
        <div class="ml-auto flex items-center gap-1 lg:gap-2">
            <!-- Settings -->
            <div class="hidden lg:flex">
                <TooltipProvider delayDuration={0}>
                    <Tooltip>
                        <TooltipTrigger>
                            {#snippet child({ props })}
                                <Link
                                    href={toUrl(edit())}
                                    {...props}
                                    class="inline-flex h-9 w-9 cursor-pointer items-center justify-center rounded-lg text-sm font-medium transition-colors hover:bg-zinc-100 dark:hover:bg-zinc-800 group"
                                >
                                    <span class="sr-only">Настройки</span>
                                    <Settings
                                        class="size-[18px] text-zinc-500 group-hover:text-zinc-700 dark:text-zinc-400 dark:group-hover:text-zinc-200 transition-colors"
                                    />
                                </Link>
                            {/snippet}
                        </TooltipTrigger>
                        <TooltipContent>
                            <p>Настройки</p>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            </div>

            <!-- Repository -->
            <div class="hidden lg:flex">
                <TooltipProvider delayDuration={0}>
                    <Tooltip>
                        <TooltipTrigger>
                            {#snippet child({ props })}
                                <a
                                    href="https://github.com/laravel/svelte-starter-kit"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    {...props}
                                    class="inline-flex h-9 w-9 cursor-pointer items-center justify-center rounded-lg text-sm font-medium transition-colors hover:bg-zinc-100 dark:hover:bg-zinc-800 group"
                                >
                                    <span class="sr-only">Репозиторий</span>
                                    <Folder
                                        class="size-[18px] text-zinc-500 group-hover:text-zinc-700 dark:text-zinc-400 dark:group-hover:text-zinc-200 transition-colors"
                                    />
                                </a>
                            {/snippet}
                        </TooltipTrigger>
                        <TooltipContent>
                            <p>Репозиторий</p>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            </div>

            <!-- Documentation -->
            <div class="hidden lg:flex">
                <TooltipProvider delayDuration={0}>
                    <Tooltip>
                        <TooltipTrigger>
                            {#snippet child({ props })}
                                <Link
                                    href="/docs"
                                    {...props}
                                    class="inline-flex h-9 w-9 cursor-pointer items-center justify-center rounded-lg text-sm font-medium transition-colors hover:bg-zinc-100 dark:hover:bg-zinc-800 group"
                                >
                                    <span class="sr-only">Документация</span>
                                    <BookOpen
                                        class="size-[18px] text-zinc-500 group-hover:text-zinc-700 dark:text-zinc-400 dark:group-hover:text-zinc-200 transition-colors"
                                    />
                                </Link>
                            {/snippet}
                        </TooltipTrigger>
                        <TooltipContent>
                            <p>Документация</p>
                        </TooltipContent>
                    </Tooltip>
                </TooltipProvider>
            </div>

            <!-- Separator -->
            <div
                class="mx-1 hidden h-5 w-px bg-zinc-200 dark:bg-zinc-700 lg:block"
            ></div>

            <!-- Profile Menu -->
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    {#snippet children(props)}
                        <Button
                            variant="ghost"
                            size="icon"
                            class="relative size-9 w-auto rounded-full p-0.5 transition-all hover:ring-2 hover:ring-[#00AAD9]/30"
                            onclick={props.onclick}
                            aria-expanded={props['aria-expanded']}
                            data-state={props['data-state']}
                        >
                            <Avatar
                                class="size-8 overflow-hidden rounded-full ring-2 ring-zinc-200 dark:ring-zinc-700"
                            >
                                {#if auth.user?.avatar}
                                    <AvatarImage
                                        src={auth.user.avatar}
                                        alt={auth.user?.name}
                                    />
                                {/if}
                                <AvatarFallback
                                    class="rounded-full bg-[#00AAD9] text-xs font-bold text-white"
                                >
                                    {getInitials(auth.user?.name ?? '')}
                                </AvatarFallback>
                            </Avatar>
                        </Button>
                    {/snippet}
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" class="w-56">
                    {#if auth.user}
                        <UserMenuContent user={auth.user} />
                    {/if}
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
    </div>
</header>
