<script lang="ts">
	import { logout } from '$lib/stores/auth.svelte';

	async function handleLogout() {
		await logout();
	}
</script>

<button onclick={handleLogout} class="btn btn-logout">Log Out</button>

<style>
	.btn {
		padding: 0.45rem 1.2rem;
		font-size: 0.85rem;
		font-weight: 500;
		border-radius: 6px;
		border: none;
		cursor: pointer;
		transition: background 0.15s;
	}
	.btn-logout {
		background: #e03131;
		color: #fff;
	}
	.btn-logout:hover {
		background: #c92a2a;
	}
</style>