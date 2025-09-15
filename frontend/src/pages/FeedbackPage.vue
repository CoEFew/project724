<template>
  <div class="min-h-screen relative overflow-x-hidden theme-modern">
    <!-- Gradient background: เหมือน DocumentsPage.vue -->
    <div class="pointer-events-none absolute inset-0 -z-10" aria-hidden="true">
      <div
        class="absolute inset-0 bg-[radial-gradient(90%_70%_at_70%_100%,rgba(99,102,241,0.45),transparent_60%),radial-gradient(60%_60%_at_0%_0%,rgba(59,130,246,0.35),transparent_60%),linear-gradient(180deg,#0b1020,#0b1120)]"
      />
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
    <div class="w-full max-w-3xl mx-auto px-4 py-8">
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

          <h1 class="text-2xl md:text-4xl font-extrabold tracking-wide text-indigo-300/90 uppercase text-center flex-1 drop-shadow-sm">
            Otter•Feedback
          </h1>

          <div class="w-[90px] sm:w-[120px]" />
        </div>
        <p class="text-slate-300/80 text-xs md:text-sm text-center">
            <span class="text-lg">🦦</span>
        </p>
      </header>

      <!-- Glass card -->
      <section
        class="w-full mx-auto rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-6 space-y-5"
      >
        <form class="space-y-4" @submit.prevent="submitFeedback">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-medium text-slate-200 mb-1" for="fbName">ชื่อ (ไม่บังคับ)</label>
              <input
                id="fbName"
                v-model="name"
                type="text"
                maxlength="64"
                placeholder="เช่น น้องหมา"
                class="px-4 py-2.5 rounded-xl w-full bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60"
                autocomplete="name"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-200 mb-1" for="fbContact">ติดต่อกลับ (อีเมล/ไลน์) (ไม่บังคับ)</label>
              <input
                id="fbContact"
                v-model="contact"
                type="text"
                maxlength="128"
                placeholder="เช่น my@mail.com หรือ @lineid"
                class="px-4 py-2.5 rounded-xl w-full bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60"
                autocomplete="email"
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-200 mb-1" for="fbMsg">ข้อความ (จำเป็น)</label>
            <textarea
              id="fbMsg"
              v-model="message"
              :maxlength="MAX_LEN"
              rows="6"
              placeholder="พิมพ์ข้อเสนอแนะ/ปัญหาที่พบ รายละเอียดที่เป็นประโยชน์"
              class="px-4 py-3 rounded-xl w-full bg-white/5 border border-white/15 text-slate-100 placeholder:slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-400/60"
            ></textarea>
            <div class="mt-1 text-[12px] text-slate-300/80 flex items-center justify-between">
              <span :class="minValid ? 'text-emerald-300' : 'text-rose-300'">
                {{ message.trim().length }} / {{ MAX_LEN }} ตัวอักษร
              </span>
              <span v-if="error" class="text-rose-300">{{ error }}</span>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <button
              type="submit"
              class="px-5 py-3 rounded-xl font-semibold transition bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
              :disabled="sending || !minValid"
              :aria-busy="sending ? 'true' : 'false'"
            >
              <span v-if="!sending">ส่ง Feedback</span>
              <span v-else class="inline-flex items-center">
                <span class="mr-2 inline-block animate-spin h-5 w-5 border-2 border-white/80 border-t-transparent rounded-full"></span>
                กำลังส่ง...
              </span>
            </button>

            <button
              type="button"
              class="px-5 py-3 rounded-xl font-semibold bg-white/10 text-indigo-100 border border-white/15 hover:bg-white/20"
              @click="clearForm"
              :disabled="sending"
            >
              ล้างฟอร์ม
            </button>
          </div>

          <p v-if="success" class="text-emerald-300 text-sm">
            ✔️ ขอบคุณสำหรับ Feedback! เราบันทึกข้อมูลของคุณเรียบร้อยแล้ว
          </p>
        </form>
      </section>

      <!-- ล่าสุด (ไม่บังคับ) -->
      <section
        class="w-full mx-auto mt-6 rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-5"
      >
        <div class="flex items-center justify-between">
          <h4 class="text-lg font-bold text-indigo-100">Feedback ล่าสุด</h4>
          <button
            @click="loadRecent"
            class="text-xs px-3 py-1 rounded-lg border border-white/10 bg-white/5 text-slate-200 hover:bg-white/10 transition"
            title="รีเฟรช"
            type="button"
          >
            รีเฟรช
          </button>
        </div>

        <p v-if="recent.length === 0" class="text-slate-300/70 text-sm mt-2">
          ยังไม่มีรายการ แสดงเมื่อมีการส่ง Feedback
        </p>

        <ul v-else class="mt-2 divide-y divide-white/10">
          <li v-for="(f, idx) in recent" :key="f.created_at + '_' + idx" class="py-2">
            <div class="text-slate-100 text-sm">
              <div class="flex items-center justify-between">
                <span class="font-semibold">
                  {{ f.name || 'ไม่ระบุชื่อ' }}
                  <span v-if="f.contact" class="text-slate-300/80 font-normal"> ({{ f.contact }})</span>
                </span>
                <span class="text-[11px] text-slate-300/70">{{ new Date(f.created_at).toLocaleString() }}</span>
              </div>
              <p class="mt-1 text-slate-200/90 whitespace-pre-line">{{ f.message }}</p>
            </div>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import catwalk from '../assets/images/catwalk.png'
import catwalk2 from '../assets/images/catwalk2.png'
import api from '../services/api'
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const loading = ref(true)
const catwalkImages = [catwalk, catwalk2]
const catwalkIndex = ref(0)
let catwalkInterval: number | undefined

const name = ref('')
const contact = ref('')
const message = ref('')
const sending = ref(false)
const success = ref(false)
const error = ref('')
const recent = ref<Array<{ name: string; contact: string; message: string; created_at: string }>>([])

const MAX_LEN = 1000
const MIN_LEN = 1
const minValid = computed(() => message.value.trim().length >= MIN_LEN && message.value.trim().length <= MAX_LEN)

function clearForm() {
  name.value = ''
  contact.value = ''
  message.value = ''
  success.value = false
  error.value = ''
}

async function submitFeedback() {
  error.value = ''
  success.value = false
  if (!minValid.value) {
    error.value = `โปรดพิมพ์อย่างน้อย ${MIN_LEN} ตัวอักษร`
    return
  }
  sending.value = true
  try {
    await api.post('/api/feedback', {
      name: name.value.trim(),
      contact: contact.value.trim(),
      message: message.value.trim(),
      source: 'FeedbackPage', // ติด tag ต้นทางหน้า
    })
    success.value = true
    await loadRecent()
    message.value = ''
  } catch (e: any) {
    error.value = e?.response?.data || 'เกิดข้อผิดพลาดในการส่งข้อมูล'
    console.error(e)
  } finally {
    sending.value = false
  }
}

async function loadRecent() {
  try {
    const res = await api.get('/api/feedback', { params: { limit: 5 } })
    recent.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const goBack = () => router.back()

onMounted(() => {
  document.title = 'PETTEXT - Feedback'
  catwalkInterval = setInterval(() => {
    catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length
  }, 200)
  setTimeout(() => (loading.value = false), 800)

  loadRecent()
})
</script>

<style scoped>
.tabular-nums { font-variant-numeric: tabular-nums; }
.theme-modern { color-scheme: dark; }
@keyframes pulseSoft {
  0%, 100% { transform: scale(1); filter: saturate(1); }
  50% { transform: scale(1.02); filter: saturate(1.05); }
}
.animate-pulse-soft { animation: pulseSoft 1.2s ease-in-out infinite; }
</style>
