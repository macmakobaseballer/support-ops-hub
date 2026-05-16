<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useApiError } from '~/composables/useApiError'

const authStore = useAuthStore()
const { error: apiError, clearError } = useApiError()
const route = useRoute()

const navItems = [
  { label: 'ダッシュボード', to: '/', icon: '📊' },
  { label: 'チケット一覧', to: '/tickets', icon: '🎫' },
]

const masterItems = [
  { label: '顧客マスタ', to: '/customers', icon: '🏢' },
  { label: 'システムマスタ', to: '/systems', icon: '⚙️' },
]

const isActive = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-page-bg">
    <!-- Sidebar -->
    <aside class="w-[240px] flex-shrink-0 bg-primary flex flex-col">
      <!-- Logo -->
      <div class="px-5 py-5 border-b border-white/10">
        <NuxtLink to="/" class="block">
          <h1 class="text-white text-base font-semibold leading-tight">
            Support Ops Hub
          </h1>
        </NuxtLink>
      </div>

      <!-- Navigation -->
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

        <!-- Master section -->
        <div class="mt-4 px-3">
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

            <!-- admin only -->
            <li v-if="authStore.isAdmin">
              <NuxtLink
                to="/admin/users"
                class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
                :class="isActive('/admin/users')
                  ? 'bg-secondary text-white font-medium'
                  : 'text-white/70 hover:bg-white/10 hover:text-white'"
              >
                <span class="text-base leading-none">👤</span>
                <span>ユーザー管理</span>
              </NuxtLink>
            </li>
          </ul>
        </div>
      </nav>

      <!-- User footer -->
      <div class="px-5 py-4 border-t border-white/10">
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
            <p class="text-white/50 text-xs truncate">
              {{ authStore.user?.role ?? '' }}
            </p>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Page content -->
      <main class="flex-1 overflow-y-auto">
        <slot />
      </main>
    </div>

    <!-- Global error toast -->
    <Transition name="toast">
      <div
        v-if="apiError"
        class="fixed bottom-6 right-6 z-50 flex items-center gap-3 bg-red-600 text-white px-4 py-3 rounded-lg shadow-lg max-w-sm"
      >
        <span class="text-sm leading-snug flex-1">{{ apiError }}</span>
        <button
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
