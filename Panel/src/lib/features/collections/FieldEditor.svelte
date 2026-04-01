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

	$effect(() => {
		if (field.searchable === undefined) field.searchable = true;
		if (field.internal === undefined) field.internal = false;
	});

	function addOption() {
		field.options = [...(field.options ?? []), { label: '', value: '' }];
	}
	function removeOption(i: number) {
		field.options = field.options?.filter((_, idx) => idx !== i);
	}
	function addSubField() {
		field.fields = [
			...(field.fields ?? []),
			{
				key: '',
				label: '',
				type: 'text',
				required: false,
				unique: false,
				is_array: false,
				array_config: { unique_items: false },
				searchable: true,
				internal: false
			}
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
	<div class="gap-3 px-4 py-3 flex items-center">
		<!-- Drag handle (visual only) -->
		<span class="text-lg cursor-grab select-none" style="color: var(--text-muted)">⠿</span>

		<!-- Key -->
		<input
			bind:value={field.key}
			placeholder="field_key"
			class="h-8 rounded-lg px-3 text-sm font-mono min-w-0 flex-1 border"
			style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
		/>

		<!-- Label -->
		<input
			bind:value={field.label}
			placeholder="Label"
			class="h-8 rounded-lg px-3 text-sm min-w-0 flex-1 border"
			style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
		/>

		<!-- Type selector -->
		<div class="w-44 flex-shrink-0">
			<FieldTypeSelector bind:value={field.type} />
		</div>

		<!-- Required toggle -->
		<label class="gap-1.5 flex flex-shrink-0 cursor-pointer items-center" title="Required">
			<input type="checkbox" bind:checked={field.required} class="w-4 h-4 accent-[--brand]" />
			<span class="text-xs" style="color: var(--text-secondary)">Req</span>
		</label>

		<!-- Expand -->
		<button
			type="button"
			onclick={() => (expanded = !expanded)}
			class="h-7 w-7 rounded-lg flex flex-shrink-0 items-center justify-center transition-all hover:bg-[--surface-alt]"
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
				class="h-7 w-7 rounded-lg text-red-400 hover:bg-red-50 hover:text-red-600 flex flex-shrink-0 items-center justify-center transition-all"
				title="Remove field">✕</button
			>
		{/if}
	</div>

	<!-- Expanded options -->
	{#if expanded}
		<div
			class="gap-4 px-4 py-4 sm:grid-cols-2 grid grid-cols-1 border-t"
			style="border-color: var(--border)"
		>
			<!-- Placeholder -->
			<div class="gap-1 flex flex-col">
				<label class="text-xs font-medium" style="color: var(--text-secondary)">Placeholder</label>
				<input
					bind:value={field.placeholder}
					placeholder="Enter placeholder text"
					class="h-8 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
			</div>

			<!-- Description -->
			<div class="gap-1 flex flex-col">
				<label class="text-xs font-medium" style="color: var(--text-secondary)">Description</label>
				<input
					bind:value={field.description}
					placeholder="Helper text shown under the field"
					class="h-8 rounded-lg px-3 text-sm w-full border"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				/>
			</div>

			<!-- Unique -->
			<div class="gap-2 flex items-center">
				<input
					type="checkbox"
					bind:checked={field.unique}
					id="unique-{field.key}"
					class="w-4 h-4 accent-[--brand]"
				/>
				<label
					for="unique-{field.key}"
					class="text-sm cursor-pointer"
					style="color: var(--text-secondary)"
				>
					Unique value
				</label>
			</div>

			<!-- Query Engine Flags -->
			<div class="gap-4 flex flex-wrap items-center">
				<div class="gap-2 flex items-center">
					<input
						type="checkbox"
						bind:checked={field.searchable}
						id="searchable-{field.key}"
						class="w-4 h-4 accent-[--brand]"
					/>
					<label
						for="searchable-{field.key}"
						class="text-sm cursor-pointer"
						style="color: var(--text-secondary)"
					>
						Searchable
					</label>
				</div>
				<div class="gap-2 flex items-center">
					<input
						type="checkbox"
						bind:checked={field.internal}
						id="internal-{field.key}"
						class="w-4 h-4 accent-[--brand]"
					/>
					<label
						for="internal-{field.key}"
						class="text-sm cursor-pointer"
						style="color: var(--text-secondary)"
					>
						Internal
					</label>
				</div>
			</div>

			<!-- Options (select/radio/multiselect/checkbox) -->
			{#if hasOptions(field.type)}
				<div class="gap-2 col-span-full flex flex-col">
					<div class="flex items-center justify-between">
						<label class="text-xs font-medium" style="color: var(--text-secondary)">Options</label>
						<Button variant="ghost" size="sm" onclick={addOption}>+ Add Option</Button>
					</div>
					{#each field.options ?? [] as opt, i}
						<div class="gap-2 flex">
							<input
								bind:value={opt.label}
								placeholder="Label"
								class="h-8 rounded-lg px-3 text-sm flex-1 border"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
							<input
								bind:value={opt.value}
								placeholder="Value"
								class="h-8 rounded-lg px-3 text-sm font-mono flex-1 border"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
							<button
								type="button"
								onclick={() => removeOption(i)}
								class="h-8 w-8 rounded-lg text-red-400 hover:bg-red-50 flex items-center justify-center"
								>✕</button
							>
						</div>
					{/each}
				</div>
			{/if}

			<!-- Relation config -->
			{#if field.type === 'relation' || field.type === 'autocomplete'}
				<div class="gap-4 col-span-full grid grid-cols-2">
					<div class="gap-1 flex flex-col">
						<label class="text-xs font-medium" style="color: var(--text-secondary)"
							>Related Collection</label
						>
						<input
							bind:value={field.relation!.collection}
							placeholder="collection-slug"
							class="h-8 rounded-lg px-3 text-sm font-mono w-full border"
							style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
						/>
					</div>
					<div class="gap-1 flex flex-col">
						<label class="text-xs font-medium" style="color: var(--text-secondary)"
							>Display Field</label
						>
						<input
							bind:value={field.relation!.field}
							placeholder="field_key"
							class="h-8 rounded-lg px-3 text-sm font-mono w-full border"
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
					<div
						class="mt-3 gap-4 rounded-lg p-4 grid grid-cols-2"
						style="background: var(--surface-alt)"
					>
						<!-- show_if -->
						<div class="gap-1 flex flex-col">
							<label class="text-xs font-medium" style="color: var(--text-secondary)"
								>Show if field key</label
							>
							<input
								value={field.logic?.show_if?.field ?? ''}
								oninput={(e) => {
									const v = (e.target as HTMLInputElement).value;
									field.logic = {
										show_if: { field: v, equals: field.logic?.show_if?.equals ?? '' }
									};
								}}
								placeholder="other_field_key"
								class="h-8 rounded-lg px-3 text-sm w-full border"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
						<div class="gap-1 flex flex-col">
							<label class="text-xs font-medium" style="color: var(--text-secondary)"
								>Equals value</label
							>
							<input
								value={field.logic?.show_if?.equals ?? ''}
								oninput={(e) => {
									const v = (e.target as HTMLInputElement).value;
									field.logic = {
										show_if: { field: field.logic?.show_if?.field ?? '', equals: v }
									};
								}}
								placeholder="true / myValue"
								class="h-8 rounded-lg px-3 text-sm w-full border"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
						<!-- Validation -->
						<div class="gap-1 flex flex-col">
							<label class="text-xs font-medium" style="color: var(--text-secondary)"
								>Min length</label
							>
							<input
								type="number"
								value={field.validation?.min_length ?? ''}
								oninput={(e) => {
									const v = parseInt((e.target as HTMLInputElement).value);
									field.validation = { ...field.validation, min_length: isNaN(v) ? undefined : v };
								}}
								class="h-8 rounded-lg px-3 text-sm w-full border"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
						<div class="gap-1 flex flex-col">
							<label class="text-xs font-medium" style="color: var(--text-secondary)"
								>Max length</label
							>
							<input
								type="number"
								value={field.validation?.max_length ?? ''}
								oninput={(e) => {
									const v = parseInt((e.target as HTMLInputElement).value);
									field.validation = { ...field.validation, max_length: isNaN(v) ? undefined : v };
								}}
								class="h-8 rounded-lg px-3 text-sm w-full border"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
						<div
							class="gap-4 pt-2 mt-2 col-span-full flex flex-col border-t border-[var(--border)]"
						>
							<label class="gap-3 group flex cursor-pointer items-center">
								<div
									class="h-5 w-5 rounded-lg flex items-center justify-center border-2 transition-all {field.is_array
										? 'border-[var(--brand)] bg-[var(--brand)]'
										: 'border-[var(--border)] group-hover:border-[var(--brand)]/50'}"
								>
									<input type="checkbox" bind:checked={field.is_array} class="sr-only" />
									{#if field.is_array}
										<span class="text-white text-xs font-black">✓</span>
									{/if}
								</div>
								<div class="flex flex-col">
									<span class="text-sm font-black tracking-tight" style="color: var(--text-primary)"
										>Allow Multiple Values (Array)</span
									>
									<span class="font-bold tracking-widest text-[10px] uppercase opacity-40"
										>Store a list of items instead of a single value</span
									>
								</div>
							</label>

							{#if field.is_array}
								<div class="gap-4 animate-fade-in pl-8 grid grid-cols-3">
									<div class="gap-1 flex flex-col">
										<label class="font-black tracking-widest text-[10px] uppercase opacity-40"
											>Min Items</label
										>
										<input
											type="number"
											bind:value={field.array_config!.min_items}
											class="h-8 rounded-lg px-3 text-sm w-full border border-[var(--border)] bg-[var(--surface)]"
										/>
									</div>
									<div class="gap-1 flex flex-col">
										<label class="font-black tracking-widest text-[10px] uppercase opacity-40"
											>Max Items</label
										>
										<input
											type="number"
											bind:value={field.array_config!.max_items}
											class="h-8 rounded-lg px-3 text-sm w-full border border-[var(--border)] bg-[var(--surface)]"
										/>
									</div>
									<div class="gap-2 pt-4 flex items-center">
										<input
											type="checkbox"
											bind:checked={field.array_config!.unique_items}
											id="unique-items-{field.key}"
											class="w-4 h-4 accent-[--brand]"
										/>
										<label
											for="unique-items-{field.key}"
											class="font-black tracking-widest cursor-pointer text-[10px] uppercase opacity-40"
											>Unique Items</label
										>
									</div>
								</div>
							{/if}
						</div>

						<div
							class="gap-1 pt-4 mt-2 col-span-full flex flex-col border-t border-[var(--border)]"
						>
							<label class="text-xs font-medium" style="color: var(--text-secondary)"
								>Regex pattern</label
							>
							<input
								value={field.validation?.regex ?? ''}
								oninput={(e) => {
									field.validation = {
										...field.validation,
										regex: (e.target as HTMLInputElement).value
									};
								}}
								placeholder="^[a-z]+$"
								class="h-8 rounded-lg px-3 text-sm font-mono w-full border"
								style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
							/>
						</div>
					</div>
				{/if}
			</div>
		</div>

		<!-- Sub-fields for group / repeater -->
		{#if isStructuralType(field.type) && field.type !== 'section'}
			<div class="px-4 py-4 border-t" style="border-color: var(--border)">
				<div class="pb-3 flex items-center justify-between">
					<span class="text-xs font-medium" style="color: var(--text-secondary)">
						Sub-fields ({field.fields?.length ?? 0})
					</span>
					<Button variant="ghost" size="sm" onclick={addSubField}>+ Add Sub-field</Button>
				</div>
				<div class="gap-2 flex flex-col">
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
