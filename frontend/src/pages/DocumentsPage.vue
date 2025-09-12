<template>
  <div class="flex flex-col items-center justify-start min-h-screen bg-transparent py-8 relative">
    <!-- Loading overlay -->
    <div v-if="loading" class="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-lg text-blue-100 font-semibold">กำลังโหลด...</span>
      </div>
    </div>

    <!-- กล่องเกม -->
    <section class="w-full max-w-xl bg-white/90 rounded-2xl shadow-lg p-6 space-y-5">
      <!-- สรุปสถานะ -->
      <header class="space-y-3">
        <h2 class="text-2xl font-extrabold text-center text-gray-800 tracking-wide">เก่งจริงก็ทายมาดิ!</h2>

        <div class="grid grid-cols-3 gap-3">
          <div class="bg-blue-50 rounded-xl py-2.5 px-4 text-center">
            <div class="text-xs font-semibold text-blue-600 tracking-wide">คะแนน</div>
            <div class="mt-0.5 text-2xl font-bold text-blue-700 tabular-nums">{{ score }}</div>
          </div>
          <div class="bg-gray-50 rounded-xl py-2.5 px-4 text-center">
            <div class="text-xs font-semibold text-gray-600 tracking-wide">เวลาคงเหลือ</div>
            <div class="mt-0.5 text-2xl font-bold text-gray-800 tabular-nums">{{ timer }} วินาที</div>
          </div>
          <div class="bg-violet-50 rounded-xl py-2.5 px-4 text-center">
            <div class="text-xs font-semibold text-violet-600 tracking-wide">ระดับ</div>
            <div class="mt-0.5 text-2xl font-bold text-violet-700 tabular-nums">Lv. {{ currentLevel }}</div>
          </div>
        </div>

        <!-- Progress bar เวลา -->
        <div class="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
          <div class="h-full bg-blue-500 transition-all duration-300" :style="{ width: timerPercent + '%' }"
            role="progressbar" :aria-valuenow="timer" aria-label="ตัวนับเวลา" />
        </div>
      </header>

      <!-- แบบฟอร์มตอบ -->
      <form class="w-full flex flex-col sm:flex-row items-stretch sm:items-center gap-2" @submit.prevent="handleSubmit">
        <input ref="answerInput" v-model="guess" type="text" placeholder="พิมพ์คำตอบที่นี่..."
          class="border border-gray-300 rounded-xl px-4 py-2.5 text-base focus:outline-none focus:ring-2 focus:ring-blue-400/70 focus:border-blue-400 w-full"
          :disabled="showModal" autocomplete="off" />
        <button type="submit"
          class="px-4 py-2.5 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto"
          :disabled="showModal || !guess.trim()">
          ส่งคำตอบ
        </button>
        <button type="button" @click="showHint"
          class="relative px-4 py-2.5 bg-yellow-500 text-white rounded-xl hover:bg-yellow-600 transition disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto"
          :disabled="showModal || hintCount >= maxHints">
          <span
            class="absolute -top-2 -right-3 bg-white text-yellow-700 rounded-full px-2 py-0.5 text-xs font-bold shadow">
            {{ Math.min(hintCount, maxHints) }}/{{ maxHints }}
          </span>
          คำใบ้
        </button>
        <button type="button" @click="() => fetchQuiz()"
          class="relative px-4 py-2.5 bg-green-600 text-white rounded-xl hover:bg-green-700 transition disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto"
          :disabled="showModal || changeCount >= maxChange">
          <span
            class="absolute -top-2 -right-3 bg-white text-green-700 rounded-full px-2 py-0.5 text-xs font-bold shadow">
            {{ changeCount }}/{{ maxChange }}
          </span>
          เปลี่ยนคำ
        </button>
      </form>

      <!-- ผลลัพธ์ -->
      <div v-if="result !== null" class="text-center" aria-live="polite">
        <span v-if="result === true" class="text-green-600 text-xl font-semibold">✔️ ถูกต้อง!</span>
        <span v-else-if="result === false" class="text-red-600 text-xl font-semibold">❌ ผิด ลองใหม่!</span>
      </div>

      <!-- แจ้งเตือนหมดอายุ token -->
      <div v-if="expiredNotice"
        class="text-center text-sm text-orange-700 bg-orange-50 border border-orange-200 rounded-lg px-3 py-2">
        หมดเวลาตอบสำหรับคำนี้แล้ว ระบบจะสุ่มคำใหม่ให้อัตโนมัติ
      </div>

      <!-- คำใบ้ -->
      <div class="flex flex-wrap items-center justify-center gap-2">
        <span v-if="hint1"
          class="inline-flex items-center gap-2 bg-purple-50 text-purple-700 border border-purple-200 px-3 py-1.5 rounded-full text-sm">
          <strong class="font-semibold">คำใบ้ 1:</strong> <span>{{ hint1 }}</span>
        </span>
        <span v-if="hint2"
          class="inline-flex items-center gap-2 bg-purple-50 text-purple-700 border border-purple-200 px-3 py-1.5 rounded-full text-sm">
          <strong class="font-semibold">คำใบ้ 2:</strong> <span>{{ hint2 }}</span>
        </span>
      </div>

      <!-- แถบปุ่มล่าง -->
      <div class="pt-2">
        <button @click="goBack"
          class="w-full sm:w-auto px-4 py-2.5 bg-gray-600 text-white rounded-xl hover:bg-gray-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="showModal" type="button">
          ย้อนกลับ
        </button>
      </div>
    </section>

    <!-- สถิติ TOP 10 -->
    <section class="w-full max-w-xl mx-auto mt-6 bg-white/90 rounded-2xl shadow p-5">
      <div class="flex items-center justify-between">
        <h4 class="text-lg font-bold text-blue-700">สถิติผู้เล่น TOP 10</h4>
        <button @click="loadScores" class="text-xs px-3 py-1 border rounded-lg hover:bg-gray-50 transition"
          title="รีเฟรชสถิติ" type="button">
          รีเฟรช
        </button>
      </div>

      <p v-if="savedScores.length === 0" class="text-gray-500 text-sm mt-2">
        ยังไม่มีสถิติ แสดงเมื่อมีการบันทึกคะแนน
      </p>

      <ul v-else class="mt-2 divide-y">
        <li v-for="(item, idx) in savedScores" :key="item.name + '_' + item.score + '_' + idx"
          class="py-2 flex items-center justify-between text-sm">
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

    <!-- โมดัล: หมดเวลา + บันทึกคะแนน + เฉลย + เลเวล (ดีไซน์ใหม่) -->
    <div v-if="showModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 px-4"
      role="dialog" aria-modal="true" aria-labelledby="timeoutTitle" aria-describedby="timeoutDesc">
      <div class="w-full max-w-xl">
        <div class="relative bg-white rounded-3xl shadow-2xl ring-1 ring-gray-100 overflow-hidden">
          <!-- Top banner -->
          <div class="bg-[#111827] px-6 py-5 text-white">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-full bg-white/20 flex items-center justify-center shrink-0">
                <!-- clock icon -->
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
          <div class="px-6 pt-5 pb-6 space-y-5">
            <!-- Summary cards -->
            <div class="grid grid-cols-2 gap-3">
              <div class="rounded-2xl border border-blue-100 bg-blue-50/60 px-4 py-3">
                <div class="text-xs font-semibold text-blue-600">คะแนนรวม</div>
                <div class="mt-1 text-3xl font-extrabold text-blue-700 tabular-nums">{{ finalScore }}</div>
              </div>
              <div class="rounded-2xl border border-violet-100 bg-violet-50/60 px-4 py-3">
                <div class="text-xs font-semibold text-violet-600">ระดับที่ไปถึง</div>
                <div class="mt-1 text-3xl font-extrabold text-violet-700">Lv. {{ finalLevel }}</div>
              </div>
            </div>

            <!-- Answer reveal -->
            <div class="rounded-2xl border border-emerald-100 bg-emerald-50/60 px-4 py-3">
              <div class="flex items-center gap-2">
                <svg viewBox="0 0 24 24" class="h-5 w-5 text-emerald-700" fill="none" stroke="currentColor"
                  stroke-width="2">
                  <path d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7S2 12 2 12Z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
                <div class="text-sm font-semibold text-emerald-800">เฉลยข้อสุดท้าย</div>
              </div>
              <div class="mt-2">
                <template v-if="revealedAnswer">
                  <span class="inline-flex items-center gap-2 text-emerald-800 font-bold text-lg">
                    {{ revealedAnswer }}
                    <span
                      class="inline-block rounded-full bg-emerald-600 text-white text-xs px-2 py-0.5">ยืนยันแล้ว</span>
                  </span>
                </template>
                <template v-else>
                  <div class="h-6 w-44 rounded-full bg-emerald-200 animate-pulse"></div>
                  <p class="text-xs text-emerald-700 mt-1">กำลังดึงเฉลย…</p>
                </template>
              </div>
            </div>

            <!-- Name + actions -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1"
                for="playerName">กรอกชื่อเพื่อบันทึกสถิติ</label>
              <input id="playerName" ref="nameInput" v-model="playerName" type="text" placeholder="เช่น น้องแมว"
                class="px-4 py-2.5 border border-gray-300 rounded-xl w-full focus:outline-none focus:ring-2 focus:ring-blue-400/70" />
              <p class="text-xs text-gray-500 mt-1">* คุณสามารถเริ่มรอบใหม่ได้ทันที หรือบันทึกคะแนนก่อน</p>

              <div class="mt-4 flex flex-col sm:flex-row gap-2">
                <button @click="saveScore"
                  class="inline-flex justify-center items-center px-5 py-3 rounded-xl bg-blue-600 text-white font-semibold hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-1/2"
                  :disabled="!playerName.trim() || isSaving" :aria-busy="isSaving ? 'true' : 'false'" type="button">
                  <svg v-if="!isSaving" viewBox="0 0 24 24" class="h-5 w-5 mr-2" fill="none" stroke="currentColor"
                    stroke-width="2">
                    <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
                    <path d="M17 21v-8H7v8"></path>
                    <path d="M7 3v5h8"></path>
                  </svg>
                  <span v-else
                    class="mr-2 inline-block animate-spin h-5 w-5 border-2 border-white border-t-transparent rounded-full"></span>
                  {{ isSaving ? 'กำลังบันทึก...' : 'บันทึกคะแนน' }}
                </button>


                <button @click="restartGame"
                  class="inline-flex justify-center items-center px-5 py-3 rounded-xl bg-white text-blue-700 font-semibold border border-blue-200 hover:bg-blue-50 transition w-full sm:w-1/2"
                  type="button">
                  <svg viewBox="0 0 24 24" class="h-5 w-5 mr-2" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
                    <path d="M21 3v7h-7"></path>
                  </svg>
                  เริ่มรอบใหม่
                </button>
              </div>
            </div>
          </div>

          <!-- subtle bottom bar -->
          <div class="px-6 py-3 bg-gray-50/60 border-t border-gray-100 text-xs text-gray-500">
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

// loading animation
const loading = ref(true)
const catwalkImages = [catwalk, catwalk2]
const catwalkIndex = ref(0)
let catwalkInterval: number | undefined

// game states
const quizId = ref('')
const quizToken = ref('') // token จาก API
const quizExp = ref(0)    // exp (unix seconds) จาก API

const guess = ref('')
const result = ref<null | boolean>(null)
const hint1 = ref('')
const hint2 = ref('')
const score = ref(0)
const timer = ref(60)
const hintCount = ref(0)     // ขอไปแล้วกี่ใบ้
const maxHints = ref(2)      // รับมาจาก backend ใน /api/quiz // จำนวนคำใบ้สูงสุดต่อข้อ
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

// Level: ทุก ๆ 10 คะแนน เพิ่มระดับ (Lv.1–3 ใช้ชุดง่าย, Lv.4+ ใช้ชุดยาก)
const currentLevel = computed(() => Math.floor(score.value / 10) + 1)

// intervals
let intervalId: number | undefined

const timerPercent = computed(() => Math.max(0, Math.min(100, (timer.value / 60) * 100)))

// ขอควิซใหม่ (ไม่มีคำใบ้มากับ response)
async function fetchQuiz(isAuto = false) {
  if (changeCount.value >= maxChange && !isAuto) return

  expiredNotice.value = false
  revealedAnswer.value = '' // เคลียร์เฉลยของข้อก่อนหน้า

  // ส่ง level ปัจจุบันไปขอชุดคำตามระดับ
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
      revealAnswer() // ดึงเฉลยของข้อนี้ (จะรอจน exp ถึงก่อน)
      nextTick(() => nameInput.value?.focus())
    }
  }, 1000)

  if (!isAuto) changeCount.value++

  // โฟกัสช่องตอบเมื่อพร้อม
  nextTick(() => answerInput.value?.focus())
  //โชว์ "คำใบ้แรก" อัตโนมัติทันทีที่เริ่มรอบ
  await requestHint(1)
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
      // กติกา: ขอใบ้ที่สอง หักเวลา 10 วิ
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

// ดึงเฉลย (อนุญาตหลัง exp เท่านั้น)
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
    await api.post('/api/scores', { name: playerName.value.trim(), score: scoreToSave })
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
  const res = await api.get('/api/scores')
  savedScores.value = (res.data || []).sort((a: any, b: any) => b.score - a.score).slice(0, 10)
}

const goBack = () => router.back()

onMounted(() => {
  document.title = 'PETTEXT - Documents'
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
