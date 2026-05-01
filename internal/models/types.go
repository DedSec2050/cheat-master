package models

type Lecture struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"is_completed"`
}

type Lesson struct {
	Name     string    `json:"name"`
	Lectures []Lecture `json:"lectures"`
}

type Course struct {
	Title    string   `json:"title"`
	Progress string   `json:"progress_bar"`
	Lessons  []Lesson `json:"lessons"`
}

type APIResponse struct {
	Data Course `json:"data"`
}