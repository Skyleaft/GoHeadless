import type { Field } from '$lib/types/collection';

/** Leaf fields with dotted paths for nested groups (matches backend flatten). */
export function flattenFieldsForFilter(fields: Field[], prefix = ''): Field[] {
	const out: Field[] = [];
	for (const f of fields) {
		if (['section', 'tabs', 'grid'].includes(f.type)) {
			out.push(...flattenFieldsForFilter(f.fields ?? [], prefix));
		} else if (f.type === 'group') {
			const p = prefix ? `${prefix}.${f.key}` : f.key;
			out.push(...flattenFieldsForFilter(f.fields ?? [], p));
		} else if (f.type === 'repeater') {
			continue;
		} else {
			out.push({
				...f,
				key: prefix ? `${prefix}.${f.key}` : f.key
			});
		}
	}
	return out;
}
