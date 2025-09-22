<!-- src/components/OptimizedImage.vue -->
<template>
  <div class="optimized-image-container" :class="containerClass">
    <!-- Placeholder while loading -->
    <div 
      v-if="showPlaceholder && !imageLoaded"
      class="absolute inset-0 bg-slate-800 animate-pulse rounded"
      :style="{ backgroundImage: placeholder ? `url(${placeholder})` : undefined }"
    />
    
    <!-- Loading spinner -->
    <div 
      v-if="showSpinner && !imageLoaded && !imageError"
      class="absolute inset-0 flex items-center justify-center"
    >
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-400"></div>
    </div>
    
    <!-- Main image -->
    <img
      ref="imageRef"
      :src="optimizedSrc"
      :alt="alt"
      :class="imageClass"
      :loading="loading"
      :sizes="sizes"
      :srcset="srcset"
      @load="onImageLoad"
      @error="onImageError"
      v-show="imageLoaded"
    />
    
    <!-- Fallback image on error -->
    <img
      v-if="imageError && fallback"
      :src="fallback"
      :alt="alt"
      :class="imageClass"
      @load="onFallbackLoad"
      @error="onFallbackError"
    />
    
    <!-- Error state -->
    <div 
      v-if="imageError && !fallback"
      class="absolute inset-0 flex items-center justify-center bg-slate-800 text-slate-400"
    >
      <div class="text-center">
        <div class="text-2xl mb-2">🖼️</div>
        <div class="text-sm">ไม่สามารถโหลดรูปภาพได้</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { 
  getOptimizedSrc, 
  generateSrcSet, 
  generateSizes, 
  createPlaceholder,
  type ImageConfig 
} from '../utils/imageOptimization'

/**
 * OptimizedImage.vue
 * 
 * A high-performance image component that provides:
 * - Automatic format optimization (WebP support detection)
 * - Responsive image loading with srcset
 * - Lazy loading with intersection observer
 * - Placeholder and loading states
 * - Error handling with fallback images
 * - Caching and performance optimization
 * 
 * Usage:
 * <OptimizedImage 
 *   :src="imageSrc" 
 *   alt="Description" 
 *   :sizes="[320, 640, 1024]"
 *   loading="lazy"
 *   class="w-full h-64"
 * />
 */

// Props
interface Props {
  src: string
  alt: string
  containerClass?: string
  imageClass?: string
  sizes?: number[]
  loading?: 'lazy' | 'eager'
  showPlaceholder?: boolean
  showSpinner?: boolean
  quality?: number
  format?: 'webp' | 'png' | 'jpg' | 'auto'
  lazy?: boolean
  placeholder?: string
  fallback?: string
}

const props = withDefaults(defineProps<Props>(), {
  containerClass: '',
  imageClass: '',
  sizes: () => [320, 640, 768, 1024, 1280],
  loading: 'lazy',
  showPlaceholder: true,
  showSpinner: true,
  quality: 75,
  format: 'auto',
  lazy: true
})

// Emits
const emit = defineEmits<{
  load: [event: Event]
  error: [event: Event]
}>()

// Refs
const imageRef = ref<HTMLImageElement | null>(null)
const imageLoaded = ref(false)
const imageError = ref(false)
const optimizedSrc = ref('')
const srcset = ref('')
const sizes = ref('')
const placeholder = ref('')

// Computed
const fallback = computed(() => props.fallback || '')

/**
 * Initialize image optimization
 */
async function initializeImage() {
  try {
    // Get optimized source
    optimizedSrc.value = await getOptimizedSrc(props.src, props.format as any)
    
    // Generate responsive srcset
    if (props.sizes && props.sizes.length > 0) {
      srcset.value = generateSrcSet(optimizedSrc.value, props.sizes, props.format as any)
      sizes.value = generateSizes()
    }
    
    // Create placeholder
    if (props.showPlaceholder) {
      placeholder.value = createPlaceholder(100, 100, '#1e293b')
    }
  } catch (error) {
    console.warn('Error initializing image:', error)
    optimizedSrc.value = props.src
  }
}

/**
 * Handle successful image load
 */
function onImageLoad(event: Event) {
  imageLoaded.value = true
  imageError.value = false
  emit('load', event)
}

/**
 * Handle image load error
 */
function onImageError(event: Event) {
  console.warn('Image load error:', props.src)
  imageError.value = true
  emit('error', event)
}

/**
 * Handle fallback image load
 */
function onFallbackLoad(event: Event) {
  imageLoaded.value = true
  imageError.value = false
  emit('load', event)
}

/**
 * Handle fallback image error
 */
function onFallbackError(event: Event) {
  console.warn('Fallback image also failed:', props.fallback)
  emit('error', event)
}

/**
 * Setup lazy loading with intersection observer
 */
function setupLazyLoading() {
  if (!props.lazy || props.loading === 'eager') return
  
  try {
    if (!('IntersectionObserver' in window)) {
      // Fallback for older browsers
      if (imageRef.value) {
        imageRef.value.src = optimizedSrc.value
      }
      return
    }

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            const img = entry.target as HTMLImageElement
            if (img.dataset.src) {
              img.src = img.dataset.src
              img.removeAttribute('data-src')
              observer.unobserve(img)
            }
          }
        })
      },
      {
        rootMargin: '50px 0px',
        threshold: 0.01
      }
    )

    if (imageRef.value) {
      imageRef.value.dataset.src = optimizedSrc.value
      observer.observe(imageRef.value)
    }
  } catch (error) {
    console.warn('Error setting up lazy loading:', error)
  }
}

// Lifecycle
onMounted(async () => {
  await initializeImage()
  setupLazyLoading()
})

// Watch for src changes
watch(() => props.src, async () => {
  imageLoaded.value = false
  imageError.value = false
  await initializeImage()
  setupLazyLoading()
})
</script>

<style scoped>
.optimized-image-container {
  position: relative;
  overflow: hidden;
}

.optimized-image-container img {
  transition: opacity 0.3s ease;
}

.optimized-image-container img[data-src] {
  opacity: 0;
}

.optimized-image-container img:not([data-src]) {
  opacity: 1;
}

/* Loading spinner animation */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Placeholder animation */
.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
