package courses

import (
	"cheat-master/internal/client"
	"cheat-master/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCourse(c *client.VTUClient, slug string) (*models.APIResponse, error) {
	url := fmt.Sprintf("%s/student/my-courses/%s/", c.Base, slug)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.APIResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return &result, nil
}