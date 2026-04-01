<script lang="ts">
	import type { Field, FieldType } from '$lib/types/collection';
	import Button from '$lib/shared/Button.svelte';
	import { flattenFieldsForFilter } from '$lib/utils/flattenFields';

	interface Props {
		open: boolean;
		fields: Field[];
		onclose: () => void;
		/** Current page URL — filters are read/applied via searchParams */
		currentUrl: URL;
		/** Navigate with updated query string */
		navigate: (pathWithSearch: string) => void;
	}

	let { open, fields, onclose, currentUrl, navigate }: Props = $props();

	let flat = $derived(flattenFieldsForFilter(fields));

	function filterable(f: Field): boolean {
		if (f.internal) return false;
		const t = f.type as FieldType;
		return !['submit', 'reset', 'action', 'file', 'image', 'relation', 'autocomplete'].includes(t);
	}

	function rfc3339ToLocalDatetime(v: string): string {
		if (!v) return '';
		const d = new Date(v);
		if (Number.isNaN(d.getTime())) return '';
		const pad = (n: number) => String(n).padStart(2, '0');
		return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
	}

	function localInputToRFC3339(v: string): string {
		if (!v) return '';
		const t = new Date(v);
		if (Number.isNaN(t.getTime())) return v;
		return t.toISOString();
	}

	let draft = $state<Record<string, string>>({});

	$effect(() => {
		if (!open) return;
		const next: Record<string, string> = {};
		for (const f of flat) {
			if (!filterable(f)) continue;
			const v = currentUrl.searchParams.get(`filter[${f.key}]`);
			const gte = currentUrl.searchParams.get(`filter[${f.key}][gte]`);
			const lte = currentUrl.searchParams.get(`filter[${f.key}][lte]`);
			if (f.type === 'datepicker' || f.type === 'datetimepicker' || f.type === 'timepicker') {
				next[`${f.key}:gte`] = gte ? rfc3339ToLocalDatetime(gte) : '';
				next[`${f.key}:lte`] = lte ? rfc3339ToLocalDatetime(lte) : '';
			} else if (f.type === 'bool' || f.type === 'toggle') {
				next[f.key] =
					v === 'true' || v === '1' ? 'true' : v === 'false' || v === '0' ? 'false' : '';
			} else {
				next[f.key] = v ?? '';
			}
		}
		draft = next;
	});

	function stripFilterParams(url: URL): URL {
		const u = new URL(url);
		for (const key of [...u.searchParams.keys()]) {
			if (key.startsWith('filter[')) u.searchParams.delete(key);
		}
		return u;
	}

	function apply() {
		const u = stripFilterParams(currentUrl);
		u.searchParams.set('page', '1');

		for (const f of flat) {
			if (!filterable(f)) continue;
			if (f.type === 'datepicker' || f.type === 'datetimepicker' || f.type === 'timepicker') {
				const gte = draft[`${f.key}:gte`]?.trim();
				const lte = draft[`${f.key}:lte`]?.trim();
				if (gte) u.searchParams.set(`filter[${f.key}][gte]`, localInputToRFC3339(gte));
				if (lte) u.searchParams.set(`filter[${f.key}][lte]`, localInputToRFC3339(lte));
			} else if (f.type === 'bool' || f.type === 'toggle') {
				if (draft[f.key] === 'true') u.searchParams.set(`filter[${f.key}]`, 'true');
				else if (draft[f.key] === 'false') u.searchParams.set(`filter[${f.key}]`, 'false');
			} else {
				const v = draft[f.key]?.trim();
				if (v) u.searchParams.set(`filter[${f.key}]`, v);
			}
		}

		navigate(u.pathname + u.search);
		onclose();
	}

	function clearFilters() {
		const u = stripFilterParams(currentUrl);
		u.searchParams.set('page', '1');
		navigate(u.pathname + u.search);
		onclose();
	}
</script>

{#if open}
	<button
		type="button"
		class="fixed inset-0 z-40 bg-black/30 transition-opacity"
		aria-label="Close filters"
		onclick={onclose}
	></button>
	<aside
		class="fixed inset-y-0 right-0 z-50 flex w-full max-w-md flex-col border-l shadow-xl"
		style="background: var(--surface); border-color: var(--border)"
	>
		<div class="flex items-center justify-between border-b px-4 py-3" style="border-color: var(--border)">
			<h2 class="text-lg font-semibold" style="color: var(--text-primary)">Filters</h2>
			<button
				type="button"
				class="rounded-lg p-2 text-xl leading-none hover:bg-[--surface-hover]"
				onclick={onclose}
				aria-label="Close">×</button
			>
		</div>

		<div class="flex-1 overflow-y-auto px-4 py-4 gap-4 flex flex-col">
			{#each flat as f}
				{#if filterable(f)}
					<div class="gap-1.5 flex flex-col">
						<label class="text-xs font-medium uppercase tracking-wide" style="color: var(--text-muted)">
							{f.label || f.key}
						</label>
						{#if f.type === 'bool' || f.type === 'toggle'}
							<label class="gap-2 flex cursor-pointer items-center text-sm">
								<input
									type="checkbox"
									checked={draft[f.key] === 'true'}
									onchange={(e) => {
										draft[f.key] = (e.currentTarget as HTMLInputElement).checked ? 'true' : 'false';
									}}
								/>
								Yes
							</label>
						{:else if f.type === 'number' || f.type === 'slider' || f.type === 'rating'}
							<input
								type="number"
								class="w-full rounded-lg border px-3 py-2 text-sm"
								style="border-color: var(--border); background: var(--surface-alt); color: var(--text-primary)"
								value={draft[f.key] ?? ''}
								oninput={(e) => (draft[f.key] = (e.currentTarget as HTMLInputElement).value)}
							/>
						{:else if f.type === 'datepicker' || f.type === 'datetimepicker' || f.type === 'timepicker'}
							<div class="gap-2 grid grid-cols-2">
								<div class="gap-1 flex flex-col">
									<span class="text-xs" style="color: var(--text-muted)">From</span>
									<input
										type="datetime-local"
										class="w-full rounded-lg border px-2 py-1.5 text-xs"
										style="border-color: var(--border); background: var(--surface-alt)"
										value={draft[`${f.key}:gte`] ?? ''}
										oninput={(e) =>
											(draft[`${f.key}:gte`] = (e.currentTarget as HTMLInputElement).value)}
									/>
								</div>
								<div class="gap-1 flex flex-col">
									<span class="text-xs" style="color: var(--text-muted)">To</span>
									<input
										type="datetime-local"
										class="w-full rounded-lg border px-2 py-1.5 text-xs"
										style="border-color: var(--border); background: var(--surface-alt)"
										value={draft[`${f.key}:lte`] ?? ''}
										oninput={(e) =>
											(draft[`${f.key}:lte`] = (e.currentTarget as HTMLInputElement).value)}
									/>
								</div>
							</div>
						{:else}
							<input
								type="text"
								class="w-full rounded-lg border px-3 py-2 text-sm"
								style="border-color: var(--border); background: var(--surface-alt); color: var(--text-primary)"
								placeholder={f.placeholder ?? ''}
								value={draft[f.key] ?? ''}
								oninput={(e) => (draft[f.key] = (e.currentTarget as HTMLInputElement).value)}
							/>
						{/if}
					</div>
				{/if}
			{/each}
		</div>

		<div
			class="gap-2 border-t px-4 py-3 flex flex-wrap items-center justify-end"
			style="border-color: var(--border); background: var(--surface-alt)"
		>
			<Button variant="ghost" onclick={clearFilters}>Clear</Button>
			<Button variant="primary" onclick={apply}>Apply</Button>
		</div>
	</aside>
{/if}
