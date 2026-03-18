import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

import Home             from '@/views/Home.vue'
import Login            from '@/views/Login.vue'
import Register         from '@/views/Register.vue'
import BookAppointment  from '@/views/BookAppointment.vue'
import MyAppointments   from '@/views/MyAppointments.vue'
import QueueStatus      from '@/views/QueueStatus.vue'
import Profile          from '@/views/Profile.vue'
import AdminDashboard   from '@/views/admin/Dashboard.vue'
import AdminDoctors     from '@/views/admin/Doctors.vue'
import AdminAppointments from '@/views/admin/Appointments.vue'
import AdminQueue       from '@/views/admin/QueueOverview.vue'

const routes = [
  { path: '/',          name: 'Home',     component: Home },
  { path: '/login',     name: 'Login',    component: Login },
  { path: '/register',  name: 'Register', component: Register },

  { path: '/book',         name: 'Book',         component: BookAppointment,  meta: { requiresAuth: true } },
  { path: '/appointments', name: 'Appointments', component: MyAppointments,   meta: { requiresAuth: true } },
  { path: '/queue',        name: 'Queue',        component: QueueStatus,      meta: { requiresAuth: true } },
  { path: '/profile',      name: 'Profile',      component: Profile,          meta: { requiresAuth: true } },

  { path: '/admin',              name: 'AdminDashboard',    component: AdminDashboard,    meta: { requiresAuth: true, adminOnly: true } },
  { path: '/admin/doctors',      name: 'AdminDoctors',      component: AdminDoctors,      meta: { requiresAuth: true, adminOnly: true } },
  { path: '/admin/appointments', name: 'AdminAppointments', component: AdminAppointments, meta: { requiresAuth: true, adminOnly: true } },
  { path: '/admin/queue',        name: 'AdminQueue',        component: AdminQueue,        meta: { requiresAuth: true, staffOnly: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isLoggedIn) return next({ name: 'Login', query: { redirect: to.fullPath } })
  if (to.meta.adminOnly  && !auth.isAdmin)    return next({ name: 'Home' })
  if (to.meta.staffOnly  && !auth.isStaff)    return next({ name: 'Home' })
  next()
})

export default router

