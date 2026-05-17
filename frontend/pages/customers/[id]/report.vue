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

definePageMeta({ title: '顧客レポート' })

const route = useRoute()
const customerId = Number(route.params.id)

const analyticsStore = useAnalyticsStore()

onMounted(async () => {
  await analyticsStore.fetchCustomerReport(customerId)
})

const report = computed(() => analyticsStore.customerReport)

const statusCards = computed(() => {
  if (!report.value) return []
  const c = report.value.counts_by_status
  return [
    { label: '新規受付', value: c.new, color: '#3b82f6' },
    { label: '対応中', value: c.in_progress, color: '#7c4d33' },
    { label: '確認待ち', value: c.waiting, color: '#8b5cf6' },
    { label: '完了', value: c.done, color: '#22c55e' },
  ]
})

const barChartData = computed(() => {
  if (!report.value) return { labels: [], datasets: [] }
  return {
    labels: report.value.monthly_trend.map(t => t.month),
    datasets: [{
      label: '件数',
      data: report.value.monthly_trend.map(t => t.count),
      backgroundColor: '#00B7B530',
      borderColor: '#00B7B5',
      borderWidth: 2,
      borderRadius: 4,
    }],
  }
})

const donutChartData = computed(() => {
  if (!report.value) return { labels: [], datasets: [] }
  const tc = report.value.open_counts_by_type
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
</script>

<template>
  <div v-if="analyticsStore.loading" class="p-8 text-slate-500">読み込み中...</div>

  <div v-else-if="report" class="p-6 space-y-6">
    <h1 class="text-xl font-semibold text-slate-800">{{ report.customer_name }} レポート</h1>

    <!-- ステータスカード -->
    <div class="grid grid-cols-4 gap-4">
      <div
        v-for="card in statusCards"
        :key="card.label"
        class="stat-card"
        :style="`border-color: ${card.color};`"
      >
        <div class="stat-label">{{ card.label }}</div>
        <div class="stat-value" :style="`color: ${card.color};`">{{ card.value }}</div>
      </div>
    </div>

    <!-- グラフ -->
    <div class="grid grid-cols-3 gap-4">
      <div class="chart-card col-span-2">
        <div class="chart-card-title">月次件数推移（過去6ヶ月）</div>
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

    <!-- システム別集計テーブル -->
    <div class="data-table">
      <div class="px-4 py-3 border-b border-slate-100">
        <h3 class="text-sm font-semibold text-slate-700">システム別ステータス集計</h3>
      </div>
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-slate-50">
            <th class="px-4 py-2 text-left text-xs text-slate-500 font-medium">システム名</th>
            <th class="px-4 py-2 text-center text-xs text-slate-500 font-medium">新規受付</th>
            <th class="px-4 py-2 text-center text-xs text-slate-500 font-medium">対応中</th>
            <th class="px-4 py-2 text-center text-xs text-slate-500 font-medium">確認待ち</th>
            <th class="px-4 py-2 text-center text-xs text-slate-500 font-medium">完了</th>
            <th class="px-4 py-2 text-center text-xs text-slate-500 font-medium">合計</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="sys in report.system_breakdown"
            :key="sys.system_id"
            class="border-t border-slate-100"
          >
            <td class="px-4 py-2 text-slate-700 font-medium">{{ sys.system_name }}</td>
            <td class="px-4 py-2 text-center text-slate-600">{{ sys.counts.new }}</td>
            <td class="px-4 py-2 text-center text-slate-600">{{ sys.counts.in_progress }}</td>
            <td class="px-4 py-2 text-center text-slate-600">{{ sys.counts.waiting }}</td>
            <td class="px-4 py-2 text-center text-slate-600">{{ sys.counts.done }}</td>
            <td class="px-4 py-2 text-center font-semibold text-slate-800">{{ sys.total }}</td>
          </tr>
          <tr v-if="report.system_breakdown.length === 0">
            <td colspan="6" class="px-4 py-4 text-center text-slate-400 text-sm">データなし</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <div v-else class="p-8 text-slate-400">データを取得できませんでした。</div>
</template>
