<script lang="ts">
	import { collectionsStore } from '$lib/stores/collections';
	import Spinner from '$lib/shared/Spinner.svelte';
	import Button from '$lib/shared/Button.svelte';

	let collections = $derived($collectionsStore);
	let loading = $state(true);
	let totalRecords = $state<Record<string, number>>({});

	$effect(() => {
		collectionsStore.load().then(async () => {
			loading = false;
			// Fetch record counts per collection
			const counts: Record<string, number> = {};
			await Promise.all(
				collections.map(async (col) => {
					try {
						const res = await fetch(`/api/v1/content/${col.slug}`);
						const data = await res.json();
						counts[col.slug] = Array.isArray(data) ? data.length : 0;
					} catch {
						counts[col.slug] = 0;
					}
				})
			);
			totalRecords = counts;
		});
	});

	let totalRecordCount = $derived(Object.values(totalRecords).reduce((a, b) => a + b, 0));
</script>

<svelte:head>
	<title>Dashboard — GoHeadless CMS</title>
</svelte:head>

<div class="flex flex-col gap-8 animate-fade-in">
	<!-- Header -->
	<div class="flex items-start justify-between">
		<div>
			<h1 class="text-2xl font-bold" style="color: var(--text-primary)">Dashboard</h1>
			<p class="mt-1 text-sm" style="color: var(--text-muted)">
				Welcome to your GoHeadless CMS control panel.
			</p>
		</div>
		<Button variant="primary" onclick={() => (window.location.href = '/collections/new')}>
			+ New Collection
		</Button>
	</div>

	<!-- Stats -->
	{#if loading}
		<div class="flex justify-center py-12"><Spinner /></div>
	{:else}
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
			<!-- Collections stat -->
			<div class="card p-5 flex items-center gap-4">
				<div
					class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl text-2xl"
					style="background: var(--brand-light)"
				>🗂</div>
				<div>
					<p class="text-2xl font-bold" style="color: var(--text-primary)">{collections.length}</p>
					<p class="text-sm" style="color: var(--text-muted)">Collections</p>
				</div>
			</div>

			<!-- Records stat -->
			<div class="card p-5 flex items-center gap-4">
				<div
					class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl text-2xl"
					style="background: #f0f9ff"
				>📄</div>
				<div>
					<p class="text-2xl font-bold" style="color: var(--text-primary)">{totalRecordCount}</p>
					<p class="text-sm" style="color: var(--text-muted)">Total Records</p>
				</div>
			</div>

			<!-- API status -->
			<div class="card p-5 flex items-center gap-4">
				<div
					class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl text-2xl"
					style="background: #f0fdf4"
				>⚡</div>
				<div>
					<p class="text-2xl font-bold text-green-500">Online</p>
					<p class="text-sm" style="color: var(--text-muted)">API Status</p>
				</div>
			</div>
		</div>

		<!-- Collections quick-list -->
		{#if collections.length > 0}
			<div class="card">
				<div class="card-header flex items-center justify-between">
					<h2 class="text-base font-semibold" style="color: var(--text-primary)">Collections</h2>
					<a href="/collections" class="text-sm transition hover:underline" style="color: var(--brand)">
						View all →
					</a>
				</div>
				<div class="divide-y" style="border-color: var(--border)">
					{#each collections as col}
						<div class="flex items-center gap-4 px-6 py-4">
							<div
								class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg text-base"
								style="background: var(--brand-light); color: var(--brand)"
							>🗂</div>
							<div class="flex flex-col flex-1 min-w-0">
								<span class="text-sm font-medium truncate" style="color: var(--text-primary)">{col.name}</span>
								<span class="text-xs font-mono" style="color: var(--text-muted)">/{col.slug}</span>
							</div>
							<div class="flex items-center gap-3">
								<span
									class="rounded-full px-2.5 py-0.5 text-xs font-medium"
									style="background: var(--surface-alt); color: var(--text-muted)"
								>
									{totalRecords[col.slug] ?? 0} records
								</span>
								<a
									href="/content/{col.slug}"
									class="text-sm font-medium transition hover:underline"
									style="color: var(--brand)"
								>Open →</a>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{:else}
			<div class="card p-12 text-center">
				<p class="text-4xl mb-4">🚀</p>
				<h3 class="text-base font-semibold mb-2" style="color: var(--text-primary)">
					Get started with GoHeadless
				</h3>
				<p class="text-sm mb-6" style="color: var(--text-muted)">
					Create your first collection to define a schema and start managing content.
				</p>
				<a
					href="/collections/new"
					class="inline-flex items-center gap-2 rounded-lg px-5 py-2.5 text-sm font-medium text-white transition hover:opacity-90"
					style="background: var(--brand)"
				>+ Create First Collection</a>
			</div>
		{/if}
	{/if}
</div>
