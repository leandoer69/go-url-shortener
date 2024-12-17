package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang-url-shortener/internal/http-server/handlers/url/save"
	"golang-url-shortener/internal/http-server/handlers/url/update"
	"golang-url-shortener/internal/lib/api/response"
	"golang-url-shortener/internal/storage"
	"io"
	"net/http"
)

const contentType = "application/json"

func (s *UrlShortenerSuite) TestSaveSuccess() {
	url := fmt.Sprintf("%s/url", s.server.URL)

	testURL := "https://mail.google.com/"
	testAlias := "mail"

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

	actualURL, err := s.storage.GetURL(testAlias)
	s.test.NoError(err)
	s.test.Equal(testURL, actualURL)
}

func (s *UrlShortenerSuite) TestSaveFailed_ErrorAlreadyExists() {
	url := fmt.Sprintf("%s/url", s.server.URL)

	testURL := "https://mail.google.com/"
	testAlias := "mail"

	req := save.Request{
		URL:   testURL,
		Alias: testAlias,
	}

	marshalledReq, err := json.Marshal(req)
	s.test.NoError(err)

	// Первая попытка - сохранили успешно
	firstSaveReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalledReq))
	s.test.NoError(err)
	firstSaveReq.Header.Set("Content-Type", contentType)

	firstSaveResp, err := s.httpClient.Do(firstSaveReq)
	s.test.NoError(err)
	s.test.Equal(http.StatusOK, firstSaveResp.StatusCode)
	defer firstSaveResp.Body.Close()

	// Вторая попытка - неудачная вставка
	secondSaveReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalledReq))
	s.test.NoError(err)
	secondSaveReq.Header.Set("Content-Type", contentType)

	secondSaveResp, err := s.httpClient.Do(secondSaveReq)
	s.test.NoError(err)
	s.test.Equal(http.StatusOK, secondSaveResp.StatusCode)
	defer secondSaveResp.Body.Close()

	body, err := io.ReadAll(secondSaveResp.Body)
	s.test.NoError(err)

	respCore := &response.Response{}
	s.test.NoError(json.Unmarshal(body, respCore))
	s.test.Equal(respCore.Status, response.StatusError)
	s.test.Equal(respCore.Error, "url already exists")

	actualURL, err := s.storage.GetURL(testAlias)
	s.test.NoError(err)
	s.test.Equal(testURL, actualURL)
}

func (s *UrlShortenerSuite) TestUpdateSuccess() {
	url := fmt.Sprintf("%s/url", s.server.URL)

	testURL := "https://mail.google.com/"
	testAlias := "mail"
	testNewAlias := "moil"

	req := save.Request{
		URL:   testURL,
		Alias: testAlias,
	}

	marshalledSaveReq, err := json.Marshal(req)
	s.test.NoError(err)

	// Вставляем url и alias
	saveReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalledSaveReq))
	s.test.NoError(err)
	saveReq.Header.Set("Content-Type", contentType)

	saveResp, err := s.httpClient.Do(saveReq)
	s.test.NoError(err)
	s.test.Equal(http.StatusOK, saveResp.StatusCode)
	defer saveResp.Body.Close()

	// Проверяем, что url и alias вставились
	actualURL, err := s.storage.GetURL(testAlias)
	s.test.NoError(err)
	s.test.Equal(testURL, actualURL)

	updateRequest := update.Request{
		URL:      testURL,
		OldAlias: testAlias,
		NewAlias: testNewAlias,
	}

	marshalledUpdateReq, err := json.Marshal(updateRequest)
	s.test.NoError(err)

	// Обновляем alias
	updateReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(marshalledUpdateReq))
	s.test.NoError(err)
	updateReq.Header.Set("Content-Type", contentType)

	updateResp, err := s.httpClient.Do(updateReq)
	s.test.NoError(err)

	body, err := io.ReadAll(updateResp.Body)
	s.test.NoError(err)

	respCore := &response.Response{}
	s.test.NoError(json.Unmarshal(body, respCore))
	s.test.Equal(respCore.Status, response.StatusOK)
	s.test.Equal(respCore.Error, "")

	// Проверяем, что alias обновился
	_, err = s.storage.GetURL(testAlias)
	s.test.ErrorIs(err, storage.ErrUrlNotFound)

	actualURL, err = s.storage.GetURL(testNewAlias)
	s.test.NoError(err)
	s.test.Equal(testURL, actualURL)
}

func (s *UrlShortenerSuite) TestUpdateFailed_ErrorUrlNotFound() {
	url := fmt.Sprintf("%s/url", s.server.URL)

	testURL := "https://mail.google.com/"
	testAlias := "mail"
	testNewAlias := "moil"

	updateRequest := update.Request{
		URL:      testURL,
		OldAlias: testAlias,
		NewAlias: testNewAlias,
	}

	marshalledUpdateReq, err := json.Marshal(updateRequest)
	s.test.NoError(err)

	// Обновляем alias
	updateReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(marshalledUpdateReq))
	s.test.NoError(err)
	updateReq.Header.Set("Content-Type", contentType)

	updateResp, err := s.httpClient.Do(updateReq)
	s.test.NoError(err)

	body, err := io.ReadAll(updateResp.Body)
	s.test.NoError(err)

	respCore := &response.Response{}
	s.test.NoError(json.Unmarshal(body, respCore))
	s.test.Equal(respCore.Status, response.StatusError)
	s.test.Equal(respCore.Error, "url with this alias not found")
}

func (s *UrlShortenerSuite) TestDeleteSuccess() {
	url := fmt.Sprintf("%s/url", s.server.URL)

	testURL := "https://mail.google.com/"
	testAlias := "mail"

	req := save.Request{
		URL:   testURL,
		Alias: testAlias,
	}

	marshalledSaveReq, err := json.Marshal(req)
	s.test.NoError(err)

	// Вставляем url и alias
	saveReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalledSaveReq))
	s.test.NoError(err)
	saveReq.Header.Set("Content-Type", contentType)

	saveResp, err := s.httpClient.Do(saveReq)
	s.test.NoError(err)
	s.test.Equal(http.StatusOK, saveResp.StatusCode)
	defer saveResp.Body.Close()

	// Проверяем, что url и alias вставились
	actualURL, err := s.storage.GetURL(testAlias)
	s.test.NoError(err)
	s.test.Equal(testURL, actualURL)

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

	// Проверяем, что alias удалился
	_, err = s.storage.GetURL(testAlias)
	s.test.ErrorIs(err, storage.ErrUrlNotFound)
}
