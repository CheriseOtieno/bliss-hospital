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

    <!-- STEP 1: BRANCH -->
    <div v-if="step === 1" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-6">Select Branch</h2>

      <div v-if="store.branches.length === 0" class="text-center py-8 text-slate-400">Loading branches…</div>

      <div class="grid sm:grid-cols-2 gap-3">
        <button
          v-for="branch in store.branches"
          :key="branch.branch_id"
          @click="selectBranch(branch)"
          class="text-left p-4 rounded-xl border-2 transition-all"
          :class="selected.branch?.branch_id === branch.branch_id ? 'border-blue-600 bg-blue-50' : 'border-slate-200 hover:border-blue-400'"
        >
          <div class="font-medium text-slate-800">{{ branch.name }}</div>
          <div class="text-xs text-slate-500 mt-0.5">{{ branch.branch_name }}</div>
        </button>
      </div>
    </div>

    <!-- STEP 2: DEPARTMENT -->
    <div v-if="step === 2" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-6">Select Department</h2>

      <div v-if="store.departments.length === 0" class="text-center py-8 text-slate-400">Loading departments…</div>

      <div class="grid sm:grid-cols-2 gap-3">
        <button
          v-for="dept in store.departments"
          :key="dept.department_id"
          @click="selectDepartment(dept)"
          class="text-left p-4 rounded-xl border-2 transition-all"
          :class="selected.department?.department_id === dept.department_id ? 'border-blue-600 bg-blue-50' : 'border-slate-200 hover:border-blue-400'"
        >
          <div class="font-medium text-slate-800">{{ dept.department_name }}</div>
          <div class="text-xs text-slate-500 mt-0.5">{{ dept.description }}</div>
        </button>
      </div>

      <button @click="step--" class="btn-secondary mt-4">← Back</button>
    </div>

    <!-- STEP 3: DOCTOR -->
    <div v-if="step === 3" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-1">Select a Doctor</h2>
      <p class="text-slate-500 text-sm mb-6">{{ selected.department?.department_name }} specialists</p>

      <div v-if="store.doctors.length === 0" class="text-center py-8 text-slate-400">
        No doctors available in this department.
      </div>

      <div class="space-y-3">
        <button
          v-for="doc in store.doctors"
          :key="doc.doctor_id"
          @click="selectDoctor(doc)"
          class="w-full text-left p-4 rounded-xl border-2 transition-all"
          :class="selected.doctor?.doctor_id === doc.doctor_id ? 'border-blue-600 bg-blue-50' : 'border-slate-200 hover:border-blue-400'"
        >
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

    <!-- STEP 4: DATE & TIME -->
    <div v-if="step === 4" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-6">Choose Date & Time</h2>

      <div class="mb-5">
        <label class="block text-sm font-medium text-slate-700 mb-1.5">Select Date</label>
        <input v-model="selected.date" type="date" class="input max-w-xs" :min="today" @change="loadSlots" />
      </div>

      <div v-if="selected.date">
        <p class="text-sm font-medium text-slate-700 mb-3">Available Time Slots</p>

        <div v-if="store.slots.length === 0" class="text-slate-400 text-sm py-4">
          No slots available on this date.
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
          <button
            v-for="slot in store.slots"
            :key="slot.slot_id"
            @click="selected.slot = slot"
            class="py-2.5 px-3 rounded-xl border-2 text-sm font-medium transition-all"
            :class="selected.slot?.slot_id === slot.slot_id ? 'border-blue-600 bg-blue-600 text-white' : 'border-slate-200 hover:border-blue-400 text-slate-700'"
          >
            {{ slot.start_time.slice(0,5) }} – {{ slot.end_time.slice(0,5) }}
          </button>
        </div>
      </div>

      <div class="flex gap-3 mt-6">
        <button @click="step--" class="btn-secondary">← Back</button>
        <button @click="step++" :disabled="!selected.slot" class="btn-primary">Continue →</button>
      </div>
    </div>

    <!-- STEP 5: CONFIRM -->
    <div v-if="step === 5" class="card">
      <h2 class="font-display font-semibold text-xl text-slate-800 mb-6">Confirm Booking</h2>

      <div class="bg-slate-50 rounded-xl p-5 space-y-3 mb-6">
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Branch</span>
          <span class="font-medium text-slate-800">{{ selected.branch?.name }}</span>
        </div>

        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Department</span>
          <span class="font-medium text-slate-800">{{ selected.department?.department_name }}</span>
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
          <span class="font-medium text-slate-800">
            {{ selected.slot?.start_time.slice(0,5) }} – {{ selected.slot?.end_time.slice(0,5) }}
          </span>
        </div>
      </div>

      <div class="mb-5">
        <label class="block text-sm font-medium text-slate-700 mb-1.5">Reason (optional)</label>
        <textarea v-model="selected.reason" class="input" rows="3" placeholder="Briefly describe your symptoms…" />
      </div>

      <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 rounded-xl px-4 py-3 text-sm mb-4">
        {{ error }}
      </div>

      <div class="flex gap-3">
        <button @click="step--" class="btn-secondary">← Back</button>
        <button @click="confirmBooking" :disabled="loading" class="btn-primary flex-1">
          {{ loading ? 'Booking…' : 'Confirm Appointment' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAppointmentStore } from '@/stores/appointments'

const store = useAppointmentStore()
const step = ref(1)
const loading = ref(false)
const error = ref(null)
const today = new Date().toISOString().split('T')[0]

const steps = ['Branch', 'Department', 'Doctor', 'Date & Time', 'Confirm']

const selected = ref({
  branch: null,
  department: null,
  doctor: null,
  date: '',
  slot: null,
  reason: ''
})

onMounted(() => {
  store.fetchBranches()
})

async function selectBranch(branch) {
  selected.value.branch = branch
  selected.value.department = null
  selected.value.doctor = null
  selected.value.slot = null

  await store.fetchDepartments(branch.branch_id)
  step.value = 2
}

async function selectDepartment(dept) {
  selected.value.department = dept
  selected.value.doctor = null
  selected.value.slot = null

  await store.fetchDoctors(dept.department_id)
  step.value = 3
}

function selectDoctor(doc) {
  selected.value.doctor = doc
  selected.value.slot = null
  step.value = 4
}

async function loadSlots() {
  if (!selected.value.date) return
  await store.fetchSlots(selected.value.doctor.doctor_id, selected.value.date)
}

async function confirmBooking() {
  loading.value = true
  error.value = null

  try {
    await store.bookAppointment({
      doctor_id: selected.value.doctor.doctor_id,
      department_id: selected.value.department.department_id,
      branch_id: selected.value.branch.branch_id,
      slot_id: selected.value.slot.slot_id,
      reason: selected.value.reason
    })

    step.value = 6
  } catch (e) {
    error.value = e.response?.data?.error || 'Booking failed. Try again.'
  } finally {
    loading.value = false
  }
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-KE', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>                                                                                                                                                                                                                          