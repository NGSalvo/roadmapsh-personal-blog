package utils

import "fmt"

func PostSSEWithCSRF(url string, csrfToken string) string {
	return fmt.Sprintf("sse('%s', {method: 'post', headers: {'x-csrf-token': '%s'}})", url, csrfToken)
}
