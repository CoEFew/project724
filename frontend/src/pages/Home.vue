<template>
  <div class="p-8 relative min-h-screen overflow-hidden">
    <transition name="float-text">
      <h1
        v-if="showFloatingText"
        :style="floatingTextStyle"
        class="text-2xl font-bold mb-8 text-center text-pink-500 fixed z-40 pointer-events-auto select-none cursor-pointer"
        @click="handleFloatingTextClick"
      >ลองเลือกดูซิจ๊ะ!</h1>
    </transition>
    <div v-if="showTooltip" :style="tooltipStyle" class="fixed z-50 bg-black text-white px-4 py-2 rounded shadow-lg text-lg">
      อย่ามาจับฉันนะ
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
      <div v-for="folder in folders" :key="folder.name" class="flex flex-col items-center bg-white rounded-lg shadow-md p-6 hover:bg-blue-50 transition cursor-pointer"
        @click="goToFolder(folder.name)">
        <svg class="w-16 h-16 text-yellow-400 mb-2" fill="currentColor" viewBox="0 0 24 24">
          <path d="M10 4H2v16h20V6H12l-2-2z" />
        </svg>
        <span class="text-lg font-semibold text-gray-700">{{ folder.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ref, onMounted } from 'vue';
import BaseButton from '../components/BaseButton.vue';

const folders = [
  { name: 'ANIMAL' },
  { name: 'Pictures' },
  { name: 'Music' },
  { name: 'Videos' }
];

const router = useRouter();
const goToFolder = (name: string) => {
  if (name === 'ANIMAL') {
    router.push({ name: 'DocumentsPage' });
  } else if (name === 'Pictures') {
    router.push({ name: 'PicturesPage' });
  } else if (name === 'Music') {
    router.push({ name: 'MusicPage' });
  } else if (name === 'Videos') {
    router.push({ name: 'VideosPage' });
  }
};

const showFloatingText = ref(true);
const floatingTextStyle = ref({
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  transition: 'top 1.2s linear, left 1.2s linear',
});
const showTooltip = ref(false);
const tooltipStyle = ref({
  top: '0px',
  left: '0px',
});
function randomPosition() {
  const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0);
  const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0);
  const top = Math.random() * (vh - 60);
  const left = Math.random() * (vw - 300);
  floatingTextStyle.value = {
    top: `${top}px`,
    left: `${left}px`,
    transform: 'translate(0, 0)',
    transition: 'top 1.2s linear, left 1.2s linear',
  };
  // อัปเดตตำแหน่ง tooltip ให้ตามข้อความ
  tooltipStyle.value = {
    top: `${top + 40}px`,
    left: `${left + 80}px`,
  };
}
let floatInterval: number | undefined;
onMounted(() => {
  randomPosition();
  floatInterval = setInterval(randomPosition, 1500);
});
function handleFloatingTextClick() {
  showTooltip.value = true;
  setTimeout(() => {
    showTooltip.value = false;
  }, 5000);
}
</script>

<style scoped>
.float-text-enter-active, .float-text-leave-active {
  transition: opacity 0.5s;
}
.float-text-enter-from, .float-text-leave-to {
  opacity: 0;
}
.float-text-enter-to, .float-text-leave-from {
  opacity: 1;
}
</style>
