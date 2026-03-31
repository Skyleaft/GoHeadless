import { collectionsApi } from '$lib/api/collections';
import { contentApi } from '$lib/api/content';
import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params }) => {
	try {
		const [collection, records] = await Promise.all([
			collectionsApi.get(params.slug),
			contentApi.list(params.slug)
		]);

		return {
			collection,
			records: records ?? [],
			slug: params.slug
		};
	} catch (err: any) {
		throw error(404, err.message || 'Collection not found');
	}
};
