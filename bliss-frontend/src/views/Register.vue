<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-slate-100 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <RouterLink to="/" class="inline-flex items-center gap-2">
          <div class="w-10 h-10 bg-blue-600 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
            </svg>
          </div>
          <span class="font-display font-bold text-xl text-slate-800">Bliss Hospital</span>
        </RouterLink>
        <h1 class="text-2xl font-display font-bold text-slate-800 mt-6">Create your account</h1>
        <p class="text-slate-500 text-sm mt-1">Book appointments and track your queue — for free</p>
      </div>

      <div class="card shadow-lg border-0">
        <form @submit.prevent="handleRegister" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Full Name</label>
            <input v-model="form.full_name" type="text" class="input" placeholder="Jane Doe" required />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Email Address</label>
            <input v-model="form.email" type="email" class="input" placeholder="you@example.com" required />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Phone Number</label>
            <input v-model="form.phone" type="tel" class="input" placeholder="+254 7XX XXX XXX" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-1.5">Password</label>
            <input v-model="form.password" type="password" class="input" placeholder="At least 6 characters" required minlength="6" />
          </div>
          <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 rounded-xl px-4 py-3 text-sm">{{ error }}</div>
          <button type="submit" class="btn-primary w-full mt-2" :disabled="loading">
            {{ loading ? 'Creating account…' : 'Create Account' }}
          </button>
        </form>
        <p class="text-center text-sm text-slate-500 mt-6">
          Already have an account?
          <RouterLink to="/login" class="text-blue-600 font-medium hover:underline">Sign in</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth   = useAuthStore()
const router = useRouter()

const form    = ref({ full_name: '', email: '', phone: '', password: '' })
const loading = ref(false)
const error   = ref(null)

async function handleRegister() {
  loading.value = true
  error.value   = null
  try {
    await auth.register({ ...form.value, role: 'patient' })
    router.push('/appointments')
  } catch (e) {
    error.value = e.response?.data?.error || 'Registration failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

