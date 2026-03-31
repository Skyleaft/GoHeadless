<script lang="ts">
	import { goto } from '$app/navigation';
	import { contentApi } from '$lib/api/content';
	import { toast } from '$lib/stores/toast';
	import { auth } from '$lib/stores/auth';
	import RecordTable from '$lib/features/content/RecordTable.svelte';
	import Button from '$lib/shared/Button.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let slug = $derived(data.slug);
	let collection = $derived(data.collection);

	// Create a reactive state that updates when the data prop changes
	let records = $state(data.records);
	$effect(() => {
		records = data.records;
	});

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
			records = records.filter((r: any) => (r._id ?? r.id) !== id);
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
				<Button variant="primary" onclick={handleNew}>+ New Record</Button>
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
