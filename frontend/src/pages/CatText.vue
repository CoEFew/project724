<!-- อยากเปิดใช้งานเมื่อไหร่ แค่เปลี่ยน const ENABLE_AI = false → true ก็กลับมาใช้งานได้ทันที
ปุ่มและช่องกรอกถูก disabled ไว้แล้ว + แสดงข้อความเตือนเพื่อกันค่าใช้จ่าย -->
<!-- CatText.vue -->
<template>
  <div class="min-h-screen bg-paw-pattern flex flex-col uppercase">
    <!-- Banner (เป็นบล็อกปกติ ไม่ absolute เพื่อไม่ให้ทับ) -->
    <header class="w-full">
      <div class="mx-auto w-full max-w-3xl px-4 sm:px-6 lg:px-8 pt-4">
  <div class="flex items-start gap-3 bg-white/90 backdrop-blur rounded-2xl shadow-lg p-4 border border-yellow-200 uppercase">
          <img :src="paw" alt="paw" class="w-8 h-8 mt-0.5" />
          <p class="text-sm sm:text-base leading-relaxed text-gray-800 uppercase">
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
        <div class="bg-white/90 backdrop-blur rounded-2xl shadow-xl p-6 border border-blue-100">
          <div class="flex items-center gap-2 mb-3">
            <img :src="catBadge" alt="cat" class="w-8 h-8" />
            <h2 class="text-lg font-bold text-blue-700 uppercase">Chat กับ น้อนแมว</h2>
          </div>

          <textarea
            v-model="chatInput"
            placeholder="พิมพ์ข้อความ..."
            rows="1"
            @keydown="onKeydown"
            class="border rounded-xl px-4 py-3 w-full mb-3 resize-none outline-none focus:ring-2 focus:ring-blue-300 disabled:opacity-60"
            :disabled="true"
          />

          <button
            @click="sendChat"
            :disabled="true"
            class="px-4 py-3 bg-blue-400 text-white rounded-xl w-full font-semibold disabled:opacity-60 disabled:cursor-not-allowed uppercase"
            title="ปิดไว้ชั่วคราว"
          >
            ส่งข้อความ (ปิดไว้ชั่วคราว)
          </button>

          <button
            @click="router.back()"
            class="mt-3 px-4 py-3 bg-gray-300 hover:bg-gray-400 text-gray-800 rounded-xl w-full font-semibold transition uppercase"
          >
            ย้อนกลับ
          </button>

          <p class="mt-3 text-sm text-gray-600 flex items-start gap-2 uppercase">
            <img :src="lock" alt="lock" class="w-4 h-4 mt-0.5" />
            ค่าเปียกแมว
          </p>

          <div
            v-if="chatResponse"
            class="bg-gray-50 rounded-xl p-3 w-full text-gray-800 mt-3 border uppercase"
          >
            {{ chatResponse }}
          </div>
        </div>
      </div>
    </main>

    <!-- Loader (ทับทั้งจอ แต่เฉพาะตอนโหลด) -->
    <div v-if="loading" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-lg text-blue-200 font-semibold">กำลังโหลด...</span>
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
/* ลายพื้นหลังอุ้งเท้า (สัตว์เลี้ยง) */
.bg-paw-pattern {
  --c1: rgba(255, 182, 193, 0.2);
  --c2: rgba(173, 216, 230, 0.2);
  background:
    radial-gradient(circle at 20px 20px, var(--c1) 6px, transparent 7px) 0 0 / 40px 40px,
    radial-gradient(circle at 40px 30px, var(--c2) 6px, transparent 7px) 0 0 / 40px 40px,
    #f9fafb;
}
</style>
