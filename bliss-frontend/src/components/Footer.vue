<template>
  <footer class="bg-white border-t border-slate-100 mt-auto py-8">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col md:flex-row justify-between items-center gap-4">
      <div class="flex items-center gap-2">
        <div class="w-6 h-6 bg-blue-600 rounded-md flex items-center justify-center">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
          </svg>
        </div>
        <span class="text-sm font-medium text-slate-600">Bliss Hospital · Smart Care, Smart Queues</span>
      </div>
      <p class="text-xs text-slate-400">© {{ new Date().getFullYear() }} Bliss Hospital. All rights reserved.</p>
    </div>
  </footer>
</template>

