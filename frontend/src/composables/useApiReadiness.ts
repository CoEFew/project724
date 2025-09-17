import { useNetworkStore } from '../store/useNetworkStore'
import api from '../services/api'

function sleep(ms: number) {
  return new Promise((r) => setTimeout(r, ms))
}

async function pingHealthWithBackoff(maxAttempts = 6) {
  const net = useNetworkStore()
  let delay = 1000 // 1s
  for (let i = 1; i <= maxAttempts; i++) {
    try {
      const res = await api.get('/health', { timeout: 8000 })
      const ok = res?.status === 200 && (res.data?.status === 'ok' || res.data === 'ok')
      net.setHealth(!!ok)
      if (ok) return true
    } catch {
      net.setHealth(false)
    }
    await sleep(delay)
    delay = Math.min(delay * 2, 15000) // 1s → 2s → 4s … cap 15s
  }
  return false
}

/**
 * initialLoad: เรียก API ที่ต้องมีก่อนเข้าใช้งาน
 * - ถ้ายัง pending → อย่าปิด loading
 * - ปิด loading ได้เมื่อ health ok และ initial APIs สำเร็จ
 */
export async function waitApiReadyAndLoadInitial() {
  const net = useNetworkStore()

  // 1) รอ health ให้กลับมา ok ก่อน (รองรับ DB ตื่นช้า)
  await pingHealthWithBackoff()

  // 2) โหลดข้อมูลตั้งต้น (อย่าลืมปรับให้ตรงกับของคุณ)
  //    ใช้ Promise.allSettled เพื่อไม่ให้ throw รวม แต่เราจะเช็คผลเอง
  const results = await Promise.allSettled([
    api.get('/api/quiz', { params: { level: 1 } }),
    api.get('/api/scores', { params: { limit: 10, gamename: 'DogPuzzle' } }),
  ])

  const allOk = results.every(r => r.status === 'fulfilled')
  if (!allOk) {
    // ยังไม่พร้อมจริง: ให้คง loading และพยายาม ping health ต่อเป็นระยะ
    // (ไม่โยน error เพื่อกันแอปดับ แต่คุณจะยังอยู่ใต้ loading overlay)
    // อาจจะตั้ง interval ภายนอกให้โหลดซ้ำเมื่อ health โอเค
  }

  // คืนสถานะที่ต้องใช้ตัดสินการปิด loading ภายนอก
  return {
    healthOk: net.healthOk,
    initialOk: allOk,
  }
}
