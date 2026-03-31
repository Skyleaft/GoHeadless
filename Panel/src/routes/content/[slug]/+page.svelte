<script lang="ts">
	import { page } from '$app/state';
	import { collectionsApi } from '$lib/api/collections';
	import { contentApi } from '$lib/api/content';
	import { toast } from '$lib/stores/toast';
	import { auth } from '$lib/stores/auth';
	import RecordTable from '$lib/features/content/RecordTable.svelte';
	import Button from '$lib/shared/Button.svelte';
	import Spinner from '$lib/shared/Spinner.svelte';
	import type { Collection, ContentRecord } from '$lib/types/collection';

	let slug = $derived(page.params.slug as string);
	let collection = $state<Collection | null>(null);
	let records = $state<ContentRecord[]>([]);
	let loading = $state(true);
	let permissions = $derived(() => {
		if (!$auth.user) return { create: false, read: false, update: false, delete: false };
		if ($auth.user.is_initial_admin)
			return { create: true, read: true, update: true, delete: true };

		const access = collection?.access;
		if (!access) return { create: true, read: true, update: true, delete: true }; // Default to open if no policy

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

	$effect(() => {
		loading = true;
		Promise.all([collectionsApi.get(slug), contentApi.list(slug)])
			.then(([coll, recs]) => {
				collection = coll;
				records = recs ?? [];
				loading = false;
			})
			.catch((err) => {
				toast.error(err.message ?? 'Failed to load data');
				loading = false;
			});
	});

	async function handleDelete(id: string) {
		try {
			await contentApi.delete(slug, id);
			records = records.filter((r) => (r._id ?? r.id) !== id);
			toast.success('Record deleted');
		} catch (err: any) {
			toast.error(err.message ?? 'Failed to delete record');
		}
	}
</script>

<svelte:head>
	<title>{collection?.name ?? slug} Content — GoHeadless CMS</title>
</svelte:head>

{#if loading}
	<div class="py-24 flex justify-center"><Spinner /></div>
{:else if collection}
	<div class="gap-6 animate-fade-in flex flex-col">
		<!-- Header -->
		<div class="gap-4 flex items-start justify-between">
			<div class="gap-1 flex flex-col">
				<div class="gap-3 flex items-center">
					<h1 class="text-2xl font-bold" style="color: var(--text-primary)">{collection.name}</h1>
					<span
						class="px-2.5 py-0.5 text-xs font-medium rounded-full"
						style="background: var(--surface-alt); color: var(--text-muted)"
						>{records.length} records</span
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
				<Button variant="primary" onclick={() => (window.location.href = `/content/${slug}/new`)}>
					+ New Record
				</Button>
			{/if}
		</div>

		<RecordTable
			fields={collection.fields}
			{records}
			collectionSlug={slug}
			ondelete={handleDelete}
			canUpdate={permissions().update}
			canDelete={permissions().delete}
		/>
	</div>
{/if}
