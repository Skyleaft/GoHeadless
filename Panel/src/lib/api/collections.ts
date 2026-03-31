import type { Collection } from '$lib/types/collection';
import { get, post, del } from './client';

export const collectionsApi = {
	list: (): Promise<Collection[]> => get<Collection[]>('/collections'),

	get: (slug: string): Promise<Collection> => get<Collection>(`/collections/${slug}`),

	create: (payload: Omit<Collection, 'id'>): Promise<Collection> =>
		post<Collection>('/collections', payload),

	delete: (slug: string): Promise<void> => del<void>(`/collections/${slug}`)
};
