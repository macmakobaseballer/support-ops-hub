<script setup lang="ts">
definePageMeta({
  layout: 'blank',
  title: 'ログイン',
})

const { login } = useAuth()
const { setError } = useApiError()
const router = useRouter()

const devUsers = [
  { label: '管理者 (admin@example.com)', email: 'admin@example.com' },
  { label: 'メンバー1 (member1@example.com)', email: 'member1@example.com' },
  { label: 'メンバー2 (member2@example.com)', email: 'member2@example.com' },
]

const selectedEmail = ref(devUsers[0]!.email)
const isLoading = ref(false)

async function handleLogin() {
  isLoading.value = true
  try {
    await login(selectedEmail.value)
    await router.push('/')
  }
  catch (err: unknown) {
    const apiErr = err as { message?: string }
    setError(apiErr.message ?? 'ログインに失敗しました')
  }
  finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-container min-h-screen flex items-center justify-center px-4">
    <div class="login-card w-full max-w-md p-8">
      <div class="text-center">
        <h1 class="text-xl font-semibold text-primary">
          Support Ops Hub
        </h1>
        <p class="mt-1 text-sm text-slate-500">
          問い合わせ管理システム
        </p>
        <span class="mt-2 inline-block rounded bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700">
          開発モード（dev-auth）
        </span>
      </div>

      <div class="mt-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">
            テストユーザー
          </label>
          <select
            v-model="selectedEmail"
            class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            data-testid="dev-user-select"
          >
            <option
              v-for="u in devUsers"
              :key="u.email"
              :value="u.email"
            >
              {{ u.label }}
            </option>
          </select>
        </div>

        <button
          type="button"
          class="btn-primary w-full rounded-lg py-2.5 text-sm font-medium"
          :disabled="isLoading"
          data-testid="login-button"
          @click="handleLogin"
        >
          {{ isLoading ? 'ログイン中...' : 'ログイン' }}
        </button>
      </div>
    </div>
  </div>
</template>
