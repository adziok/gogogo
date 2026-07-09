<script lang="ts">
	import { onMount } from 'svelte';
	import favicon from '$lib/assets/favicon.svg';
	import { auth, initializeAuth } from '$lib/stores/auth.svelte';
	import LoginButton from '$lib/components/LoginButton.svelte';
	import LogoutButton from '$lib/components/LogoutButton.svelte';
	import Profile from '$lib/components/Profile.svelte';

	let { children } = $props();

	onMount(() => {
		initializeAuth();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<nav class="navbar">
	<div class="nav-inner">
		<a href="/" class="logo">Feature Flag Manager</a>
		<div class="nav-right">
			{#if auth.isLoading}
				<span class="nav-loading">Loading...</span>
			{:else if auth.isAuthenticated}
				<Profile />
				<LogoutButton />
			{:else if !auth.isLoading && !auth.error}
				<LoginButton />
			{/if}
		</div>
	</div>
</nav>

<main>
	{@render children()}
</main>

<style>
	:global(*) {
		box-sizing: border-box;
		margin: 0;
		padding: 0;
	}
	:global(body) {
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
		background: #f5f7fa;
		color: #1a1a2e;
		line-height: 1.5;
	}
	.navbar {
		background: #fff;
		border-bottom: 1px solid #e9ecef;
		padding: 0 2rem;
	}
	.nav-inner {
		max-width: 1100px;
		margin: 0 auto;
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 56px;
	}
	.logo {
		font-size: 1.1rem;
		font-weight: 700;
		color: #1a1a2e;
		text-decoration: none;
	}
	.logo:hover {
		color: #4263eb;
	}
	.nav-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}
	.nav-loading {
		font-size: 0.85rem;
		color: #868e96;
	}
	main {
		max-width: 1100px;
		margin: 0 auto;
		padding: 2rem;
	}
</style>