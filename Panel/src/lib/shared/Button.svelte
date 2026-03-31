<script lang="ts">
	interface Props {
		variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
		size?: 'sm' | 'md' | 'lg';
		loading?: boolean;
		disabled?: boolean;
		type?: 'button' | 'submit' | 'reset';
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
		onclick,
		class: cls = '',
		children
	}: Props = $props();

	const base =
		'inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 cursor-pointer select-none';

	const sizes = {
		sm: 'h-8 px-3 text-xs',
		md: 'h-9 px-4 text-sm',
		lg: 'h-11 px-6 text-base'
	};

	const variants = {
		primary:
			'bg-[--brand] text-white hover:bg-[--brand-hover] focus-visible:ring-[--brand] shadow-sm active:scale-[0.98]',
		secondary:
			'bg-[--surface-alt] text-[--text-primary] border border-[--border] hover:bg-[--surface-hover] focus-visible:ring-[--brand]',
		ghost: 'text-[--text-secondary] hover:bg-[--surface-alt] hover:text-[--text-primary] focus-visible:ring-[--brand]',
		danger:
			'bg-red-500 text-white hover:bg-red-600 focus-visible:ring-red-500 shadow-sm active:scale-[0.98]'
	};
</script>

<button
	{type}
	class="{base} {sizes[size]} {variants[variant]} {cls}"
	disabled={disabled || loading}
	{onclick}
>
	{#if loading}
		<span
			class="h-4 w-4 rounded-full border-2 border-current border-t-transparent animate-spin"
		></span>
	{/if}
	{@render children()}
</button>

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
