<script setup lang="ts">
import { Bar, Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  BarElement,
  ArcElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend,
} from 'chart.js'

ChartJS.register(BarElement, ArcElement, CategoryScale, LinearScale, Tooltip, Legend)

definePageMeta({ title: 'ダッシュボード' })

const analyticsStore = useAnalyticsStore()
const router = useRouter()

onMounted(async () => {
  await analyticsStore.fetchDashboard()
})

const d = computed(() => analyticsStore.dashboard)

const statusCards = computed(() => {
  if (!d.value) return []
  const c = d.value.counts_by_status
  return [
    { label: '新規受付', value: c.new, color: '#3b82f6', status: 'new' },
    { label: '対応中', value: c.in_progress, color: '#7c4d33', status: 'in_progress' },
    { label: '確認待ち', value: c.waiting, color: '#8b5cf6', status: 'waiting' },
    { label: '完了', value: c.done, color: '#22c55e', status: 'done' },
  ]
})

const barChartData = computed(() => {
  if (!d.value) return { labels: [], datasets: [] }
  const trend = d.value.monthly_trend
  return {
    labels: trend.map(t => t.month),
    datasets: [{
      label: '件数',
      data: trend.map(t => t.count),
      backgroundColor: '#00B7B530',
      borderColor: '#00B7B5',
      borderWidth: 2,
      borderRadius: 4,
    }],
  }
})

const donutChartData = computed(() => {
  if (!d.value) return { labels: [], datasets: [] }
  const tc = d.value.open_counts_by_type
  return {
    labels: ['操作方法の質問', 'バグ報告', '設定変更依頼', 'データ修正依頼'],
    datasets: [{
      data: [tc.question, tc.bug, tc.config, tc.data],
      backgroundColor: ['#00B7B5', '#ef4444', '#f59e0b', '#018790'],
      borderWidth: 2,
      borderColor: '#fff',
    }],
  }
})

const barOptions = {
  responsive: true,
  plugins: { legend: { display: false } },
  scales: { y: { beginAtZero: true, ticks: { stepSize: 1 } } },
}
const donutOptions = {
  responsive: true,
  plugins: { legend: { position: 'bottom' as const } },
}

function goToList(status: string) {
  router.push({ path: '/tickets', query: { status } })
}
</script>

<template>
  <div v-if="analyticsStore.loading" class="p-8 text-slate-500">読み込み中...</div>

  <div v-else-if="d" class="p-6 space-y-6">
    <!-- ステータスカード -->
    <div class="grid grid-cols-4 gap-4">
      <div
        v-for="card in statusCards"
        :key="card.status"
        class="stat-card cursor-pointer hover:shadow-md transition-shadow"
        :style="`border-color: ${card.color};`"
        @click="goToList(card.status)"
      >
        <div class="stat-label">{{ card.label }}</div>
        <div class="stat-value" :style="`color: ${card.color};`">{{ card.value }}</div>
      </div>
    </div>

    <!-- グラフ行 -->
    <div class="grid grid-cols-3 gap-4">
      <div class="chart-card col-span-2">
        <div class="chart-card-title">月次問い合わせ件数（過去6ヶ月）</div>
        <ClientOnly>
          <Bar :data="barChartData" :options="barOptions" style="max-height:260px" />
        </ClientOnly>
      </div>
      <div class="chart-card">
        <div class="chart-card-title">種別内訳（未完了）</div>
        <ClientOnly>
          <Doughnut :data="donutChartData" :options="donutOptions" style="max-height:260px" />
        </ClientOnly>
      </div>
    </div>

    <!-- テーブル行 -->
    <div class="grid grid-cols-2 gap-4">
      <!-- システム別未完了件数 -->
      <div class="data-table">
        <div class="px-4 py-3 border-b border-slate-100">
          <h3 class="text-sm font-semibold text-slate-700">システム別未完了件数</h3>
        </div>
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-slate-50">
              <th class="px-4 py-2 text-left text-xs text-slate-500 font-medium">システム名</th>
              <th class="px-4 py-2 text-right text-xs text-slate-500 font-medium">件数</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="sys in d.top_systems_by_open_count"
              :key="sys.system_id"
              class="border-t border-slate-100"
            >
              <td class="px-4 py-2 text-slate-700">{{ sys.system_name }}</td>
              <td class="px-4 py-2 text-right font-semibold text-slate-800">{{ sys.count }}</td>
            </tr>
            <tr v-if="d.top_systems_by_open_count.length === 0">
              <td colspan="2" class="px-4 py-4 text-center text-slate-400 text-sm">データなし</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 優先度別件数 -->
      <div class="data-table">
        <div class="px-4 py-3 border-b border-slate-100">
          <h3 class="text-sm font-semibold text-slate-700">優先度別未完了件数</h3>
        </div>
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-slate-50">
              <th class="px-4 py-2 text-left text-xs text-slate-500 font-medium">優先度</th>
              <th class="px-4 py-2 text-right text-xs text-slate-500 font-medium">件数</th>
            </tr>
          </thead>
          <tbody>
            <tr class="border-t border-slate-100">
              <td class="px-4 py-2"><span class="badge-priority-high px-2 py-0.5 rounded text-xs font-medium">高</span></td>
              <td class="px-4 py-2 text-right font-semibold">{{ d.open_counts_by_priority.high }}</td>
            </tr>
            <tr class="border-t border-slate-100">
              <td class="px-4 py-2"><span class="badge-priority-medium px-2 py-0.5 rounded text-xs font-medium">中</span></td>
              <td class="px-4 py-2 text-right font-semibold">{{ d.open_counts_by_priority.medium }}</td>
            </tr>
            <tr class="border-t border-slate-100">
              <td class="px-4 py-2"><span class="badge-priority-low px-2 py-0.5 rounded text-xs font-medium">低</span></td>
              <td class="px-4 py-2 text-right font-semibold">{{ d.open_counts_by_priority.low }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <div v-else class="p-8 text-slate-400">データを取得できませんでした。</div>
</template>
