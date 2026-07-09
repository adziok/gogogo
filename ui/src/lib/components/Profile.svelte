<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';

	const placeholderImage = `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='32' height='32' viewBox='0 0 32 32'%3E%3Ccircle cx='16' cy='16' r='16' fill='%234263eb'/%3E%3Cpath d='M16 14c2.4 0 4.36-1.96 4.36-4.36S18.4 5.27 16 5.27s-4.36 1.96-4.36 4.36S13.6 14 16 14zm0 2.18c-2.91 0-8.73 1.46-8.73 4.36v1.09c0 0.6 0.49 1.09 1.09 1.09h15.28c0.6 0 1.09-0.49 1.09-1.09v-1.09c0-2.91-5.82-4.36-8.73-4.36z' fill='%23fff'/%3E%3C/svg%3E`;
</script>

{#if auth.user}
	<div class="profile">
		<img
			src={auth.user.picture || placeholderImage}
			alt={auth.user.name || 'User'}
			class="avatar"
			onerror={(e) => { (e.target as HTMLImageElement).src = placeholderImage; }}
		/>
		<div class="info">
			<div class="name">{auth.user.name || 'Unknown'}</div>
			<div class="email">{auth.user.email || ''}</div>
		</div>
	</div>
{/if}

<style>
	.profile {
		display: flex;
		align-items: center;
		gap: 0.6rem;
	}
	.avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		object-fit: cover;
	}
	.info {
		line-height: 1.3;
	}
	.name {
		font-size: 0.85rem;
		font-weight: 600;
		color: #1a1a2e;
	}
	.email {
		font-size: 0.75rem;
		color: #868e96;
	}
</style>