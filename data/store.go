package data

import (
	"fmt"
	"sync"
	"time"
)

var (
	// using a sync.RWMutex to handle concurrent read/write operations
	mu    sync.RWMutex
	Teams = make(map[string]*TeamData)
	Characters = make(map[string]*ApiCharacter)
)

func init() {
	// Initialize with some mock data
	startTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00.000Z")
	Teams["team-alpha-123"] = &TeamData{
		ID:                 "team-alpha-123",
		TeamName:           "Team Alpha",
		TeamColor:          "#FF5733",
		ChallengeStartTime: startTime,
		TotalSolves:        5,
		SolvedCharacters:   []string{"char-001", "char-007", "char-012", "char-023", "char-034"},
		FastestSolve:       12500,
		TotalScore:         1250,
		Password:           "password123",
	}

	// Initialize with some mock characters
	for i := 1; i <= 40; i++ {
		id := "char-" + fmt.Sprintf("%03d", i)
		Characters[id] = &ApiCharacter{
			ID:          id,
			ImageUrl:    "/characters/" + id + ".png",
			SolvedByTeams: []struct {
				TeamID string `json:"teamId"`
				Color  string `json:"color"`
			}{},
		}
	}
}

// GetTeams returns a copy of the teams map
func GetTeams() map[string]*TeamData {
	mu.RLock()
	defer mu.RUnlock()
	// Return a copy to prevent modification of the original map
	teamsCopy := make(map[string]*TeamData)
	for k, v := range Teams {
		teamsCopy[k] = v
	}
	return teamsCopy
}

// GetTeamByID returns a team by its ID
func GetTeamByID(id string) (*TeamData, bool) {
	mu.RLock()
	defer mu.RUnlock()
	team, ok := Teams[id]
	return team, ok
}

// GetTeamByName returns a team by its name
func GetTeamByName(name string) (*TeamData, bool) {
	mu.RLock()
	defer mu.RUnlock()
	for _, team := range Teams {
		if team.TeamName == name {
			return team, true
		}
	}
	return nil, false
}

// AddTeam adds a new team
func AddTeam(team *TeamData) {
	mu.Lock()
	defer mu.Unlock()
	Teams[team.ID] = team
}

// UpdateTeam updates an existing team
func UpdateTeam(team *TeamData) {
	mu.Lock()
	defer mu.Unlock()
	Teams[team.ID] = team
}

// GetCharacters returns a copy of the characters map
func GetCharacters() map[string]*ApiCharacter {
	mu.RLock()
	defer mu.RUnlock()
	// Return a copy to prevent modification of the original map
	charsCopy := make(map[string]*ApiCharacter)
	for k, v := range Characters {
		charsCopy[k] = v
	}
	return charsCopy
}

// UpdateCharacter updates an existing character
func UpdateCharacter(char *ApiCharacter) {
	mu.Lock()
	defer mu.Unlock()
	Characters[char.ID] = char
}