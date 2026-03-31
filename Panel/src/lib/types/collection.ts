// ============================================================
// TypeScript interfaces mirroring Go domain models
// ============================================================

export type FieldType =
	// Basic Input
	| 'text'
	| 'textarea'
	| 'number'
	| 'email'
	| 'password'
	| 'phone'
	| 'url'
	// Selection
	| 'select'
	| 'radio'
	| 'checkbox'
	| 'multiselect'
	// Date & Time
	| 'datepicker'
	| 'timepicker'
	| 'datetimepicker'
	// File & Media
	| 'file'
	| 'image'
	// Advanced Input
	| 'bool'
	| 'toggle'
	| 'slider'
	| 'rating'
	| 'colorpicker'
	// Layout & Structure
	| 'section'
	| 'group'
	| 'tabs'
	| 'grid'
	// Repeater
	| 'repeater'
	// Relational
	| 'relation'
	| 'autocomplete'
	// Actions
	| 'submit'
	| 'reset'
	| 'action';

export interface ValidationRules {
	min?: number;
	max?: number;
	min_length?: number;
	max_length?: number;
	regex?: string;
}

export interface ArrayConfig {
	min_items?: number;
	max_items?: number;
	unique_items: boolean;
}

export interface Condition {
	field: string;
	equals: unknown;
}

export interface ConditionalLogic {
	show_if?: Condition;
}

export interface Option {
	label: string;
	value: unknown;
}

export interface RelationConfig {
	collection: string;
	field: string;
}

export interface Field {
	key: string;
	label: string;
	type: FieldType;
	placeholder?: string;
	default_value?: unknown;
	required: boolean;
	unique: boolean;
	description?: string;
	options?: Option[];
	validation?: ValidationRules;
	logic?: ConditionalLogic;
	fields?: Field[]; // recursive: group / repeater sub-fields
	relation?: RelationConfig;
	computed_by?: string;

	// Array Configuration
	is_array: boolean;
	array_config?: ArrayConfig;

	props?: Record<string, unknown>;
}

export interface CRUDPolicy {
	create: string[];
	read: string[];
	update: string[];
	delete: string[];
}

export interface AccessControl {
	is_public: boolean;
	allowed_roles?: string[];
	crud_policy?: CRUDPolicy;
}

export interface Collection {
	id?: string;
	name: string;
	slug: string;
	description: string;
	fields: Field[];
	access?: AccessControl;
}

// Record is a dynamic map from the content API
export type ContentRecord = Record<string, unknown> & { _id?: string; id?: string };

// Upload API response
export interface UploadResponse {
	path: string;
}

// ── Field type metadata for the UI ───────────────────────────

export interface FieldTypeGroup {
	label: string;
	types: { value: FieldType; label: string; icon: string }[];
}

export const FIELD_TYPE_GROUPS: FieldTypeGroup[] = [
	{
		label: 'Basic Input',
		types: [
			{ value: 'text', label: 'Text', icon: 'T' },
			{ value: 'textarea', label: 'Textarea', icon: '¶' },
			{ value: 'number', label: 'Number', icon: '#' },
			{ value: 'email', label: 'Email', icon: '@' },
			{ value: 'password', label: 'Password', icon: '🔒' },
			{ value: 'phone', label: 'Phone', icon: '☎' },
			{ value: 'url', label: 'URL', icon: '🔗' }
		]
	},
	{
		label: 'Selection',
		types: [
			{ value: 'select', label: 'Select', icon: '▾' },
			{ value: 'radio', label: 'Radio', icon: '◎' },
			{ value: 'checkbox', label: 'Checkbox', icon: '☑' },
			{ value: 'multiselect', label: 'Multi-Select', icon: '☰' }
		]
	},
	{
		label: 'Date & Time',
		types: [
			{ value: 'datepicker', label: 'Date', icon: '📅' },
			{ value: 'timepicker', label: 'Time', icon: '🕐' },
			{ value: 'datetimepicker', label: 'DateTime', icon: '📆' }
		]
	},
	{
		label: 'File & Media',
		types: [
			{ value: 'file', label: 'File', icon: '📄' },
			{ value: 'image', label: 'Image', icon: '🖼' }
		]
	},
	{
		label: 'Advanced',
		types: [
			{ value: 'bool', label: 'Boolean', icon: '⊙' },
			{ value: 'toggle', label: 'Toggle', icon: '⏻' },
			{ value: 'slider', label: 'Slider', icon: '⎓' },
			{ value: 'rating', label: 'Rating', icon: '★' },
			{ value: 'colorpicker', label: 'Color', icon: '🎨' }
		]
	},
	{
		label: 'Layout',
		types: [
			{ value: 'section', label: 'Section', icon: '▬' },
			{ value: 'group', label: 'Group', icon: '⊞' },
			{ value: 'tabs', label: 'Tabs', icon: '⊟' },
			{ value: 'grid', label: 'Grid', icon: '⊞' },
			{ value: 'repeater', label: 'Repeater', icon: '⟳' }
		]
	},
	{
		label: 'Relational',
		types: [
			{ value: 'relation', label: 'Relation', icon: '↔' },
			{ value: 'autocomplete', label: 'Autocomplete', icon: '⌕' }
		]
	}
];

export function getFieldTypeLabel(type: FieldType): string {
	for (const group of FIELD_TYPE_GROUPS) {
		const found = group.types.find((t) => t.value === type);
		if (found) return found.label;
	}
	return type;
}

export function isStructuralType(type: FieldType): boolean {
	return ['group', 'repeater', 'tabs', 'grid', 'section'].includes(type);
}

export function hasOptions(type: FieldType): boolean {
	return ['select', 'radio', 'multiselect', 'checkbox'].includes(type);
}

export function isFileType(type: FieldType): boolean {
	return ['file', 'image'].includes(type);
}
