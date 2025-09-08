<template>
  <div class="flex items-center justify-center min-h-screen bg-transparent">
    <div class="w-full max-w-lg bg-white bg-opacity-80 rounded-xl shadow-lg p-6 flex flex-col items-center">
      <div class="w-full mb-4">
        <div class="bg-blue-100 rounded-lg py-2 px-4 text-center text-blue-700 font-bold text-lg">
          คะแนน: {{ score }}
        </div>
        <div class="bg-gray-100 rounded-lg py-2 px-4 text-center text-gray-700 font-bold text-lg mt-2">
          เวลาคงเหลือ: {{ timer }} วินาที
        </div>
      </div>
      <h2 class="text-3xl font-bold mb-4 text-center">เก่งจริงก็ทายมาดิ!</h2>
      <div class="mb-6 w-full flex flex-col sm:flex-row items-center justify-center gap-2">
        <input v-model="guess" type="text" placeholder="พิมพ์คำตอบที่นี่..." class="border rounded px-4 py-2 text-lg focus:outline-none focus:ring-2 focus:ring-blue-400 w-full sm:w-auto" @input="checkAnswer" :disabled="showModal" />
        <button @click="showHint" class="px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600 transition w-full sm:w-auto" :disabled="showModal">คำใบ้</button>
        <button @click="fetchQuiz" class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition w-full sm:w-auto" :disabled="showModal">เปลี่ยนคำ</button>
      </div>
      <div v-if="result !== null" class="mb-4 text-center">
        <span v-if="result === true" class="text-green-600 text-2xl">✔️ ถูกต้อง!</span>
        <span v-else-if="result === false" class="text-red-600 text-2xl">❌ ผิด ลองใหม่!</span>
      </div>
      <div v-if="result === true" class="mb-2 text-center">
        <p class="text-xl font-bold text-blue-700">คำปริศนา: เฉลยจะแสดงเมื่อทายถูก</p>
      </div>
      <div v-if="hint1" class="mt-2 text-lg text-purple-700 text-center">คำใบ้ 1: {{ hint1 }}</div>
      <div v-if="hint2" class="mt-2 text-lg text-purple-700 text-center">คำใบ้ 2: {{ hint2 }}</div>
      <button @click="goBack" class="mt-8 px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition w-full sm:w-auto" :disabled="showModal">ย้อนกลับ</button>
    </div>
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl shadow-2xl p-8 flex flex-col items-center">
        <h3 class="text-2xl font-bold text-red-600 mb-4">หมดเวลา!</h3>
        <p class="mb-6 text-lg text-gray-700">คะแนนของคุณถูกรีเซตแล้ว</p>
        <button @click="restartGame" class="px-6 py-3 bg-blue-600 text-white rounded-lg font-bold text-lg">เริ่มใหม่</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import api from '../services/api';
import { ref, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
const router = useRouter();

const quizId = ref('');
const hints = ref<string[]>([]);
const guess = ref('');
const result = ref<null | boolean>(null);
const hint1 = ref('');
const hint2 = ref('');
const score = ref(0);
const timer = ref(60);
const hintCount = ref(0);
const showModal = ref(false);
let intervalId: number | undefined;

async function fetchQuiz() {
  const res = await api.get('/quiz');
  quizId.value = res.data.id;
  hints.value = res.data.hints;
  result.value = null;
  guess.value = '';
  hint1.value = '';
  hint2.value = '';
  hintCount.value = 0;
  timer.value = 60;
  showModal.value = false;
  if (intervalId) clearInterval(intervalId);
  intervalId = setInterval(() => {
    if (timer.value > 0) timer.value--;
    if (timer.value === 0) {
      clearInterval(intervalId);
      result.value = false;
      score.value = 0;
      showModal.value = true;
    }
  }, 1000);
}

async function checkAnswer() {
  if (!quizId.value) return;
  const res = await api.post('/quiz/check', { id: quizId.value, guess: guess.value.trim() });
  result.value = res.data.correct;
}

watch(result, (val, oldVal) => {
  if (val === true && oldVal !== true) {
    score.value++;
    setTimeout(() => fetchQuiz(), 1200);
  }
});

function showHint() {
  if (hintCount.value === 0) {
    hint1.value = hints.value[0];
    hintCount.value++;
  } else if (hintCount.value === 1) {
    hint2.value = hints.value[1];
    timer.value = Math.max(timer.value - 10, 0);
    hintCount.value++;
  }
}
function restartGame() {
  score.value = 0;
  fetchQuiz();
}
const goBack = () => {
  router.back();
};
onMounted(() => {
  fetchQuiz();
});
</script>