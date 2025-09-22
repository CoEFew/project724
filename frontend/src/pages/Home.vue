<template>
  <div class="p-8 relative min-h-screen overflow-visible uppercase">
    <!-- Loading overlay -->
    <LoadingOverlay :loading="loading" message="•Few Surasak•" />

    <!-- ผึ้งบินแบบ sine + draggable -->
    <transition name="float-text">
      <div v-if="showFloatingText" :style="beeStyle"
        class="mb-8 text-center fixed z-40 pointer-events-auto select-none cursor-grab active:cursor-grabbing"
        @pointerdown="onBeePointerDown" @click="handleFloatingTextClick">
        <div class="relative w-32 h-32 mx-auto">
          <img :src="bee" alt="bee A" class="sprite" draggable="false" :class="beeFrontIsA ? 'front' : 'back'" />
          <img :src="bee2" alt="bee B" class="sprite" draggable="false" :class="beeFrontIsA ? 'back' : 'front'" />
        </div>

        <span v-if="showTooltip"
          class="absolute left-1/2 top-full -translate-x-1/2 mt-2 bg-black text-white px-4 py-2 rounded shadow-lg text-lg whitespace-nowrap"
          :style="{ left: '50%', top: '100%' }">
          อย่ามาจับฉันนะ
        </span>
      </div>
    </transition>

    <!-- Catwalk (เดินสมูทด้วย rAF + translate3d) -->
    <div class="catwalk-container">
      <div v-if="showCat" class="fixed left-0 bottom-10 w-screen h-20 pointer-events-none z-30">
        <div :style="catwalkWrapperStyle" class="pointer-events-auto cursor-pointer will-change-transform"
          @click="handleCatClick">
          <img :src="catwalk" class="sprite" alt="catwalk A" :class="catFrontIsA ? 'front' : 'back'" />
          <img :src="catwalk2" class="sprite" alt="catwalk B" :class="catFrontIsA ? 'back' : 'front'" />
        </div>

        <span v-if="showCatTooltip" :style="catTooltipStyle">เมี๊ยววว</span>
      </div>
    </div>

    <!-- เนื้อหา -->
    <div class="w-full flex justify-center mb-10">
      <h1 class="text-3xl font-bold text-blue-700 uppercase">PETTEXT</h1>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">
      <div v-for="folder in folders" :key="folder.name"
        class="relative overflow-visible flex flex-col items-center bg-white rounded-lg shadow-md p-6 hover:bg-blue-50 transition cursor-pointer"
        @click="goToFolder(folder.name)"
        @mouseenter="folder.name === 'DogPuzzle' ? animalImg = dog2 : folder.name === 'CatText' ? pictureImg = picture2 : folder.name === 'CatGame' ? catImg = picture2 : folder.name === 'PolaJigsaw' ? polabearImg = polabear2 : folder.name === 'Otterfeedback' ? otterImg = otter2 : null"
        @mouseleave="folder.name === 'DogPuzzle' ? animalImg = dog : folder.name === 'CatText' ? pictureImg = picture : folder.name === 'CatGame' ? catImg = picture : folder.name === 'PolaJigsaw' ? polabearImg = polabear : folder.name === 'Otterfeedback' ? otterImg = otter : null">
        <template v-if="folder.name === 'DogPuzzle'">
          <div class="relative w-32 h-10 sm:w-40 sm:h-12 md:w-48 md:h-14 lg:w-56 lg:h-16 mb-2">
            <img :src="animalImg" alt="DogPuzzle" class="absolute left-1/2 -top-20 -translate-x-1/2
                     w-36 h-36 sm:w-44 sm:h-44 md:w-52 md:h-52 lg:w-60 lg:h-60
                     object-contain drop-shadow-lg pointer-events-none z-10" />
          </div>
        </template>
        <template v-else-if="folder.name === 'PolaJigsaw'">
          <div class="relative w-32 h-10 sm:w-40 sm:h-12 md:w-48 md:h-14 lg:w-56 lg:h-16 mb-2">
            <img :src="polabearImg" alt="PolaJigsaw" class="absolute left-1/2 -top-20 -translate-x-1/2
                     w-36 h-36 sm:w-44 sm:h-44 md:w-52 md:h-52 lg:w-60 lg:h-60
                     object-contain drop-shadow-lg pointer-events-none z-10" />
          </div>
        </template>
        <!-- <template v-else-if="folder.name === 'CatText'">
          <div class="relative w-32 h-10 sm:w-40 sm:h-12 md:w-48 md:h-14 lg:w-56 lg:h-16 mb-2">
            <img
              :src="pictureImg"
              alt="CatText"
              class="absolute left-1/2 -top-20 -translate-x-1/2
                     w-36 h-36 sm:w-44 sm:h-44 md:w-52 md:h-52 lg:w-60 lg:h-60
                     object-contain drop-shadow-lg pointer-events-none z-10"
            />
          </div>
        </template> -->
        <template v-else-if="folder.name === 'CatGame'">
          <div class="relative w-32 h-10 sm:w-40 sm:h-12 md:w-48 md:h-14 lg:w-56 lg:h-16 mb-2">
            <img :src="catImg" alt="CatGame" class="absolute left-1/2 -top-20 -translate-x-1/2
                     w-36 h-36 sm:w-44 sm:h-44 md:w-52 md:h-52 lg:w-60 lg:h-60
                     object-contain drop-shadow-lg pointer-events-none z-10" />
          </div>
        </template>
        <template v-else-if="folder.name === 'Otterfeedback'">
          <div class="relative w-32 h-10 sm:w-40 sm:h-12 md:w-48 md:h-14 lg:w-56 lg:h-16 mb-2">
            <img :src="otterImg" alt="Otterfeedback" class="absolute left-1/2 -top-20 -translate-x-1/2
                     w-36 h-36 sm:w-44 sm:h-44 md:w-52 md:h-52 lg:w-60 lg:h-60
                     object-contain drop-shadow-lg pointer-events-none z-10" />
          </div>
        </template>
        <template v-else>
          <svg class="w-16 h-16 text-yellow-400 mb-2" fill="currentColor" viewBox="0 0 24 24">
            <path d="M10 4H2v16h20V6H12l-2-2z" />
          </svg>
        </template>

        <span class="text-lg font-semibold text-gray-700 block md:hidden uppercase">{{ folder.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed, type CSSProperties } from 'vue'
import { useRouter } from 'vue-router'
import { waitApiReadyAndLoadInitial } from '../composables/useApiReadiness'
import { useNetworkStore } from '../store/useNetworkStore'

import bee from '../assets/images/bee.png'
import bee2 from '../assets/images/bee2.png'
import catwalk from '../assets/images/catwalk.png'
import catwalk2 from '../assets/images/catwalk2.png'
import LoadingOverlay from '../components/LoadingOverlay.vue'
import dog from '../assets/images/dog.png'
import dog2 from '../assets/images/dog2.png'
import picture from '../assets/images/cat2.png'
import picture2 from '../assets/images/cat4.png'
import polabear from '../assets/images/polabear.png'
import polabear2 from '../assets/images/polabear3.png'
import otter from '../assets/images/otter.png'
import otter2 from '../assets/images/otter2.png'

const loading = ref(true)

const folders = [{ name: 'DogPuzzle' }, { name: 'PolaJigsaw' }, { name: 'CatGame' }, { name: 'Otterfeedback' }]

const animalImg = ref(dog)
const pictureImg = ref(picture)
const catImg = ref(picture)
const polabearImg = ref(polabear)
const otterImg = ref(otter)
const router = useRouter()
const goToFolder = async (name: string) => {
  try {
    if (name === 'DogPuzzle') {
      await router.push({ name: 'DogAll' })
      return
    }
    if (name === 'CatText') await router.push({ name: 'CatText' })
    else if (name === 'CatGame') await router.push({ name: 'CatGame' })
    else if (name === 'PolaJigsaw') await router.push({ name: 'JigsawPage' })
    else if (name === 'Otterfeedback') await router.push({ name: 'FeedbackPage' })
  } catch (e) {
    // ป้องกัน unhandled promise จาก NavigationDuplicated ฯลฯ
    console.warn('navigation error', e)
  }
}

const showFloatingText = ref(true)
const showTooltip = ref(false)

// Network store and health check
const net = useNetworkStore()

/* ---------- Utilities ---------- */
function preload(srcs: string[]) { srcs.forEach(s => { const i = new Image(); i.src = s }) }

/* ---------- Bee motion (rAF + translate3d) ---------- */
const margin = 16
const beeSize = 128
const beeX = ref(0)
const beeY = ref(0)
const baseY = ref(0)
const amplitude = ref(60)
const wavelength = ref(280)
const phase = ref(0)
const dir = ref(1)
const speedX = ref(180)
let rafBee: number | undefined
let lastBeeT: number | null = null
const beeFrontIsA = ref(true)
let beeSwapAcc = 0
const BEE_SWAP_EVERY = 0.05

const beeStyle = computed<CSSProperties>(() => {
  const flip = dir.value < 0 ? ' scaleX(-1)' : ''
  return {
    position: 'fixed',
    top: '0px',
    left: '0px',
    transform: `translate3d(${Math.round(beeX.value)}px, ${Math.round(beeY.value)}px, 0)` + flip,
    willChange: 'transform',
    transition: isDragging.value ? 'none' : 'transform 60ms linear',
    touchAction: 'none',
    zIndex: 40,
  }
})

function placeBeeCenter() {
  const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
  const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0)
  beeX.value = Math.min(Math.max(vw * 0.15, margin), vw - beeSize - margin)
  baseY.value = Math.min(Math.max(vh * 0.28, margin), vh - beeSize - margin)
  beeY.value = baseY.value
  phase.value = -(2 * Math.PI * beeX.value) / wavelength.value
}

function animateBee(t: number) {
  if (lastBeeT === null) lastBeeT = t
  const dt = (t - lastBeeT) / 1000
  lastBeeT = t

  if (!isDragging.value) {
    const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
    const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0)
    beeX.value += dir.value * speedX.value * dt

    const minX = margin
    const maxX = vw - beeSize - margin
    if (beeX.value <= minX) { beeX.value = minX; dir.value = 1; randomizeWave() }
    else if (beeX.value >= maxX) { beeX.value = maxX; dir.value = -1; randomizeWave() }

    const y = baseY.value + amplitude.value * Math.sin((2 * Math.PI * beeX.value) / wavelength.value + phase.value)
    const minY = margin
    const maxY = vh - beeSize - margin
    beeY.value = Math.min(Math.max(y, minY), maxY)
  }

  beeSwapAcc += dt
  if (beeSwapAcc >= BEE_SWAP_EVERY) { beeFrontIsA.value = !beeFrontIsA.value; beeSwapAcc = 0 }

  rafBee = requestAnimationFrame(animateBee)
}

function randomizeWave() {
  const amp = 50 + Math.random() * 40
  const wave = 240 + Math.random() * 160
  amplitude.value = amp
  wavelength.value = wave
  phase.value = (Math.asin(clamp((beeY.value - baseY.value) / amplitude.value, -1, 1)) - (2 * Math.PI * beeX.value) / wavelength.value)
}
function clamp(n: number, a: number, b: number) { return Math.min(Math.max(n, a), b) }

/* drag */
const isDragging = ref(false)
let dragOffsetX = 0, dragOffsetY = 0
function onBeePointerDown(e: PointerEvent) {
  const vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)
  const vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0)
  isDragging.value = true
    ; (e.currentTarget as HTMLElement).setPointerCapture?.(e.pointerId)
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  dragOffsetX = e.clientX - rect.left
  dragOffsetY = e.clientY - rect.top

  const onMove = (ev: PointerEvent) => {
    const nx = clamp(ev.clientX - dragOffsetX, margin, vw - beeSize - margin)
    const ny = clamp(ev.clientY - dragOffsetY, margin, vh - beeSize - margin)
    dir.value = nx > beeX.value ? 1 : nx < beeX.value ? -1 : dir.value
    beeX.value = nx; beeY.value = ny
  }
  const onUp = () => {
    isDragging.value = false
    window.removeEventListener('pointermove', onMove)
    window.removeEventListener('pointerup', onUp)
    window.removeEventListener('pointercancel', onUp)
    baseY.value = beeY.value
    phase.value = -(2 * Math.PI * beeX.value) / wavelength.value
  }
  window.addEventListener('pointermove', onMove)
  window.addEventListener('pointerup', onUp)
  window.addEventListener('pointercancel', onUp)
}

function handleFloatingTextClick() {
  showTooltip.value = true
  setTimeout(() => (showTooltip.value = false), 5000)
}

/* ---------- Catwalk (rAF + translate3d) ---------- */
const CAT_W = 200
const catX = ref(window.innerWidth + 100)
const catSpeed = ref(220)
const catFrontIsA = ref(true)
let rafCat: number | undefined
let lastCatT: number | null = null
let catSwapAcc = 0
const CAT_SWAP_EVERY = 0.20
let catPauseUntil = 0

const catwalkWrapperStyle = computed<CSSProperties>(() => ({
  position: 'fixed',
  bottom: '40px',
  width: `${CAT_W}px`,
  height: `${CAT_W}px`,
  transform: `translate3d(${Math.round(catX.value)}px, 0, 0)`,
  willChange: 'transform',
  zIndex: 31
}))
const catTooltipStyle = computed<CSSProperties>(() => ({
  position: 'fixed',
  left: `${Math.round(catX.value)}px`,
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
}))
const showCat = ref(true)
const showCatTooltip = ref(false)

function animateCat(t: number) {
  if (lastCatT === null) lastCatT = t
  const dt = (t - lastCatT) / 1000
  lastCatT = t

  if (t < catPauseUntil) { rafCat = requestAnimationFrame(animateCat); return }

  catX.value -= catSpeed.value * dt
  catSwapAcc += dt
  if (catSwapAcc >= CAT_SWAP_EVERY) { catFrontIsA.value = !catFrontIsA.value; catSwapAcc = 0 }

  if (catX.value < -CAT_W) {
    catX.value = window.innerWidth + 100
    catPauseUntil = t + 400 // ms
  }
  rafCat = requestAnimationFrame(animateCat)
}

function handleCatClick() {
  showCatTooltip.value = true
  setTimeout(() => (showCatTooltip.value = false), 5000)
}

/* ---------- Loading state ---------- */

/* ---------- lifecycle ---------- */
onMounted(async () => {
  document.title = 'PETTEXT - Home'
  preload([bee, bee2, dog, dog2, picture, picture2, polabear, polabear2, otter, otter2])

  // Initialize loading state

  placeBeeCenter()
  rafBee = requestAnimationFrame(animateBee)
  rafCat = requestAnimationFrame(animateCat)

  // Wait for API to be ready
  const { healthOk, initialOk } = await waitApiReadyAndLoadInitial()
  if (!healthOk || !initialOk) {
    loading.value = false
    return
  }

  setTimeout(() => (loading.value = false), 800)

  window.addEventListener('resize', () => {
    if (catX.value > window.innerWidth) catX.value = window.innerWidth + 100
    placeBeeCenter()
  })
})

onBeforeUnmount(() => {
  if (rafBee) cancelAnimationFrame(rafBee)
  if (rafCat) cancelAnimationFrame(rafCat)
})
</script>

<style scoped>
/* transition สำหรับ v-transition เท่านั้น */
.float-text-enter-active,
.float-text-leave-active {
  transition: opacity 0.5s;
}

.float-text-enter-from,
.float-text-leave-to {
  opacity: 0;
}

.float-text-enter-to,
.float-text-leave-from {
  opacity: 1;
}

/* ซ้อนรูป (โหลดครั้งเดียว ไม่เปลี่ยน src) */
.sprite {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.front {
  z-index: 2;
  visibility: visible;
}

.back {
  z-index: 1;
  visibility: hidden;
}

/* พื้นที่ catwalk */
.catwalk-container {
  pointer-events: none;
  z-index: 30;
}

/* ลดงาน layout ของภาพ */
.will-change-transform {
  will-change: transform;
}
</style>
