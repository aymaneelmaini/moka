package handlers

import (
	"html/template"
	"moka/internal/application"
	"net/http"
	"strconv"
	"time"
)

type LoanHandler struct {
	borrowMoneyUC *application.BorrowMoneyUseCase
	payLoanUC     *application.PayLoanUseCase
	templates     *template.Template
}

func NewLoanHandler(
	borrowMoneyUC *application.BorrowMoneyUseCase,
	payLoanUC *application.PayLoanUseCase,
	templates *template.Template,
) *LoanHandler {
	return &LoanHandler{
		borrowMoneyUC: borrowMoneyUC,
		payLoanUC:     payLoanUC,
		templates:     templates,
	}
}

func (h *LoanHandler) BorrowMoney(w http.ResponseWriter, r *http.Request) {
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

	lenderName := r.FormValue("lender_name")
	description := r.FormValue("description")

	output, err := h.borrowMoneyUC.Execute(application.BorrowMoneyInput{
		LenderName:  lenderName,
		Amount:      amount,
		Description: description,
		Date:        time.Now(),
	})

	if err != nil {
		http.Error(w, "Failed to borrow money", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Success": true,
		"Output":  output,
	}

	if err := h.templates.ExecuteTemplate(w, "borrow_success.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *LoanHandler) PayLoan(w http.ResponseWriter, r *http.Request) {
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

	loanID := r.FormValue("loan_id")

	output, err := h.payLoanUC.Execute(application.PayLoanInput{
		LoanID: loanID,
		Amount: amount,
		Date:   time.Now(),
	})

	if err != nil {
		http.Error(w, "Failed to pay loan", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Success": true,
		"Output":  output,
	}

	if err := h.templates.ExecuteTemplate(w, "pay_loan_success.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
