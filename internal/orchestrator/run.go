package orchestrator

import (
	"fmt"

	"cheat-master/internal/client"
	"cheat-master/internal/courses"
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

	// Group weeks
	weeks := courses.GroupByWeek(course)

	// Select week
	lectures := selectWeek(weeks)

	fmt.Println("\n🎯 Starting execution...")

	for _, lec := range lectures {
		if lec.IsCompleted {
			continue
		}

		fmt.Println("\n▶ Watching:", lec.Title)

		// simulate watch
		executor.WatchLecture(c, selected, lec.ID, 200)

		// verify
		updated, _ := courses.GetCourse(c, selected)

		if courses.IsLectureCompleted(updated, lec.ID) {
			fmt.Println("✅ Completed:", lec.Title)
		} else {
			fmt.Println("⚠ Still incomplete:", lec.Title)
		}
	}
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

