<script lang="ts">
	import type { Snippet } from 'svelte';
	import Modal from './Modal.svelte';
	import Button from './Button.svelte';

	interface Props {
		open?: boolean;
		title?: string;
		message?: string;
		confirmLabel?: string;
		cancelLabel?: string;
		loading?: boolean;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		onconfirm?: () => void;
		oncancel?: () => void;
		children?: Snippet;
	}

	let {
		open = $bindable(false),
		title = 'Are you sure?',
		message = 'This action cannot be undone.',
		confirmLabel = 'Confirm',
		cancelLabel = 'Cancel',
		loading = false,
		size = 'sm', // width is now dynamic via parameter, default to `sm`
		onconfirm,
		oncancel,
		children
	}: Props = $props();

	function cancel() {
		open = false;
		oncancel?.();
	}
	function confirm() {
		onconfirm?.();
	}
</script>

<Modal bind:open {size} onclose={cancel}>
	<!-- Body of the dialog -->
	<div class="gap-5 sm:flex-row sm:items-start sm:text-left flex flex-col text-center">
		<div
			class="h-14 w-14 rounded-2xl shadow-sm sm:mx-0 sm:h-12 sm:w-12 mx-auto flex flex-shrink-0 items-center justify-center border"
			style="background: var(--surface-alt); border-color: var(--border);"
		>
			<!-- Modern alert triangle icon -->
			<svg
				class="h-6 w-6 text-red-500"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="2"
				stroke="currentColor"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
				/>
			</svg>
		</div>
		<div class="sm:mt-0.5 mt-2 flex-1">
			<h3 class="text-xl font-bold tracking-tight" style="color: var(--text-primary)">
				{title}
			</h3>
			<p class="mt-2 text-sm leading-relaxed" style="color: var(--text-secondary)">{message}</p>
			{#if children}
				<div class="mt-4">
					{@render children()}
				</div>
			{/if}
		</div>
	</div>

	{#snippet footer()}
		<Button variant="ghost" onclick={cancel}>{cancelLabel}</Button>
		<Button variant="danger" {loading} onclick={confirm}>{confirmLabel}</Button>
	{/snippet}
</Modal>
