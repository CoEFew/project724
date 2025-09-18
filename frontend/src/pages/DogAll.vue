<!-- src/pages/DogAll.vue -->
<template>
  <div class="min-h-screen px-4 py-10 bg-gradient-to-b from-slate-900 to-slate-950 text-slate-100">
    <div class="max-w-5xl mx-auto">
      <header class="mb-8 flex items-center justify-between">
        <h1 class="text-2xl md:text-3xl font-extrabold tracking-wide">
          DOG • PUZZLE <span class="text-indigo-300">เลือกโหมดการเล่น</span>
        </h1>
        <button
          @click="goBack"
          class="px-3 py-2 rounded-xl border border-white/10 bg-white/5 hover:bg-white/10 transition"
          aria-label="กลับหน้าหลัก"
        >
          ย้อนกลับ
        </button>
      </header>

      <!-- เพิ่มเป็น 3 คอลัมน์สำหรับหน้าจอกว้าง -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <!-- โหมดเดี่ยว (มีเวลา) -->
        <div class="rounded-2xl border border-white/10 bg-white/5 p-6 flex flex-col">
          <h2 class="text-xl font-bold mb-2">โหมดเดี่ยว</h2>
          <p class="text-sm text-slate-300 mb-4">เล่นคนเดียว มีจับเวลาต่อข้อ</p>
          <ul class="text-sm text-slate-300 list-disc ml-5 mb-6 space-y-1">
            <li>ระบบใบ้คำ</li>
            <li>บันทึกสกอร์</li>
          </ul>
          <button
            @click="goSolo"
            class="mt-auto inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 hover:bg-indigo-400 text-white shadow"
          >
            เล่นโหมดเดี่ยว
            <span>→</span>
          </button>
        </div>

            <!-- ✅ โหมดเดี่ยว (ไม่จับเวลา) -->
        <div class="rounded-2xl border border-white/10 bg-white/5 p-6 flex flex-col">
          <h2 class="text-xl font-bold mb-2">โหมดเดี่ยว (ไม่จับเวลา)</h2>
          <p class="text-sm text-slate-300 mb-4">เล่นสบาย ๆ ไม่ต้องกังวลเรื่องเวลา</p>
          <ul class="text-sm text-slate-300 list-disc ml-5 mb-6 space-y-1">
            <li>ไม่มีหมดเวลา</li>
            <li>ข้ามข้อ/ขอใบ้ได้ตามต้องการ</li>
          </ul>
          <button
            @click="goSoloNoTimer"
            class="mt-auto inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl font-semibold bg-emerald-500 hover:bg-emerald-400 text-white shadow"
          >
            เล่นโหมดเดี่ยว (ไม่จับเวลา)
            <span>→</span>
          </button>
        </div>

        <!-- โหมดปาร์ตี้ -->
        <div class="rounded-2xl border border-white/10 bg-white/5 p-6 flex flex-col">
          <h2 class="text-xl font-bold mb-2">โหมดปาร์ตี้</h2>
          <p class="text-sm text-slate-300 mb-4">เล่นพร้อมกันทั้งห้อง แข่งตอบให้ไว</p>
          <ul class="text-sm text-slate-300 list-disc ml-5 mb-6 space-y-1">
            <li>เจ้าของห้องเริ่มเกมได้ทันที</li>
            <li>นับเวลาต่อข้อ + สรุปสกอร์</li>
          </ul>
          <button
            @click="goParty"
            class="mt-auto inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl font-semibold bg-fuchsia-500 hover:bg-fuchsia-400 text-white shadow"
          >
            เล่นโหมดปาร์ตี้
            <span>→</span>
          </button>
        </div>

      </div>

      <transition name="toast" tag="div">
        <div v-if="err" class="mt-6 rounded-xl px-4 py-3 text-sm shadow border border-rose-300/30 bg-rose-400/10" role="alert">
          {{ err }}
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const err = ref<string|null>(null)

async function safePush(to: any) {
  err.value = null
  try { await router.push(to) }
  catch (e:any) { err.value = e?.message || 'ไม่สามารถนำทางได้' }
}

const goBack = () => router.back()
const goSolo        = () => safePush({ name: 'DocumentsPage' }) // เดี่ยว (มีเวลา)
const goParty       = () => safePush({ name: 'DocumentsPageAlls' }) // ปาร์ตี้
// ✅ เดี่ยว (ไม่จับเวลา): ใช้ query แทนการแยก route
const goSoloNoTimer = () => safePush({ name: 'DocumentsPage', query: { noTimer: '1' } })
</script>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: opacity .25s ease, transform .25s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(6px); }
</style>
