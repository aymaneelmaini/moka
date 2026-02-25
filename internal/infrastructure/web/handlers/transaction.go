package handlers

import (
	"html/template"
	"github.com/aymaneelmaini/moka/internal/application"
	"net/http"
	"strconv"
	"time"
)

type TransactionHandler struct {
	addSalaryUC       *application.AddSalaryUseCase
	recordExpenseUC   *application.RecordExpenseUseCase
	templates         *template.Template
}

func NewTransactionHandler(
	addSalaryUC *application.AddSalaryUseCase,
	recordExpenseUC *application.RecordExpenseUseCase,
	templates *template.Template,
) *TransactionHandler {
	return &TransactionHandler{
		addSalaryUC:     addSalaryUC,
		recordExpenseUC: recordExpenseUC,
		templates:       templates,
	}
}

func (h *TransactionHandler) AddSalary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")

	output, err := h.addSalaryUC.Execute(application.AddSalaryInput{
		Amount:      amount,
		Description: description,
		Date:        time.Now(),
	})

	if err != nil {
		http.Error(w, "Failed to add salary", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Success": true,
		"Output":  output,
	}

	if err := h.templates.ExecuteTemplate(w, "salary_success.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *TransactionHandler) RecordExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	categoryName := r.FormValue("category")
	description := r.FormValue("description")

	output, err := h.recordExpenseUC.Execute(application.RecordExpenseInput{
		Amount:       amount,
		CategoryName: categoryName,
		Description:  description,
		Date:         time.Now(),
	})

	if err != nil {
		http.Error(w, "Failed to record expense", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Success": true,
		"Output":  output,
	}

	if err := h.templates.ExecuteTemplate(w, "expense_success.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
