<template>
  <aside class="w-64 bg-slate-900 text-white min-h-screen flex flex-col flex-shrink-0">
    <div class="px-6 py-5 border-b border-slate-800">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
          </svg>
        </div>
        <div>
          <div class="font-display font-bold text-sm text-white">Bliss Hospital</div>
          <div class="text-xs text-slate-400">Admin Panel</div>
        </div>
      </div>
    </div>

    <nav class="flex-1 px-3 py-4 space-y-1">
      <RouterLink v-for="item in navItems" :key="item.to" :to="item.to"
        class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all"
        :class="$route.path === item.to ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'">
        <span class="text-lg">{{ item.icon }}</span>
        {{ item.label }}
      </RouterLink>
    </nav>

    <div class="px-4 py-4 border-t border-slate-800">
      <div class="flex items-center gap-3 mb-3">
        <div class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-xs font-bold">
          {{ initials }}
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium text-white truncate">{{ auth.user?.full_name }}</div>
          <div class="text-xs text-slate-400 capitalize">{{ auth.user?.role }}</div>
        </div>
      </div>
      <button @click="handleLogout" class="w-full text-left text-xs text-slate-400 hover:text-white transition px-1 py-1">
        → Sign out
      </button>
    </div>
  </aside>
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

const navItems = [
  { to: '/admin',              icon: '📊', label: 'Dashboard' },
  { to: '/admin/queue',        icon: '🔢', label: 'Queue' },
  { to: '/admin/appointments', icon: '📅', label: 'Appointments' },
  { to: '/admin/doctors',      icon: '🩺', label: 'Doctors' },
]

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

