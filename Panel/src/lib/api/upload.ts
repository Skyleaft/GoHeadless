import type { UploadResponse } from '$lib/types/collection';
import { postForm, del, putForm } from './client';

export const uploadApi = {
	upload: async (file: File): Promise<UploadResponse> => {
		const fd = new FormData();
		fd.append('file', file);
		return postForm<UploadResponse>('/upload', fd);
	},

	delete: async (path: string): Promise<void> => {
		return del<void>(`/upload?path=${encodeURIComponent(path)}`);
	},

	replace: async (oldPath: string, file: File): Promise<UploadResponse> => {
		const fd = new FormData();
		fd.append('file', file);
		return putForm<UploadResponse>(`/upload?oldPath=${encodeURIComponent(oldPath)}`, fd);
	}
};
