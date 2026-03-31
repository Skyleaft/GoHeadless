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
	<title>Access Control | GoHeadless CMS</title>
</svelte:head>

<div class="inset-0 p-6 fixed flex items-center justify-center overflow-hidden">
	<!-- Dynamic Background Elements -->
	<div class="inset-0 absolute z-0 bg-[var(--bg)]">
		<div
			class="absolute top-[-10%] left-[-10%] h-[40%] w-[40%] rounded-full bg-[var(--brand)] opacity-[0.03] blur-[120px]"
		></div>
		<div
			class="bg-purple-500 absolute right-[-10%] bottom-[-10%] h-[40%] w-[40%] rounded-full opacity-[0.03] blur-[120px]"
		></div>
	</div>

	<div class="relative z-10 w-full max-w-[420px]">
		{#if checking}
			<div class="gap-6 animate-fade-in flex flex-col items-center text-center">
				<div class="h-16 w-16 relative">
					<div
						class="inset-0 rounded-2xl animate-pulse absolute bg-[var(--brand)] opacity-20"
					></div>
					<div class="inset-0 absolute flex items-center justify-center">
						<div
							class="h-8 w-8 animate-spin shadow-lg rounded-full border-[3px] border-[var(--brand)] border-t-transparent shadow-[var(--brand)]/20"
						></div>
					</div>
				</div>
				<div class="space-y-1">
					<p class="text-sm font-black tracking-[0.2em] uppercase opacity-40">System Security</p>
					<p
						class="text-xs font-bold animate-pulse transition-all"
						style="color: var(--text-muted)"
					>
						Establishing encrypted tunnel...
					</p>
				</div>
			</div>
		{:else}
			<div class="login-card animate-scale-in group relative overflow-hidden">
				<!-- Header Section -->
				<div class="p-10 pb-6 relative z-10 text-center">
					<div
						class="mb-6 h-16 w-16 text-white text-2xl font-black shadow-2xl animate-fade-in group relative mx-auto flex items-center justify-center rounded-[1.5rem] shadow-[var(--brand)]/40 transition-transform hover:scale-110 active:scale-95"
						style="background: linear-gradient(135deg, var(--brand), #7c3aed)"
					>
						GH
						<div
							class="-inset-2 to-purple-600 blur-xl absolute bg-gradient-to-br from-[var(--brand)] opacity-20 transition-opacity group-hover:opacity-40"
						></div>
					</div>
					<h1 class="text-3xl font-black tracking-tighter mb-2" style="color: var(--text-primary)">
						Welcome Back
					</h1>
					<p class="text-sm font-bold opacity-50" style="color: var(--text-muted)">
						Access your GoHeadless management suite
					</p>
				</div>

				<!-- Form Section -->
				<form onsubmit={handleSubmit} class="px-10 pb-10 gap-6 relative z-10 flex flex-col">
					<div class="space-y-4">
						<Input
							id="username"
							label="Identifier"
							placeholder="Your username"
							required
							bind:value={username}
							disabled={loading}
							class="h-12 rounded-2xl shadow-inner font-bold border-transparent bg-[var(--bg)]"
						/>

						<div class="gap-1.5 flex flex-col">
							<div class="px-1 flex items-center justify-between">
								<label
									class="font-black tracking-widest text-[10px] uppercase opacity-40"
									style="color: var(--text-muted)"
									for="password"
								>
									Security Key
								</label>
								<button
									type="button"
									class="font-black tracking-widest text-[10px] text-[var(--brand)] uppercase hover:underline"
									tabindex="-1"
								>
									Reset Pin
								</button>
							</div>
							<Input
								id="password"
								type="password"
								placeholder="••••••••"
								required
								bind:value={password}
								disabled={loading}
								class="h-12 rounded-2xl shadow-inner font-bold border-transparent bg-[var(--bg)]"
							/>
						</div>
					</div>

					<div class="px-1 flex items-center justify-between">
						<label class="gap-3 group flex cursor-pointer items-center">
							<div
								class="h-5 w-5 rounded-lg relative flex items-center justify-center overflow-hidden border-2 border-[var(--border)] transition-all group-hover:border-[var(--brand)]"
							>
								<input type="checkbox" id="remember" class="peer sr-only" />
								<div
									class="inset-0 absolute bg-[var(--brand)] opacity-0 transition-opacity peer-checked:opacity-100"
								></div>
								<div
									class="h-2 w-2 bg-white z-10 rounded-full opacity-0 transition-opacity peer-checked:opacity-100"
								></div>
							</div>
							<span
								class="text-xs font-bold opacity-60 transition-opacity group-hover:opacity-100"
								style="color: var(--text-secondary)">Keep me active</span
							>
						</label>
					</div>

					<div class="pt-4">
						<Button
							type="submit"
							variant="primary"
							class="h-14 rounded-2xl font-black text-base shadow-xl gap-2 group flex w-full items-center justify-center transition-all active:scale-95"
							size="lg"
							{loading}
						>
							{#if !loading}
								Unlock Suite
								<span class="text-lg group-hover:translate-x-1 transition-transform">→</span>
							{:else}
								Unlocking...
							{/if}
						</Button>
					</div>
				</form>

				<!-- Decoration -->
				<div
					class="w-40 h-40 blur-3xl pointer-events-none absolute top-[-50px] right-[-50px] rounded-full bg-[var(--brand)] opacity-[0.02]"
				></div>
				<div
					class="w-40 h-40 bg-purple-500 blur-3xl pointer-events-none absolute bottom-[-50px] left-[-50px] rounded-full opacity-[0.02]"
				></div>
			</div>

			<div class="mt-8 animate-fade-in gap-2 flex flex-col items-center text-center">
				<p class="font-black pb-1 text-[10px] tracking-[0.3em] uppercase opacity-30">
					AES-256 Cloud Infrastructure
				</p>
				<div class="gap-4 flex">
					<button class="text-xs font-bold opacity-40 transition-opacity hover:opacity-100"
						>Cloud Panel</button
					>
					<span class="text-[10px] opacity-20">•</span>
					<button class="text-xs font-bold opacity-40 transition-opacity hover:opacity-100"
						>System Metrics</button
					>
					<span class="text-[10px] opacity-20">•</span>
					<button class="text-xs font-bold opacity-40 transition-opacity hover:opacity-100"
						>Documentation</button
					>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.login-card {
		background: rgba(var(--surface-rgb), 0.7);
		backdrop-filter: blur(40px);
		border-radius: 3rem;
		box-shadow: var(--shadow-xl);
		border: 1px solid var(--border);
	}

	:global(.dark) .login-card {
		background: rgba(var(--surface-rgb), 0.4);
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 0.2;
		}
		50% {
			opacity: 0.3;
		}
	}
	.animate-pulse {
		animation: pulse 2s infinite ease-in-out;
	}
</style>
