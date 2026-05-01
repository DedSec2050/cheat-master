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

func GetIncompleteLectures(course *models.APIResponse) []models.Lecture {
	var incomplete []models.Lecture

	for _, lesson := range course.Data.Lessons {
		for _, lec := range lesson.Lectures {
			if !lec.IsCompleted {
				incomplete = append(incomplete, lec)
			}
		}
	}

	return incomplete
}

func GetIncompleteByWeek(course *models.APIResponse) map[string][]models.Lecture {
	result := make(map[string][]models.Lecture)

	for _, lesson := range course.Data.Lessons {
		var incompleteLecs []models.Lecture
		for _, lec := range lesson.Lectures {
			if !lec.IsCompleted {
				incompleteLecs = append(incompleteLecs, lec)
			}
		}
		if len(incompleteLecs) > 0 {
			result[lesson.Name] = incompleteLecs
		}
	}

	return result
}

func GroupByWeek(course *models.APIResponse) map[string][]models.Lecture {
	result := make(map[string][]models.Lecture)

	for _, lesson := range course.Data.Lessons {
		result[lesson.Name] = lesson.Lectures
	}

	return result
}

func GetTotalLectures(course *models.APIResponse) int {
	count := 0
	for _, lesson := range course.Data.Lessons {
		count += len(lesson.Lectures)
	}
	return count
}

func GetCompletionStats(course *models.APIResponse) (completed, total int) {
	total = GetTotalLectures(course)
	completed = total - len(GetPendingLectures(course))
	return
}