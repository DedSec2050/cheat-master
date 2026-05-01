package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"cheat-master/internal/client"
)

func WatchLecture(c *client.VTUClient, slug string, lectureID int, total int, email string, password string) error {
	url := fmt.Sprintf("%s/student/my-courses/%s/lectures/%d/progress",
		c.Base, slug, lectureID)

	chunk := total / 5 // simulate chunks
	current := 0

	for i := 0; i < 10; i++ { // polling loop
		current += chunk
		if current > total {
			current = total
		}

		payload := map[string]int{
			"current_time_seconds":  total,  // Incrementally increasing
			"total_duration_seconds": total,
			"seconds_just_watched":  total,   // Amount watched in this poll
		}

		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "Go-Http-Client")

		resp, err := c.Client.Do(req)
		if err != nil {
			return err
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Log response for debugging
		fmt.Printf("📤 Response Status: %d\n", resp.StatusCode)
		if len(respBody) > 0 {
			fmt.Printf("📥 Response Body: %s\n", string(respBody))
		}

		// Check for rate limiting (429 status code)
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			waitTime := 60 // default wait time in seconds

			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					waitTime = seconds
				}
			}

			fmt.Printf("⚠ Rate Limited! Status: %d - Waiting %d seconds...\n", resp.StatusCode, waitTime)
			fmt.Printf("   Retry-After: %s\n", retryAfter)
			time.Sleep(time.Duration(waitTime) * time.Second)
			continue // Retry this poll
		}

		// Check for other HTTP errors
		if resp.StatusCode >= 400 {
			fmt.Printf("❌ HTTP Error %d: %s\n", resp.StatusCode, string(respBody))

			// Check rate limit headers even on other errors
			remaining := resp.Header.Get("X-RateLimit-Remaining")
			limit := resp.Header.Get("X-RateLimit-Limit")
			reset := resp.Header.Get("X-RateLimit-Reset")

			if remaining != "" || limit != "" || reset != "" {
				fmt.Printf("📊 Rate Limit Info:\n")
				if limit != "" {
					fmt.Printf("   Limit: %s requests\n", limit)
				}
				if remaining != "" {
					fmt.Printf("   Remaining: %s requests\n", remaining)
				}
				if reset != "" {
					fmt.Printf("   Reset at: %s\n", reset)
				}
			}

			if resp.StatusCode == http.StatusTooManyRequests {
				return fmt.Errorf("rate limited")
			}
			return fmt.Errorf("HTTP %d error", resp.StatusCode)
		}

		// Check rate limit headers on successful responses
		remaining := resp.Header.Get("X-RateLimit-Remaining")
		if remaining != "" {
			fmt.Printf("📊 Requests Remaining: %s\n", remaining)
		}

		fmt.Printf("⏱ Lecture %d progress: %d/%d\n", lectureID, current, total)

		if current == total {
			break
		}

		// Re-login to get fresh access token
		fmt.Printf("🔄 Getting fresh access token...\n")
		if err := c.Login(email, password); err != nil {
			fmt.Printf("⚠ Re-login failed: %v\n", err)
			// Continue anyway with current token
		} else {
			fmt.Printf("✅ Fresh token obtained\n")
		}

		// Random delay between 1000-3000 milliseconds (1-3 seconds)
		randomDelay := time.Duration(rand.Intn(2001)+1000) * time.Millisecond
		fmt.Printf("⏳ Random delay: %dms\n", randomDelay.Milliseconds())
		time.Sleep(randomDelay)
	}

	return nil
}
