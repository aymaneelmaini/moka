package shared

import "errors"

type Category struct {
	name string
	typ  CategoryType
}

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

var (
	CategorySalary        = Category{name: "Salary", typ: CategoryTypeIncome}
	CategoryBorrowed      = Category{name: "Borrowed (Salaf)", typ: CategoryTypeIncome}
	CategoryFood          = Category{name: "Food", typ: CategoryTypeExpense}
	CategoryTransport     = Category{name: "Transport", typ: CategoryTypeExpense}
	CategoryEntertainment = Category{name: "Entertainment", typ: CategoryTypeExpense}
	CategoryShopping      = Category{name: "Shopping", typ: CategoryTypeExpense}
	CategoryHealth        = Category{name: "Health", typ: CategoryTypeExpense}
	CategoryOther         = Category{name: "Other", typ: CategoryTypeExpense}
)

func NewCategory(name string, typ CategoryType) (Category, error) {
	if name == "" {
		return Category{}, errors.New("category name cannot be empty")
	}

	if typ != CategoryTypeIncome && typ != CategoryTypeExpense {
		return Category{}, errors.New("invalid category type")
	}

	return Category{
		name: name,
		typ:  typ,
	}, nil
}

func (c Category) Name() string {
	return c.name
}

func (c Category) Type() CategoryType {
	return c.typ
}

func (c Category) IsIncome() bool {
	return c.typ == CategoryTypeIncome
}

func (c Category) IsExpense() bool {
	return c.typ == CategoryTypeExpense
}

func (c Category) String() string {
	return c.name
}
