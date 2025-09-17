import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useNetworkStore = defineStore('network', () => {
  const pendingCount = ref(0)
  const oldestPendingStartedAt = ref<number | null>(null)
  const lastError = ref<string | null>(null)
  const healthOk = ref(false)

  function onRequestStart() {
    pendingCount.value += 1
    if (!oldestPendingStartedAt.value) oldestPendingStartedAt.value = Date.now()
  }

  function onRequestEnd() {
    pendingCount.value = Math.max(0, pendingCount.value - 1)
    if (pendingCount.value === 0) oldestPendingStartedAt.value = null
  }

  function setHealth(ok: boolean) {
    healthOk.value = ok
  }

  function setError(msg: string | null) {
    lastError.value = msg
  }

  const hasPending = computed(() => pendingCount.value > 0)
  const isStalled = computed(() => {
    // ถ้าค้างเกิน 20 วินาทีถือว่า stall แต่เรายังคง “ไม่ปิด” loading
    if (!oldestPendingStartedAt.value) return false
    return Date.now() - oldestPendingStartedAt.value > 20_000
  })

  return {
    pendingCount,
    oldestPendingStartedAt,
    lastError,
    healthOk,
    hasPending,
    isStalled,
    onRequestStart,
    onRequestEnd,
    setHealth,
    setError,
  }
})
