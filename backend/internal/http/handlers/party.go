// backend/internal/http/handlers/party.go
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

/*
Party mode (in-memory)
REST:
  /api/rooms, /api/rooms/{code}/join|ready|start|guess|leave
WS:
  /ws/rooms/{code}

สเปกเกมรอบ 3:
- ใครตอบถูก → ทั้งห้องไปข้อถัดไปทันที (รีเซ็ตเวลา) ผู้ตอบถูกได้ +1
- ถ้าหมดเวลา และ “ไม่มีใครตอบถูก” ในข้อนั้น → เกมจบทันที พร้อม leaderboard
- ทุกข้อ โชว์ใบ้แรกได้จากฝั่ง FE (เราแค่เริ่มรอบใหม่)
*/

func init() { rand.Seed(time.Now().UnixNano()) }

type roomStatus string

const (
	statusWaiting  roomStatus = "waiting"
	statusPlaying  roomStatus = "playing"
	statusFinished roomStatus = "finished"
)

type Room struct {
	ID         int64      `json:"id"`
	Code       string     `json:"code"`
	OwnerName  string     `json:"owner_name"`
	Status     roomStatus `json:"status"`
	MaxPlayers int        `json:"max_players"`
	CreatedAt  time.Time  `json:"-"`
}

type Player struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	IsOwner bool   `json:"is_owner"`
	IsReady bool   `json:"is_ready"`
	Score   int    `json:"score"`
	IsOut   bool   `json:"is_out"`
}

type RoundPayload struct {
	RoundNo   int    `json:"round_no"`
	QuizID    string `json:"quiz_id"`
	QuizToken string `json:"quiz_token"`
	QuizExp   int64  `json:"quiz_exp"`
	Seconds   int    `json:"seconds"`
	Level     int    `json:"level"`
}

// leaderboard named type
type LeaderItem struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type hubMsg struct {
	Type    string        `json:"type"`
	Room    *Room         `json:"room,omitempty"`
	Players []*Player     `json:"players,omitempty"`
	Round   *RoundPayload `json:"round,omitempty"`

	// events
	Name        string       `json:"name,omitempty"`
	Guess       string       `json:"guess,omitempty"` // ✅ ส่งคำที่เดาครั้งล่าสุด (FE เก็บแสดงให้ "ฉัน" เห็นเท่านั้น)
	Correct     *bool        `json:"correct,omitempty"`
	Seconds     int          `json:"seconds,omitempty"`
	Winner      *Player      `json:"winner,omitempty"`
	Leaderboard []LeaderItem `json:"leaderboard,omitempty"`
}

type roomState struct {
	mu          sync.Mutex
	room        *Room
	players     []*Player
	round       *RoundPayload
	seconds     int
	roundSolved bool // ✅ มีคนตอบถูกในรอบนี้แล้วหรือยัง
}

var rooms sync.Map // code -> *roomState

// ---------- WS hub ----------
var (
	wsMu   sync.Mutex
	wsConn = map[string]map[*websocket.Conn]struct{}{}
)

func wsHubAdd(code string, c *websocket.Conn) {
	wsMu.Lock()
	defer wsMu.Unlock()
	if wsConn[code] == nil {
		wsConn[code] = map[*websocket.Conn]struct{}{}
	}
	wsConn[code][c] = struct{}{}
}

func wsHubRemove(code string, c *websocket.Conn) {
	wsMu.Lock()
	defer wsMu.Unlock()
	if m := wsConn[code]; m != nil {
		delete(m, c)
		if len(m) == 0 {
			delete(wsConn, code)
		}
	}
}

func wsHubBroadcast(code string, msg hubMsg) {
	wsMu.Lock()
	defer wsMu.Unlock()
	for c := range wsConn[code] {
		_ = c.WriteJSON(msg) // best-effort
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // FE/BE origin ต่างกันได้
}

// ---------- helpers ----------
func roomCode(r *http.Request) string {
	return strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "code")))
}
func codeGen() string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func getState(code string) *roomState {
	if v, ok := rooms.Load(strings.ToUpper(code)); ok {
		return v.(*roomState)
	}
	return nil
}
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
func findPlayer(arr []*Player, name string) *Player {
	for _, p := range arr {
		if strings.EqualFold(p.Name, name) {
			return p
		}
	}
	return nil
}
func clonePlayers(arr []*Player) []*Player {
	out := make([]*Player, 0, len(arr))
	for _, p := range arr {
		cp := *p
		out = append(out, &cp)
	}
	return out
}

// ---------- REST ----------

// POST /api/rooms {ownerName, maxPlayers}
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var in struct {
		OwnerName  string `json:"ownerName"`
		MaxPlayers int    `json:"maxPlayers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	in.OwnerName = strings.TrimSpace(in.OwnerName)
	if in.OwnerName == "" {
		http.Error(w, "ownerName required", http.StatusBadRequest)
		return
	}
	if in.MaxPlayers < 2 || in.MaxPlayers > 4 {
		in.MaxPlayers = 4
	}

	code := codeGen()
	room := &Room{
		Code:       code,
		OwnerName:  in.OwnerName,
		Status:     statusWaiting,
		MaxPlayers: in.MaxPlayers,
		CreatedAt:  time.Now(),
	}
	st := &roomState{room: room, players: []*Player{}}
	rooms.Store(code, st)

	writeJSON(w, http.StatusOK, map[string]any{"room": room})
}

// POST /api/rooms/{code}/join {name}
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	code := roomCode(r)
	st := getState(code)
	if st == nil {
		http.NotFound(w, r)
		return
	}
	var in struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	if st.room.Status != statusWaiting {
		http.Error(w, "already started", http.StatusConflict)
		return
	}
	if len(st.players) >= st.room.MaxPlayers {
		http.Error(w, "room full", http.StatusConflict)
		return
	}
	if findPlayer(st.players, in.Name) != nil {
		http.Error(w, "name taken", http.StatusConflict)
		return
	}

	isOwner := len(st.players) == 0 && strings.EqualFold(st.room.OwnerName, in.Name)
	p := &Player{
		ID:      int64(len(st.players) + 1),
		Name:    in.Name,
		IsOwner: isOwner,
	}
	st.players = append(st.players, p)

	wsHubBroadcast(st.room.Code, hubMsg{Type: "player_joined", Players: clonePlayers(st.players)})
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// POST /api/rooms/{code}/ready {name, ready}
func ReadyRoom(w http.ResponseWriter, r *http.Request) {
	code := roomCode(r)
	st := getState(code)
	if st == nil {
		http.NotFound(w, r)
		return
	}
	var in struct {
		Name  string `json:"name"`
		Ready bool   `json:"ready"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	if st.room.Status != statusWaiting {
		http.Error(w, "already started", http.StatusConflict)
		return
	}
	p := findPlayer(st.players, in.Name)
	if p == nil {
		http.Error(w, "not in room", http.StatusNotFound)
		return
	}
	p.IsReady = in.Ready

	wsHubBroadcast(st.room.Code, hubMsg{Type: "ready_changed", Players: clonePlayers(st.players)})
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// POST /api/rooms/{code}/start {ownerName}
// POST /api/rooms/{code}/start {ownerName}
func StartRoom(w http.ResponseWriter, r *http.Request) {
	code := roomCode(r)
	st := getState(code)
	if st == nil {
		http.NotFound(w, r)
		return
	}
	var in struct {
		OwnerName string `json:"ownerName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	owner := findPlayer(st.players, in.OwnerName)
	if owner == nil || !owner.IsOwner {
		http.Error(w, "only owner can start", http.StatusForbidden)
		return
	}

	// ✅ ถ้าเริ่มไปแล้ว ให้ตอบ 200 เฉย ๆ (idempotent) ไม่ส่ง 409
	if st.room.Status == statusPlaying {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":   true,
			"note": "already_started",
		})
		return
	}
	if st.room.Status == statusFinished {
		http.Error(w, "game finished", http.StatusConflict)
		return
	}

	st.room.Status = statusPlaying
	startRoundLocked(r.Context(), st, 1)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// GET /api/rooms/{code}/snapshot
func RoomSnapshot(w http.ResponseWriter, r *http.Request) {
	code := roomCode(r)
	st := getState(code)
	if st == nil {
		http.NotFound(w, r)
		return
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	writeJSON(w, http.StatusOK, hubMsg{
		Type:    "snapshot",
		Room:    st.room,
		Players: clonePlayers(st.players),
		Round:   st.round,
	})
}

// POST /api/rooms/{code}/guess {name, guess}
func GuessRoom(w http.ResponseWriter, r *http.Request) {
	code := roomCode(r)
	st := getState(code)
	if st == nil {
		http.NotFound(w, r)
		return
	}
	var in struct {
		Name  string `json:"name"`
		Guess string `json:"guess"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Guess = strings.TrimSpace(in.Guess)
	if in.Name == "" || in.Guess == "" {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	if st.room.Status != statusPlaying || st.round == nil {
		http.Error(w, "not playing", http.StatusConflict)
		return
	}
	p := findPlayer(st.players, in.Name)
	if p == nil {
		http.Error(w, "not in room", http.StatusNotFound)
		return
	}
	if p.IsOut {
		http.Error(w, "player is out", http.StatusConflict)
		return
	}

	ok, err := checkQuiz(r.Context(), st.round.QuizID, st.round.QuizToken, st.round.QuizExp, in.Guess)
	if err != nil {
		http.Error(w, "quiz check error", http.StatusInternalServerError)
		return
	}
	if ok {
		// ✅ ผู้เล่นได้คะแนน และถือว่า "ข้อนี้มีคนตอบถูกแล้ว"
		p.Score++
		st.roundSolved = true

		// แจ้งผลกับ client (รวมคำเดาที่เพิ่งส่ง)
		wsHubBroadcast(st.room.Code, hubMsg{
			Type:    "guess_result",
			Name:    p.Name,
			Guess:   in.Guess,
			Correct: &ok,
			Players: clonePlayers(st.players),
		})

		// ✅ ไปข้อถัดไปทันที (รีเซ็ตเวลา)
		nextNo := st.round.RoundNo + 1
		startRoundLocked(r.Context(), st, nextNo)
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}

	// ผิด → แจ้งผล (แนบ guess กลับ)
	wsHubBroadcast(st.room.Code, hubMsg{
		Type:    "guess_result",
		Name:    p.Name,
		Guess:   in.Guess,
		Correct: &ok,
		Players: clonePlayers(st.players),
	})
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// POST /api/rooms/{code}/leave {name}
func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	code := roomCode(r)
	st := getState(code)
	if st == nil {
		http.NotFound(w, r)
		return
	}
	var in struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	in.Name = strings.TrimSpace(in.Name)

	st.mu.Lock()
	defer st.mu.Unlock()

	for i, p := range st.players {
		if strings.EqualFold(p.Name, in.Name) {
			st.players = append(st.players[:i], st.players[i+1:]...)
			break
		}
	}
	wsHubBroadcast(st.room.Code, hubMsg{Type: "player_joined", Players: clonePlayers(st.players)})
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// ---------- WS ----------
func RoomWS(w http.ResponseWriter, r *http.Request) {
	code := roomCode(r)
	st := getState(code)
	if st == nil {
		http.NotFound(w, r)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	wsHubAdd(code, conn)
	defer wsHubRemove(code, conn)

	// snapshot แรก
	st.mu.Lock()
	_ = conn.WriteJSON(hubMsg{
		Type:    "snapshot",
		Room:    st.room,
		Players: clonePlayers(st.players),
		Round:   st.round,
	})
	st.mu.Unlock()

	// read pump (keepalive)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

// ---------- Round & Timer ----------

func startRoundLocked(ctx context.Context, st *roomState, round int) {
	q, err := getQuiz(ctx, 1) // level = 1 (ปรับได้)
	if err != nil {
		log.Println("getQuiz:", err)
		return
	}
	st.round = &RoundPayload{
		RoundNo:   round,
		QuizID:    q.ID,
		QuizToken: q.Token,
		QuizExp:   q.Exp,
		Seconds:   60,
		Level:     1,
	}
	st.seconds = st.round.Seconds
	st.roundSolved = false // ✅ เริ่มรอบใหม่ ต้องรีเซ็ตสถานะ

	wsHubBroadcast(st.room.Code, hubMsg{Type: "round_started", Round: st.round})
	go tickTimer(st, round) // ✅ ปะ round guard กัน timer เก่าทับ
}

func tickTimer(st *roomState, roundNo int) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for range t.C {
		st.mu.Lock()

		// ถ้าห้องเปลี่ยนรอบไปแล้ว/เกมไม่ playing → ยุติกอรุตีนนี้
		if st.room.Status != statusPlaying || st.round == nil || st.round.RoundNo != roundNo {
			st.mu.Unlock()
			return
		}

		st.seconds--
		sec := st.seconds
		wsHubBroadcast(st.room.Code, hubMsg{Type: "timer_tick", Seconds: sec})

		if sec <= 0 {
			// ⛔ หมดเวลา — ถ้า "ไม่มีใครตอบถูก" ในรอบนี้ → จบทันที
			if !st.roundSolved {
				st.room.Status = statusFinished

				leader := make([]LeaderItem, 0, len(st.players))
				for _, p := range st.players {
					leader = append(leader, LeaderItem{Name: p.Name, Score: p.Score})
				}
				sort.Slice(leader, func(i, j int) bool { return leader[i].Score > leader[j].Score })

				var champ *Player
				if len(leader) > 0 {
					// หา Player จากชื่อ top
					topName := leader[0].Name
					for _, p := range st.players {
						if p.Name == topName {
							champ = p
							break
						}
					}
				}

				wsHubBroadcast(st.room.Code, hubMsg{
					Type:        "game_over",
					Winner:      champ,
					Leaderboard: leader,
				})
				st.mu.Unlock()
				return
			}

			// ถ้ามีคนตอบถูกแล้ว รอบถัดไปถูกเริ่มจาก GuessRoom แล้ว
			// timer นี้จะถูกยุติโดย guard ด้านบน (roundNo ไม่ตรง)
			st.mu.Unlock()
			return
		}
		st.mu.Unlock()
	}
}

// ---------- quiz HTTP client (reuse single-player) ----------
type quizResp struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	Exp       int64  `json:"exp"`
	HintCount int    `json:"hintCount"`
}

func getQuiz(ctx context.Context, level int) (*quizResp, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:8080/api/quiz?level="+strconv.Itoa(level), nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("quiz server error")
	}
	var out quizResp
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func checkQuiz(ctx context.Context, id, token string, exp int64, guess string) (bool, error) {
	body := strings.NewReader(`{"id":"` + id + `","token":"` + token + `","exp":` + strconv.FormatInt(exp, 10) + `,"guess":"` + jsonEscape(guess) + `"}`)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://127.0.0.1:8080/api/quiz/check", body)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return false, errors.New("quiz check error")
	}
	var out struct {
		Correct bool `json:"correct"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return false, err
	}
	return out.Correct, nil
}

func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	return strings.Trim(string(b), `"`)
}
