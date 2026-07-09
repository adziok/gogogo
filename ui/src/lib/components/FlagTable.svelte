<script lang="ts">
	import type { FeatureFlag } from '$lib/types';

	let {
		flags,
		onEdit,
		onDelete
	}: {
		flags: FeatureFlag[];
		onEdit: (flag: FeatureFlag) => void;
		onDelete: (flag: FeatureFlag) => void;
	} = $props();
</script>

<table>
	<thead>
		<tr>
			<th>Name</th>
			<th>Description</th>
			<th>Enabled</th>
			<th>Created</th>
			<th>Actions</th>
		</tr>
	</thead>
	<tbody>
		{#each flags as flag (flag.id)}
			<tr>
				<td class="name">{flag.name}</td>
				<td>{flag.description}</td>
				<td>
					<span class="badge" class:enabled={flag.enabled} class:disabled={!flag.enabled}>
						{flag.enabled ? 'Enabled' : 'Disabled'}
					</span>
				</td>
				<td class="date">{new Date(flag.created_at).toLocaleDateString()}</td>
				<td>
					<div class="actions">
						<button class="btn btn-sm btn-secondary" onclick={() => onEdit(flag)}>Edit</button>
						<button class="btn btn-sm btn-danger" onclick={() => onDelete(flag)}>Delete</button>
					</div>
				</td>
			</tr>
		{/each}
	</tbody>
</table>

<style>
	table {
		width: 100%;
		border-collapse: collapse;
	}
	th, td {
		text-align: left;
		padding: 0.6rem 0.5rem;
		border-bottom: 1px solid #e9ecef;
		font-size: 0.9rem;
	}
	th {
		font-weight: 600;
		color: #495057;
		white-space: nowrap;
	}
	tr:hover { background: #f8f9fa; }
	.name { font-weight: 600; }
	.date { font-size: 0.8rem; color: #868e96; white-space: nowrap; }
	.badge {
		display: inline-block;
		padding: 0.15rem 0.5rem;
		border-radius: 12px;
		font-size: 0.75rem;
		font-weight: 600;
	}
	.enabled { background: #d3f9d8; color: #2b8a3e; }
	.disabled { background: #ffe3e3; color: #c92a2a; }
	.actions { display: flex; gap: 0.4rem; }
	.btn {
		display: inline-flex;
		align-items: center;
		padding: 0.25rem 0.5rem;
		border: none;
		border-radius: 6px;
		font-size: 0.8rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.15s;
	}
	.btn-secondary { background: #e9ecef; color: #495057; }
	.btn-secondary:hover { background: #dee2e6; }
	.btn-danger { background: #e03131; color: #fff; }
	.btn-danger:hover { background: #c92a2a; }
</style>