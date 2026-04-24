<template>
  <div class="flex min-h-screen bg-slate-50">
    <AdminSidebar />
    <main class="flex-1 p-8">
      <div class="mb-8">
        <h1 class="text-3xl font-display font-bold text-slate-800">All Appointments</h1>
        <p class="text-slate-500 mt-1">View and manage all patient appointments.</p>
      </div>

      <div class="flex gap-2 mb-6 flex-wrap">
        <button v-for="tab in tabs" :key="tab" @click="activeTab = tab"
          class="px-4 py-2 rounded-full text-sm font-medium capitalize transition-all"
          :class="activeTab === tab ? 'bg-blue-600 text-white' : 'bg-white border border-slate-200 text-slate-600 hover:border-blue-300'">
          {{ tab }}
        </button>
      </div>

      <div class="card overflow-hidden p-0">
        <div v-if="loading" class="text-center py-16 text-slate-400">Loading…</div>

        <div v-else-if="filtered.length === 0" class="text-center py-16 text-slate-400">
          No {{ activeTab !== 'all' ? activeTab : '' }} appointments found.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-slate-50">
              <tr>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Patient</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Doctor</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Branch</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Department</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Date & Time</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Status</th>
                <th class="text-left px-6 py-3 font-medium text-slate-500">Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr v-for="appt in filtered" :key="appt.appointment_id" class="hover:bg-slate-50">
                <td class="px-6 py-4 font-medium text-slate-800">{{ appt.patient_name }}</td>
                <td class="px-6 py-4 text-slate-600">Dr. {{ appt.doctor_name }}</td>
                <td class="px-6 py-4 text-slate-500">{{ appt.branch_name }}</td>
                <td class="px-6 py-4 text-slate-500">{{ appt.department_name }}</td>
                <td class="px-6 py-4 text-slate-600">
                  <div>{{ formatDate(appt.appointment_date) }}</div>
                  <div class="text-xs text-slate-400">{{ appt.appointment_time?.slice(0,5) }}</div>
                </td>
                <td class="px-6 py-4"><span :class="`badge-${appt.status}`">{{ appt.status }}</span></td>
                <td class="px-6 py-4">
                  <select v-if="['pending','confirmed'].includes(appt.status)"
                    @change="e => changeStatus(appt.appointment_id, e.target.value)"
                    class="text-xs border border-slate-200 rounded-lg px-2 py-1.5 text-slate-600 focus:outline-none focus:ring-1 focus:ring-blue-400">
                    <option value="">Change status…</option>
                    <option value="confirmed">Confirm</option>
                    <option value="completed">Complete</option>
                    <option value="cancelled">Cancel</option>
                  </select>
                  <span v-else class="text-xs text-slate-400">—</span>
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
import { ref, computed, onMounted } from 'vue'
import AdminSidebar from '@/components/AdminSidebar.vue'
import api from '@/services/api'

const appointments = ref([])
const loading      = ref(false)
const activeTab    = ref('all')
const tabs         = ['all', 'pending', 'confirmed', 'completed', 'cancelled']

const filtered = computed(() =>
  activeTab.value === 'all' ? appointments.value : appointments.value.filter(a => a.status === activeTab.value)
)

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await api.get('/appointments')
    appointments.value = data || []
  } finally {
    loading.value = false
  }
})

async function changeStatus(id, status) {
  if (!status) return
  await api.patch(`/appointments/${id}`, { status })
  const appt = appointments.value.find(a => a.appointment_id === id)
  if (appt) appt.status = status
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-KE', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

