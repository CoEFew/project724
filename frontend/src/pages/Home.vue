<template>
  <div class="p-8 relative min-h-screen overflow-visible">
    <div v-if="loading" class="fixed inset-0 bg-white bg-opacity-70 flex items-center justify-center z-50">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-lg text-blue-700 font-semibold">กำลังโหลด...</span>
      </div>
    </div>
    <transition name="float-text">
      <div
        v-if="showFloatingText"
        :style="floatingTextStyle"
        class="mb-8 text-center fixed z-40 pointer-events-auto select-none cursor-pointer"
        @click="handleFloatingTextClick"
      >
  <img :src="beeImages[beeIndex]" alt="ลองเลือกดูซิจ๊ะ!" class="w-32 h-32 mx-auto" />
        <span v-if="showTooltip" class="absolute left-1/2 top-full -translate-x-1/2 mt-2 bg-black text-white px-4 py-2 rounded shadow-lg text-lg whitespace-nowrap" :style="{left: '50%', top: '100%'}">
          อย่ามาจับฉันนะ
        </span>
      </div>
    </transition>
    <div class="catwalk-container">
      <div v-if="showCat" style="position:fixed;left:0;bottom:40px;width:100vw;height:80px;pointer-events:none;z-index:30;">
        <img
          :src="catwalkImages[catwalkIndex]"
          class="catwalk cursor-pointer"
          :style="catwalkStyle"
          alt="catwalk"
          @click="handleCatClick"
          style="pointer-events:auto;z-index:31;"
        />
        <span v-if="showCatTooltip"
          :style="{
            position: 'fixed',
            left: catwalkStyle.left,
            bottom: '200px',
            zIndex: 50,
            background: 'black',
            color: 'white',
            padding: '6px 16px',
            borderRadius: '8px',
            boxShadow: '0 2px 8px rgba(0,0,0,0.2)',
            fontSize: '1.1rem',
            whiteSpace: 'nowrap',
            pointerEvents: 'none',
          }"
        >
          เมี๊ยววว
        </span>
      </div>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
    <div
  v-for="folder in folders"
  :key="folder.name"
  class="relative overflow-visible flex flex-col items-center bg-white rounded-lg shadow-md p-6 hover:bg-blue-50 transition cursor-pointer"
  @click="goToFolder(folder.name)"
  @mouseenter="folder.name === 'ANIMAL' ? animalImg = dog2 : folder.name === 'Pictures' ? pictureImg = picture2 : null"
  @mouseleave="folder.name === 'ANIMAL' ? animalImg = dog : folder.name === 'Pictures' ? pictureImg = picture : null"
>

  <template v-if="folder.name === 'ANIMAL'">
    <div class="relative w-32 h-10 sm:w-40 sm:h-12 md:w-48 md:h-14 lg:w-56 lg:h-16 mb-2">
      <img
        :src="animalImg"
        alt="ANIMAL"
        class="absolute left-1/2 -top-20 -translate-x-1/2
               w-36 h-36 sm:w-44 sm:h-44 md:w-52 md:h-52 lg:w-60 lg:h-60
               object-contain drop-shadow-lg pointer-events-none z-10"
      />
    </div>
  </template>
  <template v-else-if="folder.name === 'Pictures'">
    <div class="relative w-32 h-10 sm:w-40 sm:h-12 md:w-48 md:h-14 lg:w-56 lg:h-16 mb-2">
      <img
        :src="pictureImg"
        alt="Pictures"
        class="absolute left-1/2 -top-20 -translate-x-1/2
               w-36 h-36 sm:w-44 sm:h-44 md:w-52 md:h-52 lg:w-60 lg:h-60
               object-contain drop-shadow-lg pointer-events-none z-10"
      />
    </div>
  </template>
  <template v-else>
    <svg class="w-16 h-16 text-yellow-400 mb-2" fill="currentColor" viewBox="0 0 24 24">
      <path d="M10 4H2v16h20V6H12l-2-2z" />
    </svg>
  </template>

  <span class="text-lg font-semibold text-gray-700 block md:hidden lg:block">{{ folder.name }}</span>
</div>

    </div>
  </div>
</template>

<script setup lang="ts">
const loading = ref(true);
import { useRouter } from 'vue-router';
import { ref, onMounted } from 'vue';
import type { CSSProperties } from 'vue';
import BaseButton from '../components/BaseButton.vue';
import bee from '../assets/images/bee.png';
import bee2 from '../assets/images/bee2.png';
import catwalk from '../assets/images/catwalk.png';
import catwalk2 from '../assets/images/catwalk2.png';
import dog from '../assets/images/dog.png';
import dog2 from '../assets/images/dog2.png';
import picture from '../assets/images/dog.png';
import picture2 from '../assets/images/dog2.png';

const folders = [
  { name: 'ANIMAL' },
  { name: 'Pictures' },
  { name: 'Music' },
  { name: 'Videos' }
];

const beeImages = [bee, bee2];
const beeIndex = ref(0);
const animalImg = ref(dog);
const pictureImg = ref(picture);
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

// Catwalk animation (right to left, small size, hide when out of left)
const catwalkImages = [catwalk, catwalk2];
const catwalkIndex = ref(0);
const catwalkStyle = ref<CSSProperties>({ left: '100vw', bottom: '40px', position: 'fixed' });
let catwalkPos = window.innerWidth + 100;
const showCat = ref(true);
const showCatTooltip = ref(false);
function animateCatwalk() {
  catwalkPos -= 16;
  if (catwalkPos < -80) {
    showCat.value = false;
    setTimeout(() => {
      catwalkPos = window.innerWidth + 100;
      catwalkStyle.value = {
        left: `${catwalkPos}px`,
        bottom: '40px',
        position: 'fixed',
      };
      showCat.value = true;
    }, 400);
    return;
  }
  catwalkStyle.value = {
    left: `${catwalkPos}px`,
    bottom: '40px',
    position: 'fixed',
  };
  catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length;
}
function handleCatClick() {
  showCatTooltip.value = true;
  setTimeout(() => {
    showCatTooltip.value = false;
  }, 5000);
}
let catwalkInterval: number | undefined;
onMounted(() => {
  randomPosition();
  floatInterval = setInterval(randomPosition, 1500);
  catwalkInterval = setInterval(() => {
    animateCatwalk();
    beeIndex.value = (beeIndex.value + 1) % beeImages.length;
  }, 200);
  setTimeout(() => {
    loading.value = false;
  }, 800);
});
function handleFloatingTextClick() {
  showTooltip.value = true;
  setTimeout(() => {
    showTooltip.value = false;
  }, 5000);
}
let floatInterval: number | undefined;
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
.catwalk-container {
  position: fixed;
  left: 0;
  bottom: 40px;
  width: 100vw;
  height: 80px;
  pointer-events: none;
  z-index: 30;
}
.catwalk {
  width: 200px;
  height: 200px;
  position: fixed;
  bottom: 40px;
  left: 100vw;
  transition: left 0.3s linear;
}
.folder-animal-img {
  position: relative;
  z-index: 2;
  box-shadow: 0 4px 24px rgba(0,0,0,0.12);
  margin-bottom: 0.5rem;
  /* allow overflow for visual pop */
  overflow: visible;
}
</style>
