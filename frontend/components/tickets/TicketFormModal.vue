<script setup lang="ts">
import type { components } from '~/types/api'

type Ticket = components['schemas']['Ticket']
type TicketCreate = components['schemas']['TicketCreate']
type TicketUpdate = components['schemas']['TicketUpdate']
type TicketType = components['schemas']['TicketType']
type TicketPriority = components['schemas']['TicketPriority']
type UserSummary = components['schemas']['UserSummary']

const props = defineProps<{
  modelValue: boolean
  ticket?: Ticket | null
}>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: [ticket: Ticket]
}>()

const isEdit = computed(() => !!props.ticket)

const ticketsStore = useTicketsStore()
const customersStore = useCustomersStore()
const systemsStore = useSystemsStore()

const form = reactive({
  customer_id: 0,
  system_id: 0,
  type: '' as TicketType | '',
  priority: '' as TicketPriority | '',
  title: '',
  description: '',
  assignee_id: null as number | null,
  received_at: '',
})
const formErrors = reactive({
  customer_id: '',
  system_id: '',
  type: '',
  priority: '',
  title: '',
})
const assignees = ref<UserSummary[]>([])
const saving = ref(false)

function resetForm() {
  const t = props.ticket
  if (t) {
    form.customer_id = t.customer_id
    form.system_id = t.system_id
    form.type = t.type
    form.priority = t.priority
    form.title = t.title
    form.description = t.description ?? ''
    form.assignee_id = t.assignee_id ?? null
    form.received_at = t.received_at.slice(0, 16)
  }
  else {
    form.customer_id = 0
    form.system_id = 0
    form.type = ''
    form.priority = ''
    form.title = ''
    form.description = ''
    form.assignee_id = null
    form.received_at = new Date().toISOString().slice(0, 16)
  }
  formErrors.customer_id = ''
  formErrors.system_id = ''
  formErrors.type = ''
  formErrors.priority = ''
  formErrors.title = ''
}

watch(() => props.modelValue, async (open) => {
  if (!open) return
  resetForm()
  await customersStore.fetchList(true)
  if (isEdit.value && form.system_id) {
    assignees.value = await systemsStore.fetchAssignees(form.system_id)
  }
})

watch(() => form.customer_id, async (cid) => {
  if (!cid) {
    systemsStore.systems = []
    return
  }
  await systemsStore.fetchList({ customerId: cid, isActive: true })
  if (!isEdit.value) {
    form.system_id = 0
    form.assignee_id = null
    assignees.value = []
  }
})

watch(() => form.system_id, async (sid) => {
  if (!sid) {
    assignees.value = []
    return
  }
  assignees.value = await systemsStore.fetchAssignees(sid)
  if (!isEdit.value) form.assignee_id = null
})

watch(() => form.type, (type) => {
  if (!type) return
  const defaults: Record<TicketType, TicketPriority> = {
    question: 'medium',
    bug: 'high',
    config: 'medium',
    data: 'medium',
  }
  form.priority = defaults[type as TicketType]
})

function validate(): boolean {
  formErrors.customer_id = ''
  formErrors.system_id = ''
  formErrors.type = ''
  formErrors.priority = ''
  formErrors.title = ''
  if (!form.title) formErrors.title = '必須項目です'
  if (form.title.length > 255) formErrors.title = '255文字以内で入力してください'
  if (!form.customer_id) formErrors.customer_id = '必須項目です'
  if (!form.system_id) formErrors.system_id = '必須項目です'
  if (!form.type) formErrors.type = '必須項目です'
  if (!form.priority) formErrors.priority = '必須項目です'
  return !formErrors.customer_id && !formErrors.system_id && !formErrors.type && !formErrors.priority && !formErrors.title
}

async function submit() {
  if (!validate()) return
  saving.value = true
  try {
    let ticket: Ticket
    if (isEdit.value && props.ticket) {
      const payload: TicketUpdate = {
        type: form.type as TicketType,
        priority: form.priority as TicketPriority,
        title: form.title,
        description: form.description || null,
        assignee_id: form.assignee_id,
        received_at: new Date(form.received_at).toISOString(),
      }
      ticket = await ticketsStore.update(props.ticket.id, payload)
    }
    else {
      const payload: TicketCreate = {
        customer_id: form.customer_id,
        system_id: form.system_id,
        type: form.type as TicketType,
        priority: form.priority as TicketPriority,
        title: form.title,
        description: form.description || null,
        assignee_id: form.assignee_id,
        received_at: new Date(form.received_at).toISOString(),
      }
      ticket = await ticketsStore.create(payload)
    }
    emit('saved', ticket)
    emit('update:modelValue', false)
  }
  finally {
    saving.value = false
  }
}

function close() {
  emit('update:modelValue', false)
}
</script>

<template>
  <Teleport to="body">
    <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="close" />
      <div class="relative bg-white rounded-xl shadow-2xl w-full max-w-2xl mx-4 max-h-[90vh] flex flex-col">
        <!-- ヘッダー -->
        <div class="flex items-center justify-between px-6 py-4 border-b border-slate-200">
          <h2 class="text-base font-semibold text-slate-800">
            {{ isEdit ? '問い合わせ編集' : '問い合わせ登録' }}
          </h2>
          <button class="text-slate-400 hover:text-slate-600" data-testid="modal-close" @click="close">
            <i class="bi bi-x-lg" />
          </button>
        </div>

        <!-- フォーム -->
        <div class="overflow-y-auto flex-1 px-6 py-4 space-y-4">
          <!-- 顧客企業 -->
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">顧客企業 <span class="text-red-500">*</span></label>
            <select
              v-model="form.customer_id"
              :disabled="isEdit"
              class="w-full text-sm border border-slate-200 rounded px-3 py-2 disabled:bg-slate-50 disabled:text-slate-400"
              data-testid="select-customer"
            >
              <option :value="0">選択してください</option>
              <option v-for="c in customersStore.customers" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
            <p v-if="formErrors.customer_id" class="text-xs text-red-500 mt-1">{{ formErrors.customer_id }}</p>
          </div>

          <!-- システム -->
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">システム <span class="text-red-500">*</span></label>
            <select
              v-model="form.system_id"
              :disabled="isEdit || !form.customer_id"
              class="w-full text-sm border border-slate-200 rounded px-3 py-2 disabled:bg-slate-50 disabled:text-slate-400"
              data-testid="select-system"
            >
              <option :value="0">選択してください</option>
              <option v-for="s in systemsStore.systems" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
            <p v-if="formErrors.system_id" class="text-xs text-red-500 mt-1">{{ formErrors.system_id }}</p>
          </div>

          <!-- 種別 -->
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">種別 <span class="text-red-500">*</span></label>
            <select v-model="form.type" class="w-full text-sm border border-slate-200 rounded px-3 py-2" data-testid="select-type">
              <option value="">選択してください</option>
              <option value="question">操作方法の質問</option>
              <option value="bug">バグ報告</option>
              <option value="config">設定変更依頼</option>
              <option value="data">データ修正依頼</option>
            </select>
            <p v-if="formErrors.type" class="text-xs text-red-500 mt-1">{{ formErrors.type }}</p>
          </div>

          <!-- 優先度 -->
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">優先度 <span class="text-red-500">*</span></label>
            <select v-model="form.priority" class="w-full text-sm border border-slate-200 rounded px-3 py-2" data-testid="select-priority">
              <option value="">選択してください</option>
              <option value="high">高</option>
              <option value="medium">中</option>
              <option value="low">低</option>
            </select>
            <p v-if="formErrors.priority" class="text-xs text-red-500 mt-1">{{ formErrors.priority }}</p>
          </div>

          <!-- タイトル -->
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">タイトル <span class="text-red-500">*</span></label>
            <input
              v-model="form.title"
              type="text"
              maxlength="255"
              placeholder="問い合わせの概要を入力してください"
              class="w-full text-sm border border-slate-200 rounded px-3 py-2"
              data-testid="input-title"
            >
            <p v-if="formErrors.title" class="text-xs text-red-500 mt-1" data-testid="error-title">{{ formErrors.title }}</p>
          </div>

          <!-- 詳細内容 -->
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">詳細内容</label>
            <textarea
              v-model="form.description"
              rows="4"
              placeholder="詳細な内容を入力してください"
              class="w-full text-sm border border-slate-200 rounded px-3 py-2 resize-none"
            />
          </div>

          <!-- 担当者 -->
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">担当者</label>
            <select v-model="form.assignee_id" class="w-full text-sm border border-slate-200 rounded px-3 py-2">
              <option :value="null">未アサイン</option>
              <option v-for="a in assignees" :key="a.id" :value="a.id">{{ a.name }}</option>
            </select>
          </div>

          <!-- 受付日時 -->
          <div>
            <label class="block text-xs font-medium text-slate-600 mb-1">受付日時</label>
            <input
              v-model="form.received_at"
              type="datetime-local"
              class="w-full text-sm border border-slate-200 rounded px-3 py-2"
            >
          </div>
        </div>

        <!-- フッター -->
        <div class="flex items-center justify-end gap-3 px-6 py-4 border-t border-slate-200">
          <button class="btn-outline-primary text-sm px-4 py-2 rounded" @click="close">キャンセル</button>
          <button
            class="btn-primary text-sm px-4 py-2 rounded text-white"
            :disabled="saving"
            data-testid="submit-button"
            @click="submit"
          >
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
