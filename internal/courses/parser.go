package courses

import "cheat-master/internal/models"

func GetPendingLectures(course *models.APIResponse) []int {
	var ids []int

	for _, lesson := range course.Data.Lessons {
		for _, lec := range lesson.Lectures {
			if !lec.IsCompleted {
				ids = append(ids, lec.ID)
			}
		}
	}

	return ids
}

func GroupByWeek(course *models.APIResponse) map[string][]models.Lecture {
	result := make(map[string][]models.Lecture)

	for _, lesson := range course.Data.Lessons {
		result[lesson.Name] = lesson.Lectures
	}

	return result
}