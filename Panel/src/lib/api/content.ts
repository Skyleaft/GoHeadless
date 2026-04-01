import type { ContentRecord } from '$lib/types/collection';
import { get, post, put, del } from './client';

export interface ListRecordsResult {
	data: ContentRecord[];
	total: number;
	page: number;
	limit: number;
}

function buildContentQuery(searchParams?: URLSearchParams): string {
	if (!searchParams || [...searchParams.keys()].length === 0) return '';
	return `?${searchParams.toString()}`;
}

export const contentApi = {
	list: (slug: string, searchParams?: URLSearchParams): Promise<ListRecordsResult> =>
		get<ListRecordsResult>(`/content/${slug}${buildContentQuery(searchParams)}`),

	get: (slug: string, id: string): Promise<ContentRecord> =>
		get<ContentRecord>(`/content/${slug}/${id}`),

	create: (slug: string, data: ContentRecord): Promise<{ id: string }> =>
		post<{ id: string }>(`/content/${slug}`, data),

	update: (slug: string, id: string, data: ContentRecord): Promise<void> =>
		put<void>(`/content/${slug}/${id}`, data),

	delete: (slug: string, id: string): Promise<void> => del<void>(`/content/${slug}/${id}`)
};
