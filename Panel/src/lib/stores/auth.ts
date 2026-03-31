import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface User {
    id: string;
    username: string;
    email?: string;
    role_id: string;
    active_status: boolean;
    is_initial_admin: boolean;
}

interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
}

const initial: AuthState = {
    user: null,
    token: null,
    isAuthenticated: false
};

function createAuthStore() {
    // Load from localStorage if in browser
    const stored = browser ? localStorage.getItem('auth') : null;
    const data = stored ? JSON.parse(stored) : initial;

    const { subscribe, set, update } = writable<AuthState>(data);

    return {
        subscribe,
        login: (user: User, token: string) => {
            const state = { user, token, isAuthenticated: true };
            if (browser) localStorage.setItem('auth', JSON.stringify(state));
            set(state);
        },
        logout: () => {
            if (browser) localStorage.removeItem('auth');
            set(initial);
        }
    };
}

export const auth = createAuthStore();
