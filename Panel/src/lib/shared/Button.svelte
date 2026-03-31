<script lang="ts">
	interface Props {
		variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
		size?: 'sm' | 'md' | 'lg';
		loading?: boolean;
		disabled?: boolean;
		type?: 'button' | 'submit' | 'reset';
		href?: string;
		onclick?: (e: MouseEvent) => void;
		class?: string;
		children: import('svelte').Snippet;
	}

	let {
		variant = 'primary',
		size = 'md',
		loading = false,
		disabled = false,
		type = 'button',
		href,
		onclick,
		class: cls = '',
		children
	}: Props = $props();

	const base =
		'inline-flex items-center justify-center gap-3 font-black rounded-2xl transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 cursor-pointer select-none';

	const sizes = {
		sm: 'h-9 px-4 text-xs',
		md: 'h-11 px-6 text-sm',
		lg: 'h-14 px-8 text-base tracking-tight'
	};

	const variants = {
		primary:
			'bg-[var(--brand)] text-white hover:brightness-110 shadow-lg shadow-[var(--brand)]/20 active:scale-95 focus-visible:ring-[var(--brand)]',
		secondary:
			'bg-[var(--surface-alt)] text-[var(--text-primary)] shadow-sm hover:scale-[1.02] active:scale-95 focus-visible:ring-[var(--brand)]',
		ghost:
			'text-[var(--text-secondary)] hover:bg-[var(--surface-alt)] hover:text-[var(--text-primary)] active:scale-95 focus-visible:ring-[var(--brand)]',
		danger:
			'bg-red-500 text-white hover:bg-red-600 shadow-lg shadow-red-500/20 active:scale-95 focus-visible:ring-red-500'
	};
</script>

{#if href}
	<a
		{href}
		class="{base} {sizes[size]} {variants[variant]} {cls}"
		{onclick}
	>
		{@render children()}
	</a>
{:else}
	<button
		{type}
		class="{base} {sizes[size]} {variants[variant]} {cls}"
		disabled={disabled || loading}
		{onclick}
	>
		{#if loading}
			<span class="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
			></span>
		{/if}
		{@render children()}
	</button>
{/if}

<style>
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	.animate-spin {
		animation: spin 0.7s linear infinite;
	}
</style>
