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
    <div v-if="loading"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex flex-col items-center justify-center z-[90]">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-2 animate-bounce" />
        <span class="text-base md:text-lg text-indigo-100 font-semibold">กำลังโหลด...</span>
        <span class="mt-1 text-xs text-indigo-100/70" v-if="net.hasPending">กำลังเชื่อมต่อเซิร์ฟเวอร์…</span>
        <span class="mt-1 text-xs text-amber-200/80" v-if="net.isStalled">เซิร์ฟเวอร์กำลังเริ่มทำงาน
          ช้ากว่าปกติเล็กน้อย</span>
      </div>
    </div>

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
            DOG • Puzzle
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

      <!-- กล่องเกม (Glass card) -->
      <section
        class="w-full max-w-xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-6 space-y-5">
        <!-- สรุปสถานะ -->
        <div class="space-y-3">
          <h2 class="text-xl md:text-2xl font-extrabold text-indigo-100 tracking-wide text-center">เก่งจริงก็ทายมาดิ!
          </h2>

          <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
            <InfoCard label="คะแนน" :value="score" accent="indigo" />
            <InfoCard label="เวลาคงเหลือ" :value="timer + ' วินาที'" accent="sky" />
            <InfoCard label="ระดับ" :value="'Lv. ' + currentLevel" accent="fuchsia" />
          </div>

          <!-- Progress & proximity meter -->
          <div class="space-y-2">
            <div class="w-full h-2 rounded-full overflow-hidden bg-white/10" role="progressbar" :aria-valuenow="timer">
              <div
                class="h-full transition-all duration-300 bg-gradient-to-r from-sky-400 via-indigo-400 to-fuchsia-400"
                :style="{ width: timerPercent + '%' }" />
            </div>
            <div class="flex items-center justify-between text-[11px] text-slate-300/80">
              <div class="inline-flex items-center gap-2">
                <span class="h-2 w-2 rounded-full" :class="heatDotClass"></span>
                <span>ความใกล้เคียงคำตอบ: <strong class="tabular-nums">{{ heatLabel }}</strong></span>
              </div>

            </div>
          </div>
        </div>

        <!-- ฟอร์มตอบ -->
        <form class="w-full flex flex-col sm:flex-row items-stretch sm:items-center gap-2"
          @submit.prevent="handleSubmit">
          <div class="relative flex-1">
            <input ref="answerInput" v-model="guess" type="text" :placeholder="placeholder"
              class="rounded-xl px-4 py-2.5 text-base bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60 focus:border-indigo-400/60 w-full"
              :disabled="showModal" autocomplete="off" @keydown.enter.prevent="handleSubmit"
              @keydown.ctrl.h.prevent="showHint" @keydown.ctrl.l.prevent="() => fetchQuiz()" />
            <!-- duplicate guess warning -->
            <p v-if="duplicateWarning" class="absolute -bottom-5 left-1 text-[11px] text-amber-300">
              กรอกคำซ้ำกับที่เคยเดาแล้ว</p>
          </div>
          <button type="submit"
            class="px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
            :disabled="showModal || !guessSanitized">
            ส่งคำตอบ
          </button>
          <button type="button" @click="showHint"
            class="relative px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition bg-amber-500 text-slate-900 hover:bg-amber-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
            :disabled="showModal || hintCount >= maxHints">
            <span
              class="absolute -top-2 -right-3 bg-white/90 text-amber-800 rounded-full px-2 py-0.5 text-xs font-bold shadow">{{
      Math.min(hintCount, maxHints) }}/{{ maxHints }}</span>
            คำใบ้
          </button>
          <button type="button" @click="() => fetchQuiz()"
            class="relative px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition bg-emerald-500 text-white hover:bg-emerald-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
            :disabled="showModal || changeCount >= maxChange">
            <span
              class="absolute -top-2 -right-3 bg-white/90 text-emerald-800 rounded-full px-2 py-0.5 text-xs font-bold shadow">{{
      changeCount }}/{{ maxChange }}</span>
            เปลี่ยนคำ
          </button>

          <!-- ✅ ปุ่มจบเกม (เฉพาะโหมดไม่จับเวลา) -->
          <button v-if="noTimer" type="button" @click="endNoTimer"
            class="relative px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition bg-rose-500 text-white hover:bg-rose-400 disabled:opacity-50 disabled:cursor-not-allowed shadow">
            จบเกม
          </button>
        </form>

        <!-- ผลลัพธ์ + last guesses -->
        <div class="space-y-2">
          <div v-if="result !== null" class="text-center" aria-live="polite">
            <span v-if="result === true" class="text-emerald-300 text-xl font-semibold">✔️ ถูกต้อง!</span>
            <span v-else-if="result === false" class="text-rose-300 text-xl font-semibold">❌ ผิด ลองใหม่!</span>
          </div>
          <div v-if="recentGuesses.length" class="flex flex-wrap gap-2 justify-center">
            <span v-for="(g, i) in recentGuesses" :key="g.word + '_' + i"
              class="inline-flex items-center gap-2 px-2.5 py-1 rounded-full border text-xs"
              :class="g.correct ? 'bg-emerald-400/10 border-emerald-300/20 text-emerald-100' : 'bg-white/5 border-white/10 text-slate-200'">
              <span class="tabular-nums">#{{ i + 1 }}</span>
              <span class="font-medium">{{ g.word }}</span>
              <span v-if="typeof g.heat === 'number'" class="text-[10px] opacity-80">(ใกล้: {{ heatText(g.heat)
                }})</span>
            </span>
          </div>
        </div>

        <!-- แจ้งเตือนหมดอายุ token -->
        <div v-if="expiredNotice"
          class="text-center text-sm text-amber-200 bg-amber-500/10 border border-amber-400/30 rounded-lg px-3 py-2">
          หมดเวลาตอบสำหรับคำนี้แล้ว ระบบจะสุ่มคำใหม่ให้อัตโนมัติ
        </div>

        <!-- คำใบ้ -->
        <div class="flex flex-wrap items-center justify-center gap-2">
          <Chip v-if="hint1" label="คำใบ้ 1" :value="hint1" color="fuchsia" />
          <Chip v-if="hint2" label="คำใบ้ 2" :value="hint2" color="fuchsia" />

        </div>
      </section>

      <!-- สถิติ TOP 10 -->
      <section
        class="w-full max-w-xl mx-auto mt-6 rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-5">
        <div class="flex items-center justify-between">
          <h4 class="text-lg font-bold text-indigo-100">สถิติผู้เล่น TOP 10</h4>
          <button @click="loadScores"
            class="text-xs px-3 py-1 rounded-lg border border-white/10 bg-white/5 text-slate-200 hover:bg-white/10 transition"
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

      <!-- Toasts -->
      <transition-group name="toast" tag="div" class="fixed bottom-4 right-4 flex flex-col gap-2 z-[96]">
        <div v-for="t in toasts" :key="t.id" class="rounded-xl px-4 py-3 text-sm shadow border"
          :class="toastClass(t.type)">
          <strong class="mr-1">{{ t.title }}</strong>{{ t.message }}
        </div>
      </transition-group>
    </div>

    <!-- โมดัล: หมดเวลา + บันทึกคะแนน -->
    <div v-if="showModal"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[95] px-4" role="dialog"
      aria-modal="true" aria-labelledby="timeoutTitle" aria-describedby="timeoutDesc">
      <div class="w-full max-w-xl">
        <div
          class="relative rounded-3xl overflow-hidden shadow-[0_30px_80px_rgba(0,0,0,0.55)] border border-white/10 bg-white/5 backdrop-blur-xl">
          <div class="px-6 py-5 text-white bg-gradient-to-r from-indigo-600/90 via-indigo-500/90 to-fuchsia-600/90">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-full bg-white/20 flex items-center justify-center shrink-0">
                <svg viewBox="0 0 24 24" class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="9" />
                  <path d="M12 7v5l3 2" />
                </svg>
              </div>
              <div>
                <h3 id="timeoutTitle" class="text-xl font-extrabold tracking-wide">หมดเวลา!</h3>
                <p id="timeoutDesc" class="text-white/90 text-sm">รอบนี้จบแล้ว มาดูผลกัน แล้วเริ่มรอบใหม่ได้เลย</p>
              </div>
            </div>
          </div>

          <div class="px-6 pt-5 pb-6 space-y-5 text-slate-100">
            <div class="grid grid-cols-2 gap-3">
              <SummaryCard label="คะแนนรวม" :value="finalScore" accent="indigo" />
              <SummaryCard label="ระดับที่ไปถึง" :value="'Lv. ' + finalLevel" accent="fuchsia" />
            </div>

            <div class="rounded-2xl border border-emerald-300/25 bg-emerald-400/10 px-4 py-3">
              <div class="flex items-center gap-2 text-emerald-100">
                <svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7S2 12 2 12Z" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
                <div class="text-sm font-semibold">เฉลยข้อสุดท้าย</div>
              </div>
              <div class="mt-2">
                <template v-if="revealedAnswer">
                  <span class="inline-flex items-center gap-2 text-emerald-100 font-bold text-lg">{{ revealedAnswer
                    }}</span>
                </template>
                <template v-else>
                  <div class="h-6 w-44 rounded-full bg-emerald-200/40 animate-pulse" />
                  <p class="text-xs text-emerald-200/90 mt-1">กำลังดึงเฉลย…</p>
                </template>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-200 mb-1"
                for="playerName">กรอกชื่อเพื่อบันทึกสถิติ</label>
              <input id="playerName" ref="nameInput" v-model="playerName" type="text" placeholder="เช่น น้องหมา"
                class="px-4 py-2.5 rounded-xl w-full bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60" />
              <p class="text-xs text-slate-300 mt-1">* คุณสามารถเริ่มรอบใหม่ได้ทันที หรือบันทึกคะแนนก่อน</p>

              <div class="mt-4 flex flex-col sm:flex-row gap-2">
                <button @click="saveScore" type="button"
                  class="inline-flex justify-center items-center px-5 py-3 rounded-xl font-semibold w-full sm:w-1/2 transition bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
                  :disabled="!playerName.trim() || isSaving" :aria-busy="isSaving ? 'true' : 'false'">
                  <svg v-if="!isSaving" viewBox="0 0 24 24" class="h-5 w-5 mr-2" fill="none" stroke="currentColor"
                    stroke-width="2">
                    <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
                    <path d="M17 21v-8H7v8" />
                    <path d="M7 3v5h8" />
                  </svg>
                  <span v-else
                    class="mr-2 inline-block animate-spin h-5 w-5 border-2 border-white/80 border-t-transparent rounded-full" />
                  {{ isSaving ? 'กำลังบันทึก...' : 'บันทึกคะแนน' }}
                </button>

                <button @click="restartGame" type="button"
                  class="inline-flex justify-center items-center px-5 py-3 rounded-xl font-semibold w-full sm:w-1/2 transition bg-white/10 text-indigo-100 border border-white/15 hover:bg-white/20">
                  <svg viewBox="0 0 24 24" class="h-5 w-5 mr-2" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12a9 9 0 1 1-6.219-8.56" />
                    <path d="M21 3v7h-7" />
                  </svg>
                  เริ่มรอบใหม่
                </button>
              </div>
            </div>
          </div>

          <div class="px-6 py-3 bg-white/5 border-t border-white/10 text-[12px] text-slate-300">
            เคล็ดลับ: ตอบถูกเร็วได้โบนัส และสตรีคต่อเนื่องเพิ่มคอมโบ!
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/* ===================== Imports ===================== */
import catwalk from '../assets/images/catwalk.png'
import catwalk2 from '../assets/images/catwalk2.png'
import api from '../services/api'
import {
  ref, onMounted, watch, computed, onBeforeUnmount, nextTick,
  defineComponent, h, type PropType
} from 'vue'
import { useRouter, useRoute } from 'vue-router' // ✅ เพิ่ม useRoute
import { useNetworkStore } from '../store/useNetworkStore'
import { waitApiReadyAndLoadInitial } from '../composables/useApiReadiness'

/* ===================== Config / State ===================== */
const net = useNetworkStore()
const router = useRouter()
const route = useRoute() // ✅ อ่าน query

// โหมดไม่จับเวลา จาก query: ?noTimer=1
const noTimer = computed(() => String(route.query.noTimer || '') === '1') // ✅

const GAME_NAME = 'DogPuzzle'
const TOP_LIMIT = 10
const ROUND_SECONDS = 60
const CHANGE_LIMIT = 5

// loading animation
const loading = ref(true)
const catwalkImages = [catwalk, catwalk2]
const catwalkIndex = ref(0)
let catwalkInterval: number | undefined

// game states
const quizId = ref('')
const quizToken = ref('')
const quizExp = ref(0)
const placeholder = computed(() => `พิมพ์คำตอบที่นี่…`)

const guess = ref('')
const result = ref<null | boolean>(null)
const hint1 = ref('')
const hint2 = ref('')
const score = ref(0)
const timer = ref(ROUND_SECONDS)
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

// recent guesses (รีเซ็ตเมื่อเริ่มคำใหม่)
type GuessItem = { word: string; correct: boolean; ts: number; heat?: number }
const recentGuesses = ref<GuessItem[]>([])

const answerInput = ref<HTMLInputElement | null>(null)
const nameInput = ref<HTMLInputElement | null>(null)

// Level: ทุก ๆ 10 คะแนน เพิ่มระดับ (คงไว้เพื่อปรับความยาก/เวลา)
const currentLevel = computed(() => Math.floor(score.value / 10) + 1)

// intervals
let intervalId: number | undefined
let visibilityPausedAt = 0

const timerPercent = computed(() =>
  Math.max(0, Math.min(100, (timer.value / ROUND_SECONDS) * 100))
)

const guessSanitized = computed(() => sanitize(guess.value))
const duplicateWarning = computed(() =>
  recentGuesses.value.some(
    g => g.word.toLowerCase() === guessSanitized.value.toLowerCase()
  )
)

// proximity/heat meter
const heat = ref(0) // 0..100
const heatLabel = computed(() => heatText(heat.value))
const heatDotClass = computed(() => {
  if (heat.value >= 80) return 'bg-rose-400'
  if (heat.value >= 50) return 'bg-amber-300'
  if (heat.value >= 25) return 'bg-sky-300'
  return 'bg-slate-400'
})

/* ===================== Tiny Functional Components ===================== */
const InfoCard = defineComponent({
  name: 'InfoCard',
  props: {
    label: { type: String, required: true },
    value: { type: [String, Number], required: true },
    accent: { type: String as PropType<'indigo' | 'sky' | 'fuchsia' | 'emerald'>, required: true }
  },
  setup(props) {
    const base = 'rounded-xl py-2.5 px-4 text-center border '
    const cls =
      props.accent === 'indigo' ? base + 'bg-indigo-400/10 border-indigo-300/20' :
        props.accent === 'sky' ? base + 'bg-sky-400/10 border-sky-300/20' :
          props.accent === 'fuchsia' ? base + 'bg-fuchsia-400/10 border-fuchsia-300/20' :
            base + 'bg-emerald-400/10 border-emerald-300/20'
    return () => h('div', { class: cls }, [
      h('div', { class: 'text-[11px] font-semibold text-slate-200/90 tracking-wide' }, props.label),
      h('div', { class: 'mt-0.5 text-2xl font-bold text-slate-100 tabular-nums' }, String(props.value)),
    ])
  }
})

const SummaryCard = defineComponent({
  name: 'SummaryCard',
  props: {
    label: { type: String, required: true },
    value: { type: [String, Number], required: true },
    accent: { type: String as PropType<'indigo' | 'fuchsia'>, required: true }
  },
  setup(props) {
    const cls = props.accent === 'indigo'
      ? 'rounded-2xl border px-4 py-3 border-indigo-300/25 bg-indigo-400/10'
      : 'rounded-2xl border px-4 py-3 border-fuchsia-300/25 bg-fuchsia-400/10'
    return () => h('div', { class: cls }, [
      h('div', { class: 'text-xs font-semibold text-slate-200' }, props.label),
      h('div', { class: 'mt-1 text-3xl font-extrabold text-slate-100 tabular-nums' }, String(props.value)),
    ])
  }
})

const Chip = defineComponent({
  name: 'Chip',
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true },
    color: { type: String as PropType<'fuchsia' | 'sky'>, required: true }
  },
  setup(props) {
    const cls = props.color === 'fuchsia'
      ? 'inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm border bg-fuchsia-500/10 text-fuchsia-100 border-fuchsia-300/20'
      : 'inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm border bg-sky-500/10 text-sky-100 border-sky-300/20'
    return () => h('span', { class: cls }, [
      h('strong', { class: 'font-semibold' }, `${props.label}:`),
      h('span', null, props.value),
    ])
  }
})

/* ===================== Utils ===================== */
function sanitize(s: string) {
  return s.replace(/\s+/g, ' ').trim()
}

// เพิ่มระดับเสียงให้ดังขึ้น และแยกค่า correct/incorrect
const soundGain = { ok: 0.06, no: 0.05 } // เดิม ~0.02
function playTick(success: boolean) {
  if (!soundOn.value) return
  try {
    const ctx = new (window.AudioContext || (window as any).webkitAudioContext)()
    const o = ctx.createOscillator()
    const g = ctx.createGain()
    o.type = success ? 'triangle' : 'sawtooth'
    o.frequency.value = success ? 760 : 220
    g.gain.value = success ? soundGain.ok : soundGain.no
    o.connect(g); g.connect(ctx.destination); o.start()
    setTimeout(() => { o.stop(); ctx.close() }, 140)
  } catch { /* ignore */ }
}

function toast(title: string, message: string, type: 'info' | 'success' | 'error' = 'info') {
  const id = Math.random().toString(36).slice(2)
  toasts.value.push({ id, title, message, type })
  setTimeout(() => { toasts.value = toasts.value.filter(t => t.id !== id) }, 3500)
}

function heatText(v: number) {
  if (v >= 100) return 'ตรงเป๊ะ!'
  if (v >= 90) return 'ใกล้มาก'
  if (v >= 70) return 'ใกล้'
  if (v >= 50) return 'ปานกลาง'
  if (v >= 30) return 'ห่าง'
  return 'ห่างมาก'
}

/* ===================== Network Helpers (timeout + retry) ===================== */
const DEFAULT_TIMEOUT = 10000
function withTimeout<T>(p: Promise<T>, ms = DEFAULT_TIMEOUT): Promise<T> {
  return new Promise((resolve, reject) => {
    const t = setTimeout(() => reject(new Error('timeout')), ms)
    p.then(v => { clearTimeout(t); resolve(v) })
      .catch(e => { clearTimeout(t); reject(e) })
  })
}
async function apiGet(path: string, params?: any, retry = 1) {
  try {
    return await withTimeout(api.get(path, { params }))
  } catch (e: any) {
    if (retry > 0) {
      await new Promise(r => setTimeout(r, 500))
      return apiGet(path, params, retry - 1)
    }
    throw e
  }
}
async function apiPost(path: string, body?: any, retry = 1) {
  try {
    return await withTimeout(api.post(path, body))
  } catch (e: any) {
    if (retry > 0) {
      await new Promise(r => setTimeout(r, 500))
      return apiPost(path, body, retry - 1)
    }
    throw e
  }
}

/* ===================== Core Game Flows ===================== */
async function fetchQuiz(isAuto = false) {
  if (changeCount.value >= maxChange && !isAuto) return
  try {
    expiredNotice.value = false
    revealedAnswer.value = ''

    const res = await apiGet('/api/quiz', { level: currentLevel.value })
    quizId.value = res.data.id
    quizToken.value = res.data.token
    quizExp.value = res.data.exp
    maxHints.value = typeof res.data.hintCount === 'number' ? res.data.hintCount : 2

    // reset รอบใหม่
    result.value = null
    guess.value = ''
    hint1.value = ''
    hint2.value = ''
    hintCount.value = 0
    timer.value = Math.max(ROUND_SECONDS - Math.max(0, currentLevel.value - 3) * 3, 30)
    heat.value = 0
    recentGuesses.value = [] // ✅ รีเซ็ตประวัติคำเดา (กันคำซ้ำ)
    if (!isAuto) changeCount.value++

    // ✅ จัดการตัวจับเวลา: ถ้า noTimer → ไม่ตั้ง interval ใด ๆ
    if (intervalId) clearInterval(intervalId as number)
    intervalId = undefined
    if (!noTimer.value) {
      intervalId = window.setInterval(() => {
        if (timer.value > 0) timer.value--
        if (timer.value === 0) onTimeout()
      }, 1000) as unknown as number
    }

    nextTick(() => answerInput.value?.focus())
    await requestHint(1) // auto show first hint
  } catch (e: any) {
    toast('ดึงคำถามล้มเหลว', e.message || 'network error', 'error')
  }
}

function onTimeout() {
  clearInterval(intervalId as number)
  result.value = false
  finalScore.value = score.value
  finalLevel.value = currentLevel.value
  showModal.value = true
  revealAnswer()
  nextTick(() => nameInput.value?.focus())
}

// ✅ จบเกมด้วยมือ (เฉพาะโหมดไม่จับเวลา)
function endNoTimer() {
  try {
    // ไม่เปลี่ยน UI/logic โมดัลเดิม ใช้เหมือนหมดเวลา
    if (intervalId) clearInterval(intervalId as number)
    finalScore.value = score.value
    finalLevel.value = currentLevel.value
    showModal.value = true
    revealAnswer()
    nextTick(() => nameInput.value?.focus())
  } catch (e: any) {
    toast('จบเกมไม่สำเร็จ', e?.message || 'unexpected error', 'error')
  }
}

// debounced submit guard
let lastSubmitAt = 0
function handleSubmit() {
  if (!guessSanitized.value) return
  const now = Date.now()
  if (now - lastSubmitAt < 600) return // rate-limit
  lastSubmitAt = now
  if (duplicateWarning.value) {
    toast('คำซ้ำ', 'คุณเดาคำนี้ไปแล้ว', 'info')
    return
  }
  checkAnswer()
}

async function checkAnswer() {
  if (!quizId.value) return
  try {
    const res = await apiPost('/api/quiz/check', {
      id: quizId.value,
      guess: guessSanitized.value,
      token: quizToken.value,
      exp: quizExp.value,
    })
    if ((res.data as any)?.reason === 'expired') {
      expiredNotice.value = true
      if (intervalId) clearInterval(intervalId as number)
      timer.value = 0
      result.value = false
      finalScore.value = score.value
      finalLevel.value = currentLevel.value
      showModal.value = true
      revealAnswer()
      nextTick(() => nameInput.value?.focus())
      return
    }
    const ok = !!res.data.correct
    result.value = ok

    // proximity heuristic update
    updateHeat(guessSanitized.value, ok)

    recentGuesses.value.unshift({ word: guessSanitized.value, correct: ok, ts: Date.now(), heat: heat.value })
    recentGuesses.value = recentGuesses.value.slice(0, 10)

    if (ok) {
      playTick(true)
      score.value += 1 // ✅ ให้คะแนนแบบเรียบง่าย
      setTimeout(() => fetchQuiz(true), 700)
    } else {
      playTick(false)
    }
    guess.value = ''
  } catch (e: any) {
    console.error(e)
    toast('ส่งคำตอบล้มเหลว', e.message || 'network error', 'error')
  }
}

// heuristic heat (client-side)
function updateHeat(word: string, correct: boolean) {
  let v = 0
  v += Math.min(30, Math.max(0, 30 - Math.abs(word.length - 5) * 5)) // closer to length ~5
  const lower = word.toLowerCase()
  if (hint1.value) v += 10
  if (hint2.value) v += 10
  if (/^[ก-ฮa-z]/i.test(lower)) v += 5
  if (/(\w)\1/.test(lower)) v += 2
  v += Math.max(0, 20 - recentGuesses.value.length * 2)
  if (correct) v = 100
  heat.value = Math.max(0, Math.min(100, v))
}

// hints
async function requestHint(nextIndex: 1 | 2) {
  if (!quizId.value || !quizToken.value || !quizExp.value) return
  try {
    const res = await apiPost('/api/quiz/hint', {
      id: quizId.value,
      token: quizToken.value,
      exp: quizExp.value,
      index: nextIndex,
    })
    if ((res.data as any)?.error === 'expired') {
      expiredNotice.value = true
      if (intervalId) clearInterval(intervalId as number)
      timer.value = 0
      result.value = false
      finalScore.value = score.value
      finalLevel.value = currentLevel.value
      showModal.value = true
      revealAnswer()
      nextTick(() => nameInput.value?.focus())
      return
    }
    const text = (res.data && res.data.hint) || ''
    if (nextIndex === 1) {
      hint1.value = text
      hintCount.value = 1
    } else if (nextIndex === 2) {
      hint2.value = text
      hintCount.value = 2
      timer.value = Math.max(timer.value - 0, 0)
    }
  } catch (e: any) {
    toast('ขอคำใบ้ล้มเหลว', e.message || 'network error', 'error')
  }
}

// ดึงเฉลยหลัง exp
async function revealAnswer() {
  if (!quizId.value || !quizToken.value || !quizExp.value) return
  const now = Math.floor(Date.now() / 1000)
  const delay = Math.max(0, (quizExp.value - now) * 1000 + 50)
  setTimeout(async () => {
    try {
      const res = await apiPost('/api/quiz/reveal', { id: quizId.value, token: quizToken.value, exp: quizExp.value })
      revealedAnswer.value = (res.data && res.data.answer) || ''
    } catch (e: any) {
      toast('ดึงเฉลยล้มเหลว', e.message || 'network error', 'error')
      revealedAnswer.value = ''
    }
  }, delay)
}

async function restartGame() {
  score.value = 0
  finalScore.value = 0
  finalLevel.value = 1
  changeCount.value = 0
  revealedAnswer.value = ''
  showModal.value = false
  recentGuesses.value = []
  await fetchQuiz()
  await loadScores()
}

/* ===================== Offline Score Queue ===================== */
const pendingScoresKey = 'pendingScores'
function pushPendingScore(name: string, score: number) {
  const raw = localStorage.getItem(pendingScoresKey)
  const arr = raw ? JSON.parse(raw) as any[] : []
  arr.push({ name, score, gamename: GAME_NAME, at: Date.now() })
  localStorage.setItem(pendingScoresKey, JSON.stringify(arr))
}
async function flushPendingScores() {
  const raw = localStorage.getItem(pendingScoresKey)
  if (!raw) return
  const arr = JSON.parse(raw) as any[]
  if (!arr.length) return
  try {
    for (const s of arr) {
      await apiPost('/api/scores', s)
    }
    localStorage.removeItem(pendingScoresKey)
    toast('อัปโหลดสถิติ', 'ส่งคะแนนที่ค้างไว้สำเร็จ', 'success')
  } catch {
    // keep for next try
  }
}

async function saveScore() {
  if (!playerName.value.trim() || isSaving.value) return
  isSaving.value = true
  const scoreToSave = finalScore.value || score.value
  try {
    await apiPost('/api/scores', { name: playerName.value.trim(), score: scoreToSave, gamename: GAME_NAME })
    await loadScores()

    // reset for new game
    score.value = 0
    finalScore.value = 0
    finalLevel.value = 1
    changeCount.value = 0
    revealedAnswer.value = ''
    showModal.value = false
    playerName.value = ''
    fetchQuiz()
  } catch {
    pushPendingScore(playerName.value.trim(), scoreToSave)
    toast('บันทึกแบบออฟไลน์', 'เครือข่ายล้มเหลว — เก็บคะแนนไว้แล้ว จะส่งให้อัตโนมัติเมื่อออนไลน์', 'info')
    showModal.value = false
  } finally {
    isSaving.value = false
  }
}

async function loadScores() {
  try {
    const res = await apiGet('/api/scores', { limit: TOP_LIMIT, gamename: GAME_NAME })
    savedScores.value = (res.data || []).slice(0, TOP_LIMIT)
  } catch (e: any) {
    toast('โหลดสถิติล้มเหลว', e.message || 'network error', 'error')
  }
}

const goBack = () => router.back()

// ปุ่มคำใบ้
function showHint() {
  if (!quizId.value || !quizToken.value || !quizExp.value) return
  if (hintCount.value >= maxHints.value) return
  const nextIndex = Math.min(hintCount.value + 1, 2) as 1 | 2
  requestHint(nextIndex)
}

/* ===================== Lifecycle ===================== */
onMounted(async () => {
  document.title = 'PETTEXT - Context Quest'
  catwalkInterval = window.setInterval(() => {
    catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length
  }, 200)

  const { healthOk, initialOk } = await waitApiReadyAndLoadInitial()
  if (healthOk && initialOk) {
    loading.value = false
    await flushPendingScores()
    fetchQuiz()
    loadScores()
  } else {
    const checkTimer = window.setInterval(async () => {
      const h = await waitApiReadyAndLoadInitial()
      if (h.healthOk && h.initialOk) {
        window.clearInterval(checkTimer)
        loading.value = false
        await flushPendingScores()
        fetchQuiz()
        loadScores()
      }
    }, 5000)
  }

  window.addEventListener('visibilitychange', onVisibilityChange)
})

onBeforeUnmount(() => {
  if (intervalId) clearInterval(intervalId as number)
  if (catwalkInterval) clearInterval(catwalkInterval as number)
  window.removeEventListener('visibilitychange', onVisibilityChange)
})

function onVisibilityChange() {
  if (document.hidden) {
    visibilityPausedAt = Date.now()
  } else if (visibilityPausedAt && intervalId) {
    const diff = Math.floor((Date.now() - visibilityPausedAt) / 1000)
    timer.value = Math.max(0, timer.value - diff)
    if (timer.value === 0) onTimeout()
    visibilityPausedAt = 0
  }
}

/* ===================== Watches & Toasts ===================== */
watch(result, (val, oldVal) => {
  if (val === true && oldVal !== true) {
    const el = document.createElement('div')
    el.className = 'confetti'
    document.body.appendChild(el)
    setTimeout(() => el.remove(), 800)
  }
})

const toasts = ref<{ id: string; title: string; message: string; type: 'info' | 'success' | 'error' }[]>([])
function toastClass(type: 'info' | 'success' | 'error') {
  const base = 'bg-white/10 border-white/15 text-slate-100 backdrop-blur'
  if (type === 'success') return base + ' border-emerald-300/30 bg-emerald-400/10'
  if (type === 'error') return base + ' border-rose-300/30 bg-rose-400/10'
  return base
}
</script>
