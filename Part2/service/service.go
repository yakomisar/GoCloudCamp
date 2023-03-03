package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"musicplayer/models"
	"net/http"
)

type Service struct {
	DB *gorm.DB
}

func (s *Service) Open() error {
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error when loading .env file")
		return err
	}
	// Пробуем открыть БД, используя файл окружения
	//sqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", env["DB_USER"], env["DB_PASSWORD"], env["DB_HOST"], env["DB_PORT"], env["DB_NAME"])
	//fmt.Println(sqlInfo)
	fmt.Println(env["DSN"])
	db, err := gorm.Open(postgres.Open(env["DSN"]), &gorm.Config{})
	if err != nil {
		log.Fatal("Error while openning database")
		return err
	}
	// Записываем БД в нашу структуру
	s.DB = db
	// Запускаем наш роутер, который отслеживает работу с запросами
	s.Routes()
	// Drop table if exists (will ignore or delete foreign key constraints when dropping)
	s.DB.Migrator().DropTable(&models.Song{})
	return s.DB.AutoMigrate(&models.Song{})
}

func (s *Service) Routes() http.Handler {
	r := mux.NewRouter()
	// Создадим маршруты
	// 1. Добавление песен
	r.HandleFunc("/api/addSong", s.createSong).Methods("POST")
	// 2. Получение списка песен
	r.HandleFunc("/api/getSongs", s.getSongs).Methods("POST")
	return r
}
