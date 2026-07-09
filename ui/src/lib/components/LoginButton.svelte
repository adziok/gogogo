<script lang="ts">
	import { login } from '$lib/stores/auth.svelte';

	async function handleLogin() {
		await login();
	}
</script>

<button onclick={handleLogin} class="btn btn-login">Log In</button>

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
	.btn-login {
		background: #4263eb;
		color: #fff;
	}
	.btn-login:hover {
		background: #3b5bdb;
	}
</style>