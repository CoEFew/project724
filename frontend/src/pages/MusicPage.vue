<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="bg-white rounded-xl shadow-lg p-8 flex flex-col items-center">
      <h2 class="text-3xl font-bold mb-4 text-center text-green-600">Music</h2>
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
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const showModal = ref(false);
const router = useRouter();
function closeModal() {
  showModal.value = false;
  router.push('/');
}
</script>
