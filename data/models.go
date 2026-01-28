package data

import "time"

// TeamData is the central model for a team's private state.
type TeamData struct {
	ID                 string    `json:"id"`
	TeamName           string    `json:"teamName"`
	TeamColor          string    `json:"teamColor"`
	ChallengeStartTime time.Time `json:"challengeStartTime"`
	TotalSolves        int       `json:"totalSolves"`
	SolvedCharacters   []string  `json:"solvedCharacters"`
	FastestSolve       int       `json:"fastestSolve"` // Duration in milliseconds
	TotalScore         int       `json:"totalScore"`
	Password           string    `json:"-"` // This will not be exposed in the API
}

// ApiCharacter represents a character on the game board.
type ApiCharacter struct {
	ID          string `json:"id"`
	ImageUrl    string `json:"imageUrl"`
	SolvedByTeams []struct {
		TeamID string `json:"teamId"`
		Color  string `json:"color"`
	} `json:"solvedByTeams"`
}

// ApiLeaderboardEntry represents a single entry on the leaderboard.
type ApiLeaderboardEntry struct {
	Rank          int    `json:"rank"`
	TeamName      string `json:"teamName"`
	Score         int    `json:"score"`
	Solves        int    `json:"solves"`
	QuickestSolve int    `json:"quickestSolve"` // duration in ms
	TeamColor     string `json:"teamColor"`
}

// SignupRequest is the request body for the signup endpoint.
type SignupRequest struct {
	TeamName string `json:"teamName"`
	Password string `json:"password"`
}

// SignupResponse is the response body for the signup endpoint.
type SignupResponse struct {
	Message string `json:"message"`
}

// LoginRequest is the request body for the login endpoint.
type LoginRequest struct {
	TeamName string `json:"teamName"`
	Password string `json:"password"`
}

// LoginResponse is the response body for the login endpoint.
type LoginResponse struct {
	Token string `json:"token"`
	Team  struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"team"`
}

// GetGameStateResponse is the response for the /api/game/state endpoint.
type GetGameStateResponse struct {
	Characters []ApiCharacter      `json:"characters"`
	Leaderboard []ApiLeaderboardEntry `json:"leaderboard"`
}

// ResetResponse is the response for the /api/team/reset endpoint.
type ResetResponse struct {
	Message string `json:"message"`
}