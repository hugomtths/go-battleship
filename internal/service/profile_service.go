package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/allanjose001/go-battleship/internal/entity"
	"github.com/allanjose001/go-battleship/internal/medal"
)

const defaultPath string = "internal/data/profiles.json"

var profiles []entity.Profile

// init carrega profiles em memoria ao iniciar jogo
func init() {
	var err error
	err = loadProfiles() //caso não carregue arquivos o jogo pode continuar normalmente
	if err != nil {
		fmt.Println("Erro carregando profiles:", err) // remover apos integração
		profiles = []entity.Profile{}
	}
}

// GetProfiles retorna lista de profiles
func GetProfiles() []entity.Profile {
	return profiles
}

// loadProfiles carrega profiles do arquivo para a variável em memória
func loadProfiles() error {
	data, err := os.ReadFile(defaultPath)
	if err != nil {
		if os.IsNotExist(err) {
			profiles = []entity.Profile{}
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &profiles)
	if err != nil {
		return err
	}
	return nil
}

// SaveProfile é basicamente um alias para update, pois update acaba salvando o profile mesmo assim caso não exista
func SaveProfile(profile entity.Profile) error {
	return UpdateProfile(profile)
}

// UpdateProfile atualiza perfil em memoria e no arquivo json
func UpdateProfile(profile entity.Profile) error {
	updated := false

	for i, p := range profiles {
		if p.Username == profile.Username {
			profiles[i] = profile
			updated = true
			break
		}
	}

	if !updated {
		profiles = append(profiles, profile)
	}

	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(defaultPath, data, 0644)
}

// RemoveProfile apaga perfil em memoria e reescreve json com lista resultante
func RemoveProfile(username string) error {
	// procura e remove o profile da lista em memória
	removed := false
	newProfiles := make([]entity.Profile, 0, len(profiles))
	for _, p := range profiles {
		if p.Username != username {
			newProfiles = append(newProfiles, p)
		} else {
			removed = true
		}
	}

	if !removed {
		return nil // não encontrou profile, nada a fazer
	}

	// atualiza a lista em memória
	profiles = newProfiles

	// grava no arquivo
	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(defaultPath, data, 0644)
}

// GetProfileMedals retorna medalhas em forma de struct -> [está aqui para evitar cyclic imports]
func GetProfileMedals(p entity.Profile) []*medal.Medal {
	return medal.GetMedals(p.MedalsNames)
}

// AddMatchToProfile repassa salvamento para profile e atualiza saves no arquivo,
// retorna numero de medalhas ganhas apos partida
func AddMatchToProfile(profile *entity.Profile, result entity.MatchResult) (int, error) {
	profile.History = append(profile.History, result)

	profile.Stats.ApplyMatch(result)

	newMedals := checkNewMedals(profile)

	e := UpdateProfile(*profile)

	if e != nil {
		return 0, e
	}
	return newMedals, nil
}

// checkNewMedals verifica medalhas do player, atualiza, e retorna o n de novas medalhas
// (seria method de profile, mas deu cyclic import)
func checkNewMedals(p *entity.Profile) int {
	gained := 0

	for name, m := range medal.MedalsMap {
		if p.HasMedal(name) { // se tem medalha, passa p prox iteração
			continue
		}

		if m.Verification(p.Stats) {
			p.MedalsNames = append(p.MedalsNames, name) //adiciona nova medalha
			gained++
		}
	}
	return gained
}
