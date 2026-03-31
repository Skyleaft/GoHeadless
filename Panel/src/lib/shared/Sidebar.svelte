<script lang="ts">
	import { page } from '$app/state';
	import { collectionsStore } from '$lib/stores/collections';
	import { auth } from '$lib/stores/auth';
	import Spinner from './Spinner.svelte';

	interface Props {
		collapsed?: boolean;
	}
	let { collapsed = $bindable(false) }: Props = $props();

	let collections = $derived($collectionsStore);

	const navItems = [
		{ href: '/', icon: '⊞', label: 'Dashboard' },
		{ href: '/collections', icon: '🗂', label: 'Collections' }
	];

	const adminItems = [
		{ href: '/admin/users', icon: '👥', label: 'Users' },
		{ href: '/admin/roles', icon: '🛡️', label: 'Roles' }
	];

	function isActive(href: string) {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname.startsWith(href);
	}
</script>

<aside
	class="sidebar flex flex-col border-r transition-all duration-300"
	style="
		width: {collapsed ? '64px' : 'var(--sidebar-width)'};
		background: var(--surface);
		border-color: var(--border);
		height: 100vh;
		overflow: hidden;
		flex-shrink: 0;
	"
>
	<!-- Brand -->
	<div
		class="flex items-center gap-3 border-b px-4"
		style="height: var(--topbar-height); border-color: var(--border); min-height: var(--topbar-height);"
	>
		<div
			class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg text-sm font-bold text-white"
			style="background: var(--brand)"
		>
			GH
		</div>
		{#if !collapsed}
			<div class="flex flex-col leading-tight overflow-hidden">
				<span class="text-sm font-semibold truncate" style="color: var(--text-primary)">GoHeadless</span>
				<span class="text-xs" style="color: var(--text-muted)">CMS Panel</span>
			</div>
		{/if}
	</div>

	<!-- Nav -->
	<nav class="flex flex-col gap-1 p-3 flex-shrink-0">
		{#each navItems as item}
			<a
				href={item.href}
				class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-150"
				style="
					color: {isActive(item.href) ? 'var(--brand)' : 'var(--text-secondary)'};
					background: {isActive(item.href) ? 'var(--brand-light)' : 'transparent'};
				"
				onmouseenter={(e) => { if (!isActive(item.href)) (e.currentTarget as HTMLElement).style.background = 'var(--surface-alt)'; }}
				onmouseleave={(e) => { if (!isActive(item.href)) (e.currentTarget as HTMLElement).style.background = 'transparent'; }}
			>
				<span class="text-base flex-shrink-0">{item.icon}</span>
				{#if !collapsed}
					<span class="truncate">{item.label}</span>
				{/if}
			</a>
		{/each}
	</nav>

	{#if $auth.user?.is_initial_admin}
		<div class="px-4 py-2 mt-2">
			<span class="section-title">Administration</span>
		</div>
		<nav class="flex flex-col gap-1 px-3 flex-shrink-0">
			{#each adminItems as item}
				<a
					href={item.href}
					class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-150"
					style="
						color: {isActive(item.href) ? 'var(--brand)' : 'var(--text-secondary)'};
						background: {isActive(item.href) ? 'var(--brand-light)' : 'transparent'};
					"
					onmouseenter={(e) => { if (!isActive(item.href)) (e.currentTarget as HTMLElement).style.background = 'var(--surface-alt)'; }}
					onmouseleave={(e) => { if (!isActive(item.href)) (e.currentTarget as HTMLElement).style.background = 'transparent'; }}
				>
					<span class="text-base flex-shrink-0">{item.icon}</span>
					{#if !collapsed}
						<span class="truncate">{item.label}</span>
					{/if}
				</a>
			{/each}
		</nav>
	{/if}

	<hr style="border-color: var(--border); margin: 0" />

	<!-- Collections list -->
	{#if !collapsed}
		<div class="flex flex-col flex-1 overflow-hidden pt-3">
			<div class="flex items-center justify-between px-4 pb-2">
				<span class="section-title">Collections</span>
				<a
					href="/collections/new"
					title="New Collection"
					class="flex h-6 w-6 items-center justify-center rounded-md text-base transition-all hover:text-[--brand]"
					style="color: var(--text-muted)"
				>+</a>
			</div>
			<div class="flex-1 overflow-y-auto px-3 pb-3">
				{#each collections as col}
					<a
						href="/content/{col.slug}"
						class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-all duration-150"
						style="
							color: {page.url.pathname.includes(col.slug) ? 'var(--brand)' : 'var(--text-secondary)'};
							background: {page.url.pathname.includes(col.slug) ? 'var(--brand-light)' : 'transparent'};
						"
						onmouseenter={(e) => { if (!page.url.pathname.includes(col.slug)) (e.currentTarget as HTMLElement).style.background = 'var(--surface-alt)'; }}
						onmouseleave={(e) => { if (!page.url.pathname.includes(col.slug)) (e.currentTarget as HTMLElement).style.background = 'transparent'; }}
					>
						<span class="flex-shrink-0 text-xs" style="color: var(--text-muted)">▸</span>
						<span class="truncate">{col.name}</span>
					</a>
				{:else}
					<p class="px-3 py-2 text-xs" style="color: var(--text-muted)">No collections yet</p>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Collapse toggle -->
	<div class="border-t p-3" style="border-color: var(--border)">
		<button
			onclick={() => (collapsed = !collapsed)}
			class="flex w-full items-center justify-center gap-2 rounded-lg px-3 py-2 text-sm transition-all hover:bg-[--surface-alt]"
			style="color: var(--text-muted)"
			title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
		>
			<span class="text-base">{collapsed ? '→' : '←'}</span>
			{#if !collapsed}<span>Collapse</span>{/if}
		</button>
	</div>
</aside>
