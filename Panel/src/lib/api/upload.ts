import type { UploadResponse, UploadMultipleResponse } from '$lib/types/collection';
import { postForm, del, putForm } from './client';

export const uploadApi = {
	upload: async (file: File): Promise<UploadResponse> => {
		const fd = new FormData();
		fd.append('file', file);
		return postForm<UploadResponse>('/upload', fd);
	},

	uploadMultiple: async (files: File[]): Promise<UploadMultipleResponse> => {
		const fd = new FormData();
		for (const file of files) fd.append('files', file);
		return postForm<UploadMultipleResponse>('/upload/multiple', fd);
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
