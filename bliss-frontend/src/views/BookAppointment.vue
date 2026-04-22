<template>
  <div class="max-w-3xl mx-auto px-4 py-10">
    <div class="mb-8">
      <h1 class="text-3xl font-display font-bold text-slate-800">Book an Appointment</h1>
      <p class="text-slate-500 mt-1">Follow the steps below to schedule your visit.</p>
    </div>

    <!-- Step indicator -->
    <div class="flex items-center gap-2 mb-10">
      <div v-for="(s, i) in steps" :key="i" class="flex items-center gap-2">
        <div class="flex items-center gap-1.5">
          <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold transition-all"
            :class="step > i+1 ? 'bg-green-500 text-white' : step === i+1 ? 'bg-blue-600 text-white' : 'bg-slate-200 text-slate-400'">
            <span v-if="step > i+1">✓</span>
            <span v-else>{{ i + 1 }}</span>
          </div>
          <span class="text-sm font-medium hidden sm:block"
            :class="step === i+1 ? 'text-blue-600' : step > i+1 ? 'text-green-600' : 'text-slate-400'">
            {{ s }}
          </span>
        </div>
        <div v-if="i < steps.length - 1" class="flex-1 h-0.5 w-6 sm:w-10"
          :class="step > i+1 ? 'bg-green-400' : 'bg-slate-200'" />
      </div>
    </div>

    <!-- Step 1: Department -->
    <div v-if="step === 1" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-6">Select Department</h2>
      <div v-if="store.departments.length === 0" class="text-slate-400 text-center py-8">Loading departments…</div>
      <div class="grid sm:grid-cols-2 gap-3">
        <button v-for="dept in store.departments" :key="dept.department_id"
          @click="selectDepartment(dept)"
          class="text-left p-4 rounded-xl border-2 transition-all hover:border-blue-400 hover:bg-blue-50"
          :class="selected.department?.department_id === dept.department_id ? 'border-blue-600 bg-blue-50' : 'border-slate-200 bg-white'">
          <div class="font-medium text-slate-800">{{ dept.name }}</div>
          <div class="text-xs text-slate-500 mt-0.5">{{ dept.description }}</div>
        </button>
      </div>
    </div>

    <!-- Step 2: Doctor -->
    <div v-if="step === 2" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-1">Select a Doctor</h2>
      <p class="text-slate-500 text-sm mb-6">{{ selected.department?.name }} specialists</p>
      <div v-if="store.doctors.length === 0" class="text-center py-8 text-slate-400">No doctors available in this department.</div>
      <div class="space-y-3">
        <button v-for="doc in store.doctors" :key="doc.doctor_id"
          @click="selectDoctor(doc)"
          class="w-full text-left p-4 rounded-xl border-2 transition-all hover:border-blue-400"
          :class="selected.doctor?.doctor_id === doc.doctor_id ? 'border-blue-600 bg-blue-50' : 'border-slate-200 bg-white'">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center font-bold text-blue-700">
              {{ doc.full_name.split(' ').map(n => n[0]).join('').slice(0,2) }}
            </div>
            <div>
              <div class="font-medium text-slate-800">Dr. {{ doc.full_name }}</div>
              <div class="text-xs text-slate-500">{{ doc.specialty }}</div>
            </div>
          </div>
          <p v-if="doc.bio" class="text-xs text-slate-500 mt-2">{{ doc.bio }}</p>
        </button>
      </div>
      <button @click="step--" class="btn-secondary mt-4">← Back</button>
    </div>

    <!-- Step 3: Date & Slot -->
    <div v-if="step === 3" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-6">Choose Date & Time</h2>
      <div class="mb-5">
        <label class="block text-sm font-medium text-slate-700 mb-1.5">Select Date</label>
        <input v-model="selected.date" type="date" class="input max-w-xs" :min="today" @change="loadSlots" />
      </div>
      <div v-if="selected.date">
        <p class="text-sm font-medium text-slate-700 mb-3">Available Time Slots</p>
        <div v-if="store.slots.length === 0" class="text-slate-400 text-sm py-4">No slots available on this date. Try another date.</div>
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
          <button v-for="slot in store.slots" :key="slot.slot_id"
            @click="selected.slot = slot"
            class="py-2.5 px-3 rounded-xl border-2 text-sm font-medium transition-all"
            :class="selected.slot?.slot_id === slot.slot_id ? 'border-blue-600 bg-blue-600 text-white' : 'border-slate-200 hover:border-blue-400 text-slate-700'">
            {{ slot.start_time.slice(0,5) }} – {{ slot.end_time.slice(0,5) }}
          </button>
        </div>
      </div>
      <div class="flex gap-3 mt-6">
        <button @click="step--" class="btn-secondary">← Back</button>
        <button @click="step++" :disabled="!selected.slot" class="btn-primary">Continue →</button>
      </div>
    </div>

    <!-- Step 4: Confirm -->
    <div v-if="step === 4" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-6">Confirm Booking</h2>
      <div class="bg-slate-50 rounded-xl p-5 space-y-3 mb-6">
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Branch</span>
          <span class="font-medium text-slate-800">{{ selected.branch?.name }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Department</span>
          <span class="font-medium text-slate-800">{{ selected.department?.name }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Doctor</span>
          <span class="font-medium text-slate-800">Dr. {{ selected.doctor?.full_name }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Date</span>
          <span class="font-medium text-slate-800">{{ formatDate(selected.date) }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Time</span>
          <span class="font-medium text-slate-800">{{ selected.slot?.start_time?.slice(0,5) }} – {{ selected.slot?.end_time?.slice(0,5) }}</span>
        </div>
      </div>
      <div class="mb-5">
        <label class="block text-sm font-medium text-slate-700 mb-1.5">Reason for visit (optional)</label>
        <textarea v-model="selected.reason" class="input" rows="3" placeholder="Briefly describe your symptoms or reason…" />
      </div>
      <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 rounded-xl px-4 py-3 text-sm mb-4">{{ error }}</div>
      <div class="flex gap-3">
        <button @click="step--" class="btn-secondary">← Back</button>
        <button @click="confirmBooking" :disabled="loading" class="btn-primary flex-1">
          {{ loading ? 'Booking…' : 'Confirm Appointment' }}
        </button>
      </div>
    </div>

    <!-- Step 5: Success -->
    <div v-if="step === 5" class="card text-center py-10">
      <div class="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-5">
        <svg class="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
      </div>
      <h2 class="text-2xl font-display font-bold text-slate-800 mb-2">Appointment Booked!</h2>
      <p class="text-slate-500 mb-2">Your appointment with <strong>Dr. {{ selected.doctor?.full_name }}</strong> is confirmed.</p>
      <p class="text-slate-500 text-sm mb-8">A confirmation notification has been sent to you.</p>
      <div class="flex gap-3 justify-center">
        <RouterLink to="/appointments" class="btn-primary">View My Appointments</RouterLink>
        <button @click="resetForm" class="btn-secondary">Book Another</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAppointmentStore } from '@/stores/appointments'

const store   = useAppointmentStore()
const step    = ref(1)
const loading = ref(false)
const error   = ref(null)
const today   = new Date().toISOString().split('T')[0]
const steps   = ['Department', 'Doctor', 'Date & Time', 'Confirm']
const selected = ref({ department: null, doctor: null, date: '', slot: null, reason: '' })

onMounted(() => store.fetchDepartments())

async function selectDepartment(dept) {
  selected.value.department = dept
  selected.value.doctor = null
  selected.value.slot   = null
  await store.fetchDoctors(dept.department_id)
  step.value = 2
}

function selectDoctor(doc) {
  selected.value.doctor = doc
  selected.value.slot   = null
  step.value = 3
}

async function loadSlots() {
  console.log("LOAD SLOTS TRIGGERED")

  if (!selected.value.date) return

  console.log("Doctor:", selected.value.doctor?.doctor_id)
  console.log("Date:", selected.value.date)

  await store.fetchSlots(
    selected.value.doctor.doctor_id,
    selected.value.date
  )
}

async function confirmBooking() {
  loading.value = true
  error.value   = null
  try {
    await store.bookAppointment({
      doctor_id:     selected.value.doctor.doctor_id,
      department_id: selected.value.department.department_id,
      slot_id:       selected.value.slot.slot_id,
      reason:        selected.value.reason,
    })
    step.value = 5
  } catch (e) {
    error.value = e.response?.data?.error || 'Booking failed. Please try again.'
  } finally {
    loading.value = false
  }
}

function resetForm() {
  step.value = 1
  selected.value = { department: null, doctor: null, date: '', slot: null, reason: '' }
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-KE', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

