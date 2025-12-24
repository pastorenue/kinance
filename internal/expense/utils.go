package expense

import (
	"time"
)

// Helper methods for recurring expenses
func (re *RecurringExpense) IsDue(currentDate time.Time) bool {
	return !re.NextDueDate.After(currentDate)
}

func (re *RecurringExpense) ShouldProcess() bool {
	if !re.IsActive {
		return false
	}

	if re.EndDate != nil && time.Now().After(*re.EndDate) {
		return false
	}

	return time.Now().After(re.NextDueDate) || time.Now().Equal(re.NextDueDate)
}

func (re *RecurringExpense) CalculateNextDueDate() {
	switch re.Frequency {
	case Daily:
		re.NextDueDate = re.NextDueDate.AddDate(0, 0, 1)
	case Weekly:
		re.NextDueDate = re.NextDueDate.AddDate(0, 0, 7)
	case Monthly:
		re.NextDueDate = re.NextDueDate.AddDate(0, 1, 0)
	case Yearly:
		re.NextDueDate = re.NextDueDate.AddDate(1, 0, 0)
	default:
		// Default to monthly if frequency is unknown
		re.NextDueDate = re.NextDueDate.AddDate(0, 1, 0)
	}
}

func (re *RecurringExpense) DaysUntilNextDue(currentDate time.Time) int {
	if re.NextDueDate.Before(currentDate) {
		return 0
	}
	return int(re.NextDueDate.Sub(currentDate).Hours() / 24)
}

func (re *RecurringExpense) HasEnded(currentDate time.Time) bool {
	if re.EndDate == nil {
		return false
	}
	return currentDate.After(*re.EndDate)
}

func (re *RecurringExpense) Activate() {
	re.IsActive = true
}

func (re *RecurringExpense) Deactivate() {
	re.IsActive = false
}

func (re *RecurringExpense) CancelReccurence() {
	re.EndDate = &time.Time{}
	re.IsActive = false
}
