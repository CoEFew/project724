// internal/db/ping.go
package db

import "context"

// ภายในแพ็กเกจนี้ควรมีตัวแปร pool อยู่แล้วจาก Init()
// สมมติว่าเป็น: var pool *pgxpool.Pool
// ถ้าใช้ lib อื่น ปรับตามจริงได้เลย

func Ping(ctx context.Context) error {
	if pool == nil {
		return ErrNotInitialized // ถ้ายังไม่มี ให้สร้าง error นี้ไว้ หรือใช้ errors.New
	}
	return pool.Ping(ctx)
}
