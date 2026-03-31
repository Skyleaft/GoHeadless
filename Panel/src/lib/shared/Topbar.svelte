<script lang="ts">
	import { page } from '$app/state';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	// Build breadcrumbs from current path
	let breadcrumbs = $derived(() => {
		const parts = page.url.pathname.split('/').filter(Boolean);
		const crumbs: { label: string; href: string }[] = [{ label: 'Home', href: '/' }];
		let path = '';
		for (const part of parts) {
			path += `/${part}`;
			crumbs.push({
				label: part.charAt(0).toUpperCase() + part.slice(1).replace(/-/g, ' '),
				href: path
			});
		}
		return crumbs;
	});

	let showUserMenu = $state(false);

	function handleLogout() {
		auth.logout();
		goto('/login');
	}

	function getInitials(name: string) {
		return name
			.split(/[_\s]/)
			.map((n) => n[0].toUpperCase())
			.join('')
			.slice(0, 2);
	}
</script>

<header
	class="top-0 px-8 sticky z-40 flex items-center justify-between transition-all duration-300"
	style="
		height: var(--topbar-height);
		background: rgba(var(--bg-rgb), 0.7);
		backdrop-filter: blur(24px);
		box-shadow: var(--shadow-sm);
	"
>
	<!-- Left: Breadcrumb -->
	<div class="gap-3 text-sm flex items-center">
		{#each breadcrumbs() as crumb, i}
			{#if i > 0}
				<span style="color: var(--text-muted)" class="opacity-40">/</span>
			{/if}
			{#if i === breadcrumbs().length - 1}
				<span class="font-bold tracking-tight" style="color: var(--text-primary)"
					>{crumb.label}</span
				>
			{:else}
				<a
					href={crumb.href}
					class="font-medium transition-all hover:text-[--brand]"
					style="color: var(--text-muted)">{crumb.label}</a
				>
			{/if}
		{/each}
	</div>

	<!-- Right: Actions -->
	<div class="gap-4 flex items-center">
		<!-- API Docs -->
		<a
			href="http://localhost:3030/docs"
			target="_blank"
			rel="noopener"
			class="md:flex h-10 rounded-xl px-4 text-sm font-bold hidden items-center justify-center bg-[--surface-alt] transition-all hover:scale-105 active:scale-95"
			style="color: var(--text-secondary)"
		>
			API Explorer
		</a>

		<!-- Theme Toggle -->
		<button
			onclick={toggleTheme}
			class="h-10 w-10 rounded-xl text-lg flex items-center justify-center bg-[--surface-alt] transition-all hover:scale-105 active:scale-95"
			style="color: var(--text-secondary)"
			title="Switch Appearance"
		>
			{$theme === 'dark' ? '🌙' : $theme === 'light' ? '🔆' : '🌗'}
		</button>

		<!-- Account Section -->
		<div class="relative">
			<button
				onclick={() => (showUserMenu = !showUserMenu)}
				class="gap-3 pl-3 pr-1 py-1 rounded-2xl hover:shadow-md flex items-center border border-[--border] bg-[--surface-alt] transition-all active:scale-95"
			>
				<div class="sm:flex hidden flex-col items-end leading-none">
					<span class="text-xs font-bold" style="color: var(--text-primary)"
						>{$auth.user?.username}</span
					>
					<span
						class="font-bold tracking-widest text-[10px] uppercase opacity-50"
						style="color: var(--text-muted)"
					>
						{$auth.user?.is_initial_admin ? 'Superadmin' : 'Staff'}
					</span>
				</div>
				<div
					class="h-8 w-8 rounded-xl text-sm font-black text-white shadow-lg flex items-center justify-center"
					style="background: linear-gradient(135deg, var(--brand), #7c3aed)"
				>
					{getInitials($auth.user?.username || 'U')}
				</div>
			</button>

			{#if showUserMenu}
				<div
					class="right-0 mt-4 w-64 rounded-2xl p-2 shadow-xl animate-scale-in absolute"
					style="background: rgba(var(--surface-rgb), 0.9); backdrop-filter: blur(24px); border: 1px solid var(--border);"
				>
					<div class="p-4 mb-1">
						<p
							class="font-black mb-1.5 text-[10px] tracking-[0.2em] text-[--text-muted] uppercase opacity-60"
						>
							Personal Profile
						</p>
						<p class="text-base font-black tracking-tight truncate">{$auth.user?.username}</p>
					</div>

					<div class="mx-2 mb-2 h-px bg-current opacity-[0.03]"></div>

					<button
						class="gap-3 p-3 rounded-xl text-sm font-bold flex w-full items-center transition-all hover:bg-[--surface-alt]"
						style="color: var(--text-secondary)"
					>
						<span>⚙️</span> Settings
					</button>

					<button
						onclick={handleLogout}
						class="gap-3 p-3 rounded-xl text-sm font-bold hover:bg-red-50 hover:text-red-500 flex w-full items-center transition-all"
					>
						<span>🚪</span> Sign Out
					</button>
				</div>
			{/if}
		</div>
	</div>
</header>
