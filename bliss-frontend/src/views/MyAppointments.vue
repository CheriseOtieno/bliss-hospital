<template>
  <div class="max-w-5xl mx-auto px-4 py-10">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
      <div>
        <h1 class="text-3xl font-display font-bold text-slate-800">My Appointments</h1>
        <p class="text-slate-500 mt-1">Manage and track all your bookings.</p>
      </div>
      <RouterLink to="/book" class="btn-primary self-start sm:self-auto">+ Book New</RouterLink>
    </div>

    <div class="flex gap-2 mb-6 flex-wrap">
      <button v-for="tab in tabs" :key="tab.value" @click="activeTab = tab.value"
        class="px-4 py-2 rounded-full text-sm font-medium transition-all"
        :class="activeTab === tab.value ? 'bg-blue-600 text-white' : 'bg-white text-slate-600 border border-slate-200 hover:border-blue-300'">
        {{ tab.label }} ({{ countByStatus(tab.value) }})
      </button>
    </div>

    <div v-if="store.loading" class="text-center py-20 text-slate-400">Loading appointments…</div>

    <div v-else-if="filtered.length === 0" class="card text-center py-16">
      <div class="text-5xl mb-4">📋</div>
      <h3 class="font-display font-semibold text-slate-700 mb-2">No appointments found</h3>
      <p class="text-slate-400 text-sm mb-6">You don't have any {{ activeTab !== 'all' ? activeTab : '' }} appointments yet.</p>
      <RouterLink to="/book" class="btn-primary">Book Your First Appointment</RouterLink>
    </div>

    <div v-else class="space-y-4">
      <div v-for="appt in filtered" :key="appt.appointment_id" class="card hover:shadow-md transition-shadow">
        <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
          <div class="flex items-start gap-4">
            <div class="w-12 h-12 rounded-xl bg-blue-100 flex items-center justify-center flex-shrink-0">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
            </div>
            <div>
              <div class="font-semibold text-slate-800">Dr. {{ appt.doctor_name || 'N/A' }}</div>
              <div class="text-sm text-slate-500">{{ appt.department_name || 'General' }}</div>
              <div class="flex items-center gap-3 mt-2 text-sm text-slate-600">
                <span>📅 {{ formatDate(appt.appointment_date) }}</span>
                <span>🕐 {{ appt.appointment_time?.slice(0,5) }}</span>
              </div>
              <p v-if="appt.reason" class="text-sm text-slate-500 mt-1">{{ appt.reason }}</p>
            </div>
          </div>
          <div class="flex flex-col items-start sm:items-end gap-3">
            <span :class="`badge-${appt.status}`">{{ appt.status }}</span>
            <div class="flex gap-2">
              <RouterLink v-if="isToday(appt.appointment_date) && ['pending','confirmed'].includes(appt.status)"
                to="/queue" class="btn-primary text-xs py-1.5 px-3">
                Check In →
              </RouterLink>
              <button v-if="['pending','confirmed'].includes(appt.status)"
                @click="cancelAppt(appt.appointment_id)" class="btn-danger text-xs py-1.5 px-3">
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAppointmentStore } from '@/stores/appointments'

const store     = useAppointmentStore()
const activeTab = ref('all')
const tabs = [
  { label: 'All',       value: 'all' },
  { label: 'Pending',   value: 'pending' },
  { label: 'Confirmed', value: 'confirmed' },
  { label: 'Completed', value: 'completed' },
  { label: 'Cancelled', value: 'cancelled' },
]

onMounted(() => store.fetchMyAppointments())

const filtered = computed(() =>
  activeTab.value === 'all' ? store.appointments : store.appointments.filter(a => a.status === activeTab.value)
)

function countByStatus(status) {
  return status === 'all' ? store.appointments.length : store.appointments.filter(a => a.status === status).length
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-KE', { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })
}

function isToday(dateStr) {
  return new Date(dateStr).toDateString() === new Date().toDateString()
}

async function cancelAppt(id) {
  if (!confirm('Cancel this appointment?')) return
  await store.cancelAppointment(id)
}
</script>

