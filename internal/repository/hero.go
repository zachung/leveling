package repository

import (
	"encoding/json"
	"leveling/internal/entity"
)

func GetHeroData() (heroesEntity []entity.Hero) {
	heroJsonData := []string{
		`{"name": "Taras", "Health": 100, "Strength": 8, "mainHand": 2}`,
		`{"name": "Sin", "Health": 100, "Strength": 2, "mainHand": 1}`,
	}
	for _, jsonDatum := range heroJsonData {
		data := entity.Hero{}
		err := json.Unmarshal([]byte(jsonDatum), &data)
		if err != nil {
			panic(err)
		}
		heroesEntity = append(heroesEntity, data)
	}
	return
}

func GetHeroByName(name string) (heroEntity entity.Hero) {
	jsonDatum := `{"name": "Brian", "Health": 100, "Strength": 6, "mainHand": 0}`
	data := entity.Hero{}
	json.Unmarshal([]byte(jsonDatum), &data)

	return data
}
