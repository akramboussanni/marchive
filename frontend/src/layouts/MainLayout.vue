<template>
  <div class="main-layout">
    <AppHeader @open-search="searchModalOpen = true" />
    <router-view :key="routerViewKey" />
    <SearchModal 
      :is-open="searchModalOpen" 
      @close="searchModalOpen = false"
      @book-added="handleBookAdded"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import SearchModal from '@/components/SearchModal.vue'

const router = useRouter()
const searchModalOpen = ref(false)
const routerViewKey = ref(0)

const handleBookAdded = () => {
  searchModalOpen.value = false
  // Navigate to home if not already there
  if (router.currentRoute.value.path !== '/') {
    router.push('/')
  } else {
    // Force refresh by updating key
    routerViewKey.value++
  }
}
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0f1e 0%, #1a1f3a 100%);
}
</style>
