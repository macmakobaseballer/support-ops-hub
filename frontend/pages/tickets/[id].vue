<script setup lang="ts">
import type { components } from '~/types/api'

type TicketStatus = components['schemas']['TicketStatus']
type TicketPriority = components['schemas']['TicketPriority']
type UserSummary = components['schemas']['UserSummary']

definePageMeta({ title: 'チケット詳細' })

const route = useRoute()
const ticketId = Number(route.params.id)

const ticketsStore = useTicketsStore()
const systemsStore = useSystemsStore()

const showEditModal = ref(false)
const assignees = ref<UserSummary[]>([])
const updatingStatus = ref(false)
const updatingAssignee = ref(false)

const ticket = computed(() => ticketsStore.current)

onMounted(async () => {
  await ticketsStore.fetchById(ticketId)
  if (ticket.value) {
    assignees.value = await systemsStore.fetchAssignees(ticket.value.system_id)
  }
})

// Transitions available from current status (ルール F4)
const availableTransitions = computed<TicketStatus[]>(() => {
  if (!ticket.value) return []
  return VALID_STATUS_TRANSITIONS[ticket.value.status as TicketStatus] ?? []
})

function transitionLabel(to: TicketStatus, currentStatus: TicketStatus): string {
  if (currentStatus === 'waiting') {
    return STATUS_TRANSITION_LABELS_FROM_WAITING[to] ?? STATUS_TRANSITION_LABELS[to] ?? to
  }
  return STATUS_TRANSITION_LABELS[to] ?? to
}

async function changeStatus(status: TicketStatus) {
  if (!ticket.value) return
  updatingStatus.value = true
  try {
    await ticketsStore.updateStatus(ticketId, status)
  }
  finally {
    updatingStatus.value = false
  }
}

async function changeAssignee(event: Event) {
  const select = event.target as HTMLSelectElement
  const val = select.value
  const assigneeId = val === '' ? null : Number(val)
  updatingAssignee.value = true
  try {
    await ticketsStore.updateAssignee(ticketId, assigneeId)
  }
  finally {
    updatingAssignee.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('ja-JP')
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
  <div>
    <div v-if="!ticket" class="p-8 text-slate-400">読み込み中...</div>

    <div v-else class="p-6 grid grid-cols-12 gap-6">
      <!-- 左カラム（メイン情報） -->
      <div class="col-span-8 space-y-4">
        <div class="form-card">
          <!-- 編集ボタン -->
          <div class="flex justify-end mb-3">
            <button class="btn-outline-primary text-sm px-3 py-1.5 rounded" @click="showEditModal = true">
              <i class="bi bi-pencil mr-1" /> 編集
            </button>
          </div>
          <!-- タイトル・メタ情報 -->
          <div class="mb-4">
            <div class="flex items-start gap-3">
              <span class="text-slate-400 text-sm mt-0.5">#{{ ticket.id }}</span>
              <h1 class="text-lg font-semibold text-slate-800 leading-tight">{{ ticket.title }}</h1>
            </div>
            <div class="flex items-center gap-3 mt-2">
              <span :class="`${TICKET_STATUS_BADGE_CLASS[ticket.status as TicketStatus]} px-2 py-0.5 rounded text-xs font-medium`">
                {{ statusLabels[ticket.status] }}
              </span>
              <span :class="`${TICKET_PRIORITY_BADGE_CLASS[ticket.priority as TicketPriority]} px-2 py-0.5 rounded text-xs font-medium`">
                {{ priorityLabels[ticket.priority] }}
              </span>
              <span class="text-xs text-slate-500">{{ typeLabels[ticket.type] }}</span>
            </div>
          </div>

          <!-- 詳細内容 -->
          <div class="mb-4">
            <div class="form-section-title mb-2">詳細内容</div>
            <div class="text-sm text-slate-700 whitespace-pre-wrap min-h-12">
              {{ ticket.description ?? '（内容なし）' }}
            </div>
          </div>

          <!-- ステータス変更ボタン -->
          <div v-if="availableTransitions.length > 0" class="flex gap-2 flex-wrap">
            <button
              v-for="to in availableTransitions"
              :key="to"
              class="btn-primary text-sm px-4 py-2 rounded text-white"
              :disabled="updatingStatus"
              :data-testid="`status-btn-${to}`"
              @click="changeStatus(to)"
            >
              {{ transitionLabel(to, ticket.status as TicketStatus) }}
            </button>
          </div>
        </div>

        <!-- M4 プレースホルダー（コメント・添付ファイル） -->
        <div class="form-card opacity-60">
          <div class="form-section-title mb-2">対応履歴・コメント</div>
          <p class="text-sm text-slate-400 py-4 text-center">M4 で実装予定</p>
        </div>
        <div class="form-card opacity-60">
          <div class="form-section-title mb-2">添付ファイル</div>
          <p class="text-sm text-slate-400 py-4 text-center">M4 で実装予定</p>
        </div>
      </div>

      <!-- 右カラム（サイドバー情報） -->
      <div class="col-span-4 space-y-4">
        <div class="form-card">
          <div class="form-section-title mb-3">チケット情報</div>
          <dl class="space-y-3 text-sm">
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">顧客企業</dt>
              <dd class="text-slate-800 font-medium">{{ ticket.customer_name }}</dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">システム</dt>
              <dd class="text-slate-800">{{ ticket.system_name }}</dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">種別</dt>
              <dd class="text-slate-800">{{ typeLabels[ticket.type] }}</dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">優先度</dt>
              <dd>
                <span :class="`${TICKET_PRIORITY_BADGE_CLASS[ticket.priority as TicketPriority]} px-2 py-0.5 rounded text-xs font-medium`">
                  {{ priorityLabels[ticket.priority] }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">ステータス</dt>
              <dd>
                <span :class="`${TICKET_STATUS_BADGE_CLASS[ticket.status as TicketStatus]} px-2 py-0.5 rounded text-xs font-medium`">
                  {{ statusLabels[ticket.status] }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">担当者</dt>
              <dd>
                <select
                  class="w-full text-sm border border-slate-200 rounded px-2 py-1.5 mt-0.5"
                  :value="ticket.assignee_id ?? ''"
                  :disabled="updatingAssignee"
                  data-testid="select-assignee"
                  @change="changeAssignee"
                >
                  <option value="">未アサイン</option>
                  <option v-for="a in assignees" :key="a.id" :value="a.id">{{ a.name }}</option>
                </select>
              </dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">受付日時</dt>
              <dd class="text-slate-800 text-xs">{{ formatDate(ticket.received_at) }}</dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">登録者</dt>
              <dd class="text-slate-800 text-xs">{{ ticket.created_by_name }}</dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500 mb-0.5">最終更新</dt>
              <dd class="text-slate-800 text-xs">{{ formatDate(ticket.updated_at) }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>

    <!-- 編集モーダル -->
    <TicketsTicketFormModal
      v-model="showEditModal"
      :ticket="ticket"
    />
  </div>
</template>
