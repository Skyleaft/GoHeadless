<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getSetupStatus, initializeSystem, login } from '$lib/api/auth';
	import { auth } from '$lib/stores/auth';
	import { toast } from '$lib/stores/toast';
	import Button from '$lib/shared/Button.svelte';
	import Input from '$lib/shared/Input.svelte';

	let username = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let loading = $state(false);
	let checking = $state(true);

	onMount(async () => {
		try {
			const status = await getSetupStatus();
			if (!status.setup_required) {
				goto('/login');
			}
		} catch (err) {
			toast.error('Failed to connect to backend server');
		} finally {
			checking = false;
		}
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (password !== confirmPassword) {
			toast.error('Passwords do not match');
			return;
		}

		if (password.length < 8) {
			toast.error('Password must be at least 8 characters');
			return;
		}

		loading = true;
		try {
			await initializeSystem({ username, password });
			toast.success('System initialized successfully!');
			
			// Auto login
			const res = await login({ username, password });
			auth.login(res.user, res.token);
			
			goto('/');
		} catch (err: any) {
			toast.error(err.message || 'Failed to initialize system');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>System Setup | GoHeadless CMS</title>
</svelte:head>

<div class="fixed inset-0 flex items-center justify-center p-4 bg-slate-50 dark:bg-slate-950">
	<div class="w-full max-w-md">
		{#if checking}
			<div class="flex flex-col items-center gap-4 animate-fade-in text-center">
				<div class="h-10 w-10 rounded-full border-4 border-[var(--brand)] border-t-transparent animate-spin"></div>
				<p class="text-sm text-[var(--text-secondary)]">Checking system status...</p>
			</div>
		{:else}
			<div class="card animate-scale-in overflow-hidden shadow-2xl">
				<div class="bg-[var(--brand)] p-8 text-center text-white">
					<div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-white/20 text-3xl font-bold backdrop-blur-sm">
						G
					</div>
					<h1 class="text-2xl font-bold text-white mb-1">Welcome to GoHeadless</h1>
					<p class="text-white/80 text-sm">Create your superadmin account to get started</p>
				</div>

				<form onsubmit={handleSubmit} class="card-body flex flex-col gap-5">
					<Input
						label="Superadmin Username"
						placeholder="e.g. admin"
						required
						bind:value={username}
						disabled={loading}
					/>

					<Input
						label="Password"
						type="password"
						placeholder="••••••••"
						required
						bind:value={password}
						disabled={loading}
						description="Minimum 8 characters"
					/>

					<Input
						label="Confirm Password"
						type="password"
						placeholder="••••••••"
						required
						bind:value={confirmPassword}
						disabled={loading}
					/>

					<div class="mt-4">
						<Button type="submit" class="w-full" size="lg" {loading}>
							Initialize System
						</Button>
					</div>

					<p class="text-center text-xs text-slate-400 mt-2">
						This account will have full administrative access to all system features.
					</p>
				</form>
			</div>
		{/if}
	</div>
</div>

<style>
	.card {
		background: var(--surface);
		border-radius: var(--radius-xl);
	}
</style>
