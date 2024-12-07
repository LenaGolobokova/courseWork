package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func connectDB() (*pgx.Conn, error) {
	connStr := "postgres://postgres:lena@localhost:5432/database?search_path=data"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func registerUser(conn *pgx.Conn, login, password string) error {
	// Хеширование пароля
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %v", err)
	}

	// Сохранение в базу данных
	_, err = conn.Exec(context.Background(),
		"INSERT INTO data.user_accounts (login, password_hash) VALUES ($1, $2)", login, hashedPassword)
	if err != nil {
		return fmt.Errorf("ошибка регистрации: %v", err)
	}

	return nil
}

func validateUserCredentials(conn *pgx.Conn, login, password string) (bool, error) {
	if login == "" || password == "" {
		return false, fmt.Errorf("логин и пароль не должны быть пустыми")
	}

	var passwordHash string
	err := conn.QueryRow(context.Background(),
		"SELECT password_hash FROM data.user_accounts WHERE login = $1", login).Scan(&passwordHash)
	if err != nil {
		return false, fmt.Errorf("неверный логин")
	}

	if !checkPasswordHash(password, passwordHash) {
		return false, fmt.Errorf("неверный пароль")
	}

	return true, nil
}

func addRecordToDB(conn *pgx.Conn, record Record) error {
	// Преобразуем первую букву каждого поля в заглавную
	record.LastName = capitalizeFirstLetter(record.LastName)
	record.FirstName = capitalizeFirstLetter(record.FirstName)
	record.MiddleName = capitalizeFirstLetter(record.MiddleName)
	record.ZodiacSign = capitalizeFirstLetter(record.ZodiacSign)
	record.Profession = capitalizeFirstLetter(record.Profession)

	query := `INSERT INTO data.users (last_name, first_name, middle_name, gender, age, height, weight, zodiac_sign, profession)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := conn.Exec(context.Background(), query,
		record.LastName, record.FirstName, record.MiddleName, record.Gender, record.Age,
		record.Height, record.Weight, record.ZodiacSign, record.Profession)
	return err
}

// Фильтрация записей по полю и значению
func filterRecordsByField(conn *pgx.Conn, field, value string) ([]Record, error) {
	query := fmt.Sprintf(`SELECT last_name, first_name, middle_name, gender, age, height, weight, zodiac_sign, profession FROM data.users WHERE %s = $1`, field)
	rows, err := conn.Query(context.Background(), query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var record Record
		err = rows.Scan(&record.LastName, &record.FirstName, &record.MiddleName, &record.Gender, &record.Age, &record.Height, &record.Weight, &record.ZodiacSign, &record.Profession)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	// Проверка на отсутствие записей
	if len(records) == 0 {
		return nil, fmt.Errorf("Записи '%s = %s' не найдены", field, value)
	}

	return records, nil
}

// Удаление записей по полю и значению
func deleteRecordsByField(conn *pgx.Conn, field, value string) error {
	query := fmt.Sprintf(`DELETE FROM data.users WHERE %s = $1`, field)
	_, err := conn.Exec(context.Background(), query, value)
	return err
}

// Редактирование записи по ФИО
func editRecordByName(conn *pgx.Conn, lastName, firstName, middleName, field, newValue string) error {
	var query string
	var args []interface{}

	// Формируем запрос в зависимости от наличия отчества
	if middleName == "" {
		query = fmt.Sprintf(`UPDATE data.users SET %s = $1 WHERE last_name = $2 AND first_name = $3 AND middle_name IS NULL`, field)
		args = []interface{}{newValue, lastName, firstName}
	} else {
		query = fmt.Sprintf(`UPDATE data.users SET %s = $1 WHERE last_name = $2 AND first_name = $3 AND middle_name = $4`, field)
		args = []interface{}{newValue, lastName, firstName, middleName}
	}

	_, err := conn.Exec(context.Background(), query, args...)
	return err
}

func getRecordsFromDB(conn *pgx.Conn) ([]Record, error) {
	rows, err := conn.Query(context.Background(), "SELECT last_name, first_name, middle_name, gender, age, height, weight, zodiac_sign, profession FROM data.users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var record Record
		err = rows.Scan(&record.LastName, &record.FirstName, &record.MiddleName, &record.Gender, &record.Age, &record.Height, &record.Weight, &record.ZodiacSign, &record.Profession)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
