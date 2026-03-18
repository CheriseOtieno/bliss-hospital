<template>
  <div class="min-h-screen flex flex-col">
    <Navbar v-if="!isAuthPage" />
    <main class="flex-1">
      <RouterView />
    </main>
    <Footer v-if="!isAuthPage && !isAdminPage" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Navbar from '@/components/Navbar.vue'
import Footer from '@/components/Footer.vue'

const route       = useRoute()
const isAuthPage  = computed(() => ['Login', 'Register'].includes(route.name))
const isAdminPage = computed(() => route.path.startsWith('/admin'))
</script>

