package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang-url-shortener/internal/http-server/handlers/url/save"
	"golang-url-shortener/internal/lib/api/response"
	"io"
	"net/http"
)

const contentType = "application/json"

func (s *UrlShortenerE2ESuite) TestSaveAndRedirectSuccess() {
	url := fmt.Sprintf("%s/url", s.server.URL)

	testURL := "https://api.kanye.rest"
	testAlias := "kanye"

	req := save.Request{
		URL:   testURL,
		Alias: testAlias,
	}

	marshalledReq, err := json.Marshal(req)
	s.test.NoError(err)

	saveReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalledReq))
	s.test.NoError(err)
	saveReq.Header.Set("Content-Type", contentType)

	saveResp, err := s.httpClient.Do(saveReq)
	s.test.NoError(err)
	s.test.Equal(http.StatusOK, saveResp.StatusCode)
	defer saveResp.Body.Close()

	getURL := fmt.Sprintf("%s/%s", s.server.URL, testAlias)
	getResp, err := s.httpClient.Get(getURL)
	s.test.NoError(err)
	defer getResp.Body.Close()

	// Проверяем, что произошел редирект
	s.test.Equal(getResp.Request.URL.String(), testURL)
}

func (s *UrlShortenerE2ESuite) TestSaveAndDeleteSuccess() {
	url := fmt.Sprintf("%s/url", s.server.URL)

	testURL := "https://api.kanye.rest"
	testAlias := "kanye"

	req := save.Request{
		URL:   testURL,
		Alias: testAlias,
	}

	marshalledReq, err := json.Marshal(req)
	s.test.NoError(err)

	saveReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalledReq))
	s.test.NoError(err)
	saveReq.Header.Set("Content-Type", contentType)

	saveResp, err := s.httpClient.Do(saveReq)
	s.test.NoError(err)
	s.test.Equal(http.StatusOK, saveResp.StatusCode)
	defer saveResp.Body.Close()

	// Удаляем url и alias
	deleteURL := fmt.Sprintf("%s/%s", url, testAlias)
	deleteReq, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	s.test.NoError(err)
	deleteReq.Header.Set("Content-Type", contentType)

	deleteResp, err := s.httpClient.Do(deleteReq)
	s.test.NoError(err)

	body, err := io.ReadAll(deleteResp.Body)
	s.test.NoError(err)

	respCore := &response.Response{}
	s.test.NoError(json.Unmarshal(body, respCore))
	s.test.Equal(respCore.Status, response.StatusOK)
	s.test.Equal(respCore.Error, "")

	// Пытаемся перейти по удаленному alias - получаем ошибку
	getURL := fmt.Sprintf("%s/%s", s.server.URL, testAlias)
	getResp, err := s.httpClient.Get(getURL)
	s.test.NoError(err)
	defer getResp.Body.Close()

	respBody, err := io.ReadAll(getResp.Body)
	s.test.NoError(err)
	s.test.Equal("\"url not found\"\n", string(respBody))
}
