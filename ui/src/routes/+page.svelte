<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import type { FeatureFlag, CreateFlagPayload, UpdateFlagPayload } from '$lib/types';
	import { fetchFlags, createFlag, updateFlag, deleteFlag } from '$lib/api';
	import FlagTable from '$lib/components/FlagTable.svelte';
	import CreateModal from '$lib/components/CreateModal.svelte';
	import EditModal from '$lib/components/EditModal.svelte';
	import DeleteModal from '$lib/components/DeleteModal.svelte';
	import Toast from '$lib/components/Toast.svelte';

	let flags = $state<FeatureFlag[]>([]);
	let loading = $state(true);
	let error = $state('');

	let showCreate = $state(false);
	let editTarget = $state<FeatureFlag | null>(null);
	let deleteTarget = $state<FeatureFlag | null>(null);

	let toastMsg = $state('');
	let toastType: 'success' | 'error' = $state('success');
	let toastVisible = $state(false);
	let toastTimer: ReturnType<typeof setTimeout> | null = null;

	function showToast(msg: string, type: 'success' | 'error' = 'success') {
		if (toastTimer) clearTimeout(toastTimer);
		toastMsg = msg;
		toastType = type;
		toastVisible = true;
		toastTimer = setTimeout(() => { toastVisible = false; }, 3000);
	}

	async function load() {
		loading = true;
		error = '';
		try {
			flags = await fetchFlags();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	async function handleCreate(payload: CreateFlagPayload) {
		await createFlag(payload);
		showCreate = false;
		showToast('Flag created successfully');
		load();
	}

	async function handleUpdate(id: string, payload: UpdateFlagPayload) {
		await updateFlag(id, payload);
		editTarget = null;
		showToast('Flag updated successfully');
		load();
	}

	async function handleDelete(id: string) {
		await deleteFlag(id);
		deleteTarget = null;
		showToast('Flag deleted successfully');
		load();
	}

	$effect(() => {
		if (auth.isAuthenticated && !auth.isLoading) {
			load();
		}
	});
</script>

<svelte:head>
	<title>Feature Flag Manager</title>
</svelte:head>

{#if auth.isLoading}
	<div class="state-msg">Loading...</div>
{:else if !auth.isAuthenticated}
	<div class="welcome-card">
		<h1>Feature Flag Manager</h1>
		<p>Please log in to manage your feature flags.</p>
	</div>
{:else}
	<div class="card">
		<div class="card-header">
			<h2>All Flags</h2>
			<button class="btn btn-primary" onclick={() => showCreate = true}>+ New Flag</button>
		</div>

		{#if loading}
			<div class="state-msg">Loading flags...</div>
		{:else if error}
			<div class="state-msg error">Failed to load flags: {error}</div>
		{:else if flags.length === 0}
			<div class="state-msg">No feature flags yet. Create one!</div>
		{:else}
			<FlagTable {flags} onEdit={(f) => editTarget = f} onDelete={(f) => deleteTarget = f} />
		{/if}
	</div>

	<CreateModal open={showCreate} onClose={() => showCreate = false} onCreate={handleCreate} />
	<EditModal flag={editTarget} open={editTarget !== null} onClose={() => editTarget = null} onUpdate={handleUpdate} />
	<DeleteModal flag={deleteTarget} open={deleteTarget !== null} onClose={() => deleteTarget = null} onDeleteConfirm={handleDelete} />
	<Toast message={toastMsg} type={toastType} visible={toastVisible} />
{/if}

<style>
	.card {
		background: #fff;
		border-radius: 8px;
		box-shadow: 0 1px 3px rgba(0,0,0,0.08);
		padding: 1.5rem;
	}
	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}
	.card-header h2 {
		font-size: 1.1rem;
		font-weight: 600;
		margin: 0;
	}
	.state-msg {
		text-align: center;
		padding: 2rem;
		color: #868e96;
		font-size: 0.95rem;
	}
	.state-msg.error {
		color: #c92a2a;
	}
	.welcome-card {
		background: #fff;
		border-radius: 8px;
		box-shadow: 0 1px 3px rgba(0,0,0,0.08);
		padding: 3rem 2rem;
		text-align: center;
	}
	.welcome-card h1 {
		font-size: 1.75rem;
		margin-bottom: 0.75rem;
		color: #1a1a2e;
	}
	.welcome-card p {
		color: #868e96;
		font-size: 1rem;
	}
	.btn {
		display: inline-flex;
		align-items: center;
		gap: 0.3rem;
		padding: 0.4rem 0.85rem;
		border: none;
		border-radius: 6px;
		font-size: 0.85rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.15s;
	}
	.btn-primary { background: #4263eb; color: #fff; }
	.btn-primary:hover { background: #3b5bdb; }
</style>