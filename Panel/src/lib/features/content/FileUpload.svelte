<script lang="ts">
	import { uploadApi } from '$lib/api/upload';
	import { getFileUrl } from '$lib/utils/path';

	// ─── Single-file mode ──────────────────────────────────────
	interface SingleProps {
		multiple?: false;
		value?: string;
		accept?: string;
		isImage?: boolean;
		onchange?: (path: string) => void;
	}

	// ─── Multi-file mode ───────────────────────────────────────
	interface MultiProps {
		multiple: true;
		value?: string[];
		accept?: string;
		isImage?: boolean;
		onchange?: (paths: string[]) => void;
	}

	type Props = SingleProps | MultiProps;

	let {
		multiple = false,
		value = $bindable(multiple ? [] : ''),
		accept = '*',
		isImage = false,
		onchange
	}: Props = $props();

	let uploading = $state(false);
	let error = $state('');
	let dragover = $state(false);
	let input: HTMLInputElement;

	// ── Single-file helpers ─────────────────────────────────────
	async function uploadSingle(file: File) {
		uploading = true;
		error = '';
		try {
			const res =
				typeof value === 'string' && value
					? await uploadApi.replace(value as string, file)
					: await uploadApi.upload(file);
			(value as any) = res.path;
			(onchange as any)?.(res.path);
		} catch (err: any) {
			error = err.message ?? 'Upload failed';
		} finally {
			uploading = false;
		}
	}

	async function removeSingle() {
		if (!value) return;
		try {
			await uploadApi.delete(value as string);
		} catch {
			// silently ignore
		}
		(value as any) = '';
		(onchange as any)?.('');
	}

	// ── Multi-file helpers ──────────────────────────────────────
	async function uploadMultiple(files: File[]) {
		if (!files.length) return;
		uploading = true;
		error = '';
		try {
			const res = await uploadApi.uploadMultiple(files);
			(value as any) = [...((value as string[]) ?? []), ...res.paths];
			(onchange as any)?.(value as string[]);
		} catch (err: any) {
			error = err.message ?? 'Upload failed';
		} finally {
			uploading = false;
		}
	}

	async function removeMultiItem(path: string) {
		try {
			await uploadApi.delete(path);
		} catch {
			// silently ignore
		}
		(value as any) = ((value as string[]) ?? []).filter((p) => p !== path);
		(onchange as any)?.(value as string[]);
	}

	// ── Shared event handlers ───────────────────────────────────
	async function handleChange(e: Event) {
		const files = Array.from((e.target as HTMLInputElement).files ?? []);
		if (!files.length) return;
		if (multiple) {
			await uploadMultiple(files);
		} else {
			await uploadSingle(files[0]);
		}
		// reset input so the same file can be re-selected
		(e.target as HTMLInputElement).value = '';
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragover = false;
		const files = Array.from(e.dataTransfer?.files ?? []);
		if (!files.length) return;
		if (multiple) {
			await uploadMultiple(files);
		} else {
			await uploadSingle(files[0]);
		}
	}

	// ── Derived ─────────────────────────────────────────────────
	const singleValue = $derived(multiple ? '' : (value as string));
	const multiValues = $derived(multiple ? (value as string[]) ?? [] : []);
	const hasFiles = $derived(multiple ? multiValues.length > 0 : !!singleValue);
</script>

<div class="upload-root">
	{#if !multiple}
		<!-- ═══════════════════ SINGLE MODE ═══════════════════ -->
		{#if singleValue && isImage}
			<div class="preview-single">
				<img src={getFileUrl(singleValue)} alt="Preview" class="img-preview" />
				<button type="button" onclick={removeSingle} class="remove-btn remove-btn--abs" title="Remove">✕</button>
			</div>
		{:else if singleValue}
			<div class="file-pill">
				<span class="file-icon">📄</span>
				<span class="file-path">{singleValue}</span>
				<button type="button" onclick={removeSingle} class="remove-inline" title="Remove">✕</button>
			</div>
		{/if}

		{#if !singleValue}
			<button
				type="button"
				onclick={() => input.click()}
				ondragover={(e) => { e.preventDefault(); dragover = true; }}
				ondragleave={() => (dragover = false)}
				ondrop={handleDrop}
				class="dropzone"
				class:dragover
				disabled={uploading}
			>
				{#if uploading}
					<span class="spinner"></span>
					<span class="hint">Uploading…</span>
				{:else}
					<span class="drop-icon">{isImage ? '🖼' : '📁'}</span>
					<p class="drop-label">Drop {isImage ? 'an image' : 'a file'} here, or <u>click to browse</u></p>
					<p class="drop-hint">{isImage ? 'Converted to WebP automatically' : 'Any file type accepted'}</p>
				{/if}
			</button>
		{:else if !isImage}
			<button type="button" onclick={() => input.click()} class="replace-btn">Replace file</button>
		{/if}

	{:else}
		<!-- ═══════════════════ MULTI MODE ════════════════════ -->
		{#if isImage && multiValues.length > 0}
			<!-- Image grid -->
			<div class="img-grid">
				{#each multiValues as path (path)}
					<div class="img-grid-item">
						<img src={getFileUrl(path)} alt="Preview" class="img-thumb" />
						<button
							type="button"
							onclick={() => removeMultiItem(path)}
							class="remove-btn remove-btn--abs"
							title="Remove"
						>✕</button>
					</div>
				{/each}
				<!-- Add-more tile -->
				<button
					type="button"
					onclick={() => input.click()}
					class="img-grid-add"
					disabled={uploading}
					title="Add more images"
				>
					{#if uploading}
						<span class="spinner spinner--sm"></span>
					{:else}
						<span class="add-icon">+</span>
					{/if}
				</button>
			</div>
		{:else if !isImage && multiValues.length > 0}
			<!-- File list -->
			<div class="file-list">
				{#each multiValues as path (path)}
					<div class="file-pill">
						<span class="file-icon">📄</span>
						<span class="file-path">{path}</span>
						<button
							type="button"
							onclick={() => removeMultiItem(path)}
							class="remove-inline"
							title="Remove"
						>✕</button>
					</div>
				{/each}
			</div>
		{/if}

		<!-- Drop zone: always visible in multi mode (to add more) -->
		<button
			type="button"
			onclick={() => input.click()}
			ondragover={(e) => { e.preventDefault(); dragover = true; }}
			ondragleave={() => (dragover = false)}
			ondrop={handleDrop}
			class="dropzone"
			class:dragover
			class:dropzone--compact={hasFiles}
			disabled={uploading}
		>
			{#if uploading}
				<span class="spinner"></span>
				<span class="hint">Uploading…</span>
			{:else if hasFiles}
				<span class="drop-icon" style="font-size:1.25rem">{isImage ? '🖼' : '📁'}</span>
				<p class="drop-label" style="font-size:0.8rem">Add more {isImage ? 'images' : 'files'}</p>
			{:else}
				<span class="drop-icon">{isImage ? '🖼' : '📁'}</span>
				<p class="drop-label">Drop {isImage ? 'images' : 'files'} here, or <u>click to browse</u></p>
				<p class="drop-hint">You can select multiple files at once</p>
			{/if}
		</button>
	{/if}

	<!-- Hidden file input -->
	<input
		bind:this={input}
		type="file"
		{accept}
		multiple={multiple || undefined}
		onchange={handleChange}
		class="hidden-input"
		disabled={uploading}
	/>

	{#if error}
		<p class="error-msg">{error}</p>
	{/if}
</div>

<style>
	.upload-root {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		width: 100%;
	}

	/* ── Drop zone ── */
	.dropzone {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 2rem 1.5rem;
		border-radius: 0.75rem;
		border: 2px dashed var(--border);
		background: var(--surface-alt);
		cursor: pointer;
		text-align: center;
		transition: border-color 0.15s, background 0.15s;
		width: 100%;
	}
	.dropzone.dragover {
		border-color: var(--brand);
		background: var(--brand-light, color-mix(in srgb, var(--brand) 10%, transparent));
	}
	.dropzone:disabled {
		cursor: not-allowed;
		opacity: 0.7;
	}
	.dropzone--compact {
		padding: 0.75rem 1rem;
		flex-direction: row;
		gap: 0.5rem;
	}

	.drop-icon { font-size: 2rem; }
	.drop-label { font-size: 0.875rem; font-weight: 500; color: var(--text-primary); margin: 0; }
	.drop-label u { text-underline-offset: 2px; }
	.drop-hint { font-size: 0.75rem; color: var(--text-muted); margin: 0; }

	/* ── Single image preview ── */
	.preview-single {
		position: relative;
		display: inline-flex;
	}
	.img-preview {
		width: 8rem;
		height: 8rem;
		object-fit: cover;
		border-radius: 0.75rem;
		border: 1px solid var(--border);
	}

	/* ── Image grid (multi) ── */
	.img-grid {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}
	.img-grid-item {
		position: relative;
	}
	.img-thumb {
		width: 6rem;
		height: 6rem;
		object-fit: cover;
		border-radius: 0.625rem;
		border: 1px solid var(--border);
		display: block;
	}
	.img-grid-add {
		width: 6rem;
		height: 6rem;
		border-radius: 0.625rem;
		border: 2px dashed var(--border);
		background: var(--surface-alt);
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: border-color 0.15s, background 0.15s;
		flex-shrink: 0;
	}
	.img-grid-add:hover:not(:disabled) {
		border-color: var(--brand);
		background: var(--brand-light, color-mix(in srgb, var(--brand) 10%, transparent));
	}
	.add-icon {
		font-size: 1.75rem;
		color: var(--text-muted);
		line-height: 1;
	}

	/* ── File list / pill ── */
	.file-list {
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
	}
	.file-pill {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.4rem 0.75rem;
		border-radius: 0.5rem;
		border: 1px solid var(--border);
		background: var(--surface-alt);
	}
	.file-icon { font-size: 1.1rem; flex-shrink: 0; }
	.file-path {
		flex: 1;
		font-size: 0.8125rem;
		font-family: monospace;
		color: var(--text-primary);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	/* ── Remove buttons ── */
	.remove-btn {
		width: 1.4rem;
		height: 1.4rem;
		border-radius: 50%;
		background: #ef4444;
		color: #fff;
		font-size: 0.7rem;
		font-weight: 700;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		border: none;
		transition: background 0.15s;
		flex-shrink: 0;
	}
	.remove-btn:hover { background: #dc2626; }
	.remove-btn--abs {
		position: absolute;
		top: -0.5rem;
		right: -0.5rem;
	}
	.remove-inline {
		background: none;
		border: none;
		color: #f87171;
		font-size: 0.875rem;
		cursor: pointer;
		padding: 0;
		transition: color 0.15s;
	}
	.remove-inline:hover { color: #ef4444; }

	/* ── Replace btn ── */
	.replace-btn {
		background: none;
		border: none;
		font-size: 0.875rem;
		color: var(--text-muted);
		cursor: pointer;
		transition: color 0.15s;
		align-self: flex-start;
		padding: 0;
	}
	.replace-btn:hover { color: var(--brand); }

	/* ── Spinner ── */
	.spinner {
		display: block;
		width: 1.75rem;
		height: 1.75rem;
		border-radius: 50%;
		border: 2px solid var(--brand, #14b8a6);
		border-top-color: transparent;
		animation: spin 0.7s linear infinite;
	}
	.spinner--sm {
		width: 1.25rem;
		height: 1.25rem;
	}
	.hint { font-size: 0.875rem; color: var(--text-muted); }

	.hidden-input { display: none; }

	.error-msg { font-size: 0.75rem; color: #ef4444; margin: 0; }

	@keyframes spin { to { transform: rotate(360deg); } }
</style>
