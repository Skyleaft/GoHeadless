<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import BaseField from './BaseField.svelte';

	interface Props {
		field: Field;
		data: any;
		depth?: number;
		// Recursion support
		renderer: any; 
	}

	let { field, data = $bindable(), depth = 0, renderer: FieldRenderer }: Props = $props();

	// Repeater items local state
	let repeaterItems = $state<Record<string, any>[]>(
		Array.isArray(data) ? [...(data as Record<string, any>[])] : []
	);

	$effect(() => {
		if (field.type === 'repeater' && data !== repeaterItems) {
			data = repeaterItems;
		}
	});

	function addRepeaterItem() {
		repeaterItems = [...repeaterItems, {}];
	}

	function removeRepeaterItem(i: number) {
		repeaterItems = repeaterItems.filter((_, idx) => idx !== i);
	}

	function ensureGroupData() {
		if (field.type === 'group' && (!data || typeof data !== 'object' || Array.isArray(data))) {
			data = {};
		}
	}

	$effect(() => {
		if (field.type === 'group') ensureGroupData();
	});

	// Generate a unique ID for the toggle/switch if needed
	const toggleId = $derived(field.key + '-' + Math.random().toString(36).substring(2, 9));
</script>

{#if field.type === 'section'}
	<div class="section-divider">
		<span class="section-label">{field.label}</span>
		<hr class="divider-line" />
	</div>

{:else if field.type === 'group'}
	<div class="group-container">
		{#if field.label}
			<p class="group-label">{field.label}</p>
		{/if}
		<div class="group-fields">
			{#each field.fields ?? [] as subField}
				<FieldRenderer field={subField} bind:data={data} depth={depth + 1} />
			{/each}
		</div>
	</div>

{:else if field.type === 'repeater'}
	<div class="repeater-container">
		{#if field.label}
			<div class="repeater-header">
				<label class="repeater-title">{field.label}</label>
				<button 
					type="button" 
					onclick={addRepeaterItem} 
					class="add-btn"
					title="Add new item"
					aria-label="Add new item to {field.label}"
				>+ Add Item</button>
			</div>
		{/if}

		<div class="repeater-items">
			{#each repeaterItems as _, i}
				<div class="repeater-item">
					<div class="item-header">
						<span class="item-index">Item {i + 1}</span>
						<button type="button" onclick={() => removeRepeaterItem(i)} class="remove-btn">Remove</button>
					</div>
					<div class="item-fields">
						{#each field.fields ?? [] as subField}
							<FieldRenderer field={subField} bind:data={repeaterItems[i]} depth={depth + 1} />
						{/each}
					</div>
				</div>
			{/each}
		</div>

		{#if repeaterItems.length === 0}
			<button type="button" onclick={addRepeaterItem} class="empty-state-btn">
				+ Add First Item
			</button>
		{/if}
	</div>

{:else if field.type === 'bool' || field.type === 'toggle'}
	<div class="toggle-container">
		<label class="toggle-info" for={toggleId}>
			{#if field.label}
				<span class="toggle-label">{field.label}</span>
			{/if}
			{#if field.description}
				<span class="toggle-desc">{field.description}</span>
			{/if}
		</label>
		<button
			id={toggleId}
			type="button"
			role="switch"
			aria-checked={!!data}
			onclick={() => (data = !data)}
			class="switch-btn"
			data-active={!!data}
		>
			<span class="switch-thumb" data-active={!!data}></span>
		</button>
	</div>
{/if}

<style>
	.section-divider {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem 0;
	}

	.section-label {
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--text-primary);
	}

	.divider-line {
		flex: 1;
		border: 0;
		border-top: 1px solid var(--border);
	}

	.group-container {
		border-radius: 0.75rem;
		padding: 1rem;
		border: 1px solid var(--border);
		background: var(--surface-alt);
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.group-label {
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--text-primary);
	}

	.group-fields {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.repeater-container {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.repeater-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.repeater-title {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--text-primary);
	}

	.add-btn {
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--brand);
		background: transparent;
		border: 0;
		cursor: pointer;
		transition: opacity 0.2s;
	}

	.add-btn:hover {
		opacity: 0.8;
	}

	.repeater-items {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.repeater-item {
		border-radius: 0.75rem;
		padding: 1rem;
		border: 1px solid var(--border);
		background: var(--surface-alt);
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.item-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.item-index {
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--text-muted);
	}

	.remove-btn {
		font-size: 0.75rem;
		color: #f87171;
		background: transparent;
		border: 0;
		cursor: pointer;
		transition: color 0.2s;
	}

	.remove-btn:hover {
		color: #ef4444;
	}

	.item-fields {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.empty-state-btn {
		border-radius: 0.75rem;
		padding: 1.5rem 0;
		font-size: 0.875rem;
		font-weight: 500;
		border: 2px dashed var(--border);
		background: transparent;
		color: var(--text-muted);
		cursor: pointer;
		transition: all 0.2s;
	}

	.empty-state-btn:hover {
		border-color: var(--brand);
		color: var(--brand);
	}

	.toggle-container {
		border-radius: 0.5rem;
		padding: 0.625rem 0.75rem;
		display: flex;
		align-items: center;
		justify-content: space-between;
		border: 1px solid var(--border);
		background: var(--surface-alt);
	}

	.toggle-info {
		display: flex;
		flex-direction: column;
	}

	.toggle-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--text-primary);
	}

	.toggle-desc {
		font-size: 0.75rem;
		color: var(--text-muted);
	}

	.switch-btn {
		width: 2.75rem;
		height: 1.5rem;
		border-radius: 9999px;
		position: relative;
		border: 0;
		cursor: pointer;
		transition: background-color 0.25s;
		background-color: var(--border);
	}

	.switch-btn[data-active=true] {
		background-color: var(--brand);
	}

	.switch-thumb {
		width: 1rem;
		height: 1rem;
		background-color: white;
		border-radius: 9999px;
		position: absolute;
		top: 0.25rem;
		left: 0.25rem;
		transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 1px 2px rgba(0,0,0,0.1);
	}

	.switch-thumb[data-active=true] {
		transform: translateX(1.25rem);
	}
</style>
