package database

import (
	"log"
	"renault-backend/models"
)

func seedCarsData() error {
	repo := NewCarRepository()

	// Проверяем, есть ли уже данные
	cars, err := repo.GetAllCars()
	if err != nil {
		return err
	}

	if len(cars) > 0 {
		log.Println("Cars data already seeded")
		return nil
	}

	// Легковые автомобили
	lightCars := []struct {
		car     models.Car
		details models.CarDetails
	}{
		{
			car: models.Car{
				Model:       "logan",
				Title:       "Renault Logan",
				Price:       "от 950 000 ₽",
				Category:    "light-cars",
				Image:       "images/renault_logan.jpeg",
				Description: "Renault Logan - это надежный и практичный седан, который идеально подходит для городских поездок и длительных путешествий. Сочетает в себе комфорт, экономичность и доступную цену.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.6 л, 82 л.с."},
					{Name: "Коробка передач", Value: "5-ступенчатая МКПП"},
					{Name: "Расход топлива", Value: "6.1 л/100км"},
					{Name: "Разгон 0-100 км/ч", Value: "11.9 сек"},
					{Name: "Объем багажника", Value: "510 л"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Климат-контроль", Value: "Есть"},
					{Name: "Электростеклоподъемники", Value: "Все"},
					{Name: "Аудиосистема", Value: "Radio Media Nav"},
					{Name: "Кондиционер", Value: "Есть"},
					{Name: "Круиз-контроль", Value: "Опция"},
					{Name: "Парктроник", Value: "Опция"},
				},
				Features: []string{
					"Система ABS+EBD",
					"Подушки безопасности",
					"Центральный замок",
					"Иммобилайзер",
					"Регулируемый руль",
					"Складывающиеся задние сиденья",
				},
			},
		},
		// Добавьте остальные автомобили аналогичным образом
		// ...
	}

	// Кроссоверы
	crossovers := []struct {
		car     models.Car
		details models.CarDetails
	}{
		{
			car: models.Car{
				Model:       "duster",
				Title:       "Renault Duster",
				Price:       "от 1 450 000 ₽",
				Category:    "crossovers",
				Image:       "images/duster.jpeg",
				Description: "Легендарный внедорожник Renault Duster с полным приводом готов покорить любые дороги. Проходимость, надежность и современный дизайн.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.6 л, 114 л.с."},
					{Name: "Привод", Value: "Полный 4x4"},
					{Name: "Клиренс", Value: "210 мм"},
					{Name: "Расход топлива", Value: "7.2 л/100км"},
					{Name: "Объем багажника", Value: "475 л"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Климат-контроль", Value: "2-зонный"},
					{Name: "Мультимедиа", Value: "Media Nav"},
					{Name: "Круиз-контроль", Value: "Есть"},
					{Name: "Парктроник", Value: "Передний и задний"},
					{Name: "Камера заднего вида", Value: "Есть"},
					{Name: "Подогрев сидений", Value: "Передние"},
				},
				Features: []string{
					"Система полного привода",
					"Защита картера",
					"Рейлинги на крыше",
					"Запасное колесо",
					"Защита бампера",
					"Противотуманные фары",
				},
			},
		},
		// Добавьте остальные кроссоверы
		// ...
	}

	// Заполняем базу данных
	log.Println("Seeding cars data...")

	// Добавляем легковые автомобили
	for _, data := range lightCars {
		if err := repo.CreateCar(&data.car, &data.details); err != nil {
			log.Printf("Error creating car %s: %v", data.car.Model, err)
		}
	}

	// Добавляем кроссоверы
	for _, data := range crossovers {
		if err := repo.CreateCar(&data.car, &data.details); err != nil {
			log.Printf("Error creating car %s: %v", data.car.Model, err)
		}
	}

	// Добавьте остальные категории аналогичным образом

	log.Println("Cars data seeded successfully")
	return nil
}
