<template>
  <div class="min-h-screen relative overflow-x-hidden theme-modern">
    <!-- Modern gradient background -->
    <div class="pointer-events-none absolute inset-0 -z-10" aria-hidden="true">
      <div
        class="absolute inset-0 bg-[radial-gradient(90%_70%_at_70%_100%,rgba(99,102,241,0.45),transparent_60%),radial-gradient(60%_60%_at_0%_0%,rgba(59,130,246,0.35),transparent_60%),linear-gradient(180deg,#0b1020,#0b1120)]"
      />
      <div class="absolute -bottom-16 right-10 h-80 w-80 rounded-full blur-3xl opacity-40 bg-indigo-500/30" />
      <div class="absolute -top-12 left-[-4rem] h-72 w-72 rounded-full blur-3xl opacity-30 bg-fuchsia-500/25" />
    </div>

    <!-- Loading overlay -->
    <div v-if="loading" class="fixed inset-0 bg-black/60 flex items-center justify-center z-[90]">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-base md:text-lg text-indigo-100 font-semibold">กำลังโหลด...</span>
      </div>
    </div>

    <div class="w-full max-w-6xl px-4 py-8 mx-auto">
      <!-- Header -->
      <header class="flex flex-col gap-3 items-center mb-6">
        <div class="w-full flex items-center justify-between">
          <button
            @click="goBack"
            type="button"
            class="inline-flex items-center gap-2 px-3 py-2 rounded-xl border bg-white/5 border-white/10 text-slate-100 hover:bg-white/10 transition shadow-sm"
          >
            <svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 18l-6-6 6-6"></path>
            </svg>
            ย้อนกลับ
          </button>

          <h1 class="text-3xl md:text-4xl font-extrabold tracking-wide text-indigo-300/90 uppercase text-center flex-1 drop-shadow-sm">
            PolaJigsaw
          </h1>

          <div class="w-[90px] sm:w-[120px]" />
        </div>
        <p class="text-slate-300/80 text-sm md:text-base text-center">
          ต่อจิ๊กซอนะ • ไม่ใช่ • ต่อ ธนภพ
        </p>
      </header>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Left: Control / Wizard -->
        <section class="lg:col-span-1">
          <div class="rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-5 space-y-5">
            <!-- Step indicator -->
            <div class="flex flex-wrap items-center gap-2 text-xs">
              <span class="px-2 py-1 rounded border"
                    :class="stepBadgeClass(1)">1 เลือกรูป</span>
              <span class="text-slate-400">→</span>
              <span class="px-2 py-1 rounded border"
                    :class="stepBadgeClass(2)">2 เลือกระดับ</span>
              <span class="text-slate-400">→</span>
              <span class="px-2 py-1 rounded border"
                    :class="stepBadgeClass(3)">3 ใส่ชื่อ</span>
              <span class="text-slate-400">→</span>
              <span class="px-2 py-1 rounded border"
                    :class="inPlay ? 'bg-emerald-600 text-white border-emerald-600' : 'bg-white/10 text-slate-300 border-white/15'">
                4 เริ่มเล่น
              </span>
            </div>

            <!-- Step 1: choose image (4 random each round) -->
            <div class="space-y-2 rounded-xl p-3 transition"
                 :class="step === 1 ? 'ring-2 ring-indigo-400/70 step-glow' : 'ring-0'">
              <div class="flex items-center justify-between">
                <label class="block text-sm font-semibold text-slate-100">เลือกรูป (สุ่ม 4 ภาพ)</label>
                <button @click="reshuffleRoundImages" class="text-xs text-slate-300 hover:underline">เปลี่ยนชุดรูป</button>
              </div>
              <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
                <button
                  v-for="(img, idx) in roundImages"
                  :key="'r-'+idx"
                  @click="chooseSample(img)"
                  class="rounded-xl overflow-hidden border-2 transition hover:scale-[1.02] focus:outline-none focus:ring-2 focus:ring-indigo-400/60"
                  :class="imageSrc===img ? 'border-indigo-400' : 'border-transparent'"
                >
                  <img :src="img" class="h-24 w-full object-cover" />
                </button>
              </div>

              <div v-if="apiImages.length" class="pt-2">
                <div class="text-xs text-slate-400 mb-1">จาก API</div>
                <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
                  <button
                    v-for="(img, idx) in apiImages"
                    :key="'api-'+idx"
                    @click="chooseSample(img)"
                    class="rounded-xl overflow-hidden border-2 transition hover:scale-[1.02] focus:outline-none focus:ring-2 focus:ring-indigo-400/60"
                    :class="imageSrc===img ? 'border-indigo-400' : 'border-transparent'"
                  >
                    <img :src="img" class="h-24 w-full object-cover" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Step 2: difficulty -->
            <div class="space-y-3 pt-2 rounded-xl p-3 transition"
                 :class="step === 2 ? 'ring-2 ring-indigo-400/70 step-glow' : 'ring-0'">
              <label class="block text-sm font-semibold text-slate-100">ระดับความยาก</label>
              <div class="grid grid-cols-2 gap-2">
                <button
                  v-for="d in difficulties"
                  :key="d.key"
                  @click="setDifficulty(d)"
                  class="px-3 py-2 rounded-xl border text-sm font-semibold transition hover:shadow"
                  :class="gridSize===d.size
                    ? 'border-indigo-400 bg-indigo-400/10 text-indigo-100'
                    : 'border-white/15 bg-white/5 text-slate-200'"
                >
                  {{ d.label }} ({{ d.size }}×{{ d.size }})
                </button>
              </div>
            </div>

            <!-- Step 3: player name -->
            <div class="space-y-2 pt-2 rounded-xl p-3 transition"
                 :class="step === 3 ? 'ring-2 ring-indigo-400/70 step-glow' : 'ring-0'">
              <label class="block text-sm font-semibold text-slate-100">ชื่อผู้เล่น (สำหรับบันทึกคะแนน)</label>
              <input
                v-model="playerName"
                maxlength="64"
                placeholder="พิมพ์ชื่อของคุณ"
                @focus="step = 3"
                class="w-full rounded-lg border px-3 py-2 text-sm bg-white/5 border-white/15 text-slate-100 placeholder:slate-400 focus:ring-2 focus:ring-indigo-400/60 outline-none"
              />
              <p class="text-xs text-slate-400">* ใช้สำหรับแสดงบนตารางสถิติ</p>
            </div>

            <!-- Actions -->
            <div class="flex flex-wrap gap-2 pt-1">
              <button
                @click="startGame"
                :disabled="!canStart"
                class="flex-1 px-4 py-2 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-400 shadow disabled:opacity-50"
                title="ต้องเลือกรูปและระดับ และกรอกชื่อก่อน"
              >
                เริ่มเกม
              </button>
              <button @click="shuffle" :disabled="!inPlay"
                      class="px-4 py-2 rounded-xl border font-semibold hover:bg-white/10 disabled:opacity-50
                             bg-white/5 border-white/15 text-slate-200">
                สุ่มใหม่
              </button>
              <button @click="togglePause" :disabled="!inPlay"
                      class="px-4 py-2 rounded-xl border font-semibold hover:bg-white/10 disabled:opacity-50
                             bg-white/5 border-white/15 text-slate-200">
                {{ paused ? 'เล่นต่อ' : 'พัก' }}
              </button>
              <button @click="peek" :disabled="!inPlay || hintsLeft<=0"
                      class="px-4 py-2 rounded-xl bg-amber-500 text-slate-900 font-bold hover:bg-amber-400 disabled:opacity-50">
                Hint ({{ hintsLeft }})
              </button>
            </div>

            <!-- Stats header -->
            <header class="space-y-3 mt-2">
              <div class="grid grid-cols-3 gap-3">
                <div class="rounded-xl py-2.5 px-4 text-center bg-sky-400/10 border border-sky-300/20">
                  <div class="text-xs font-semibold text-sky-200 tracking-wide">เวลา</div>
                  <div class="mt-0.5 text-2xl font-bold text-sky-100 tabular-nums">{{ formattedTime }}</div>
                </div>
                <div class="rounded-xl py-2.5 px-4 text-center bg-white/10 border border-white/15">
                  <div class="text-xs font-semibold text-slate-200 tracking-wide">จำนวนสลับ</div>
                  <div class="mt-0.5 text-2xl font-bold text-slate-100 tabular-nums">{{ moves }}</div>
                </div>
                <div class="rounded-xl py-2.5 px-4 text-center bg-fuchsia-400/10 border border-fuchsia-300/20">
                  <div class="text-xs font-semibold text-fuchsia-200 tracking-wide">คะแนน</div>
                  <div class="mt-0.5 text-2xl font-bold text-fuchsia-100 tabular-nums">{{ score }}</div>
                </div>
              </div>
              <div class="w-full h-2 bg-white/10 rounded-full overflow-hidden">
                <div
                  class="h-full transition-all duration-300 bg-gradient-to-r from-sky-400 via-indigo-400 to-fuchsia-400"
                  :style="{ width: inPlay ? Math.min(100, (seconds/(gridSize*30 || 1))*100) + '%' : '0%' }"
                  role="progressbar"
                  :aria-valuenow="seconds"
                  aria-label="ความคืบหน้าเวลา"
                />
              </div>
            </header>
          </div>

          <!-- Leaderboard -->
          <section class="mt-6 w-full rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-5">
            <div class="flex items-center justify-between">
              <h4 class="text-lg font-bold text-indigo-100">สถิติผู้เล่น TOP 10</h4>
              <button @click="loadTopScores"
                      class="text-xs px-3 py-1 rounded-lg border border-white/15 bg-white/5 text-slate-200 hover:bg-white/10 transition"
                      title="รีเฟรชสถิติ" type="button">
                รีเฟรช
              </button>
            </div>

            <p v-if="topScores.length === 0" class="text-slate-300/70 text-sm mt-2">
              ยังไม่มีสถิติ แสดงเมื่อมีการบันทึกคะแนน
            </p>

            <ul v-else class="mt-2 divide-y divide-white/10">
              <li v-for="(item, idx) in topScores"
                  :key="item.name + '_' + item.score + '_' + idx"
                  class="py-2 flex items-center justify-between text-sm text-slate-100">
                <div class="flex items-center gap-2 min-w-0">
                  <span class="w-6 text-center">
                    <template v-if="idx === 0">🥇</template>
                    <template v-else-if="idx === 1">🥈</template>
                    <template v-else-if="idx === 2">🥉</template>
                    <template v-else>{{ idx + 1 }}.</template>
                  </span>
                  <span class="font-medium truncate max-w-[10rem] sm:max-w-[14rem]" :title="item.name">
                    {{ item.name }}
                  </span>
                </div>
                <div class="tabular-nums font-semibold">{{ item.score }} คะแนน</div>
              </li>
            </ul>
          </section>
        </section>

        <!-- Right: Board -->
        <section class="lg:col-span-2">
          <div class="rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-4 md:p-6">
            <div class="flex items-center justify-between mb-3">
              <h2 class="font-bold text-indigo-100">กระดาน</h2>
              <div class="text-xs text-slate-300/80">ลากสลับชิ้นส่วน หรือแตะเลือก 2 ชิ้นเพื่อสลับ</div>
            </div>

            <!-- Responsive square board -->
            <div ref="boardWrapperRef" class="mx-auto w-full max-w-[820px]">
              <div
                class="relative rounded-2xl overflow-hidden shadow-inner ring-2 ring-white/15 bg-slate-900/20"
                :style="boardStyle"
              >
                <!-- Tiles -->
                <button
                  v-for="(tile, idx) in tiles"
                  :key="idx"
                  class="tile focus:outline-none"
                  :class="tileClass(idx)"
                  :style="tileStyle(tile)"
                  draggable="true"
                  @dragstart="onDragStart(idx)"
                  @dragenter.prevent="onDragEnter(idx)"
                  @dragover.prevent
                  @dragleave="onDragLeave(idx)"
                  @drop="onDrop(idx)"
                  @click="onTap(idx)"
                />

                <!-- Paused overlay -->
                <div v-if="paused && inPlay" class="absolute inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-40">
                  <div class="text-white text-lg font-bold">หยุดชั่วคราว</div>
                </div>

                <!-- Hint preview overlay -->
                <transition name="fade">
                  <img
                    v-if="peeking && imageSrc"
                    :src="imageSrc"
                    class="absolute inset-0 w-full h-full object-cover opacity-80 z-40"
                    alt="hint"
                  />
                </transition>
              </div>
            </div>

            <div class="mt-4 flex flex-col md:flex-row items-center justify-between gap-2">
              <div class="text-xs text-slate-400">Tip: ยิ่งสลับน้อยและเร็ว คะแนนยิ่งสูง</div>
            </div>
          </div>
        </section>
      </div>
    </div>

    <!-- Win Modal -->
    <transition name="pop">
      <div v-if="winModal" class="fixed inset-0 bg-black/60 flex items-center justify-center z-[95]">
        <div class="rounded-2xl overflow-hidden shadow-[0_30px_80px_rgba(0,0,0,0.55)] w-[min(92vw,520px)] border border-white/10 bg-white/5 backdrop-blur-xl">
          <div class="px-6 py-5 text-white bg-gradient-to-r from-indigo-600/90 via-indigo-500/90 to-fuchsia-600/90">
            <h3 class="text-2xl font-extrabold text-white mb-1">เยี่ยมมาก! ต่อสำเร็จ 🎉</h3>
            <p class="text-indigo-100/90">สรุปผลงานของคุณ</p>
          </div>

          <div class="p-6 text-slate-100">
            <div class="grid grid-cols-2 gap-3 mb-4">
              <div class="rounded-2xl border border-sky-300/25 bg-sky-400/10 px-4 py-3">
                <div class="text-xs font-semibold text-sky-200">เวลา</div>
                <div class="mt-1 text-2xl md:text-3xl font-extrabold text-sky-100 tabular-nums">{{ formattedTime }}</div>
              </div>
              <div class="rounded-2xl border border-white/15 bg-white/10 px-4 py-3">
                <div class="text-xs font-semibold text-slate-200">จำนวนสลับ</div>
                <div class="mt-1 text-2xl md:text-3xl font-extrabold text-slate-100 tabular-nums">{{ moves }}</div>
              </div>
              <div class="rounded-2xl border border-fuchsia-300/25 bg-fuchsia-400/10 px-4 py-3">
                <div class="text-xs font-semibold text-fuchsia-200">คะแนน</div>
                <div class="mt-1 text-2xl md:text-3xl font-extrabold text-fuchsia-100 tabular-nums">{{ score }}</div>
              </div>
              <div class="rounded-2xl border border-emerald-300/25 bg-emerald-400/10 px-4 py-3">
                <div class="text-xs font-semibold text-emerald-200">โหมด</div>
                <div class="mt-1 text-2xl md:text-3xl font-extrabold text-emerald-100">{{ currentDiff?.label }}</div>
              </div>
            </div>
            <div class="flex items-center justify-end gap-2">
              <button @click="winModal=false"
                      class="px-4 py-2 rounded-xl border bg-white/5 border-white/15 text-slate-100 hover:bg-white/10">
                ปิด
              </button>
              <button @click="startGame"
                      class="px-4 py-2 rounded-xl bg-indigo-500 text-white font-bold hover:bg-indigo-400">
                เล่นอีกครั้ง
              </button>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import type { CSSProperties } from 'vue'
import { useRouter } from 'vue-router'
import catwalk from '../assets/images/catwalk.png'
import catwalk2 from '../assets/images/catwalk2.png'
import api from '../services/api'

/** ---------- Config & API ---------- */
const GAME_NAME = 'PolaJigsaw'
const TOP_LIMIT = 10

/** ---------- Loading ---------- */
const loading = ref(true)
const catwalkImages = [catwalk, catwalk2]
const catwalkIndex = ref(0)
let catwalkInterval: number | undefined

/** ---------- Router ---------- */
const router = useRouter()
const goBack = () => router.back()

/** ---------- Assets: load .png ---------- */
const allAssetImages = ref<string[]>([])
const roundImages = ref<string[]>([])
const apiImages = ref<string[]>([])

function loadAssetsPng() {
  const modules = import.meta.glob('../assets/images/**/*.png', { eager: true, as: 'url' }) as Record<string, string>
  allAssetImages.value = Object.values(modules)
}

function pickRoundImages() {
  const pool = allAssetImages.value.slice()
  shuffleArray(pool)
  roundImages.value = pool.slice(0, 4)
}

function reshuffleRoundImages() {
  pickRoundImages()
  imageSrc.value = ''
  gridSize.value = 0
  step.value = 1
}

async function loadImagesFromApiIfAny() {
  try {
    const res = await api.get('/api/jigsaw/images')
    if (Array.isArray(res.data)) {
      apiImages.value = res.data.filter(Boolean)
    }
  } catch {}
}

/** ---------- Wizard ---------- */
const imageSrc = ref<string>('')
function chooseSample(src: string) {
  imageSrc.value = src
  if (step.value < 2) step.value = 2
}

type Diff = { key: string; label: string; size: number; multiplier: number; hints: number }
const difficulties = ref<Diff[]>([
  { key: 'easy',   label: 'ง่าย',     size: 3, multiplier: 1.0, hints: 3 },
  { key: 'medium', label: 'ปานกลาง', size: 4, multiplier: 1.4, hints: 3 },
  { key: 'hard',   label: 'ยาก',     size: 5, multiplier: 1.9, hints: 2 },
  { key: 'insane', label: 'โหด',     size: 6, multiplier: 2.6, hints: 1 },
])
const gridSize = ref<number>(0)
const currentDiff = computed<Diff|undefined>(() => difficulties.value.find(d => d.size === gridSize.value))
function setDifficulty(d: Diff) { gridSize.value = d.size; if (step.value < 3) step.value = 3 }

const playerName = ref<string>('')

const step = ref<number>(1)
const canStart = computed(() => !!imageSrc.value && gridSize.value > 0 && playerName.value.trim().length > 0)
function stepBadgeClass(n: number) {
  const baseInactive = 'bg-white/10 text-slate-300 border-white/15'
  const active = 'bg-indigo-500 text-white border-indigo-500'
  const glow = 'animate-pulse-soft'
  if (n === 1) return step.value === 1 ? `${active} ${glow}` : (step.value > 1 ? active : baseInactive)
  if (n === 2) return step.value === 2 ? `${active} ${glow}` : (step.value > 2 ? active : baseInactive)
  if (n === 3) return step.value === 3 && !inPlay.value ? `${active} ${glow}` : (step.value > 3 || inPlay.value ? active : baseInactive)
  return baseInactive
}

/** ---------- Board / Tiles ---------- */
const tiles = ref<number[]>([])
const inPlay = ref(false)
const paused = ref(false)
const selectedIndex = ref<number|null>(null)
const draggingIndex = ref<number|null>(null)
const dragOverIndex = ref<number|null>(null)

const boardWrapperRef = ref<HTMLElement | null>(null)
const boardSidePx = ref<number>(560)

const boardStyle = computed<CSSProperties>(() => ({
  width: '100%',
  maxWidth: 'min(100%, 820px)',
  aspectRatio: '1 / 1',
  display: 'grid',
  gridTemplateColumns: `repeat(${gridSize.value || 1}, 1fr)`,
  gridTemplateRows: `repeat(${gridSize.value || 1}, 1fr)`,
  position: 'relative',
  height: 'auto',
  minHeight: `${Math.max(280, Math.min(boardSidePx.value, 820))}px`,
}))

let ro: ResizeObserver | null = null
function setupResizeObserver() {
  if (!boardWrapperRef.value) return
  ro = new ResizeObserver(entries => {
    for (const e of entries) {
      const w = e.contentRect.width
      boardSidePx.value = w
    }
  })
  ro.observe(boardWrapperRef.value)
}

function tileStyle(tileIndex: number): CSSProperties {
  const n = gridSize.value || 1
  const row = Math.floor(tileIndex / n)
  const col = tileIndex % n
  return {
    backgroundImage: imageSrc.value ? `url('${imageSrc.value}')` : 'none',
    backgroundSize: `${n * 100}% ${n * 100}%`,
    backgroundPosition: `${(col / (n - 1 || 1)) * 100}% ${(row / (n - 1 || 1)) * 100}%`,
    aspectRatio: '1 / 1',
    transition: 'transform 120ms ease, box-shadow 120ms ease, border-color 120ms ease',
  }
}

function tileClass(idx: number) {
  const isCorrect = tileCorrect(idx)
  const isSelected = selectedIndex.value === idx
  const isDragging = draggingIndex.value === idx
  const isDropTarget = dragOverIndex.value === idx && draggingIndex.value !== null

  return [
    isCorrect ? 'ring-2 ring-emerald-400 z-20' : 'ring-1 ring-white/40',
    (isSelected || isDragging) ? 'border-2 border-amber-500 shadow-xl z-30 animate-pulse-soft bg-white/5' : 'border border-transparent',
    isDropTarget ? 'outline outline-2 outline-sky-500 outline-offset-[-2px] z-30' : '',
    'rounded-none'
  ].join(' ')
}

function tileCorrect(posIdx: number) { return tiles.value[posIdx] === posIdx }

/** ---------- Timer / Moves / Score ---------- */
const seconds = ref(0)
const timerHandle = ref<number|undefined>(undefined)
const moves = ref(0)
const hintsLeft = ref(3)
const peeking = ref(false)

const formattedTime = computed(() => {
  const s = seconds.value
  const mm = Math.floor(s / 60).toString().padStart(2, '0')
  const ss = (s % 60).toString().padStart(2, '0')
  return `${mm}:${ss}`
})

const score = computed(() => {
  if (!inPlay.value && !winModal.value) return 0
  const base = 10000
  const diffMul = (currentDiff.value?.multiplier ?? 1)
  const t = Math.max(1, seconds.value)
  const m = Math.max(1, moves.value)
  const hintPenalty = (initialHints.value - hintsLeft.value) * 5
  const raw = Math.round((base * diffMul) / (t + m + hintPenalty))
  return Math.max(1, raw)
})

/** ---------- Leaderboard (DB) ---------- */
type TopScore = { name: string; score: number; gamename?: string }
const topScores = ref<TopScore[]>([])

async function saveScoreToDb(name: string, s: number) {
  await api.post('/api/scores', { name: name.trim(), score: s, gamename: GAME_NAME })
}
async function loadTopScores() {
  const res = await api.get('/api/scores', { params: { limit: TOP_LIMIT, gamename: GAME_NAME } })
  topScores.value = Array.isArray(res.data) ? res.data.slice(0, TOP_LIMIT) : []
}

/** ---------- Game Flow ---------- */
const winModal = ref(false)
const initialHints = ref(0)

function buildSolved(n: number) { return Array.from({ length: n * n }, (_, i) => i) }
function shuffleArray<T>(arr: T[]) {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[arr[i], arr[j]] = [arr[j], arr[i]]
  }
}

function startTimer() {
  stopTimer()
  timerHandle.value = window.setInterval(() => { if (!paused.value) seconds.value += 1 }, 1000)
}
function stopTimer() {
  if (timerHandle.value) { clearInterval(timerHandle.value); timerHandle.value = undefined }
}

function startGame() {
  if (!canStart.value) return
  inPlay.value = true
  step.value = 4
  paused.value = false
  winModal.value = false
  seconds.value = 0
  moves.value = 0
  initialHints.value = currentDiff.value?.hints ?? 2
  hintsLeft.value = initialHints.value
  selectedIndex.value = null
  draggingIndex.value = null
  dragOverIndex.value = null

  const n = gridSize.value
  const arr = buildSolved(n)
  do { shuffleArray(arr) } while (arr.every((v, i) => v === i))
  tiles.value = arr
  startTimer()
}

function shuffle() {
  if (!inPlay.value) return
  const arr = tiles.value.slice()
  shuffleArray(arr)
  if (arr.every((v, i) => v === i)) shuffleArray(arr)
  tiles.value = arr
  moves.value += 1
}

function togglePause() { if (inPlay.value) paused.value = !paused.value }

async function peek() {
  if (!inPlay.value || hintsLeft.value <= 0) return
  hintsLeft.value -= 1
  peeking.value = true
  await new Promise(res => setTimeout(res, 2000))
  peeking.value = false
}

/** Drag & Drop logic + Drop target highlight */
let dragFrom = -1
function onDragStart(idx: number) {
  if (!inPlay.value || paused.value) return
  dragFrom = idx
  draggingIndex.value = idx
  dragOverIndex.value = null
}
function onDragEnter(idx: number) {
  if (!inPlay.value || paused.value) return
  if (draggingIndex.value === null || idx === draggingIndex.value) { dragOverIndex.value = null; return }
  dragOverIndex.value = idx
}
function onDragLeave(idx: number) {
  if (dragOverIndex.value === idx) dragOverIndex.value = null
}
function onDrop(idx: number) {
  if (!inPlay.value || paused.value) return
  if (dragFrom === -1 || dragFrom === idx) { draggingIndex.value = null; dragOverIndex.value = null; return }
  swapTiles(dragFrom, idx)
  dragFrom = -1
  draggingIndex.value = null
  dragOverIndex.value = null
}
function onTap(idx: number) {
  if (!inPlay.value || paused.value) return
  if (selectedIndex.value === null) {
    selectedIndex.value = idx
  } else if (selectedIndex.value === idx) {
    selectedIndex.value = null
  } else {
    swapTiles(selectedIndex.value, idx)
    selectedIndex.value = null
  }
}
function swapTiles(a: number, b: number) {
  const arr = tiles.value.slice()
  ;[arr[a], arr[b]] = [arr[b], arr[a]]
  tiles.value = arr
  moves.value += 1
  checkWin()
}
async function checkWin() {
  if (tiles.value.every((v, i) => v === i)) {
    inPlay.value = false
    stopTimer()
    winModal.value = true
    try { await saveScoreToDb(playerName.value, score.value) } catch {}
    loadTopScores()
  }
}

/** ---------- Mount / Unmount ---------- */
onMounted(() => {
  document.title = 'PETTEXT - PolaJigsaw'

  catwalkInterval = window.setInterval(() => {
    catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length
  }, 200)
  setTimeout(() => { loading.value = false }, 600)

  loadAssetsPng()
  pickRoundImages()
  loadImagesFromApiIfAny()
  loadTopScores()

  setupResizeObserver()
})

onBeforeUnmount(() => {
  stopTimer()
  if (catwalkInterval) clearInterval(catwalkInterval)
  if (ro && boardWrapperRef.value) ro.unobserve(boardWrapperRef.value)
})
</script>

<style scoped>
.tile {
  position: relative;
  width: 100%;
  height: 100%;
  cursor: grab;
  user-select: none;
  background-repeat: no-repeat;
  outline: none;
}
.tile:active { cursor: grabbing; }
.tabular-nums { font-variant-numeric: tabular-nums; }

/* Step highlight (border glow) */
@keyframes glow {
  0% { box-shadow: 0 0 0 0 rgba(99,102,241,.55); }
  70% { box-shadow: 0 0 0 8px rgba(99,102,241,0); }
  100% { box-shadow: 0 0 0 0 rgba(99,102,241,0); }
}
.step-glow { animation: glow 1.6s infinite; }

/* Soft pulse utility */
@keyframes pulseSoft {
  0%, 100% { transform: scale(1); filter: saturate(1); }
  50% { transform: scale(1.02); filter: saturate(1.05); }
}
.animate-pulse-soft { animation: pulseSoft 1.3s ease-in-out infinite; }

/* Transitions */
.fade-enter-active,.fade-leave-active{ transition: opacity .2s ease }
.fade-enter-from,.fade-leave-to{ opacity: 0 }
.pop-enter-active,.pop-leave-active{ transition: transform .2s ease, opacity .2s ease }
.pop-enter-from,.pop-leave-to{ transform: scale(.96); opacity: 0 }

/* Dark-friendly theme */
.theme-modern { color-scheme: dark; }
</style>
