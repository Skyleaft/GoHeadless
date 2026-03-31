<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getSetupStatus, login } from '$lib/api/auth';
	import { auth } from '$lib/stores/auth';
	import { toast } from '$lib/stores/toast';
	import Button from '$lib/shared/Button.svelte';
	import Input from '$lib/shared/Input.svelte';

	let username = $state('');
	let password = $state('');
	let loading = $state(false);
	let checking = $state(true);

	onMount(async () => {
		try {
			const status = await getSetupStatus();
			if (status.setup_required) {
				goto('/setup');
			}
		} catch (err) {
			// ignore, wait for manual login
		} finally {
			checking = false;
		}
	});

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		loading = true;
		try {
			const res = await login({ username, password });
			auth.login(res.user, res.token);
			toast.success(`Welcome back, ${res.user.username}!`);
			goto('/');
		} catch (err: any) {
			toast.error(err.message || 'Login failed. Please check your credentials.');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Sign In | GoHeadless CMS</title>
</svelte:head>

<div class="inset-0 p-4 bg-slate-50 dark:bg-slate-950 fixed flex items-center justify-center">
	<div class="max-w-md w-full">
		{#if checking}
			<div class="gap-4 animate-fade-in flex flex-col items-center text-center">
				<div
					class="h-10 w-10 animate-spin rounded-full border-4 border-[--brand] border-t-transparent"
				></div>
				<p class="text-sm text-[--text-secondary]">Connecting to server...</p>
			</div>
		{:else}
			<div class="card animate-scale-in shadow-xl overflow-hidden">
				<div class="bg-slate-100 dark:bg-slate-900/50 p-8 border-b border-[--border] text-center">
					<div
						class="mb-4 h-14 w-14 rounded-xl text-white text-2xl font-bold shadow-lg animate-fade-in mx-auto flex items-center justify-center bg-[--brand] shadow-[--brand]/20"
					>
						G
					</div>
					<h1 class="text-xl font-bold mb-1">Sign In to Dashboard</h1>
					<p class="text-sm text-[--text-secondary]">
						Enter your credentials to manage your content
					</p>
				</div>

				<form onsubmit={handleSubmit} class="card-body gap-6 p-8 flex flex-col">
					<Input
						label="Username"
						placeholder="admin"
						required
						bind:value={username}
						disabled={loading}
					/>

					<div class="gap-1.5 flex flex-col">
						<div class="flex items-center justify-between">
							<label class="text-sm font-medium" style="color: var(--text-primary)" for="password">
								Password
							</label>
							<button type="button" class="text-xs text-[--brand] hover:underline">Forgot password?</button>
						</div>
						<Input
							id="password"
							type="password"
							placeholder="••••••••"
							required
							bind:value={password}
							disabled={loading}
						/>
					</div>

					<div class="gap-2 flex items-center">
						<input
							type="checkbox"
							id="remember"
							class="h-4 w-4 rounded border-slate-300 text-[--brand] focus:ring-[--brand]"
						/>
						<label for="remember" class="text-sm cursor-pointer text-[--text-secondary]"
							>Remember me for 30 days</label
						>
					</div>

					<div class="pt-2">
						<Button type="submit" class="h-11 w-full" size="lg" {loading}>Authenticate</Button>
					</div>
				</form>

				<div
					class="bg-slate-50 dark:bg-slate-900/30 px-8 py-4 border-t border-[--border] text-center"
				>
					<p class="text-xs text-[--text-muted]">Secured by AES-256 JWT Authentication</p>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.card {
		background: var(--surface);
		border-radius: var(--radius-xl);
		border: 1px solid var(--border);
	}
</style>
