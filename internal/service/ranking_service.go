package service

import (
	"github.com/allanjose001/go-battleship/internal/entity"
	"sort"
)

func GetTopScores(limit int) []entity.Profile {
	profiles := GetProfiles()

	ranking := make([]entity.Profile, len(profiles))
	copy(ranking, profiles)

	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].Stats.TotalScore > ranking[j].Stats.TotalScore
	})

	if len(ranking) > limit {
		ranking = ranking[:limit]
	}

	return ranking
}
