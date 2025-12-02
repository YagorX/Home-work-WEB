package database

import (
	"log"
	"renault-backend/models"
)

// SeedCarsData заполняет базу данных всеми автомобилями из index5.html
func SeedCarsData() error {
	repo := NewCarRepository()

	// Проверяем, есть ли уже данные
	cars, err := repo.GetAllCars()
	if err != nil {
		return err
	}

	if len(cars) > 0 {
		log.Println("Cars data already exists, skipping seeding")
		return nil
	}

	log.Println("Starting to seed cars data...")

	// Легковые автомобили (Light Cars)
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
		{
			car: models.Car{
				Model:       "sandero",
				Title:       "Renault Sandero",
				Price:       "от 890 000 ₽",
				Category:    "light-cars",
				Image:       "images/renault_sander.jpg",
				Description: "Компактный хэтчбек Renault Sandero предлагает просторный салон и отличную маневренность в городских условиях. Идеальный выбор для повседневных поездок.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.6 л, 75 л.с."},
					{Name: "Коробка передач", Value: "5-ступенчатая МКПП"},
					{Name: "Расход топлива", Value: "5.8 л/100км"},
					{Name: "Разгон 0-100 км/ч", Value: "12.5 сек"},
					{Name: "Объем багажника", Value: "320 л"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Кондиционер", Value: "Есть"},
					{Name: "Электростеклоподъемники", Value: "Передние"},
					{Name: "Аудиосистема", Value: "Radio Media Nav"},
					{Name: "Круиз-контроль", Value: "Опция"},
					{Name: "Давление в шинах", Value: "Контроль"},
					{Name: "Сигнализация", Value: "Есть"},
				},
				Features: []string{
					"Система стабилизации",
					"Передние подушки безопасности",
					"Регулировка руля",
					"Подогрев зеркал",
					"Противотуманные фары",
					"Легкосплавные диски",
				},
			},
		},
		{
			car: models.Car{
				Model:       "stepway",
				Title:       "Renault Sandero Stepway",
				Price:       "от 1 100 000 ₽",
				Category:    "light-cars",
				Image:       "images/renault_sander_stepway.jpeg",
				Description: "Renault Sandero Stepway - это хэтчбек в кросс-кузове с увеличенным клиренсом и стильным дизайном. Идеален для городских приключений.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.6 л, 90 л.с."},
					{Name: "Коробка передач", Value: "5-ступенчатая МКПП"},
					{Name: "Клиренс", Value: "195 мм"},
					{Name: "Расход топлива", Value: "6.2 л/100км"},
					{Name: "Объем багажника", Value: "320 л"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Климат-контроль", Value: "Есть"},
					{Name: "Электростеклоподъемники", Value: "Все"},
					{Name: "Мультимедиа", Value: "Media Nav"},
					{Name: "Защита бампера", Value: "Есть"},
					{Name: "Рейлинги на крыше", Value: "Есть"},
					{Name: "Легкосплавные диски", Value: "16\""},
				},
				Features: []string{
					"Увеличенный клиренс",
					"Защита бампера",
					"Рейлинги на крыше",
					"Противотуманные фары",
					"Система стабилизации",
					"Подогрев зеркал",
				},
			},
		},
	}

	// Кроссоверы (Crossovers)
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
		{
			car: models.Car{
				Model:       "kaptur",
				Title:       "Renault Kaptur",
				Price:       "от 1 350 000 ₽",
				Category:    "crossovers",
				Image:       "images/kapture.jpeg",
				Description: "Стильный компактный кроссовер с передовыми технологиями безопасности. Идеальное сочетание городского комфорта и внедорожных возможностей.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.6 л, 113 л.с."},
					{Name: "Коробка передач", Value: "Вариатор"},
					{Name: "Клиренс", Value: "204 мм"},
					{Name: "Расход топлива", Value: "6.7 л/100км"},
					{Name: "Объем багажника", Value: "387 л"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Климат-контроль", Value: "2-зонный"},
					{Name: "Мультимедиа", Value: "R-Link 2"},
					{Name: "Система ESP", Value: "Есть"},
					{Name: "Круиз-контроль", Value: "Есть"},
					{Name: "Камера 360°", Value: "Опция"},
					{Name: "Бесключевой доступ", Value: "Есть"},
				},
				Features: []string{
					"Мультимедиа R-Link",
					"Система ESP",
					"Камера заднего вида",
					"Датчики парковки",
					"Светодиодные фары",
					"Бесключевой доступ",
				},
			},
		},
		{
			car: models.Car{
				Model:       "arkana",
				Title:       "Renault Arkana",
				Price:       "от 1 650 000 ₽",
				Category:    "crossovers",
				Image:       "images/arkana.jpeg",
				Description: "Элегантное кросс-купе с динамичным характером и просторным салоном. Уникальный дизайн и передовые технологии.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.3 л, 150 л.с."},
					{Name: "Коробка передач", Value: "Вариатор X-Tronic"},
					{Name: "Клиренс", Value: "198 мм"},
					{Name: "Расход топлива", Value: "6.4 л/100км"},
					{Name: "Объем багажника", Value: "480 л"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Климат-контроль", Value: "2-зонный"},
					{Name: "Мультимедиа", Value: "EASY LINK"},
					{Name: "Цифровая панель", Value: "7\""},
					{Name: "Круиз-контроль", Value: "Адаптивный"},
					{Name: "Подогрев руля", Value: "Есть"},
					{Name: "Панорамная крыша", Value: "Опция"},
				},
				Features: []string{
					"Купе-форма кузова",
					"Вариатор X-Tronic",
					"Цифровая приборная панель",
					"Мультимедиа EASY LINK",
					"Светодиодная оптика",
					"Адаптивный круиз-контроль",
				},
			},
		},
	}

	// Коммерческие автомобили (Commercial)
	commercial := []struct {
		car     models.Car
		details models.CarDetails
	}{
		{
			car: models.Car{
				Model:       "loganvan",
				Title:       "Renault Logan Van",
				Price:       "от 1 000 000 ₽",
				Category:    "commercial",
				Image:       "images/van.jpeg",
				Description: "Коммерческая версия Logan с увеличенным багажным отделением. Надежность и экономичность для бизнеса.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.6 л, 82 л.с."},
					{Name: "Коробка передач", Value: "5-ступенчатая МКПП"},
					{Name: "Объем багажника", Value: "800 л"},
					{Name: "Грузоподъемность", Value: "500 кг"},
					{Name: "Расход топлива", Value: "6.3 л/100км"},
					{Name: "Количество мест", Value: "2"},
				},
				Equipment: []models.CarSpec{
					{Name: "Кондиционер", Value: "Есть"},
					{Name: "Аудиосистема", Value: "Radio"},
					{Name: "Электростеклоподъемники", Value: "Передние"},
					{Name: "Центральный замок", Value: "Есть"},
					{Name: "Сигнализация", Value: "Есть"},
					{Name: "Грузовая перегородка", Value: "Опция"},
				},
				Features: []string{
					"Увеличенный багажник",
					"Низкий расход топлива",
					"Прочная подвеска",
					"Грузовая перегородка",
					"Защита грузового отсека",
					"Доступная цена",
				},
			},
		},
		{
			car: models.Car{
				Model:       "kangoo",
				Title:       "Renault Kangoo",
				Price:       "от 1 300 000 ₽",
				Category:    "commercial",
				Image:       "images/kangoo.jpeg",
				Description: "Компактный коммерческий автомобиль с отличной маневренностью. Идеален для городских перевозок.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.6 л, 90 л.с."},
					{Name: "Коробка передач", Value: "5-ступенчатая МКПП"},
					{Name: "Объем грузового отсека", Value: "3-4.6 м³"},
					{Name: "Грузоподъемность", Value: "650 кг"},
					{Name: "Расход топлива", Value: "6.8 л/100км"},
					{Name: "Сдвижные двери", Value: "2"},
				},
				Equipment: []models.CarSpec{
					{Name: "Кондиционер", Value: "Есть"},
					{Name: "Аудиосистема", Value: "Radio"},
					{Name: "Электростеклоподъемники", Value: "Передние"},
					{Name: "Центральный замок", Value: "Есть"},
					{Name: "Сигнализация", Value: "Есть"},
					{Name: "Регулируемые сиденья", Value: "Есть"},
				},
				Features: []string{
					"Сдвижные боковые двери",
					"Большой грузовой объем",
					"Хорошая маневренность",
					"Низкий расход топлива",
					"Простая эксплуатация",
					"Доступное обслуживание",
				},
			},
		},
		{
			car: models.Car{
				Model:       "trafic",
				Title:       "Renault Trafic",
				Price:       "от 1 800 000 ₽",
				Category:    "commercial",
				Image:       "images/trafic.jpg",
				Description: "Универсальный коммерческий автомобиль для перевозки грузов. Надежность и вместительность.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "1.6 л дизель, 120 л.с."},
					{Name: "Коробка передач", Value: "6-ступенчатая МКПП"},
					{Name: "Объем грузового отсека", Value: "5.2-8.6 м³"},
					{Name: "Грузоподъемность", Value: "1-1.5 т"},
					{Name: "Расход топлива", Value: "7.1 л/100км"},
					{Name: "Количество мест", Value: "3"},
				},
				Equipment: []models.CarSpec{
					{Name: "Кондиционер", Value: "Есть"},
					{Name: "Мультимедиа", Value: "Media Nav"},
					{Name: "Круиз-контроль", Value: "Есть"},
					{Name: "Электростеклоподъемники", Value: "Все"},
					{Name: "Центральный замок", Value: "Есть"},
					{Name: "Система ESP", Value: "Есть"},
				},
				Features: []string{
					"Дизельный двигатель",
					"Большой грузовой объем",
					"Экономичный расход",
					"Комфортная кабина",
					"Надежная конструкция",
					"Легкое управление",
				},
			},
		},
	}

	// Электромобили и гибриды (Electro)
	electro := []struct {
		car     models.Car
		details models.CarDetails
	}{
		{
			car: models.Car{
				Model:       "zoe",
				Title:       "Renault ZOE",
				Price:       "от 2 200 000 ₽",
				Category:    "electro",
				Image:       "images/zoe.jpeg",
				Description: "Компактный электромобиль для города с впечатляющим запасом хода. Экологичность и современные технологии.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "Электромотор, 135 л.с."},
					{Name: "Запас хода", Value: "395 км"},
					{Name: "Батарея", Value: "52 кВт·ч"},
					{Name: "Разгон 0-100 км/ч", Value: "9.5 сек"},
					{Name: "Быстрая зарядка", Value: "30 мин до 80%"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Климат-контроль", Value: "Есть"},
					{Name: "Мультимедиа", Value: "EASY LINK"},
					{Name: "Режимы вождения", Value: "3 режима"},
					{Name: "Регенеративное торможение", Value: "Есть"},
					{Name: "Мобильное приложение", Value: "Есть"},
					{Name: "Бесключевой доступ", Value: "Есть"},
				},
				Features: []string{
					"Большой запас хода",
					"Быстрая зарядка",
					"Тихая работа",
					"Нулевые выбросы",
					"Регенеративное торможение",
					"Умное приложение",
				},
			},
		},
		{
			car: models.Car{
				Model:       "megane",
				Title:       "Renault Megane E-Tech",
				Price:       "от 3 500 000 ₽",
				Category:    "electro",
				Image:       "images/megane e.jpg",
				Description: "Современный электрокроссовер с технологиями нового поколения. Инновации и премиальный комфорт.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "Электромотор, 220 л.с."},
					{Name: "Запас хода", Value: "470 км"},
					{Name: "Батарея", Value: "60 кВт·ч"},
					{Name: "Разгон 0-100 км/ч", Value: "7.4 сек"},
					{Name: "Быстрая зарядка", Value: "130 кВт"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Климат-контроль", Value: "2-зонный"},
					{Name: "Мультимедиа", Value: "OpenR Link 12\""},
					{Name: "Цифровая панель", Value: "12,3\""},
					{Name: "Круиз-контроль", Value: "Адаптивный"},
					{Name: "Панорамная крыша", Value: "Есть"},
					{Name: "Массаж сидений", Value: "Опция"},
				},
				Features: []string{
					"Большой запас хода",
					"Мощный электромотор",
					"Цифровая панель 12,3\"",
					"Мультимедиа OpenR Link",
					"Панорамная крыша",
					"Премиальная отделка",
				},
			},
		},
		{
			car: models.Car{
				Model:       "captur",
				Title:       "Renault Captur E-Tech",
				Price:       "от 1 900 000 ₽",
				Category:    "electro",
				Image:       "images/captur e.jpg",
				Description: "Гибридный кроссовер с экономичным расходом и отличной динамикой. Эффективность и стиль.",
			},
			details: models.CarDetails{
				TechSpecs: []models.CarSpec{
					{Name: "Двигатель", Value: "Гибрид, 140 л.с."},
					{Name: "Коробка передач", Value: "Автоматическая"},
					{Name: "Расход топлива", Value: "4.5 л/100км"},
					{Name: "Электро-привод", Value: "На малых скоростях"},
					{Name: "Объем багажника", Value: "536 л"},
					{Name: "Количество мест", Value: "5"},
				},
				Equipment: []models.CarSpec{
					{Name: "Климат-контроль", Value: "2-зонный"},
					{Name: "Мультимедиа", Value: "EASY LINK"},
					{Name: "Режимы вождения", Value: "4 режима"},
					{Name: "Регенеративное торможение", Value: "Есть"},
					{Name: "Круиз-контроль", Value: "Есть"},
					{Name: "Подогрев сидений", Value: "Передние"},
				},
				Features: []string{
					"Гибридная система",
					"Низкий расход топлива",
					"Электро-привод в городе",
					"Регенеративное торможение",
					"Несколько режимов вождения",
					"Экологичность",
				},
			},
		},
	}

	// Добавляем все автомобили в базу данных
	log.Println("Adding light cars...")
	for _, data := range lightCars {
		if err := repo.CreateCar(&data.car, &data.details); err != nil {
			log.Printf("Error creating car %s: %v", data.car.Model, err)
		} else {
			log.Printf("✓ Added car: %s", data.car.Title)
		}
	}

	log.Println("Adding crossovers...")
	for _, data := range crossovers {
		if err := repo.CreateCar(&data.car, &data.details); err != nil {
			log.Printf("Error creating car %s: %v", data.car.Model, err)
		} else {
			log.Printf("✓ Added car: %s", data.car.Title)
		}
	}

	log.Println("Adding commercial vehicles...")
	for _, data := range commercial {
		if err := repo.CreateCar(&data.car, &data.details); err != nil {
			log.Printf("Error creating car %s: %v", data.car.Model, err)
		} else {
			log.Printf("✓ Added car: %s", data.car.Title)
		}
	}

	log.Println("Adding electric/hybrid vehicles...")
	for _, data := range electro {
		if err := repo.CreateCar(&data.car, &data.details); err != nil {
			log.Printf("Error creating car %s: %v", data.car.Model, err)
		} else {
			log.Printf("✓ Added car: %s", data.car.Title)
		}
	}

	// Проверяем сколько автомобилей добавлено
	finalCars, err := repo.GetAllCars()
	if err != nil {
		log.Printf("Error checking final count: %v", err)
	} else {
		log.Printf("✅ Successfully seeded %d cars into the database", len(finalCars))
	}

	return nil
}
