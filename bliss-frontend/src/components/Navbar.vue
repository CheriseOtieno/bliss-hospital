<template>
  <nav class="bg-white border-b border-slate-100 sticky top-0 z-50 shadow-sm">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16 items-center">

        <RouterLink to="/" class="flex items-center gap-2">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
            </svg>
          </div>
          <span class="font-display font-bold text-lg text-slate-800">Bliss <span class="text-blue-600">Hospital</span></span>
        </RouterLink>

        <div class="hidden md:flex items-center gap-1">
          <template v-if="auth.isLoggedIn">
            <RouterLink v-if="auth.isStaff" to="/admin"
              class="px-4 py-2 rounded-lg text-sm font-medium text-slate-600 hover:bg-slate-100 transition">
              Dashboard
            </RouterLink>
            <template v-else>
              <RouterLink to="/book"         class="px-4 py-2 rounded-lg text-sm font-medium text-slate-600 hover:bg-slate-100 transition">Book Appointment</RouterLink>
              <RouterLink to="/appointments" class="px-4 py-2 rounded-lg text-sm font-medium text-slate-600 hover:bg-slate-100 transition">My Appointments</RouterLink>
              <RouterLink to="/queue"        class="px-4 py-2 rounded-lg text-sm font-medium text-slate-600 hover:bg-slate-100 transition">Queue Status</RouterLink>
            </template>
          </template>
        </div>

        <div class="flex items-center gap-3">
          <template v-if="auth.isLoggedIn">
            <RouterLink to="/profile" class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-slate-100 hover:bg-slate-200 transition">
              <div class="w-7 h-7 rounded-full bg-blue-600 flex items-center justify-center text-white text-xs font-bold">
                {{ initials }}
              </div>
              <span class="text-sm font-medium text-slate-700 hidden sm:block">{{ auth.user?.full_name?.split(' ')[0] }}</span>
            </RouterLink>
            <button @click="handleLogout" class="btn-secondary text-sm py-2">Logout</button>
          </template>
          <template v-else>
            <RouterLink to="/login"    class="btn-secondary text-sm py-2">Login</RouterLink>
            <RouterLink to="/register" class="btn-primary  text-sm py-2">Register</RouterLink>
          </template>
        </div>

      </div>
    </div>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth   = useAuthStore()
const router = useRouter()

const initials = computed(() =>
  (auth.user?.full_name || '').split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
)

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

