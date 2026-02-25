package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aymaneelmaini/moka/internal/application"
	"github.com/aymaneelmaini/moka/internal/infrastructure/persistence/sqlite"
	"github.com/aymaneelmaini/moka/internal/infrastructure/web/handlers"
)

//go:embed ../internal/infrastructure/web/templates/*.html
var templatesFS embed.FS

//go:embed ../static/*
var staticFS embed.FS

//go:embed ../migrations/*.sql
var migrationsFS embed.FS

func main() {
	dataDir := os.Getenv("MOKA_DATA_DIR")
	if dataDir == "" {
		dataDir = "."
	}

	dbPath := filepath.Join(dataDir, "moka.db")

	log.Println("Initializing database...")
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Running migrations...")
	if err := db.RunMigrationsFromFS(migrationsFS, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Initializing repositories...")
	transactionRepo := sqlite.NewTransactionRepository(db)
	budgetRepo := sqlite.NewBudgetRepository(db)
	fixedChargeRepo := sqlite.NewFixedChargeRepository(db)
	loanRepo := sqlite.NewLoanRepository(db)

	log.Println("Initializing use cases...")
	addSalaryUC := application.NewAddSalaryUseCase(transactionRepo, fixedChargeRepo)
	recordExpenseUC := application.NewRecordExpenseUseCase(transactionRepo, budgetRepo)
	borrowMoneyUC := application.NewBorrowMoneyUseCase(loanRepo, transactionRepo)
	payLoanUC := application.NewPayLoanUseCase(loanRepo, transactionRepo)
	getMonthlySummaryUC := application.NewGetMonthlySummaryUseCase(transactionRepo, budgetRepo, loanRepo, fixedChargeRepo)

	log.Println("Loading templates...")
	tmpl, err := template.ParseFS(templatesFS, "internal/infrastructure/web/templates/*.html")
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	log.Println("Initializing handlers...")
	dashboardHandler := handlers.NewDashboardHandler(getMonthlySummaryUC, tmpl)
	transactionHandler := handlers.NewTransactionHandler(addSalaryUC, recordExpenseUC, tmpl)
	loanHandler := handlers.NewLoanHandler(borrowMoneyUC, payLoanUC, tmpl)
	fixedChargeHandler := handlers.NewFixedChargeHandler(fixedChargeRepo, tmpl)

	log.Println("Setting up routes...")
	mux := http.NewServeMux()

	staticFiles := http.FileServer(http.FS(staticFS))
	mux.Handle("/static/", staticFiles)

	mux.HandleFunc("/", dashboardHandler.ShowDashboard)

	mux.HandleFunc("/salary", transactionHandler.AddSalary)
	mux.HandleFunc("/expense", transactionHandler.RecordExpense)
	mux.HandleFunc("/loan/borrow", loanHandler.BorrowMoney)
	mux.HandleFunc("/loan/pay", loanHandler.PayLoan)
	mux.HandleFunc("/fixed-charges", fixedChargeHandler.ListFixedCharges)
	mux.HandleFunc("/fixed-charge/add", fixedChargeHandler.AddFixedCharge)

	port := ":9876"
	log.Printf("✨ Moka is running on http://moka.local%s", port)
	log.Println("")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
