package main

import (
	"html/template"
	"log"
	"net/http"

	"moka/internal/application"
	"moka/internal/infrastructure/persistence/sqlite"
	"moka/internal/infrastructure/web/handlers"
)

func main() {

	log.Println("Initializing database...")
	db, err := sqlite.NewDB("./moka.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Running migrations...")
	if err := db.RunMigrations("./migrations"); err != nil {
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
	getMonthlySummaryUC := application.NewGetMonthlySummaryUseCase(transactionRepo, budgetRepo, loanRepo)

	log.Println("Loading templates...")
	tmpl, err := template.ParseGlob("./internal/infrastructure/web/templates/*.html")
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

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", dashboardHandler.ShowDashboard)

	mux.HandleFunc("/salary", transactionHandler.AddSalary)
	mux.HandleFunc("/expense", transactionHandler.RecordExpense)
	mux.HandleFunc("/loan/borrow", loanHandler.BorrowMoney)
	mux.HandleFunc("/loan/pay", loanHandler.PayLoan)
	mux.HandleFunc("/fixed-charges", fixedChargeHandler.ListFixedCharges)
	mux.HandleFunc("/fixed-charge/add", fixedChargeHandler.AddFixedCharge)

	port := ":8080"
	log.Printf("✨ Moka is running on http://localhost%s", port)
	log.Println("")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
