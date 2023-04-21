package activity

type CreateActivity struct {
	Title string `json:"title"`
	Email string `json:"email"`
}

type UpdateActivity struct {
	Title string `json:"title"`
}

type Activity struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
