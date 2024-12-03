package repository

import (
	"encoding/json"
	"leveling/internal/server/entity"
)

var enemiesJsonData = map[string]string{
	"Enemy": `{"name": "Enemy", "Health": 100, "Strength": 8, "mainHand": 0}`,
}

var heroJsonData = map[string]string{
	"Taras": `{"name": "Taras", "Health": 100, "Strength": 8, "mainHand": 2}`,
	"Sin":   `{"name": "Sin", "Health": 100, "Strength": 2, "mainHand": 1}`,
	"Brian": `{"name": "Brian", "Health": 100, "Strength": 6, "mainHand": 0}`,
}

func GetHeroData() (heroesEntity []entity.Hero) {
	for _, jsonDatum := range enemiesJsonData {
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
	data := entity.Hero{}
	json.Unmarshal([]byte(heroJsonData[name]), &data)

	return data
}
