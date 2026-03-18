<template>
  <div>
    <section class="bg-gradient-to-br from-blue-700 to-blue-900 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 flex flex-col lg:flex-row items-center gap-12">
        <div class="flex-1 text-center lg:text-left">
          <span class="inline-block bg-blue-500/30 text-blue-200 text-xs font-semibold px-3 py-1 rounded-full mb-4 uppercase tracking-wide">
            Smart Healthcare · Bliss Hospitals
          </span>
          <h1 class="text-4xl lg:text-5xl font-display font-bold leading-tight mb-5">
            Skip the Queue.<br>Book Your Visit<br>
            <span class="text-blue-300">in 60 Seconds.</span>
          </h1>
          <p class="text-blue-100 text-lg mb-8 max-w-xl">
            Book appointments online, track your queue position in real time, and receive instant notifications — all in one place.
          </p>
          <div class="flex flex-col sm:flex-row gap-3 justify-center lg:justify-start">
            <RouterLink to="/register" class="bg-white text-blue-700 font-semibold px-6 py-3 rounded-xl hover:bg-blue-50 transition text-center">
              Get Started — It's Free
            </RouterLink>
            <RouterLink to="/login" class="border border-blue-400 text-white font-medium px-6 py-3 rounded-xl hover:bg-blue-800 transition text-center">
              Sign In
            </RouterLink>
          </div>
        </div>
        <div class="flex-1 grid grid-cols-2 gap-4 w-full max-w-sm">
          <div v-for="stat in stats" :key="stat.label"
            class="bg-white/10 backdrop-blur rounded-2xl p-5 text-center border border-white/20">
            <div class="text-3xl font-display font-bold text-white">{{ stat.value }}</div>
            <div class="text-blue-200 text-sm mt-1">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="text-center mb-14">
        <h2 class="text-3xl font-display font-bold text-slate-800">Everything You Need</h2>
        <p class="text-slate-500 mt-2">A smarter way to manage your hospital visits.</p>
      </div>
      <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="feature in features" :key="feature.title" class="card hover:shadow-md transition-shadow">
          <div class="w-12 h-12 rounded-xl flex items-center justify-center mb-4" :class="feature.bg">
            <span class="text-2xl">{{ feature.icon }}</span>
          </div>
          <h3 class="font-display font-semibold text-slate-800 mb-2">{{ feature.title }}</h3>
          <p class="text-slate-500 text-sm leading-relaxed">{{ feature.desc }}</p>
        </div>
      </div>
    </section>

    <section class="bg-blue-600 text-white py-16">
      <div class="max-w-3xl mx-auto px-4 text-center">
        <h2 class="text-3xl font-display font-bold mb-4">Ready to experience smarter healthcare?</h2>
        <p class="text-blue-100 mb-8">Join thousands of patients at Bliss Hospital already saving time.</p>
        <RouterLink to="/register" class="bg-white text-blue-700 font-semibold px-8 py-3 rounded-xl hover:bg-blue-50 transition inline-block">
          Create Your Account
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<script setup>
const stats = [
  { value: '10+',   label: 'Departments' },
  { value: '<5min', label: 'Avg. Booking Time' },
  { value: '24/7',  label: 'Online Access' },
  { value: '100%',  label: 'Paperless' },
]
const features = [
  { icon: '📅', title: 'Online Booking',      bg: 'bg-blue-100',   desc: 'Book appointments with any specialist across all Bliss branches — any time, from anywhere.' },
  { icon: '🔢', title: 'Live Queue Tracking', bg: 'bg-purple-100', desc: 'See your real-time queue position on your phone. No more guessing how long you will wait.' },
  { icon: '🔔', title: 'Smart Notifications', bg: 'bg-yellow-100', desc: 'Get email alerts for confirmations, reminders, and when it is your turn.' },
  { icon: '🏥', title: 'Multi-Department',    bg: 'bg-green-100',  desc: 'From General Medicine to Cardiology — choose from 10 specialist departments.' },
  { icon: '📊', title: 'Admin Dashboard',     bg: 'bg-red-100',    desc: 'Staff get powerful tools to manage queues, schedules, and patient flow in real time.' },
  { icon: '🔒', title: 'Secure & Private',    bg: 'bg-teal-100',   desc: 'Your health data is protected with JWT authentication and encrypted storage.' },
]
</script>

