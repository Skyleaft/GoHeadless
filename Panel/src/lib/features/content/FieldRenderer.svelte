<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import { collectionsApi } from '$lib/api/collections';
	import FileUpload from './FileUpload.svelte';

	interface Props {
		field: Field;
		// Use Record<string,any> so bind: targets on data[field.key] are writable
		data: Record<string, any>;
		depth?: number;
	}

	let { field, data = $bindable(), depth = 0 }: Props = $props();

	// Evaluate show_if conditional logic
	let visible = $derived(
		field.logic?.show_if
			? String(data[field.logic.show_if.field]) === String(field.logic.show_if.equals)
			: true
	);

	// Rating hover state
	let ratingHover = $state(0);

	// Relation options loaded on focus
	let relationOptions = $state<{ label: string; value: string }[]>([]);
	let relationLoading = $state(false);
	let relationSearched = $state(false);

	async function loadRelation() {
		if (relationSearched || !field.relation) return;
		relationLoading = true;
		try {
			const res = await fetch(`/api/v1/content/${field.relation.collection}`);
			const records: Record<string, any>[] = await res.json();
			const displayField = field.relation.field;
			relationOptions = records.map((r) => ({
				label: String(r[displayField] ?? r._id ?? ''),
				value: String(r._id ?? r.id ?? '')
			}));
		} catch {
			// ignore
		} finally {
			relationLoading = false;
			relationSearched = true;
		}
	}

	// Repeater — keep a local array state synced to data[field.key]
	let repeaterItems = $state<Record<string, any>[]>(
		Array.isArray(data[field.key]) ? [...(data[field.key] as Record<string, any>[])] : []
	);

	$effect(() => {
		data[field.key] = repeaterItems;
	});

	function addRepeaterItem() {
		repeaterItems = [...repeaterItems, {}];
	}

	function removeRepeaterItem(i: number) {
		repeaterItems = repeaterItems.filter((_, idx) => idx !== i);
	}

	// Group — ensure data[field.key] is an object
	function ensureGroupData(): Record<string, any> {
		if (!data[field.key] || typeof data[field.key] !== 'object' || Array.isArray(data[field.key])) {
			data[field.key] = {};
		}
		return data[field.key] as Record<string, any>;
	}
</script>

{#if visible}
	<div style="margin-left: {depth > 0 ? '8px' : '0'}">
		{#if field.type === 'section'}
			<!-- ── Section divider ──────────────────────────────── -->
			<div class="gap-3 py-2 flex items-center">
				<span class="text-sm font-semibold" style="color: var(--text-primary)">{field.label}</span>
				<hr class="flex-1" style="border-color: var(--border)" />
			</div>
		{:else if field.type === 'group'}
			<!-- ── Group container ──────────────────────────────── -->
			<div
				class="rounded-xl p-4 gap-4 flex flex-col border"
				style="border-color: var(--border); background: var(--surface-alt)"
			>
				{#if field.label}
					<p class="text-sm font-semibold" style="color: var(--text-primary)">{field.label}</p>
				{/if}
				{ensureGroupData() && ''}
				{#each field.fields ?? [] as subField}
					<svelte:self field={subField} bind:data={data[field.key]} depth={depth + 1} />
				{/each}
			</div>
		{:else if field.type === 'repeater'}
			<!-- ── Repeater ─────────────────────────────────────── -->
			<div class="gap-2 flex flex-col">
				{#if field.label}
					<div class="flex items-center justify-between">
						<label class="text-sm font-medium" style="color: var(--text-primary)"
							>{field.label}</label
						>
						<button
							type="button"
							onclick={addRepeaterItem}
							class="text-xs font-medium transition hover:opacity-80"
							style="color: var(--brand)">+ Add Item</button
						>
					</div>
				{/if}
				{#each repeaterItems as _, i}
					<div
						class="rounded-xl p-4 gap-4 flex flex-col border"
						style="border-color: var(--border); background: var(--surface-alt)"
					>
						<div class="flex items-center justify-between">
							<span class="text-xs font-medium" style="color: var(--text-muted)">Item {i + 1}</span>
							<button
								type="button"
								onclick={() => removeRepeaterItem(i)}
								class="text-xs text-red-400 hover:text-red-600 transition">Remove</button
							>
						</div>
						{#each field.fields ?? [] as subField}
							<svelte:self field={subField} bind:data={repeaterItems[i]} depth={depth + 1} />
						{/each}
					</div>
				{/each}
				{#if repeaterItems.length === 0}
					<button
						type="button"
						onclick={addRepeaterItem}
						class="rounded-xl py-6 text-sm font-medium border-2 border-dashed transition"
						style="border-color: var(--border); color: var(--text-muted)"
						onmouseenter={(e) => {
							(e.currentTarget as HTMLElement).style.borderColor = 'var(--brand)';
							(e.currentTarget as HTMLElement).style.color = 'var(--brand)';
						}}
						onmouseleave={(e) => {
							(e.currentTarget as HTMLElement).style.borderColor = 'var(--border)';
							(e.currentTarget as HTMLElement).style.color = 'var(--text-muted)';
						}}>+ Add First Item</button
					>
				{/if}
			</div>
		{:else if field.type === 'file'}
			<!-- ── File upload ──────────────────────────────────── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<FileUpload bind:value={data[field.key]} isImage={false} accept="*" />
			</div>
		{:else if field.type === 'image'}
			<!-- ── Image upload ─────────────────────────────────── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<FileUpload bind:value={data[field.key]} isImage={true} accept="image/*" />
			</div>
		{:else if field.type === 'bool' || field.type === 'toggle'}
			<!-- ── Toggle switch ────────────────────────────────── -->
			<div
				class="rounded-lg px-3 py-2.5 flex items-center justify-between border"
				style="border-color: var(--border); background: var(--surface-alt)"
			>
				<div class="flex flex-col">
					{#if field.label}
						<span class="text-sm font-medium" style="color: var(--text-primary)">{field.label}</span
						>
					{/if}
					{#if field.description}
						<span class="text-xs" style="color: var(--text-muted)">{field.description}</span>
					{/if}
				</div>
				<button
					type="button"
					role="switch"
					aria-checked={!!data[field.key]}
					onclick={() => (data[field.key] = !data[field.key])}
					class="h-6 w-11 relative flex flex-shrink-0 rounded-full transition-colors duration-200"
					style="background: {data[field.key] ? 'var(--brand)' : 'var(--border)'}"
				>
					<span
						class="top-1 h-4 w-4 bg-white shadow absolute rounded-full transition-transform duration-200"
						style="left: {data[field.key] ? '1.375rem' : '0.25rem'}"
					></span>
				</button>
			</div>
		{:else if field.type === 'textarea'}
			<!-- ── Textarea ─────────────────────────────────────── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<textarea
					bind:value={data[field.key]}
					placeholder={field.placeholder}
					rows="4"
					class="rounded-lg px-3 py-2 text-sm w-full resize-y border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				></textarea>
				{@render fieldDescription()}
			</div>
		{:else if field.type === 'select'}
			<!-- ── Select dropdown ──────────────────────────────── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<select
					bind:value={data[field.key]}
					class="h-9 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				>
					<option value="">Select an option…</option>
					{#each field.options ?? [] as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
				{@render fieldDescription()}
			</div>
		{:else if field.type === 'radio'}
			<!-- ── Radio buttons ────────────────────────────────── -->
			<div class="gap-2 flex flex-col">
				{@render fieldLabel()}
				<div class="gap-2 flex flex-col">
					{#each field.options ?? [] as opt}
						<label class="gap-2 flex cursor-pointer items-center">
							<input
								type="radio"
								name={field.key}
								value={opt.value}
								checked={data[field.key] === opt.value}
								onchange={() => (data[field.key] = opt.value)}
								class="accent-[--brand]"
							/>
							<span class="text-sm" style="color: var(--text-primary)">{opt.label}</span>
						</label>
					{/each}
				</div>
			</div>
		{:else if field.type === 'checkbox'}
			<!-- ── Single checkbox ──────────────────────────────── -->
			<label class="gap-3 flex cursor-pointer items-center">
				<input
					type="checkbox"
					checked={!!data[field.key]}
					onchange={(e) => (data[field.key] = (e.target as HTMLInputElement).checked)}
					class="h-4 w-4 accent-[--brand]"
				/>
				<span class="text-sm font-medium" style="color: var(--text-primary)">{field.label}</span>
			</label>
		{:else if field.type === 'multiselect'}
			<!-- ── Multi-select checkboxes ───────────────────────── -->
			<div class="gap-2 flex flex-col">
				{@render fieldLabel()}
				{#each field.options ?? [] as opt}
					<label class="gap-2 flex cursor-pointer items-center">
						<input
							type="checkbox"
							checked={Array.isArray(data[field.key]) && data[field.key].includes(opt.value)}
							onchange={(e) => {
								const arr: unknown[] = Array.isArray(data[field.key]) ? [...data[field.key]] : [];
								if ((e.target as HTMLInputElement).checked) {
									data[field.key] = [...arr, opt.value];
								} else {
									data[field.key] = arr.filter((v) => v !== opt.value);
								}
							}}
							class="h-4 w-4 accent-[--brand]"
						/>
						<span class="text-sm" style="color: var(--text-primary)">{opt.label}</span>
					</label>
				{/each}
			</div>
		{:else if field.type === 'slider'}
			<!-- ── Range slider ─────────────────────────────────── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<div class="gap-3 flex items-center">
					<input
						type="range"
						bind:value={data[field.key]}
						min={field.validation?.min ?? 0}
						max={field.validation?.max ?? 100}
						class="flex-1 accent-[--brand]"
					/>
					<span class="w-10 text-sm text-right tabular-nums" style="color: var(--text-primary)">
						{data[field.key] ?? 0}
					</span>
				</div>
			</div>
		{:else if field.type === 'rating'}
			<!-- ── Star rating ──────────────────────────────────── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<div class="gap-1 flex">
					{#each [1, 2, 3, 4, 5] as star}
						<button
							type="button"
							onclick={() => (data[field.key] = star)}
							onmouseenter={() => (ratingHover = star)}
							onmouseleave={() => (ratingHover = 0)}
							class="text-2xl transition-colors"
							style="color: {star <= (ratingHover || data[field.key] || 0)
								? '#f59e0b'
								: 'var(--border)'}">★</button
						>
					{/each}
				</div>
			</div>
		{:else if field.type === 'colorpicker'}
			<!-- ── Color picker ─────────────────────────────────── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<div class="gap-3 flex items-center">
					<input
						type="color"
						bind:value={data[field.key]}
						class="h-9 w-16 rounded-lg px-1 cursor-pointer border"
						style="border-color: var(--border)"
					/>
					<input
						type="text"
						bind:value={data[field.key]}
						placeholder="#14b8a6"
						class="h-9 rounded-lg px-3 text-sm font-mono flex-1 border"
						style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
					/>
				</div>
			</div>
		{:else if field.type === 'datepicker'}
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<input
					type="date"
					bind:value={data[field.key]}
					class="h-9 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
			</div>
		{:else if field.type === 'timepicker'}
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<input
					type="time"
					bind:value={data[field.key]}
					class="h-9 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
			</div>
		{:else if field.type === 'datetimepicker'}
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<input
					type="datetime-local"
					bind:value={data[field.key]}
					class="h-9 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
			</div>
		{:else if field.type === 'relation' || field.type === 'autocomplete'}
			<!-- ── Relation selector ─────────────────────────────── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<select
					bind:value={data[field.key]}
					onfocus={loadRelation}
					class="h-9 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				>
					<option value="">
						{relationLoading ? 'Loading…' : 'Select related record…'}
					</option>
					{#each relationOptions as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
				{@render fieldDescription()}
			</div>
		{:else if field.type === 'number'}
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<input
					type="number"
					bind:value={data[field.key]}
					placeholder={field.placeholder}
					min={field.validation?.min}
					max={field.validation?.max}
					class="h-9 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
				{@render fieldDescription()}
			</div>
		{:else}
			<!-- ── Default: text-like (text, email, url, phone, password) ── -->
			<div class="gap-1.5 flex flex-col">
				{@render fieldLabel()}
				<input
					type={field.type === 'email'
						? 'email'
						: field.type === 'password'
							? 'password'
							: field.type === 'url'
								? 'url'
								: 'text'}
					bind:value={data[field.key]}
					placeholder={field.placeholder}
					required={field.required}
					class="h-9 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
				{@render fieldDescription()}
			</div>
		{/if}
	</div>
{/if}

{#snippet fieldLabel()}
	{#if field.label && field.type !== 'checkbox' && field.type !== 'bool' && field.type !== 'toggle'}
		<label class="text-sm font-medium" style="color: var(--text-primary)">
			{field.label}{#if field.required}<span class="text-red-500 ml-0.5">*</span>{/if}
		</label>
	{/if}
{/snippet}

{#snippet fieldDescription()}
	{#if field.description}
		<p class="text-xs" style="color: var(--text-muted)">{field.description}</p>
	{/if}
{/snippet}
