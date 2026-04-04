<script lang="ts">
    let {
        class: className = '',
        ...rest
    }: {
        class?: string;
        [key: string]: unknown;
    } = $props();
</script>

<svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 100 100"
    class={className}
    {...rest}
>
    <!-- Stylized Train / Railway KTZ Logo -->
    <!-- Rail tracks -->
    <rect x="10" y="72" width="80" height="4" rx="2" fill="currentColor" opacity="0.9" />
    <rect x="18" y="78" width="4" height="8" rx="1" fill="currentColor" opacity="0.7" />
    <rect x="30" y="78" width="4" height="8" rx="1" fill="currentColor" opacity="0.7" />
    <rect x="42" y="78" width="4" height="8" rx="1" fill="currentColor" opacity="0.7" />
    <rect x="54" y="78" width="4" height="8" rx="1" fill="currentColor" opacity="0.7" />
    <rect x="66" y="78" width="4" height="8" rx="1" fill="currentColor" opacity="0.7" />
    <rect x="78" y="78" width="4" height="8" rx="1" fill="currentColor" opacity="0.7" />

    <!-- Train body -->
    <rect x="20" y="40" width="55" height="32" rx="6" fill="currentColor" opacity="0.95" />
    
    <!-- Train front (locomotive nose) -->
    <path d="M75 40 L85 52 L85 66 Q85 72 79 72 L75 72 Z" fill="currentColor" opacity="0.85" />
    
    <!-- Train windows -->
    <rect x="26" y="46" width="10" height="10" rx="2" fill="currentColor" opacity="0.2" />
    <rect x="40" y="46" width="10" height="10" rx="2" fill="currentColor" opacity="0.2" />
    <rect x="54" y="46" width="10" height="10" rx="2" fill="currentColor" opacity="0.2" />
    
    <!-- Front windshield -->
    <path d="M73 46 L80 52 L80 58 L73 58 Z" fill="currentColor" opacity="0.2" />
    
    <!-- Headlight -->
    <circle cx="83" cy="60" r="3" fill="currentColor" opacity="0.3" />
    
    <!-- Wheels -->
    <circle cx="32" cy="72" r="6" fill="currentColor" opacity="0.8" />
    <circle cx="32" cy="72" r="3" fill="currentColor" opacity="0.2" />
    <circle cx="52" cy="72" r="6" fill="currentColor" opacity="0.8" />
    <circle cx="52" cy="72" r="3" fill="currentColor" opacity="0.2" />
    <circle cx="72" cy="72" r="5" fill="currentColor" opacity="0.8" />
    <circle cx="72" cy="72" r="2.5" fill="currentColor" opacity="0.2" />
    
    <!-- Telemetry signal waves (monitoring aspect) -->
    <path d="M38 32 Q42 20 50 32" stroke="currentColor" fill="none" stroke-width="2.5" stroke-linecap="round" opacity="0.6" />
    <path d="M34 26 Q44 10 54 26" stroke="currentColor" fill="none" stroke-width="2" stroke-linecap="round" opacity="0.4" />
    <path d="M30 20 Q46 0 58 20" stroke="currentColor" fill="none" stroke-width="1.5" stroke-linecap="round" opacity="0.25" />
    
    <!-- Antenna on roof -->
    <rect x="44" y="32" width="2" height="8" rx="1" fill="currentColor" opacity="0.7" />
    <circle cx="45" cy="32" r="2" fill="currentColor" opacity="0.7" />
</svg>
