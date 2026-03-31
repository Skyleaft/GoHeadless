<script lang="ts">
	interface Props {
		id?: string;
		label?: string;
		description?: string;
		error?: string;
		type?: string;
		value?: string | number;
		placeholder?: string;
		required?: boolean;
		disabled?: boolean;
		class?: string;
		onchange?: (e: Event) => void;
		oninput?: (e: Event) => void;
	}

	let {
		id,
		label,
		description,
		error,
		type = 'text',
		value = $bindable(''),
		placeholder,
		required = false,
		disabled = false,
		class: cls = '',
		onchange,
		oninput
	}: Props = $props();

	const inputId = id ?? `input-${Math.random().toString(36).slice(2)}`;
</script>

<div class="flex flex-col gap-1.5 {cls}">
	{#if label}
		<label
			for={inputId}
			class="text-sm font-medium"
			style="color: var(--text-primary)"
		>
			{label}
			{#if required}<span class="text-red-500 ml-0.5">*</span>{/if}
		</label>
	{/if}

	<input
		{id}
		{type}
		{placeholder}
		{required}
		{disabled}
		bind:value
		{onchange}
		{oninput}
		class="w-full rounded-lg border px-3 py-2 text-sm transition-all duration-150"
		style="
			background: var(--surface);
			color: var(--text-primary);
			border-color: {error ? '#ef4444' : 'var(--border)'};
		"
	/>

	{#if description && !error}
		<p class="text-xs" style="color: var(--text-muted)">{description}</p>
	{/if}
	{#if error}
		<p class="text-xs text-red-500">{error}</p>
	{/if}
</div>
