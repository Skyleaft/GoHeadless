<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { contentApi } from '$lib/api/content';
	import { toast } from '$lib/stores/toast';
	import { auth } from '$lib/stores/auth';
	import ContentFilterSidebar from '$lib/features/content/ContentFilterSidebar.svelte';
	import RecordTable from '$lib/features/content/RecordTable.svelte';
	import Button from '$lib/shared/Button.svelte';
	import Input from '$lib/shared/Input.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let slug = $derived(data.slug);
	let collection = $derived(data.collection);

	let records = $derived(data.records);
	let total = $derived(data.total);
	let pageNum = $derived(data.page);
	let limitNum = $derived(data.limit);

	let filterOpen = $state(false);
	let searchDraft = $state('');

	$effect(() => {
		searchDraft = $page.url.searchParams.get('search') ?? '';
	});

	let searchTimer: ReturnType<typeof setTimeout> | null = null;

	function onSearchInput() {
		if (searchTimer) clearTimeout(searchTimer);
		searchTimer = setTimeout(() => {
			const u = new URL($page.url);
			const t = searchDraft.trim();
			if (t) u.searchParams.set('search', t);
			else u.searchParams.delete('search');
			u.searchParams.set('page', '1');
			goto(u.pathname + u.search, { replaceState: true, noScroll: true });
		}, 300);
	}

	function navigateWithQuery(pathWithSearch: string) {
		goto(pathWithSearch);
	}

	function onSortClick(fieldKey: string) {
		const u = new URL($page.url);
		const cur = u.searchParams.get('sort');
		if (cur === fieldKey) u.searchParams.set('sort', `-${fieldKey}`);
		else if (cur === `-${fieldKey}`) u.searchParams.delete('sort');
		else u.searchParams.set('sort', fieldKey);
		u.searchParams.set('page', '1');
		goto(u.pathname + u.search, { keepFocus: true, noScroll: true });
	}

	let sortParam = $derived($page.url.searchParams.get('sort'));
	let totalPages = $derived(Math.max(1, Math.ceil(total / limitNum)));

	function goPage(next: number) {
		const u = new URL($page.url);
		u.searchParams.set('page', String(next));
		goto(u.pathname + u.search, { keepFocus: true, noScroll: true });
	}

	let permissions = $derived(() => {
		if (!$auth.user) return { create: false, read: false, update: false, delete: false };
		if ($auth.user.is_initial_admin)
			return { create: true, read: true, update: true, delete: true };

		const access = collection?.access;
		if (!access) return { create: true, read: true, update: true, delete: true };

		const roleId = $auth.user.role_id;
		const check = (action: 'create' | 'read' | 'update' | 'delete') => {
			const roles = access.crud_policy?.[action] ?? [];
			return roles.length === 0 || roles.includes(roleId) || roles.includes('Admin');
		};

		return {
			create: check('create'),
			read: check('read'),
			update: check('update'),
			delete: check('delete')
		};
	});

	async function handleDelete(id: string) {
		try {
			await contentApi.delete(slug, id);
			await invalidateAll();
			toast.success('Record deleted');
		} catch (err: any) {
			toast.error(err.message ?? 'Failed to delete record');
		}
	}

	function handleNew() {
		goto(`/content/${slug}/new`);
	}
</script>

<svelte:head>
	<title>{collection?.name ?? slug} Content — GoHeadless CMS</title>
</svelte:head>

{#if collection}
	<div class="gap-6 animate-fade-in flex flex-col">
		<div class="gap-4 flex items-start justify-between">
			<div class="gap-1 flex flex-col">
				<div class="gap-3 flex items-center">
					<h1 class="text-2xl font-bold" style="color: var(--text-primary)">{collection.name}</h1>
					<span
						class="px-2.5 py-0.5 text-xs font-medium rounded-full"
						style="background: var(--surface-alt); color: var(--text-muted)"
						>{total} records</span
					>
				</div>
				<div class="gap-3 flex items-center">
					<code class="text-sm font-mono" style="color: var(--text-muted)">/{collection.slug}</code>
					<a
						href="/collections/{slug}"
						class="text-xs transition hover:underline"
						style="color: var(--text-muted)">View schema →</a
					>
				</div>
			</div>
			{#if permissions().create}
				<Button variant="primary" onclick={handleNew}>+ New Record</Button>
			{/if}
		</div>

		<div class="gap-3 flex flex-col sm:flex-row sm:items-end sm:justify-between">
			<div class="max-w-md flex-1">
				<Input
					label="Search"
					placeholder="Search searchable fields…"
					bind:value={searchDraft}
					oninput={onSearchInput}
				/>
			</div>
			<Button variant="secondary" onclick={() => (filterOpen = true)}>Filters</Button>
		</div>

		<RecordTable
			fields={collection.fields}
			{records}
			collectionSlug={slug}
			ondelete={handleDelete}
			canUpdate={permissions().update}
			canDelete={permissions().delete}
			sortParam={sortParam}
			onSortClick={onSortClick}
			totalCount={total}
		/>

		{#if totalPages > 1}
			<div
				class="gap-3 px-1 flex flex-wrap items-center justify-between text-sm"
				style="color: var(--text-muted)"
			>
				<button
					type="button"
					class="rounded-lg px-3 py-1.5 font-medium transition disabled:opacity-40"
					style="background: var(--surface-alt); color: var(--text-primary)"
					disabled={pageNum <= 1}
					onclick={() => goPage(pageNum - 1)}
				>
					Previous
				</button>
				<span>Page {pageNum} of {totalPages}</span>
				<button
					type="button"
					class="rounded-lg px-3 py-1.5 font-medium transition disabled:opacity-40"
					style="background: var(--surface-alt); color: var(--text-primary)"
					disabled={pageNum >= totalPages}
					onclick={() => goPage(pageNum + 1)}
				>
					Next
				</button>
			</div>
		{/if}
	</div>

	<ContentFilterSidebar
		open={filterOpen}
		fields={collection.fields}
		onclose={() => (filterOpen = false)}
		currentUrl={$page.url}
		navigate={navigateWithQuery}
	/>
{/if}
