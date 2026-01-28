package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"guesswhoclientapi/data"
	"guesswhoclientapi/sse"

	"github.com/google/uuid"
)

var broker *sse.Broker

func init() {
	broker = sse.NewBroker()
}

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// authMiddleware is a placeholder for JWT authentication
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement JWT validation
		// For now, we'll just pass the request through
		next.ServeHTTP(w, r)
	})
}

// SignupHandler handles team creation
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req data.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, exists := data.GetTeamByName(req.TeamName); exists {
		http.Error(w, "Team name already exists", http.StatusConflict)
		return
	}

	newTeam := &data.TeamData{
		ID:                 "team-" + uuid.New().String(),
		TeamName:           req.TeamName,
		Password:           req.Password, // In a real app, hash this!
		TeamColor:          fmt.Sprintf("#%06x", uuid.New().ID()%0xFFFFFF),
		ChallengeStartTime: time.Now(),
	}
	data.AddTeam(newTeam)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data.SignupResponse{Message: "Team created successfully"})
}

// LoginHandler handles team authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req data.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	team, ok := data.GetTeamByName(req.TeamName)
	if !ok || team.Password != req.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// TODO: Generate a real JWT
	token := "fake-jwt-token-for-" + team.ID

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data.LoginResponse{
		Token: token,
		Team: struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		}{
			ID:    team.ID,
			Name:  team.TeamName,
			Color: team.TeamColor,
		},
	})
}

// GetGameStateHandler returns the public game state
func GetGameStateHandler(w http.ResponseWriter, r *http.Request) {
	characters := data.GetCharacters()
	teams := data.GetTeams()

	// In a real app, you'd calculate the leaderboard based on scores
	leaderboard := []data.ApiLeaderboardEntry{}
	rank := 1
	for _, team := range teams {
		leaderboard = append(leaderboard, data.ApiLeaderboardEntry{
			Rank:          rank,
			TeamName:      team.TeamName,
			Score:         team.TotalScore,
			Solves:        team.TotalSolves,
			QuickestSolve: team.FastestSolve,
			TeamColor:     team.TeamColor,
		})
		rank++
	}

	charList := []data.ApiCharacter{}
	for _, char := range characters {
		charList = append(charList, *char)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data.GetGameStateResponse{
		Characters:  charList,
		Leaderboard: leaderboard,
	})
}

// GetTeamProgressHandler returns the private progress for the authenticated team
func GetTeamProgressHandler(w http.ResponseWriter, r *http.Request) {
	// In a real app, you'd get the team ID from the JWT
	teamID := "team-alpha-123" // Hardcoded for now
	team, ok := data.GetTeamByID(teamID)
	if !ok {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}

// ResetHandler resets the solvedCharacters array for the authenticated team
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	// In a real app, you'd get the team ID from the JWT
	teamID := "team-alpha-123" // Hardcoded for now
	team, ok := data.GetTeamByID(teamID)
	if !ok {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}
	team.SolvedCharacters = []string{}
	data.UpdateTeam(team)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data.ResetResponse{Message: "Progress reset successfully"})
}

// UpdateDataHandler simulates the game API updating the data
func UpdateDataHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simple simulation. In a real app, this would be triggered by the game engine.
	// For example, let's say a team solves a character.
	team, _ := data.GetTeamByID("team-alpha-123")
	charID := "char-002"
	team.SolvedCharacters = append(team.SolvedCharacters, charID)
	team.TotalSolves++
	team.TotalScore += 100
	data.UpdateTeam(team)

	char, _ := data.GetCharacters()[charID]
	char.SolvedByTeams = append(char.SolvedByTeams, struct {
		TeamID string `json:"teamId"`
		Color  string `json:"color"`
	}{TeamID: team.ID, Color: team.TeamColor})
	data.UpdateCharacter(char)

	// Broadcast a message to all clients to refetch data
	broker.Broadcast("update")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data updated and event broadcasted")
}

// EventsHandler is the handler for the SSE endpoint
func EventsHandler(w http.ResponseWriter, r *http.Request) {
	broker.ServeHTTP(w, r)
}

// NewRouterAndBroker creates a new router, registers the handlers, and returns the router and the SSE broker
func NewRouterAndBroker() (http.Handler, *sse.Broker) {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/signup", SignupHandler)
	mux.HandleFunc("/api/auth/login", LoginHandler)
	mux.HandleFunc("/api/game/state", GetGameStateHandler)
	mux.HandleFunc("/api/team/progress", GetTeamProgressHandler)
	mux.HandleFunc("/api/team/reset", ResetHandler)
	mux.HandleFunc("/api/teams/update-data", UpdateDataHandler)
	mux.HandleFunc("/events", EventsHandler)

	return corsMiddleware(mux), broker
}