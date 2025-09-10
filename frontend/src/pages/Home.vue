<template>
  <div class="p-8 relative min-h-screen overflow-visible uppercase">
    <div v-if="loading" class="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50">
      <div class="flex flex-col items-center">
        <img :src="catwalkImages[catwalkIndex]" alt="loading cat" class="h-24 w-24 mb-4 animate-bounce" />
        <span class="text-lg text-blue-700 font-semibold">กำลังโหลด...</span>
      </div>
    </div>

    <!-- ผึ้งบินแบบ sine + draggable -->
    <transition name="float-text">
      <div
        v-if="showFloatingText"
        :style="beeStyle"
        class="mb-8 text-center fixed z-40 pointer-events-auto select-none cursor-grab active:cursor-grabbing"
        @pointerdown="onBeePointerDown"
        @click="handleFloatingTextClick"
      >
        <img
          :src="beeImages[beeIndex]"
          alt="ลองเลือกดูซิจ๊ะ!"
          class="w-32 h-32 mx-auto"
          draggable="false"
        />
        <span
          v-if="showTooltip"
          class="absolute left-1/2 top-full -translate-x-1/2 mt-2 bg-black text-white px-4 py-2 rounded shadow-lg text-lg whitespace-nowrap"
          :style="{left: '50%', top: '100%'}"
        >
          อย่ามาจับฉันนะ
        </span>
      </div>
    </transition>

    <!-- Catwalk -->
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

    <div class="w-full flex justify-center mb-6">
      <h1 class="text-3xl font-bold text-blue-700 uppercase">PETTEXT</h1>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
      <div
        v-for="folder in folders"
        :key="folder.name"
        class="relative overflow-visible flex flex-col items-center bg-white rounded-lg shadow-md p-6 hover:bg-blue-50 transition cursor-pointer"
        @click="goToFolder(folder.name)"
        @mouseenter="folder.name === 'ANIMAL' ? animalImg = dog2 : folder.name === 'CatText' ? pictureImg = picture2 : null"
        @mouseleave="folder.name === 'ANIMAL' ? animalImg = dog : folder.name === 'CatText' ? pictureImg = picture : null"
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
        <template v-else-if="folder.name === 'CatText'">
          <div class="relative w-32 h-10 sm:w-40 sm:h-12 md:w-48 md:h-14 lg:w-56 lg:h-16 mb-2">
            <img
              :src="pictureImg"
              alt="CatText"
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

        <span class="text-lg font-semibold text-gray-700 block md:hidden lg:block uppercase">{{ folder.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const loading = ref(true)
import { useRouter } from 'vue-router'
import { ref, onMounted, onBeforeUnmount, computed, type CSSProperties } from 'vue'
import BaseButton from '../components/BaseButton.vue'
import bee from '../assets/images/bee.png'
import bee2 from '../assets/images/bee2.png'
import catwalk from '../assets/images/catwalk.png'
import catwalk2 from '../assets/images/catwalk2.png'
import dog from '../assets/images/dog.png'
import dog2 from '../assets/images/dog2.png'
import picture from '../assets/images/cat4.png'
import picture2 from '../assets/images/cat2.png'

const folders = [
  { name: 'ANIMAL' },
  { name: 'CatText' },
  { name: 'bird' },
  { name: 'fish' }
]

const beeImages = [bee, bee2]
const beeIndex = ref(0)
const animalImg = ref(dog)
const pictureImg = ref(picture)
const router = useRouter()
const goToFolder = (name: string) => {
  if (name === 'ANIMAL') {
    router.push({ name: 'DocumentsPage' })
  } else if (name === 'CatText') {
    router.push({ name: 'CatText' })
  } else if (name === 'bird') {
    router.push({ name: 'MusicPage' })
  } else if (name === 'fish') {
    router.push({ name: 'VideosPage' })
  }
}

const showFloatingText = ref(true)
const showTooltip = ref(false)

// --------- Bee motion (sine) + drag ----------
const margin = 16
const beeSize = 128 // w-32 h-32
const beeX = ref(0)
const beeY = ref(0)
const baseY = ref(0)
const amplitude = ref(60)      // px
const wavelength = ref(280)    // px
const phase = ref(0)           // rad
const dir = ref(1)             // 1 → right, -1 → left
const speedX = ref(180)        // px/s
let rafId: number | undefined
let lastT: number | null = null

// flap wings via the existing catwalk interval; no change needed here

// style (use transform for smooth animation)
const beeStyle = computed<CSSProperties>(() => {
  // flip horizontally when moving left
  const flip = dir.value < 0 ? ' scaleX(-1)' : ''
  return {
    position: 'fixed',
    top: '0px',
    left: '0px',
    transform: `translate3d(${Math.round(beeX.value)}px, ${Math.round(beeY.value)}px, 0)` + flip,
    willChange: 'transform',
    transition: isDragging.value ? 'none' : 'transform 60ms linear',
    touchAction: 'none', // allow drag on touch
    zIndex: 40,
  }
})

// init position & responsive baseline
function placeBeeCenter() {
  const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
  const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0)
  beeX.value = Math.min(Math.max(vw * 0.15, margin), vw - beeSize - margin)
  baseY.value = Math.min(Math.max(vh * 0.28, margin), vh - beeSize - margin)
  beeY.value = baseY.value
  // set phase so that sin(...) = 0 at current x
  phase.value = -(2 * Math.PI * beeX.value) / wavelength.value
}

// main animation loop
function animateBee(t: number) {
  if (lastT === null) lastT = t
  const dt = (t - lastT) / 1000
  lastT = t

  if (!isDragging.value) {
    const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
    const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0)

    // move in X
    beeX.value += dir.value * speedX.value * dt

    // bounce at edges
    const minX = margin
    const maxX = vw - beeSize - margin
    if (beeX.value <= minX) {
      beeX.value = minX
      dir.value = 1
      randomizeWave()
    } else if (beeX.value >= maxX) {
      beeX.value = maxX
      dir.value = -1
      randomizeWave()
    }

    // sine Y based on X (S-curve)
    const y = baseY.value + amplitude.value * Math.sin((2 * Math.PI * beeX.value) / wavelength.value + phase.value)
    const minY = margin
    const maxY = vh - beeSize - margin
    beeY.value = Math.min(Math.max(y, minY), maxY)
  }

  rafId = requestAnimationFrame(animateBee)
}

function randomizeWave() {
  // เล็กน้อยพอให้ดูมีชีวิตชีวา (แต่ไม่หลุดธีม)
  const amp = 50 + Math.random() * 40        // 50–90
  const wave = 240 + Math.random() * 160     // 240–400
  amplitude.value = amp
  wavelength.value = wave
  // รักษาความต่อเนื่อง: ปรับ phase เพื่อให้ตำแหน่งปัจจุบันยังไม่กระโดด
  phase.value = (Math.asin(clamp((beeY.value - baseY.value) / amplitude.value, -1, 1)) - (2 * Math.PI * beeX.value) / wavelength.value)
}

function clamp(n: number, a: number, b: number) {
  return Math.min(Math.max(n, a), b)
}

// drag handling
const isDragging = ref(false)
let dragOffsetX = 0
let dragOffsetY = 0

function onBeePointerDown(e: PointerEvent) {
  const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
  const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0)

  isDragging.value = true
  ;(e.currentTarget as HTMLElement).setPointerCapture?.(e.pointerId)

  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  dragOffsetX = e.clientX - rect.left
  dragOffsetY = e.clientY - rect.top

  const onMove = (ev: PointerEvent) => {
    const nx = clamp(ev.clientX - dragOffsetX, margin, vw - beeSize - margin)
    const ny = clamp(ev.clientY - dragOffsetY, margin, vh - beeSize - margin)
    // update direction according to movement
    dir.value = nx > beeX.value ? 1 : nx < beeX.value ? -1 : dir.value
    beeX.value = nx
    beeY.value = ny
  }
  const onUp = () => {
    isDragging.value = false
    window.removeEventListener('pointermove', onMove)
    window.removeEventListener('pointerup', onUp)
    window.removeEventListener('pointercancel', onUp)
    // resume sine smoothly from current position:
    // set baseY to current Y and phase so sin(...) = 0 at x
    baseY.value = beeY.value
    phase.value = -(2 * Math.PI * beeX.value) / wavelength.value
  }

  window.addEventListener('pointermove', onMove)
  window.addEventListener('pointerup', onUp)
  window.addEventListener('pointercancel', onUp)
}

// tooltip click
function handleFloatingTextClick() {
  showTooltip.value = true
  setTimeout(() => (showTooltip.value = false), 5000)
}

// --------- Catwalk animation ---------
const catwalkImages = [catwalk, catwalk2]
const catwalkIndex = ref(0)
const catwalkStyle = ref<CSSProperties>({ left: '100vw', bottom: '40px', position: 'fixed' })
let catwalkPos = window.innerWidth + 100
const showCat = ref(true)
const showCatTooltip = ref(false)

function animateCatwalk() {
  catwalkPos -= 16
  if (catwalkPos < -80) {
    showCat.value = false
    setTimeout(() => {
      catwalkPos = window.innerWidth + 100
      catwalkStyle.value = { left: `${catwalkPos}px`, bottom: '40px', position: 'fixed' }
      showCat.value = true
    }, 400)
    return
  }
  catwalkStyle.value = { left: `${catwalkPos}px`, bottom: '40px', position: 'fixed' }
  catwalkIndex.value = (catwalkIndex.value + 1) % catwalkImages.length
}

function handleCatClick() {
  showCatTooltip.value = true
  setTimeout(() => (showCatTooltip.value = false), 5000)
}

let catwalkInterval: number | undefined

// lifecycle
onMounted(() => {
  document.title = 'PETTEXT - Home'
  placeBeeCenter()
  rafId = requestAnimationFrame(animateBee)

  catwalkInterval = setInterval(() => {
    animateCatwalk()
    beeIndex.value = (beeIndex.value + 1) % beeImages.length // ใช้เดิมเป็น wing flap
  }, 200)

  setTimeout(() => (loading.value = false), 800)

  // update baseline & bounds on resize
  window.addEventListener('resize', placeBeeCenter)
})

onBeforeUnmount(() => {
  if (rafId) cancelAnimationFrame(rafId)
  if (catwalkInterval) clearInterval(catwalkInterval as number)
  window.removeEventListener('resize', placeBeeCenter)
})
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

/* keep original catwalk look */
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

/* for the folder image card drop-shadow */
.folder-animal-img {
  position: relative;
  z-index: 2;
  box-shadow: 0 4px 24px rgba(0,0,0,0.12);
  margin-bottom: 0.5rem;
  overflow: visible;
}
</style>
