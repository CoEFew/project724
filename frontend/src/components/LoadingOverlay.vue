<!-- src/components/LoadingOverlay.vue -->
<template>
  <div v-if="loading"
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex flex-col items-center justify-center z-[90]">
    <div class="flex flex-col items-center">
      <!-- Cat loading animation -->
      <div class="relative h-24 w-24 mb-2">
        <img 
          :src="catwalkImages[catwalkIndex]" 
          alt="loading cat" 
          class="h-24 w-24 animate-bounce"
          loading="eager"
          @load="onImageLoad"
          @error="onImageError"
        />
        <!-- Fallback loading spinner if image fails -->
        <div v-if="imageError" class="h-24 w-24 flex items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-400"></div>
        </div>
      </div>
      
      <!-- Loading text -->
      <span class="text-base md:text-lg text-indigo-100 font-semibold">กำลังโหลด...</span>
      
      <!-- Network status messages -->
      <span v-if="net?.hasPending" class="mt-1 text-xs text-indigo-100/70">
        กำลังเชื่อมต่อเซิร์ฟเวอร์…
      </span>
      <span v-if="net?.isStalled" class="mt-1 text-xs text-amber-200/80">
        เซิร์ฟเวอร์กำลังเริ่มทำงาน
      </span>
      
      <!-- Custom loading message -->
      <span v-if="message" class="mt-1 text-xs text-slate-300/80">
        {{ message }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useNetworkStore } from '../store/useNetworkStore'

/**
 * LoadingOverlay.vue
 * 
 * A reusable loading component that provides consistent loading experience
 * across all pages with the animated cat visual from DocumentsPage.vue
 * 
 * Features:
 * - Animated cat loading visual
 * - Network status indicators
 * - Custom loading messages
 * - Fallback spinner for image errors
 * - Optimized image loading
 * - Robust error handling
 */

// Props
interface Props {
  loading: boolean
  message?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  message: ''
})

// Composables
const net = useNetworkStore()

// State
const catwalkIndex = ref(0)
const imageError = ref(false)
let catwalkInterval: number | undefined

// Optimized catwalk images - using smaller, compressed versions
const catwalkImages = [
  '/src/assets/images/catwalk.png', // Will be created
  '/src/assets/images/catwalk2.png' // Will be created
]

/**
 * Handle successful image load
 */
function onImageLoad() {
  imageError.value = false
}

/**
 * Handle image load error - fallback to spinner
 */
function onImageError() {
  console.warn('Failed to load catwalk image, using fallback spinner')
  imageError.value = true
}

/**
 * Start the catwalk animation interval
 */
function startAnimation() {
  if (catwalkInterval) return
  
  catwalkInterval = window.setInterval(() => {
    catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length
  }, 200)
}

/**
 * Stop the catwalk animation interval
 */
function stopAnimation() {
  if (catwalkInterval) {
    clearInterval(catwalkInterval)
    catwalkInterval = undefined
  }
}

// Lifecycle
onMounted(() => {
  if (props.loading) {
    startAnimation()
  }
})

onBeforeUnmount(() => {
  stopAnimation()
})

// Watch for loading state changes
import { watch } from 'vue'
watch(() => props.loading, (newLoading) => {
  if (newLoading) {
    startAnimation()
  } else {
    stopAnimation()
  }
})
</script>

<style scoped>
/* Ensure smooth animations */
.animate-bounce {
  animation: bounce 1s infinite;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes bounce {
  0%, 20%, 53%, 80%, 100% {
    transform: translate3d(0, 0, 0);
  }
  40%, 43% {
    transform: translate3d(0, -8px, 0);
  }
  70% {
    transform: translate3d(0, -4px, 0);
  }
  90% {
    transform: translate3d(0, -2px, 0);
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
