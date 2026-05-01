package courses

import (
	"cheat-master/internal/client"
	"cheat-master/internal/models"
	"encoding/json"
	"net/http"
)

type EnrollmentResponse struct {
	Data []struct {
		Details struct {
			Slug string `json:"slug"`
		} `json:"details"`
	} `json:"data"`
}

func GetEnrollments(c *client.VTUClient) ([]string, error) {
	req, _ := http.NewRequest("GET", c.Base+"/student/my-enrollments", nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data EnrollmentResponse
	json.NewDecoder(resp.Body).Decode(&data)

	var slugs []string
	for _, e := range data.Data {
		slugs = append(slugs, e.Details.Slug)
	}

	return slugs, nil
}


func IsLectureCompleted(course *models.APIResponse, lectureID int) bool {
	for _, lesson := range course.Data.Lessons {
		for _, lec := range lesson.Lectures {
			if lec.ID == lectureID {
				return lec.IsCompleted
			}
		}
	}
	return false
}


