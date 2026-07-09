<script lang="ts">
	import type { CreateFlagPayload } from '$lib/types';

	let {
		open,
		onClose,
		onCreate
	}: {
		open: boolean;
		onClose: () => void;
		onCreate: (payload: CreateFlagPayload) => Promise<void>;
	} = $props();

	let name = $state('');
	let description = $state('');
	let enabled = $state(false);
	let submitting = $state(false);

	async function submit(e: Event) {
		e.preventDefault();
		submitting = true;
		try {
			await onCreate({ name: name.trim(), description: description.trim(), enabled });
		} finally {
			submitting = false;
		}
	}

	function resetForm() {
		name = '';
		description = '';
		enabled = false;
	}
</script>

{#if open}
	<div class="overlay" onclick={(e) => { if (e.target === e.currentTarget) onClose(); }} onkeydown={(e) => { if (e.key === 'Escape') onClose(); }} role="dialog" tabindex="-1">
		<div class="modal">
			<h3>New Feature Flag</h3>
			<p>Create a new feature flag.</p>
			<form onsubmit={submit}>
				<div class="field">
					<label for="new-name">Name *</label>
					<input
						id="new-name"
						type="text"
						bind:value={name}
						required
						minlength={3}
						pattern="[a-zA-Z0-9]+"
						title="Alphanumeric, min 3 characters"
						disabled={submitting}
					/>
				</div>
				<div class="field">
					<label for="new-desc">Description *</label>
					<textarea
						id="new-desc"
						bind:value={description}
						required
						maxlength={100}
						disabled={submitting}
					></textarea>
				</div>
				<div class="field checkbox">
					<input id="new-enabled" type="checkbox" bind:checked={enabled} disabled={submitting} />
					<label for="new-enabled">Enabled</label>
				</div>
				<div class="form-actions">
					<button type="button" class="btn btn-secondary" onclick={onClose} disabled={submitting}>Cancel</button>
					<button type="submit" class="btn btn-primary" disabled={submitting}>
						{submitting ? 'Creating...' : 'Create'}
					</button>
				</div>
			</form>
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
	p { color: #868e96; font-size: 0.9rem; margin-bottom: 1rem; }
	.field { margin-bottom: 1rem; }
	.field label { display: block; font-size: 0.85rem; font-weight: 500; margin-bottom: 0.3rem; color: #495057; }
	.field input[type="text"], .field textarea {
		width: 100%;
		padding: 0.5rem 0.7rem;
		border: 1px solid #ced4da;
		border-radius: 6px;
		font-size: 0.9rem;
		font-family: inherit;
		transition: border-color 0.15s;
	}
	.field input[type="text"]:focus, .field textarea:focus {
		outline: none;
		border-color: #4263eb;
		box-shadow: 0 0 0 3px rgba(66,99,235,0.15);
	}
	.field textarea { resize: vertical; min-height: 60px; }
	.checkbox { display: flex; align-items: center; gap: 0.5rem; }
	.checkbox label { margin-bottom: 0; }
	.checkbox input { width: 1.1rem; height: 1.1rem; cursor: pointer; }
	.form-actions { display: flex; gap: 0.5rem; justify-content: flex-end; margin-top: 1rem; }
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
	.btn-primary { background: #4263eb; color: #fff; }
	.btn-primary:hover:not(:disabled) { background: #3b5bdb; }
	.btn-secondary { background: #e9ecef; color: #495057; }
	.btn-secondary:hover:not(:disabled) { background: #dee2e6; }
</style>