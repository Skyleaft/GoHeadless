import type { UploadResponse } from '$lib/types/collection';

const UPLOAD_BASE = '/api/v1/upload';

export const uploadApi = {
	upload: async (file: File): Promise<UploadResponse> => {
		const fd = new FormData();
		fd.append('file', file);
		const res = await fetch(UPLOAD_BASE, { method: 'POST', body: fd });
		if (!res.ok) {
			const body = await res.json().catch(() => ({}));
			throw new Error(body.error ?? `Upload failed: ${res.status}`);
		}
		return res.json();
	},

	delete: async (path: string): Promise<void> => {
		const res = await fetch(`${UPLOAD_BASE}?path=${encodeURIComponent(path)}`, {
			method: 'DELETE'
		});
		if (!res.ok) throw new Error(`Delete failed: ${res.status}`);
	},

	replace: async (oldPath: string, file: File): Promise<UploadResponse> => {
		const fd = new FormData();
		fd.append('file', file);
		const res = await fetch(`${UPLOAD_BASE}?oldPath=${encodeURIComponent(oldPath)}`, {
			method: 'PUT',
			body: fd
		});
		if (!res.ok) {
			const body = await res.json().catch(() => ({}));
			throw new Error(body.error ?? `Replace failed: ${res.status}`);
		}
		return res.json();
	}
};
