import type { ContentRecord } from '$lib/types/collection';
import { get, post, put, del } from './client';

export const contentApi = {
	list: (slug: string): Promise<ContentRecord[]> => get<ContentRecord[]>(`/content/${slug}`),

	get: (slug: string, id: string): Promise<ContentRecord> =>
		get<ContentRecord>(`/content/${slug}/${id}`),

	create: (slug: string, data: ContentRecord): Promise<{ id: string }> =>
		post<{ id: string }>(`/content/${slug}`, data),

	update: (slug: string, id: string, data: ContentRecord): Promise<void> =>
		put<void>(`/content/${slug}/${id}`, data),

	delete: (slug: string, id: string): Promise<void> => del<void>(`/content/${slug}/${id}`)
};
