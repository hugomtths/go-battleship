package service

import (
	"encoding/json"
	"os"
	"fmt"
)

type Profile struct {
	Username string;
	TotalScore int;
	HighestScore int;
	GamesPlayed int;
	MedalsEarned int;
}

const defaultPath string = "internal/service/saves/profiles.json";

func FindProfile(username string) (*Profile, error) {
    profiles, err := LoadProfiles()
    if err != nil {
        return nil, err
    }

    for _, p := range profiles {
        if p.Username == username {
            // cria uma cópia em memória
            profile := p
            return &profile, nil
        }
    }

    return nil, fmt.Errorf("profile '%s' not found", username)
}


func SaveProfile(profile Profile) error {
    if _, err := os.Stat(defaultPath); err == nil { // arquivo existe, atualiza Json
        return UpdateProfile(profile)

    } else if !os.IsNotExist(err) { 
        return err
    }

    // arquivo não existe, cria um com o profile
    profiles := []Profile{profile}

    data, err := json.MarshalIndent(profiles, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(defaultPath, data, 0644)
}

func UpdateProfile(profile Profile) error {
    profiles, err := LoadProfiles()
    if err != nil {
        return err
    }

    updated := false
    for i, p := range profiles {
        if p.Username == profile.Username { // verifica se username já existe
            profiles[i] = profile
            updated = true
            break
        }
    }

    if !updated { // adiciona novo profile se username não existe
        profiles = append(profiles, profile)
    }

    data, err := json.MarshalIndent(profiles, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(defaultPath, data, 0644)
}

func LoadProfiles() ([]Profile, error) {
	
	data, err := os.ReadFile(defaultPath)
    if err != nil {
		if os.IsNotExist(err) {
			return []Profile{}, nil // arquivo ainda não existe
        }
        return nil, err
    }

    var profiles []Profile
    err = json.Unmarshal(data, &profiles)
    if err != nil {
		return nil, err
    }
	
    return profiles, nil
}

func RemoveProfile(username string) error {
    // se o arquivo não existe, não há o que remover
    if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
        return nil
    } else if err != nil {
        return err
    }

    profiles, err := LoadProfiles()
    if err != nil {
        return err
    }

    newProfiles := make([]Profile, 0, len(profiles))
    removed := false

    for _, p := range profiles {
        if p.Username != username {
            newProfiles = append(newProfiles, p)
        } else {
            removed = true
        }
    }

    // se não encontrou o profile
    if !removed {
        return nil
    }

    data, err := json.MarshalIndent(newProfiles, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(defaultPath, data, 0644)
}
