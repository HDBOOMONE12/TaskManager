package taskclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type TaskClient struct {
	BaseURL string
	Client  *http.Client
}

func NewTaskClient(baseURL string) *TaskClient {
	return &TaskClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *TaskClient) HasUserWithEmail(ctx context.Context, email string) (bool, error) {
	u, _ := url.Parse(c.BaseURL + "/users")
	q := u.Query()
	q.Set("email", email)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("GET %s → status: %d", u.String(), resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)

}
