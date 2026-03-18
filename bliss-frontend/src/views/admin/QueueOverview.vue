<template>
  <div class="flex min-h-screen bg-slate-50">
    <AdminSidebar />
    <main class="flex-1 p-8">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-display font-bold text-slate-800">Queue Management</h1>
          <p class="text-slate-500 mt-1">Manage today's patient queue in real time.</p>
        </div>
        <button @click="refresh" class="btn-secondary text-sm">🔄 Refresh</button>
      </div>

      <div v-if="queueStore.stats" class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div v-for="s in statItems" :key="s.label" class="card text-center">
          <div class="text-3xl font-display font-bold" :class="s.color">{{ queueStore.stats[s.key] }}</div>
          <div class="text-sm text-slate-500 mt-1">{{ s.label }}</div>
        </div>
      </div>

      <div class="card overflow-hidden p-0">
        <div class="px-6 py-4 border-b border-slate-100">
          <h2 class="font-display font-semibold text-slate-800">Today's Queue</h2>
        </div>

        <div v-if="queueStore.fullQueue.length === 0" class="text-center py-16 text-slate-400">
          No patients in the queue today.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-slate-50">
              <tr>
                <th class="text-left px-6 py-3 font-medium text-slate-500">#</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Patient</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Doctor</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Checked In</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Status</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr v-for="entry in queueStore.fullQueue" :key="entry.queue_id"
                class="hover:bg-slate-50 transition-colors"
                :class="entry.status === 'called' ? 'bg-purple-50' : ''">
                <td class="px-6 py-4 font-bold text-blue-600 text-lg">{{ entry.queue_number }}</td>
                <td class="px-6 py-4 font-medium text-slate-800">{{ entry.patient_name }}</td>
                <td class="px-6 py-4 text-slate-600">Dr. {{ entry.doctor_name }}</td>
                <td class="px-6 py-4 text-slate-500">{{ formatTime(entry.checked_in_at) }}</td>
                <td class="px-6 py-4"><span :class="`badge-${entry.status}`">{{ entry.status }}</span></td>
                <td class="px-6 py-4">
                  <div class="flex gap-2">
                    <button v-if="entry.status === 'waiting'"
                      @click="updateEntry(entry.queue_id, 'called')"
                      class="text-xs bg-purple-100 text-purple-700 hover:bg-purple-200 px-3 py-1.5 rounded-lg font-medium transition">
                      Call Patient
                    </button>
                    <button v-if="entry.status === 'called'"
                      @click="updateEntry(entry.queue_id, 'serving')"
                      class="text-xs bg-teal-100 text-teal-700 hover:bg-teal-200 px-3 py-1.5 rounded-lg font-medium transition">
                      Start Serving
                    </button>
                    <button v-if="entry.status === 'serving'"
                      @click="updateEntry(entry.queue_id, 'done')"
                      class="text-xs bg-green-100 text-green-700 hover:bg-green-200 px-3 py-1.5 rounded-lg font-medium transition">
                      Mark Done
                    </button>
                    <button v-if="['waiting','called'].includes(entry.status)"
                      @click="updateEntry(entry.queue_id, 'skipped')"
                      class="text-xs bg-slate-100 text-slate-600 hover:bg-slate-200 px-3 py-1.5 rounded-lg font-medium transition">
                      Skip
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useQueueStore } from '@/stores/queue'
import AdminSidebar from '@/components/AdminSidebar.vue'

const queueStore = useQueueStore()
let interval = null

const statItems = [
  { key: 'waiting', label: 'Waiting', color: 'text-orange-600' },
  { key: 'serving', label: 'Serving', color: 'text-teal-600' },
  { key: 'done',    label: 'Done',    color: 'text-green-600' },
  { key: 'total',   label: 'Total',   color: 'text-slate-800' },
]

onMounted(async () => {
  await refresh()
  interval = setInterval(refresh, 20000)
})

onUnmounted(() => clearInterval(interval))

async function refresh() {
  await queueStore.fetchFullQueue()
  await queueStore.fetchStats()
}

async function updateEntry(id, status) {
  await queueStore.updateQueueEntry(id, status)
  await queueStore.fetchStats()
}

function formatTime(ts) {
  if (!ts) return '—'
  return new Date(ts).toLocaleTimeString('en-KE', { hour: '2-digit', minute: '2-digit' })
}
</script>

