package utils

import "errors"

// Sentinel error domain service-order — dipusatkan di sini agar service dan
// handler memakai instance yang sama (mendukung errors.Is lintas layer).
var (
	// ErrNotFound = order tidak ditemukan / bukan milik company.
	ErrNotFound = errors.New("order not found")
	// ErrPriorityNotFound = priority_uuid tidak dikenal.
	ErrPriorityNotFound = errors.New("priority not found")
	// ErrStatusNotFound = status_uuid tidak dikenal.
	ErrStatusNotFound = errors.New("status not found")
	// ErrCategoryNotFound = salah satu category_uuid tidak dikenal / bukan milik company.
	ErrCategoryNotFound = errors.New("category not found")
	// ErrCategoryEntityNotFound = kategori tidak ditemukan atau bukan milik company.
	ErrCategoryEntityNotFound = errors.New("category not found")
	// ErrCategoryNameExists = nama kategori sudah ada dalam company yang sama.
	ErrCategoryNameExists = errors.New("category name already exists")
	// ErrDayNotFound = day tidak ditemukan.
	ErrDayNotFound = errors.New("day not found")
	// ErrDateNotFound = date tidak ditemukan.
	ErrDateNotFound = errors.New("date not found")
	// ErrMasterDateNotFound = master date tidak ditemukan.
	ErrMasterDateNotFound = errors.New("master date not found")
	// ErrOrderPriorityNotFound = order priority tidak ditemukan.
	ErrOrderPriorityNotFound = errors.New("order priority not found")
	// ErrOrderStatusNotFound = order status tidak ditemukan.
	ErrOrderStatusNotFound = errors.New("order status not found")
	// ErrScheduleNotFound = schedule tidak ditemukan.
	ErrScheduleNotFound = errors.New("schedule not found")
)
