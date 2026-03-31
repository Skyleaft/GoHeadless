<script lang="ts">
	import type { Field, FieldType, Option } from '$lib/types/collection';
	import { hasOptions, isStructuralType } from '$lib/types/collection';
	import FieldTypeSelector from './FieldTypeSelector.svelte';
	import Button from '$lib/shared/Button.svelte';

	interface Props {
		field: Field;
		depth?: number;
		onremove?: () => void;
	}

	let { field = $bindable(), depth = 0, onremove }: Props = $props();

	let expanded = $state(true);
	let showAdvanced = $state(false);

	function addOption() {
		field.options = [...(field.options ?? []), { label: '', value: '' }];
	}
	function removeOption(i: number) {
		field.options = field.options?.filter((_, idx) => idx !== i);
	}
	function addSubField() {
		field.fields = [
			...(field.fields ?? []),
			{ key: '', label: '', type: 'text', required: false, unique: false }
		];
	}
	function removeSubField(i: number) {
		field.fields = field.fields?.filter((_, idx) => idx !== i);
	}
</script>

<div
	class="rounded-xl border transition-all"
	style="
		border-color: var(--border);
		background: {depth > 0 ? 'var(--surface-alt)' : 'var(--surface)'};
		margin-left: {depth * 16}px;
	"
>
	<!-- Field Header -->
	<div class="flex items-center gap-3 px-4 py-3">
		<!-- Drag handle (visual only) -->
		<span class="cursor-grab text-lg select-none" style="color: var(--text-muted)">⠿</span>

		<!-- Key -->
		<input
			bind:value={field.key}
			placeholder="field_key"
			class="h-8 flex-1 rounded-lg border px-3 text-sm font-mono min-w-0"
			style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
		/>

		<!-- Label -->
		<input
			bind:value={field.label}
			placeholder="Label"
			class="h-8 flex-1 rounded-lg border px-3 text-sm min-w-0"
			style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
		/>

		<!-- Type selector -->
		<div class="w-44 flex-shrink-0">
			<FieldTypeSelector bind:value={field.type} />
		</div>

		<!-- Required toggle -->
		<label class="flex items-center gap-1.5 cursor-pointer flex-shrink-0" title="Required">
			<input type="checkbox" bind:checked={field.required} class="accent-[--brand] w-4 h-4" />
			<span class="text-xs" style="color: var(--text-secondary)">Req</span>
		</label>

		<!-- Expand -->
		<button
			type="button"
			onclick={() => (expanded = !expanded)}
			class="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg transition-all hover:bg-[--surface-alt]"
			style="color: var(--text-muted)"
			title="Expand field options"
		>
			{expanded ? '▲' : '▼'}
		</button>

		<!-- Remove -->
		{#if onremove}
			<button
				type="button"
				onclick={onremove}
				class="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-red-400 transition-all hover:bg-red-50 hover:text-red-600"
				title="Remove field"
			>✕</button>
		{/if}
	</div>

	<!-- Expanded options -->
	{#if expanded}
		<div
			class="grid grid-cols-1 gap-4 border-t px-4 py-4 sm:grid-cols-2"
			style="border-color: var(--border)"
		>
			<!-- Placeholder -->
			<div class="flex flex-col gap-1">
				<label class="text-xs font-medium" style="color: var(--text-secondary)">Placeholder</label>
				<input
					bind:value={field.placeholder}
					placeholder="Enter placeholder text"
					class="h-8 rounded-lg border px-3 text-sm w-full"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
			</div>

			<!-- Description -->
			<div class="flex flex-col gap-1">
				<label class="text-xs font-medium" style="color: var(--text-secondary)">Description</label>
				<input
					bind:value={field.description}
					placeholder="Helper text shown under the field"
					class="h-8 rounded-lg border px-3 text-sm w-full"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
			</div>

			<!-- Unique -->
			<div class="flex items-center gap-2">
				<input type="checkbox" bind:checked={field.unique} id="unique-{field.key}" class="accent-[--brand] w-4 h-4" />
				<label for="unique-{field.key}" class="text-sm cursor-pointer" style="color: var(--text-secondary)">
					Unique value
				</label>
			</div>

			<!-- Options (select/radio/multiselect/checkbox) -->
			{#if hasOptions(field.type)}
				<div class="col-span-full flex flex-col gap-2">
					<div class="flex items-center justify-between">
						<label class="text-xs font-medium" style="color: var(--text-secondary)">Options</label>
						<Button variant="ghost" size="sm" onclick={addOption}>+ Add Option</Button>
					</div>
					{#each field.options ?? [] as opt, i}
						<div class="flex gap-2">
							<input
								bind:value={opt.label}
								placeholder="Label"
								class="h-8 flex-1 rounded-lg border px-3 text-sm"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
							<input
								bind:value={opt.value}
								placeholder="Value"
								class="h-8 flex-1 rounded-lg border px-3 text-sm font-mono"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
							<button
								type="button"
								onclick={() => removeOption(i)}
								class="flex h-8 w-8 items-center justify-center rounded-lg text-red-400 hover:bg-red-50"
							>✕</button>
						</div>
					{/each}
				</div>
			{/if}

			<!-- Relation config -->
			{#if field.type === 'relation' || field.type === 'autocomplete'}
				<div class="col-span-full grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1">
						<label class="text-xs font-medium" style="color: var(--text-secondary)">Related Collection</label>
						<input
							bind:value={field.relation!.collection}
							placeholder="collection-slug"
							class="h-8 rounded-lg border px-3 text-sm font-mono w-full"
							style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
						/>
					</div>
					<div class="flex flex-col gap-1">
						<label class="text-xs font-medium" style="color: var(--text-secondary)">Display Field</label>
						<input
							bind:value={field.relation!.field}
							placeholder="field_key"
							class="h-8 rounded-lg border px-3 text-sm font-mono w-full"
							style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
						/>
					</div>
				</div>
			{/if}

			<!-- Conditional logic -->
			<div class="col-span-full">
				<button
					type="button"
					onclick={() => (showAdvanced = !showAdvanced)}
					class="text-xs font-medium transition-colors hover:text-[--brand]"
					style="color: var(--text-muted)"
				>
					{showAdvanced ? '▲' : '▼'} Conditional Logic & Validation
				</button>

				{#if showAdvanced}
					<div class="mt-3 grid grid-cols-2 gap-4 rounded-lg p-4" style="background: var(--surface-alt)">
						<!-- show_if -->
						<div class="flex flex-col gap-1">
							<label class="text-xs font-medium" style="color: var(--text-secondary)">Show if field key</label>
							<input
								value={field.logic?.show_if?.field ?? ''}
								oninput={(e) => {
									const v = (e.target as HTMLInputElement).value;
									field.logic = { show_if: { field: v, equals: field.logic?.show_if?.equals ?? '' } };
								}}
								placeholder="other_field_key"
								class="h-8 rounded-lg border px-3 text-sm w-full"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
						<div class="flex flex-col gap-1">
							<label class="text-xs font-medium" style="color: var(--text-secondary)">Equals value</label>
							<input
								value={field.logic?.show_if?.equals ?? ''}
								oninput={(e) => {
									const v = (e.target as HTMLInputElement).value;
									field.logic = { show_if: { field: field.logic?.show_if?.field ?? '', equals: v } };
								}}
								placeholder="true / myValue"
								class="h-8 rounded-lg border px-3 text-sm w-full"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
						<!-- Validation -->
						<div class="flex flex-col gap-1">
							<label class="text-xs font-medium" style="color: var(--text-secondary)">Min length</label>
							<input
								type="number"
								value={field.validation?.min_length ?? ''}
								oninput={(e) => {
									const v = parseInt((e.target as HTMLInputElement).value);
									field.validation = { ...field.validation, min_length: isNaN(v) ? undefined : v };
								}}
								class="h-8 rounded-lg border px-3 text-sm w-full"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
						<div class="flex flex-col gap-1">
							<label class="text-xs font-medium" style="color: var(--text-secondary)">Max length</label>
							<input
								type="number"
								value={field.validation?.max_length ?? ''}
								oninput={(e) => {
									const v = parseInt((e.target as HTMLInputElement).value);
									field.validation = { ...field.validation, max_length: isNaN(v) ? undefined : v };
								}}
								class="h-8 rounded-lg border px-3 text-sm w-full"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
						<div class="col-span-full flex flex-col gap-1">
							<label class="text-xs font-medium" style="color: var(--text-secondary)">Regex pattern</label>
							<input
								value={field.validation?.regex ?? ''}
								oninput={(e) => {
									field.validation = { ...field.validation, regex: (e.target as HTMLInputElement).value };
								}}
								placeholder="^[a-z]+$"
								class="h-8 rounded-lg border px-3 text-sm font-mono w-full"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
					</div>
				{/if}
			</div>
		</div>

		<!-- Sub-fields for group / repeater -->
		{#if isStructuralType(field.type) && field.type !== 'section'}
			<div class="border-t px-4 py-4" style="border-color: var(--border)">
				<div class="flex items-center justify-between pb-3">
					<span class="text-xs font-medium" style="color: var(--text-secondary)">
						Sub-fields ({field.fields?.length ?? 0})
					</span>
					<Button variant="ghost" size="sm" onclick={addSubField}>+ Add Sub-field</Button>
				</div>
				<div class="flex flex-col gap-2">
					{#each field.fields ?? [] as subField, i}
						<svelte:self
							bind:field={field.fields![i]}
							depth={depth + 1}
							onremove={() => removeSubField(i)}
						/>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>
