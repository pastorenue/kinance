package expense

type Handler struct {
	*ExpenseHandler
	*CategoryHandler
}

func NewHandler(expenseService *ExpenseService, categoryService *CategoryService) *Handler {
	return &Handler{
		ExpenseHandler:  NewExpenseHandler(expenseService),
		CategoryHandler: NewCategoryHandler(categoryService),
	}
}
