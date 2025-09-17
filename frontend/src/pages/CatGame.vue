<template>
  <div class="min-h-screen relative overflow-x-hidden theme-modern">
    <!-- Background -->
    <div class="pointer-events-none absolute inset-0 -z-10" aria-hidden="true">
      <div
        class="absolute inset-0 bg-[radial-gradient(90%_70%_at_70%_100%,rgba(99,102,241,0.45),transparent_60%),radial-gradient(60%_60%_at_0%_0%,rgba(59,130,246,0.35),transparent_60%),linear-gradient(180deg,#0b1020,#0b1120)]">
      </div>
      <div class="absolute -bottom-16 right-10 h-80 w-80 rounded-full blur-3xl opacity-40 bg-indigo-500/30"></div>
      <div class="absolute -top-12 left-[-4rem] h-72 w-72 rounded-full blur-3xl opacity-30 bg-fuchsia-500/25"></div>
    </div>

    <!-- Loading overlay -->
    <div v-if="loading"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex flex-col items-center justify-center z-[120]">
      <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-2 animate-bounce" />
      <span class="text-base md:text-lg text-indigo-100 font-semibold">กำลังโหลด...</span>
      <span class="mt-1 text-xs text-indigo-100/70" v-if="net.hasPending">กำลังเชื่อมต่อเซิร์ฟเวอร์…</span>
      <span class="mt-1 text-xs text-amber-200/80" v-if="net.isStalled">เซิร์ฟเวอร์กำลังเริ่มทำงาน
        ช้ากว่าปกติเล็กน้อย</span>
      <span class="mt-1 text-xs text-rose-200/80" v-if="net.lastError">พบข้อผิดพลาดเครือข่าย: {{ net.lastError }}</span>
    </div>


    <!-- Rotate-warning overlay for small screens (force landscape UX) -->
    <div v-if="isPortrait"
      class="fixed inset-0 z-[100] bg-black/80 backdrop-blur-sm grid place-items-center text-center p-6">
      <div class="max-w-sm text-slate-100">
        <div class="text-5xl mb-3">📱↻</div>
        <h2 class="text-xl font-bold mb-2">กรุณาหมุนอุปกรณ์เป็นแนวนอน</h2>
        <p class="text-slate-300 text-sm">เพื่อประสบการณ์ที่ดีที่สุดในการเล่นไพ่ 8 ใบ
          ระบบจะแสดงผลเมื่ออุปกรณ์อยู่ในแนวนอน</p>
      </div>
    </div>

    <!-- Container -->
    <div class="w-full max-w-7xl mx-auto px-4 pt-8 pb-44"><!-- bottom padding for sticky controls -->
      <!-- Header -->
      <header class="flex flex-col gap-3 items-center mb-6">
        <div class="w-full flex items-center justify-between">
          <!-- ปุ่มย้อนกลับ -->
          <button @click="goBack" type="button"
            class="inline-flex items-center gap-2 px-3 py-2 rounded-xl border bg-white/5 border-white/10 text-slate-100 hover:bg-white/10 transition shadow-sm">
            <svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 18l-6-6 6-6"></path>
            </svg>
            ย้อนกลับ
          </button>

          <h1
            class="text-2xl md:text-4xl font-extrabold tracking-wide text-indigo-300/90 uppercase text-center flex-1 drop-shadow-sm">
            Cat•Demon
          </h1>

          <div class="flex items-center gap-2">
            <button @click="hardReset()" type="button"
              class="inline-flex items-center gap-2 px-3 py-2 rounded-xl border bg-white/5 border-white/10 text-slate-100 hover:bg-white/10 transition shadow-sm">
              <svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M15 18l-6-6 6-6"></path>
              </svg>
              รีเซ็ตทั้งหมด
            </button>
            <button type="button"
              class="inline-flex items-center gap-2 px-3 py-2 rounded-xl border bg-white/5 border-white/10 text-slate-100 hover:bg-white/10 transition shadow-sm"
              @click="openDeckModal">
              🗂️ ดูการ์ด
            </button>
            <button type="button"
              class="inline-flex items-center gap-2 px-3 py-2 rounded-xl border bg-white/5 border-white/10 text-slate-100 hover:bg-white/10 transition shadow-sm"
              @click="openHowToModal">
              ❓ วิธีเล่น & คอมโบ
            </button>
          </div>
        </div>
      </header>

      <!-- Boss Arena (centered) -->
      <section
        class="rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-4 flex flex-col items-center justify-center">
        <div class="w-full flex items-center justify-between text-slate-200 text-sm mb-2">
          <div class="flex items-center gap-2">
            <span class="text-indigo-200 font-semibold">Boss:</span>
            <span class="font-bold">{{ boss.name }}</span>
          </div>
          <div class="text-xs text-slate-300 flex items-center gap-2">
            <span class="px-2 py-0.5 rounded-full bg-white/10 border border-white/10">โต้กลับในอีก <b
                class="tabular-nums">{{ turnsUntilCounter }}</b> เทิร์น</span>
            <span class="px-2 py-0.5 rounded-full bg-white/10 border border-white/10">โต้กลับแล้ว <b
                class="tabular-nums">{{ Math.floor(turnCount / bossAttackEvery) }}</b> ครั้ง</span>
            <span class="px-2 py-0.5 rounded-full bg-white/10 border border-white/10">ดาเมจโต้กลับครั้งถัดไป: <b
                class="tabular-nums">{{ nextCounterDamage ?? '—' }}</b></span>
          </div>
        </div>
        <div
          class="w-56 h-56 md:w-72 md:h-72 rounded-2xl border border-white/10 bg-gradient-to-b from-slate-800/60 to-slate-900/60 relative overflow-hidden grid place-items-center">
          <img v-if="bossImg" :src="bossImg" alt="boss" class="object-cover w-full h-full" />
          <div v-else class="text-6xl select-none">👹</div>

          <!-- Attack overlay shown ONLY during actual attack -->
          <div v-if="playerAttacking" class="attack-overlay">
            <div class="slash"></div>
            <div class="dmg-float">-{{ lastDamage }}</div>
          </div>
        </div>
        <div class="w-full mt-3">
          <div class="flex items-center justify-between text-slate-200 text-sm mb-1">
            <span>HP</span>
            <span class="tabular-nums">{{ boss.hp }} / {{ boss.maxHp }}</span>
          </div>
          <div class="hpbar">
            <div class="hpfill boss" :style="{ width: bossHpPct + '%' }"></div>
          </div>
        </div>
        <p v-if="victory" class="mt-2 text-emerald-300 text-sm font-semibold">🎉 คุณชนะแล้ว!</p>
        <p v-if="defeat" class="mt-2 text-rose-300 text-sm font-semibold">💀 พ่ายแพ้!</p>
      </section>

      <!-- Hand area (always single row on md+, wraps on small screens) -->
      <section class="mt-6">
        <div class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-8 gap-3">
          <div v-for="(card, idx) in hand" :key="idx + '_' + card.rank + card.suit" class="relative">
            <button class="card group" :class="{
      'ring-2 ring-rose-400/70 selected': selectedIndices.has(idx)
    }" @click="onCardClick(idx)" :disabled="turnEnded">
              <div class="card-head">
                <span class="rank">{{ card.rank }}</span>
                <span class="suit" :class="suitTextClass(card.suit)">{{ card.suit }}</span>
              </div>
              <div class="card-body">
                <span class="suit-big" :class="suitTextClass(card.suit)">{{ suitAnimal(card.suit) }}</span>
              </div>
            </button>
          </div>
        </div>
      </section>

      <!-- สถิติ TOP 10 -->
      <section
        class="w-full max-w-3xl mx-auto mt-6 rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-5">
        <div class="flex items-center justify-between">
          <h4 class="text-lg font-bold text-indigo-100">สถิติผู้เล่น TOP 10</h4>
          <button @click="loadScores"
            class="text-xs px-3 py-1 rounded-lg border border-white/10 bg-white/5 text-slate-200 hover:bg-white/10 transition"
            title="รีเฟรชสถิติ" type="button">
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

    <!-- Sticky bottom control bar (full 12-grid, centered) -->
    <footer class="fixed bottom-0 left-0 right-0 z-50">
      <div class="max-w-7xl mx-auto px-4 pb-4">
        <div
          class="rounded-t-2xl border border-white/10 bg-white/10 backdrop-blur-md p-3 shadow-[0_-10px_30px_rgba(0,0,0,0.25)]">
          <div class="grid grid-cols-12 gap-3 items-center">
            <!-- Player box (col-span responsive) -->
            <div class="col-span-12 md:col-span-5 lg:col-span-4">
              <div class="fighter-card relative" :class="playerBeingHit ? 'shake' : ''">
                <div class="flex items-center justify-between text-slate-200 text-sm mb-1">
                  <span>🧑 Player</span>
                  <span class="tabular-nums">{{ player.hp }} / {{ player.maxHp }}</span>
                </div>
                <div class="hpbar">
                  <div class="hpfill" :style="{ width: playerHpPct + '%' }"></div>
                </div>
                <div v-if="playerBeingHit" class="attack-overlay">
                  <div class="slash reverse"></div>
                  <div class="dmg-float text-rose-300">-{{ bossLastDamage }}</div>
                </div>
                <div class="mt-2 text-xs text-slate-300/80">
                  ดาเมจที่จะเกิดตอนนี้: <b class="tabular-nums text-emerald-300">{{ previewDamage }}</b> <span
                    class="text-slate-400">({{ previewLabel }})</span>
                </div>
                <!-- คะแนนสด -->
                <div class="mt-1 text-xs text-slate-300/80">
                  คะแนนสะสม: <b class="tabular-nums text-indigo-200">{{ score }}</b>
                </div>
              </div>
            </div>

            <!-- Controls -->
            <div
              class="col-span-12 md:col-span-7 lg:col-span-8 flex flex-wrap items-center justify-center md:justify-end gap-2">
              <div class="text-slate-200 text-xs md:text-sm mr-2">
                <span class="mr-3">ดึงได้อีก: <b class="tabular-nums">{{ discardsLeft }}</b> / {{ MAX_DISCARDS }}</span>
                <span>เลือกทิ้งรอบนี้: <b class="tabular-nums">{{ selectedIndices.size }}</b> / 5</span>
              </div>

              <button type="button"
                class="px-4 py-2 rounded-xl font-semibold transition bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
                @click="discardAndDraw" :disabled="turnEnded || discardsLeft === 0 || selectedIndices.size === 0">ทิ้ง &
                จั่ว ({{ selectedIndices.size }})</button>

              <button type="button"
                class="px-4 py-2 rounded-xl font-semibold bg-white/10 text-indigo-100 border border-white/15 hover:bg-white/20 disabled:opacity-50"
                @click="endTurn" :disabled="turnEnded">จบเทิร์น / โจมตี</button>
            </div>
          </div>
        </div>
      </div>
    </footer>

    <!-- Deck modal -->
    <div v-if="showDeck" class="fixed inset-0 z-[60] bg-black/70 backdrop-blur-sm grid place-items-center p-4">
      <div class="w-full max-w-4xl rounded-2xl border border-white/10 bg-white/10 p-4 max-h-[85vh] overflow-y-auto">
        <div
          class="flex items-center justify-between mb-3 sticky top-0 -m-4 p-3 pl-4 pr-4 bg-white/10 backdrop-blur-md border-b border-white/10 rounded-t-2xl z-10">
          <h4 class="text-lg font-bold text-indigo-100">สำรับทั้งหมด (52 ใบ)</h4>
          <button class="px-3 py-1 rounded-lg border border-white/10 bg-white/10 text-slate-200 hover:bg-white/20"
            @click="showDeck = false">ปิด</button>
        </div>
        <div class="grid grid-cols-4 sm:grid-cols-8 md:grid-cols-13 gap-2">
          <div v-for="c in fullDeck" :key="serialize(c)" class="relative">
            <div class="card mini" :class="cardStateClass(c)">
              <div class="card-head">
                <span class="rank">{{ c.rank }}</span>
                <span class="suit" :class="suitTextClass(c.suit)">{{ c.suit }}</span>
              </div>
              <div class="card-body">
                <span class="suit-big" :class="suitTextClass(c.suit)">{{ suitAnimal(c.suit) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="mt-3 text-xs text-slate-300/80 space-y-1">
          <div>สถานะสีเทา = ถูกเปิดใช้/ทิ้งแล้วในเทิร์นนี้</div>
          <div>ยังมีสี = ยังอยู่ในสำรับหรือบนมือ</div>
        </div>
      </div>
    </div>

    <!-- How-to modal -->
    <div v-if="showHowTo" class="fixed inset-0 z-[70] bg-black/70 backdrop-blur-sm grid place-items-center p-4">
      <div class="w-full max-w-4xl rounded-2xl border border-white/10 bg-white/10 p-4 max-h-[85vh] overflow-y-auto">
        <div
          class="flex items-center justify-between mb-3 sticky top-0 -m-4 p-3 pl-4 pr-4 bg-white/10 backdrop-blur-md border-b border-white/10 rounded-t-2xl z-10">
          <h4 class="text-lg font-bold text-indigo-100">วิธีเล่น & คอมโบ</h4>
          <button class="px-3 py-1 rounded-lg border border-white/10 bg-white/10 text-slate-200 hover:bg-white/20"
            @click="showHowTo = false">ปิด</button>
        </div>

        <div class="space-y-4 text-slate-200 text-sm">
          <div>
            <h5 class="font-semibold text-indigo-200 mb-1">วิธีเล่นย่อ</h5>
            <ul class="list-disc list-inside text-slate-300">
              <li>แจก 8 ใบ/เทิร์น</li>
              <li>ทิ้งได้ครั้งละ ≤5 ใบ รวมสูงสุด {{ MAX_DISCARDS }} ครั้ง/เทิร์น</li>
              <li>กดปุ่ม <b>จบเทิร์น / โจมตี</b> เพื่อคำนวณคอมโบและโจมตี (การคลิกการ์ดจะไม่โจมตี)</li>
              <li>บอสโต้กลับทุก {{ bossAttackEvery }} เทิร์น</li>
            </ul>
          </div>

          <div>
            <h5 class="font-semibold text-indigo-200 mb-2">ตัวอย่างคอมโบ</h5>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div v-for="ex in HOWTO_EXAMPLES" :key="ex.key" class="rounded-xl border border-white/10 bg-white/5 p-3">
                <div class="flex items-center justify-between mb-2">
                  <span class="font-bold">{{ ex.label }}</span>
                  <span class="text-xs text-slate-400">Base {{ ex.base }}</span>
                </div>
                <div class="grid grid-cols-5 gap-2">
                  <div v-for="(c, i) in ex.sample" :key="i" class="relative">
                    <div class="card mini">
                      <div class="card-head">
                        <span class="rank">{{ c.rank }}</span>
                        <span class="suit" :class="suitTextClass(c.suit)">{{ c.suit }}</span>
                      </div>
                      <div class="card-body">
                        <span class="suit-big" :class="suitTextClass(c.suit)">{{ suitAnimal(c.suit) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Score Modal -->
    <div v-if="showScoreModal"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[95] px-4" role="dialog"
      aria-modal="true" aria-labelledby="scoreTitle">
      <div class="w-full max-w-xl">
        <div
          class="relative rounded-3xl overflow-hidden shadow-[0_30px_80px_rgba(0,0,0,0.55)] border border-white/10 bg-white/5 backdrop-blur-xl">
          <div class="px-6 py-5 text-white bg-gradient-to-r from-indigo-600/90 via-indigo-500/90 to-fuchsia-600/90">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-full bg-white/20 flex items-center justify-center shrink-0">🏆</div>
              <div>
                <h3 id="scoreTitle" class="text-xl font-extrabold tracking-wide">บันทึกคะแนน</h3>
                <p class="text-white/90 text-sm">รอบนี้คุณได้ <b class="tabular-nums">{{ finalScore }}</b> คะแนน</p>
              </div>
            </div>
          </div>

          <div class="px-6 pt-5 pb-6 space-y-4 text-slate-100">
            <div>
              <label class="block text-sm font-medium text-slate-200 mb-1"
                for="playerName">กรอกชื่อเพื่อบันทึกสถิติ</label>
              <input id="playerName" v-model="playerName" type="text" placeholder="เช่น น้องแมว"
                class="px-4 py-2.5 rounded-xl w-full bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60" />
            </div>

            <div class="mt-2 flex flex-col sm:flex-row gap-2">
              <button @click="saveScore"
                class="inline-flex justify-center items-center px-5 py-3 rounded-xl font-semibold w-full sm:w-1/2 transition
                       bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
                :disabled="!playerName.trim() || isSaving" type="button">
                <span v-if="!isSaving">บันทึกคะแนน</span>
                <span v-else
                  class="inline-block animate-spin h-5 w-5 border-2 border-white/80 border-t-transparent rounded-full"></span>
              </button>

              <button @click="restartGame" class="inline-flex justify-center items-center px-5 py-3 rounded-xl font-semibold w-full sm:w-1/2 transition
                       bg-white/10 text-indigo-100 border border-white/15 hover:bg-white/20" type="button">
                เล่นใหม่
              </button>
            </div>
          </div>

          <div class="px-6 py-3 bg-white/5 border-t border-white/10 text-[12px] text-slate-300">
            บันทึกแล้วจะแสดงใน “สถิติผู้เล่น TOP 10” ด้านล่าง
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import api from '../services/api'
import { useRouter } from 'vue-router'
import { useNetworkStore } from '../store/useNetworkStore'
import { waitApiReadyAndLoadInitial } from '../composables/useApiReadiness'
import catwalk from '../assets/images/catwalk.png'
import catwalk2 from '../assets/images/catwalk2.png'

const net = useNetworkStore()

const loading = ref(true)
const catwalkImages = [catwalk, catwalk2]
const catwalkIndex = ref(0)
let catwalkInterval: number | undefined
let readinessTimer: number | undefined



const router = useRouter()
const goBack = () => router.back()

/** ------- คะแนน/บอร์ดสกอร์ ------- */
const GAME_NAME = 'CatGame'
const TOP_LIMIT = 10
const score = ref(0)
const savedScores = ref<{ name: string; score: number }[]>([])
const showScoreModal = ref(false)
const finalScore = ref(0)
const playerName = ref('')
const isSaving = ref(false)

async function loadScores() {
  const res = await api.get('/api/scores', { params: { limit: TOP_LIMIT, gamename: GAME_NAME } })
  savedScores.value = (res.data || []).slice(0, TOP_LIMIT)
}
async function saveScore() {
  if (!playerName.value.trim() || isSaving.value) return
  isSaving.value = true
  try {
    await api.post('/api/scores', {
      name: playerName.value.trim(),
      score: finalScore.value || score.value,
      gamename: GAME_NAME,
    })
    await loadScores()
    showScoreModal.value = false
    playerName.value = ''
  } catch (e) {
    console.error(e)
  } finally {
    isSaving.value = false
  }
}
function openScoreModal() {
  finalScore.value = score.value
  showScoreModal.value = true
}
function restartGame() {
  showScoreModal.value = false
  playerName.value = ''
  score.value = 0
  hardReset()
}

/** ------- Config ------- */
const MAX_DISCARDS = 3
const bossAttackEvery = 3 // lv.1

/** ------- Types ------- */
export type Suit = '♠' | '♥' | '♦' | '♣'
export type Rank = 'A' | 'K' | 'Q' | 'J' | '10' | '9' | '8' | '7' | '6' | '5' | '4' | '3' | '2'
interface Card { suit: Suit; rank: Rank; value: number }
type ComboKey = 'DEMONS_HAND' | 'MARCHING_HORDE' | 'TETRAD' | 'GRAND_WARHOST' | 'HORDE' | 'MARCH' | 'TRIAD' | 'DYAD_SET' | 'DYAD' | 'SOLO'
interface ComboMeta { key: ComboKey; label: string; base: number }
interface BonusLine { label: string; value: number }
interface EvalResult { key: ComboKey; comboLabel: string; base: number; bonuses: BonusLine[]; multiplier: number; total: number; usedIdx: number[]; note?: string }
interface Fighter { name: string; maxHp: number; hp: number }
interface Sigil { id: string; name: string; desc: string; add?: number; mult?: number; combos?: ComboKey[] }

/** ------- Constants ------- */
const SUITS: Suit[] = ['♠', '♥', '♦', '♣']
const RANKS: Rank[] = ['A', 'K', 'Q', 'J', '10', '9', '8', '7', '6', '5', '4', '3', '2']
const RANK_VALUE: Record<Rank, number> = { '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, '10': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14 }
const COMBOS: ComboMeta[] = [
  { key: 'DEMONS_HAND', label: "The Demon’s Hand (Royal Flush)", base: 220 },
  { key: 'MARCHING_HORDE', label: 'Marching Horde (Straight Flush)', base: 160 },
  { key: 'TETRAD', label: 'Tetrad (Four of a Kind)', base: 120 },
  { key: 'GRAND_WARHOST', label: 'Grand Warhost (Full House)', base: 90 },
  { key: 'HORDE', label: 'Horde (Flush)', base: 70 },
  { key: 'MARCH', label: 'March (Straight)', base: 60 },
  { key: 'TRIAD', label: 'Triad (Three of a Kind)', base: 45 },
  { key: 'DYAD_SET', label: 'Dyad Set (Two Pair)', base: 30 },
  { key: 'DYAD', label: 'Dyad (Pair)', base: 20 },
  { key: 'SOLO', label: 'Solo (High Card)', base: 10 },
]
const SIGIL_POOL: Sigil[] = [
  { id: 'sg_horde', name: 'Horde Specialist', desc: 'Flush ได้ +30 แต้ม', add: 30, combos: ['HORDE'] },
  { id: 'sg_march', name: 'March General', desc: 'Straight/Str. Flush ได้ ×1.2', mult: 1.2, combos: ['MARCH', 'MARCHING_HORDE'] },
  { id: 'sg_tetrad', name: 'Tetrad Surge', desc: 'Four of a Kind ได้ +25', add: 25, combos: ['TETRAD'] },
  { id: 'sg_warhost', name: 'Warhost Commander', desc: 'Full House ได้ ×1.15', mult: 1.15, combos: ['GRAND_WARHOST'] },
  { id: 'sg_unholy', name: 'Unholy Favor', desc: 'ทุกคอมโบได้ +10', add: 10 },
  { id: 'sg_royal', name: 'Demon’s Blessing', desc: 'Royal Flush ได้ ×1.5', mult: 1.5, combos: ['DEMONS_HAND'] },
]

/** ------- Reactive State ------- */
const deck = ref<Card[]>([])
const hand = ref<Card[]>([])
const fullDeck = ref<Card[]>([])
const usedThisTurn = ref<Set<string>>(new Set())
const selectedIndices = ref<Set<number>>(new Set())
const discardsLeft = ref<number>(MAX_DISCARDS)
const turnEnded = ref(false)
const result = ref<EvalResult | null>(null)

// Preview state
const previewDamage = ref(0)
const previewLabel = ref('')

// Battle State
const player = reactive<Fighter>({ name: 'Player', maxHp: 1000, hp: 1000 })
const boss = reactive<Fighter>({ name: 'Demon Lord', maxHp: 1200, hp: 1200 })
const playerAttacking = ref(false)
const playerBeingHit = ref(false)
const lastDamage = ref(0)
const bossLastDamage = ref(0)
const victory = ref(false)
const defeat = ref(false)
const turnCount = ref(0)

const bossImg = ref<string>('')

const bossHpPct = computed(() => (Math.max(0, boss.hp) / boss.maxHp) * 100)
const playerHpPct = computed(() => (Math.max(0, player.hp) / player.maxHp) * 100)
const turnsUntilCounter = computed(() => { const mod = turnCount.value % bossAttackEvery; return mod === 0 ? bossAttackEvery : bossAttackEvery - mod })

// Forecast next counter damage (deterministic per cycle)
const nextCounterDamage = ref<number | null>(null)
function ensureCounterForecast() { if (nextCounterDamage.value == null) nextCounterDamage.value = calcBossDamage() }

// Sigils
const availableSigils = ref<Sigil[]>(SIGIL_POOL)
const equippedSigils = ref<Sigil[]>([])
function equipSigil(s: Sigil) { if (equippedSigils.value.length < 3 && !equippedSigils.value.some(e => e.id === s.id)) equippedSigils.value.push(s) }
function unequipSigil(id: string) { equippedSigils.value = equippedSigils.value.filter(s => s.id !== id) }
function applySigils(ev: EvalResult): EvalResult {
  let addTotal = 0, multTotal = 1; const lines: BonusLine[] = []
  for (const s of equippedSigils.value) { const ok = !s.combos || s.combos.includes(ev.key); if (!ok) continue; if (s.add) { addTotal += s.add; lines.push({ label: s.name, value: s.add }) } if (s.mult) { multTotal *= s.mult; lines.push({ label: s.name + ' (×)', value: Math.round(ev.base * (s.mult - 1)) }) } }
  const total = Math.round((ev.base + addTotal) * multTotal)
  return { ...ev, bonuses: lines, multiplier: multTotal, total }
}

// Orientation (force landscape UX with overlay)
const isPortrait = ref(false)
function updateOrientation() { isPortrait.value = window.innerHeight > window.innerWidth }
onMounted(async () => {
  // orientation overlay
  updateOrientation()
  window.addEventListener('resize', updateOrientation)

  // ชื่อหน้า
  document.title = 'PETTEXT - CatGame'

  // ไอคอนโหลดสลับรูป
  catwalkInterval = window.setInterval(() => {
    catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length
  }, 200)

  // ✅ รอความพร้อมของ API แบบเดียวกับ DocumentsPage.vue
  const { healthOk, initialOk } = await waitApiReadyAndLoadInitial()

  const boot = () => {
    loading.value = false
    // เดิมคุณทำพวกนี้ใน onMounted — ย้ายมาเรียกหลัง “พร้อมจริง”
    hardReset()
    loadScores()
  }

  if (healthOk && initialOk) {
    boot()
  } else {
    // ตรวจซ้ำทุก 5 วิจนกว่าจะพร้อม
    readinessTimer = window.setInterval(async () => {
      const h = await waitApiReadyAndLoadInitial()
      if (h.healthOk && h.initialOk) {
        if (readinessTimer) window.clearInterval(readinessTimer)
        boot()
      }
    }, 5000)
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateOrientation)
  if (readinessTimer) clearInterval(readinessTimer)
  if (catwalkInterval) clearInterval(catwalkInterval)
})

// Suit helpers (colors + animals)
function suitTextClass(s: Suit) { return (s === '♥' || s === '♦') ? 'text-rose-400' : 'text-slate-100' }
function suitAnimal(s: Suit) { switch (s) { case '♥': return '🐶'; case '♠': return '🐱'; case '♦': return '🐻'; case '♣': return '🦦'; } }

// Deck helpers
function mkDeck(): Card[] { const out: Card[] = []; for (const s of SUITS) for (const r of RANKS) out.push({ suit: s, rank: r, value: RANK_VALUE[r] }); return out }
function serialize(c: Card) { return `${c.rank}|${c.suit}` }
function shuffle<T>(arr: T[]): T[] { const a = arr.slice(); for (let i = a.length - 1; i > 0; i--) { const j = Math.floor(Math.random() * (i + 1));[a[i], a[j]] = [a[j], a[i]] } return a }
function draw(n: number) { const take = deck.value.splice(0, n); hand.value.push(...take) }

// Deck modal
const showDeck = ref(false)
function openDeckModal() { showDeck.value = true }
function cardStateClass(c: Card) {
  if (usedThisTurn.value.has(serialize(c))) return 'is-used'
  const inDeck = deck.value.some(d => d.rank === c.rank && d.suit === c.suit)
  const inHand = hand.value.some(h => h.rank === c.rank && h.suit === c.suit)
  return (!inDeck && !inHand) ? 'is-used' : ''
}

// How-to modal
const showHowTo = ref(false)
function openHowToModal() { showHowTo.value = true }

// ======== SFX (optional) ========
const sfxAttackUrl = ref<string>('') // set to your file path if available
const sfxBossUrl = ref<string>('')   // set to your file path if available
let audioCtx: (AudioContext) = null as any

function playBeep(freq = 620) {
  try {
    const Ctx = (window as any).AudioContext || (window as any).webkitAudioContext
    if (!Ctx) return
    audioCtx = audioCtx || new Ctx()
    const o = audioCtx.createOscillator()
    const g = audioCtx.createGain()
    o.type = 'sawtooth'; o.frequency.value = freq
    o.connect(g); g.connect(audioCtx.destination)
    g.gain.setValueAtTime(0.12, audioCtx.currentTime)
    g.gain.exponentialRampToValueAtTime(0.0001, audioCtx.currentTime + 0.12)
    o.start(); o.stop(audioCtx.currentTime + 0.12)
  } catch (e) { /* noop */ }
}
function playSfxAttack() {
  if (sfxAttackUrl.value) { const a = new Audio(sfxAttackUrl.value); a.volume = 0.4; a.play().catch(() => { }) }
  else { playBeep(620) }
}
function playSfxBoss() {
  if (sfxBossUrl.value) { const a = new Audio(sfxBossUrl.value); a.volume = 0.4; a.play().catch(() => { }) }
  else { playBeep(380) }
}

// Battle helpers (attack only from endTurn)
function performAttack(dmg: number, onDone?: () => void) {
  if (!result.value) return
  lastDamage.value = dmg
  playerAttacking.value = true
  playSfxAttack()

  // เพิ่มคะแนนตามดาเมจที่ทำได้
  score.value += Math.max(0, dmg)

  animateHp(boss, dmg, () => {
    playerAttacking.value = false
    if (boss.hp <= 0) {
      victory.value = true
      openScoreModal()
      onDone && onDone()
      return
    }
    onDone && onDone()
  })
}
function bossCounterAttack() {
  if (victory.value || defeat.value) return
  const dmg = (nextCounterDamage.value ?? calcBossDamage())
  bossLastDamage.value = dmg
  nextCounterDamage.value = null
  playerBeingHit.value = true

  playSfxBoss()

  animateHp(player, dmg, () => {
    playerBeingHit.value = false
    if (player.hp <= 0) {
      defeat.value = true
      openScoreModal()
    }
    ensureCounterForecast()
  })
}
function calcBossDamage() { const base = 140, variance = 40; return Math.max(50, Math.round(base + (Math.random() * variance - variance / 2))) }
function animateHp(target: Fighter, dmg: number, onDone?: () => void) { const s = target.hp, e = Math.max(0, s - dmg); const frames = 24; let f = 0; const t = setInterval(() => { f++; target.hp = Math.round(s - ((s - e) * f) / frames); if (f >= frames) { clearInterval(t); target.hp = e; onDone && onDone() } }, 20) }

function hardReset() { victory.value = false; defeat.value = false; turnCount.value = 0; boss.hp = boss.maxHp; player.hp = player.maxHp; startNewTurn(true) }
function startNewTurn(resetAll = false) {
  usedThisTurn.value.clear()
  deck.value = shuffle(mkDeck())
  fullDeck.value = mkDeck()
  hand.value = []
  selectedIndices.value.clear()
  discardsLeft.value = MAX_DISCARDS
  turnEnded.value = false
  result.value = null
  previewDamage.value = 0; previewLabel.value = ''
  draw(8)
  if (resetAll) equippedSigils.value = equippedSigils.value.slice(0, 3)
  updatePreview()
  ensureCounterForecast()
}

function toggleSelect(idx: number) { if (turnEnded.value) return; if (selectedIndices.value.has(idx)) selectedIndices.value.delete(idx); else if (selectedIndices.value.size < 5) selectedIndices.value.add(idx) }
function onCardClick(idx: number) { toggleSelect(idx); updatePreview() }

function discardAndDraw() {
  if (turnEnded.value) return
  const cnt = selectedIndices.value.size
  if (cnt === 0 || discardsLeft.value <= 0) return
  const sorted = Array.from(selectedIndices.value).sort((a, b) => b - a)
  const removed = sorted.map(i => hand.value[i])
  for (const c of removed) usedThisTurn.value.add(serialize(c))
  const need = sorted.length
  for (const i of sorted) hand.value.splice(i, 1)
  selectedIndices.value.clear()
  draw(need)
  discardsLeft.value -= 1
  updatePreview()
}

function endTurn() {
  if (turnEnded.value || victory.value || defeat.value) return
  const ev = evaluateBest(hand.value)
  const withSigils = applySigils(ev)
  result.value = withSigils
  turnEnded.value = true

  // mark current hand as used for deck modal before turn resets
  for (const c of hand.value) usedThisTurn.value.add(serialize(c))

  // Attack boss, then handle counter (if any), and start next turn
  performAttack(withSigils.total, () => {
    turnCount.value += 1
    if (turnCount.value % bossAttackEvery === 0) { setTimeout(bossCounterAttack, 200) }
    startNewTurn()
  })
}

// Evaluation
function evaluateBest(cards: Card[]): EvalResult {
  const pickIdx = (picked: Card[]): number[] => { const idxs: number[] = []; const remain = picked.slice(); cards.forEach((c, i) => { const k = remain.findIndex(rc => rc.rank === c.rank && rc.suit === c.suit); if (k !== -1) { idxs.push(i); remain.splice(k, 1) } }); return idxs }
  const bySuit: Record<Suit, Card[]> = { '♠': [], '♥': [], '♦': [], '♣': [] }; for (const c of cards) bySuit[c.suit].push(c); for (const s of SUITS) bySuit[s].sort((a, b) => b.value - a.value)
  const isRoyal = (arr: Card[]) => { const vals = arr.map(c => c.value); const need = [10, 11, 12, 13, 14]; return need.every(v => vals.includes(v)) }
  const findStraight = (arr: Card[]): Card[] | null => {
    if (arr.length < 5) return null
    const u = Array.from(new Set(arr.map(c => c.value))).sort((a, b) => b - a)
    if (u[0] === 14) u.push(1)
    let run: number[] = []
    for (let i = 0; i < u.length; i++) {
      if (i === 0 || u[i] === u[i - 1] - 1) run.push(u[i]); else run = [u[i]]
      if (run.length >= 5) break
    }
    if (run.length >= 5) {
      const need = new Set(run.slice(0, 5))
      const take: Card[] = []
      for (const v of Array.from(need).sort((a, b) => b - a)) {
        const asVal = v === 1 ? 14 : v
        const found = arr.find(c => c.value === asVal && !take.includes(c))
        if (found) take.push(found)
      }
      return take.length === 5 ? take : null
    }
    return null
  }
  // Royal / Straight Flush
  for (const s of SUITS) {
    const suited = bySuit[s]
    const sfive = findStraight(suited)
    if (sfive) {
      const royal = isRoyal(sfive)
      const meta = COMBOS.find(c => c.key === (royal ? 'DEMONS_HAND' : 'MARCHING_HORDE'))!
      const used = pickIdx(sfive)
      return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used }
    }
  }
  // Four
  { const map = new Map<number, Card[]>(); for (const c of cards) { const arr = map.get(c.value) || []; arr.push(c); map.set(c.value, arr) } const quad = Array.from(map.values()).find(v => v.length >= 4); if (quad) { const meta = COMBOS.find(c => c.key === 'TETRAD')!; const used = pickIdx(quad.slice(0, 4)); return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used } } }
  // Full House
  { const map = new Map<number, Card[]>(); for (const c of cards) { const arr = map.get(c.value) || []; arr.push(c); map.set(c.value, arr) } const trips = Array.from(map.entries()).filter(([, v]) => v.length >= 3).sort((a, b) => b[0] - a[0]); const pairs = Array.from(map.entries()).filter(([, v]) => v.length >= 2).sort((a, b) => b[0] - a[0]); if (trips.length) { const t = trips[0][1].slice(0, 3); const p = pairs.find(([val]) => val !== trips[0][0])?.[1].slice(0, 2); if (p) { const meta = COMBOS.find(c => c.key === 'GRAND_WARHOST')!; const used = pickIdx([...t, ...p]); return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used } } } }
  // Flush
  for (const s of SUITS) {
    const suited = bySuit[s]
    if (suited.length >= 5) {
      const meta = COMBOS.find(c => c.key === 'HORDE')!
      const used = pickIdx(suited.slice(0, 5))
      return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used }
    }
  }
  // Straight
  { const sfive = findStraight(cards); if (sfive) { const meta = COMBOS.find(c => c.key === 'MARCH')!; const used = pickIdx(sfive); return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used } } }
  // Three
  { const map = new Map<number, Card[]>(); for (const c of cards) { const arr = map.get(c.value) || []; arr.push(c); map.set(c.value, arr) } const tri = Array.from(map.values()).find(v => v.length >= 3); if (tri) { const meta = COMBOS.find(c => c.key === 'TRIAD')!; const used = pickIdx(tri.slice(0, 3)); return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used } } }
  // Two Pair
  { const map = new Map<number, Card[]>(); for (const c of cards) { const arr = map.get(c.value) || []; arr.push(c); map.set(c.value, arr) } const pairs = Array.from(map.entries()).filter(([, v]) => v.length >= 2).sort((a, b) => b[0] - a[0]); if (pairs.length >= 2) { const p1 = pairs[0][1].slice(0, 2); const p2 = pairs[1][1].slice(0, 2); const meta = COMBOS.find(c => c.key === 'DYAD_SET')!; const used = pickIdx([...p1, ...p2]); return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used } } }
  // Pair
  { const map = new Map<number, Card[]>(); for (const c of cards) { const arr = map.get(c.value) || []; arr.push(c); map.set(c.value, arr) } const pair = Array.from(map.values()).find(v => v.length >= 2); if (pair) { const meta = COMBOS.find(c => c.key === 'DYAD')!; const used = pickIdx(pair.slice(0, 2)); return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used } } }
  // High
  { const high = cards.slice().sort((a, b) => b.value - a.value)[0]; const meta = COMBOS.find(c => c.key === 'SOLO')!; const used = high ? [cards.findIndex(c => c === high)] : []; return { key: meta.key, comboLabel: meta.label, base: meta.base, bonuses: [], multiplier: 1, total: meta.base, usedIdx: used } }
}

function updatePreview() {
  const ev = evaluateBest(hand.value)
  const withSigils = applySigils(ev)
  previewDamage.value = withSigils.total
  previewLabel.value = withSigils.comboLabel
}

// ===== How-to examples =====
function makeCard(s: Suit, r: Rank): Card { return { suit: s, rank: r, value: RANK_VALUE[r] } }
const HOWTO_EXAMPLES = [
  { key: 'DEMONS_HAND', label: COMBOS.find(c => c.key === 'DEMONS_HAND')!.label, base: 220, sample: [makeCard('♥', '10'), makeCard('♥', 'J'), makeCard('♥', 'Q'), makeCard('♥', 'K'), makeCard('♥', 'A')] },
  { key: 'MARCHING_HORDE', label: COMBOS.find(c => c.key === 'MARCHING_HORDE')!.label, base: 160, sample: [makeCard('♠', '5'), makeCard('♠', '6'), makeCard('♠', '7'), makeCard('♠', '8'), makeCard('♠', '9')] },
  { key: 'TETRAD', label: COMBOS.find(c => c.key === 'TETRAD')!.label, base: 120, sample: [makeCard('♥', '7'), makeCard('♠', '7'), makeCard('♦', '7'), makeCard('♣', '7'), makeCard('♣', 'A')] },
  { key: 'GRAND_WARHOST', label: COMBOS.find(c => c.key === 'GRAND_WARHOST')!.label, base: 90, sample: [makeCard('♣', 'Q'), makeCard('♥', 'Q'), makeCard('♠', 'Q'), makeCard('♦', '9'), makeCard('♣', '9')] },
  { key: 'HORDE', label: COMBOS.find(c => c.key === 'HORDE')!.label, base: 70, sample: [makeCard('♥', 'A'), makeCard('♥', '9'), makeCard('♥', '7'), makeCard('♥', '4'), makeCard('♥', '2')] },
  { key: 'MARCH', label: COMBOS.find(c => c.key === 'MARCH')!.label, base: 60, sample: [makeCard('♣', '6'), makeCard('♦', '7'), makeCard('♠', '8'), makeCard('♥', '9'), makeCard('♣', '10')] },
  { key: 'TRIAD', label: COMBOS.find(c => c.key === 'TRIAD')!.label, base: 45, sample: [makeCard('♦', '8'), makeCard('♣', '8'), makeCard('♥', '8'), makeCard('♠', 'K'), makeCard('♦', '4')] },
  { key: 'DYAD_SET', label: COMBOS.find(c => c.key === 'DYAD_SET')!.label, base: 30, sample: [makeCard('♠', 'K'), makeCard('♥', 'K'), makeCard('♦', '3'), makeCard('♣', '3'), makeCard('♥', '7')] },
  { key: 'DYAD', label: COMBOS.find(c => c.key === 'DYAD')!.label, base: 20, sample: [makeCard('♣', 'A'), makeCard('♥', 'A'), makeCard('♦', '5'), makeCard('♠', '9'), makeCard('♦', '2')] },
  { key: 'SOLO', label: COMBOS.find(c => c.key === 'SOLO')!.label, base: 10, sample: [makeCard('♠', 'A'), makeCard('♥', '7'), makeCard('♦', '4'), makeCard('♣', '9'), makeCard('♣', '2')] },
] as const
</script>

<style scoped>
.theme-modern {
  color-scheme: dark;
}

.tabular-nums {
  font-variant-numeric: tabular-nums;
}

.card {
  @apply w-full aspect-[5/7] rounded-2xl bg-gradient-to-b from-slate-800/80 to-slate-900/80 border border-white/10 p-3 text-left shadow hover:shadow-lg transition relative;
  transition: transform .18s ease, box-shadow .18s ease;
  will-change: transform;
}

.card.selected {
  transform: translateY(-6px) scale(1.2);
  z-index: 10;
}

.card.mini {
  @apply aspect-[3/4] p-2 rounded-xl;
}

.card-head {
  @apply flex items-center justify-between font-bold;
  font-size: 1.25rem;
  /* ~text-xl */
}

.card.mini .card-head {
  font-size: .9rem;
  /* bigger than text-xs */
}

.card-body {
  @apply flex items-center justify-center h-[calc(100%-1.9rem)];
  font-size: 3.75rem;
  /* ~text-6xl */
}

.card.mini .card-body {
  font-size: 1.9rem;
  /* ~text-3xl */
}

.rank {
  @apply text-slate-100;
}

.suit {
  @apply text-slate-200;
}

.suit-big {
  @apply opacity-90;
}

/* used/discarded state in deck modal */
.is-used {
  @apply opacity-50 grayscale;
}

.fighter-card {
  @apply rounded-xl border border-white/10 bg-white/5 p-3 relative overflow-hidden;
}

.hpbar {
  @apply w-full h-3 bg-slate-700/60 rounded-full overflow-hidden;
}

.hpfill {
  @apply h-full bg-emerald-400 transition-all duration-500;
}

.hpfill.boss {
  @apply bg-rose-400;
}

.attack-overlay {
  @apply absolute inset-0 pointer-events-none;
}

.slash {
  position: absolute;
  left: 10%;
  top: 10%;
  width: 80%;
  height: 3px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.9), transparent);
  transform: rotate(-20deg);
  animation: slashIn 0.35s ease-out;
}

.slash.reverse {
  transform: rotate(20deg);
  left: 10%;
  top: 70%;
}

@keyframes slashIn {
  from {
    opacity: 0;
    transform: translateY(-10px) rotate(-20deg);
  }

  to {
    opacity: 1;
    transform: translateY(20px) rotate(-20deg);
  }
}

.dmg-float {
  position: absolute;
  left: 50%;
  top: 0.75rem;
  transform: translateX(-50%);
  font-size: 1.5rem;
  font-weight: 900;
  color: rgb(252 165 165);
  animation: dmgUp 0.8s ease-out forwards;
  text-shadow: 0 2px 6px rgba(0, 0, 0, .35);
}

@keyframes dmgUp {
  0% {
    opacity: 0;
    transform: translate(-50%, 0);
  }

  20% {
    opacity: 1;
  }

  100% {
    opacity: 0;
    transform: translate(-50%, -40px);
  }
}

.shake {
  animation: shake 0.25s linear 2;
}

@keyframes shake {

  0%,
  100% {
    transform: translateX(0);
  }

  25% {
    transform: translateX(-2px);
  }

  75% {
    transform: translateX(2px);
  }
}

/* suit colors */
.text-rose-400 {
  color: rgb(251 113 133);
}

.text-slate-100 {
  color: rgb(241 245 249);
}
</style>
