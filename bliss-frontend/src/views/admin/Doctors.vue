<template>
  <div class="flex min-h-screen bg-slate-50">
    <AdminSidebar />
    <main class="flex-1 p-8">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-display font-bold text-slate-800">Doctors</h1>
          <p class="text-slate-500 mt-1">Manage specialist doctors and availability.</p>
        </div>
        <button @click="showModal = true" class="btn-primary">+ Add Doctor</button>
      </div>

      <div v-if="loading" class="text-center py-16 text-slate-400">Loading doctors…</div>

      <div v-else class="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
        <div v-for="doc in doctors" :key="doc.doctor_id" class="card hover:shadow-md transition-shadow">
          <div class="flex items-start gap-3 mb-3">
            <div class="w-12 h-12 rounded-xl bg-blue-100 flex items-center justify-center font-bold text-blue-700 text-lg flex-shrink-0">
              {{ doc.full_name.split(' ').map(n => n[0]).join('').slice(0,2) }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="font-semibold text-slate-800">Dr. {{ doc.full_name }}</div>
              <div class="text-sm text-slate-500">{{ doc.specialty }}</div>
              <div class="text-xs text-blue-600 mt-0.5">{{ doc.department }}</div>
            </div>
            <span class="text-xs px-2 py-1 rounded-full font-medium"
              :class="doc.available ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-600'">
              {{ doc.available ? 'Available' : 'Unavailable' }}
            </span>
          </div>
          <p v-if="doc.bio" class="text-xs text-slate-500 mb-3">{{ doc.bio }}</p>
          <div class="flex gap-2 pt-2 border-t border-slate-100">
            <button @click="toggleAvailability(doc)"
              class="flex-1 text-xs py-1.5 rounded-lg font-medium transition"
              :class="doc.available ? 'bg-red-50 text-red-600 hover:bg-red-100' : 'bg-green-50 text-green-600 hover:bg-green-100'">
              {{ doc.available ? 'Set Unavailable' : 'Set Available' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Add Doctor Modal -->
      <div v-if="showModal" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-2xl shadow-xl w-full max-w-md p-6">
          <h2 class="font-display font-semibold text-xl text-slate-800 mb-5">Add New Doctor</h2>
          <form @submit.prevent="createDoctor" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">Full Name</label>
              <input v-model="form.full_name" class="input" placeholder="John Kamau" required />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">Specialty</label>
              <input v-model="form.specialty" class="input" placeholder="e.g. Cardiologist" required />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">Department</label>
              <select v-model="form.department_id" class="input" required>
                <option value="">Select department…</option>
                <option v-for="d in departments" :key="d.department_id" :value="d.department_id">{{ d.name }}</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-700 mb-1.5">Bio (optional)</label>
              <textarea v-model="form.bio" class="input" rows="2" placeholder="Brief description…" />
            </div>
            <div v-if="formError" class="bg-red-50 border border-red-200 text-red-700 rounded-xl px-4 py-3 text-sm">{{ formError }}</div>
            <div class="flex gap-3 pt-2">
              <button type="button" @click="closeModal" class="btn-secondary flex-1">Cancel</button>
              <button type="submit" :disabled="saving" class="btn-primary flex-1">{{ saving ? 'Saving…' : 'Add Doctor' }}</button>
            </div>
          </form>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AdminSidebar from '@/components/AdminSidebar.vue'
import api from '@/services/api'

const doctors     = ref([])
const departments = ref([])
const loading     = ref(false)
const showModal   = ref(false)
const saving      = ref(false)
const formError   = ref(null)
const form        = ref({ full_name: '', specialty: '', department_id: '', bio: '' })

onMounted(async () => {
  loading.value = true
  try {
    const [docRes, deptRes] = await Promise.all([api.get('/doctors'), api.get('/departments')])
    doctors.value     = docRes.data  || []
    departments.value = deptRes.data || []
  } finally {
    loading.value = false
  }
})

async function createDoctor() {
  saving.value    = true
  formError.value = null
  try {
    const { data } = await api.post('/admin/doctors', form.value)
    doctors.value.unshift(data.data)
    closeModal()
  } catch (e) {
    formError.value = e.response?.data?.error || 'Failed to add doctor.'
  } finally {
    saving.value = false
  }
}

async function toggleAvailability(doc) {
  await api.patch(`/admin/doctors/${doc.doctor_id}`, { available: !doc.available, specialty: doc.specialty, bio: doc.bio })
  doc.available = !doc.available
}

function closeModal() {
  showModal.value = false
  form.value = { full_name: '', specialty: '', department_id: '', bio: '' }
  formError.value = null
}
</script>

