import api from './api'

export type RoomStatus = 'waiting' | 'playing' | 'finished'
export interface Room {
  id: number; code: string; status: RoomStatus; max_players: number; owner_name: string;
}
export interface Player {
  id: number; name: string; is_owner: boolean; is_ready: boolean; score: number; is_out: boolean;
}
export interface RoundPayload {
  round_no: number;
  quiz_id: string; quiz_token: string; quiz_exp: number;
  seconds: number;        // เวลารอบนี้ (ฝั่ง server ส่งลงมา)
  level: number;          // สำหรับปรับยาก
}

export async function createRoom(ownerName: string, maxPlayers = 4) {
  return api.post('/api/rooms', { ownerName, maxPlayers })
}
export async function joinRoom(code: string, name: string) {
  return api.post(`/api/rooms/${encodeURIComponent(code)}/join`, { name })
}
export async function setReady(code: string, name: string, ready: boolean) {
  return api.post(`/api/rooms/${encodeURIComponent(code)}/ready`, { name, ready })
}
export async function startRoom(code: string, ownerName: string) {
  return api.post(`/api/rooms/${encodeURIComponent(code)}/start`, { ownerName })
}
export async function submitGuess(code: string, name: string, guess: string) {
  return api.post(`/api/rooms/${encodeURIComponent(code)}/guess`, { name, guess })
}
export async function leaveRoom(code: string, name: string) {
  return api.post(`/api/rooms/${encodeURIComponent(code)}/leave`, { name })
}

// WS helper
export function openRoomSocket(code: string): WebSocket {
  // อิง origin ของ backend จาก axios baseURL เพื่อไม่พลาดพอร์ต
  const base = new URL((api as any).defaults?.baseURL || `${location.protocol}//${location.host}`)
  const wsProto = base.protocol === 'https:' ? 'wss:' : 'ws:'
  const url = `${wsProto}//${base.host}/ws/rooms/${encodeURIComponent(code)}`
  return new WebSocket(url)
}