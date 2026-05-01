package display

import (
	"fmt"
	"strings"

	"cheat-master/internal/models"
)

// IncompleteInfo holds a group of incomplete lectures by week
type IncompleteInfo struct {
	Week     string
	Lectures []models.Lecture
}

// DisplayIncompleteByWeek shows all incomplete lectures organized by week with icons
func DisplayIncompleteByWeek(course *models.APIResponse) {
	incomplete := getIncompleteByWeek(course)

	if len(incomplete) == 0 {
		fmt.Println("\n✅ All lectures completed!")
		return
	}

	fmt.Printf("\n📚 Incomplete Lectures for: %s\n", course.Data.Title)
	fmt.Printf("⏳ Progress: %s%%\n", course.Data.Progress)
	fmt.Println(strings.Repeat("─", 80))

	for _, group := range incomplete {
		fmt.Printf("\n📆 %s\n", group.Week)
		fmt.Println(strings.Repeat("  ─", 25))

		for idx, lec := range group.Lectures {
			icon := "⏲"
			numStr := fmt.Sprintf("[%02d]", idx+1)
			fmt.Printf("  %s %s 📹 %s\n", icon, numStr, lec.Title)
		}

		fmt.Printf("  📊 Total pending: %d\n", len(group.Lectures))
	}

	totalIncomplete := getTotalIncomplete(incomplete)
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("\n🎯 Total Incomplete Lectures: %d\n\n", totalIncomplete)
}

// DisplayIncompleteFull shows a detailed table format of incomplete lectures
func DisplayIncompleteFull(course *models.APIResponse) {
	incomplete := getIncompleteByWeek(course)

	if len(incomplete) == 0 {
		fmt.Println("\n✅ All lectures completed!")
		return
	}

	fmt.Printf("\n📚 Incomplete Lectures for: %s\n", course.Data.Title)
	fmt.Printf("⏳ Progress: %s%%\n\n", course.Data.Progress)

	// Table header
	fmt.Printf("%-10s %-40s %-30s\n", "ID", "Title", "Week")
	fmt.Println(strings.Repeat("─", 80))

	idx := 1
	for _, group := range incomplete {
		for _, lec := range group.Lectures {
			title := truncateString(lec.Title, 38)
			week := truncateString(group.Week, 28)
			fmt.Printf("%-10d %-40s %-30s\n", lec.ID, title, week)
			idx++
		}
	}

	totalIncomplete := getTotalIncomplete(incomplete)
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("Total: %d incomplete lectures\n\n", totalIncomplete)
}

// DisplayIncompleteCompact shows a compact list of incomplete lectures
func DisplayIncompleteCompact(course *models.APIResponse) {
	incomplete := getIncompleteByWeek(course)

	if len(incomplete) == 0 {
		fmt.Println("✅ All lectures completed!")
		return
	}

	fmt.Printf("📚 %s - ⏳ %s%%\n", course.Data.Title, course.Data.Progress)

	for _, group := range incomplete {
		fmt.Printf("  📆 %s (%d pending)\n", group.Week, len(group.Lectures))
		for i, lec := range group.Lectures {
			prefix := "├─"
			if i == len(group.Lectures)-1 {
				prefix = "└─"
			}
			fmt.Printf("    %s 📹 %s\n", prefix, lec.Title)
		}
	}

	totalIncomplete := getTotalIncomplete(incomplete)
	fmt.Printf("\n🎯 %d incomplete | ✅ %d total\n", totalIncomplete, getTotalLectures(course))
}

// getIncompleteByWeek returns incomplete lectures grouped by week
func getIncompleteByWeek(course *models.APIResponse) []IncompleteInfo {
	var result []IncompleteInfo

	for _, lesson := range course.Data.Lessons {
		var incompleteLecs []models.Lecture

		for _, lec := range lesson.Lectures {
			if !lec.IsCompleted {
				incompleteLecs = append(incompleteLecs, lec)
			}
		}

		if len(incompleteLecs) > 0 {
			result = append(result, IncompleteInfo{
				Week:     lesson.Name,
				Lectures: incompleteLecs,
			})
		}
	}

	return result
}

// getTotalIncomplete counts all incomplete lectures
func getTotalIncomplete(incomplete []IncompleteInfo) int {
	count := 0
	for _, group := range incomplete {
		count += len(group.Lectures)
	}
	return count
}

// getTotalLectures counts all lectures in the course
func getTotalLectures(course *models.APIResponse) int {
	count := 0
	for _, lesson := range course.Data.Lessons {
		count += len(lesson.Lectures)
	}
	return count
}

// truncateString truncates a string to a maximum length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
