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
	class="sidebar fixed z-50 flex flex-col rounded-[2rem] transition-all duration-300"
	style="
		left: 1rem;
		top: 1rem;
		bottom: 1rem;
		width: {collapsed ? '80px' : 'var(--sidebar-width)'};
		background: rgba(var(--surface-rgb), 0.7);
		backdrop-filter: blur(24px);
		box-shadow: var(--shadow-xl);
		border: 1px solid var(--border);
		overflow: hidden;
	"
>
	<!-- Brand -->
	<div
		class="gap-3 flex items-center {collapsed ? 'justify-center' : 'px-6'}"
		style="height: 80px; min-height: 80px;"
	>
		<div
			class="h-11 w-11 rounded-2xl text-sm font-black text-white shadow-lg flex flex-shrink-0 items-center justify-center transition-transform hover:scale-105 active:scale-95"
			style="background: linear-gradient(135deg, var(--brand), #7c3aed)"
		>
			GH
		</div>
		{#if !collapsed}
			<div class="leading-tight animate-fade-in flex flex-col overflow-hidden">
				<span class="text-base font-bold truncate" style="color: var(--text-primary)"
					>GoHeadless</span
				>
				<span
					class="text-xs font-medium tracking-wide uppercase opacity-60"
					style="color: var(--text-muted)">Admin Panel</span
				>
			</div>
		{/if}
	</div>

	<!-- Nav -->
	<nav class="gap-1.5 p-3 flex flex-shrink-0 flex-col">
		{#each navItems as item}
			<a
				href={item.href}
				class="rounded-2xl px-3 py-2.5 group flex items-center transition-all duration-150 {collapsed
					? 'justify-center'
					: 'gap-3'}"
				style="
					color: {isActive(item.href) ? 'var(--brand)' : 'var(--text-secondary)'};
					background: {isActive(item.href) ? 'var(--brand-light)' : 'transparent'};
				"
				onmouseenter={(e) => {
					if (!isActive(item.href))
						(e.currentTarget as HTMLElement).style.background = 'var(--surface-alt)';
				}}
				onmouseleave={(e) => {
					if (!isActive(item.href))
						(e.currentTarget as HTMLElement).style.background = 'transparent';
				}}
			>
				<span class="text-xl flex-shrink-0 transition-transform group-hover:scale-110"
					>{item.icon}</span
				>
				{#if !collapsed}
					<span class="font-bold tracking-tight truncate">{item.label}</span>
				{/if}
			</a>
		{/each}
	</nav>

	{#if $auth.user?.is_initial_admin}
		<div class="px-6 py-2 mt-2 {collapsed ? 'hidden' : ''}">
			<span class="section-title font-black text-[10px] opacity-40">Management</span>
		</div>
		<nav class="gap-1.5 px-3 flex flex-shrink-0 flex-col">
			{#each adminItems as item}
				<a
					href={item.href}
					class="rounded-2xl px-3 py-2.5 group flex items-center transition-all duration-150 {collapsed
						? 'justify-center'
						: 'gap-3'}"
					style="
						color: {isActive(item.href) ? 'var(--brand)' : 'var(--text-secondary)'};
						background: {isActive(item.href) ? 'var(--brand-light)' : 'transparent'};
					"
					onmouseenter={(e) => {
						if (!isActive(item.href))
							(e.currentTarget as HTMLElement).style.background = 'var(--surface-alt)';
					}}
					onmouseleave={(e) => {
						if (!isActive(item.href))
							(e.currentTarget as HTMLElement).style.background = 'transparent';
					}}
				>
					<span class="text-xl flex-shrink-0 transition-transform group-hover:scale-110"
						>{item.icon}</span
					>
					{#if !collapsed}
						<span class="font-bold tracking-tight truncate">{item.label}</span>
					{/if}
				</a>
			{/each}
		</nav>
	{/if}

	<hr style="border-color: var(--border); margin: 0" />

	<!-- Collections list -->
	{#if !collapsed}
		<div class="pt-3 flex flex-1 flex-col overflow-hidden">
			<div class="px-4 pb-2 flex items-center justify-between">
				<span class="section-title">Collections</span>
				<a
					href="/collections/new"
					title="New Collection"
					class="h-6 w-6 rounded-md text-base flex items-center justify-center transition-all hover:text-[var(--brand)]"
					style="color: var(--text-muted)">+</a
				>
			</div>
			<div class="px-3 pb-3 flex-1 overflow-y-auto">
				{#each collections as col}
					<a
						href="/content/{col.slug}"
						class="gap-2 rounded-lg px-3 py-2 text-sm flex items-center transition-all duration-150"
						style="
							color: {page.url.pathname.includes(col.slug) ? 'var(--brand)' : 'var(--text-secondary)'};
							background: {page.url.pathname.includes(col.slug) ? 'var(--brand-light)' : 'transparent'};
						"
						onmouseenter={(e) => {
							if (!page.url.pathname.includes(col.slug))
								(e.currentTarget as HTMLElement).style.background = 'var(--surface-alt)';
						}}
						onmouseleave={(e) => {
							if (!page.url.pathname.includes(col.slug))
								(e.currentTarget as HTMLElement).style.background = 'transparent';
						}}
					>
						<span class="text-xs flex-shrink-0" style="color: var(--text-muted)">▸</span>
						<span class="truncate">{col.name}</span>
					</a>
				{:else}
					<p class="px-3 py-2 text-xs" style="color: var(--text-muted)">No collections yet</p>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Collapse toggle -->
	<div class="p-4">
		<button
			onclick={() => (collapsed = !collapsed)}
			class="gap-3 rounded-xl px-3 py-2.5 text-sm font-medium group flex w-full items-center justify-center transition-all hover:bg-[var(--surface-alt)] hover:text-[var(--brand)]"
			style="color: var(--text-secondary)"
			title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
		>
			<span class="text-lg transition-transform group-hover:scale-110">{collapsed ? '→' : '←'}</span
			>
			{#if !collapsed}<span class="animate-fade-in">Collapse Sidebar</span>{/if}
		</button>
	</div>
</aside>
