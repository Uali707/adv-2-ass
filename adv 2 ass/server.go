package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template" // Добавьте этот импорт для шаблонов
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Структура для данных о товаре
type Device struct {
	ID      uint    `gorm:"primaryKey"`
	Name    string  `gorm:"type:varchar(100)"`
	Price   float64 `gorm:"type:decimal(10,2)"`
	Catalog string  `gorm:"type:varchar(50)"`
}

var db *gorm.DB
var logger = logrus.New()
var limiter = rate.NewLimiter(rate.Every(1*time.Second), 5) // Лимит 5 запросов в секунду для всех
var clientLimiter = make(map[string]*rate.Limiter)
var mu sync.Mutex

// Инициализация логирования
func init() {
	// Устанавливаем формат логов
	logger.SetFormatter(&logrus.JSONFormatter{})
	// Выводим логи в стандартный вывод
	logger.SetOutput(os.Stdout)
	// Уровень логирования
	logger.SetLevel(logrus.InfoLevel)
}

// Инициализация базы данных
func initDB() {
	var err error
	dsn := "host=localhost user=postgres password=newpassword dbname=advprog port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"dsn": dsn,
		}).Fatal("Failed to connect to database")
	}

	logger.Info("Database connected successfully")

	// Миграция модели Device
	if err := db.AutoMigrate(&Device{}); err != nil {
		logger.Fatal("Database migration failed: ", err)
	}
}

// Обработка ошибок
func handleError(w http.ResponseWriter, err error, message string, statusCode int) {
	logger.WithFields(logrus.Fields{
		"error": err,
	}).Error(message)
	http.Error(w, message, statusCode)
}

// Обработчик с rate limiting
func rateLimitedHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr

	mu.Lock()
	// Проверяем лимит для IP
	limiter, exists := clientLimiter[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(1*time.Second), 5)
		clientLimiter[ip] = limiter
	}
	mu.Unlock()

	if !limiter.Allow() {
		http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
		logger.WithFields(logrus.Fields{
			"ip": ip,
		}).Warn("Rate limit exceeded")
		return
	}

	// Ваш код обработчика
}

// Обработчик для отображения товаров
func productsHandler(w http.ResponseWriter, r *http.Request) {
	var devices []Device
	query := db

	// Фильтрация по каталогу
	catalog := r.URL.Query().Get("catalog")
	if catalog != "" {
		query = query.Where("catalog = ?", catalog)
	}

	// Фильтрация по цене
	minPrice := r.URL.Query().Get("min_price")
	if minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query = query.Where("price >= ?", min)
		}
	}
	maxPrice := r.URL.Query().Get("max_price")
	if maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query = query.Where("price <= ?", max)
		}
	}

	// Сортировка
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy != "" {
		sortOrder := r.URL.Query().Get("sort_order")
		if sortOrder != "desc" {
			sortOrder = "asc"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	}

	// Пагинация
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 5
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Получение данных из базы
	if err := query.Find(&devices).Error; err != nil {
		handleError(w, err, "Failed to fetch devices", http.StatusInternalServerError)
		return
	}

	// Рендеринг HTML-шаблона
	tmpl, err := template.ParseFiles("./public/products.html")
	if err != nil {
		handleError(w, err, "Failed to load template", http.StatusInternalServerError)
		return
	}

	// Передача данных в шаблон
	if err := tmpl.Execute(w, devices); err != nil {
		handleError(w, err, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// Мягкое завершение работы
func gracefulShutdown(srv *http.Server) {
	// Канал для получения сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала
	<-stop

	logger.Info("Shutting down server...")

	// Устанавливаем тайм-аут на завершение работы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server shutdown failed")
	}
	logger.Info("Server gracefully stopped")
}

func main() {
	// Инициализация базы данных
	initDB()

	// Создание сервера
	srv := &http.Server{
		Addr:    ":3000",
		Handler: http.DefaultServeMux,
	}

	// Обслуживание статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	// Маршрут для отображения продуктов
	http.HandleFunc("/products", productsHandler)

	// Запуск сервера в отдельной горутине
	go func() {
		logger.Info("Server is running on port 3000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Обработка мягкого завершения
	gracefulShutdown(srv)
}
