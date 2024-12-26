package repository

import (
	"encoding/json"
	"leveling/internal/server/repository/dao"
)

var enemiesJsonData = map[string]string{
	"Enemy1": `{"name": "Enemy1", "Health": 100, "Strength": 1, "mainHand": 0, "position": [50, 50]}`,
	"Enemy2": `{"name": "Enemy2", "Health": 100, "Strength": 1, "mainHand": 0, "position": [-100, 150]}`,
	"Enemy3": `{"name": "Enemy3", "Health": 100, "Strength": 1, "mainHand": 0, "position": [50, 100]}`,
}

var heroJsonData = map[string]string{
	"Taras": `{"name": "Taras", "Health": 100, "Strength": 8, "mainHand": 2, "position": [0, 0]}`,
	"Sin":   `{"name": "Sin", "Health": 100, "Strength": 2, "mainHand": 1, "position": [0, 0]}`,
	"Brian": `{"name": "Brian", "Health": 100, "Strength": 6, "mainHand": 0, "position": [0, 0]}`,
}

func GetHeroData() (heroesEntity []dao.Hero) {
	for _, jsonDatum := range enemiesJsonData {
		data := dao.Hero{}
		err := json.Unmarshal([]byte(jsonDatum), &data)
		if err != nil {
			panic(err)
		}
		heroesEntity = append(heroesEntity, data)
	}
	return
}

func GetHeroByName(name string) (heroEntity dao.Hero) {
	data := dao.Hero{}
	json.Unmarshal([]byte(heroJsonData[name]), &data)

	return data
}
