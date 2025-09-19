// backend/internal/http/handlers/party.go
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

/*
Party mode (in-memory multiplayer quiz game)

This module implements a real-time multiplayer quiz game system with the following features:
- Room management with unique 6-character codes
- Player joining/leaving with WebSocket real-time updates
- Ready state management before game start
- Real-time quiz gameplay with automatic progression
- Score tracking and leaderboard display
- Automatic room cleanup when empty

REST API Endpoints:
  GET  /api/rooms              - List all available rooms
  POST /api/rooms              - Create a new room
  POST /api/rooms/{code}/join  - Join a room by code
  POST /api/rooms/{code}/ready - Toggle player ready state
  POST /api/rooms/{code}/start - Start the game (owner only)
  POST /api/rooms/{code}/guess - Submit answer during gameplay
  POST /api/rooms/{code}/leave - Leave the room
  GET  /api/rooms/{code}/snapshot - Get current room state

WebSocket:
  /ws/rooms/{code} - Real-time room updates

Game Flow:
1. Host creates room → gets unique code
2. Players join room using code
3. All players must be ready before game can start
4. Game starts with quiz questions
5. First correct answer advances to next question
6. Game ends when time runs out or all questions answered
7. Scores are saved and leaderboard displayed

Technical Notes:
- Uses in-memory storage (rooms sync.Map) for simplicity
- WebSocket hub for real-time communication
- Automatic cleanup of empty rooms
- Robust error handling and validation
- Thread-safe operations with mutex locks
*/

func init() {
	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	} // dev ในเครื่อง
	// เราเรียก endpoint /api/quiz ของแอพตัวเอง (same process)
	quizBaseURL = "http://127.0.0.1:" + port
	log.Printf("[party] quizBaseURL = %s", quizBaseURL)
}

// roomStatus represents the current state of a game room
type roomStatus string

const (
	statusWaiting  roomStatus = "waiting"  // Room is waiting for players to join and get ready
	statusPlaying  roomStatus = "playing"  // Game is currently in progress
	statusFinished roomStatus = "finished" // Game has ended, showing results
)

// Room represents a multiplayer game room
type Room struct {
	ID         int64      `json:"id"`          // Unique room identifier
	Code       string     `json:"code"`        // 6-character room code for joining
	OwnerName  string     `json:"owner_name"`  // Name of the room creator
	Status     roomStatus `json:"status"`      // Current room status
	MaxPlayers int        `json:"max_players"` // Maximum number of players allowed
	CreatedAt  time.Time  `json:"-"`           // Room creation timestamp (not sent to client)
}

// Player represents a participant in a game room
type Player struct {
	ID      int64  `json:"id"`       // Unique player identifier within the room
	Name    string `json:"name"`     // Player's display name
	IsOwner bool   `json:"is_owner"` // Whether this player is the room owner
	IsReady bool   `json:"is_ready"` // Whether the player is ready to start
	Score   int    `json:"score"`    // Player's current score
	IsOut   bool   `json:"is_out"`   // Whether the player is eliminated
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
	Guess       string       `json:"guess,omitempty"`
	Correct     *bool        `json:"correct,omitempty"`
	Seconds     int          `json:"seconds,omitempty"`
	Winner      *Player      `json:"winner,omitempty"`
	Leaderboard []LeaderItem `json:"leaderboard,omitempty"`

	// ✅ แจ้ง error ของรอบ/ระบบ
	Error string `json:"error,omitempty"`
}

type roomState struct {
	mu          sync.Mutex
	room        *Room
	players     []*Player
	round       *RoundPayload
	seconds     int
	roundSolved bool   // ✅ มีคนตอบถูกในรอบนี้แล้วหรือยัง
	category    string // ✅ หมวดหมู่ของเกม
}

var rooms sync.Map // code -> *roomState

// ---------- WS hub ----------
var (
	wsMu   sync.Mutex
	wsConn = map[string]map[*websocket.Conn]struct{}{}
)

var (
	quizBaseURL    string
	httpQuizClient = &http.Client{Timeout: 5 * time.Second} // ✅ กันค้าง
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

	connections, exists := wsConn[code]
	if !exists {
		return
	}

	// Remove failed connections while broadcasting
	var failedConnections []*websocket.Conn
	for c := range connections {
		if err := c.WriteJSON(msg); err != nil {
			log.Printf("WebSocket broadcast error for room %s: %v", code, err)
			failedConnections = append(failedConnections, c)
		}
	}

	// Clean up failed connections
	for _, c := range failedConnections {
		delete(connections, c)
		if err := c.Close(); err != nil {
			log.Printf("Error closing failed WebSocket connection: %v", err)
		}
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

// Helper function to write error responses with consistent format
func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
		"code":  http.StatusText(code),
	})
}

// Helper function to validate room state
func validateRoomState(st *roomState) error {
	if st == nil {
		return errors.New("room not found")
	}
	if st.room == nil {
		return errors.New("room data corrupted")
	}
	return nil
}

// Helper function to clean up a room after game ends
func cleanupRoom(code string) {
	// Remove room from active rooms
	rooms.Delete(code)

	// Broadcast room closed event to any remaining connections
	wsHubBroadcast(code, hubMsg{Type: "room_closed"})

	log.Printf("Room %s cleaned up after game ended", code)
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

// GET /api/rooms - List all available rooms
func ListRooms(w http.ResponseWriter, r *http.Request) {
	roomsList := make([]map[string]any, 0)
	rooms.Range(func(key, value interface{}) bool {
		state := value.(*roomState)
		state.mu.Lock()
		room := state.room
		playerCount := len(state.players)
		state.mu.Unlock()

		// Only show waiting rooms
		if room.Status == statusWaiting {
			roomInfo := map[string]any{
				"id":           room.ID,
				"code":         room.Code,
				"owner_name":   room.OwnerName,
				"status":       room.Status,
				"max_players":  room.MaxPlayers,
				"player_count": playerCount, // Add current player count
			}
			roomsList = append(roomsList, roomInfo)
		}
		return true
	})

	writeJSON(w, http.StatusOK, map[string]any{"rooms": roomsList})
}

// POST /api/rooms {ownerName, maxPlayers}
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var in struct {
		OwnerName  string `json:"ownerName"`
		MaxPlayers int    `json:"maxPlayers"`
		Category   string `json:"category"`
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
	if in.Category == "" {
		in.Category = "สัตว์" // Default to animals
	}

	code := codeGen()
	room := &Room{
		Code:       code,
		OwnerName:  in.OwnerName,
		Status:     statusWaiting,
		MaxPlayers: in.MaxPlayers,
		CreatedAt:  time.Now(),
	}
	st := &roomState{room: room, players: []*Player{}, category: in.Category}
	rooms.Store(code, st)

	writeJSON(w, http.StatusOK, map[string]any{"room": room})
}

// POST /api/rooms/{code}/join {name}
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	code := roomCode(r)
	st := getState(code)
	if err := validateRoomState(st); err != nil {
		writeError(w, http.StatusNotFound, "Room not found")
		return
	}

	var in struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		writeError(w, http.StatusBadRequest, "Player name is required")
		return
	}

	// Validate name length and characters
	if len(in.Name) > 20 {
		writeError(w, http.StatusBadRequest, "Player name too long (max 20 characters)")
		return
	}

	st.mu.Lock()
	defer st.mu.Unlock()

	if st.room.Status != statusWaiting {
		writeError(w, http.StatusConflict, "Game has already started")
		return
	}

	// Remove any existing player with the same name (handle re-joining)
	for i, p := range st.players {
		if strings.EqualFold(p.Name, in.Name) {
			st.players = append(st.players[:i], st.players[i+1:]...)
			break
		}
	}

	if len(st.players) >= st.room.MaxPlayers {
		writeError(w, http.StatusConflict, "Room is full")
		return
	}

	isOwner := len(st.players) == 0 && strings.EqualFold(st.room.OwnerName, in.Name)
	p := &Player{
		ID:      int64(len(st.players) + 1),
		Name:    in.Name,
		IsOwner: isOwner,
	}
	st.players = append(st.players, p)

	// Broadcast player joined event with error handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error broadcasting player_joined: %v", r)
		}
	}()
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

	// ✅ ตรวจสอบว่าผู้เล่นทุกคนพร้อมแล้ว
	if len(st.players) == 0 {
		http.Error(w, "no players in room", http.StatusBadRequest)
		return
	}

	// ตรวจสอบว่าผู้เล่นทุกคนพร้อมแล้ว (ยกเว้นเจ้าของห้องที่สามารถเริ่มได้คนเดียว)
	allReady := true
	for _, p := range st.players {
		if !p.IsReady {
			allReady = false
			break
		}
	}

	// ถ้ามีผู้เล่นมากกว่า 1 คน ต้องให้ทุกคนพร้อม
	if len(st.players) > 1 && !allReady {
		http.Error(w, "all players must be ready", http.StatusBadRequest)
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

	// Find and remove the player from the room
	for i, p := range st.players {
		if strings.EqualFold(p.Name, in.Name) {
			st.players = append(st.players[:i], st.players[i+1:]...)
			break
		}
	}

	// If no players left, clean up the room
	if len(st.players) == 0 {
		rooms.Delete(code)
		wsHubBroadcast(code, hubMsg{Type: "room_closed"})
	} else {
		// Broadcast player left event to remaining players
		wsHubBroadcast(st.room.Code, hubMsg{
			Type:    "player_left",
			Name:    in.Name,
			Players: clonePlayers(st.players),
		})
	}

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
		log.Printf("WebSocket upgrade failed for room %s: %v", code, err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing WebSocket connection for room %s: %v", code, err)
		}
	}()

	wsHubAdd(code, conn)
	defer wsHubRemove(code, conn)

	// snapshot แรก
	st.mu.Lock()
	if err := conn.WriteJSON(hubMsg{
		Type:    "snapshot",
		Room:    st.room,
		Players: clonePlayers(st.players),
		Round:   st.round,
	}); err != nil {
		log.Printf("Error sending initial snapshot to room %s: %v", code, err)
		st.mu.Unlock()
		return
	}
	st.mu.Unlock()

	// read pump (keepalive)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Printf("WebSocket read error for room %s: %v", code, err)
			return
		}
	}
}

// ---------- Round & Timer ----------

func startRoundLocked(ctx context.Context, st *roomState, round int) {
	q, err := getQuiz(ctx, 1, st.category) // level = 1 (ปรับได้)
	if err != nil {
		// ❌ ดึงคำถามไม่ได้ (เช่น พอร์ตผิด / โปรเซสยังไม่ตื่น)
		log.Printf("[party] getQuiz error: %v", err)
		// กลับไป waiting เพื่อให้กดเริ่มใหม่
		st.room.Status = statusWaiting
		st.round = nil
		st.seconds = 0
		st.roundSolved = false

		wsHubBroadcast(st.room.Code, hubMsg{
			Type:  "round_failed",
			Room:  st.room,
			Error: "quiz_unavailable",
		})
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
	st.roundSolved = false

	wsHubBroadcast(st.room.Code, hubMsg{Type: "round_started", Round: st.round})
	go tickTimer(st, round) // guard ด้วย roundNo กัน timer เก่าทับ
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

				// Broadcast game over event
				wsHubBroadcast(st.room.Code, hubMsg{
					Type:        "game_over",
					Winner:      champ,
					Leaderboard: leader,
				})

				// Clean up the room after game ends
				roomCode := st.room.Code
				st.mu.Unlock()
				cleanupRoom(roomCode)
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

func getQuiz(ctx context.Context, level int, category string) (*quizResp, error) {
	// ✅ ใช้พอร์ตจริงของโปรเซส
	url := quizBaseURL + "/api/quiz?level=" + strconv.Itoa(level)
	if category != "" {
		url += "&category=" + category
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := httpQuizClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("quiz server error: " + res.Status)
	}
	var out quizResp
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func checkQuiz(ctx context.Context, id, token string, exp int64, guess string) (bool, error) {
	url := quizBaseURL + "/api/quiz/check"
	body := strings.NewReader(`{"id":"` + id + `","token":"` + token + `","exp":` + strconv.FormatInt(exp, 10) + `,"guess":"` + jsonEscape(guess) + `"}`)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	req.Header.Set("Content-Type", "application/json")

	res, err := httpQuizClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return false, errors.New("quiz check error: " + res.Status)
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
