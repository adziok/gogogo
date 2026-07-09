<script lang="ts">
	let {
		message,
		type = 'success',
		visible
	}: {
		message: string;
		type?: 'success' | 'error';
		visible: boolean;
	} = $props();
</script>

{#if visible}
	<div class="toast" class:success={type === 'success'} class:error={type === 'error'}>
		{message}
	</div>
{/if}

<style>
	.toast {
		position: fixed;
		top: 1rem;
		right: 1rem;
		padding: 0.75rem 1.25rem;
		border-radius: 8px;
		color: #fff;
		font-weight: 500;
		font-size: 0.9rem;
		z-index: 1000;
	}
	.success { background: #2b8a3e; }
	.error { background: #c92a2a; }
</style>