import type { Collection } from '$lib/types/collection';
import { get, post, put, del } from './client';

export const collectionsApi = {
	list: (): Promise<Collection[]> => get<Collection[]>('/collections'),

	get: (slug: string): Promise<Collection> => get<Collection>(`/collections/${slug}`),

	create: (payload: Omit<Collection, 'id'>): Promise<Collection> =>
		post<Collection>('/collections', payload),

	update: (slug: string, payload: Omit<Collection, 'id'>): Promise<Collection> =>
		put<Collection>(`/collections/${slug}`, payload),

	delete: (slug: string): Promise<void> => del<void>(`/collections/${slug}`)
};
