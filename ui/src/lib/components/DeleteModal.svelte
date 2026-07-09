<script lang="ts">
	import type { FeatureFlag } from '$lib/types';

	let {
		flag,
		open,
		onClose,
		onDeleteConfirm
	}: {
		flag: FeatureFlag | null;
		open: boolean;
		onClose: () => void;
		onDeleteConfirm: (id: string) => Promise<void>;
	} = $props();

	let submitting = $state(false);

	async function confirm() {
		if (!flag) return;
		submitting = true;
		try {
			await onDeleteConfirm(flag.id);
		} finally {
			submitting = false;
		}
	}
</script>

{#if open && flag}
	<div class="overlay" onclick={(e) => { if (e.target === e.currentTarget) onClose(); }} onkeydown={(e) => { if (e.key === 'Escape') onClose(); }} role="dialog" tabindex="-1">
		<div class="modal">
			<h3>Delete Flag</h3>
			<p>Delete "{flag.name}"? This cannot be undone.</p>
			<div class="form-actions">
				<button class="btn btn-secondary" onclick={onClose} disabled={submitting}>Cancel</button>
				<button class="btn btn-danger" onclick={confirm} disabled={submitting}>
					{submitting ? 'Deleting...' : 'Delete'}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0,0,0,0.4);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 999;
	}
	.modal {
		background: #fff;
		border-radius: 10px;
		padding: 1.5rem;
		width: 90%;
		max-width: 480px;
		box-shadow: 0 10px 25px rgba(0,0,0,0.15);
	}
	h3 { font-size: 1.1rem; margin-bottom: 0.25rem; }
	p { color: #495057; font-size: 0.9rem; margin-bottom: 1rem; }
	.form-actions { display: flex; gap: 0.5rem; justify-content: flex-end; }
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
	.btn:disabled { opacity: 0.6; cursor: not-allowed; }
	.btn-secondary { background: #e9ecef; color: #495057; }
	.btn-secondary:hover:not(:disabled) { background: #dee2e6; }
	.btn-danger { background: #e03131; color: #fff; }
	.btn-danger:hover:not(:disabled) { background: #c92a2a; }
</style>