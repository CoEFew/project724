import axios from 'axios';
import { useNetworkStore } from '../store/useNetworkStore';



const BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080' // fallback กันลืมตั้ง env

const api = axios.create({
  baseURL: BASE,
  timeout: 60000,
})

// ---- Interceptors เพื่อนับ pending ----
api.interceptors.request.use((config) => {
  // บางที store อาจยังไม่ ready ตอน import: ป้องกัน error ด้วย try/catch
  try {
    const net = useNetworkStore()
    net.onRequestStart()
  } catch {}
  return config
})

api.interceptors.response.use(
  (res) => {
    try {
      const net = useNetworkStore()
      net.onRequestEnd()
    } catch {}
    return res
  },
  (err) => {
    try {
      const net = useNetworkStore()
      net.onRequestEnd()
      net.setError(err?.message || 'Request error')
    } catch {}
    return Promise.reject(err)
  }
)

export default api;
