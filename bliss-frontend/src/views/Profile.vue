<template>
  <div class="max-w-2xl mx-auto px-4 py-10">
    <h1 class="text-3xl font-display font-bold text-slate-800 mb-8">My Profile</h1>

    <div class="card mb-6">
      <div class="flex items-center gap-5 mb-6">
        <div class="w-16 h-16 rounded-full bg-blue-600 flex items-center justify-center text-white text-xl font-bold">
          {{ initials }}
        </div>
        <div>
          <div class="font-display font-semibold text-xl text-slate-800">{{ auth.user?.full_name }}</div>
          <div class="text-slate-500 text-sm">{{ auth.user?.email }}</div>
          <span class="mt-1 inline-block px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-700 capitalize">
            {{ auth.user?.role }}
          </span>
        </div>
      </div>

      <div class="space-y-3 border-t border-slate-100 pt-5">
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Full Name</span>
          <span class="font-medium text-slate-800">{{ auth.user?.full_name }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Email</span>
          <span class="font-medium text-slate-800">{{ auth.user?.email }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Phone</span>
          <span class="font-medium text-slate-800">{{ auth.user?.phone || 'Not provided' }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-slate-500">Account type</span>
          <span class="font-medium text-slate-800 capitalize">{{ auth.user?.role }}</span>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4 mb-6">
      <RouterLink to="/appointments" class="card hover:shadow-md transition-shadow text-center py-6">
        <div class="text-3xl mb-2">📅</div>
        <div class="font-semibold text-slate-700">My Appointments</div>
      </RouterLink>
      <RouterLink to="/queue" class="card hover:shadow-md transition-shadow text-center py-6">
        <div class="text-3xl mb-2">🔢</div>
        <div class="font-semibold text-slate-700">Queue Status</div>
      </RouterLink>
    </div>

    <button @click="handleLogout" class="btn-danger w-full">Sign Out</button>
  </div>
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

