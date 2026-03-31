// Base API client — wraps fetch with error handling

const BASE_URL = '/api/v1';

export class ApiError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

async function handleResponse<T>(res: Response): Promise<T> {
	if (!res.ok) {
		let message = `HTTP ${res.status}`;
		try {
			const body = await res.json();
			message = body.error ?? body.message ?? message;
		} catch {
			// ignore parse error
		}
		throw new ApiError(res.status, message);
	}
	// 204 No Content
	if (res.status === 204) return undefined as T;
	return res.json() as Promise<T>;
}

export async function get<T>(path: string): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		headers: { Accept: 'application/json' }
	});
	return handleResponse<T>(res);
}

export async function post<T>(path: string, body: unknown): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', Accept: 'application/json' },
		body: JSON.stringify(body)
	});
	return handleResponse<T>(res);
}

export async function put<T>(path: string, body: unknown): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json', Accept: 'application/json' },
		body: JSON.stringify(body)
	});
	return handleResponse<T>(res);
}

export async function del<T>(path: string): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		method: 'DELETE',
		headers: { Accept: 'application/json' }
	});
	return handleResponse<T>(res);
}

export async function postForm<T>(path: string, formData: FormData): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		method: 'POST',
		body: formData
	});
	return handleResponse<T>(res);
}

export async function putForm<T>(path: string, formData: FormData): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		method: 'PUT',
		body: formData
	});
	return handleResponse<T>(res);
}
