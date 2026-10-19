package handlers

import (
		"log"
		"net/http"
		"html/template"
		"strings"
		"sync"
		"strconv"
		"time"

		"Project/internal/models"
)

// Хранение списка трат
var (
	expenses []models.Expense
	nextID   = 1
	locker       sync.RWMutex // Используем RWMutex для оптимизации чтения
)

// Загружаем шаблоны
func loadTemplates(page string) (*template.Template, error) {
	// Загружаем layout + конкретную страницу
	return template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/"+page,
	)
}

// Создание списка трат
func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
		tmpl, err := loadTemplates("expenses.html")
		if err != nil {
				log.Printf("Ошибка загрузки шаблона: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
		}

	locker.RLock()
	safeExpenses := make([]models.Expense, len(expenses))
	copy(safeExpenses, expenses)
	locker.RUnlock()

	// Выполняем шаблон
	err = tmpl.ExecuteTemplate(w, "layout", safeExpenses)
	if err != nil {
		log.Printf("Ошибка рендеринга: %v", err)
	}
}

// Обработка введенных данных
func AddHandler(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
				tmpl, err := loadTemplates("add.html")
				if err != nil {
						http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
						return
				}
				if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
						log.Println("Ошибка при загрузке формы:", err)
				}
				return
		}

		if r.Method == http.MethodPost {
				if err := r.ParseForm(); err != nil {
						http.Error(w, "Ошибка при чтении формы", http.StatusBadRequest)
						return
				}

				// Проверяем сумму
				amount, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("amount")), 64)
				if err != nil || amount <= 0 {
						http.Error(w, "Некорректная сумма: должно быть число больше нуля", http.StatusBadRequest)
						return
				}

				// Проверяем описание
				description := strings.TrimSpace(r.FormValue("description"))
				if description == "" || len(description) > 200 {
						http.Error(w, "Описание обязательно и не длиннее 200 символов", http.StatusBadRequest)
						return
				}

				// Проверяем дату
				date := r.FormValue("date")
				if _, err := time.Parse("2006-01-02", date); err != nil {
						http.Error(w, "Введенная дата не коректна", http.StatusBadRequest)
						return
				}

				locker.Lock()
				expenses = append(expenses, models.Expense{
						ID:			 nextID,
						Amount:		 amount,
						Description: description,
						Date:		 date,
				})
				nextID++
				locker.Unlock()

				// Возврат к списоку
				http.Redirect(w, r, "/expenses", http.StatusSeeOther)
				return
		}

		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

