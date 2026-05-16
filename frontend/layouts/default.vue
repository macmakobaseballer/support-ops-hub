<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useApiError } from '~/composables/useApiError'

const authStore = useAuthStore()
const { error: apiError, clearError } = useApiError()
const { login } = useAuth()
const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()

const navItems = [
  { label: 'ダッシュボード', to: '/', icon: '📊' },
  { label: 'チケット一覧', to: '/tickets', icon: '🎫' },
]

const masterItems = [
  { label: '顧客マスタ', to: '/customers', icon: '🏢' },
  { label: 'システムマスタ', to: '/systems', icon: '⚙️' },
  { label: 'ユーザー管理', to: '/admin/users', icon: '👤' },
]

const pageTitle = computed(() => (route.meta.title as string) || 'Support Ops Hub')

const roleBadgeClass = computed(() =>
  authStore.user?.role === 'admin'
    ? 'bg-accent/20 text-accent'
    : 'bg-white/15 text-white/70',
)

const roleBadgeLabel = computed(() =>
  authStore.user?.role === 'admin' ? '管理者' : '担当者',
)

const isActive = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

const devUsers = ['admin@example.com', 'member1@example.com', 'member2@example.com']

async function handleDevUserSwitch(e: Event) {
  const email = (e.target as HTMLSelectElement).value
  await login(email)
  router.go(0)
}

async function logout() {
  authStore.clearAuth()
  await navigateTo('/login')
}
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-page-bg">
    <!-- Sidebar -->
    <aside class="w-[240px] flex-shrink-0 bg-primary flex flex-col">
      <div class="px-5 py-5 border-b border-white/10">
        <NuxtLink to="/" class="block">
          <h1 class="text-white text-base font-semibold leading-tight">
            🎧 Support Ops Hub
          </h1>
        </NuxtLink>
      </div>

      <nav class="flex-1 overflow-y-auto py-3">
        <ul class="space-y-0.5 px-3">
          <li v-for="item in navItems" :key="item.to">
            <NuxtLink
              :to="item.to"
              class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
              :class="isActive(item.to)
                ? 'bg-secondary text-white font-medium'
                : 'text-white/70 hover:bg-white/10 hover:text-white'"
            >
              <span class="text-base leading-none">{{ item.icon }}</span>
              <span>{{ item.label }}</span>
            </NuxtLink>
          </li>
        </ul>

        <!-- Master section (admin only) -->
        <div v-if="authStore.isAdmin" class="mt-4 px-3">
          <p class="px-3 mb-1 text-xs font-medium uppercase tracking-wider text-white/35">
            マスタ管理
          </p>
          <ul class="space-y-0.5">
            <li v-for="item in masterItems" :key="item.to">
              <NuxtLink
                :to="item.to"
                class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
                :class="isActive(item.to)
                  ? 'bg-secondary text-white font-medium'
                  : 'text-white/70 hover:bg-white/10 hover:text-white'"
              >
                <span class="text-base leading-none">{{ item.icon }}</span>
                <span>{{ item.label }}</span>
              </NuxtLink>
            </li>
          </ul>
        </div>
      </nav>

      <!-- User footer -->
      <div class="px-5 py-4 border-t border-white/10 space-y-3">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-full bg-secondary flex items-center justify-center flex-shrink-0">
            <span class="text-white text-sm font-medium">
              {{ authStore.user?.name?.charAt(0) ?? '?' }}
            </span>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-white text-sm font-medium truncate">
              {{ authStore.user?.name ?? 'ゲスト' }}
            </p>
            <span
              v-if="authStore.user"
              class="inline-block mt-0.5 px-1.5 py-0.5 rounded text-[10px] font-medium"
              :class="roleBadgeClass"
            >
              {{ roleBadgeLabel }}
            </span>
          </div>
        </div>

        <!-- Dev user switcher — visible in development mode only -->
        <div v-if="config.public.isDev" class="pt-1">
          <p class="text-[10px] text-white/35 mb-1 uppercase tracking-wide">
            DEV: ユーザー切替
          </p>
          <select
            class="w-full text-xs bg-white/10 text-white rounded px-2 py-1 border border-white/20 focus:outline-none"
            :value="authStore.user?.email ?? ''"
            data-testid="dev-user-switcher"
            @change="handleDevUserSwitch"
          >
            <option
              v-for="email in devUsers"
              :key="email"
              :value="email"
            >
              {{ email }}
            </option>
          </select>
        </div>

        <button
          type="button"
          class="w-full text-left text-sm text-white/70 hover:text-white transition-colors"
          @click="logout"
        >
          ログアウト
        </button>
      </div>
    </aside>

    <!-- Main area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <header class="flex-shrink-0 bg-white border-b border-slate-200 px-6 py-4 flex items-center justify-between gap-4">
        <h2 class="text-lg font-semibold text-slate-800">
          {{ pageTitle }}
        </h2>
        <div class="flex items-center gap-2">
          <slot name="header-actions" />
        </div>
      </header>

      <main class="flex-1 overflow-y-auto">
        <slot />
      </main>
    </div>

    <Transition name="toast">
      <div
        v-if="apiError"
        class="fixed bottom-6 right-6 z-50 flex items-center gap-3 bg-red-600 text-white px-4 py-3 rounded-lg shadow-lg max-w-sm"
      >
        <span class="text-sm leading-snug flex-1">{{ apiError }}</span>
        <button
          type="button"
          class="text-white/70 hover:text-white flex-shrink-0"
          @click="clearError"
        >
          ✕
        </button>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.2s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
