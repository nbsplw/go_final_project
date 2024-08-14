package common

type AddTask struct {
	Date    string `json:"date,omitempty" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}
