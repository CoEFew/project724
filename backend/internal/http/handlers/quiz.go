package handlers

import (
	"crypto/hmac"
	// "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ==== (1) Data & Config ====

type QuizResp struct {
	ID    string   `json:"id"`
	Hints []string `json:"hints"`
	Token string   `json:"token"` // HMAC(secret, answer|id|exp)
	Exp   int64    `json:"exp"`   // unix seconds
}

var quizzes = []struct {
	Answer string
	Hints  []string
}{
	{"แมว", []string{"สัตว์เลี้ยงที่ร้องเหมียว", "มีหนวดและชอบกินปลา"}},
	{"หมา", []string{"สัตว์เลี้ยงที่เห่า", "เป็นเพื่อนที่ซื่อสัตย์"}},
	{"ช้าง", []string{"มีงวง", "ตัวใหญ่และมีงา"}},
	{"ม้า", []string{"วิ่งเร็ว", "ใช้ในการแข่ง"}},
	{"วัว", []string{"ให้นม", "มีเขา"}},
	{"กระต่าย", []string{"หูกระต่าย", "กินแครอท"}},
	{"เสือ", []string{"นักล่า", "ลายสวย"}},
	{"สิงโต", []string{"เจ้าป่า", "มีแผงคอ"}},
	{"แพนด้า", []string{"กินไผ่", "ตัวอ้วนกลม"}},
	{"โคอาล่า", []string{"อยู่บนต้นไม้", "ชอบกอดต้นยูคาลิปตัส"}},
	{"นก", []string{"บินได้", "มีปีก"}},
	{"ไก่", []string{"ขันตอนเช้า", "มีหงอน"}},
	{"เป็ด", []string{"เดินโยกเยก", "ร้องก๊าบๆ"}},
	{"นกยูง", []string{"หางสวย", "ชอบรำแพนหาง"}},
	{"กบ", []string{"กระโดดได้", "ร้องอ๊บๆ"}},
	{"งู", []string{"ไม่มีขา", "เลื้อยได้"}},
	{"เต่า", []string{"เดินช้า", "มีเปลือกแข็ง"}},
	{"มด", []string{"ตัวเล็ก", "ขยัน"}},
	{"ผึ้ง", []string{"ทำรัง", "มีเหล็กใน"}},
	{"ยุง", []string{"ดูดเลือด", "บินเสียงดัง"}},
	{"แมลงปอ", []string{"บินเร็ว", "ตาโต"}},
	{"ปลาทอง", []string{"เลี้ยงในตู้", "สีทอง"}},
	{"ปลาทู", []string{"หัวโต", "นิยมทำต้มยำ"}},
	{"ปลานิล", []string{"เลี้ยงง่าย", "เนื้ออร่อย"}},
	{"ปลาฉลาม", []string{"ฟันแหลมคม", "นักล่าในทะเล"}},
	{"โลมา", []string{"ฉลาด", "กระโดดน้ำ"}},
	{"ปลาวาฬ", []string{"ตัวใหญ่ที่สุดในทะเล", "พ่นน้ำ"}},
	{"ปู", []string{"เดินข้าง", "มีหนีบ"}},
	{"กุ้ง", []string{"มีเปลือกแข็ง", "มีหนวด"}},
	{"หอย", []string{"มีเปลือก", "อยู่ในน้ำ"}},
	{"หมู", []string{"ชอบกลิ้งโคลน", "เนื้อที่คนนิยมกิน"}},
	{"แพะ", []string{"เขาโค้ง", "กินหญ้าและใบไม้"}},
	{"แกะ", []string{"มีขนหนา", "ให้ขนมาทำเสื้อผ้า"}},
	{"ลิง", []string{"ปีนต้นไม้เก่ง", "ชอบกล้วย"}},
	{"กวาง", []string{"มีเขางาม", "วิ่งไวในป่า"}},
	{"จิ้งจอก", []string{"หางฟู", "เจ้าเล่ห์ในนิทาน"}},
	{"หมาป่า", []string{"อยู่เป็นฝูง", "หอนเสียงดัง"}},
	{"ค้างคาว", []string{"บินกลางคืน", "ห้อยหัวนอน"}},
	{"กระรอก", []string{"หางฟู", "ชอบเก็บถั่ว"}},
	{"เม่น", []string{"มีหนามแหลม", "ม้วนตัวได้"}},
	{"ยีราฟ", []string{"คอยาวมาก", "ชอบกินยอดไม้สูง"}},
	{"ม้าลาย", []string{"ลายขาวดำ", "อยู่ในทุ่งหญ้า"}},
	{"แรด", []string{"มีนอใหญ่", "ตัวหนาแรงเยอะ"}},
	{"ฮิปโป", []string{"อยู่ริมน้ำ", "ปากกว้างมาก"}},
	{"จิงโจ้", []string{"มีกระเป๋าหน้าท้อง", "กระโดดเร็ว"}},
	{"นกฮูก", []string{"ตาโต", "ตื่นกลางคืน"}},
}

// ดึง secret จาก env (ตั้งค่าเช่น HMAC_SECRET=super-secret ใน .env หรือระบบ deploy)
// ถ้าไม่ตั้งจะ fallback เป็นค่า default (แต่ควรตั้ง!)
func getSecret() []byte {
	sec := os.Getenv("HMAC_SECRET")
	if sec == "" {
		sec = "change-me-in-env"
	}
	return []byte(sec)
}

// สุ่ม id ปลอดภัย
func randomID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// สร้าง HMAC(base64) จาก s
func sign(secret []byte, s string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(s))
	sum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

// เทียบ HMAC
func equalHMAC(a, b string) bool {
	ab, _ := base64.StdEncoding.DecodeString(a)
	bb, _ := base64.StdEncoding.DecodeString(b)
	return hmac.Equal(ab, bb)
}

// ==== (2) Handlers ====

func GetQuiz(w http.ResponseWriter, r *http.Request) {
	// เลือกคำถามแบบสุ่ม
	now := time.Now()
	rand.Seed(now.UnixNano())
	idx := rand.Intn(len(quizzes))
	q := quizzes[idx]

	// ทำ id แบบสุ่ม + ใส่วันหมดอายุ (เช่น 90 วินาที)
	id, err := randomID()
	if err != nil {
		http.Error(w, "cannot generate id", http.StatusInternalServerError)
		return
	}
	exp := now.Add(90 * time.Second).Unix()

	// token = HMAC(secret, answer|id|exp)
	secret := getSecret()
	payload := q.Answer + "|" + id + "|" + strconvFormatInt(exp)
	token := sign(secret, payload)

	resp := QuizResp{
		ID:    id,
		Hints: q.Hints,
		Token: token,
		Exp:   exp,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

type CheckReq struct {
	ID    string `json:"id"`
	Guess string `json:"guess"`
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

func CheckQuiz(w http.ResponseWriter, r *http.Request) {
	var req CheckReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// หมดอายุ?
	if time.Now().Unix() > req.Exp {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"correct": false, "reason": "expired"})
		return
	}

	// เทียบ HMAC(secret, guess|id|exp) กับ token
	secret := getSecret()
	normalizedGuess := normalizeThai(req.Guess)
	signGuess := sign(secret, normalizedGuess+"|"+req.ID+"|"+strconvFormatInt(req.Exp))
	ok := equalHMAC(signGuess, req.Token)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"correct": ok})
}

// helper: ปรับรูปแบบตัวสะกด/ช่องว่าง (กันปัญหา Unicode ที่เทียบตรงๆ แล้วพลาด)
func normalizeThai(s string) string {
	// ตัดช่องว่างหัว/ท้าย + บีบช่องว่างซ้อน
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "  ", " ")
	// เพิ่มเติม: สามารถทำ NFC normalization ได้ถ้าต้องการ
	// (Go มาตรฐานไม่มีใน stdlib ต้องใช้แพ็กเกจภายนอก เช่น "golang.org/x/text/unicode/norm")
	return s
}

func strconvFormatInt(v int64) string {
	return strconv.FormatInt(v, 10)
}
