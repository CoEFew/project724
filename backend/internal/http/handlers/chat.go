package handlers

import (
	"encoding/json"
	"net/http"
)

type ChatRequest struct {
	Message string `json:"message"`
}

type OpenAIRespLite struct {
	OutputText string `json:"output_text"`
}

// --- ฟังก์ชันช่วย: ดึงข้อความจาก response รูปแบบต่างๆ ของ Responses API ---
func extractOutputText(raw map[string]any) string {
	// 1) ทางลัด: ถ้ามี output_text ใช้เลย
	if v, ok := raw["output_text"].(string); ok && v != "" {
		return v
	}
	// 2) รูปแบบหลักของ Responses API: output -> [ message ] -> content -> [ {type: output_text, text: "..."} ]
	if out, ok := raw["output"].([]any); ok {
		for _, item := range out {
			msg, _ := item.(map[string]any)
			if msg == nil {
				continue
			}
			if msg["type"] == "message" {
				if content, ok2 := msg["content"].([]any); ok2 {
					for _, c := range content {
						part, _ := c.(map[string]any)
						if part == nil {
							continue
						}
						// มาตรฐานใหม่จะเป็น type: "output_text"
						if (part["type"] == "output_text" || part["type"] == "text") && part["text"] != nil {
							if s, ok3 := part["text"].(string); ok3 && s != "" {
								return s
							}
						}
					}
				}
			}
		}
	}
	// 3) กันเหนียว: ถ้าบางรุ่นส่งเป็น choices แบบเดิม (ไม่ค่อยเกิดใน /v1/responses)
	if choices, ok := raw["choices"].([]any); ok && len(choices) > 0 {
		if ch0, ok := choices[0].(map[string]any); ok {
			if txt, ok := ch0["text"].(string); ok && txt != "" {
				return txt
			}
			// รูปแบบ chat-completions เก่า
			if msg, ok := ch0["message"].(map[string]any); ok {
				if content, ok := msg["content"].(string); ok && content != "" {
					return content
				}
			}
		}
	}
	return ""
}

// ยิง api หา chatGPT จริงต้องเสียตังเลยปิดไว้ก่อน
// func ChatHandler(w http.ResponseWriter, r *http.Request) {
// 	var req ChatRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}

// 	apiKey := os.Getenv("OPENAI_API_KEY")
// 	if apiKey == "" {
// 		http.Error(w, "Server misconfigured: OPENAI_API_KEY not set", http.StatusInternalServerError)
// 		return
// 	}

// 	payload := map[string]any{
// 		"model": "gpt-4.1-mini",
// 		"input": "ตอบสั้น กระชับ: " + req.Message,
// 	}
// 	body, _ := json.Marshal(payload)

// 	httpReq, _ := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(body))
// 	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
// 	httpReq.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(httpReq)
// 	if err != nil {
// 		http.Error(w, "OpenAI error", http.StatusBadGateway)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	respBody, _ := io.ReadAll(resp.Body)

// 	// passthrough error จาก OpenAI (จะเห็นข้อความจริง)
// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
// 		log.Printf("OpenAI %d: %s", resp.StatusCode, string(respBody))
// 		w.WriteHeader(resp.StatusCode)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(respBody)
// 		return
// 	}

// 	var raw map[string]any
// 	if err := json.Unmarshal(respBody, &raw); err != nil {
// 		http.Error(w, "Decode error", http.StatusInternalServerError)
// 		return
// 	}

// 	out := OpenAIRespLite{OutputText: extractOutputText(raw)}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(out)
// }

// ใช้ mock response แทนเรียก OpenAI
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// mock response สำหรับ dev (ไม่เสียตังค์)
	out := OpenAIRespLite{
		OutputText: "นี่คือข้อความ mock จาก backend: " + req.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
