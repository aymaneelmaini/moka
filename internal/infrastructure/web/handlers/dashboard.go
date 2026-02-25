package handlers

import (
	"html/template"
	"github.com/aymaneelmaini/moka/internal/application"
	"net/http"
	"strconv"
	"time"
)

type DashboardHandler struct {
	getMonthlySummaryUC *application.GetMonthlySummaryUseCase
	templates           *template.Template
}

func NewDashboardHandler(
	getMonthlySummaryUC *application.GetMonthlySummaryUseCase,
	templates *template.Template,
) *DashboardHandler {
	return &DashboardHandler{
		getMonthlySummaryUC: getMonthlySummaryUC,
		templates:           templates,
	}
}

func (h *DashboardHandler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	year := now.Year()
	month := now.Month()

	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	if monthStr := r.URL.Query().Get("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = time.Month(m)
		}
	}

	summary, err := h.getMonthlySummaryUC.Execute(application.GetMonthlySummaryInput{
		Year:  year,
		Month: month,
	})

	if err != nil {
		http.Error(w, "Failed to get monthly summary", http.StatusInternalServerError)
		return
	}

	// Calculate previous and next month/year
	currentDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)
	nextDate := currentDate.AddDate(0, 1, 0)

	data := map[string]interface{}{
		"Title":     "Moka - Dashboard",
		"Summary":   summary,
		"PrevYear":  prevDate.Year(),
		"PrevMonth": int(prevDate.Month()),
		"NextYear":  nextDate.Year(),
		"NextMonth": int(nextDate.Month()),
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
