package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
)

type ErrHTTP struct {
	Code   int
	Status string
}

func (e *ErrHTTP) Error() string {
	return fmt.Sprintf("bad status: %s", e.Status)
}

func CheckResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return &ErrHTTP{
			Status: resp.Status,
			Code:   resp.StatusCode,
		}
	}
	return nil
}

func HTTPGet(url string) ([]byte, error) {
	resp, err := retryablehttp.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Check server response
	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)

	return data, err
}

func DownloadFile(filepath, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	// Get the data
	resp, err := retryablehttp.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
