package entity

type Profile struct {
	Username        string        `json:"username"`
	Stats           PlayerStats   `json:"player_stats"` //evitei field promotion para facilitar jason
	MedalsNames     []string      `json:"medals"`       //armazena apenas nomes
	History         []MatchResult `json:"history"`
	CurrentCampaign *Campaign     `json:"current_campaign"`
	Campaigns       []Campaign    `json:"campaigns"`
}

// HasMedal verifica se player possui medalha
func (p *Profile) HasMedal(name string) bool {
	for _, m := range p.MedalsNames {
		if m == name {
			return true
		}
	}
	return false
}
