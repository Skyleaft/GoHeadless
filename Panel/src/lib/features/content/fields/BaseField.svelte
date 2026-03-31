<script lang="ts">
	import type { Field } from '$lib/types/collection';

	interface Props {
		field: Field;
		children: import('svelte').Snippet<[{ id: string }]>;
		showLabel?: boolean;
	}

	let { field, children, showLabel = true }: Props = $props();

	// Generate a unique ID for the input to link with the label
	const inputId = $derived(`field-${field.key}-${Math.random().toString(36).slice(2, 9)}`);

	// Some field types shouldn't show the standard label (e.g. sections or standalone checkboxes)
	const shouldShowStandardLabel = $derived(
		showLabel && 
		field.label && 
		!['checkbox', 'bool', 'toggle', 'section'].includes(field.type)
	);
</script>

<div class="field-container">
	{#if shouldShowStandardLabel}
		<label class="field-label" for={inputId}>
			{field.label}
			{#if field.required}
				<span class="required-indicator">*</span>
			{/if}
		</label>
	{/if}

	<div class="field-content">
		{@render children({ id: inputId })}
	</div>

	{#if field.description}
		<p class="field-description">
			{field.description}
		</p>
	{/if}
</div>

<style>
	.field-container {
		display: flex;
		flex-direction: column;
		gap: 0.375rem; /* gap-1.5 */
		width: 100%;
	}

	.field-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--text-primary);
		line-height: 1.25rem;
	}

	.required-indicator {
		color: #ef4444; /* red-500 */
		margin-left: 0.125rem;
	}

	.field-content {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.field-description {
		font-size: 0.75rem;
		color: var(--text-muted);
		line-height: 1rem;
	}
</style>
