package main

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var passwordHash string

var myApp fyne.App
var myWindow fyne.Window

func init() {
	hashedPassword, _ := hashPassword("12345")
	passwordHash = hashedPassword
}

func main() {
	// Создаем приложение
	myApp = app.New()
	myWindow = myApp.NewWindow("База данных сайта знакомств")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Подключение к базе данных
	conn, err := connectDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer conn.Close(context.Background())

	// Открываем окно авторизации
	loginWindow(conn)

	// Показываем основное окно
	myWindow.ShowAndRun()
}
