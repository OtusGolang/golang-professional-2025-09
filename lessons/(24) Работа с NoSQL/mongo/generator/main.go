package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const numRecords = 1000 // Количество генерируемых записей

var countries = []string{"Россия", "США", "Германия", "Франция", "Италия"}

type User struct {
	ID      string `json:"_id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

func generateRandomName(gender bool) string {
	namesMale := []string{"Игорь", "Дмитрий", "Александр", "Алексей", "Максим", "Михаил", "Андрей", "Сергей", "Евгений", "Павел"}
	namesFemale := []string{"Ольга", "Елена", "Марина", "Наталья", "Светлана", "Анна", "Татьяна", "Виктория", "Кристина", "Алёна"}
	surnamesM := []string{"Иванов", "Петров", "Васечкин", "Семенов", "Попов", "Шарапов", "Федотов", "Смирнов", "Рябцев", "Жуков"}
	surnamesF := []string{"Иванова", "Петрова", "Васечкина", "Семенова", "Попова", "Шарапова", "Федотова", "Смирнова", "Рябцева", "Жукова"}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	if gender {
		nameIndex := rand.Intn(len(namesMale))    // Индекс имени
		surnameIndex := rand.Intn(len(surnamesM)) // Индекс фамилии

		return namesMale[nameIndex] + " " + surnamesM[surnameIndex]
	} else {
		nameIndex := rand.Intn(len(namesFemale))  // Индекс имени
		surnameIndex := rand.Intn(len(surnamesF)) // Индекс фамилии

		return namesFemale[nameIndex] + " " + surnamesF[surnameIndex]
	}
}

func generateRecord(id int) User {
	gender := rand.Float64() > 0.5 // true для мужчин, false для женщин
	name := generateRandomName(gender)
	age := rand.Intn(60) + 18 // Возраст от 18 до 78 лет
	country := countries[rand.Intn(len(countries))]

	return User{ID: fmt.Sprintf("%d", id), Name: name, Age: age, Country: country}
}

func main() {
	users := make([]User, numRecords)

	for i := 0; i < numRecords; i++ {
		users[i] = generateRecord(i + 1)
	}

	file, err := os.Create("data.json")
	if err != nil {
		panic(fmt.Sprintf("Ошибка открытия файла: %v", err))
	}
	defer file.Close()

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "\t") // Для красивого отступа

	err = jsonEncoder.Encode(users)
	if err != nil {
		panic(fmt.Sprintf("Ошибка записи JSON: %v", err))
	}

	fmt.Println("Файл успешно создан.")
}
