<script lang="ts">
	import { page } from '$app/state';
	import { theme, toggleTheme } from '$lib/stores/theme';

	interface Props {
		onMenuToggle?: () => void;
	}
	let { onMenuToggle }: Props = $props();

	// Removed isDark as we'll use $theme directly

	// Build breadcrumbs from current path
	let breadcrumbs = $derived(() => {
		const parts = page.url.pathname.split('/').filter(Boolean);
		const crumbs: { label: string; href: string }[] = [{ label: 'Home', href: '/' }];
		let path = '';
		for (const part of parts) {
			path += `/${part}`;
			crumbs.push({
				label: part.charAt(0).toUpperCase() + part.slice(1).replace(/-/g, ' '),
				href: path
			});
		}
		return crumbs;
	});
</script>

<header
	class="flex items-center justify-between border-b px-6"
	style="
		height: var(--topbar-height);
		background: var(--surface);
		border-color: var(--border);
		flex-shrink: 0;
	"
>
	<!-- Left: Breadcrumb -->
	<div class="flex items-center gap-2 text-sm">
		{#each breadcrumbs() as crumb, i}
			{#if i > 0}
				<span style="color: var(--text-muted)" class="text-xs">/</span>
			{/if}
			{#if i === breadcrumbs().length - 1}
				<span class="font-medium" style="color: var(--text-primary)">{crumb.label}</span>
			{:else}
				<a
					href={crumb.href}
					class="transition-colors hover:text-[--brand]"
					style="color: var(--text-muted)"
				>{crumb.label}</a>
			{/if}
		{/each}
	</div>

	<!-- Right: Actions -->
	<div class="flex items-center gap-2">
		<!-- Theme toggle -->
		<button
			onclick={toggleTheme}
			id="theme-toggle"
			class="flex h-9 w-9 items-center justify-center rounded-lg text-base transition-all hover:bg-[--surface-alt]"
			style="color: var(--text-secondary)"
			title={$theme === 'dark' ? 'Dark Mode' : $theme === 'light' ? 'Light Mode' : 'System Mode (Auto)'}
			aria-label="Toggle theme"
		>
			{$theme === 'dark' ? '🌙' : $theme === 'light' ? '☀' : '💻'}
		</button>

		<!-- API Docs link -->
		<a
			href="http://localhost:3030/docs"
			target="_blank"
			rel="noopener"
			class="flex h-9 items-center gap-1.5 rounded-lg px-3 text-sm font-medium transition-all hover:bg-[--surface-alt]"
			style="color: var(--text-secondary)"
			title="Open API documentation"
		>
			<span>API Docs</span>
			<span class="text-xs">↗</span>
		</a>
	</div>
</header>
