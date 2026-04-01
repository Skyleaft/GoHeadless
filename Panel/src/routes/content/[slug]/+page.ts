import { collectionsApi } from '$lib/api/collections';
import { contentApi } from '$lib/api/content';
import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, url }) => {
	try {
		const [collection, result] = await Promise.all([
			collectionsApi.get(params.slug),
			contentApi.list(params.slug, url.searchParams)
		]);

		return {
			collection,
			records: result.data ?? [],
			total: result.total ?? 0,
			page: result.page ?? 1,
			limit: result.limit ?? 10,
			slug: params.slug
		};
	} catch (err: any) {
		throw error(404, err.message || 'Collection not found');
	}
};
