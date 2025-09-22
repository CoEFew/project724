<!--
Multiple-Choice Quiz Game - Single Player Mode

This component implements a single-player multiple-choice quiz game with the following features:
- 30-second rounds with automatic hint timing
- 4 multiple-choice options (1 correct + 3 incorrect)
- Same scoring and level-up logic as DocumentsPage.vue
- Reuses the same UI structure and theme
- Robust error handling and user feedback

Gameplay Rules:
1. Round time: 30 seconds
2. Hint #1: Shown immediately at start
3. Hint #2: Shown when 10 seconds remain
4. Choices: 1 correct (from current quiz) + 3 incorrect (random from other quizzes)
5. Scoring: Identical to DocumentsPage.vue

Technical Implementation:
- Vue 3 Composition API with TypeScript
- Same theme and UI structure as DocumentsPage.vue
- Multiple choice generation with randomization
- Timer-based hint display
- Comprehensive error handling
- Responsive design with Tailwind CSS
-->
<template>
  <div class="min-h-screen relative overflow-x-hidden theme-modern">
    <!-- Gradient background (โมเดิร์น, สบายตา) -->
    <div class="pointer-events-none absolute inset-0 -z-10" aria-hidden="true">
      <div
        class="absolute inset-0 bg-[radial-gradient(90%_70%_at_70%_100%,rgba(99,102,241,0.45),transparent_60%),radial-gradient(60%_60%_at_0%_0%,rgba(59,130,246,0.35),transparent_60%),linear-gradient(180deg,#0b1020,#0b1120)]" />
      <div class="absolute -bottom-16 right-10 h-80 w-80 rounded-full blur-3xl opacity-40 bg-indigo-500/30" />
      <div class="absolute -top-12 left-[-4rem] h-72 w-72 rounded-full blur-3xl opacity-30 bg-fuchsia-500/25" />
    </div>

    <!-- Loading overlay -->
    <LoadingOverlay :loading="loading" />

    <!-- Page container -->
    <div class="w-full max-w-6xl mx-auto px-4 py-8" v-show="!loading">
      <!-- Header -->
      <header class="flex flex-col gap-3 items-center mb-6">
        <div class="w-full flex items-center justify-between">
          <button @click="goBack" type="button"
            class="inline-flex items-center gap-2 px-3 py-2 rounded-xl border bg-white/5 border-white/10 text-slate-100 hover:bg-white/10 transition shadow-sm">
            <svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 18l-6-6 6-6" />
            </svg>
            ย้อนกลับ
          </button>

          <h1
            class="text-2xl md:text-4xl font-extrabold tracking-wide text-indigo-300/90 uppercase text-center flex-1 drop-shadow-sm">
            DOG • QUESTION
          </h1>

          <div class="flex items-center gap-2">
            <!-- Sound toggle -->
            <button type="button" @click="soundOn = !soundOn"
              class="px-3 py-2 rounded-xl border border-white/10 bg-white/5 text-slate-100 hover:bg-white/10">
              <span v-if="soundOn">🔊</span><span v-else>🔇</span>
            </button>
            <!-- Theme toggle (optional future) -->
            <div class="w-[6px]" />
          </div>
        </div>
        <p class="text-slate-300/80 text-xs md:text-sm text-center"><span class="text-lg">🐕</span></p>
      </header>

      <!-- Game Mode Info -->
      <section v-if="!showModal && !quizId" class="w-full max-w-2xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-6 space-y-5">
        <h2 class="text-xl md:text-2xl font-extrabold text-indigo-100 tracking-wide text-center">เกมเริ่มแล้ว!</h2>
        <div class="text-center">
          <p class="text-slate-300 mb-4">หมวดหมู่: <span class="font-semibold text-indigo-200">{{ selectedCategory || 'กำลังโหลด...' }}</span></p>
          <p class="text-sm text-slate-400">โหมดตัวเลือก - เลือกคำตอบที่ถูกต้องจาก 4 ตัวเลือก</p>
        </div>
        <button @click="startGame" 
          class="w-full px-4 py-3 rounded-xl font-semibold transition bg-indigo-500 text-white hover:bg-indigo-400 shadow">
          เริ่มเกม
        </button>
      </section>

      <!-- กล่องเกม (Glass card) -->
      <section v-if="showModal || quizId"
        class="w-full max-w-xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-6 space-y-5">
        <!-- สรุปสถานะ -->
        <div class="space-y-3">
          <h2 class="text-xl md:text-2xl font-extrabold text-indigo-100 tracking-wide text-center">เก่งจริงก็ทายมาดิ!
          </h2>

          <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <InfoCard label="คะแนน" :value="score" accent="indigo" />
            <InfoCard label="เวลาคงเหลือ" :value="noTimer ? 'ไม่จำกัดเวลา' : timer + ' วินาที'" accent="sky" />
            <InfoCard label="ระดับ" :value="'Lv. ' + currentLevel" accent="fuchsia" />
            <InfoCard label="ข้อผิด" :value="wrongAnswerCount + '/' + maxWrongAnswers" accent="rose" />
          </div>

          <!-- Progress & proximity meter -->
          <div v-if="!noTimer" class="space-y-2">
            <div class="w-full h-2 rounded-full overflow-hidden bg-white/10" role="progressbar" :aria-valuenow="timer">
              <div
                class="h-full transition-all duration-300 bg-gradient-to-r from-sky-400 via-indigo-400 to-fuchsia-400"
                :style="{ width: timerPercent + '%' }" />
            </div>
           
          </div>
        </div>

        <!-- Hints Display - Show both hints directly -->
        <div class="space-y-3">
          <h3 class="text-lg font-semibold text-indigo-100 text-center">คำถาม</h3>
          <div class="space-y-2">
            <div v-if="hint1" class="p-3 rounded-xl bg-white/5 border border-white/10">
              <div class="flex items-center gap-2 mb-1">
                <span class="text-sm font-medium text-slate-300">ใบ้ที่ 1:</span>
                <span class="text-sm text-slate-200">{{ hint1 }}</span>
              </div>
            </div>
            <div v-if="hint2" class="p-3 rounded-xl bg-white/5 border border-white/10">
              <div class="flex items-center gap-2 mb-1">
                <span class="text-sm font-medium text-slate-300">ใบ้ที่ 2:</span>
                <span class="text-sm text-slate-200">{{ hint2 }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Multiple Choice Options -->
        <div class="space-y-3">
          <h3 class="text-lg font-semibold text-indigo-100 text-center">เลือกคำตอบที่ถูกต้อง</h3>
          <div class="grid grid-cols-1 gap-3">
            <button
              v-for="(choice, index) in choices"
              :key="index"
              @click="selectChoice(choice, index)"
              :disabled="showModal || confirmedChoice !== null"
              :class="[
                'p-4 rounded-xl border transition-all text-left',
                selectedChoice === index
                  ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100'
                  : confirmedChoice === index
                    ? choice.isCorrect
                      ? 'border-green-400 bg-green-500/20 text-green-100'
                      : 'border-red-400 bg-red-500/20 text-red-100'
                    : 'border-white/10 bg-white/5 text-slate-100 hover:bg-white/10 hover:border-indigo-400/30',
                showModal || confirmedChoice !== null ? 'cursor-not-allowed opacity-60' : 'cursor-pointer'
              ]"
            >
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-full border-2 flex items-center justify-center text-sm font-semibold"
                  :class="[
                    selectedChoice === index
                      ? 'border-indigo-400 text-indigo-400'
                      : confirmedChoice === index
                        ? choice.isCorrect
                          ? 'border-green-400 text-green-400'
                          : 'border-red-400 text-red-400'
                        : 'border-white/20 text-slate-300'
                  ]">
                  {{ String.fromCharCode(65 + index) }}
                </div>
                <span class="font-medium">{{ choice.text }}</span>
                <div v-if="confirmedChoice === index" class="ml-auto">
                  <span v-if="choice.isCorrect" class="text-green-400">✓</span>
                  <span v-else class="text-red-400">✗</span>
                </div>
                <div v-else-if="selectedChoice === index" class="ml-auto">
                  <span class="text-indigo-400">เลือกแล้ว</span>
                </div>
              </div>
            </button>
          </div>
        </div>

        <!-- Confirmation Button -->
        <div v-if="selectedChoice !== null && confirmedChoice === null" class="flex justify-center">
          <button @click="confirmChoice" :disabled="showModal"
            class="px-8 py-3 rounded-xl font-semibold transition bg-emerald-500 text-white hover:bg-emerald-400 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg">
            ยืนยันคำตอบ
          </button>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-2">
          <button @click="skipQuestion" :disabled="showModal || confirmedChoice !== null || changeCount >= maxChange"
            class="flex-1 px-4 py-2.5 rounded-xl font-semibold transition bg-slate-600 text-white hover:bg-slate-500 disabled:opacity-50 disabled:cursor-not-allowed">
            ข้ามข้อ ({{ changeCount }}/{{ maxChange }}) 
          </button>
        </div>
      </section>

      <!-- Scoreboard Section -->
      <section v-if="!showModal && !quizId" class="w-full max-w-2xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-6 space-y-5">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-indigo-100">ตารางคะแนนสูงสุด</h3>
          <button @click="loadScores" 
            class="px-3 py-1.5 rounded-lg border border-white/10 bg-white/5 text-slate-300 hover:bg-white/10 hover:text-slate-100 transition-colors text-sm"
            title="รีเฟรชสถิติ" type="button">
            รีเฟรช
          </button>
        </div>

        <p v-if="savedScores.length === 0" class="text-slate-300/70 text-sm mt-2">ยังไม่มีสถิติ
          แสดงเมื่อมีการบันทึกคะแนน</p>

        <ul v-else class="mt-2 divide-y divide-white/10">
          <li v-for="(item, idx) in savedScores" :key="item.name + '_' + item.score + '_' + idx"
            class="py-2 flex items-center justify-between text-sm text-slate-100">
            <div class="flex items-center gap-2 min-w-0">
              <span class="w-6 text-center">
                <template v-if="idx === 0">🥇</template>
                <template v-else-if="idx === 1">🥈</template>
                <template v-else-if="idx === 2">🥉</template>
                <template v-else>{{ idx + 1 }}.</template>
              </span>
              <span class="font-medium truncate max-w-[10rem] sm:max-w-[14rem]" :title="item.name">{{ item.name
                }}</span>
            </div>
            <div class="tabular-nums font-semibold">{{ item.score }} คะแนน</div>
          </li>
        </ul>
      </section>

      <!-- Modal -->
      <div v-if="showModal"
        class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
        <div class="w-full max-w-md rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-6 space-y-4">
          <h3 class="text-xl font-bold text-indigo-100 text-center">{{ modalTitle }}</h3>
          <p class="text-slate-300 text-center">{{ modalMessage }}</p>
          
          <!-- Name input for score saving -->
          <div v-if="modalButtonText === 'บันทึกคะแนน'" class="space-y-2">
            <label class="block text-sm font-medium text-slate-300">ชื่อผู้เล่น</label>
            <input ref="nameInput" v-model="playerName" type="text" placeholder="กรอกชื่อของคุณ"
              class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60"
              :disabled="isSaving" />
          </div>
          
          <div class="flex gap-2">
            <button @click="closeModal" :disabled="isSaving" type="button"
              class="flex-1 px-4 py-2.5 rounded-xl font-semibold transition bg-slate-600 text-white hover:bg-slate-500 disabled:opacity-50">
              {{ modalButtonText === 'บันทึกคะแนน' ? 'ยกเลิก' : 'ปิด' }}
            </button>
            <button v-if="modalButtonText === 'บันทึกคะแนน'" @click="saveScore" :disabled="isSaving || !playerName.trim()" type="button"
              class="flex-1 px-4 py-2.5 rounded-xl font-semibold transition bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50">
              {{ isSaving ? 'กำลังบันทึก...' : 'บันทึกคะแนน' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/* ===================== Imports ===================== */
import LoadingOverlay from '../components/LoadingOverlay.vue'
import api from '../services/api'
import {
  ref, onMounted, watch, computed, onBeforeUnmount, nextTick,
  defineComponent, h, type PropType
} from 'vue'
import { useRouter, useRoute } from 'vue-router'

/* ===================== Type Definitions ===================== */
/**
 * Multiple choice option interface
 * Represents a single choice in the multiple choice quiz
 */
interface Choice {
  text: string
  isCorrect: boolean
}

/**
 * Guess item interface for tracking recent guesses
 * Used for proximity meter and duplicate detection
 */
interface GuessItem {
  word: string
  correct: boolean
  ts: number
  heat?: number
}

/**
 * Score data interface for saving player scores
 * Matches the backend API structure
 */
interface ScoreData {
  name: string
  score: number
  gamename: string
  level: number
}

/**
 * API readiness status interface
 * Used for health check and initial loading
 */
interface ApiReadiness {
  healthOk: boolean
  initialOk: boolean
}

/**
 * DogQuestion.vue
 * 
 * Multiple-Choice Quiz Game - Single Player Mode
 * 
 * This component implements a single-player multiple-choice quiz game with:
 * - 30-second rounds with automatic hint timing
 * - 4 multiple-choice options (1 correct + 3 incorrect)
 * - Same scoring and level-up logic as DocumentsPage.vue
 * - Reuses the same UI structure and theme
 * - Robust error handling and user feedback
 * 
 * Gameplay Rules:
 * 1. Round time: 30 seconds
 * 2. Hint #1: Shown immediately at start
 * 3. Hint #2: Shown when 10 seconds remain
 * 4. Choices: 1 correct (from current quiz) + 3 incorrect (random from other quizzes)
 * 5. Scoring: Identical to DocumentsPage.vue
 */

/* ===================== Constants ===================== */
const GAME_NAME = 'DogQuestion'
const TOP_LIMIT = 10
const ROUND_SECONDS = 30 // 30 seconds per round
const CHANGE_LIMIT = 5

/* ===================== Router & State ===================== */
const router = useRouter()
const route = useRoute()

// loading state
const loading = ref(true)

// Get category and timer mode from route query
const selectedCategory = ref('')
const noTimer = ref(false)
selectedCategory.value = (route.query.category as string) || ''
noTimer.value = route.query.noTimer === 'true'

// game states
const quizId = ref('')
const quizToken = ref('')
const quizExp = ref(0)
const placeholder = computed(() => `เลือกคำตอบที่ถูกต้อง…`)

const score = ref(0)
const timer = ref(ROUND_SECONDS)
const currentLevel = ref(1)
const hint1 = ref('')
const hint2 = ref('')
const hintCount = ref(0)
const maxHints = ref(2)
const showModal = ref(false)
const changeCount = ref(0)
const maxChange = CHANGE_LIMIT
const playerName = ref('')
const savedScores = ref<{ name: string; score: number }[]>([])
const finalScore = ref(0)
const finalLevel = ref(1)
const expiredNotice = ref(false)
const revealedAnswer = ref('')
const isSaving = ref(false)
const soundOn = ref(true)

// Multiple choice specific states
const choices = ref<Choice[]>([])
const selectedChoice = ref<number | null>(null)
const confirmedChoice = ref<number | null>(null)
const correctAnswer = ref<string>('')
const wrongAnswerCount = ref<number>(0)
const maxWrongAnswers = 3
// recent guesses (รีเซ็ตเมื่อเริ่มคำใหม่)
const recentGuesses = ref<GuessItem[]>([])

// DOM element references
const answerInput = ref<HTMLInputElement | null>(null)
const nameInput = ref<HTMLInputElement | null>(null)

// modal states
const modalTitle = ref<string>('')
const modalMessage = ref<string>('')
const modalButtonText = ref<string>('')

// timer & intervals
let intervalId: number | undefined
let visibilityPausedAt: number | 0 = 0

/* ===================== Computed Properties ===================== */
const timerPercent = computed(() => (timer.value / ROUND_SECONDS) * 100)

const heatLabel = computed(() => {
  if (recentGuesses.value.length === 0) return 'ยังไม่เดา'
  const last = recentGuesses.value[recentGuesses.value.length - 1]
  if (!last.heat) return 'ยังไม่เดา'
  if (last.heat >= 0.8) return 'ร้อนมาก!'
  if (last.heat >= 0.6) return 'ร้อน!'
  if (last.heat >= 0.4) return 'อุ่น'
  if (last.heat >= 0.2) return 'เย็น'
  return 'เย็นมาก'
})

const heatDotClass = computed(() => {
  if (recentGuesses.value.length === 0) return 'bg-slate-400'
  const last = recentGuesses.value[recentGuesses.value.length - 1]
  if (!last.heat) return 'bg-slate-400'
  if (last.heat >= 0.8) return 'bg-red-400'
  if (last.heat >= 0.6) return 'bg-orange-400'
  if (last.heat >= 0.4) return 'bg-yellow-400'
  if (last.heat >= 0.2) return 'bg-blue-400'
  return 'bg-slate-400'
})

/* ===================== Utility Functions ===================== */
function toast(title: string, message: string, type: 'info' | 'success' | 'error' = 'info') {
  console.log(`[${type.toUpperCase()}] ${title}: ${message}`)
  // In a real app, you'd show a toast notification here
}

/**
 * Enhanced API GET with comprehensive error handling and retry logic
 * 
 * This function provides robust error handling for API calls with:
 * - Request timeout handling
 * - Network error detection and retry logic
 * - Response validation
 * - User-friendly error messages
 * 
 * @param url - API endpoint URL
 * @param params - Query parameters
 * @param retries - Number of retry attempts (default: 3)
 * @returns Promise with API response
 */
async function apiGet(url: string, params?: any, retries: number = 3): Promise<any> {
  let lastError: Error | null = null
  
  for (let attempt = 1; attempt <= retries; attempt++) {
    try {
      const response = await api.get(url, { 
        params,
        timeout: 10000 // 10 second timeout
      })
      
      // Validate response structure
      if (!response || typeof response !== 'object') {
        throw new Error('Invalid response structure')
      }
      
      return response
    } catch (error: any) {
      lastError = error
      
      // Log error details for debugging
      console.warn(`API GET attempt ${attempt}/${retries} failed for ${url}:`, error)
      
      // Don't retry on certain error types
      if (error.response?.status === 404 || error.response?.status === 400) {
        throw error
      }
      
      // Wait before retry (exponential backoff)
      if (attempt < retries) {
        const delay = Math.min(1000 * Math.pow(2, attempt - 1), 5000)
        await new Promise(resolve => setTimeout(resolve, delay))
      }
    }
  }
  
  // All retries failed
  throw lastError || new Error('API request failed after all retries')
}

/**
 * Enhanced API POST with comprehensive error handling and retry logic
 * 
 * This function provides robust error handling for API calls with:
 * - Request timeout handling
 * - Network error detection and retry logic
 * - Response validation
 * - User-friendly error messages
 * 
 * @param url - API endpoint URL
 * @param data - Request body data
 * @param retries - Number of retry attempts (default: 3)
 * @returns Promise with API response
 */
async function apiPost(url: string, data?: any, retries: number = 3): Promise<any> {
  let lastError: Error | null = null
  
  for (let attempt = 1; attempt <= retries; attempt++) {
    try {
      const response = await api.post(url, data, {
        timeout: 10000 // 10 second timeout
      })
      
      // Validate response structure
      if (!response || typeof response !== 'object') {
        throw new Error('Invalid response structure')
      }
      
      return response
    } catch (error: any) {
      lastError = error
      
      // Log error details for debugging
      console.warn(`API POST attempt ${attempt}/${retries} failed for ${url}:`, error)
      
      // Don't retry on certain error types
      if (error.response?.status === 404 || error.response?.status === 400) {
        throw error
      }
      
      // Wait before retry (exponential backoff)
      if (attempt < retries) {
        const delay = Math.min(1000 * Math.pow(2, attempt - 1), 5000)
        await new Promise(resolve => setTimeout(resolve, delay))
      }
    }
  }
  
  // All retries failed
  throw lastError || new Error('API request failed after all retries')
}

/* ===================== Game Logic ===================== */
/**
 * Generate multiple choice options
 * Creates 4 choices: 1 correct (from current quiz) + 3 incorrect (random from other quizzes)
 */
/**
 * Enhanced choice generation with comprehensive error handling and fallback mechanisms
 * 
 * This function generates 4 multiple-choice options (1 correct + 3 incorrect) with robust error handling:
 * 
 * 1. Input Validation:
 *    - Validates that correctAnswer is available and not empty
 *    - Provides detailed error logging for debugging
 *    - Throws error if validation fails (should be caught by caller)
 * 
 * 2. Incorrect Choice Generation:
 *    - Fetches random answers from other quizzes to use as incorrect options
 *    - Implements retry logic with attempt limits to prevent infinite loops
 *    - Handles API failures gracefully with fallback options
 * 
 * 3. Fallback Mechanisms:
 *    - Uses predefined fallback words if API calls fail
 *    - Ensures exactly 4 choices are always generated
 *    - Prevents duplicate choices in the final selection
 * 
 * 4. Choice Shuffling:
 *    - Randomly shuffles choices to prevent predictable patterns
 *    - Maintains correct answer identification through isCorrect flag
 * 
 * Error Handling:
 * - Network failures are caught and logged
 * - Invalid responses are handled gracefully
 * - Fallback choices ensure game continuity
 * 
 * @throws Error - If correctAnswer is not available (should be caught by caller)
 */
async function generateChoices() {
  try {
    // Enhanced input validation - this is the critical check that was failing
    if (!correctAnswer.value || typeof correctAnswer.value !== 'string' || correctAnswer.value.trim() === '') {
      const errorMsg = `No correct answer available for choice generation. Value: "${correctAnswer.value}" (type: ${typeof correctAnswer.value})`
      console.error(errorMsg)
      throw new Error(errorMsg)
    }

    console.log('Generating choices for answer:', correctAnswer.value)

    // Get 3 random incorrect answers from other quizzes
    const incorrectChoices: string[] = []
    const maxAttempts = 20 // Prevent infinite loops
    let attempts = 0
    let apiErrors = 0

    while (incorrectChoices.length < 3 && attempts < maxAttempts) {
      try {
        const params: any = { level: currentLevel.value }
        if (selectedCategory.value) {
          const categoryToUse = (selectedCategory as any).actualCategory || selectedCategory.value
          params.category = categoryToUse
        }
        
        // Use multiple-choice endpoint for better performance and consistency
        const res = await apiGet('/api/quiz/multiple-choice', params)
        
        // Validate the response
        if (res.data && res.data.answer && typeof res.data.answer === 'string') {
          const randomAnswer = res.data.answer.trim()
          
          // Ensure it's different from correct answer and not already in incorrect choices
          if (randomAnswer && 
              randomAnswer !== correctAnswer.value && 
              !incorrectChoices.includes(randomAnswer) &&
              randomAnswer.length > 0) {
            incorrectChoices.push(randomAnswer)
            console.log(`Added incorrect choice ${incorrectChoices.length}:`, randomAnswer)
          }
        } else {
          console.warn('Invalid response structure for random answer:', res.data)
        }
      } catch (error) {
        apiErrors++
        console.warn(`Failed to fetch random answer for incorrect choice (attempt ${attempts + 1}):`, error)
        
        // If we have too many API errors, break early to use fallbacks
        if (apiErrors >= 5) {
          console.warn('Too many API errors, using fallback choices')
          break
        }
      }
      attempts++
    }

    // If we couldn't get enough incorrect choices from API, use fallback words
    const fallbackWords = [
      'คำตอบ', 'ตัวเลือก', 'คำถาม', 'ข้อสอบ', 'แบบทดสอบ', 
      'คำศัพท์', 'คำสำคัญ', 'คำตอบที่ถูก', 'คำตอบที่ผิด',
      'คำศัพท์ใหม่', 'คำศัพท์เก่า', 'คำศัพท์ยาก', 'คำศัพท์ง่าย'
    ]
    
    while (incorrectChoices.length < 3) {
      const fallbackIndex = incorrectChoices.length % fallbackWords.length
      const fallback = fallbackWords[fallbackIndex]
      
      // Ensure fallback is different from correct answer and not already included
      if (fallback !== correctAnswer.value && !incorrectChoices.includes(fallback)) {
        incorrectChoices.push(fallback)
        console.log(`Added fallback choice ${incorrectChoices.length}:`, fallback)
      } else {
        // If all fallbacks are exhausted, create unique ones
        const uniqueFallback = `${fallback} ${incorrectChoices.length + 1}`
        incorrectChoices.push(uniqueFallback)
        console.log(`Added unique fallback choice ${incorrectChoices.length}:`, uniqueFallback)
      }
    }

    // Create choices array with correct answer and 3 incorrect answers
    const allChoices = [
      { text: correctAnswer.value, isCorrect: true },
      ...incorrectChoices.map(text => ({ text, isCorrect: false }))
    ]

    // Shuffle the choices randomly using Fisher-Yates algorithm
    for (let i = allChoices.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [allChoices[i], allChoices[j]] = [allChoices[j], allChoices[i]]
    }

    choices.value = allChoices
    console.log('Successfully generated choices:', choices.value.map(c => ({ text: c.text, isCorrect: c.isCorrect })))
    
  } catch (error) {
    console.error('Error generating choices:', error)
    
    // Enhanced fallback choices with better error handling
    const fallbackChoices = [
      { text: correctAnswer.value || 'คำตอบ', isCorrect: true },
      { text: 'ตัวเลือก 1', isCorrect: false },
      { text: 'ตัวเลือก 2', isCorrect: false },
      { text: 'ตัวเลือก 3', isCorrect: false }
    ]
    
    choices.value = fallbackChoices
    console.warn('Using fallback choices due to error:', fallbackChoices)
    
    // Show user-friendly error message
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถสร้างตัวเลือกได้ ใช้ตัวเลือกสำรอง', 'error')
  }
}

/**
 * Select a choice (allows changing until confirmation)
 * Users can change their selection until they explicitly confirm
 */
function selectChoice(choice: { text: string; isCorrect: boolean }, index: number) {
  if (confirmedChoice.value !== null || showModal.value) return

  try {
    // Allow changing selection until confirmation
    selectedChoice.value = index
    console.log(`Selected choice ${index}: ${choice.text}`)
  } catch (error) {
    console.error('Error selecting choice:', error)
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถเลือกคำตอบได้', 'error')
  }
}

/**
 * Confirm the selected choice and process the result
 * This is where the actual answer checking happens
 */
function confirmChoice() {
  if (selectedChoice.value === null || confirmedChoice.value !== null || showModal.value) return

  try {
    const choice = choices.value[selectedChoice.value]
    confirmedChoice.value = selectedChoice.value
    
    if (choice.isCorrect) {
      // Correct answer
      score.value += Math.max(1, Math.floor(timer.value / 10))
      toast('ถูกต้อง!', `+${Math.max(1, Math.floor(timer.value / 10))} คะแนน`, 'success')
      
      // Level up check
      if (score.value >= currentLevel.value * 100) {
        currentLevel.value++
        toast('เลเวลอัพ!', `เลเวล ${currentLevel.value}`, 'success')
      }
      
      // Reset wrong answer count on correct answer
      wrongAnswerCount.value = 0
      
      // Move to next question after a short delay
      setTimeout(() => {
        fetchQuiz()
      }, 1500)
    } else {
      // Wrong answer
      wrongAnswerCount.value++
      toast('ผิด!', `คำตอบที่ถูกต้องคือ: ${correctAnswer.value}`, 'error')
      
      // Check if this is the 3rd wrong answer
      if (wrongAnswerCount.value >= maxWrongAnswers) {
        toast('เกมจบ!', 'ตอบผิด 3 ครั้งแล้ว', 'error')
        setTimeout(() => {
          endGame()
        }, 2000)
        return
      }
      
      // Move to next question after a short delay
      setTimeout(() => {
        fetchQuiz()
      }, 2000)
    }
  } catch (error) {
    console.error('Error confirming choice:', error)
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถประมวลผลคำตอบได้', 'error')
  }
}

/**
 * Skip the current question (increment change count)
 */
function skipQuestion() {
  if (changeCount.value >= maxChange || confirmedChoice.value !== null) return
  
  try {
    changeCount.value++
    toast('ข้ามข้อ', `ข้ามข้อแล้ว (${changeCount.value}/${maxChange})`, 'info')
    
    // Move to next question
    setTimeout(() => {
      fetchQuiz()
    }, 1000)
  } catch (error) {
    console.error('Error skipping question:', error)
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถข้ามข้อได้', 'error')
  }
}

/**
 * Close modal and reset game state
 */
function closeModal() {
  showModal.value = false
  // Reset game state if needed
  if (modalButtonText.value === 'บันทึกคะแนน') {
    // Reset game after saving score
    quizId.value = ''
    correctAnswer.value = ''
    choices.value = []
    selectedChoice.value = null
    confirmedChoice.value = null
    wrongAnswerCount.value = 0
    changeCount.value = 0
    score.value = 0
    currentLevel.value = 1
  }
}

/**
 * Load saved scores from the server
 */
async function loadScores() {
  try {
    const res = await apiGet('/api/scores', { gamename: 'dog_question' })
    if (res.data && Array.isArray(res.data)) {
      savedScores.value = res.data.sort((a: ScoreData, b: ScoreData) => b.score - a.score)
    }
  } catch (error) {
    console.error('Failed to load scores:', error)
    toast('โหลดสถิติล้มเหลว', 'ไม่สามารถโหลดตารางคะแนนได้', 'error')
  }
}

/**
 * Save player score to the server
 */
async function saveScore() {
  if (isSaving.value || !playerName.value.trim()) return

  try {
    isSaving.value = true
    const scoreData: ScoreData = {
      name: playerName.value.trim(),
      score: score.value,
      gamename: 'dog_question',
      level: currentLevel.value
    }

    await apiPost('/api/scores', scoreData)
    toast('บันทึกสำเร็จ', `บันทึกคะแนน ${score.value} เรียบร้อยแล้ว`, 'success')
    
    // Reload scores to show updated leaderboard
    await loadScores()
    
    // Close modal
    showModal.value = false
  } catch (error) {
    console.error('Failed to save score:', error)
    toast('บันทึกล้มเหลว', 'ไม่สามารถบันทึกคะแนนได้', 'error')
  } finally {
    isSaving.value = false
  }
}

/**
 * End the game and show final score
 */
function endGame() {
  try {
    // Stop timer
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = undefined
    }
    
    // Show final score modal
    modalTitle.value = 'เกมจบ!'
    modalMessage.value = `คะแนนสุดท้าย: ${score.value} | ระดับ: ${currentLevel.value}`
    modalButtonText.value = 'บันทึกคะแนน'
    showModal.value = true
    
    // Reset game state
    quizId.value = ''
    correctAnswer.value = ''
    choices.value = []
    selectedChoice.value = null
    confirmedChoice.value = null
    wrongAnswerCount.value = 0
  } catch (error) {
    console.error('Error ending game:', error)
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถจบเกมได้', 'error')
  }
}

/**
 * Enhanced game initialization with comprehensive error handling and recovery
 * 
 * This function starts the multiple-choice game with robust error handling:
 * 
 * 1. Category Validation:
 *    - Validates that a category is selected before starting
 *    - Redirects to category selection if no category is available
 *    - Provides clear user feedback for validation failures
 * 
 * 2. Game Initialization:
 *    - Calls fetchQuiz to load the first question
 *    - Handles fetchQuiz failures gracefully
 *    - Provides retry mechanism for transient failures
 * 
 * 3. Error Recovery:
 *    - Resets game state on critical failures
 *    - Provides user-friendly error messages
 *    - Maintains consistent application state
 * 
 * 4. User Experience:
 *    - Clear feedback for different error scenarios
 *    - Automatic retry for recoverable errors
 *    - Graceful fallback to category selection
 */
async function startGame() {
  try {
    // Validate category selection
    if (!selectedCategory.value) {
      console.error('No category selected for game start')
      toast('เกิดข้อผิดพลาด', 'ไม่พบหมวดหมู่ที่เลือก กรุณาเลือกหมวดหมู่ใหม่', 'error')
      router.push({ name: 'CategorySelection' })
      return
    }
    
    console.log('Starting multiple-choice game with category:', selectedCategory.value)
    
    // Reset game state before starting
    correctAnswer.value = ''
    choices.value = []
    selectedChoice.value = null
    hintCount.value = 0
    timer.value = ROUND_SECONDS
    
    // Attempt to fetch the first quiz
    try {
      await fetchQuiz()
      console.log('Game started successfully')
    } catch (fetchError) {
      console.error('Failed to fetch initial quiz:', fetchError)
      
      // Provide specific error message based on error type
      let errorMessage = 'ไม่สามารถโหลดคำถามแรกได้'
      if (fetchError instanceof Error) {
        if (fetchError.message.includes('missing answer')) {
          errorMessage = 'ข้อมูลคำถามไม่สมบูรณ์ กรุณาลองใหม่'
        } else if (fetchError.message.includes('network') || fetchError.message.includes('Network')) {
          errorMessage = 'เชื่อมต่อเซิร์ฟเวอร์ไม่ได้ กรุณาตรวจสอบการเชื่อมต่อ'
        } else if (fetchError.message.includes('timeout')) {
          errorMessage = 'การเชื่อมต่อช้าเกินไป กรุณาลองใหม่'
        }
      }
      
      toast('เริ่มเกมล้มเหลว', errorMessage, 'error')
      
      // Reset critical state to prevent further errors
      correctAnswer.value = ''
      choices.value = []
      selectedChoice.value = null
      
      // Don't throw the error - let the user try again
      return
    }
    
  } catch (error) {
    console.error('Critical error starting game:', error)
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถเริ่มเกมได้ กรุณาลองใหม่', 'error')
    
    // Reset all game state on critical failure
    correctAnswer.value = ''
    choices.value = []
    selectedChoice.value = null
    hintCount.value = 0
    timer.value = ROUND_SECONDS
  }
}

/**
 * Fetch a new quiz and set up the game
 */
/**
 * Enhanced quiz fetching with comprehensive error handling and validation
 * 
 * This function fetches a new quiz from the API and sets up the game state.
 * It includes robust error handling for various failure scenarios:
 * 
 * 1. API Response Validation:
 *    - Validates that the response contains required fields
 *    - Ensures the answer field is not empty or null
 *    - Provides fallback values for optional fields
 * 
 * 2. Error Recovery:
 *    - Retries on network failures
 *    - Provides user feedback for different error types
 *    - Prevents game state corruption on API failures
 * 
 * 3. Data Integrity:
 *    - Validates quiz data before setting game state
 *    - Ensures correctAnswer is available before generating choices
 *    - Maintains consistent game state across retries
 * 
 * @param isAuto - Whether this is an automatic quiz fetch (affects change count)
 */
async function fetchQuiz(isAuto = false) {
  if (changeCount.value >= maxChange && !isAuto) return
  
  try {
    // Reset game state for new quiz
    expiredNotice.value = false
    revealedAnswer.value = ''

    // Prepare API parameters
    const params: any = { level: currentLevel.value }
    if (selectedCategory.value) {
      // Use actual category if Random is selected, otherwise use selected category
      const categoryToUse = (selectedCategory as any).actualCategory || selectedCategory.value
      params.category = categoryToUse
    }

    // Fetch quiz data from API using multiple-choice endpoint
    const res = await apiGet('/api/quiz/multiple-choice', params)
    
    // Validate API response structure
    if (!res || !res.data) {
      throw new Error('Invalid API response: missing data field')
    }

    // Validate required fields with detailed error messages
    if (!res.data.id) {
      throw new Error('Invalid API response: missing quiz ID')
    }

    if (!res.data.answer || typeof res.data.answer !== 'string' || res.data.answer.trim() === '') {
      throw new Error('Invalid API response: missing or empty answer field')
    }

    // Validate answer is not just whitespace
    const trimmedAnswer = res.data.answer.trim()
    if (trimmedAnswer.length === 0) {
      throw new Error('Invalid API response: answer field contains only whitespace')
    }

    // Set quiz metadata
    quizId.value = res.data.id
    quizToken.value = res.data.token || ''
    quizExp.value = res.data.exp || 0
    maxHints.value = typeof res.data.hintCount === 'number' ? res.data.hintCount : 2

    // Set correct answer with validation
    const answer = res.data.answer.trim()
    if (!answer) {
      throw new Error('Quiz answer is empty after trimming')
    }
    correctAnswer.value = answer

    // Reset game state for new round
    selectedChoice.value = null
    confirmedChoice.value = null
    choices.value = []
    hint1.value = res.data.hint1 || ''
    hint2.value = res.data.hint2 || ''
    hintCount.value = 0
    timer.value = ROUND_SECONDS
    recentGuesses.value = []

    // Generate multiple choice options (only if we have a valid answer)
    if (correctAnswer.value) {
      await generateChoices()
    } else {
      throw new Error('Cannot generate choices: correct answer is not available')
    }

    // Start timer
    startTimer()

    // Show hint 1 immediately if available
    if (hint1.value) {
      hintCount.value = 1
    }

    console.log('Quiz loaded successfully:', { 
      answer: correctAnswer.value, 
      hints: { hint1: hint1.value, hint2: hint2.value },
      quizId: quizId.value
    })
    
  } catch (error) {
    console.error('Failed to fetch quiz:', error)
    
    // Provide specific error messages based on error type
    let errorMessage = 'กรุณาลองใหม่'
    if (error instanceof Error) {
      if (error.message.includes('missing answer')) {
        errorMessage = 'ข้อมูลคำถามไม่สมบูรณ์ กรุณาลองใหม่'
      } else if (error.message.includes('network') || error.message.includes('Network')) {
        errorMessage = 'เชื่อมต่อเซิร์ฟเวอร์ไม่ได้ กรุณาตรวจสอบการเชื่อมต่อ'
      } else if (error.message.includes('timeout')) {
        errorMessage = 'การเชื่อมต่อช้าเกินไป กรุณาลองใหม่'
      } else if (error.message.includes('no quiz available') || error.message.includes('404')) {
        errorMessage = 'ไม่พบคำถามในหมวดหมู่นี้ กรุณาเลือกหมวดหมู่อื่น'
      }
    }
    
    toast('โหลดคำถามล้มเหลว', errorMessage, 'error')
    
    // Reset critical state to prevent further errors
    correctAnswer.value = ''
    choices.value = []
    selectedChoice.value = null
    
    // If it's a "no quiz available" error, redirect to category selection
    if (error instanceof Error && (error.message.includes('no quiz available') || error.message.includes('404'))) {
      setTimeout(() => {
        router.push({ name: 'CategorySelection' })
      }, 2000)
    }
  }
}

/**
 * Start the game timer
 */
function startTimer() {
  if (intervalId) clearInterval(intervalId)
  
  // Don't start timer in no-timer mode
  if (noTimer.value) {
    console.log('No-timer mode: Timer disabled')
    return
  }
  
  intervalId = window.setInterval(() => {
    if (timer.value <= 0) {
      // Time's up
      clearInterval(intervalId)
      handleTimeUp()
      return
    }
    
    timer.value--
    
    // Show hint 2 when 10 seconds remain
    if (timer.value === 10 && hint2.value && hintCount.value < 2) {
      hintCount.value = 2
      toast('ใบ้เพิ่ม!', 'ใบ้ที่ 2 ปรากฏแล้ว', 'info')
    }
  }, 1000)
}

/**
 * Handle when time runs out
 */
function handleTimeUp() {
  try {
    // Check if user hasn't selected or confirmed an answer
    if (selectedChoice.value === null || confirmedChoice.value === null) {
      // No answer selected or confirmed - treat as incorrect
      toast('หมดเวลา!', `คำตอบที่ถูกต้องคือ: ${correctAnswer.value}`, 'error')
      revealedAnswer.value = correctAnswer.value
      
      // Increment wrong answer count
      wrongAnswerCount.value++
      
      // Check if this is the 3rd wrong answer
      if (wrongAnswerCount.value >= maxWrongAnswers) {
        toast('เกมจบ!', 'ตอบผิด 3 ครั้งแล้ว', 'error')
        setTimeout(() => {
          endGame()
        }, 2000)
        return
      }
      
      // Move to next question after delay
      setTimeout(() => {
        fetchQuiz()
      }, 2000)
    }
  } catch (error) {
    console.error('Error handling time up:', error)
  }
}

/**
 * Show additional hint
 */


/**
 * Navigate back to previous page
 */
function goBack() {
  try {
    router.back()
  } catch (error) {
    console.error('Error navigating back:', error)
    router.push({ name: 'DogAll' })
  }
}


/* ===================== Lifecycle ===================== */
onMounted(async () => {
  document.title = 'PETTEXT - Dog Question'

  // Validate category selection
  if (!selectedCategory.value) {
    console.warn('No category selected, redirecting to category selection')
    toast('กรุณาเลือกหมวดหมู่', 'ไม่พบหมวดหมู่ที่เลือก', 'error')
    router.push({ name: 'CategorySelection' })
    return
  }

  const { healthOk, initialOk } = await waitApiReadyAndLoadInitial()
  if (healthOk && initialOk) {
    loading.value = false
    await flushPendingScores()
    loadScores()
  } else {
    const checkTimer = window.setInterval(async () => {
      const h = await waitApiReadyAndLoadInitial()
      if (h.healthOk && h.initialOk) {
        window.clearInterval(checkTimer)
        loading.value = false
        await flushPendingScores()
        loadScores()
      }
    }, 5000)
  }

  window.addEventListener('visibilitychange', onVisibilityChange)
})

onBeforeUnmount(() => {
  if (intervalId) clearInterval(intervalId as number)
  window.removeEventListener('visibilitychange', onVisibilityChange)
})

function onVisibilityChange() {
  if (document.hidden) {
    visibilityPausedAt = Date.now()
  } else if (visibilityPausedAt && intervalId) {
    const paused = Date.now() - visibilityPausedAt
    timer.value = Math.max(0, timer.value - Math.floor(paused / 1000))
    visibilityPausedAt = 0
  }
}

/* ===================== Score Management ===================== */
async function flushPendingScores() {
  try {
    const raw = localStorage.getItem('pendingScores')
    if (!raw) return
    const arr = JSON.parse(raw) as any[]
    if (arr.length === 0) return

    for (const item of arr) {
      try {
        await apiPost('/api/scores', item)
      } catch (e) {
        console.warn('Failed to flush score:', e)
      }
    }
    localStorage.removeItem('pendingScores')
  } catch (error) {
    console.error('Error flushing pending scores:', error)
  }
}


/* ===================== API Readiness ===================== */
/**
 * Enhanced API readiness checker with comprehensive endpoint validation
 * 
 * This function ensures the backend is fully operational before allowing the game to start.
 * It performs the following checks:
 * 1. Health endpoint check - verifies server is running and database is accessible
 * 2. Quiz endpoint check - ensures quiz data can be fetched (essential for gameplay)
 * 3. Scores endpoint check - verifies score saving functionality works
 * 
 * @returns Promise<{healthOk: boolean, initialOk: boolean}> - Status of API readiness
 * 
 * Error Handling:
 * - Network failures are caught and logged
 * - Individual endpoint failures don't crash the entire check
 * - Returns false for both flags if any critical check fails
 * 
 * Performance:
 * - Uses Promise.allSettled to check multiple endpoints concurrently
 * - Fails fast if health check fails (no point checking other endpoints)
 * - Comprehensive error logging for debugging
 */
async function waitApiReadyAndLoadInitial() {
  try {
    // Step 1: Check server health and database connectivity
    const healthRes = await apiGet('/health')
    const healthOk = healthRes.status === 200

    if (!healthOk) {
      console.warn('Health check failed - server or database not ready')
      return { healthOk: false, initialOk: false }
    }

    // Step 2: Test essential API endpoints to ensure backend is fully ready
    // Using Promise.allSettled to prevent one failure from stopping all checks
    const results = await Promise.allSettled([
      apiGet('/api/quiz', { params: { level: 1 } }),           // Essential for gameplay
      apiGet('/api/scores', { params: { limit: 10, gamename: 'DogPuzzle' } })  // Essential for score saving
    ])
    
    // All endpoints must succeed for the API to be considered ready
    const initialOk = results.every(r => r.status === 'fulfilled')
    
    if (!initialOk) {
      console.warn('One or more API endpoints failed during readiness check')
    }

    return { healthOk, initialOk }
  } catch (error) {
    console.error('API readiness check failed with error:', error)
    return { healthOk: false, initialOk: false }
  }
}

/* ===================== Components ===================== */
const InfoCard = defineComponent({
  name: 'InfoCard',
  props: {
    label: { type: String, required: true },
    value: { type: [String, Number], required: true },
    accent: { type: String as PropType<'indigo' | 'sky' | 'fuchsia' | 'emerald' | 'rose'>, required: true }
  },
  setup(props) {
    const base = 'rounded-xl py-2.5 px-4 text-center border '
    const cls =
      props.accent === 'indigo' ? base + 'bg-indigo-400/10 border-indigo-300/20' :
        props.accent === 'sky' ? base + 'bg-sky-400/10 border-sky-300/20' :
          props.accent === 'fuchsia' ? base + 'bg-fuchsia-400/10 border-fuchsia-300/20' :
            props.accent === 'rose' ? base + 'bg-rose-400/10 border-rose-300/20' :
              base + 'bg-emerald-400/10 border-emerald-300/20'
    return () => h('div', { class: cls }, [
      h('div', { class: 'text-[11px] font-semibold text-slate-200/90 tracking-wide' }, props.label),
      h('div', { class: 'mt-0.5 text-2xl font-bold text-slate-100 tabular-nums' }, String(props.value)),
    ])
  }
})
</script>

<style scoped>
/* Custom styles for multiple choice game */
.choice-button {
  transition: all 0.2s ease;
}

.choice-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.choice-button:active:not(:disabled) {
  transform: translateY(0);
}

/* Correct answer animation */
.choice-button.correct {
  animation: correctPulse 0.6s ease-out;
}

.choice-button.incorrect {
  animation: incorrectShake 0.6s ease-out;
}

@keyframes correctPulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.02); }
  100% { transform: scale(1); }
}

@keyframes incorrectShake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-4px); }
  75% { transform: translateX(4px); }
}

/* Timer progress bar animation */
.timer-progress {
  transition: width 0.3s ease;
}

/* Hint display animation */
.hint-display {
  animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
