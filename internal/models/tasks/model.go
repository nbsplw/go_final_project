package tasks

type Task struct {
	ID      string `db:"id" json:"id" binding:"required"`
	Date    string `db:"date" json:"date" binding:"required"`
	Title   string `db:"title" json:"title" binding:"required"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat,omitempty"`
}
