import { get, post, put, del } from './client';
import type { User } from '$lib/stores/auth';

export interface Role {
    id?: string;
    name: string;
    description?: string;
    permissions: Array<{
        collection_slug: string;
        actions: string[]; // "create", "read", "update", "delete"
    }>;
}

export const adminApi = {
    // Users
    getUsers: () => get<User[]>('/admin/users'),
    createUser: (data: any) => post<User>('/admin/users', data),
    deleteUser: (id: string) => del(`/admin/users/${id}`),

    // Roles
    getRoles: () => get<Role[]>('/admin/roles'),
    createRole: (role: Role) => post<Role>('/admin/roles', role),
    updateRole: (id: string, role: Role) => put<Role>(`/admin/roles/${id}`, role),
    deleteRole: (id: string) => del(`/admin/roles/${id}`),

    // Stats
    getStats: () => get<Record<string, number>>('/admin/stats')
};
