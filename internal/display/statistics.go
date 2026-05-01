package display

import (
	"fmt"

	"cheat-master/internal/models"
)

// StatisticsSummary displays course statistics
func StatisticsSummary(course *models.APIResponse) {
	completed, total := getStats(course)
	percentage := float64(completed) / float64(total) * 100

	fmt.Printf("\n📊 Course Statistics: %s\n", course.Data.Title)
	fmt.Printf("   ✅ Completed: %d/%d (%.1f%%)\n", completed, total, percentage)
	fmt.Printf("   ⏲ Remaining: %d\n", total-completed)
	fmt.Printf("   ⏳ Progress: %s%%\n", course.Data.Progress)
}

// ProgressBar displays a simple ASCII progress bar
func ProgressBar(course *models.APIResponse) {
	completed, total := getStats(course)
	percentage := float64(completed) / float64(total) * 100

	barLength := 30
	filledLength := int(percentage / 100 * float64(barLength))

	bar := "["
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			bar += "="
		} else {
			bar += "-"
		}
	}
	bar += "]"

	fmt.Printf("\n%s %.1f%% (%d/%d)\n", bar, percentage, completed, total)
}

// WeekSummary shows summary for a specific week
func WeekSummary(course *models.APIResponse, weekName string) {
	for _, lesson := range course.Data.Lessons {
		if lesson.Name != weekName {
			continue
		}

		total := len(lesson.Lectures)
		completed := 0

		for _, lec := range lesson.Lectures {
			if lec.IsCompleted {
				completed++
			}
		}

		percentage := float64(completed) / float64(total) * 100

		fmt.Printf("\n📆 %s Summary\n", weekName)
		fmt.Printf("   ✅ Completed: %d/%d (%.1f%%)\n", completed, total, percentage)
		fmt.Printf("   ⏲ Remaining: %d\n", total-completed)
		return
	}
}

// getStats calculates completion statistics
func getStats(course *models.APIResponse) (completed, total int) {
	total = 0
	completed = 0

	for _, lesson := range course.Data.Lessons {
		for _, lec := range lesson.Lectures {
			total++
			if lec.IsCompleted {
				completed++
			}
		}
	}

	return
}
