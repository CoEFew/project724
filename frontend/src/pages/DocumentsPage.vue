<template>
  <div class="min-h-screen relative overflow-x-hidden theme-modern">
    <!-- Gradient background (โมเดิร์น, สบายตา) -->
    <div
      class="pointer-events-none absolute inset-0 -z-10"
      aria-hidden="true"
    >
      <!-- ชั้นที่ 1: ไล่เฉดตาม reference -->
      <div
        class="absolute inset-0 bg-[radial-gradient(90%_70%_at_70%_100%,rgba(99,102,241,0.45),transparent_60%),radial-gradient(60%_60%_at_0%_0%,rgba(59,130,246,0.35),transparent_60%),linear-gradient(180deg,#0b1020,#0b1120)]"
      />
      <!-- ชั้นที่ 2: วงเรืองรองนุ่ม ๆ -->
      <div class="absolute -bottom-16 right-10 h-80 w-80 rounded-full blur-3xl opacity-40 bg-indigo-500/30" />
      <div class="absolute -top-12 left-[-4rem] h-72 w-72 rounded-full blur-3xl opacity-30 bg-fuchsia-500/25" />
    </div>

    <!-- Loading overlay -->
    <div v-if="loading" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[90]">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-base md:text-lg text-indigo-100 font-semibold">กำลังโหลด...</span>
      </div>
    </div>

    <!-- Page container -->
    <div class="w-full max-w-6xl mx-auto px-4 py-8">
      <!-- Header แบบเดียวกับ Jigsaw + glass -->
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

          <h1 class="text-2xl md:text-4xl font-extrabold tracking-wide text-indigo-300/90 uppercase text-center flex-1 drop-shadow-sm">
            Dog•Puzzle
          </h1>

          <div class="w-[90px] sm:w-[120px]" />
        </div>
        <p class="text-slate-300/80 text-xs md:text-sm text-center">
          บ๊อก • แบ๊ก
        </p>
      </header>

      <!-- กล่องเกม (Glass card) -->
      <section class="w-full max-w-xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-6 space-y-5">
        <!-- สรุปสถานะ -->
        <div class="space-y-3">
          <h2 class="text-xl md:text-2xl font-extrabold text-indigo-100 tracking-wide text-center">เก่งจริงก็ทายมาดิ!</h2>

          <div class="grid grid-cols-3 gap-3">
            <div class="rounded-xl py-2.5 px-4 text-center bg-indigo-400/10 border border-indigo-300/20">
              <div class="text-[11px] font-semibold text-indigo-200/90 tracking-wide">คะแนน</div>
              <div class="mt-0.5 text-2xl font-bold text-indigo-100 tabular-nums">{{ score }}</div>
            </div>
            <div class="rounded-xl py-2.5 px-4 text-center bg-sky-400/10 border border-sky-300/20">
              <div class="text-[11px] font-semibold text-sky-200/90 tracking-wide">เวลาคงเหลือ</div>
              <div class="mt-0.5 text-2xl font-bold text-slate-100 tabular-nums">{{ timer }} วินาที</div>
            </div>
            <div class="rounded-xl py-2.5 px-4 text-center bg-fuchsia-400/10 border border-fuchsia-300/20">
              <div class="text-[11px] font-semibold text-fuchsia-200/90 tracking-wide">ระดับ</div>
              <div class="mt-0.5 text-2xl font-bold text-fuchsia-100 tabular-nums">Lv. {{ currentLevel }}</div>
            </div>
          </div>

          <!-- Progress bar -->
          <div class="w-full h-2 rounded-full overflow-hidden bg-white/10">
            <div
              class="h-full transition-all duration-300 bg-gradient-to-r from-sky-400 via-indigo-400 to-fuchsia-400"
              :style="{ width: timerPercent + '%' }"
              role="progressbar"
              :aria-valuenow="timer"
              aria-label="ตัวนับเวลา"
            />
          </div>
        </div>

        <!-- แบบฟอร์มตอบ -->
        <form class="w-full flex flex-col sm:flex-row items-stretch sm:items-center gap-2" @submit.prevent="handleSubmit">
          <input
            ref="answerInput"
            v-model="guess"
            type="text"
            placeholder="พิมพ์คำตอบที่นี่..."
            class="rounded-xl px-4 py-2.5 text-base bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60 focus:border-indigo-400/60 w-full"
            :disabled="showModal"
            autocomplete="off"
          />
          <button
            type="submit"
            class="px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition
                   bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
            :disabled="showModal || !guess.trim()"
          >
            ส่งคำตอบ
          </button>
          <button
            type="button"
            @click="showHint"
            class="relative px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition
                   bg-amber-500 text-slate-900 hover:bg-amber-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
            :disabled="showModal || hintCount >= maxHints"
          >
            <span class="absolute -top-2 -right-3 bg-white/90 text-amber-800 rounded-full px-2 py-0.5 text-xs font-bold shadow">
              {{ Math.min(hintCount, maxHints) }}/{{ maxHints }}
            </span>
            คำใบ้
          </button>
          <button
            type="button"
            @click="() => fetchQuiz()"
            class="relative px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition
                   bg-emerald-500 text-white hover:bg-emerald-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
            :disabled="showModal || changeCount >= maxChange"
          >
            <span class="absolute -top-2 -right-3 bg-white/90 text-emerald-800 rounded-full px-2 py-0.5 text-xs font-bold shadow">
              {{ changeCount }}/{{ maxChange }}
            </span>
            เปลี่ยนคำ
          </button>
        </form>

        <!-- ผลลัพธ์ -->
        <div v-if="result !== null" class="text-center" aria-live="polite">
          <span v-if="result === true" class="text-emerald-300 text-xl font-semibold">✔️ ถูกต้อง!</span>
          <span v-else-if="result === false" class="text-rose-300 text-xl font-semibold">❌ ผิด ลองใหม่!</span>
        </div>

        <!-- แจ้งเตือนหมดอายุ token -->
        <div v-if="expiredNotice" class="text-center text-sm text-amber-200 bg-amber-500/10 border border-amber-400/30 rounded-lg px-3 py-2">
          หมดเวลาตอบสำหรับคำนี้แล้ว ระบบจะสุ่มคำใหม่ให้อัตโนมัติ
        </div>

        <!-- คำใบ้ -->
        <div class="flex flex-wrap items-center justify-center gap-2">
          <span v-if="hint1" class="inline-flex items-center gap-2 bg-fuchsia-500/10 text-fuchsia-100 border border-fuchsia-300/20 px-3 py-1.5 rounded-full text-sm">
            <strong class="font-semibold">คำใบ้ 1:</strong> <span>{{ hint1 }}</span>
          </span>
          <span v-if="hint2" class="inline-flex items-center gap-2 bg-fuchsia-500/10 text-fuchsia-100 border border-fuchsia-300/20 px-3 py-1.5 rounded-full text-sm">
            <strong class="font-semibold">คำใบ้ 2:</strong> <span>{{ hint2 }}</span>
          </span>
        </div>
      </section>

      <!-- สถิติ TOP 10 -->
      <section class="w-full max-w-xl mx-auto mt-6 rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-5">
        <div class="flex items-center justify-between">
          <h4 class="text-lg font-bold text-indigo-100">สถิติผู้เล่น TOP 10</h4>
          <button @click="loadScores" class="text-xs px-3 py-1 rounded-lg border border-white/10 bg-white/5 text-slate-200 hover:bg-white/10 transition" title="รีเฟรชสถิติ" type="button">
            รีเฟรช
          </button>
        </div>

        <p v-if="savedScores.length === 0" class="text-slate-300/70 text-sm mt-2">
          ยังไม่มีสถิติ แสดงเมื่อมีการบันทึกคะแนน
        </p>

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
              <span class="font-medium truncate max-w-[10rem] sm:max-w-[14rem]" :title="item.name">
                {{ item.name }}
              </span>
            </div>
            <div class="tabular-nums font-semibold">{{ item.score }} คะแนน</div>
          </li>
        </ul>
      </section>
    </div>

    <!-- โมดัล: หมดเวลา + บันทึกคะแนน (z สูงกว่า overlay) -->
    <div v-if="showModal"
         class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[95] px-4"
         role="dialog" aria-modal="true" aria-labelledby="timeoutTitle" aria-describedby="timeoutDesc">
      <div class="w-full max-w-xl">
        <div class="relative rounded-3xl overflow-hidden shadow-[0_30px_80px_rgba(0,0,0,0.55)] border border-white/10 bg-white/5 backdrop-blur-xl">
          <!-- Top banner -->
          <div class="px-6 py-5 text-white bg-gradient-to-r from-indigo-600/90 via-indigo-500/90 to-fuchsia-600/90">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-full bg-white/20 flex items-center justify-center shrink-0">
                <svg viewBox="0 0 24 24" class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="9"></circle>
                  <path d="M12 7v5l3 2"></path>
                </svg>
              </div>
              <div>
                <h3 id="timeoutTitle" class="text-xl font-extrabold tracking-wide">หมดเวลา!</h3>
                <p id="timeoutDesc" class="text-white/90 text-sm">รอบนี้จบแล้ว มาดูผลกัน แล้วเริ่มรอบใหม่ได้เลย</p>
              </div>
            </div>
          </div>

          <!-- Body -->
          <div class="px-6 pt-5 pb-6 space-y-5 text-slate-100">
            <!-- Summary cards -->
            <div class="grid grid-cols-2 gap-3">
              <div class="rounded-2xl border border-indigo-300/25 bg-indigo-400/10 px-4 py-3">
                <div class="text-xs font-semibold text-indigo-200">คะแนนรวม</div>
                <div class="mt-1 text-3xl font-extrabold text-indigo-100 tabular-nums">{{ finalScore }}</div>
              </div>
              <div class="rounded-2xl border border-fuchsia-300/25 bg-fuchsia-400/10 px-4 py-3">
                <div class="text-xs font-semibold text-fuchsia-200">ระดับที่ไปถึง</div>
                <div class="mt-1 text-3xl font-extrabold text-fuchsia-100">Lv. {{ finalLevel }}</div>
              </div>
            </div>

            <!-- Answer reveal -->
            <div class="rounded-2xl border border-emerald-300/25 bg-emerald-400/10 px-4 py-3">
              <div class="flex items-center gap-2 text-emerald-100">
                <svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7S2 12 2 12Z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
                <div class="text-sm font-semibold">เฉลยข้อสุดท้าย</div>
              </div>
              <div class="mt-2">
                <template v-if="revealedAnswer">
                  <span class="inline-flex items-center gap-2 text-emerald-100 font-bold text-lg">
                    {{ revealedAnswer }}
                  </span>
                </template>
                <template v-else>
                  <div class="h-6 w-44 rounded-full bg-emerald-200/40 animate-pulse"></div>
                  <p class="text-xs text-emerald-200/90 mt-1">กำลังดึงเฉลย…</p>
                </template>
              </div>
            </div>

            <!-- Name + actions -->
            <div>
              <label class="block text-sm font-medium text-slate-200 mb-1" for="playerName">กรอกชื่อเพื่อบันทึกสถิติ</label>
              <input
                id="playerName"
                ref="nameInput"
                v-model="playerName"
                type="text"
                placeholder="เช่น น้องหมา"
                class="px-4 py-2.5 rounded-xl w-full bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60"
              />
              <p class="text-xs text-slate-300 mt-1">* คุณสามารถเริ่มรอบใหม่ได้ทันที หรือบันทึกคะแนนก่อน</p>

              <div class="mt-4 flex flex-col sm:flex-row gap-2">
                <button
                  @click="saveScore"
                  class="inline-flex justify-center items-center px-5 py-3 rounded-xl font-semibold w-full sm:w-1/2 transition
                         bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
                  :disabled="!playerName.trim() || isSaving"
                  :aria-busy="isSaving ? 'true' : 'false'"
                  type="button"
                >
                  <svg v-if="!isSaving" viewBox="0 0 24 24" class="h-5 w-5 mr-2" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
                    <path d="M17 21v-8H7v8"></path>
                    <path d="M7 3v5h8"></path>
                  </svg>
                  <span v-else class="mr-2 inline-block animate-spin h-5 w-5 border-2 border-white/80 border-t-transparent rounded-full"></span>
                  {{ isSaving ? 'กำลังบันทึก...' : 'บันทึกคะแนน' }}
                </button>

                <button
                  @click="restartGame"
                  class="inline-flex justify-center items-center px-5 py-3 rounded-xl font-semibold w-full sm:w-1/2 transition
                         bg-white/10 text-indigo-100 border border-white/15 hover:bg-white/20"
                  type="button"
                >
                  <svg viewBox="0 0 24 24" class="h-5 w-5 mr-2" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
                    <path d="M21 3v7h-7"></path>
                  </svg>
                  เริ่มรอบใหม่
                </button>
              </div>
            </div>
          </div>

          <div class="px-6 py-3 bg-white/5 border-t border-white/10 text-[12px] text-slate-300">
            เคล็ดลับ: ทุก ๆ 10 คะแนน เลเวลจะเพิ่มขึ้น (ตั้งแต่เลเวล 4 เป็นต้นไปจะเป็นชุดคำยาก)
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import catwalk from '../assets/images/catwalk.png'
import catwalk2 from '../assets/images/catwalk2.png'
import api from '../services/api'
import { ref, onMounted, watch, computed, onBeforeUnmount, nextTick } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

/** ------- Config ------- */
const GAME_NAME = 'DogPuzzle'
const TOP_LIMIT = 10

// loading animation
const loading = ref(true)
const catwalkImages = [catwalk, catwalk2]
const catwalkIndex = ref(0)
let catwalkInterval: number | undefined

// game states
const quizId = ref('')
const quizToken = ref('')
const quizExp = ref(0)

const guess = ref('')
const result = ref<null | boolean>(null)
const hint1 = ref('')
const hint2 = ref('')
const score = ref(0)
const timer = ref(60)
const hintCount = ref(0)
const maxHints = ref(2) // รับจาก backend /api/quiz
const showModal = ref(false)
const changeCount = ref(0)
const maxChange = 5
const playerName = ref('')
const savedScores = ref<{ name: string; score: number }[]>([])
const finalScore = ref(0)
const finalLevel = ref(1)
const expiredNotice = ref(false)
const revealedAnswer = ref('')
const isSaving = ref(false)

const answerInput = ref<HTMLInputElement | null>(null)
const nameInput = ref<HTMLInputElement | null>(null)

// Level: ทุก ๆ 10 คะแนน เพิ่มระดับ
const currentLevel = computed(() => Math.floor(score.value / 10) + 1)

// intervals
let intervalId: number | undefined

const timerPercent = computed(() => Math.max(0, Math.min(100, (timer.value / 60) * 100)))

// ขอควิซใหม่
async function fetchQuiz(isAuto = false) {
  if (changeCount.value >= maxChange && !isAuto) return

  expiredNotice.value = false
  revealedAnswer.value = ''

  const res = await api.get('/api/quiz', { params: { level: currentLevel.value } })
  quizId.value = res.data.id
  quizToken.value = res.data.token
  quizExp.value = res.data.exp
  maxHints.value = typeof res.data.hintCount === 'number' ? res.data.hintCount : 2

  result.value = null
  guess.value = ''
  hint1.value = ''
  hint2.value = ''
  hintCount.value = 0
  timer.value = 60
  showModal.value = false

  if (intervalId) clearInterval(intervalId as number)
  intervalId = window.setInterval(() => {
    if (timer.value > 0) timer.value--
    if (timer.value === 0) {
      clearInterval(intervalId as number)
      result.value = false
      finalScore.value = score.value
      finalLevel.value = currentLevel.value
      showModal.value = true
      revealAnswer()
      nextTick(() => nameInput.value?.focus())
    }
  }, 1000)

  if (!isAuto) changeCount.value++
  nextTick(() => answerInput.value?.focus())
  await requestHint(1) // เปิดรอบโชว์ใบ้แรกอัตโนมัติ
}

function handleSubmit() {
  if (!guess.value.trim()) return
  checkAnswer()
}

async function checkAnswer() {
  if (!quizId.value) return
  try {
    const res = await api.post('/api/quiz/check', {
      id: quizId.value,
      guess: guess.value.trim(),
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
    result.value = res.data.correct
  } catch (e) {
    console.error(e)
  }
}

watch(result, (val, oldVal) => {
  if (val === true && oldVal !== true) {
    score.value++
    setTimeout(() => fetchQuiz(true), 1000)
  }
})

// ขอคำใบ้จาก server ทีละอัน
async function requestHint(nextIndex: 1 | 2) {
  if (!quizId.value || !quizToken.value || !quizExp.value) return
  try {
    const res = await api.post('/api/quiz/hint', {
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
      timer.value = Math.max(timer.value - 10, 0)
    }
  } catch (e) {
    console.error(e)
  }
}

function showHint() {
  if (hintCount.value >= maxHints.value) return
  const nextIndex = (hintCount.value + 1) as 1 | 2
  requestHint(nextIndex)
}

// ดึงเฉลยหลัง exp
async function revealAnswer() {
  if (!quizId.value || !quizToken.value || !quizExp.value) return
  const now = Math.floor(Date.now() / 1000)
  const delay = Math.max(0, (quizExp.value - now) * 1000 + 50)

  setTimeout(async () => {
    try {
      const res = await api.post('/api/quiz/reveal', {
        id: quizId.value,
        token: quizToken.value,
        exp: quizExp.value,
      })
      revealedAnswer.value = (res.data && res.data.answer) || ''
    } catch (e) {
      console.error(e)
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
  await fetchQuiz()
  await loadScores()
}

async function saveScore() {
  if (!playerName.value.trim() || isSaving.value) return

  isSaving.value = true
  try {
    const scoreToSave = finalScore.value || score.value
    await api.post('/api/scores', {
      name: playerName.value.trim(),
      score: scoreToSave,
      gamename: GAME_NAME,
    })
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
  } catch (e) {
    console.error(e)
  } finally {
    isSaving.value = false
  }
}

async function loadScores() {
  const res = await api.get('/api/scores', { params: { limit: TOP_LIMIT, gamename: GAME_NAME } })
  savedScores.value = (res.data || []).slice(0, TOP_LIMIT)
}

const goBack = () => router.back()

onMounted(() => {
  document.title = 'PETTEXT - DogPuzzle'
  catwalkInterval = setInterval(() => {
    catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length
  }, 200)
  setTimeout(() => (loading.value = false), 800)

  fetchQuiz()
  loadScores()
})

onBeforeUnmount(() => {
  if (intervalId) clearInterval(intervalId as number)
  if (catwalkInterval) clearInterval(catwalkInterval as number)
})
</script>

<style scoped>
.tabular-nums { font-variant-numeric: tabular-nums; }

/* โหมดธีม (จูน contrast และ smoothing) */
.theme-modern { color-scheme: dark; }

/* ลดจังหวะ animation ให้สบายตา */
@keyframes pulseSoft {
  0%, 100% { transform: scale(1); filter: saturate(1); }
  50% { transform: scale(1.02); filter: saturate(1.05); }
}
.animate-pulse-soft { animation: pulseSoft 1.2s ease-in-out infinite; }
</style>
