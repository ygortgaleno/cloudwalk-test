package dtos

type GameDto struct {
	TotalKills uint32           `json:"total_kills"`
	Players    []string         `json:"players"`
	Kills      map[string]int64 `json:"kills"`
}
