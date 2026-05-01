package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cheat-master/internal/client"
)

func WatchLecture(c *client.VTUClient, slug string, lectureID int, total int) error {
	url := fmt.Sprintf("%s/student/my-courses/%s/lectures/%d/progress",
		c.Base, slug, lectureID)

	chunk := total / 50 // simulate chunks
	current := 0

	for i := 0; i < 10; i++ { // polling loop
		current += chunk
		if current > total {
			current = total
		}

		payload := map[string]int{
			"current_time_seconds":  current,
			"total_duration_seconds": total,
			"seconds_just_watched":  chunk,
		}

		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.Client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		fmt.Printf("⏱ Lecture %d progress: %d/%d\n", lectureID, current, total)

		if current == total {
			break
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}
