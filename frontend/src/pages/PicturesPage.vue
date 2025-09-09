<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100 relative">
    <div v-if="loading" class="fixed inset-0 bg-white bg-opacity-70 flex items-center justify-center z-50">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-lg text-blue-700 font-semibold">กำลังโหลด...</span>
      </div>
    </div>
    <div class="bg-white rounded-xl shadow-lg p-8 flex flex-col items-center">
      <h2 class="text-3xl font-bold mb-4 text-center text-blue-600">Pictures</h2>
      <button @click="showModal = true" class="px-6 py-3 bg-yellow-500 text-white rounded-lg font-bold text-lg">เล่นเกมส์</button>
    </div>
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl shadow-2xl p-8 flex flex-col items-center">
        <h3 class="text-2xl font-bold text-gray-800 mb-4">Coming Soon</h3>
        <p class="mb-6 text-lg text-gray-700">ฟีเจอร์นี้จะเปิดให้เล่นเร็วๆนี้</p>
         <button @click="closeModal" class="px-6 py-3 bg-blue-600 text-white rounded-lg font-bold text-lg">ปิด</button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import catwalk from '../assets/images/catwalk.png';
import catwalk2 from '../assets/images/catwalk2.png';
import { ref, onMounted } from 'vue';
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

import { useRouter } from 'vue-router';
const showModal = ref(false);
const router = useRouter();
function closeModal() {
  showModal.value = false;
  router.push('/');
}
</script>
