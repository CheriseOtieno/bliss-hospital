<template>
  <div class="max-w-2xl mx-auto px-4 py-10">
    <div class="mb-8">
      <h1 class="text-3xl font-display font-bold text-slate-800">Queue Status</h1>
      <p class="text-slate-500 mt-1">Check in and track your real-time position.</p>
    </div>

    <!-- Check-in panel (shown when not yet in queue) -->
    <div v-if="!queueStore.myQueue" class="card mb-6">
      <h2 class="font-display font-semibold text-lg text-slate-800 mb-4">Check In for Today's Appointment</h2>

      <div v-if="todayAppointments.length === 0" class="text-center py-8">
        <div class="text-4xl mb-3">📋</div>
        <p class="text-slate-500 mb-4">No appointments scheduled for today.</p>
        <RouterLink to="/book" class="btn-primary">Book an Appointment</RouterLink>
      </div>

      <div v-else class="space-y-3">
        <div v-for="appt in todayAppointments" :key="appt.appointment_id"
          class="p-4 rounded-xl border-2 border-slate-200 flex items-center justify-between gap-4">
          <div>
            <div class="font-medium text-slate-800">Dr. {{ appt.doctor_name }}</div>
            <div class="text-sm text-slate-500">{{ appt.department_name }} · {{ appt.appointment_time?.slice(0,5) }}</div>
          </div>
          <button @click="doCheckIn(appt.appointment_id)" :disabled="checking" class="btn-primary text-sm py-2 whitespace-nowrap">
            {{ checking ? 'Checking in…' : 'Check In' }}
          </button>
        </div>
      </div>

      <div v-if="checkInError" class="bg-red-50 border border-red-200 text-red-700 rounded-xl px-4 py-3 text-sm mt-4">
        {{ checkInError }}
      </div>
    </div>

    <!-- Live Queue Card -->
    <div v-if="queueStore.myQueue" class="card mb-6">
      <div class="text-center">
        <div class="inline-block mb-4">
          <span :class="`badge-${queueStore.myQueue.status} text-sm px-4 py-1.5`">
            {{ statusLabel(queueStore.myQueue.status) }}
          </span>
        </div>

        <div class="w-32 h-32 rounded-full border-4 flex items-center justify-center mx-auto mb-4"
          :class="queueStore.myQueue.status === 'called' ? 'border-purple-500 bg-purple-50 animate-pulse' : 'border-blue-500 bg-blue-50'">
          <div>
            <div class="text-xs font-medium text-slate-500 mb-1">Queue No.</div>
            <div class="text-4xl font-display font-bold"
              :class="queueStore.myQueue.status === 'called' ? 'text-purple-700' : 'text-blue-700'">
              {{ queueStore.myQueue.queue_number }}
            </div>
          </div>
        </div>

        <div v-if="queueStore.myQueue.status === 'called'"
          class="bg-purple-100 border border-purple-300 rounded-xl p-4 mb-4 text-purple-800">
          <div class="font-bold text-lg">🔔 It's Your Turn!</div>
          <p class="text-sm mt-1">Please proceed to the consultation room now.</p>
        </div>

        <div v-else class="mb-4">
          <p class="text-2xl font-display font-semibold text-slate-700">
            {{ queueStore.myQueue.position === 1 ? 'You are next!' : `${queueStore.myQueue.position - 1} patient(s) ahead` }}
          </p>
          <p class="text-slate-500 text-sm mt-1">Doctor: Dr. {{ queueStore.myQueue.doctor_name }}</p>
        </div>

        <div class="grid grid-cols-3 gap-3 text-center mt-6">
          <div class="bg-slate-50 rounded-xl p-3">
            <div class="text-xs text-slate-500 mb-1">Checked In</div>
            <div class="text-sm font-semibold">{{ formatTime(queueStore.myQueue.checked_in_at) }}</div>
          </div>
          <div class="bg-slate-50 rounded-xl p-3">
            <div class="text-xs text-slate-500 mb-1">Position</div>
            <div class="text-sm font-semibold">#{{ queueStore.myQueue.position }}</div>
          </div>
          <div class="bg-slate-50 rounded-xl p-3">
            <div class="text-xs text-slate-500 mb-1">Est. Wait</div>
            <div class="text-sm font-semibold">~{{ (queueStore.myQueue.position - 1) * 15 }} min</div>
          </div>
        </div>
        <p class="text-xs text-slate-400 mt-4">Queue updates automatically every 30 seconds.</p>
      </div>
    </div>

    <div v-if="queueStore.myQueue" class="text-center">
      <button @click="queueStore.fetchMyQueue()" class="btn-secondary text-sm">🔄 Refresh Status</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useQueueStore } from '@/stores/queue'
import { useAppointmentStore } from '@/stores/appointments'

const queueStore = useQueueStore()
const apptStore  = useAppointmentStore()
const checking     = ref(false)
const checkInError = ref(null)
let refreshInterval = null

const todayAppointments = computed(() =>
  (apptStore.appointments || []).filter(a => {
    const isToday = new Date(a.appointment_date).toDateString() === new Date().toDateString()
    return isToday && ['pending', 'confirmed'].includes(a.status)
  })
)

onMounted(async () => {
  await apptStore.fetchMyAppointments()
  await queueStore.fetchMyQueue()
  refreshInterval = setInterval(() => queueStore.fetchMyQueue(), 30000)
})

onUnmounted(() => clearInterval(refreshInterval))

async function doCheckIn(appointmentId) {
  checking.value     = true
  checkInError.value = null
  try {
    await queueStore.checkIn(appointmentId)
  } catch (e) {
    checkInError.value = e.response?.data?.error || 'Check-in failed. Please try again.'
  } finally {
    checking.value = false
  }
}

function statusLabel(status) {
  const map = { waiting: 'Waiting', called: 'Called — Go Now!', serving: 'Being Served', done: 'Done' }
  return map[status] || status
}

function formatTime(ts) {
  if (!ts) return '—'
  return new Date(ts).toLocaleTimeString('en-KE', { hour: '2-digit', minute: '2-digit' })
}
</script>

