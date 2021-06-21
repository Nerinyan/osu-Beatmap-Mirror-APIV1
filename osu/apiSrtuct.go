package osu

type BeatmapsetsSearch struct {
	Beatmapsets *[]BeatmapSetsIN `json:"beatmapsets"`
	Cursor      *struct {
		LastUpdate *string `json:"last_update"`
		Id         *string `json:"_id"`
	} `json:"cursor"`
	Search *struct {
		Sort *string `json:"sort"`
	} `json:"search"`
	RecommendedDifficulty float64      `json:"recommended_difficulty"`
	Error                 *interface{} `json:"error"`
	Total                 int          `json:"total"`
}
