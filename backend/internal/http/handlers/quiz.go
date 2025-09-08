package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Quiz struct {
	ID    string   `json:"id"`
	Hints []string `json:"hints"`
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
}

func GetQuiz(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(quizzes))
	quiz := quizzes[idx]
	id := time.Now().Format("20060102150405") + ":" + quiz.Answer // ใช้ : เป็นตัวแบ่ง
	resp := Quiz{
		ID:    id,
		Hints: quiz.Hints,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func CheckQuiz(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		ID    string `json:"id"`
		Guess string `json:"guess"`
	}
	var req Req
	json.NewDecoder(r.Body).Decode(&req)
	parts := []rune(req.ID)
	sepIdx := 0
	for i, c := range parts {
		if c == ':' {
			sepIdx = i
			break
		}
	}
	answer := string(parts[sepIdx+1:])
	result := req.Guess == answer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"correct": result})
}
