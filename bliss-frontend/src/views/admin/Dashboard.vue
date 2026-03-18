<template>
  <div class="flex min-h-screen bg-slate-50">
    <AdminSidebar />
    <main class="flex-1 p-8">
      <div class="mb-8">
        <h1 class="text-3xl font-display font-bold text-slate-800">Dashboard</h1>
        <p class="text-slate-500 mt-1">Good {{ greeting }}, {{ auth.user?.full_name?.split(' ')[0] }} 👋</p>
      </div>

      <div class="grid sm:grid-cols-2 lg:grid-cols-5 gap-4 mb-8">
        <div v-for="stat in statCards" :key="stat.label" class="card hover:shadow-md transition-shadow">
          <div class="flex items-start justify-between mb-2">
            <span class="text-2xl">{{ stat.icon }}</span>
            <span class="text-xs font-medium px-2 py-0.5 rounded-full" :class="stat.badge">{{ stat.tag }}</span>
          </div>
          <div class="text-3xl font-display font-bold text-slate-800">{{ stats ? stats[stat.key] : '—' }}</div>
          <div class="text-sm text-slate-500 mt-1">{{ stat.label }}</div>
        </div>
      </div>

      <div class="card">
        <div class="flex items-center justify-between mb-5">
          <h2 class="font-display font-semibold text-lg text-slate-800">Today's Queue</h2>
          <RouterLink to="/admin/queue" class="text-sm text-blue-600 hover:underline">Manage Queue →</RouterLink>
        </div>

        <div v-if="queueStore.fullQueue.length === 0" class="text-center py-8 text-slate-400">
          No patients in the queue today.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-slate-100">
                <th class="text-left pb-3 font-medium text-slate-500">#</th>
                <th class="text-left pb-3 font-medium text-slate-500">Patient</th>
                <th class="text-left pb-3 font-medium text-slate-500">Doctor</th>
                <th class="text-left pb-3 font-medium text-slate-500">Checked In</th>
                <th class="text-left pb-3 font-medium text-slate-500">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="entry in queueStore.fullQueue.slice(0, 8)" :key="entry.queue_id"
                class="border-b border-slate-50 hover:bg-slate-50">
                <td class="py-3 font-bold text-blue-600">{{ entry.queue_number }}</td>
                <td class="py-3 font-medium text-slate-800">{{ entry.patient_name }}</td>
                <td class="py-3 text-slate-600">Dr. {{ entry.doctor_name }}</td>
                <td class="py-3 text-slate-500">{{ formatTime(entry.checked_in_at) }}</td>
                <td class="py-3"><span :class="`badge-${entry.status}`">{{ entry.status }}</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useQueueStore } from '@/stores/queue'
import AdminSidebar from '@/components/AdminSidebar.vue'
import api from '@/services/api'

const auth       = useAuthStore()
const queueStore = useQueueStore()
const stats      = ref(null)

const greeting = computed(() => {
  const h = new Date().getHours()
  return h < 12 ? 'morning' : h < 17 ? 'afternoon' : 'evening'
})

const statCards = [
  { icon: '👥', label: 'Total Patients',       key: 'total_patients',       tag: 'All time',      badge: 'bg-blue-100 text-blue-700' },
  { icon: '🩺', label: 'Active Doctors',        key: 'total_doctors',        tag: 'Available',     badge: 'bg-green-100 text-green-700' },
  { icon: '📅', label: "Today's Appointments",  key: 'today_appointments',   tag: 'Today',         badge: 'bg-purple-100 text-purple-700' },
  { icon: '⏳', label: 'Pending Approvals',     key: 'pending_appointments', tag: 'Action needed', badge: 'bg-yellow-100 text-yellow-700' },
  { icon: '🔢', label: 'In Queue Now',          key: 'waiting_queue',        tag: 'Live',          badge: 'bg-orange-100 text-orange-700' },
]

onMounted(async () => {
  const { data } = await api.get('/admin/dashboard')
  stats.value = data
  await queueStore.fetchFullQueue()
})

function formatTime(ts) {
  if (!ts) return '—'
  return new Date(ts).toLocaleTimeString('en-KE', { hour: '2-digit', minute: '2-digit' })
}
</script>

