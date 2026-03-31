<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { adminApi } from '$lib/api/admin';
	import { toast } from '$lib/stores/toast';

	let stats = $state<Record<string, number>>({});
	let loading = $state(true);

	onMount(async () => {
		if ($auth.user?.is_initial_admin) {
			try {
				stats = await adminApi.getStats();
			} catch (err: any) {
				toast.error('Failed to load dashboard statistics');
			} finally {
				loading = false;
			}
		} else {
			loading = false;
		}
	});

	const statCards = $derived([
		{ label: 'Total Collections', value: stats.collections || 0, icon: '🗂', color: 'blue' },
		{ label: 'Total Records', value: stats.records || 0, icon: '📄', color: 'green' },
		{ label: 'System Users', value: stats.users || 0, icon: '👥', color: 'purple' },
		{ label: 'Security Roles', value: stats.roles || 0, icon: '🛡️', color: 'amber' }
	]);
</script>

<div class="flex flex-col gap-8 animate-fade-in">
	<div class="flex flex-col gap-2">
		<h1 class="text-3xl font-bold tracking-tight">Welcome back, {$auth.user?.username}!</h1>
		<p class="text-[var(--text-secondary)]">You are logged in as <span class="font-semibold text-[var(--brand)]">{$auth.user?.is_initial_admin ? 'Superadmin' : 'Staff'}</span>. Here is what's happening today.</p>
	</div>

	{#if $auth.user?.is_initial_admin}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			{#each statCards as card}
				<div class="card p-6 flex flex-col gap-4 group hover:shadow-lg transition-all border-l-4" style="border-left-color: var(--{card.color}-500, var(--brand))">
					<div class="flex items-center justify-between">
						<span class="text-2xl">{card.icon}</span>
						{#if loading}
							<div class="h-4 w-12 bg-[--surface-alt] animate-pulse rounded"></div>
						{:else}
							<span class="text-2xl font-bold leading-none">{card.value}</span>
						{/if}
					</div>
					<span class="text-sm font-medium text-[--text-muted]">{card.label}</span>
				</div>
			{/each}
		</div>
	{/if}

	<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
		<div class="lg:col-span-2 flex flex-col gap-6">
			<div class="card p-8 bg-gradient-to-br from-[var(--brand)] to-[#7c3aed] text-white overflow-hidden relative">
				<div class="relative z-10">
					<h2 class="text-xl font-bold mb-4">Start Building Your API</h2>
					<p class="text-white/80 text-sm mb-6 max-w-md">Create your first collection to define your content structure. GoHeadless will automatically generate a high-performance REST API for you.</p>
					<a href="/collections/new" class="inline-flex h-10 items-center px-6 rounded-lg bg-white text-[var(--brand)] font-bold text-sm transition-transform active:scale-95 shadow-xl">
						Create New Collection
					</a>
				</div>
				<div class="absolute -right-10 -bottom-10 text-[160px] opacity-10 pointer-events-none rotate-12">
					🏗️
				</div>
			</div>

			<div class="card p-6">
				<h3 class="font-bold mb-4">Quick Actions</h3>
				<div class="grid grid-cols-2 gap-3">
					<a href="/collections" class="p-3 rounded-lg border border-[var(--border)] hover:bg-[var(--surface-alt)] transition-colors flex items-center gap-3 text-sm">
						<span class="text-lg">🗂</span> Manage Collections
					</a>
					<a href="/admin/users" class="p-3 rounded-lg border border-[var(--border)] hover:bg-[var(--surface-alt)] transition-colors flex items-center gap-3 text-sm">
						<span class="text-lg">👥</span> Manage Users
					</a>
					<a href="/docs" class="p-3 rounded-lg border border-[var(--border)] hover:bg-[var(--surface-alt)] transition-colors flex items-center gap-3 text-sm">
						<span class="text-lg">📄</span> API Documentation
					</a>
					<button onclick={() => window.open('https://github.com/Skyleaft/GoHeadless', '_blank')} class="p-3 rounded-lg border border-[var(--border)] hover:bg-[var(--surface-alt)] transition-colors flex items-center gap-3 text-sm text-left">
						<span class="text-lg">⭐</span> GitHub Repository
					</button>
				</div>
			</div>
		</div>

		<div class="flex flex-col gap-6">
			<div class="card p-6">
				<h3 class="font-bold mb-4">System Status</h3>
				<div class="space-y-4">
					<div class="flex items-center justify-between text-sm">
						<span class="text-[--text-muted]">Backend Version</span>
						<span class="font-mono bg-[--surface-alt] px-2 py-0.5 rounded text-xs">v1.0.0</span>
					</div>
					<div class="flex items-center justify-between text-sm">
						<span class="text-[--text-muted]">Database</span>
						<span class="text-emerald-500 font-medium">Connected</span>
					</div>
					<div class="flex items-center justify-between text-sm">
						<span class="text-[--text-muted]">Self-Healing Setup</span>
						<span class="text-emerald-500 font-medium">Active</span>
					</div>
					<div class="flex items-center justify-between text-sm">
						<span class="text-[--text-muted]">Storage</span>
						<span class="text-[--text-secondary]">Local (./uploads)</span>
					</div>
				</div>
			</div>

			<div class="card p-6 bg-[var(--brand-light)]/30 border-[var(--brand-light)]">
				<h3 class="font-bold mb-2 text-[var(--brand)]">Need Help?</h3>
				<p class="text-xs text-[var(--text-secondary)] leading-relaxed mb-4">GoHeadless is designed for speed and flexibility. Check out our documentation or join our community for support.</p>
				<button class="text-xs font-bold text-[var(--brand)] hover:underline">Documentation Hub →</button>
			</div>
		</div>
	</div>
</div>
