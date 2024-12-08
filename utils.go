package main

import (
	"encoding/json"
	"fmt"
	"os"
	"unicode"

	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/bcrypt"
)



func clearEntries(entries ...*widget.Entry) {
	for _, entry := range entries {
		entry.SetText("")
	}
}

// capitalizeFirstLetter делает первую букву строки заглавной
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Функция для сохранения записей в JSON файл
func saveRecordsToFile(records []Record, filename string) error {
	// Открываем или создаем файл для записи
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Не удалось создать файл: %v", err)
	}
	defer file.Close()

	// Создаем JSON encoder
	encoder := json.NewEncoder(file)
	// Устанавливаем отступы для лучшей читаемости
	encoder.SetIndent("", "  ")

	// Записываем записи в файл в формате JSON
	err = encoder.Encode(records)
	if err != nil {
		return fmt.Errorf("Не удалось записать в файл: %v", err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
