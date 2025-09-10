package handlers

import (
	"crypto/hmac"
	crand "crypto/rand" // ใช้สำหรับสุ่มไบต์ปลอดภัย
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

type QA struct {
	Answer string
	Hints  []string
}

// -------- Level 1 (ชุดง่าย - ของเดิม) --------
var quizzesEasy = []QA{
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

// -------- Level 2+ (ชุดยาก - ใหม่ 30 คำ) --------
var quizzesHard = []QA{
	{"เลียงผา", []string{"อยู่ภูเขาสูงผาชัน", "กระโดดเหวได้คล่อง"}},
	{"ละมั่ง", []string{"กวางเอเชีย", "เขาโค้งเรียวงาม"}},
	{"เก้ง", []string{"กวางตัวเล็ก", "บางท้องถิ่นร้องเหมือนเห่า"}},
	{"กระทิง", []string{"วัวป่าตัวใหญ่", "เขาดำโค้งทรงพลัง"}},
	{"สมเสร็จ", []string{"ลายขาวดำ", "จมูกยาวคล้ายงวง"}},
	{"ชะมด", []string{"ต่อมกลิ่นใช้ทำเครื่องหอม", "หากินกลางคืน"}},
	{"อีเห็น", []string{"คล้ายชะมด", "มีลายแถบตามตัว"}},
	{"นิ่ม", []string{"เกล็ดแข็งทั่วตัว", "ม้วนตัวป้องกันภัย"}},
	{"นาก", []string{"กึ่งบกน้ำ", "ชอบจับปลาเล่น"}},
	{"เสือดาว", []string{"ลายจุดโรเซต", "ปีนต้นไม้เก่ง"}},
	{"เสือชีตาห์", []string{"วิ่งเร็วสุดบนบก", "มีรอยเส้นน้ำตาดำ"}},
	{"เฟนเนก", []string{"สุนัขจิ้งจอกทะเลทราย", "หูใหญ่มาก"}},
	{"ลามา", []string{"ญาติอัลปากา", "ถิ่นภูเขาแอนดีส"}},
	{"อัลปากา", []string{"ขนฟูนุ่ม", "ไม่พอใจชอบพ่นน้ำลาย"}},
	{"ไคแมน", []string{"ญาติจระเข้ตัวเล็ก", "ลุ่มน้ำอเมซอน"}},
	{"โอคาปิ", []string{"ญาติกับยีราฟ", "ขาลายเหมือนม้าลาย"}},
	{"ไอบริกซ์", []string{"แพะภูเขา", "เขายาวบิดเป็นเกลียว"}},
	{"วอมแบต", []string{"กระเป๋าหน้าท้อง", "อุจจาระเป็นสี่เหลี่ยม"}},
	{"โคอาติ", []string{"จมูกยาวโค้ง", "อยู่ทวีปอเมริกา"}},
	{"คาปิบารา", []string{"หนูยักษ์", "ชอบลงน้ำแช่"}},
	{"อาร์มาดิลโล", []string{"เกราะแข็ง", "ม้วนตัวกลม"}},
	{"พิรันยา", []string{"ฟันแหลมคม", "หากินเป็นฝูง"}},
	{"คาสโซวารี", []string{"นกใหญ่คอสีฟ้า", "มีหมวกแข็งอันตราย"}},
	{"กิ้งก่าใบไม้", []string{"พรางตัวสุดเนียน", "ตัวแบนเหมือนใบไม้"}},
	{"แพนเธอร์", []string{"เสือสีดำเมลานิสติก", "ลายจุดซ่อนอยู่"}},
	{"คูดู", []string{"เขายาวบิดเป็นเกลียว", "อัฟริกา"}},
	{"ชิมแปนซี", []string{"ฉลาดใช้เครื่องมือ", "สังคมซับซ้อน"}},
	{"อุรังอุตัง", []string{"ลิงใหญ่สีส้ม", "แขนยาวอยู่คบไม้"}},
	{"วาฬเพชฌฆาต", []string{"โอร์ก้า", "นักล่าทะเลยอดฝีมือ"}},
	{"ควอกกา", []string{"หน้าตายิ้ม", "เกาะทางตะวันตกออสเตรเลีย"}},
}

// จัดกลุ่มตามระดับ (รองรับอนาคต)
var quizzesByLevel = map[int][]QA{
	1: quizzesEasy,
	2: quizzesHard,
}

// ดึง secret จาก env (ตั้งค่าเช่น HMAC_SECRET=super-secret)
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
	if _, err := crand.Read(buf); err != nil {
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

func poolForLevel(level int) []QA {
	// กติกาตามที่ตกลง: ใช้ชุดยากตั้งแต่ level >= 4
	tier := 1
	if level >= 4 {
		tier = 2
	}
	if p, ok := quizzesByLevel[tier]; ok && len(p) > 0 {
		return p
	}
	// fallback ปลอดภัย
	return quizzesEasy
}

func GetQuiz(w http.ResponseWriter, r *http.Request) {
	// อ่าน level จาก query
	level := 1
	if lv := r.URL.Query().Get("level"); lv != "" {
		if n, err := strconv.Atoi(lv); err == nil && n > 0 {
			level = n
		}
	}

	// ✅ ใช้ quizzesByLevel ผ่าน helper
	pool := poolForLevel(level)

	// สุ่มคำถามจาก pool
	now := time.Now()
	rand.Seed(now.UnixNano())
	idx := rand.Intn(len(pool))
	q := pool[idx]

	// ทำ id/token/exp เหมือนเดิม...
	id, err := randomID()
	if err != nil {
		http.Error(w, "cannot generate id", http.StatusInternalServerError)
		return
	}
	exp := now.Add(90 * time.Second).Unix()

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
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}

func strconvFormatInt(v int64) string {
	return strconv.FormatInt(v, 10)
}
