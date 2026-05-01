package orchestrator

import (
	"fmt"

	"cheat-master/internal/client"
	"cheat-master/internal/courses"
	"cheat-master/internal/display"
	"cheat-master/internal/executor"
	"cheat-master/internal/models"
)

// func Run(email, password string) {
// 	c := client.NewClient()

// 	// 1. Login
// 	err := c.Login(email, password)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 2. Get enrollments
// 	slugs, err := courses.GetEnrollments(c)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Courses:", slugs)

// 	// 3. Process each course
// 	for _, slug := range slugs {
// 		fmt.Println("\nProcessing:", slug)

// 		course, err := courses.GetCourse(c, slug)
// 		if err != nil {
// 			continue
// 		}

// 		pending := courses.GetPendingLectures(course)

// 		fmt.Println("Pending lectures:", len(pending))

// 		for _, lecID := range pending {
// 			executor.MarkComplete(c, slug, lecID)
// 			time.Sleep(2 * time.Second)
// 		}
// 	}

// 	fmt.Println("\n🎯 Done")
// }


func Run(email, password string) {
	c := client.NewClient()

	// Login
	if err := c.Login(email, password); err != nil {
		panic(err)
	}

	// Get courses
	slugs, _ := courses.GetEnrollments(c)

	// Select course
	selected := selectCourse(slugs)

	// Fetch course
	course, _ := courses.GetCourse(c, selected)

	// Display incomplete lectures with icons
	display.DisplayIncompleteByWeek(course)

	// Ask if user wants to proceed
	fmt.Print("\nProceed with execution? (y/n): ")
	var proceed string
	fmt.Scanln(&proceed)

	if proceed != "y" && proceed != "Y" {
		fmt.Println("Cancelled.")
		return
	}

	// Get only weeks with incomplete lectures
	incompleteWeeks := courses.GetIncompleteByWeek(course)

	// Select week
	lectures := selectWeek(incompleteWeeks)

	fmt.Println("\n🎯 Starting execution...")

	// Keep track of lectures to process
	lecturesToProcess := lectures
	attempt := 0
	maxPolls := 6

	for len(lecturesToProcess) > 0 {
		var nextRound []models.Lecture

		for _, lec := range lecturesToProcess {
			if lec.IsCompleted {
				continue
			}

			fmt.Println("\n▶ Watching:", lec.Title)

			// Try 5-6 polls for this lecture
			for poll := 0; poll < maxPolls; poll++ {
				executor.WatchLecture(c, selected, lec.ID, 200, email, password)

				updated, _ := courses.GetCourse(c, selected)

				if courses.IsLectureCompleted(updated, lec.ID) {
					fmt.Println("✅ Completed:", lec.Title)
					break
				}

				if poll < maxPolls-1 {
					fmt.Printf("⏳ Poll %d/%d: Still incomplete, retrying...\n", poll+1, maxPolls)
				}
			}

			// Check again after polls
			updated, _ := courses.GetCourse(c, selected)
			if !courses.IsLectureCompleted(updated, lec.ID) {
				// Still incomplete, add to next round
				nextRound = append(nextRound, lec)
				fmt.Printf("⚠ Moving to next lecture, will retry: %s\n", lec.Title)
			}
		}

		// Prepare for next round
		lecturesToProcess = nextRound
		attempt++

		if len(lecturesToProcess) > 0 {
			fmt.Printf("\n📊 Round %d: %d lectures still incomplete, retrying...\n", attempt, len(lecturesToProcess))
		}
	}

	fmt.Println("\n✅ All lectures completed!")
}



func selectCourse(slugs []string) string {
	fmt.Println("\nSelect Course:")
	for i, s := range slugs {
		fmt.Printf("[%d] %s\n", i+1, s)
	}

	var choice int
	fmt.Print("Enter choice: ")
	fmt.Scanln(&choice)

	return slugs[choice-1]
}

func selectWeek(weeks map[string][]models.Lecture) []models.Lecture {
	fmt.Println("\nSelect Week:")

	keys := []string{}
	i := 1
	for k := range weeks {
		fmt.Printf("[%d] %s\n", i, k)
		keys = append(keys, k)
		i++
	}

	var choice int
	fmt.Print("Enter choice: ")
	fmt.Scanln(&choice)

	return weeks[keys[choice-1]]
}

