package cheatmaster

import (
	"fmt"

	"cheat-master/internal/courses"
	"cheat-master/internal/display"
	"cheat-master/internal/models"
)

// Example 1: Display all incomplete lectures with default view
func ExampleDisplayIncomplete(course *models.APIResponse) {
	display.DisplayIncompleteByWeek(course)
}

// Example 2: Get incomplete lectures count
func ExampleGetIncompleteCount(course *models.APIResponse) {
	incompleteLectures := courses.GetPendingLectures(course)
	fmt.Printf("Total incomplete: %d\n", len(incompleteLectures))
}

// Example 3: Work with incomplete lectures grouped by week
func ExampleGroupIncompleteByWeek(course *models.APIResponse) {
	incompleteByWeek := courses.GetIncompleteByWeek(course)

	for week, lectures := range incompleteByWeek {
		fmt.Printf("%s: %d lectures\n", week, len(lectures))
		for _, lec := range lectures {
			fmt.Printf("  - %s\n", lec.Title)
		}
	}
}

// Example 4: Get all incomplete lectures as flat list
func ExampleGetAllIncomplete(course *models.APIResponse) {
	allIncomplete := courses.GetIncompleteLectures(course)
	fmt.Printf("Incomplete lectures:\n")
	for i, lec := range allIncomplete {
		fmt.Printf("%d. %s\n", i+1, lec.Title)
	}
}

// Example 5: Get completion statistics
func ExampleGetStats(course *models.APIResponse) {
	completed, total := courses.GetCompletionStats(course)
	percentage := float64(completed) / float64(total) * 100

	fmt.Printf("Completed: %d/%d (%.1f%%)\n", completed, total, percentage)
}

// Example 6: Display compact view
func ExampleCompactDisplay(course *models.APIResponse) {
	display.DisplayIncompleteCompact(course)
}

// Example 7: Display full table view
func ExampleFullTableDisplay(course *models.APIResponse) {
	display.DisplayIncompleteFull(course)
}

// Example 8: Show statistics summary
func ExampleStatistics(course *models.APIResponse) {
	display.StatisticsSummary(course)
	display.ProgressBar(course)
}

// Example 9: Conditional display based on progress
func ExampleConditionalDisplay(course *models.APIResponse) {
	completed, total := courses.GetCompletionStats(course)

	if completed == total {
		fmt.Println("✅ All lectures completed!")
	} else if completed == 0 {
		fmt.Println("⚠ No lectures completed yet")
		display.DisplayIncompleteByWeek(course)
	} else {
		fmt.Printf("📊 Progress: %d/%d\n", completed, total)
		display.DisplayIncompleteCompact(course)
	}
}

// Example 10: Filter and display specific week
func ExampleDisplayWeek(course *models.APIResponse, weekName string) {
	incompleteByWeek := courses.GetIncompleteByWeek(course)

	if lectures, exists := incompleteByWeek[weekName]; exists {
		fmt.Printf("📆 %s - %d incomplete lectures:\n", weekName, len(lectures))
		for i, lec := range lectures {
			fmt.Printf("  %d. %s\n", i+1, lec.Title)
		}
	} else {
		fmt.Printf("Week '%s' not found or all lectures completed\n", weekName)
	}
}

// Example 11: Find lectures by title pattern
func ExampleFindLectures(course *models.APIResponse, pattern string) {
	incomplete := courses.GetIncompleteLectures(course)

	found := 0
	for _, lec := range incomplete {
		if contains(lec.Title, pattern) {
			fmt.Printf("Found: %s\n", lec.Title)
			found++
		}
	}

	if found == 0 {
		fmt.Printf("No incomplete lectures matching '%s' found\n", pattern)
	}
}

// Example 12: Count incomplete per week
func ExampleIncompletePerWeek(course *models.APIResponse) {
	incompleteByWeek := courses.GetIncompleteByWeek(course)

	fmt.Println("Incomplete lectures per week:")
	for week, lectures := range incompleteByWeek {
		fmt.Printf("  %s: %d\n", week, len(lectures))
	}
}

// Example 13: Get IDs of incomplete lectures for batch operations
func ExampleGetIncompleteIDs(course *models.APIResponse) []int {
	return courses.GetPendingLectures(course)
}

// Example 14: Check if specific lecture is incomplete
func ExampleCheckLectureStatus(course *models.APIResponse, lectureID int) bool {
	pendingIDs := courses.GetPendingLectures(course)
	for _, id := range pendingIDs {
		if id == lectureID {
			return true // incomplete
		}
	}
	return false // completed
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) > 0)
}

/*
USAGE IN MAIN CODE:

Example in run.go:

    func Run(email, password string) {
        c := client.NewClient()
        if err := c.Login(email, password); err != nil {
            panic(err)
        }

        slugs, _ := courses.GetEnrollments(c)
        selected := selectCourse(slugs)
        course, _ := courses.GetCourse(c, selected)

        // Display incomplete lectures with icons
        display.DisplayIncompleteByWeek(course)

        // Ask for confirmation
        fmt.Print("\nProceed? (y/n): ")
        var proceed string
        fmt.Scanln(&proceed)

        if proceed != "y" {
            return
        }

        // Continue with execution
        ...
    }

IMPORT STATEMENTS:

    import (
        "cheat-master/internal/courses"
        "cheat-master/internal/display"
    )

DISPLAY OPTIONS:

1. Main (recommended):
   display.DisplayIncompleteByWeek(course)

2. Compact:
   display.DisplayIncompleteCompact(course)

3. Table:
   display.DisplayIncompleteFull(course)

4. Statistics:
   display.StatisticsSummary(course)
   display.ProgressBar(course)

DATA RETRIEVAL:

1. Get incomplete lecture IDs:
   ids := courses.GetPendingLectures(course)

2. Get incomplete lectures with details:
   lectures := courses.GetIncompleteLectures(course)

3. Get incomplete grouped by week:
   weekMap := courses.GetIncompleteByWeek(course)

4. Get statistics:
   completed, total := courses.GetCompletionStats(course)
*/
