import { get, post } from './client';
import type { User } from '$lib/stores/auth';

export interface LoginResponse {
    token: string;
    user: User;
}

export async function getSetupStatus(): Promise<{ setup_required: boolean }> {
    return get('/setup/status');
}

export async function initializeSystem(data: any): Promise<any> {
    return post('/setup', data);
}

export async function login(data: any): Promise<LoginResponse> {
    return post('/auth/login', data);
}
