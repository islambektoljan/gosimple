package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB инициализирует подключение к базе данных и выполняет миграции
func InitDB() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment")
	}

	// Получаем параметры из .env файла
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	sslmode := "disable"

	// Формируем строку подключения
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, sslmode)

	// Открываем соединение с базой данных
	sqlDB, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Создаём экземпляр драйвера для миграций
	driver, err := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Получаем корневую рабочую директорию проекта (не из cmd)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}

	// Путь к папке с миграциями относительно корня проекта
	migrationsPath := filepath.Join(cwd, "internal", "db", "migrations")

	// Преобразуем в путь с прямыми слэшами для совместимости с URL
	migrationsPath = filepath.ToSlash(migrationsPath)

	// Генерируем правильный URL для миграций
	migrationsURL := fmt.Sprintf("file://%s", migrationsPath)
	log.Printf("Migrations path: %s", migrationsURL)

	// Создаём миграции
	m, err := migrate.NewWithDatabaseInstance(migrationsURL, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	// Применяем миграции с использованием golang-migrate
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}

	// Подключаемся к базе данных через GORM
	gormDB, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Сохраняем подключение в глобальную переменную DB
	DB = gormDB
}
