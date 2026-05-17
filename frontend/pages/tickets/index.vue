<script setup lang="ts">
import type { TicketFilters } from '~/stores/tickets'
import type { components } from '~/types/api'

type TicketStatus = components['schemas']['TicketStatus']
type TicketPriority = components['schemas']['TicketPriority']

definePageMeta({ title: 'チケット一覧' })

const ticketsStore = useTicketsStore()
const customersStore = useCustomersStore()
const systemsStore = useSystemsStore()
const route = useRoute()

const filters = reactive<TicketFilters>({
  status: (route.query.status as TicketStatus) || '',
  priority: '',
  type: '',
  customer_id: '',
  system_id: '',
  keyword: '',
  page: 1,
  per_page: 50,
})

onMounted(async () => {
  await Promise.all([
    ticketsStore.fetchList(filters),
    customersStore.fetchList(true),
  ])
})

watch(() => filters.customer_id, async (cid) => {
  filters.system_id = ''
  if (cid) {
    await systemsStore.fetchList({ customerId: Number(cid), isActive: true })
  }
  else {
    systemsStore.systems = []
  }
})

watch(filters, async () => {
  filters.page = 1
  await ticketsStore.fetchList(filters)
}, { deep: true, flush: 'post' })

async function goPage(p: number) {
  filters.page = p
  await ticketsStore.fetchList(filters)
}

function resetFilters() {
  filters.status = ''
  filters.priority = ''
  filters.type = ''
  filters.customer_id = ''
  filters.system_id = ''
  filters.keyword = ''
  filters.page = 1
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ja-JP')
}

const statusLabels: Record<string, string> = {
  new: '新規受付',
  in_progress: '対応中',
  waiting: '確認待ち',
  done: '完了',
}
const priorityLabels: Record<string, string> = { high: '高', medium: '中', low: '低' }
const typeLabels: Record<string, string> = {
  question: '操作方法の質問',
  bug: 'バグ報告',
  config: '設定変更依頼',
  data: 'データ修正依頼',
}
</script>

<template>
  <div class="p-6 space-y-4">
    <!-- フィルターバー -->
    <div class="filter-bar p-4 flex flex-wrap gap-3 items-end">
      <div>
        <label class="block text-xs text-slate-500 mb-1">ステータス</label>
        <select v-model="filters.status" class="text-sm border border-slate-200 rounded px-2 py-1.5">
          <option value="">すべて</option>
          <option value="new">新規受付</option>
          <option value="in_progress">対応中</option>
          <option value="waiting">確認待ち</option>
          <option value="done">完了</option>
        </select>
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1">優先度</label>
        <select v-model="filters.priority" class="text-sm border border-slate-200 rounded px-2 py-1.5">
          <option value="">すべて</option>
          <option value="high">高</option>
          <option value="medium">中</option>
          <option value="low">低</option>
        </select>
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1">種別</label>
        <select v-model="filters.type" class="text-sm border border-slate-200 rounded px-2 py-1.5">
          <option value="">すべて</option>
          <option value="question">操作方法の質問</option>
          <option value="bug">バグ報告</option>
          <option value="config">設定変更依頼</option>
          <option value="data">データ修正依頼</option>
        </select>
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1">顧客企業</label>
        <select v-model="filters.customer_id" class="text-sm border border-slate-200 rounded px-2 py-1.5">
          <option value="">すべて</option>
          <option v-for="c in customersStore.customers" :key="c.id" :value="c.id">
            {{ c.name }}
          </option>
        </select>
      </div>
      <div>
        <label class="block text-xs text-slate-500 mb-1">システム</label>
        <select v-model="filters.system_id" class="text-sm border border-slate-200 rounded px-2 py-1.5" :disabled="!filters.customer_id">
          <option value="">すべて</option>
          <option v-for="s in systemsStore.systems" :key="s.id" :value="s.id">
            {{ s.name }}
          </option>
        </select>
      </div>
      <div class="flex-1 min-w-40">
        <label class="block text-xs text-slate-500 mb-1">キーワード</label>
        <input
          v-model="filters.keyword"
          type="text"
          placeholder="タイトル・内容を検索"
          class="w-full text-sm border border-slate-200 rounded px-2 py-1.5"
        >
      </div>
      <button
        class="btn-outline-primary text-sm px-3 py-1.5 rounded"
        @click="resetFilters"
      >
        リセット
      </button>
    </div>

    <!-- テーブル -->
    <div class="data-table">
      <div v-if="ticketsStore.loading" class="p-8 text-center text-slate-400">読み込み中...</div>
      <template v-else>
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200">
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">ID</th>
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">タイトル</th>
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">顧客企業</th>
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">システム</th>
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">種別</th>
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">優先度</th>
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">ステータス</th>
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">担当者</th>
              <th class="px-4 py-3 text-left text-xs text-slate-500 font-medium">受付日</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="ticket in ticketsStore.list"
              :key="ticket.id"
              class="border-t border-slate-100 hover:bg-slate-50 cursor-pointer transition-colors"
              @click="$router.push(`/tickets/${ticket.id}`)"
            >
              <td class="px-4 py-3 text-slate-400 text-xs">#{{ ticket.id }}</td>
              <td class="px-4 py-3 text-slate-800 font-medium max-w-xs truncate">{{ ticket.title }}</td>
              <td class="px-4 py-3 text-slate-600">{{ ticket.customer_name }}</td>
              <td class="px-4 py-3 text-slate-600">{{ ticket.system_name }}</td>
              <td class="px-4 py-3 text-slate-600 text-xs">{{ typeLabels[ticket.type] }}</td>
              <td class="px-4 py-3">
                <span :class="`${TICKET_PRIORITY_BADGE_CLASS[ticket.priority as TicketPriority]} px-2 py-0.5 rounded text-xs font-medium`">
                  {{ priorityLabels[ticket.priority] }}
                </span>
              </td>
              <td class="px-4 py-3">
                <span :class="`${TICKET_STATUS_BADGE_CLASS[ticket.status as TicketStatus]} px-2 py-0.5 rounded text-xs font-medium`">
                  {{ statusLabels[ticket.status] }}
                </span>
              </td>
              <td class="px-4 py-3 text-slate-600">{{ ticket.assignee_name ?? '未アサイン' }}</td>
              <td class="px-4 py-3 text-slate-500 text-xs">{{ formatDate(ticket.received_at) }}</td>
            </tr>
            <tr v-if="ticketsStore.list.length === 0">
              <td colspan="9" class="px-4 py-8 text-center text-slate-400">
                該当する問い合わせがありません
              </td>
            </tr>
          </tbody>
        </table>

        <!-- ページネーション -->
        <div v-if="ticketsStore.pagination && ticketsStore.pagination.total_pages > 1" class="flex items-center justify-between px-4 py-3 border-t border-slate-100">
          <span class="text-xs text-slate-500">
            全 {{ ticketsStore.pagination.total }} 件
          </span>
          <div class="flex gap-1">
            <button
              v-for="p in ticketsStore.pagination.total_pages"
              :key="p"
              class="w-8 h-8 text-xs rounded"
              :class="p === ticketsStore.pagination.page ? 'bg-secondary text-white' : 'text-slate-600 hover:bg-slate-100'"
              @click="goPage(p)"
            >
              {{ p }}
            </button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
