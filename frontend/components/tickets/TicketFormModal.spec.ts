import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import type { components } from '~/types/api'
import TicketFormModal from './TicketFormModal.vue'

type Ticket = components['schemas']['Ticket']

// Mock Nuxt auto-imports
vi.mock('#app', () => ({
  useNuxtApp: () => ({ $api: { request: vi.fn() } }),
  definePageMeta: vi.fn(),
}))

// Provide Nuxt-injected composables used inside the component
const mockNuxtApp = {
  $api: { request: vi.fn().mockResolvedValue([]) },
}

vi.stubGlobal('useNuxtApp', () => mockNuxtApp)
vi.stubGlobal('useRouter', () => ({ push: vi.fn() }))
vi.stubGlobal('useRoute', () => ({ params: {}, query: {} }))

// Mock all stores to prevent Pinia-related errors in isolation
vi.stubGlobal('useCustomersStore', () => ({
  customers: [],
  loading: false,
  fetchList: vi.fn(),
}))
vi.stubGlobal('useSystemsStore', () => ({
  systems: [],
  loading: false,
  fetchList: vi.fn(),
  fetchAssignees: vi.fn().mockResolvedValue([]),
}))
vi.stubGlobal('useTicketsStore', () => ({
  create: vi.fn(),
  update: vi.fn(),
}))

interface ModalProps {
  modelValue: boolean
  ticket: Ticket | null
}

const defaultProps: ModalProps = {
  modelValue: true,
  ticket: null,
}

function createWrapper(props: ModalProps = defaultProps) {
  return mount(TicketFormModal, {
    props,
    global: {
      plugins: [createTestingPinia({ createSpy: vi.fn })],
      stubs: {
        Teleport: true,
      },
    },
  })
}

describe('TicketFormModal', () => {
  describe('バリデーション', () => {
    it('タイトルが空のまま保存しようとするとエラーが表示される', async () => {
      const wrapper = createWrapper()
      await wrapper.find('[data-testid="submit-button"]').trigger('click')
      const errorEl = wrapper.find('[data-testid="error-title"]')
      expect(errorEl.exists()).toBe(true)
      expect(errorEl.text()).toContain('必須')
    })
  })

  describe('編集モード（customer_id ロック）', () => {
    it('編集時は顧客企業セレクトが disabled になる', () => {
      const sampleTicket: Ticket = {
        id: 1,
        title: 'テスト',
        type: 'bug',
        priority: 'high',
        status: 'new',
        customer_id: 1,
        customer_name: 'テスト企業',
        system_id: 1,
        system_name: 'テストシステム',
        created_by: 1,
        created_by_name: '管理者',
        received_at: '2026-05-17T00:00:00Z',
        created_at: '2026-05-17T00:00:00Z',
        updated_at: '2026-05-17T00:00:00Z',
      }
      const wrapper = createWrapper({
        modelValue: true,
        ticket: sampleTicket,
      })
      const customerSelect = wrapper.find('[data-testid="select-customer"]')
      expect(customerSelect.attributes('disabled')).toBeDefined()
    })

    it('編集時はシステムセレクトが disabled になる', () => {
      const sampleTicket2: Ticket = {
        id: 1,
        title: 'テスト',
        type: 'bug',
        priority: 'high',
        status: 'new',
        customer_id: 1,
        customer_name: 'テスト企業',
        system_id: 1,
        system_name: 'テストシステム',
        created_by: 1,
        created_by_name: '管理者',
        received_at: '2026-05-17T00:00:00Z',
        created_at: '2026-05-17T00:00:00Z',
        updated_at: '2026-05-17T00:00:00Z',
      }
      const wrapper = createWrapper({ modelValue: true, ticket: sampleTicket2 })
      const systemSelect = wrapper.find('[data-testid="select-system"]')
      expect(systemSelect.attributes('disabled')).toBeDefined()
    })

    it('新規登録時は顧客企業セレクトが enabled になる', () => {
      const wrapper = createWrapper({ modelValue: true, ticket: null })
      const customerSelect = wrapper.find('[data-testid="select-customer"]')
      expect(customerSelect.attributes('disabled')).toBeUndefined()
    })
  })
})
