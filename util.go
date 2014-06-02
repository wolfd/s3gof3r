package s3gof3r

import (
	"bytes"

	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// convenience multipliers
const (
	_        = iota
	kb int64 = 1 << (10 * iota)
	mb
	gb
	tb
	pb
	eb
)

// Min and Max functions
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// Error type and functions for http response
// http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
type respError struct {
	Code       string
	Message    string
	Resource   string
	RequestId  string
	StatusCode int
}

func newRespError(r *http.Response) *respError {
	e := new(respError)
	e.StatusCode = r.StatusCode
	b, _ := ioutil.ReadAll(r.Body)
	xml.NewDecoder(bytes.NewReader(b)).Decode(e) // parse error from response
	r.Body.Close()
	return e
}

func (e *respError) Error() string {
	return fmt.Sprintf(
		"%d: %q",
		e.StatusCode,
		e.Message,
	)
}

func bucketFromUrl(subdomain string) string {
	s := strings.Split(subdomain, ".")
	return strings.Join(s[:len(s)-1], ".")
}

func checkClose(c io.Closer, err *error) {
	cerr := c.Close()
	if *err == nil {
		*err = cerr
	}
}
