<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import ArrayWrapper from './ArrayWrapper.svelte';
	import TagInput from '$lib/shared/TagInput.svelte';

	// Sub-renderers
	import TextField from './fields/TextField.svelte';
	import TextAreaField from './fields/TextAreaField.svelte';
	import NumberField from './fields/NumberField.svelte';
	import SelectField from './fields/SelectField.svelte';
	import DateTimeField from './fields/DateTimeField.svelte';
	import RelationField from './fields/RelationField.svelte';
	import StructuralField from './fields/StructuralField.svelte';
	import FileUpload from './FileUpload.svelte';
	import BaseField from './fields/BaseField.svelte';
	import Self from './FieldRenderer.svelte';

	interface Props {
		field: Field;
		data: Record<string, any>;
		depth?: number;
	}

	let { field, data = $bindable(), depth = 0 }: Props = $props();

	// Visibility logic
	let visible = $derived(
		field.logic?.show_if
			? String(data[field.logic.show_if.field]) === String(field.logic.show_if.equals)
			: true
	);

	// Ensure array data for array fields
	$effect(() => {
		if (field.is_array && !Array.isArray(data[field.key])) {
			data[field.key] = [];
		}
	});

	// Components mapping helpers
	const isText = (t: string) => ['text', 'email', 'password', 'url', 'phone'].includes(t);
	const isNumber = (t: string) => ['number', 'slider', 'rating'].includes(t);
	const isSelect = (t: string) => ['select', 'radio', 'checkbox', 'multiselect'].includes(t);
	const isDateTime = (t: string) => ['datepicker', 'timepicker', 'datetimepicker'].includes(t);
	const isRelation = (t: string) => ['relation', 'autocomplete'].includes(t);
	const isStructural = (t: string) => ['section', 'group', 'repeater', 'bool', 'toggle'].includes(t);
</script>

{#if visible}
	<div class="renderer-root" style="margin-left: {depth > 0 ? '0.5rem' : '0'}">
		{#if field.is_array && field.type === 'text'}
			<BaseField {field}>
				<TagInput bind:value={data[field.key]} placeholder={field.placeholder || "Add item..."} />
			</BaseField>

		{:else if field.is_array}
			<BaseField {field}>
				<ArrayWrapper bind:value={data[field.key]} {field}>
					{#each data[field.key] as _, i}
						<div class="array-item-wrapper">
							<Self 
								field={{...field, is_array: false, label: ''}} 
								bind:data={data[field.key][i]} 
								depth={depth + 1} 
							/>
						</div>
					{/each}
				</ArrayWrapper>
			</BaseField>

		{:else if isText(field.type)}
			<TextField {field} bind:data={data[field.key]} />

		{:else if field.type === 'textarea'}
			<TextAreaField {field} bind:data={data[field.key]} />

		{:else if isNumber(field.type)}
			<NumberField {field} bind:data={data[field.key]} />

		{:else if isSelect(field.type)}
			<SelectField {field} bind:data={data[field.key]} />

		{:else if isDateTime(field.type)}
			<DateTimeField {field} bind:data={data[field.key]} />

		{:else if isRelation(field.type)}
			<RelationField {field} bind:data={data[field.key]} />

		{:else if isStructural(field.type)}
			<StructuralField {field} bind:data={data[field.key]} {depth} renderer={Self} />

		{:else if field.type === 'file' || field.type === 'image'}
			<BaseField {field}>
				<FileUpload 
					bind:value={data[field.key]} 
					isImage={field.type === 'image'} 
					accept={field.type === 'image' ? 'image/*' : '*'} 
				/>
			</BaseField>

		{:else if field.type === 'colorpicker'}
			<BaseField {field}>
				<div class="color-picker-container">
					<input type="color" bind:value={data[field.key]} class="color-input" />
					<input type="text" bind:value={data[field.key]} placeholder="#14b8a6" class="hex-input" />
				</div>
			</BaseField>
		{/if}
	</div>
{/if}

<style>
	.renderer-root {
		width: 100%;
	}

	.array-item-wrapper {
		padding: 0.5rem;
		border-bottom: 1px solid var(--border);
	}

	.array-item-wrapper:last-child {
		border-bottom: none;
	}

	.color-picker-container {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.color-input {
		height: 2.25rem;
		width: 4rem;
		border-radius: 0.5rem;
		padding: 0.25rem;
		cursor: pointer;
		border: 1px solid var(--border);
		background: var(--surface);
	}

	.hex-input {
		height: 2.25rem;
		border-radius: 0.5rem;
		padding: 0 0.75rem;
		font-size: 0.875rem;
		font-family: monospace;
		flex: 1;
		border: 1px solid var(--border);
		background: var(--surface);
		color: var(--text-primary);
	}
</style>
