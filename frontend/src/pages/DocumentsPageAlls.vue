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
        <div v-if="loading"
            class="fixed inset-0 bg-black/60 backdrop-blur-sm flex flex-col items-center justify-center z-[90]">
            <div class="flex flex-col items-center">
                <span class="text-base md:text-lg text-indigo-100 font-semibold">กำลังโหลด...</span>
            </div>
        </div>

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
                        DOG • PUZZLE (PARTY)
                    </h1>

                    <div class="flex gap-2 items-center">
                        <button type="button" @click="soundOn = !soundOn"
                            class="px-3 py-2 rounded-xl border border-white/10 bg-white/5 text-slate-100 hover:bg-white/10">
                            <span v-if="soundOn">🔊</span><span v-else>🔇</span>
                        </button>
                    </div>
                </div>
            </header>

            <!-- ===== LOBBY ===== -->
            <section v-if="phase === 'lobby'"
                class="w-full max-w-3xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-5 space-y-5">
                <div class="text-slate-200 text-sm">
                    <div class="mb-2">รหัสห้อง: <b class="tabular-nums">{{ room?.code || '—' }}</b></div>
                    <div>ผู้เล่น: <b>{{ players.length }}</b> / {{ room?.max_players || 4 }}</div>
                    <p class="mt-3 text-xs text-slate-400">
                        * เจ้าของห้องสามารถกด “เริ่มเกม” ได้ทันที (รองรับเล่นคนเดียว)
                        หรือเปิดอีกแท็บเข้าด้วยลิงก์/โค้ดเดียวกัน
                    </p>
                </div>

                <!-- join / ready / start -->
                <div class="flex flex-col sm:flex-row gap-2">
                    <input v-model="playerName" type="text" placeholder="กรอกชื่อของคุณ"
                        class="flex-1 px-4 py-2.5 rounded-xl bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60" />

                    <button v-if="!joined" @click="joinRoom"
                        class="px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
                        :disabled="!playerName.trim() || joining">
                        {{ joining ? 'กำลังเข้าร่วม...' : 'สร้างห้อง' }}
                    </button>

                    <button v-else @click="toggleReady"
                        class="px-4 py-2.5 rounded-xl font-semibold bg-white/10 text-indigo-100 border border-white/15 hover:bg-white/20">
                        {{ meReady ? 'ยังไม่พร้อม' : 'พร้อมแล้ว' }}
                    </button>

                    <button type="button" @click="copyInvite"
                        class="px-3 py-2 rounded-xl border border-white/10 bg-white/5 text-slate-100 hover:bg-white/10"
                        :disabled="!room">
                        คัดลอกลิงก์ห้อง
                    </button>

                    <button v-if="isOwner && joined" @click="startGame" :disabled="starting"
                        class="px-4 py-2.5 rounded-xl font-semibold bg-emerald-500 text-white hover:bg-emerald-400 disabled:opacity-50 shadow">
                        {{ starting ? 'กำลังเริ่ม...' : 'เริ่มเกม' }}
                    </button>
                </div>

                <!-- join by code -->
                <div class="rounded-xl border border-white/10 bg-white/5 p-3">
                    <div class="text-slate-200 text-sm font-semibold mb-2">เข้าห้องด้วยโค้ด</div>
                    <div class="flex gap-2">
                        <input v-model="joinCode" placeholder="ใส่โค้ด เช่น ABC123"
                            class="flex-1 px-3 py-2 rounded-lg bg-white/5 border border-white/15 text-slate-100" />
                        <button @click="goJoinByCode"
                            class="px-4 py-2 rounded-lg bg-indigo-500 text-white hover:bg-indigo-400">เข้าร่วมด้วยโค้ด</button>
                    </div>
                </div>

                <!-- players list -->
                <div class="rounded-xl border border-white/10 bg-white/5">
                    <div class="p-3 text-slate-200 text-sm font-semibold">ผู้เล่นในห้อง</div>
                    <ul class="divide-y divide-white/10">
                        <li v-for="p in players" :key="p.name"
                            class="px-3 py-2 text-slate-200 text-sm flex items-center justify-between">
                            <div class="flex items-center gap-2">
                                <span class="inline-block w-2 h-2 rounded-full"
                                    :class="p.is_ready ? 'bg-emerald-400' : 'bg-slate-400'"></span>
                                <span class="font-medium">{{ p.name }}</span>
                                <span v-if="p.is_owner"
                                    class="text-[11px] px-2 py-0.5 rounded-full bg-white/10 border border-white/10">เจ้าของ</span>
                                <span v-if="p.is_out"
                                    class="text-[11px] px-2 py-0.5 rounded-full bg-rose-500/20 border border-rose-400/30 text-rose-200">ตกรอบ</span>
                            </div>
                            <span class="tabular-nums text-slate-300">คะแนน: {{ p.score }}</span>
                        </li>
                    </ul>
                </div>
            </section>

            <!-- ===== PLAYING (ฝังเกมของ single-player มาที่นี่) ===== -->
            <section v-else-if="phase === 'playing'"
                class="w-full max-w-3xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-6 space-y-5">

                <!-- status & header (InfoCards แบบเดียวกับหน้าเดิม) -->
                <div class="space-y-3">
                    <h2 class="text-xl md:text-2xl font-extrabold text-indigo-100 tracking-wide text-center">
                        เก่งจริงก็ทายมาดิ! • ห้อง {{ room?.code }}
                    </h2>

                    <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
                        <InfoCard label="รอบ" :value="round?.round_no ?? '-'" accent="indigo" />
                        <InfoCard label="เวลาคงเหลือ" :value="remainSeconds + ' วินาที'" accent="sky" />
                        <InfoCard label="ผู้เล่น" :value="players.length + ' / ' + (room?.max_players || 4)"
                            accent="fuchsia" />
                    </div>

                    <div class="space-y-2">
                        <div class="w-full h-2 rounded-full overflow-hidden bg-white/10">
                            <div class="h-full transition-all duration-300 bg-gradient-to-r from-sky-400 via-indigo-400 to-fuchsia-400"
                                :style="{ width: Math.max(0, Math.min(100, (remainSeconds / (round?.seconds || 60)) * 100)) + '%' }" />
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

                <!-- answer form (reuse look & feel) -->
                <form class="w-full flex flex-col sm:flex-row items-stretch sm:items-center gap-2"
                    @submit.prevent="submitGuess">
                    <div class="relative flex-1">
                        <input ref="answerInput" v-model="guess" type="text"
                            :placeholder="meOut ? 'คุณตกรอบแล้ว — รอรอบถัดไป/จบเกม' : 'พิมพ์คำตอบที่นี่…'"
                            class="rounded-xl px-4 py-2.5 text-base bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60 focus:border-indigo-400/60 w-full"
                            :disabled="meOut || sending" />
                    </div>
                    <button type="submit"
                        class="px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
                        :disabled="meOut || !guessSanitized || sending">
                        ส่งคำตอบ
                    </button>
                    <button type="button" @click="showHint"
                        class="px-4 py-2.5 rounded-xl font-semibold w-full sm:w-auto transition bg-amber-500 text-slate-900 hover:bg-amber-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
                        :disabled="!round || hintBusy">
                        {{ hintBusy ? 'กำลังดึงคำใบ้…' : 'คำใบ้' }}
                    </button>
                </form>

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

            <!-- ===== GAME OVER ===== -->
            <section v-else-if="phase === 'over'"
                class="w-full max-w-3xl mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md p-5 space-y-4">
                <h3 class="text-xl font-extrabold text-indigo-100">ผลการแข่งขัน</h3>
                <p class="text-slate-200 text-sm">ผู้ชนะ: <b>{{ winner?.name || '—' }}</b></p>
                <ul class="divide-y divide-white/10">
                    <li v-for="(l, idx) in leaderboard" :key="l.name"
                        class="py-2 flex items-center justify-between text-sm text-slate-100">
                        <div class="flex items-center gap-2 min-w-0">
                            <span class="w-6 text-center">
                                <template v-if="idx === 0">🥇</template>
                                <template v-else-if="idx === 1">🥈</template>
                                <template v-else-if="idx === 2">🥉</template>
                                <template v-else>{{ idx + 1 }}.</template>
                            </span>
                            <span class="font-medium truncate max-w-[10rem] sm:max-w-[14rem]" :title="l.name">{{ l.name
                                }}</span>
                        </div>
                        <div class="tabular-nums font-semibold">{{ l.score }} คะแนน</div>
                    </li>
                </ul>

                <div class="flex gap-2">
                    <button @click="goBack"
                        class="px-4 py-2.5 rounded-xl font-semibold bg-white/10 text-indigo-100 border border-white/15 hover:bg-white/20">
                        ออกจากห้อง
                    </button>
                    <button v-if="isOwner" @click="restartRoom"
                        class="px-4 py-2.5 rounded-xl font-semibold bg-emerald-500 text-white hover:bg-emerald-400">
                        เริ่มใหม่
                    </button>
                </div>
            </section>
        </div>

        <!-- Toasts -->
        <transition-group name="toast" tag="div" class="fixed bottom-4 right-4 flex flex-col gap-2 z-[96]">
            <div v-for="t in toasts" :key="t.id" class="rounded-xl px-4 py-3 text-sm shadow border"
                :class="toastClass(t.type)">
                <strong class="mr-1">{{ t.title }}</strong>{{ t.message }}
            </div>
        </transition-group>
    </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch, defineComponent, h, type PropType } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../services/api'

/** ---------- route / state ---------- */
const router = useRouter()
const route = useRoute()
const loading = ref(true)
const soundOn = ref(true)

type Room = { code: string; max_players: number; status: 'waiting' | 'playing' | 'finished' }
type Player = { name: string; is_owner: boolean; is_ready: boolean; score: number; is_out: boolean }
type RoundPayload = { round_no: number; quiz_id: string; quiz_token: string; quiz_exp: number; seconds: number; level: number }

const room = ref<Room | null>(null)
const players = ref<Player[]>([])
const round = ref<RoundPayload | null>(null)

const phase = ref<'lobby' | 'playing' | 'over'>('lobby')
const joined = ref(false)
const isOwner = computed(() => !!players.value.find(p => p.name === playerName.value && p.is_owner))
const me = computed(() => players.value.find(p => p.name === playerName.value))
const meReady = computed(() => !!me.value?.is_ready)
const meOut = computed(() => !!me.value?.is_out)

const playerName = ref(localStorage.getItem('party_name') || '')
watch(playerName, v => localStorage.setItem('party_name', v || ''))

/** ---------- ws ---------- */
let ws: WebSocket | null = null
const remainSeconds = ref(0)

/** ---------- ui ---------- */
const joining = ref(false)
const starting = ref(false)
const sending = ref(false)
const hintBusy = ref(false)
const hints = ref<string[]>([])
const guess = ref('')
const guessSanitized = computed(() => guess.value.trim())
const lastResult = ref<boolean | null>(null)
const answerInput = ref<HTMLInputElement | null>(null)

/** ---------- my guesses (เฉพาะของฉันต่อ "ข้อปัจจุบัน") ---------- */
const myRoundGuesses = ref<string[]>([])
const guessSet = ref<Set<string>>(new Set()) // กันส่งคำซ้ำในข้อปัจจุบัน

/** ---------- join by code (รองรับวางลิงก์เต็ม) ---------- */
const joinCode = ref('')
function extractCodeFromInput(raw: string): string | null {
    const s = raw.trim()
    const m = s.match(/([A-Za-z0-9]{6})$/)
    return m ? m[1].toUpperCase() : null
}
function goJoinByCode() {
    const code = extractCodeFromInput(joinCode.value || '')
    if (!code) { toast('โค้ดไม่ถูกต้อง', 'กรุณาใส่โค้ด 6 ตัวหรือวางลิงก์ห้อง', 'error'); return }
    router.push({ name: 'DocumentsPageRoom', params: { code } })
}

/** ---------- toast ---------- */
const toasts = ref<{ id: string; title: string; message: string; type: 'info' | 'success' | 'error' }[]>([])
function toast(title: string, message: string, type: 'info' | 'success' | 'error' = 'info') {
    const id = Math.random().toString(36).slice(2)
    toasts.value.push({ id, title, message, type })
    setTimeout(() => { toasts.value = toasts.value.filter(t => t.id !== id) }, 3500)
}
function toastClass(type: 'info' | 'success' | 'error') {
    const base = 'bg-white/10 border-white/15 text-slate-100 backdrop-blur'
    if (type === 'success') return base + ' border-emerald-300/30 bg-emerald-400/10'
    if (type === 'error') return base + ' border-rose-300/30 bg-rose-400/10'
    return base
}

/** ---------- helpers ---------- */
function goBack() {
    // 🔧 หลังจบเกมให้ย้อนกลับหน้า home เสมอ
    if (phase.value === 'over') {
        router.back()
    }
    router.back()
}
async function copyInvite() {
    if (!room.value) return
    const url = `${location.origin}/party/${room.value.code}`
    try { await navigator.clipboard.writeText(url); toast('คัดลอกแล้ว', 'ส่งลิงก์ให้เพื่อนเข้าห้องได้เลย', 'success') }
    catch { toast('คัดลอกไม่สำเร็จ', url, 'info') }
}
function playTick(ok: boolean) {
    if (!soundOn.value) return
    try {
        const ctx = new (window.AudioContext || (window as any).webkitAudioContext)()
        const o = ctx.createOscillator()
        const g = ctx.createGain()
        o.type = ok ? 'triangle' : 'sawtooth'
        o.frequency.value = ok ? 780 : 220
        g.gain.value = ok ? 0.05 : 0.04
        o.connect(g); g.connect(ctx.destination); o.start()
        setTimeout(() => { o.stop(); ctx.close() }, 140)
    } catch { }
}

/** ---------- tiny components ---------- */
const InfoCard = defineComponent({
    name: 'InfoCard',
    props: { label: { type: String, required: true }, value: { type: [String, Number], required: true }, accent: { type: String as PropType<'indigo' | 'sky' | 'fuchsia' | 'emerald'>, required: true } },
    setup(props) {
        const base = 'rounded-xl py-2.5 px-4 text-center border '
        const cls = props.accent === 'indigo' ? base + 'bg-indigo-400/10 border-indigo-300/20'
            : props.accent === 'sky' ? base + 'bg-sky-400/10 border-sky-300/20'
                : props.accent === 'fuchsia' ? base + 'bg-fuchsia-400/10 border-fuchsia-300/20'
                    : base + 'bg-emerald-400/10 border-emerald-300/20'
        return () => h('div', { class: cls }, [
            h('div', { class: 'text-[11px] font-semibold text-slate-200/90 tracking-wide' }, props.label),
            h('div', { class: 'mt-0.5 text-2xl font-bold text-slate-100 tabular-nums' }, String(props.value)),
        ])
    }
})
const Chip = defineComponent({
    name: 'Chip',
    props: { label: { type: String, required: true }, value: { type: String, required: true }, color: { type: String as PropType<'fuchsia' | 'sky'>, required: true } },
    setup(props) {
        const cls = props.color === 'fuchsia'
            ? 'inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm border bg-fuchsia-500/10 text-fuchsia-100 border-fuchsia-300/20'
            : 'inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm border bg-sky-500/10 text-sky-100 border-sky-300/20'
        return () => h('span', { class: cls }, [
            h('strong', { class: 'font-semibold' }, props.label + ':'),
            h('span', null, props.value),
        ])
    }
})

/** ---------- bootstrap ---------- */
function ensureName() {
    if (!playerName.value.trim()) {
        playerName.value = 'Player-' + Math.random().toString(36).slice(2, 6).toUpperCase()
    }
}
async function bootstrapFromCode(code: string) {
    try { ws?.close() } catch { }
    ws = null
    joined.value = false
    myRoundGuesses.value = []
    guessSet.value = new Set()
    hints.value = []
    lastResult.value = null
    phase.value = 'lobby'
    loading.value = true

    try {
        // ✅ รอ BE ตื่น (สำคัญมากบน Render)
        await waitApiUp()
        room.value = { code, max_players: 4, status: 'waiting' }
        loading.value = false
        await joinRoom()
    } catch (e: any) {
        loading.value = false
        toast('เข้าห้องไม่สำเร็จ', e?.message || 'network error', 'error')
    }
}

// เพิ่ม util ง่าย ๆ ไว้ในไฟล์เดียวกัน (หรือ import จาก composable เดิมก็ได้)
async function waitApiUp(maxWaitMs = 20000) {
  const start = Date.now()
  while (Date.now() - start < maxWaitMs) {
    try {
      const res = await api.get('/health', { timeout: 4000 })
      if (String(res.data).toLowerCase().includes('ok')) return true
    } catch {}
    await new Promise(r => setTimeout(r, 1200)) // เว้น 1.2s ให้ instance ตื่น
  }
  return false
}


onMounted(async () => {
    ensureName()

const ok = await waitApiUp()
  if (!ok) {
    toast('เครือข่ายช้า', 'กำลังปลุกเซิร์ฟเวอร์ ลองใหม่อีกครั้ง', 'error')
    // ยังไปต่อได้ ระบบจะ retry ใน action ต่อไป
  }

    const rawParam = (route.params.code as string | undefined) || ''
    const pathTail = window.location.pathname.split('/').pop() || ''
    const paramCode = (rawParam || pathTail).toUpperCase()
    const looksLikeCode = /^[A-Z0-9]{6}$/.test(paramCode)

    try {
        if (looksLikeCode) {
            await bootstrapFromCode(paramCode) // ผู้เล่นเข้าห้องด้วยลิงก์
        } else {
            const res = await api.post('/api/rooms', { ownerName: playerName.value || 'Player', maxPlayers: 4 })
            const r = res.data.room
            room.value = { code: r.code, max_players: r.maxPlayers, status: r.status }
            loading.value = false
        }
    } catch (e: any) {
        loading.value = false
        toast('สร้าง/โหลดห้องล้มเหลว', e?.message || 'network error', 'error')
    }
})

// ✅ เปลี่ยนโค้ด/ลิงก์ในช่อง -> ไม่ต้องรีเฟรช
watch(() => route.params.code, async (val) => {
    const code = (val as string | undefined)?.toUpperCase()
    if (code && /^[A-Z0-9]{6}$/.test(code)) { await bootstrapFromCode(code) }
})

onBeforeUnmount(() => { try { ws?.close() } catch { } })

/** ---------- actions ---------- */
async function joinRoom() {
    if (!room.value) return
    if (!playerName.value.trim()) return
    joining.value = true
    try {
        await api.post(`/api/rooms/${room.value.code}/join`, { name: playerName.value.trim() })
        joined.value = true
        connectWS()
        await nextTick()
        answerInput.value?.focus()
    } catch (e: any) {
        toast('เข้าร่วมล้มเหลว', e?.message || 'network error', 'error')
    } finally {
        joining.value = false
    }
}
async function toggleReady() {
    if (!room.value || !joined.value) return
    try {
        await api.post(`/api/rooms/${room.value.code}/ready`, { name: playerName.value, ready: !meReady.value })
    } catch (e: any) {
        toast('เปลี่ยนสถานะพร้อมล้มเหลว', e?.message || 'network error', 'error')
    }
}

async function safePullSnapshot() {
  try {
    if (!room.value) return
    const res = await api.get(`/api/rooms/${room.value.code}/snapshot`)
    const m = res.data
    if (m?.room) room.value = normalizeRoom(m.room)
    if (Array.isArray(m?.players)) players.value = m.players.map(normalizePlayer)
    if (m?.round) {
      round.value = normalizeRound(m.round)
      phase.value = 'playing'
      remainSeconds.value = m.round.seconds || 60
      fetchFirstHintIfAny()
    }
  } catch {/* ignore */}
}

async function startGame() {
  if (!room.value) return
  starting.value = true
  try {
    await api.post(`/api/rooms/${room.value.code}/start`, { ownerName: playerName.value })
  } catch (e:any) {
    // บน Render บางที WS มาช้า → fallback ดึง snapshot
    await safePullSnapshot()
    // ไม่เด้ง error ถ้าเป็น 409/started ไปแล้ว
    const msg = String(e?.response?.status || '')
    if (msg !== '409') toast('เริ่มเกมล้มเหลว', e?.message || 'network error', 'error')
  } finally {
    starting.value = false
  }
}

// async function hardLeave() {
//   try {
//     if (room.value && playerName.value) {
//       await api.post(`/api/rooms/${room.value.code}/leave`, { name: playerName.value })
//     }
//   } catch {}
//   router.replace({ name: 'DocumentsPageAlls' })
// }
async function restartRoom() {
    hints.value = []; lastResult.value = null; remainSeconds.value = 0; phase.value = 'lobby'
    myRoundGuesses.value = []; guessSet.value = new Set()
    try { await api.post(`/api/rooms/${room.value!.code}/start`, { ownerName: playerName.value }) }
    catch (e: any) { toast('เริ่มใหม่ล้มเหลว', e?.message || 'network error', 'error') }
}

/** ---------- auto-start: owner เริ่มเมื่อพร้อม ---------- */
const autostarted = ref(false)

const allReady = computed(() => players.value.length > 0 && players.value.every(p => p.is_ready))
watch([players, joined, () => meReady.value], async () => {
  if (!joined.value || !isOwner.value || !room.value) return
  if (starting.value || autostarted.value) return
  const soloReady = players.value.length === 1 && meReady.value
  if (soloReady || allReady.value) {
    try {
      starting.value = true
      autostarted.value = true         // ✅ กันยิงซ้ำ
      await api.post(`/api/rooms/${room.value.code}/start`, { ownerName: playerName.value })
    } catch (e:any) {
      // ถ้า 409/200 ก็ปล่อยผ่าน (idempotent แล้ว) และดึง snapshot เผื่อ WS ไม่มา
      await safePullSnapshot()
    } finally {
      starting.value = false
    }
  }
}, { deep: true })


/** ---------- WS ---------- */
function connectWS() {
  if (!room.value) return
  const base = new URL((api as any).defaults?.baseURL || `${location.protocol}//${location.host}`)
  const wsProto = base.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsOrigin = `${wsProto}//${base.host}`
  ws = new WebSocket(`${wsOrigin}/ws/rooms/${room.value.code}`)
  ws.onopen = () => console.info('[WS] connected', wsOrigin)
  ws.onclose = (ev) => console.warn('[WS] closed', ev.code, ev.reason)
  ws.onmessage = (ev) => { try { handleWs(JSON.parse(ev.data)) } catch (e){ console.error('WS parse', e) } }
  ws.onerror = () => { toast('การเชื่อมต่อเรียลไทม์ผิดพลาด', 'กำลังพยายามเชื่อมใหม่…', 'error') }
}


async function fetchFirstHintIfAny() {
    // ทุกข้อโชว์ใบ้แรกเลย
    if (!round.value) return
    try {
        const res = await api.post('/api/quiz/hint', {
            id: round.value.quiz_id,
            token: round.value.quiz_token,
            exp: round.value.quiz_exp,
            index: 1,
        })
        const text = res?.data?.hint || ''
        hints.value = text ? [text] : []
    } catch (e: any) {
        // ไม่ fail เกม เพียงแค่ไม่แสดงใบ้
        console.warn('hint error', e?.message || e)
    }
}

function handleWs(m: any) {
    switch (m.type) {
        case 'snapshot':
            if (m.room) room.value = normalizeRoom(m.room)
            if (Array.isArray(m.players)) players.value = m.players.map(normalizePlayer)
            if (m.round) { round.value = normalizeRound(m.round); phase.value = 'playing'; remainSeconds.value = m.round.seconds; fetchFirstHintIfAny() }
            else phase.value = 'lobby'
            break
        case 'player_joined':
        case 'ready_changed':
            if (Array.isArray(m.players)) players.value = m.players.map(normalizePlayer)
            break
        case 'round_started':
            round.value = normalizeRound(m.round)
            // เริ่ม "ข้อใหม่" -> รีเซ็ตสถานะต่อข้อ
            hints.value = []; lastResult.value = null
            myRoundGuesses.value = []; guessSet.value = new Set()
            phase.value = 'playing'
            remainSeconds.value = round.value?.seconds || 60
            nextTick(() => answerInput.value?.focus())
            fetchFirstHintIfAny()
            break
        case 'timer_tick':
            remainSeconds.value = m.seconds || 0
            break
        case 'guess_result':
            if (Array.isArray(m.players)) players.value = m.players.map(normalizePlayer)
            // เก็บเฉพาะของฉัน
            if (m.name === playerName.value && typeof m.guess === 'string') {
                myRoundGuesses.value.unshift(m.guess)
                myRoundGuesses.value = myRoundGuesses.value.slice(0, 15)
            }
            if (typeof m.correct === 'boolean' && m.name === playerName.value) {
                lastResult.value = m.correct; playTick(m.correct)
            }
            break
        case 'game_over':
            phase.value = 'over'
            winner.value = m.winner ? normalizePlayer(m.winner) : null
            leaderboard.value = Array.isArray(m.leaderboard) ? m.leaderboard.map((x: any) => ({ name: x.name, score: x.score })) : []
            break
    }
}
function normalizeRoom(r: any): Room { return { code: r.code, max_players: r.maxPlayers ?? r.max_players ?? 4, status: r.status } }
function normalizePlayer(p: any): Player { return { name: p.name, is_owner: !!(p.isOwner ?? p.is_owner), is_ready: !!(p.isReady ?? p.is_ready), score: p.score ?? 0, is_out: !!(p.isOut ?? p.is_out) } }
function normalizeRound(r: any): RoundPayload {
    return { round_no: r.roundNo ?? r.round_no ?? 1, quiz_id: r.quiz_id ?? r.quizId ?? r.quizID, quiz_token: r.quiz_token ?? r.quizToken, quiz_exp: r.quiz_exp ?? r.quizExp, seconds: r.seconds ?? 60, level: r.level ?? 1 }
}

/** ---------- gameplay ---------- */
async function submitGuess() {
    if (!room.value || !round.value || !guessSanitized.value || sending.value) return
    const key = guessSanitized.value.toLowerCase()
    if (guessSet.value.has(key)) { toast('คำซ้ำ', 'คุณส่งคำนี้ไปแล้วสำหรับข้อปัจจุบัน', 'info'); return }
    sending.value = true
    try {
        await api.post(`/api/rooms/${room.value.code}/guess`, { name: playerName.value, guess: guessSanitized.value })
        guessSet.value.add(key)
        guess.value = ''
    } catch (e: any) {
        toast('ส่งคำตอบล้มเหลว', e?.message || 'network error', 'error')
    } finally {
        sending.value = false
    }
}
async function showHint() {
    if (!round.value) return
    hintBusy.value = true
    try {
        const res = await api.post('/api/quiz/hint', {
            id: round.value.quiz_id, token: round.value.quiz_token, exp: round.value.quiz_exp, index: hints.value.length + 1,
        })
        const text = res?.data?.hint || ''
        if (text) hints.value.push(text)
    } catch (e: any) {
        toast('ขอคำใบ้ล้มเหลว', e?.message || 'network error', 'error')
    } finally { hintBusy.value = false }
}

/** ---------- game over ---------- */
const winner = ref<Player | null>(null)
const leaderboard = ref<{ name: string; score: number }[]>([])

</script>

<style scoped>
.theme-modern {
    color-scheme: dark;
}

.tabular-nums {
    font-variant-numeric: tabular-nums;
}

.toast-enter-active,
.toast-leave-active {
    transition: all .25s ease;
}

.toast-enter-from,
.toast-leave-to {
    opacity: 0;
    transform: translateY(6px);
}
</style>
