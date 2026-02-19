package domain

import "time"

type Discount struct {
	percent   int64
	startDate time.Time
	endDate   time.Time
}

func NewDiscount(percent int64, startDate, endDate time.Time) *Discount {
	if percent < 0 || percent > 100 {
		return nil
	}
	if endDate.Before(startDate) {
		return nil
	}
	return &Discount{
		percent:   percent,
		startDate: startDate,
		endDate:   endDate,
	}
}

func (d *Discount) Percentage() int64 {
	if d == nil {
		return 0
	}
	return d.percent
}

func (d *Discount) StartDate() time.Time {
	if d == nil {
		return time.Time{}
	}
	return d.startDate
}

func (d *Discount) EndDate() time.Time {
	if d == nil {
		return time.Time{}
	}
	return d.endDate
}

func (d *Discount) IsValidAt(t time.Time) bool {
	if d == nil {
		return false
	}
	return !t.Before(d.startDate) && !t.After(d.endDate)
}
