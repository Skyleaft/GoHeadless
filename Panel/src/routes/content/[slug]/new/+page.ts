import { collectionsApi } from '$lib/api/collections';
import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params }) => {
	try {
		const collection = await collectionsApi.get(params.slug);
		return {
			collection,
			slug: params.slug
		};
	} catch (err: any) {
		throw error(404, err.message || 'Collection not found');
	}
};
