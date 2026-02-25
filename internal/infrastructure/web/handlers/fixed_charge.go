package handlers

import (
	"html/template"
	"moka/internal/domain/fixed_charge"
	"moka/internal/shared"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type FixedChargeHandler struct {
	fixedChargeRepo fixed_charge.Repository
	templates       *template.Template
}

func NewFixedChargeHandler(
	fixedChargeRepo fixed_charge.Repository,
	templates *template.Template,
) *FixedChargeHandler {
	return &FixedChargeHandler{
		fixedChargeRepo: fixedChargeRepo,
		templates:       templates,
	}
}

func (h *FixedChargeHandler) ListFixedCharges(w http.ResponseWriter, r *http.Request) {
	charges, err := h.fixedChargeRepo.FindAll()
	if err != nil {
		http.Error(w, "Failed to get fixed charges", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Charges": charges,
	}

	if err := h.templates.ExecuteTemplate(w, "fixed_charges_list.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *FixedChargeHandler) AddFixedCharge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	money, err := shared.NewMoney(amount)
	if err != nil {
		http.Error(w, "Invalid amount: "+err.Error(), http.StatusBadRequest)
		return
	}

	charge := fixed_charge.NewFixedCharge(
		uuid.New().String(),
		name,
		money,
		description,
		true,
	)

	if err := h.fixedChargeRepo.Save(charge); err != nil {
		http.Error(w, "Failed to save fixed charge", http.StatusInternalServerError)
		return
	}

	h.ListFixedCharges(w, r)
}
