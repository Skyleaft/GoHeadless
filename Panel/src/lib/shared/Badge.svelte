<script lang="ts">
	import type { FieldType } from '$lib/types/collection';

	interface Props {
		type: FieldType | string;
		class?: string;
	}

	let { type, class: cls = '' }: Props = $props();

	const colorMap: Record<string, string> = {
		text: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
		textarea: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
		number: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
		email: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
		password: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
		select: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
		radio: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
		checkbox: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
		multiselect: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
		datepicker: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
		timepicker: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
		datetimepicker: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
		file: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
		image: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
		bool: 'bg-teal-100 text-teal-700 dark:bg-teal-900/30 dark:text-teal-400',
		toggle: 'bg-teal-100 text-teal-700 dark:bg-teal-900/30 dark:text-teal-400',
		group: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400',
		repeater: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400',
		relation: 'bg-pink-100 text-pink-700 dark:bg-pink-900/30 dark:text-pink-400',
		autocomplete: 'bg-pink-100 text-pink-700 dark:bg-pink-900/30 dark:text-pink-400'
	};

	let color = $derived(
		colorMap[type] ?? 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'
	);
</script>

<span
	class="inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium {color} {cls}"
>
	{type}
</span>
