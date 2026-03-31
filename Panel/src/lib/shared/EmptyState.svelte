<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title?: string;
		description?: string;
		icon?: string;
		children?: Snippet;
	}
	let { title = 'Nothing here yet', description = '', icon = '📭', children }: Props = $props();
</script>

<div class="flex flex-col items-center justify-center gap-4 py-16 text-center">
	<div
		class="flex h-20 w-20 items-center justify-center rounded-2xl text-4xl"
		style="background: var(--surface-alt)"
	>
		{icon}
	</div>
	<div class="flex flex-col gap-1">
		<h3 class="text-base font-semibold" style="color: var(--text-primary)">{title}</h3>
		{#if description}
			<p class="text-sm" style="color: var(--text-muted)">{description}</p>
		{/if}
	</div>
	{#if children}
		<div>{@render children()}</div>
	{/if}
</div>
