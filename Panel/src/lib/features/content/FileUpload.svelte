<script lang="ts">
	import { uploadApi } from '$lib/api/upload';
	import { getFileUrl } from '$lib/utils/path';

	interface Props {
		value?: string; // stored file path
		accept?: string;
		isImage?: boolean;
		onchange?: (path: string) => void;
	}

	let { value = $bindable(''), accept = '*', isImage = false, onchange }: Props = $props();

	let uploading = $state(false);
	let error = $state('');
	let dragover = $state(false);
	let input: HTMLInputElement;

	async function upload(file: File) {
		uploading = true;
		error = '';
		try {
			const res = value ? await uploadApi.replace(value, file) : await uploadApi.upload(file);
			value = res.path;
			onchange?.(res.path);
		} catch (err: any) {
			error = err.message ?? 'Upload failed';
		} finally {
			uploading = false;
		}
	}

	async function handleChange(e: Event) {
		const file = (e.target as HTMLInputElement).files?.[0];
		if (file) await upload(file);
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragover = false;
		const file = e.dataTransfer?.files[0];
		if (file) await upload(file);
	}

	async function remove() {
		if (!value) return;
		try {
			await uploadApi.delete(value);
			value = '';
			onchange?.('');
		} catch {
			// silently ignore
		}
	}

	let previewSrc = $derived(value ? getFileUrl(value) : '');
</script>

<div class="gap-2 flex flex-col">
	{#if value && isImage}
		<!-- Image preview -->
		<div class="relative inline-flex">
			<img
				src={previewSrc}
				alt="Preview"
				class="h-32 w-32 rounded-xl border object-cover"
				style="border-color: var(--border)"
			/>
			<button
				type="button"
				onclick={remove}
				class="-right-2 -top-2 h-6 w-6 bg-red-500 text-white text-xs shadow-md hover:bg-red-600 absolute flex items-center justify-center rounded-full transition"
				title="Remove file">✕</button
			>
		</div>
	{:else if value}
		<!-- File path -->
		<div
			class="gap-3 rounded-lg px-3 py-2 flex items-center border"
			style="background: var(--surface-alt); border-color: var(--border)"
		>
			<span class="text-xl">📄</span>
			<span class="text-sm font-mono flex-1 truncate" style="color: var(--text-primary)"
				>{value}</span
			>
			<button
				type="button"
				onclick={remove}
				class="text-red-400 hover:text-red-600 text-sm transition">✕</button
			>
		</div>
	{/if}

	<!-- Drop zone -->
	{#if !value}
		<button
			type="button"
			onclick={() => input.click()}
			ondragover={(e) => {
				e.preventDefault();
				dragover = true;
			}}
			ondragleave={() => (dragover = false)}
			ondrop={handleDrop}
			class="gap-3 rounded-xl px-6 py-8 flex flex-col items-center justify-center border-2 border-dashed text-center transition-all"
			style="
				border-color: {dragover ? 'var(--brand)' : 'var(--border)'};
				background: {dragover ? 'var(--brand-light)' : 'var(--surface-alt)'};
			"
			disabled={uploading}
		>
			{#if uploading}
				<span
					class="h-7 w-7 rounded-full border-2 border-[--brand] border-t-transparent"
					style="animation: spin 0.7s linear infinite; display: block"
				></span>
				<span class="text-sm" style="color: var(--text-muted)">Uploading…</span>
			{:else}
				<span class="text-3xl">{isImage ? '🖼' : '📁'}</span>
				<div>
					<p class="text-sm font-medium" style="color: var(--text-primary)">
						Drop {isImage ? 'an image' : 'a file'} here, or click to browse
					</p>
					<p class="text-xs mt-1" style="color: var(--text-muted)">
						{isImage ? 'Images will be converted to WebP automatically' : 'Any file type accepted'}
					</p>
				</div>
			{/if}
		</button>
	{:else if !isImage}
		<button
			type="button"
			onclick={() => input.click()}
			class="text-sm transition hover:text-[--brand]"
			style="color: var(--text-muted)">Replace file</button
		>
	{/if}

	<input
		bind:this={input}
		type="file"
		{accept}
		onchange={handleChange}
		class="hidden"
		disabled={uploading}
	/>

	{#if error}
		<p class="text-xs text-red-500">{error}</p>
	{/if}
</div>

<style>
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
