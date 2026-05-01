package orchestrator

import (
	"fmt"
	"sync"

	"cheat-master/internal/client"
	"cheat-master/internal/courses"
	"cheat-master/internal/display"
	"cheat-master/internal/executor"
	"cheat-master/internal/models"
)


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

	fmt.Println("\n🎯 Starting execution (Concurrent Mode)...")

	// Keep track of lectures to process
	lecturesToProcess := lectures
	attempt := 0
	maxPolls := 6

	for len(lecturesToProcess) > 0 {
		var nextRound []models.Lecture
		var mu sync.Mutex
		var wg sync.WaitGroup

		for _, lec := range lecturesToProcess {
			if lec.IsCompleted {
				continue
			}

			wg.Add(1)
			go func(lecture models.Lecture) {
				defer wg.Done()

				// Create separate client for this lecture
				lectureClient := client.NewClient()
				if err := lectureClient.Login(email, password); err != nil {
					fmt.Printf("❌ Login failed for %s: %v\n", lecture.Title, err)
					return
				}

				fmt.Printf("\n▶ Watching [Goroutine]: %s\n", lecture.Title)

				// Try 5-6 polls for this lecture
				for poll := 0; poll < maxPolls; poll++ {
					executor.WatchLecture(lectureClient, selected, lecture.ID, 200, email, password)

					updated, _ := courses.GetCourse(lectureClient, selected)

					if courses.IsLectureCompleted(updated, lecture.ID) {
						fmt.Printf("✅ Completed [Goroutine]: %s\n", lecture.Title)
						return
					}

					if poll < maxPolls-1 {
						fmt.Printf("⏳ Poll %d/%d: Still incomplete, retrying... [%s]\n", poll+1, maxPolls, lecture.Title)
					}
				}

				// Check again after polls
				updated, _ := courses.GetCourse(lectureClient, selected)
				if !courses.IsLectureCompleted(updated, lecture.ID) {
					// Still incomplete, add to next round
					mu.Lock()
					nextRound = append(nextRound, lecture)
					mu.Unlock()
					fmt.Printf("⚠ Moving to next round, will retry: %s\n", lecture.Title)
				}
			}(lec)
		}

		// Wait for all goroutines to complete
		wg.Wait()

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

