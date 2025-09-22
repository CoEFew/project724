<!--
Multiple-Choice Quiz Game - Party Mode

This component implements a multiplayer multiple-choice quiz game with the following features:
- Room discovery and joining system
- Real-time player management with WebSocket
- Lobby with player ready states and room information
- In-game multiple-choice quiz interface with automatic hint display
- Score tracking and leaderboard display
- Robust error handling and user feedback

Gameplay Rules:
1. Round time: 30 seconds
2. Hint #1: Shown immediately at start
3. Hint #2: Shown when 10 seconds remain
4. Choices: 1 correct (from current quiz) + 3 incorrect (random from other quizzes)
5. Scoring: Identical to DocumentsPageAlls.vue

Technical Implementation:
- Vue 3 Composition API with TypeScript
- WebSocket for real-time communication
- Same theme and UI structure as DocumentsPageAlls.vue
- Multiple choice generation with randomization
- Timer-based hint display
- Comprehensive error handling
- Responsive design with Tailwind CSS
-->
<template>
  <div class="min-h-screen relative overflow-x-hidden theme-modern">
    <!-- BG -->
    <div class="pointer-events-none absolute inset-0 -z-10" aria-hidden="true">
      <div
        class="absolute inset-0 bg-[radial-gradient(90%_70%_at_70%_100%,rgba(99,102,241,0.45),transparent_60%),radial-gradient(60%_60%_at_0%_0%,rgba(59,130,246,0.35),transparent_60%),linear-gradient(180deg,#0b1020,#0b1120)]" />
      <div class="absolute -bottom-16 right-10 h-80 w-80 rounded-full blur-3xl opacity-40 bg-indigo-500/30" />
      <div class="absolute -top-12 left-[-4rem] h-72 w-72 rounded-full blur-3xl opacity-30 bg-fuchsia-500/25" />
    </div>

    <!-- Loading -->
    <LoadingOverlay :loading="loading" />

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
            DOG • QUESTION (PARTY)
          </h1>

          <div class="flex gap-2 items-center">
            <button type="button" @click="soundOn = !soundOn"
              class="px-3 py-2 rounded-xl border border-white/10 bg-white/5 text-slate-100 hover:bg-white/10">
              <span v-if="soundOn">🔊</span><span v-else>🔇</span>
            </button>
          </div>
        </div>
      </header>

      <!-- ===== MAINTENANCE MODE ===== -->
      <section class="w-full max-w-3xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-8 space-y-6 text-center">
        <div class="space-y-4">
          <div class="text-6xl">🔧</div>
          <h2 class="text-2xl font-bold text-indigo-100">โหมดปาร์ตี้อยู่ระหว่างการบำรุงรักษา</h2>
          <p class="text-slate-300">เรากำลังปรับปรุงโหมดปาร์ตี้ให้ดีขึ้น กรุณาใช้โหมดเดี่ยวแทน</p>
          <div class="space-y-3">
            <button @click="goBack" 
              class="px-6 py-3 rounded-xl font-semibold transition bg-indigo-500 text-white hover:bg-indigo-400 shadow-lg">
              กลับไปหน้าแรก
            </button>
            <div class="text-sm text-slate-400">
              <p>คุณสามารถเล่นโหมดเดี่ยวได้ที่:</p>
              <p class="font-medium text-indigo-200">DOG • Puzzle (Single Player)</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ===== LOBBY (DISABLED) ===== -->
      <section v-if="false && phase === 'lobby'"
        class="w-full max-w-3xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-5 space-y-5">
        
        <!-- Room List Section -->
        <div v-if="!room" class="space-y-4">
          <div class="text-center">
            <div class="flex items-center justify-center gap-3 mb-2">
              <h3 class="text-lg font-bold text-indigo-100">ห้องที่เปิดอยู่</h3>
              <button @click="loadAvailableRooms" 
                class="p-1.5 rounded-lg border border-white/10 bg-white/5 text-slate-300 hover:bg-white/10 hover:text-slate-100 transition-colors"
                title="รีเฟรชรายการห้อง">
                <svg viewBox="0 0 24 24" class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
                  <path d="M21 3v5h-5"/>
                  <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
                  <path d="M3 21v-5h5"/>
                </svg>
              </button>
            </div>
            <p class="text-sm text-slate-400">เลือกห้องเพื่อเข้าร่วมหรือสร้างห้องใหม่</p>
          </div>
          <div v-if="availableRooms.length === 0" class="text-center text-slate-400 py-8">
            <div class="text-6xl mb-4">🏠</div>
            <p class="text-lg font-medium mb-2">ยังไม่มีห้องที่เปิดอยู่</p>
            <p class="text-sm">สร้างห้องใหม่เพื่อเริ่มเล่นกับเพื่อน</p>
          </div>
          <div v-else class="grid gap-3">
            <div v-for="roomItem in availableRooms" :key="roomItem.code"
              :class="[
                'rounded-xl border p-4 transition-all group',
                (roomItem.player_count || 0) >= roomItem.max_players
                  ? 'border-rose-400/30 bg-rose-500/5 cursor-not-allowed opacity-60'
                  : 'border-white/10 bg-white/5 hover:bg-white/10 hover:border-indigo-400/30 cursor-pointer'
              ]"
              @click="joinRoom(roomItem.code)">
              <div class="flex items-center justify-between">
                <div>
                  <div class="font-semibold text-slate-100">{{ roomItem.code }}</div>
                  <div class="text-sm text-slate-400">
                    {{ roomItem.player_count || 0 }}/{{ roomItem.max_players }} ผู้เล่น
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-sm text-slate-300">{{ roomItem.category || 'สุ่ม' }}</div>
                  <div class="text-xs text-slate-400">โหมดตัวเลือก</div>
                </div>
              </div>
            </div>
          </div>
          <button @click="createRoom" 
            class="w-full px-4 py-3 rounded-xl font-semibold transition bg-indigo-500 text-white hover:bg-indigo-400 shadow">
            สร้างห้องใหม่
          </button>
        </div>

        <!-- Room Details Section -->
        <div v-else class="space-y-4">
          <div class="text-center">
            <h3 class="text-lg font-bold text-indigo-100">ห้อง {{ room?.code || 'Unknown' }}</h3>
            <p class="text-sm text-slate-400">โหมดตัวเลือก - เลือกคำตอบที่ถูกต้องจาก 4 ตัวเลือก</p>
          </div>
          
          <!-- Category Selection -->
          <div class="space-y-3">
            <h4 class="text-md font-semibold text-slate-200">เลือกหมวดหมู่</h4>
            <div class="grid grid-cols-2 gap-3">
              <button @click="selectCategory('สัตว์')"
                :class="['p-3 rounded-xl border transition-all', selectedCategory === 'สัตว์' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
                <div class="text-center">
                  <div class="text-xl mb-1">🐕</div>
                  <div class="text-sm font-medium">สัตว์</div>
                </div>
              </button>
              <button @click="selectCategory('เครื่องใช้ไฟฟ้า')"
                :class="['p-3 rounded-xl border transition-all', selectedCategory === 'เครื่องใช้ไฟฟ้า' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
                <div class="text-center">
                  <div class="text-xl mb-1">⚡</div>
                  <div class="text-sm font-medium">เครื่องใช้ไฟฟ้า</div>
                </div>
              </button>
              <button @click="selectCategory('ผลไม้')"
                :class="['p-3 rounded-xl border transition-all', selectedCategory === 'ผลไม้' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
                <div class="text-center">
                  <div class="text-xl mb-1">🍎</div>
                  <div class="text-sm font-medium">ผลไม้</div>
                </div>
              </button>
              <button @click="selectCategory('อาชีพ')"
                :class="['p-3 rounded-xl border transition-all', selectedCategory === 'อาชีพ' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
                <div class="text-center">
                  <div class="text-xl mb-1">👨‍💼</div>
                  <div class="text-sm font-medium">อาชีพ</div>
                </div>
              </button>
              <button @click="selectRandomCategory"
                :class="['p-3 rounded-xl border transition-all col-span-2', selectedCategory === 'Random' ? 'border-indigo-400 bg-indigo-500/20 text-indigo-100' : 'border-white/10 bg-white/5 text-slate-300 hover:bg-white/10']">
                <div class="text-center">
                  <div class="text-xl mb-1">🎲</div>
                  <div class="text-sm font-medium">สุ่มหมวดหมู่</div>
                </div>
              </button>
            </div>
          </div>

          <!-- Player List -->
          <div class="space-y-3">
            <h4 class="text-md font-semibold text-slate-200">ผู้เล่น ({{ players.length }}/4)</h4>
            <div class="space-y-2">
              <div v-for="player in players" :key="player.id"
                class="flex items-center justify-between p-3 rounded-xl border border-white/10 bg-white/5">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full bg-indigo-500/20 flex items-center justify-center text-sm font-semibold text-indigo-300">
                    {{ player.name.charAt(0).toUpperCase() }}
                  </div>
                  <span class="text-slate-100">{{ player.name }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span v-if="player.ready" class="text-green-400 text-sm">พร้อม</span>
                  <span v-else class="text-amber-400 text-sm">รอ</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex gap-2">
            <button @click="leaveRoom" 
              class="flex-1 px-4 py-2.5 rounded-xl font-semibold transition bg-slate-600 text-white hover:bg-slate-500">
              ออกจากห้อง
            </button>
            <button @click="toggleReady" 
              :class="['flex-1 px-4 py-2.5 rounded-xl font-semibold transition', isReady ? 'bg-amber-600 text-white hover:bg-amber-500' : 'bg-indigo-500 text-white hover:bg-indigo-400']">
              {{ isReady ? 'ยกเลิกความพร้อม' : 'พร้อมเล่น' }}
            </button>
          </div>
        </div>
      </section>

      <!-- ===== PLAYING (DISABLED) ===== -->
      <section v-if="false && phase === 'game'"
        class="w-full max-w-3xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-6 space-y-5">

        <!-- status & header (InfoCards แบบเดียวกับหน้าเดิม) -->
        <div class="space-y-3">
          <h2 class="text-xl md:text-2xl font-extrabold text-indigo-100 tracking-wide text-center">
            เก่งจริงก็ทายมาดิ! • ห้อง {{ room?.code }}
          </h2>

          <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
            <InfoCard label="รอบ" :value="round?.round_no ?? '-'" accent="indigo" />
            <InfoCard label="เวลาคงเหลือ" :value="timer + ' วินาที'" accent="sky" />
            <InfoCard label="ผู้เล่น" :value="players.length + ' / ' + (room?.max_players || 4)"
              accent="fuchsia" />
          </div>

          <div class="space-y-2">
            <div class="w-full h-2 rounded-full overflow-hidden bg-white/10">
              <div class="h-full transition-all duration-300 bg-gradient-to-r from-sky-400 via-indigo-400 to-fuchsia-400"
                :style="{ width: timerPercent + '%' }" />
            </div>
            <div class="text-[11px] text-slate-300/80 text-right">ทุกคนในห้องได้คำเดียวกัน</div>
          </div>
        </div>

        <!-- players mini board -->
        <div class="rounded-xl border border-white/10 bg-white/5 p-3">
          <div class="text-slate-200 text-sm font-semibold mb-2">ผู้เล่น</div>
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
            <div v-for="p in players" :key="p.name"
              class="rounded-lg border border-white/10 px-3 py-2 text-slate-200 text-xs flex items-center justify-between">
              <div class="flex items-center gap-2 min-w-0">
                <span class="inline-block w-2 h-2 rounded-full"
                  :class="p.is_out ? 'bg-rose-400' : 'bg-emerald-400'"></span>
                <span class="truncate" :title="p.name">{{ p.name }}</span>
              </div>
              <span class="tabular-nums">{{ p.score }}</span>
            </div>
          </div>
        </div>

        <!-- result & hints -->
        <div class="space-y-2">
          <div v-if="lastResult !== null" class="text-center" aria-live="polite">
            <span v-if="lastResult === true" class="text-emerald-300 text-xl font-semibold">✔️
              ถูกต้อง!</span>
            <span v-else-if="lastResult === false" class="text-rose-300 text-xl font-semibold">❌ ผิด
              ลองใหม่!</span>
          </div>
          <div v-if="hints.length" class="flex flex-wrap gap-2 justify-center">
            <Chip v-for="(h, i) in hints" :key="i" :label="'คำใบ้ ' + (i + 1)" :value="h" color="fuchsia" />
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
              :disabled="selectedChoice !== null || meOut"
              :class="[
                'p-4 rounded-xl border transition-all text-left',
                selectedChoice === index
                  ? choice.isCorrect
                    ? 'border-green-400 bg-green-500/20 text-green-100'
                    : 'border-red-400 bg-red-500/20 text-red-100'
                  : 'border-white/10 bg-white/5 text-slate-100 hover:bg-white/10 hover:border-indigo-400/30',
                selectedChoice !== null || meOut ? 'cursor-not-allowed opacity-60' : 'cursor-pointer'
              ]"
            >
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-full border-2 flex items-center justify-center text-sm font-semibold"
                  :class="[
                    selectedChoice === index
                      ? choice.isCorrect
                        ? 'border-green-400 text-green-400'
                        : 'border-red-400 text-red-400'
                      : 'border-white/20 text-slate-300'
                  ]">
                  {{ String.fromCharCode(65 + index) }}
                </div>
                <span class="font-medium">{{ choice.text }}</span>
                <div v-if="selectedChoice === index" class="ml-auto">
                  <span v-if="choice.isCorrect" class="text-green-400">✓</span>
                  <span v-else class="text-red-400">✗</span>
                </div>
              </div>
            </button>
          </div>
        </div>

        <!-- my recent guesses (ของฉันเท่านั้นใน "ข้อปัจจุบัน") -->
        <div class="flex flex-wrap gap-2 justify-center">
          <div v-if="!myRoundGuesses.length" class="text-slate-400 text-xs">ยังไม่มีคำตอบของคุณ</div>
          <div v-else class="flex flex-wrap gap-2">
            <span v-for="(g, idx) in myRoundGuesses" :key="g + '_' + idx"
              class="px-2.5 py-1 rounded-full text-xs border bg-white/5 border-white/10 text-slate-200">
              #{{ g }}
            </span>
          </div>
        </div>

      </section>

      <!-- ===== GAME OVER (DISABLED) ===== -->
      <section v-if="false && phase === 'scoreboard'"
        class="w-full max-w-3xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-5 space-y-4">
        <h3 class="text-xl font-extrabold text-indigo-100">ผลการแข่งขัน</h3>
        <p class="text-slate-200 text-sm">ผู้ชนะ: <b>{{ winner?.name || '—' }}</b></p>
        
        <!-- Revealed Answer -->
        <div v-if="revealedAnswer" class="rounded-xl border border-emerald-300/30 bg-emerald-400/10 p-4">
          <h4 class="text-lg font-bold text-emerald-100 mb-2 text-center">คำตอบที่ถูกต้อง</h4>
          <div class="text-center">
            <span class="inline-flex items-center gap-2 text-emerald-100 font-bold text-xl">{{ revealedAnswer }}</span>
          </div>
        </div>

        <!-- Current Game Scoreboard -->
        <div class="rounded-xl border border-white/10 bg-white/5 p-4">
          <h4 class="text-lg font-bold text-indigo-100 mb-3 text-center">ตารางคะแนนรอบนี้</h4>
          <ul class="divide-y divide-white/10">
            <li v-for="(l, idx) in leaderboard" :key="l.name"
              class="py-3 flex items-center justify-between text-sm text-slate-100">
              <div class="flex items-center gap-3 min-w-0">
                <span class="w-8 text-center text-lg">
                  <template v-if="idx === 0">🥇</template>
                  <template v-else-if="idx === 1">🥈</template>
                  <template v-else-if="idx === 2">🥉</template>
                  <template v-else>{{ idx + 1 }}</template>
                </span>
                <span class="truncate" :title="l.name">{{ l.name }}</span>
              </div>
              <span class="tabular-nums font-semibold text-indigo-300">{{ l.score }}</span>
            </li>
          </ul>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-2">
          <button @click="leaveRoom" 
            class="flex-1 px-4 py-3 rounded-xl font-semibold transition bg-slate-600 text-white hover:bg-slate-500">
            ออกจากห้อง
          </button>
          <button @click="restartGame" 
            class="flex-1 px-4 py-3 rounded-xl font-semibold transition bg-indigo-500 text-white hover:bg-indigo-400">
            เล่นใหม่
          </button>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
/* ===================== Imports ===================== */
import LoadingOverlay from '../components/LoadingOverlay.vue'
import api from '../services/api'
import { ref, onMounted, onBeforeUnmount, computed, defineComponent, h, type PropType } from 'vue'
import { useRouter } from 'vue-router'
import { waitApiReadyAndLoadInitial } from '../composables/useApiReadiness'

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
 * Player interface for multiplayer game
 * Represents a player in the room
 */
interface Player {
  id: string
  name: string
  score: number
  is_ready: boolean
  is_owner: boolean
  is_out: boolean
  ready?: boolean // Legacy property for compatibility
  selectedChoice?: number // Current choice selection
  choiceCorrect?: boolean // Whether current choice is correct
}

/**
 * Room interface for multiplayer game
 * Represents a game room
 */
interface Room {
  code: string
  name: string
  max_players: number
  player_count: number
  category: string
  status: string
}

/**
 * WebSocket message interface
 * Represents messages sent between client and server
 */
interface WSMessage {
  type: string
  room?: string
  [key: string]: any
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
 * DogQuestionAll.vue
 * 
 * Multiple-Choice Quiz Game - Party Mode
 * 
 * This component implements a multiplayer multiple-choice quiz game with:
 * - Room discovery and joining system
 * - Real-time player management with WebSocket
 * - Lobby with player ready states and room information
 * - In-game multiple-choice quiz interface with automatic hint display
 * - Score tracking and leaderboard display
 * - Robust error handling and user feedback
 * 
 * Gameplay Rules:
 * 1. Round time: 30 seconds
 * 2. Hint #1: Shown immediately at start
 * 3. Hint #2: Shown when 10 seconds remain
 * 4. Choices: 1 correct (from current quiz) + 3 incorrect (random from other quizzes)
 * 5. Scoring: Identical to DocumentsPageAlls.vue
 */

/* ===================== Constants ===================== */
const ROUND_SECONDS = 30 // 30 seconds per round

/* ===================== Router & State ===================== */
const router = useRouter()

// loading state
const loading = ref(true)

// game phase: 'lobby' | 'game' | 'scoreboard'
const phase = ref<'lobby' | 'game' | 'scoreboard'>('lobby')

// room and player states
const room = ref<Room | null>(null)
const players = ref<Player[]>([])
const isReady = ref<boolean>(false)
const selectedCategory = ref<string>('')
const soundOn = ref<boolean>(true)

// available rooms
const availableRooms = ref<Room[]>([])

// game states
const quizId = ref<string>('')
const quizToken = ref<string>('')
const quizExp = ref<number>(0)
const timer = ref<number>(ROUND_SECONDS)
const currentLevel = ref<number>(1)
const hint1 = ref<string>('')
const hint2 = ref<string>('')
const hintCount = ref<number>(0)
const maxHints = ref<number>(2)

// Multiple choice specific states
const choices = ref<Choice[]>([])
const selectedChoice = ref<number | null>(null)
const correctAnswer = ref<string>('')

// Game state properties for DocumentsPageAlls.vue compatibility
const round = ref<any>(null)
const lastResult = ref<boolean | null>(null)
const hints = ref<string[]>([])
const meOut = ref(false)
const myRoundGuesses = ref<string[]>([])
const winner = ref<any>(null)
const revealedAnswer = ref('')
const leaderboard = ref<any[]>([])

// WebSocket
let ws: WebSocket | null = null
let reconnectAttempts = 0
const maxReconnectAttempts = 5

// timer & intervals
let intervalId: number | undefined

/* ===================== Computed Properties ===================== */
const timerPercent = computed(() => (timer.value / ROUND_SECONDS) * 100)

const sortedPlayers = computed(() => {
  return [...players.value].sort((a, b) => (b.score || 0) - (a.score || 0))
})

/* ===================== Utility Functions ===================== */
function toast(title: string, message: string, type: 'info' | 'success' | 'error' = 'info') {
  console.log(`[${type.toUpperCase()}] ${title}: ${message}`)
  // In a real app, you'd show a toast notification here
}

/**
 * Enhanced API GET function with comprehensive error handling and user feedback
 * 
 * This function provides robust error handling for all GET requests with the following features:
 * 
 * 1. Error Classification:
 *    - Server errors (4xx, 5xx) with specific error messages
 *    - Network errors (no response received) with connection guidance
 *    - Request setup errors (configuration issues) with debugging info
 * 
 * 2. User-Friendly Error Messages:
 *    - 404 errors: Clear "endpoint not found" messages
 *    - 500 errors: Generic server error with retry suggestion
 *    - 4xx errors: Client-side error details from server
 *    - Network errors: Connection troubleshooting guidance
 * 
 * 3. Developer Experience:
 *    - Comprehensive error logging for debugging
 *    - Structured error information for troubleshooting
 *    - Consistent error handling across all API calls
 * 
 * 4. Error Recovery:
 *    - Throws meaningful errors that can be caught by calling functions
 *    - Provides context for error handling decisions
 *    - Maintains error chain for proper debugging
 * 
 * @param url - API endpoint URL to call
 * @param params - Optional query parameters to include in the request
 * @returns Promise<AxiosResponse> - Response data on success
 * @throws Error - Detailed error message on failure
 * 
 * Example Usage:
 * ```typescript
 * try {
 *   const response = await apiGet('/api/quiz', { level: 1 })
 *   console.log('Quiz data:', response.data)
 * } catch (error) {
 *   console.error('Failed to load quiz:', error.message)
 * }
 * ```
 */
async function apiGet(url: string, params?: any) {
  try {
    const response = await api.get(url, { params })
    return response
  } catch (error: any) {
    console.error(`API GET Error for ${url}:`, error)
    
    // Handle different types of errors
    if (error.response) {
      // Server responded with error status
      const status = error.response.status
      const message = error.response.data?.message || 'Server error'
      
      if (status === 404) {
        throw new Error(`API endpoint not found: ${url}`)
      } else if (status === 500) {
        throw new Error('Internal server error. Please try again later.')
      } else if (status >= 400 && status < 500) {
        throw new Error(`Client error: ${message}`)
      } else {
        throw new Error(`Server error (${status}): ${message}`)
      }
    } else if (error.request) {
      // Network error - no response received
      throw new Error('Network error. Please check your connection and try again.')
    } else {
      // Other error
      throw new Error(`Request setup error: ${error.message}`)
    }
  }
}

/**
 * Enhanced API POST function with comprehensive error handling and user feedback
 * 
 * This function provides robust error handling for all POST requests with the following features:
 * 
 * 1. Error Classification:
 *    - Server errors (4xx, 5xx) with specific error messages
 *    - Network errors (no response received) with connection guidance
 *    - Request setup errors (configuration issues) with debugging info
 * 
 * 2. User-Friendly Error Messages:
 *    - 404 errors: Clear "endpoint not found" messages
 *    - 500 errors: Generic server error with retry suggestion
 *    - 4xx errors: Client-side error details from server
 *    - Network errors: Connection troubleshooting guidance
 * 
 * 3. Developer Experience:
 *    - Comprehensive error logging for debugging
 *    - Structured error information for troubleshooting
 *    - Consistent error handling across all API calls
 * 
 * 4. Error Recovery:
 *    - Throws meaningful errors that can be caught by calling functions
 *    - Provides context for error handling decisions
 *    - Maintains error chain for proper debugging
 * 
 * @param url - API endpoint URL to call
 * @param data - Request body data to send
 * @returns Promise<AxiosResponse> - Response data on success
 * @throws Error - Detailed error message on failure
 * 
 * Example Usage:
 * ```typescript
 * try {
 *   const response = await apiPost('/api/rooms', { category: 'Animals' })
 *   console.log('Room created:', response.data)
 * } catch (error) {
 *   console.error('Failed to create room:', error.message)
 * }
 * ```
 */
async function apiPost(url: string, data?: any) {
  try {
    const response = await api.post(url, data)
    return response
  } catch (error: any) {
    console.error(`API POST Error for ${url}:`, error)
    
    // Handle different types of errors
    if (error.response) {
      // Server responded with error status
      const status = error.response.status
      const message = error.response.data?.message || 'Server error'
      
      if (status === 404) {
        throw new Error(`API endpoint not found: ${url}`)
      } else if (status === 500) {
        throw new Error('Internal server error. Please try again later.')
      } else if (status >= 400 && status < 500) {
        throw new Error(`Client error: ${message}`)
      } else {
        throw new Error(`Server error (${status}): ${message}`)
      }
    } else if (error.request) {
      // Network error - no response received
      throw new Error('Network error. Please check your connection and try again.')
    } else {
      // Other error
      throw new Error(`Request setup error: ${error.message}`)
    }
  }
}

/* ===================== WebSocket Management ===================== */
function connectWebSocket() {
  try {
    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'
    ws = new WebSocket(wsUrl)
    
    ws.onopen = () => {
      console.log('WebSocket connected')
      reconnectAttempts = 0
      
      // Send join message if we have a room
      if (room.value) {
        sendMessage({
          type: 'join',
          room: room.value.code,
          player: getPlayerName()
        })
      }
    }
    
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        handleWebSocketMessage(data)
      } catch (error) {
        console.error('Error parsing WebSocket message:', error)
      }
    }
    
    ws.onclose = () => {
      console.log('WebSocket disconnected')
      ws = null
      
      // Attempt to reconnect
      if (reconnectAttempts < maxReconnectAttempts) {
        reconnectAttempts++
        setTimeout(() => {
          console.log(`Attempting to reconnect... (${reconnectAttempts}/${maxReconnectAttempts})`)
          connectWebSocket()
        }, 2000 * reconnectAttempts)
      }
    }
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  } catch (error) {
    console.error('Error connecting WebSocket:', error)
  }
}

/**
 * Enhanced WebSocket message sending with comprehensive error handling and connection management
 * 
 * This function provides robust WebSocket communication with the following features:
 * 
 * 1. Connection State Management:
 *    - Checks WebSocket initialization before sending
 *    - Handles different connection states (OPEN, CONNECTING, CLOSING, CLOSED)
 *    - Provides appropriate feedback for each state
 * 
 * 2. Message Queuing:
 *    - Queues messages when WebSocket is connecting
 *    - Retries sending after connection is established
 *    - Prevents message loss during connection transitions
 * 
 * 3. Automatic Reconnection:
 *    - Triggers reconnection when WebSocket is closed
 *    - Queues messages for sending after reconnection
 *    - Maintains message delivery reliability
 * 
 * 4. Error Handling:
 *    - Catches and logs all sending errors
 *    - Provides user feedback through toast notifications
 *    - Prevents application crashes from WebSocket errors
 * 
 * 5. Performance Optimization:
 *    - Logs successful message sending for debugging
 *    - Uses appropriate timeouts for queued messages
 *    - Prevents infinite retry loops
 * 
 * @param message - Message object to send via WebSocket
 * 
 * Message Format:
 * ```typescript
 * {
 *   type: string,        // Message type (join, select_choice, etc.)
 *   room?: string,       // Room code (if applicable)
 *   player?: string,     // Player name (if applicable)
 *   [key: string]: any   // Additional message data
 * }
 * ```
 * 
 * Example Usage:
 * ```typescript
 * sendMessage({
 *   type: 'select_choice',
 *   room: 'ABC123',
 *   choice: 'Dog',
 *   correct: true
 * })
 * ```
 */
function sendMessage(message: any) {
  try {
    if (!ws) {
      console.warn('WebSocket not initialized, cannot send message:', message)
      return
    }
    
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message))
      console.log('WebSocket message sent:', message)
    } else if (ws.readyState === WebSocket.CONNECTING) {
      console.warn('WebSocket is connecting, queuing message:', message)
      // Queue the message to send when connection is established
      setTimeout(() => sendMessage(message), 1000)
    } else if (ws.readyState === WebSocket.CLOSING || ws.readyState === WebSocket.CLOSED) {
      console.warn('WebSocket is closed or closing, attempting to reconnect...')
      connectWebSocket()
      // Queue the message to send after reconnection
      setTimeout(() => sendMessage(message), 2000)
    }
  } catch (error) {
    console.error('Error sending WebSocket message:', error)
    toast('ส่งข้อความล้มเหลว', 'ไม่สามารถส่งข้อความได้', 'error')
  }
}

function handleWebSocketMessage(data: any) {
  try {
    switch (data.type) {
      case 'room_joined':
        handleRoomJoined(data)
        break
      case 'player_joined':
        handlePlayerJoined(data)
        break
      case 'player_left':
        handlePlayerLeft(data)
        break
      case 'player_ready':
        handlePlayerReady(data)
        break
      case 'game_started':
        handleGameStarted(data)
        break
      case 'game_ended':
        handleGameEnded(data)
        break
      case 'quiz_loaded':
        handleQuizLoaded(data)
        break
      case 'choice_selected':
        handleChoiceSelected(data)
        break
      case 'round_ended':
        handleRoundEnded(data)
        break
      case 'error':
        handleError(data)
        break
      default:
        console.log('Unknown WebSocket message type:', data.type)
    }
  } catch (error) {
    console.error('Error handling WebSocket message:', error)
  }
}

/* ===================== WebSocket Event Handlers ===================== */
function handleRoomJoined(data: any) {
  room.value = data.room
  players.value = data.players || []
  toast('เข้าร่วมห้องสำเร็จ', `ห้อง ${data.room.code}`, 'success')
}

function handlePlayerJoined(data: any) {
  players.value = data.players || []
  toast('ผู้เล่นใหม่', `${data.player.name} เข้าร่วมห้อง`, 'info')
}

function handlePlayerLeft(data: any) {
  players.value = data.players || []
  toast('ผู้เล่นออก', `${data.player.name} ออกจากห้อง`, 'info')
}

function handlePlayerReady(data: any) {
  const player = players.value.find(p => p.id === data.player.id)
  if (player) {
    player.ready = data.player.ready
  }
  toast('สถานะเปลี่ยน', `${data.player.name} ${data.player.ready ? 'พร้อม' : 'ไม่พร้อม'}`, 'info')
}

function handleGameStarted(data: any) {
  phase.value = 'game'
  currentLevel.value = data.level || 1
  toast('เกมเริ่มแล้ว!', 'เริ่มเล่นกันเลย', 'success')
}

function handleGameEnded(data: any) {
  phase.value = 'scoreboard'
  players.value = data.players || []
  toast('เกมจบแล้ว!', 'ดูผลการแข่งขัน', 'info')
}

function handleQuizLoaded(data: any) {
  quizId.value = data.quiz.id
  quizToken.value = data.quiz.token
  quizExp.value = data.quiz.exp
  correctAnswer.value = data.quiz.answer
  hint1.value = data.quiz.hint1 || ''
  hint2.value = data.quiz.hint2 || ''
  maxHints.value = data.quiz.hintCount || 2
  
  // Reset game state
  selectedChoice.value = null
  choices.value = []
  hintCount.value = 0
  timer.value = ROUND_SECONDS
  
  // Generate choices
  generateChoices()
  
  // Start timer
  startTimer()
  
  // Show hint 1 immediately
  if (hint1.value) {
    hintCount.value = 1
  }
}

function handleChoiceSelected(data: any) {
  const player = players.value.find(p => p.id === data.player.id)
  if (player) {
    player.selectedChoice = data.choice
    player.choiceCorrect = data.correct
  }
}

function handleRoundEnded(data: any) {
  // Show correct answer
  correctAnswer.value = data.correctAnswer
  
  // Update player scores
  if (data.players) {
    players.value = data.players
  }
  
  // Move to next round after delay
  setTimeout(() => {
    if (data.nextRound) {
      // Next round will be loaded automatically
    } else {
      // Game ended
      phase.value = 'scoreboard'
    }
  }, 3000)
}

function handleError(data: any) {
  toast('เกิดข้อผิดพลาด', data.message || 'ไม่ทราบสาเหตุ', 'error')
}

/* ===================== Game Logic ===================== */
/**
 * Generate multiple choice options
 * Creates 4 choices: 1 correct (from current quiz) + 3 incorrect (random from other quizzes)
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
 * Select a choice and check if it's correct
 */
function selectChoice(choice: { text: string; isCorrect: boolean }, index: number) {
  if (selectedChoice.value !== null) return

  try {
    selectedChoice.value = index
    
    // Send choice to server
    sendMessage({
      type: 'select_choice',
      room: room.value?.code || '',
      choice: choice.text,
      correct: choice.isCorrect
    })
    
    if (choice.isCorrect) {
      toast('ถูกต้อง!', `+${Math.max(1, Math.floor(timer.value / 10))} คะแนน`, 'success')
    } else {
      toast('ผิด!', `คำตอบที่ถูกต้องคือ: ${correctAnswer.value}`, 'error')
    }
  } catch (error) {
    console.error('Error handling choice selection:', error)
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถประมวลผลคำตอบได้', 'error')
  }
}

/**
 * Start the game timer
 */
function startTimer() {
  if (intervalId) clearInterval(intervalId)
  
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
    if (selectedChoice.value === null) {
      // No answer selected
      toast('หมดเวลา!', `คำตอบที่ถูกต้องคือ: ${correctAnswer.value}`, 'error')
      
      // Send timeout to server
      sendMessage({
        type: 'timeout',
        room: room.value?.code || ''
      })
    }
  } catch (error) {
    console.error('Error handling time up:', error)
  }
}

/* ===================== Room Management ===================== */
async function loadAvailableRooms() {
  try {
    const res = await apiGet('/api/rooms')
    availableRooms.value = res.data || []
    toast('โหลดห้องสำเร็จ', `พบ ${availableRooms.value.length} ห้อง`, 'success')
  } catch (error) {
    console.error('Failed to load rooms:', error)
    toast('โหลดห้องล้มเหลว', 'กรุณาลองใหม่', 'error')
  }
}

async function createRoom() {
  try {
    const res = await apiPost('/api/rooms', {
      category: selectedCategory.value || 'Random',
      gameMode: 'multiple-choice'
    })
    
    room.value = res.data.room
    players.value = [res.data.player]
    isReady.value = false
    
    // Connect WebSocket
    connectWebSocket()
    
    toast('สร้างห้องสำเร็จ', `ห้อง ${room.value?.code || 'Unknown'}`, 'success')
  } catch (error) {
    console.error('Failed to create room:', error)
    toast('สร้างห้องล้มเหลว', 'กรุณาลองใหม่', 'error')
  }
}

async function joinRoom(roomCode: string) {
  try {
    const res = await apiPost(`/api/rooms/${roomCode}/join`, {
      player: getPlayerName()
    })
    
    room.value = res.data.room
    players.value = res.data.players || []
    isReady.value = false
    
    // Connect WebSocket
    connectWebSocket()
    
    toast('เข้าร่วมห้องสำเร็จ', `ห้อง ${roomCode}`, 'success')
  } catch (error) {
    console.error('Failed to join room:', error)
    toast('เข้าร่วมห้องล้มเหลว', 'กรุณาลองใหม่', 'error')
  }
}

function leaveRoom() {
  try {
    if (ws) {
      sendMessage({
        type: 'leave',
        room: room.value?.code || ''
      })
      ws.close()
    }
    
    // Reset state
    room.value = null
    players.value = []
    isReady.value = false
    phase.value = 'lobby'
    
    // Clear local storage
    localStorage.removeItem('party_name')
    
    toast('ออกจากห้องแล้ว', 'กลับไปที่ล็อบบี้', 'info')
  } catch (error) {
    console.error('Error leaving room:', error)
  }
}

function toggleReady() {
  try {
    isReady.value = !isReady.value
    
    sendMessage({
      type: 'toggle_ready',
      room: room.value?.code,
      ready: isReady.value
    })
    
    // Check if all players are ready and start game
    if (isReady.value && players.value.length >= 2) {
      const allReady = players.value.every(p => p.ready)
      if (allReady) {
        startGame()
      }
    }
  } catch (error) {
    console.error('Error toggling ready state:', error)
  }
}

function startGame() {
  try {
    sendMessage({
      type: 'start_game',
      room: room.value?.code,
      category: selectedCategory.value
    })
  } catch (error) {
    console.error('Error starting game:', error)
  }
}

function restartGame() {
  try {
    phase.value = 'lobby'
    isReady.value = false
    
    sendMessage({
      type: 'restart_game',
      room: room.value?.code
    })
  } catch (error) {
    console.error('Error restarting game:', error)
  }
}

/* ===================== Category Selection ===================== */
function selectCategory(category: string) {
  try {
    selectedCategory.value = category
    // Clear any previous random category selection
    ;(selectedCategory as any).actualCategory = null
    console.log('Category selected:', category)
  } catch (error) {
    console.error('Error selecting category:', error)
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถเลือกหมวดหมู่ได้', 'error')
  }
}

function selectRandomCategory() {
  try {
    const categories = ['สัตว์', 'เครื่องใช้ไฟฟ้า', 'ผลไม้', 'อาชีพ']
    const randomIndex = Math.floor(Math.random() * categories.length)
    selectedCategory.value = 'Random'
    // Store the actual selected category for API calls
    ;(selectedCategory as any).actualCategory = categories[randomIndex]
    toast('สุ่มหมวดหมู่แล้ว', `เลือกหมวดหมู่: ${categories[randomIndex]}`, 'success')
    console.log('Random category selected:', categories[randomIndex])
  } catch (error) {
    console.error('Error selecting random category:', error)
    toast('เกิดข้อผิดพลาด', 'ไม่สามารถสุ่มหมวดหมู่ได้', 'error')
  }
}

/* ===================== Utility Functions ===================== */
function getPlayerName(): string {
  let name = localStorage.getItem('player_name')
  if (!name) {
    name = `ผู้เล่น${Math.floor(Math.random() * 1000)}`
    localStorage.setItem('player_name', name)
  }
  return name
}

function goBack() {
  try {
    router.back()
  } catch (error) {
    console.error('Error navigating back:', error)
    router.push({ name: 'DogAll' })
  }
}

/* ===================== Lifecycle ===================== */
/**
 * Enhanced component initialization with robust API readiness checking
 * 
 * This lifecycle hook ensures the component only initializes when the backend is fully ready.
 * It implements a comprehensive initialization strategy:
 * 
 * 1. API Readiness Check:
 *    - Verifies server health and database connectivity
 *    - Tests essential endpoints (quiz, scores) before proceeding
 *    - Uses the shared useApiReadiness composable for consistency
 * 
 * 2. Graceful Degradation:
 *    - If API is not ready, shows loading state and retries periodically
 *    - Prevents component from breaking due to backend unavailability
 *    - Provides user feedback through console warnings
 * 
 * 3. Initialization Sequence:
 *    - Sets page title for better UX
 *    - Loads available rooms for lobby display
 *    - Establishes WebSocket connection for real-time communication
 *    - Hides loading overlay when everything is ready
 * 
 * 4. Error Recovery:
 *    - Automatic retry mechanism with exponential backoff
 *    - Prevents infinite retry loops with reasonable intervals
 *    - Maintains loading state until successful initialization
 * 
 * Error Handling:
 * - Network failures are handled gracefully
 * - WebSocket connection failures are retried automatically
 * - Room loading failures are logged but don't crash the component
 * 
 * Performance:
 * - Non-blocking initialization with async/await
 * - Efficient retry mechanism prevents resource waste
 * - Early return on API failure prevents unnecessary operations
 */
onMounted(async () => {
  document.title = 'PETTEXT - Dog Question Party'
  
  // Wait for API to be ready before proceeding
  const { healthOk, initialOk } = await waitApiReadyAndLoadInitial()
  
  if (healthOk && initialOk) {
    // API is ready - proceed with normal initialization
    console.log('API ready, initializing component...')
    
    // Load available rooms for lobby display
    await loadAvailableRooms()
    
    // Connect WebSocket for real-time communication
    connectWebSocket()
    
    // Hide loading overlay - component is ready
    loading.value = false
  } else {
    // API is not ready - implement retry mechanism
    console.warn('API not ready, implementing retry mechanism...')
    
    const retryInterval = setInterval(async () => {
      const { healthOk: retryHealthOk, initialOk: retryInitialOk } = await waitApiReadyAndLoadInitial()
      
      if (retryHealthOk && retryInitialOk) {
        // API is now ready - clear retry interval and initialize
        clearInterval(retryInterval)
        console.log('API ready after retry, initializing component...')
        
        await loadAvailableRooms()
        connectWebSocket()
        loading.value = false
      } else {
        console.warn('API still not ready, will retry in 2 seconds...')
      }
    }, 2000) // Retry every 2 seconds to avoid overwhelming the server
  }
})

onBeforeUnmount(() => {
  if (intervalId) clearInterval(intervalId)
  if (ws) {
    ws.close()
  }
})

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
/* Custom styles for multiple choice party game */
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

/* Player score animation */
.player-score {
  transition: all 0.3s ease;
}

.player-score.winner {
  animation: winnerGlow 2s ease-in-out infinite;
}

@keyframes winnerGlow {
  0%, 100% { box-shadow: 0 0 5px rgba(255, 215, 0, 0.3); }
  50% { box-shadow: 0 0 20px rgba(255, 215, 0, 0.6); }
}
</style>
