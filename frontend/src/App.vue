<script setup lang="ts">
import { RouterView, RouterLink } from 'vue-router'
</script>

<template>
  <div class="min-h-screen bg-gray-50 font-sans">
    <nav class="bg-indigo-600 text-white shadow-lg">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <RouterLink to="/" class="text-xl font-bold tracking-tight hover:text-indigo-100 uppercase">
              Stock Tracker
            </RouterLink>
            <div class="ml-10 flex items-baseline space-x-4">
              <RouterLink to="/" class="px-3 py-2 rounded-md text-sm font-medium hover:bg-indigo-500 transition-colors" active-class="bg-indigo-700">
                Dashboard
              </RouterLink>
              <RouterLink to="/recommendation" class="px-3 py-2 rounded-md text-sm font-medium hover:bg-indigo-500 transition-colors" active-class="bg-indigo-700">
                Insights
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <RouterView v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </RouterView>
    </main>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
