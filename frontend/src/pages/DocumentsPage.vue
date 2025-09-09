<template>
  <div class="flex flex-col items-center justify-start min-h-screen bg-transparent py-8 relative">
  <div v-if="loading" class="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-lg text-blue-700 font-semibold">กำลังโหลด...</span>
      </div>
    </div>
    <!-- กล่องเกม -->
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
        <input
          v-model="guess"
          type="text"
          placeholder="พิมพ์คำตอบที่นี่..."
          class="border rounded px-4 py-2 text-lg focus:outline-none focus:ring-2 focus:ring-blue-400 w-full sm:w-auto"
          @input="checkAnswer"
          :disabled="showModal"
        />
        <button
          @click="showHint"
          class="px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600 transition w-full sm:w-auto"
          :disabled="showModal"
        >
          คำใบ้
        </button>
        <button
          @click="() => fetchQuiz()"
          class="relative px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition w-full sm:w-auto"
          :disabled="showModal || changeCount >= maxChange"
        >
          <span
            class="absolute -top-2 -right-3 bg-white text-green-600 rounded-full px-2 py-0.5 text-xs font-bold shadow"
          >
            {{ changeCount }}/{{ maxChange }}
          </span>
          เปลี่ยนคำ
        </button>
      </div>

      <div v-if="result !== null" class="mb-4 text-center">
        <span v-if="result === true" class="text-green-600 text-2xl">✔️ ถูกต้อง!</span>
        <span v-else-if="result === false" class="text-red-600 text-2xl">❌ ผิด ลองใหม่!</span>
      </div>

      <div v-if="hint1" class="mt-2 text-lg text-purple-700 text-center">คำใบ้ 1: {{ hint1 }}</div>
      <div v-if="hint2" class="mt-2 text-lg text-purple-700 text-center">คำใบ้ 2: {{ hint2 }}</div>

      <div class="mt-8 flex gap-2">
        <button
          @click="goBack"
          class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition w-full sm:w-auto"
          :disabled="showModal"
        >
          ย้อนกลับ
        </button>
      </div>
    </div>

    <!-- โมดัล: หมดเวลา + บันทึกคะแนน -->
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl shadow-2xl p-8 flex flex-col items-center w-full max-w-lg">
        <h3 class="text-2xl font-bold text-red-600 mb-4">หมดเวลา!</h3>
        <p class="mb-2 text-lg text-gray-700">คะแนนของคุณ: {{ finalScore }}</p>
        <p class="mb-6 text-sm text-gray-500">*จะเริ่มใหม่หลังจากบันทึก หรือกดเริ่มใหม่</p>
        <input
          v-model="playerName"
          type="text"
          placeholder="กรอกชื่อผู้เล่น"
          class="mb-4 px-4 py-2 border rounded w-full text-center"
        />
        <div class="flex gap-2 w-full">
          <button
            @click="saveScore"
            class="px-6 py-3 bg-green-600 text-white rounded-lg font-bold text-lg w-1/2"
          >
            บันทึก
          </button>
          <button
            @click="restartGame"
            class="px-6 py-3 bg-blue-600 text-white rounded-lg font-bold text-lg w-1/2"
          >
            เริ่มใหม่
          </button>
        </div>
      </div>
    </div>

    <!-- สถิติผู้เล่น TOP 10 แสดงใต้กล่องเกม -->
    <div class="w-full max-w-lg mx-auto mt-6 bg-white bg-opacity-80 rounded-xl shadow p-4">
      <div class="flex items-center justify-between">
        <h4 class="text-lg font-bold text-blue-700">สถิติผู้เล่น TOP 10</h4>
        <button
          @click="loadScores"
          class="text-xs px-3 py-1 border rounded hover:bg-gray-50"
          title="รีเฟรชสถิติ"
        >
          รีเฟรช
        </button>
      </div>

      <div v-if="savedScores.length === 0" class="text-gray-500 text-sm mt-2">
        ยังไม่มีสถิติ แสดงเมื่อมีการบันทึกคะแนน
      </div>

      <div v-else class="mt-2 space-y-1">
        <div
          v-for="(item, idx) in savedScores"
          :key="item.name + '_' + item.score + '_' + idx"
          class="text-sm text-gray-700 flex items-center justify-between border-b last:border-b-0 py-1"
        >
          <div class="flex items-center gap-2">
            <span class="w-6 text-center">
              <template v-if="idx === 0">🥇</template>
              <template v-else-if="idx === 1">🥈</template>
              <template v-else-if="idx === 2">🥉</template>
              <template v-else> {{ idx + 1 }}.</template>
            </span>
            <span class="font-medium truncate max-w-[9rem] sm:max-w-[12rem]" :title="item.name">
              {{ item.name }}
            </span>
          </div>
          <div class="tabular-nums font-semibold">{{ item.score }} คะแนน</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import catwalk from '../assets/images/catwalk.png';
import catwalk2 from '../assets/images/catwalk2.png';
const loading = ref(true);
const catwalkImages = [catwalk, catwalk2];
const catwalkIndex = ref(0);
let catwalkInterval: number | undefined;
onMounted(() => {
  catwalkInterval = setInterval(() => {
    catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length;
  }, 200);
  setTimeout(() => {
    loading.value = false;
  }, 800);
});
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
const changeCount = ref(0);
const maxChange = 5;
const playerName = ref('');
const savedScores = ref<{ name: string; score: number }[]>([]);
const finalScore = ref(0);

// หมายเหตุ: ในเบราว์เซอร์ setInterval คืนค่า number
let intervalId: number | undefined;

async function fetchQuiz(isAuto = false) {
  if (changeCount.value >= maxChange && !isAuto) return;

  const res = await api.get('/api/quiz');
  quizId.value = res.data.id;
  hints.value = res.data.hints;

  result.value = null;
  guess.value = '';
  hint1.value = '';
  hint2.value = '';
  hintCount.value = 0;
  timer.value = 60;
  showModal.value = false;

  if (intervalId) clearInterval(intervalId as number);
  intervalId = window.setInterval(() => {
    if (timer.value > 0) timer.value--;
    if (timer.value === 0) {
      clearInterval(intervalId as number);
      result.value = false;
      finalScore.value = score.value; // เก็บคะแนนสุดท้ายไว้ในโมดัล
      showModal.value = true;
      // อย่า reset score ที่นี่ เพื่อไม่ให้คะแนนหายก่อนบันทึก
    }
  }, 1000);

  if (!isAuto) changeCount.value++;
}

async function checkAnswer() {
  if (!quizId.value) return;
  const res = await api.post('/api/quiz/check', { id: quizId.value, guess: guess.value.trim() });
  result.value = res.data.correct;
}

watch(result, (val, oldVal) => {
  if (val === true && oldVal !== true) {
    score.value++;
    setTimeout(() => fetchQuiz(true), 1200);
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

async function restartGame() {
  score.value = 0;
  finalScore.value = 0;
  changeCount.value = 0;
  showModal.value = false;
  await fetchQuiz();
  await loadScores();
}

async function saveScore() {
  if (playerName.value.trim()) {
    const scoreToSave = finalScore.value || score.value;
    await api.post('/api/scores', { name: playerName.value.trim(), score: scoreToSave });
    await loadScores();

    // รีเซ็ตสำหรับเกมใหม่
    score.value = 0;
    finalScore.value = 0;
    changeCount.value = 0;
    showModal.value = false;
    playerName.value = '';
    fetchQuiz();
  }
}

async function loadScores() {
  const res = await api.get('/api/scores');
  // คาดว่า backend คืนมาเป็น array ที่ sort แล้ว (สูง -> ต่ำ)
  // ถ้ายังไม่ sort ให้จัดการที่นี่:
  savedScores.value = (res.data || []).sort((a: any, b: any) => b.score - a.score).slice(0, 10);
}

const goBack = () => router.back();

onMounted(() => {
  fetchQuiz();
  loadScores();
});
</script>
