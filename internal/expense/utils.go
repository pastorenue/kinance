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
	
	if re.EndDate == nil && time.Now().After(*re.EndDate) {
		return false
	}

	return time.Now().After(re.NextDueDate) || time.Now().Equal(re.NextDueDate)
}

func (re *RecurringExpense) CalculateNextDueDate() time.Time {
	switch re.Frequency {
	case Daily:
		return re.NextDueDate.AddDate(0, 0, 1)
	case Weekly:
		return re.NextDueDate.AddDate(0, 0, 7)
	case Monthly:
		return re.NextDueDate.AddDate(0, 1, 0)
	case Yearly:
		return re.NextDueDate.AddDate(1, 0, 0)
	default:
		return re.NextDueDate // No change for unknown frequency
	}
}
