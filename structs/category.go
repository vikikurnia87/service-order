package structs

// CategoryCreateRequest = payload buat kategori baru (ter-scope company dari JWT).
type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

// CategoryUpdateRequest = payload ubah nama dan/atau deskripsi kategori.
type CategoryUpdateRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}
