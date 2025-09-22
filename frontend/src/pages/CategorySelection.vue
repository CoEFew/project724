<!-- src/pages/CategorySelection.vue -->
<template>
  <div class="min-h-screen px-4 py-10 bg-gradient-to-b from-slate-900 to-slate-950 text-slate-100">
    <div class="max-w-2xl mx-auto">
      <!-- Header with back button -->
      <header class="mb-8 flex items-center justify-between">
        <h1 class="text-2xl md:text-3xl font-extrabold tracking-wide">
          DOG • PUZZLE <span class="text-indigo-300">เลือกหมวดหมู่</span>
        </h1>
        <button
          @click="goBack"
          class="px-3 py-2 rounded-xl border border-white/10 bg-white/5 hover:bg-white/10 transition"
          aria-label="กลับหน้าหลัก"
        >
          ย้อนกลับ
        </button>
      </header>

      <!-- Game mode info -->
      <div class="mb-6 rounded-xl border border-white/10 bg-white/5 p-4">
        <h2 class="text-lg font-semibold text-indigo-100 mb-2">{{ gameModeTitle }}</h2>
        <p class="text-sm text-slate-300">{{ gameModeDescription }}</p>
      </div>

      <!-- Category Selection -->
      <section class="w-full max-w-2xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-6 space-y-5">
        <h2 class="text-xl md:text-2xl font-extrabold text-indigo-100 tracking-wide text-center">เลือกหมวดหมู่</h2>
        <div class="grid grid-cols-2 gap-3">
          <button @click="selectCategory('สัตว์')" 
            :class="['p-4 rounded-xl border transition-all', selectedCategory === 'สัตว์' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
            <div class="text-center">
              <div class="text-2xl mb-2">🐕</div>
              <div class="font-semibold">สัตว์</div>
              <div class="text-xs opacity-80">สัตว์ต่างๆ ในโลก</div>
            </div>
          </button>
          <button @click="selectCategory('เครื่องใช้ไฟฟ้า')" 
            :class="['p-4 rounded-xl border transition-all', selectedCategory === 'เครื่องใช้ไฟฟ้า' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
            <div class="text-center">
              <div class="text-2xl mb-2">⚡</div>
              <div class="font-semibold">เครื่องใช้ไฟฟ้า</div>
              <div class="text-xs opacity-80">อุปกรณ์ไฟฟ้าต่างๆ</div>
            </div>
          </button>
          <button @click="selectCategory('ผลไม้')" 
            :class="['p-4 rounded-xl border transition-all', selectedCategory === 'ผลไม้' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
            <div class="text-center">
              <div class="text-2xl mb-2">🍎</div>
              <div class="font-semibold">ผลไม้</div>
              <div class="text-xs opacity-80">ผลไม้ต่างๆ ในโลก</div>
            </div>
          </button>
          <button @click="selectCategory('อาชีพ')" 
            :class="['p-4 rounded-xl border transition-all', selectedCategory === 'อาชีพ' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
            <div class="text-center">
              <div class="text-2xl mb-2">👨‍💼</div>
              <div class="font-semibold">อาชีพ</div>
              <div class="text-xs opacity-80">อาชีพต่างๆ ในสังคม</div>
            </div>
          </button>
          <button @click="selectRandomCategory" 
            :class="['p-4 rounded-xl border transition-all col-span-2', selectedCategory === 'Random' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
            <div class="text-center">
              <div class="text-2xl mb-2">🎲</div>
              <div class="font-semibold">สุ่มหมวดหมู่</div>
              <div class="text-xs opacity-80">สุ่มจากสัตว์, เครื่องใช้ไฟฟ้า, ผลไม้, อาชีพ</div>
            </div>
          </button>
        </div>
        <button @click="startGame" :disabled="!selectedCategory"
          class="w-full px-4 py-3 rounded-xl font-semibold transition bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow">
          เริ่มเกม
        </button>
      </section>

      <!-- Toast notifications -->
      <transition name="toast" tag="div">
        <div v-if="toasts.length > 0" class="fixed top-4 right-4 z-50 space-y-2">
          <div v-for="toast in toasts" :key="toast.id" 
            :class="[
              'px-4 py-3 rounded-xl shadow-lg border max-w-sm',
              toast.type === 'success' ? 'bg-emerald-500/90 border-emerald-300/30 text-white' :
              toast.type === 'error' ? 'bg-rose-500/90 border-rose-300/30 text-white' :
              'bg-blue-500/90 border-blue-300/30 text-white'
            ]">
            <div class="font-semibold text-sm">{{ toast.title }}</div>
            <div class="text-xs opacity-90">{{ toast.message }}</div>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

/**
 * CategorySelection.vue
 * 
 * This component handles category selection for single-player game modes.
 * It provides a clean interface for users to choose from available categories
 * or select a random category before starting their game.
 * 
 * Features:
 * - Category selection with visual feedback
 * - Random category selection
 * - Game mode information display
 * - Toast notifications for user feedback
 * - Robust error handling
 */

const router = useRouter()
const route = useRoute()

// Reactive state
const selectedCategory = ref('')
const toasts = ref<{ id: string; title: string; message: string; type: 'info' | 'success' | 'error' }[]>([])

/**
 * Computed properties for game mode information
 * Determines the title and description based on the game mode from route query
 */
const gameModeTitle = computed(() => {
  const noTimer = route.query.noTimer === '1'
  return noTimer ? 'โหมดเดี่ยว (ไม่จับเวลา)' : 'โหมดเดี่ยว (มีเวลา)'
})

const gameModeDescription = computed(() => {
  const noTimer = route.query.noTimer === '1'
  return noTimer 
    ? 'เล่นสบาย ๆ ไม่ต้องกังวลเรื่องเวลา ข้ามข้อ/ขอใบ้ได้ตามต้องการ'
    : 'เล่นคนเดียว มีจับเวลาต่อข้อ ระบบใบ้คำและบันทึกสกอร์'
})

/**
 * Select a specific category
 * @param category - The category to select
 */
function selectCategory(category: string) {
  try {
    selectedCategory.value = category
    // Clear any previous random category selection
    ;(selectedCategory as any).actualCategory = null
    console.log('Category selected:', category)
  } catch (error) {
    console.error('Error selecting category:', error)
    showToast('เกิดข้อผิดพลาด', 'ไม่สามารถเลือกหมวดหมู่ได้', 'error')
  }
}

/**
 * Select a random category from the available categories
 * This function randomly picks one of: Animals, Electronics, Fruits, Jobs
 */
function selectRandomCategory() {
  try {
    const categories = ['สัตว์', 'เครื่องใช้ไฟฟ้า', 'ผลไม้', 'อาชีพ']
    const randomIndex = Math.floor(Math.random() * categories.length)
    selectedCategory.value = 'Random'
    // Store the actual selected category for API calls
    ;(selectedCategory as any).actualCategory = categories[randomIndex]
    showToast('สุ่มหมวดหมู่แล้ว', `เลือกหมวดหมู่: ${categories[randomIndex]}`, 'success')
    console.log('Random category selected:', categories[randomIndex])
  } catch (error) {
    console.error('Error selecting random category:', error)
    showToast('เกิดข้อผิดพลาด', 'ไม่สามารถสุ่มหมวดหมู่ได้', 'error')
  }
}

/**
 * Start the game with the selected category
 * Navigates to the appropriate page based on game mode
 */
async function startGame() {
  try {
    if (!selectedCategory.value) {
      showToast('กรุณาเลือกหมวดหมู่', 'ต้องเลือกหมวดหมู่ก่อนเริ่มเกม', 'error')
      return
    }

    // Prepare the actual category for the game
    const actualCategory = (selectedCategory as any).actualCategory || selectedCategory.value
    
    // Check if multiple-choice mode is requested
    const isMultipleChoiceMode = route.query.mode === 'multiple-choice'
    
    // Navigate to the appropriate page with category and game mode
    const query: any = {}
    if (route.query.noTimer === '1') {
      query.noTimer = '1'
    }
    query.category = actualCategory

    if (isMultipleChoiceMode) {
      // Navigate to DogQuestion for multiple-choice mode
      await router.push({
        name: 'DogQuestion',
        query
      })
      console.log('Starting multiple-choice game with category:', actualCategory)
    } else {
      // Navigate to DocumentsPage for typing mode
      await router.push({
        name: 'DocumentsPage',
        query
      })
      console.log('Starting typing game with category:', actualCategory)
    }
  } catch (error) {
    console.error('Error starting game:', error)
    showToast('เกิดข้อผิดพลาด', 'ไม่สามารถเริ่มเกมได้', 'error')
  }
}

/**
 * Navigate back to the previous page
 */
function goBack() {
  try {
    router.back()
  } catch (error) {
    console.error('Error navigating back:', error)
    // Fallback to home page if back navigation fails
    router.push({ name: 'DogAll' })
  }
}

/**
 * Show a toast notification
 * @param title - The title of the toast
 * @param message - The message of the toast
 * @param type - The type of the toast (info, success, error)
 */
function showToast(title: string, message: string, type: 'info' | 'success' | 'error' = 'info') {
  try {
    const id = Math.random().toString(36).slice(2)
    toasts.value.push({ id, title, message, type })
    
    // Auto-remove toast after 3.5 seconds
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, 3500)
  } catch (error) {
    console.error('Error showing toast:', error)
  }
}

/**
 * Component mounted lifecycle hook
 * Initialize the component and show welcome message
 */
onMounted(() => {
  try {
    console.log('CategorySelection component mounted')
    showToast('ยินดีต้อนรับ', 'เลือกหมวดหมู่ที่คุณต้องการเล่น', 'info')
  } catch (error) {
    console.error('Error in component mount:', error)
  }
})
</script>

<style scoped>
/* Toast animation styles */
.toast-enter-active, .toast-leave-active { 
  transition: opacity .25s ease, transform .25s ease; 
}
.toast-enter-from, .toast-leave-to { 
  opacity: 0; 
  transform: translateY(6px); 
}

/* Button hover effects */
button:not(:disabled):hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* Disabled button styles */
button:disabled {
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* Category button selection animation */
button[class*="border-indigo-400"] {
  animation: selectPulse 0.3s ease-out;
}

@keyframes selectPulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.02); }
  100% { transform: scale(1); }
}
</style>
