<!-- อยากเปิดใช้งานเมื่อไหร่ แค่เปลี่ยน const ENABLE_AI = false → true ก็กลับมาใช้งานได้ทันที
ปุ่มและช่องกรอกถูก disabled ไว้แล้ว + แสดงข้อความเตือนเพื่อกันค่าใช้จ่าย -->
<!-- CatText.vue -->
<template>
  <div class="min-h-screen relative overflow-x-hidden theme-modern flex flex-col uppercase">
    <!-- Gradient background (โมเดิร์น, สบายตา) -->
    <div class="pointer-events-none absolute inset-0 -z-10" aria-hidden="true">
      <!-- ชั้นที่ 1: ไล่เฉดตาม reference -->
      <div
        class="absolute inset-0 bg-[radial-gradient(90%_70%_at_70%_100%,rgba(99,102,241,0.45),transparent_60%),radial-gradient(60%_60%_at_0%_0%,rgba(59,130,246,0.35),transparent_60%),linear-gradient(180deg,#0b1020,#0b1120)]"
      />
      <!-- ชั้นที่ 2: วงเรืองรองนุ่ม ๆ -->
      <div class="absolute -bottom-16 right-10 h-80 w-80 rounded-full blur-3xl opacity-40 bg-indigo-500/30" />
      <div class="absolute -top-12 left-[-4rem] h-72 w-72 rounded-full blur-3xl opacity-30 bg-fuchsia-500/25" />
    </div>

    <!-- Banner (เป็นบล็อกปกติ ไม่ absolute เพื่อไม่ให้ทับ) -->
    <header class="w-full">
      <div class="mx-auto w-full max-w-3xl px-4 sm:px-6 lg:px-8 pt-6">
        <div class="flex items-start gap-3 rounded-2xl shadow-lg p-4 border bg-white/5 border-white/10 backdrop-blur">
          <img :src="paw" alt="paw" class="w-8 h-8 mt-0.5" />
          <p class="text-sm sm:text-base leading-relaxed text-indigo-100/90">
            ถ้าอยากคุยกับ น้อนนนแมว สนทบทุน ค่าเปียกแมวมาให้หน่อยยย!
            <br class="hidden sm:block" />
            ตอนนี้ไม่มีตัง ขอปิดฟังก์ชันไว้ก่อน
          </p>
        </div>
      </div>
    </header>

    <!-- Main content: ใช้ container + spacing ป้องกันการทับกัน -->
    <main class="flex-1 w-full">
      <div class="mx-auto w-full max-w-3xl px-4 sm:px-6 lg:px-8 py-6 sm:py-10">
        <div class="rounded-2xl border border-white/10 bg-white/5 backdrop-blur-md shadow-[0_10px_30px_rgba(0,0,0,0.35)] p-6">
          <!-- Header แบบเดียวกับ Jigsaw + glass -->
          <div class="flex items-center justify-between mb-5">
            <button
              @click="router.back()"
              type="button"
              class="inline-flex items-center gap-2 px-3 py-2 rounded-xl border bg-white/5 border-white/10 text-slate-100 hover:bg-white/10 transition shadow-sm"
            >
              <svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M15 18l-6-6 6-6"></path>
              </svg>
              ย้อนกลับ
            </button>

            <h1 class="text-2xl md:text-3xl font-extrabold tracking-wide text-indigo-300/90 uppercase text-center flex-1 drop-shadow-sm">
              Cat•Text
            </h1>

            <div class="w-[90px] sm:w-[120px]" />
          </div>

          <div class="flex items-center gap-2 mb-3">
            <img :src="catBadge" alt="cat" class="w-8 h-8" />
            <h2 class="text-lg font-bold text-indigo-100">Chat กับ น้อนแมว</h2>
          </div>

          <textarea
            v-model="chatInput"
            placeholder="พิมพ์ข้อความ..."
            rows="1"
            @keydown="onKeydown"
            class="rounded-xl px-4 py-3 w-full mb-3 resize-none outline-none
                   bg-white/5 border border-white/15 text-slate-100 placeholder:text-slate-400
                   focus:ring-2 focus:ring-indigo-400/60 focus:border-indigo-400/60 disabled:opacity-60"
            :disabled="true"
          />

          <button
            @click="sendChat"
            :disabled="true"
            class="px-4 py-3 rounded-xl w-full font-semibold transition
                   bg-indigo-500 text-white hover:bg-indigo-400 disabled:opacity-50 disabled:cursor-not-allowed shadow"
            title="ปิดไว้ชั่วคราว"
          >
            ส่งข้อความ (ปิดไว้ชั่วคราว)
          </button>

          <button
            @click="router.back()"
            class="mt-3 px-4 py-3 rounded-xl w-full font-semibold transition
                   bg-white/10 text-indigo-100 border border-white/15 hover:bg-white/20"
            type="button"
          >
            ย้อนกลับ
          </button>

          <p class="mt-3 text-sm text-indigo-200/90 flex items-start gap-2">
            <img :src="lock" alt="lock" class="w-4 h-4 mt-0.5 opacity-90" />
            ค่าเปียกแมว
          </p>

          <div
            v-if="chatResponse"
            class="rounded-xl p-3 w-full mt-3 border
                   bg-white/5 border-white/10 text-slate-100"
          >
            {{ chatResponse }}
          </div>
        </div>
      </div>
    </main>

    <!-- Loader (ทับทั้งจอ แต่เฉพาะตอนโหลด) -->
    <div v-if="loading" class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-base md:text-lg text-indigo-100 font-semibold">กำลังโหลด...</span>
      </div>
    </div>

    <!-- ตกแต่งมุม (ซ่อนบนจอเล็ก และลด z เพื่อไม่บังคอนเทนต์) -->
    <img
      :src="catCorner"
      alt="decor cat"
      class="hidden md:block w-28 lg:w-32 h-auto fixed bottom-4 left-4 opacity-80 select-none pointer-events-none z-10"
    />
    <img
      :src="dogCorner"
      alt="decor dog"
      class="hidden md:block w-28 lg:w-32 h-auto fixed top-4 right-4 opacity-80 select-none pointer-events-none z-10"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api';

const ENABLE_AI = false;

const chatInput = ref('');
const chatResponse = ref('');
const chatLoading = ref(false);

type ChatAPIResponse = {
  output_text?: string;
  error?: { message?: string };
};

async function sendChat() {
  if (!ENABLE_AI) {
    chatResponse.value = 'ฟังก์ชันคุยกับ AI ถูกปิดชั่วคราวเพื่อกันค่าใช้จ่าย';
    return;
  }

  if (!chatInput.value.trim() || chatLoading.value) return;
  chatLoading.value = true;
  chatResponse.value = '';
  try {
    const res = await api.post<ChatAPIResponse>('/api/chat', { message: chatInput.value });
    const text = res.data?.output_text?.trim();
    if (text && text.length > 0) {
      chatResponse.value = text;
    } else if (res.data?.error?.message) {
      chatResponse.value = `เกิดข้อผิดพลาดจากโมเดล: ${res.data.error.message}`;
    } else {
      chatResponse.value = res.data?.output_text?.trim() || 'ไม่มีข้อความตอบกลับจากโมเดล';
    }
  } catch (e: any) {
    chatResponse.value = e?.response?.data?.error?.message
      ? `เกิดข้อผิดพลาด: ${e.response.data.error.message}`
      : 'เกิดข้อผิดพลาดในการเชื่อมต่อ AI';
  } finally {
    chatLoading.value = false;
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    e.preventDefault();
    sendChat();
  }
}

/** assets (ระวังชื่อไฟล์ให้ตรงกับของจริงในโปรเจกต์) */
import catwalk from '../assets/images/catwalk.png';
import catwalk2 from '../assets/images/catwalk2.png';
import paw from '../assets/images/cat.png';
import catBadge from '../assets/images/catwalk.png';
import dogBadge from '../assets/images/catwalk2.png';
import catCorner from '../assets/images/dog3.png';
import dogCorner from '../assets/images/dog4.png';
import lock from '../assets/images/dog.png';

const loading = ref(true);
const catwalkImages = [catwalk, catwalk2];
const catwalkIndex = ref(0);
let catwalkInterval: number | undefined;

onMounted(() => {
  document.title = 'PETTEXT - CatText';
  catwalkInterval = setInterval(() => {
    catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length;
  }, 200);
  setTimeout(() => {
    loading.value = false;
  }, 800);
});

const showModal = ref(false);
const router = useRouter();
function closeModal() {
  showModal.value = false;
  router.push('/');
}
</script>

<style scoped>
.tabular-nums { font-variant-numeric: tabular-nums; }

/* ลายพื้นหลังอุ้งเท้า (สัตว์เลี้ยง) */
/* เก็บคอมเมนต์และคลาสเดิมไว้ตามคำขอ แม้หน้าใหม่จะใช้ gradient เป็นหลัก */
.bg-paw-pattern {
  --c1: rgba(255, 182, 193, 0.2);
  --c2: rgba(173, 216, 230, 0.2);
  background:
    radial-gradient(circle at 20px 20px, var(--c1) 6px, transparent 7px) 0 0 / 40px 40px,
    radial-gradient(circle at 40px 30px, var(--c2) 6px, transparent 7px) 0 0 / 40px 40px,
    #f9fafb;
}

/* โหมดธีม (จูน contrast และ smoothing) */
.theme-modern { color-scheme: dark; }

/* ลดจังหวะ animation ให้สบายตา */
@keyframes pulseSoft {
  0%, 100% { transform: scale(1); filter: saturate(1); }
  50% { transform: scale(1.02); filter: saturate(1.05); }
}
.animate-pulse-soft { animation: pulseSoft 1.2s ease-in-out infinite; }
</style>
