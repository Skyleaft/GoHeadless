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
			class="text-[10px] font-black uppercase tracking-[0.15em] opacity-50 px-1"
			style="color: var(--text-primary)"
		>
			{label}
			{#if required}<span class="text-[var(--brand)] ml-1">●</span>{/if}
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
		class="w-full h-12 rounded-2xl px-4 text-sm font-bold transition-all duration-200 border-2 focus:outline-none focus:ring-4 focus:ring-[var(--brand)]/10"
		style="
			background: var(--surface-alt);
			color: var(--text-primary);
			border-color: {error ? '#ef4444' : 'transparent'};
			box-shadow: {error ? 'var(--shadow-sm)' : 'inset 0 2px 4px rgba(0,0,0,0.02)'};
		"
	/>

	{#if description && !error}
		<p class="text-xs" style="color: var(--text-muted)">{description}</p>
	{/if}
	{#if error}
		<p class="text-xs text-red-500">{error}</p>
	{/if}
</div>
