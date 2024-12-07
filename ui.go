package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5"
)

func loginWindow(conn *pgx.Conn) {
	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Введите логин")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	loginButton := widget.NewButton("Войти", func() {
		// Валидация логина и пароля
		isValid, err := validateUserCredentials(conn, loginEntry.Text, passwordEntry.Text)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		if isValid {
			dialog.ShowInformation("Успех", "Вы вошли!", myWindow)
			showDatabaseMenu(conn)
		}
	})

	registerButton := widget.NewButton("Зарегистрироваться", func() {
		showRegistrationWindow(conn)
	})

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Войдите, чтобы продолжить"),
		loginEntry,
		passwordEntry,
		loginButton,
		registerButton,
	))
}

func showRegistrationWindow(conn *pgx.Conn) {
	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Введите логин")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	adminCodeEntry := widget.NewEntry()
	adminCodeEntry.SetPlaceHolder("Введите код администратора")

	registerButton := widget.NewButton("Зарегистрироваться", func() {
		login := loginEntry.Text
		password := passwordEntry.Text
		adminCode := adminCodeEntry.Text

		if login == "" || password == "" || adminCode == "" {
			dialog.ShowError(fmt.Errorf("Все поля обязательны для заполнения"), myWindow)
			return
		}

		if adminCode != "H$7mZp!xQ9kLb3G#" {
			dialog.ShowError(fmt.Errorf("Неверный код администратора"), myWindow)
			return
		}

		// Вызов функции для регистрации
		err := registerUser(conn, login, password)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		dialog.ShowInformation("Успех", "Вы успешно зарегистрировались!", myWindow)
		loginWindow(conn)
	})

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Регистрация"),
		loginEntry,
		passwordEntry,
		adminCodeEntry,
		registerButton,
		widget.NewButton("Назад", func() {
			loginWindow(conn)
		}),
	))
}

func showAddRecordWindow(conn *pgx.Conn) {
	form := createRecordForm(func(record Record) {
		err := addRecordToDB(conn, record)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		dialog.ShowInformation("Успех", "Запись добавлена", myWindow)
		showDatabaseMenu(conn)
	})

	myWindow.SetContent(container.NewVBox(
		form,
		widget.NewButton("Назад", func() {
			showDatabaseMenu(conn)
		}),
	))
}

func showTableInNewWindow(records []Record) {
	tableWindow := myApp.NewWindow("Таблица записей")
	tableWindow.Resize(fyne.NewSize(800, 600))

	// Создаем таблицу
	table := createTableFromRecords(records)

	// Кнопка для сохранения записей в файл
	saveButton := widget.NewButton("Сохранить в файл", func() {
		err := saveRecordsToFile(records, "all_records.txt")
		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка сохранения в файл: %v", err), tableWindow)
			return
		}
		dialog.ShowInformation("Успех", "Записи сохранены в 'all_records.txt'", tableWindow)
	})

	// Кнопка для закрытия окна
	closeButton := widget.NewButton("Закрыть", func() {
		tableWindow.Close()
	})

	// Устанавливаем содержимое окна
	tableWindow.SetContent(container.NewBorder(nil, container.NewHBox(saveButton, closeButton), nil, nil, table))
	tableWindow.Show()
}

func showDatabaseMenu(conn *pgx.Conn) {
	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Работа с базой данных:"),
		widget.NewButton("Добавить запись", func() {
			showAddRecordWindow(conn)
		}),
		widget.NewButton("Посмотреть записи", func() {
			showAllRecords(conn)
		}),
		widget.NewButton("Фильтр записей по полю", func() {
			showFilterRecordsWindow(conn)
		}),
		widget.NewButton("Удалить записи", func() {
			showDeleteRecordsWindow(conn)
		}),
		widget.NewButton("Редактировать запись по ФИО", func() {
			showEditRecordWindow(conn)
		}),
		widget.NewButton("Выход", func() {
			loginWindow(conn)
		}),
	))
}

func createRecordForm(onSubmit func(record Record)) fyne.CanvasObject {
	lastNameEntry := widget.NewEntry()
	lastNameEntry.SetPlaceHolder("Фамилия")

	firstNameEntry := widget.NewEntry()
	firstNameEntry.SetPlaceHolder("Имя")

	middleNameEntry := widget.NewEntry()
	middleNameEntry.SetPlaceHolder("Отчество (необязательно)")

	genderSelect := widget.NewSelect([]string{"М", "Ж"}, nil)
	genderSelect.PlaceHolder = "Выберите пол"

	ageEntry := widget.NewEntry()
	ageEntry.SetPlaceHolder("Возраст")

	heightEntry := widget.NewEntry()
	heightEntry.SetPlaceHolder("Рост (см)")

	weightEntry := widget.NewEntry()
	weightEntry.SetPlaceHolder("Вес (кг)")

	zodiacSignEntry := widget.NewEntry()
	zodiacSignEntry.SetPlaceHolder("Знак Зодиака")

	professionEntry := widget.NewEntry()
	professionEntry.SetPlaceHolder("Профессия")

	submitButton := widget.NewButton("Добавить запись", func() {
		// Проверяем обязательные поля
		if lastNameEntry.Text == "" || firstNameEntry.Text == "" ||
			genderSelect.Selected == "" || ageEntry.Text == "" ||
			heightEntry.Text == "" || weightEntry.Text == "" ||
			zodiacSignEntry.Text == "" || professionEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Заполните все обязательные поля"), myWindow)
			return
		}

		// Проверяем длины текстовых полей
		if len(ageEntry.Text) > 3 || len(heightEntry.Text) > 3 || len(weightEntry.Text) > 3 || len(lastNameEntry.Text) > 100 || len(firstNameEntry.Text) > 100 || len(professionEntry.Text) > 100 || len(middleNameEntry.Text) > 100 {
			dialog.ShowError(fmt.Errorf("Некорректные данные"), myWindow)
			return
		}

		// Проверяем числовые поля
		age, err := strconv.Atoi(ageEntry.Text)
		if err != nil || age <= 0 {
			dialog.ShowError(fmt.Errorf("Возраст должен быть положительным числом"), myWindow)
			return
		}

		height, err := strconv.ParseFloat(heightEntry.Text, 64)
		if err != nil || height <= 0 {
			dialog.ShowError(fmt.Errorf("Рост должен быть положительным числом"), myWindow)
			return
		}

		weight, err := strconv.ParseFloat(weightEntry.Text, 64)
		if err != nil || weight <= 0 {
			dialog.ShowError(fmt.Errorf("Вес должен быть положительным числом"), myWindow)
			return
		}

		// Создаем запись
		record := Record{
			LastName:   capitalizeFirstLetter(lastNameEntry.Text),
			FirstName:  capitalizeFirstLetter(firstNameEntry.Text),
			MiddleName: capitalizeFirstLetter(middleNameEntry.Text),
			Gender:     genderSelect.Selected,
			Age:        age,
			Height:     height,
			Weight:     weight,
			ZodiacSign: capitalizeFirstLetter(zodiacSignEntry.Text),
			Profession: capitalizeFirstLetter(professionEntry.Text),
		}

		// Передаем запись в обработчик
		onSubmit(record)

		// Очищаем поля формы
		lastNameEntry.SetText("")
		firstNameEntry.SetText("")
		middleNameEntry.SetText("")
		ageEntry.SetText("")
		heightEntry.SetText("")
		weightEntry.SetText("")
		zodiacSignEntry.SetText("")
		professionEntry.SetText("")
		genderSelect.SetSelected("")
		professionEntry.SetText("")
	})

	return container.NewVBox(
		lastNameEntry,
		firstNameEntry,
		middleNameEntry,
		genderSelect,
		ageEntry,
		heightEntry,
		weightEntry,
		zodiacSignEntry,
		professionEntry,
		submitButton,
	)
}

// Просмотр всех записей
func showAllRecords(conn *pgx.Conn) {
	records, err := getRecordsFromDB(conn)
	if err != nil {
		dialog.ShowError(err, myWindow)
		return
	}
	showTableInNewWindow(records)
}

func createTableFromRecords(records []Record) fyne.CanvasObject {
	// Заголовки таблицы
	headers := []string{"Фамилия", "Имя", "Отчество", "Пол", "Возраст", "Рост", "Вес", "Знак Зодиака", "Профессия"}

	// Создаем данные таблицы
	data := [][]string{}
	for _, record := range records {
		// Преобразуем первую букву в заглавную для всех полей
		record.LastName = capitalizeFirstLetter(record.LastName)
		record.FirstName = capitalizeFirstLetter(record.FirstName)
		record.MiddleName = capitalizeFirstLetter(record.MiddleName)
		record.ZodiacSign = capitalizeFirstLetter(record.ZodiacSign)
		record.Profession = capitalizeFirstLetter(record.Profession)

		row := []string{
			record.LastName,
			record.FirstName,
			record.MiddleName,
			record.Gender,
			strconv.Itoa(record.Age),
			fmt.Sprintf("%.2f", record.Height),
			fmt.Sprintf("%.2f", record.Weight),
			record.ZodiacSign,
			record.Profession,
		}
		data = append(data, row)
	}

	// Создаем таблицу
	table := widget.NewTable(
		// Размер таблицы
		func() (int, int) { return len(data) + 1, len(headers) },
		// Создание объекта для ячейки
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Создаем пустую ячейку
		},
		// Наполнение данных
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row == 0 {
				// Если это первая строка, то выводим заголовки
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				// Для остальных строк выводим данные
				label.SetText(data[id.Row-1][id.Col]) // -1 из-за заголовков
			}
		},
	)

	// Настраиваем ширину колонок
	for i := range headers {
		table.SetColumnWidth(i, 100)
	}

	return table
}

func showFilterRecordsWindow(conn *pgx.Conn) {
	fieldSelect := widget.NewSelect([]string{"gender", "age", "height", "weight", "zodiac_sign", "profession"}, nil)
	fieldSelect.PlaceHolder = "Выберите поле"
	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("Введите значение")

	filterButton := widget.NewButton("Фильтровать", func() {
		field := fieldSelect.Selected
		value := capitalizeFirstLetter(valueEntry.Text)
		if field == "" || value == "" {
			dialog.ShowError(fmt.Errorf("Не выбрано поле или значение"), myWindow)
			return
		}
		// Выполняем фильтрацию
		records, err := filterRecordsByField(conn, field, value)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		// Сохраняем результаты в файл
		err = saveRecordsToFile(records, "filtered_records.txt")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		// Показываем уведомление об успехе
		dialog.ShowInformation("Успех", "Результаты фильтрации сохранены в 'filtered_records.txt'", myWindow)
		// Очищаем поля
		fieldSelect.Selected = ""
		valueEntry.SetText("")
		// Обновляем таблицу
		showTableInNewWindow(records)
	})

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Фильтр"),
		fieldSelect,
		valueEntry,
		container.NewHBox(
			filterButton,
			widget.NewButton("Назад", func() {
				showDatabaseMenu(conn)
			}),
		),
	))
}

func showDeleteRecordsWindow(conn *pgx.Conn) {
	fieldSelect := widget.NewSelect([]string{"gender", "age", "height", "weight", "zodiac_sign", "profession"}, nil)
	fieldSelect.PlaceHolder = "Выберите поле"
	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("Введите значение")

	deleteButton := widget.NewButton("Удалить", func() {
		field := fieldSelect.Selected
		value := capitalizeFirstLetter(valueEntry.Text)
		if field == "" || value == "" {
			dialog.ShowError(fmt.Errorf("Не выбрано поле или значение"), myWindow)
			return
		}

		// Получаем записи, которые будут удалены, для сохранения в файл
		records, err := filterRecordsByField(conn, field, value)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка получения записей: %v", err), myWindow)
			return
		}
		if len(records) == 0 {
			dialog.ShowError(fmt.Errorf("Записей не найдено для удаления"), myWindow)
			return
		}

		// Сохраняем записи в файл
		err = saveRecordsToFile(records, "deleted_records.txt")
		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка сохранения записей в файл: %v", err), myWindow)
			return
		}

		// Удаляем записи
		err = deleteRecordsByField(conn, field, value)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка удаления записей: %v", err), myWindow)
			return
		}

		// Показываем уведомление об успехе
		dialog.ShowInformation("Успех", "Записи удалены и сохранены в 'deleted_records.txt'", myWindow)
		// Очищаем поля
		fieldSelect.Selected = ""
		valueEntry.SetText("")
		// Обновляем таблицу
		showAllRecords(conn)
	})

	backButton := widget.NewButton("Назад", func() {
		showDatabaseMenu(conn)
	})

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Удаление записей"),
		fieldSelect,
		valueEntry,
		container.NewHBox(deleteButton, backButton),
	))
}

func showEditRecordWindow(conn *pgx.Conn) {
	lastNameEntry := widget.NewEntry()
	lastNameEntry.SetPlaceHolder("Фамилия")

	firstNameEntry := widget.NewEntry()
	firstNameEntry.SetPlaceHolder("Имя")

	middleNameEntry := widget.NewEntry()
	middleNameEntry.SetPlaceHolder("Отчество (необязательно)")

	fieldSelect := widget.NewSelect([]string{"gender", "age", "height", "weight", "zodiac_sign", "profession"}, nil)
	fieldSelect.PlaceHolder = "Выберите поле"

	newValueEntry := widget.NewEntry()
	newValueEntry.SetPlaceHolder("Введите новое значение")

	editButton := widget.NewButton("Редактировать", func() {
		lastName := capitalizeFirstLetter(lastNameEntry.Text)
		firstName := capitalizeFirstLetter(firstNameEntry.Text)
		middleName := capitalizeFirstLetter(middleNameEntry.Text)
		field := fieldSelect.Selected
		newValue := capitalizeFirstLetter(newValueEntry.Text)

		if lastName == "" || firstName == "" || field == "" || newValue == "" {
			dialog.ShowError(fmt.Errorf("Заполните все обязательные поля"), myWindow)
			return
		}

		// Получаем текущую запись перед редактированием
		records, err := filterRecordsByField(conn, "last_name", lastName)
		if err != nil || len(records) == 0 {
			dialog.ShowError(fmt.Errorf("Ошибка получения записи: %v", err), myWindow)
			return
		}
		// Сохраняем текущую запись в файл
		err = saveRecordsToFile(records, "edited_record.txt")
		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка сохранения записи в файл: %v", err), myWindow)
			return
		}

		// Редактируем запись
		err = editRecordByName(conn, lastName, firstName, middleName, field, newValue)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка редактирования записи: %v", err), myWindow)
			return
		}

		// Показываем уведомление об успехе
		dialog.ShowInformation("Успех", "Запись обновлена и исходные данные сохранены в 'edited_record.txt'", myWindow)
		// Очищаем поля
		lastNameEntry.SetText("")
		firstNameEntry.SetText("")
		middleNameEntry.SetText("")
		fieldSelect.Selected = ""
		newValueEntry.SetText("")
		// Обновляем таблицу
		showAllRecords(conn)
	})

	backButton := widget.NewButton("Назад", func() {
		showDatabaseMenu(conn)
	})

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Редактировать запись"),
		lastNameEntry,
		firstNameEntry,
		middleNameEntry,
		fieldSelect,
		newValueEntry,
		container.NewHBox(editButton, backButton),
	))
}
