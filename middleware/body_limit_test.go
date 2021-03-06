package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/zqzca/echo"
	"github.com/zqzca/echo/test"
	"github.com/stretchr/testify/assert"
)

func TestBodyLimit(t *testing.T) {
	e := echo.New()
	hw := []byte("Hello, World!")
	req := test.NewRequest(echo.POST, "/", bytes.NewReader(hw))
	rec := test.NewResponseRecorder()
	c := e.NewContext(req, rec)
	h := func(c echo.Context) error {
		body, err := ioutil.ReadAll(c.Request().Body())
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, string(body))
	}

	// Based on content length (within limit)
	if assert.NoError(t, BodyLimit("2M")(h)(c)) {
		assert.Equal(t, http.StatusOK, rec.Status())
		assert.Equal(t, hw, rec.Body.Bytes())
	}

	// Based on content read (overlimit)
	he := BodyLimit("2B")(h)(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusRequestEntityTooLarge, he.Code)

	// Based on content read (within limit)
	req = test.NewRequest(echo.POST, "/", bytes.NewReader(hw))
	rec = test.NewResponseRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, BodyLimit("2M")(h)(c)) {
		assert.Equal(t, http.StatusOK, rec.Status())
		assert.Equal(t, "Hello, World!", rec.Body.String())
	}

	// Based on content read (overlimit)
	req = test.NewRequest(echo.POST, "/", bytes.NewReader(hw))
	rec = test.NewResponseRecorder()
	c = e.NewContext(req, rec)
	he = BodyLimit("2B")(h)(c).(*echo.HTTPError)
	assert.Equal(t, http.StatusRequestEntityTooLarge, he.Code)
}
