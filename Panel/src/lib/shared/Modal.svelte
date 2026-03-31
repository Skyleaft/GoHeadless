<script lang="ts">
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		open?: boolean;
		title?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		onclose?: () => void;
		children: Snippet;
		footer?: Snippet;
	}

	let { open = $bindable(false), title, size = 'md', onclose, children, footer }: Props = $props();

	let dialog: HTMLDialogElement;

	const sizes = {
		sm: 'max-w-sm',
		md: 'max-w-lg',
		lg: 'max-w-2xl',
		xl: 'max-w-4xl'
	};

	$effect(() => {
		if (!dialog) return;
		if (open) {
			dialog.showModal();
		} else {
			dialog.close();
		}
	});

	function handleClose() {
		open = false;
		onclose?.();
	}

	function handleBackdrop(e: MouseEvent) {
		if (e.target === dialog) handleClose();
	}

	onMount(() => {
		dialog.addEventListener('cancel', handleClose);
		return () => dialog.removeEventListener('cancel', handleClose);
	});
</script>

<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<dialog
	bind:this={dialog}
	onclick={handleBackdrop}
	class="modal-dialog rounded-2xl p-0 shadow-xl m-auto w-[calc(100%-2rem)] {sizes[
		size
	]} animate-scale-in"
	style="
		background: var(--surface);
		border-color: var(--border);
		color: var(--text-primary);
	"
	class:hidden={!open}
>
	<!-- Header -->
	{#if title}
		<div class="px-8 py-6 flex items-center justify-between">
			<h3 class="text-base font-semibold" style="color: var(--text-primary)">{title}</h3>
			<button
				onclick={handleClose}
				class="w-8 h-8 rounded-lg text-lg flex items-center justify-center transition-colors hover:bg-[--surface-alt]"
				style="color: var(--text-muted)"
				aria-label="Close modal"
			>
				✕
			</button>
		</div>
	{/if}

	<!-- Body -->
	<div class="p-6">
		{@render children()}
	</div>

	<!-- Footer -->
	{#if footer}
		<div
			class="gap-3 px-8 py-6 rounded-b-2xl flex items-center justify-end"
			style="background: var(--surface-alt)"
		>
			{@render footer()}
		</div>
	{/if}
</dialog>

<style>
	.modal-dialog {
		position: fixed;
		inset: 0;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-dialog::backdrop {
		background: rgb(0 0 0 / 0.5);
		backdrop-filter: blur(2px);
	}

	@keyframes scale-in {
		from {
			opacity: 0;
			transform: scale(0.95) translateY(-8px);
		}
		to {
			opacity: 1;
			transform: scale(1) translateY(0);
		}
	}

	.animate-scale-in {
		animation: scale-in 0.18s cubic-bezier(0.16, 1, 0.3, 1) forwards;
	}
</style>
