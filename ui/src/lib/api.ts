import type { FeatureFlag, CreateFlagPayload, UpdateFlagPayload } from './types';
import { getToken } from '$lib/stores/auth.svelte';

async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
	const token = await getToken();
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...options.headers as Record<string, string>,
	};
	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}
	const res = await fetch(path, { ...options, headers });
	if (!res.ok) {
		let msg = `Request failed (${res.status})`;
		try { const body = await res.json(); msg = body.error || msg; } catch { /* ignore */ }
		throw new Error(msg);
	}
	if (res.status === 204) return null as T;
	return res.json();
}

export function fetchFlags(): Promise<FeatureFlag[]> {
	return api('/feature-flag');
}

export function createFlag(payload: CreateFlagPayload): Promise<FeatureFlag> {
	return api('/feature-flag', { method: 'POST', body: JSON.stringify(payload) });
}

export function updateFlag(id: string, payload: UpdateFlagPayload): Promise<void> {
	return api(`/feature-flag/${id}`, { method: 'PUT', body: JSON.stringify(payload) });
}

export function deleteFlag(id: string): Promise<void> {
	return api(`/feature-flag/${id}`, { method: 'DELETE' });
}