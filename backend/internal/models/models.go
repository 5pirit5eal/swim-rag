package models

type User struct {
	ID string `db:"user_id"`
	// Name is the user's name as provided by the oauth provider
	Name string `db:"name"`
	// Email is the user's email address as provided by the oauth provider
	Email string `db:"email"`
	// CreatedAt is the time the user was created as a datetime string
	CreatedAt string `db:"created_at"`
	// Last time the user was active as a datetime string
	LastActive string `db:"last_active"`
}

type Feedback struct {
	UserID    string `db:"user_id"`
	PlanID    string `db:"plan_id"`
	Rating    int    `db:"rating"`
	Comment   string `db:"comment"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type ChoiceResult struct {
	Idx         int    `json:"index" example:"1"`
	Description string `json:"description" example:"Selected plan based on your requirements"`
}
