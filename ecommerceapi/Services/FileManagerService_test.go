package services

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockMinIO(handler http.HandlerFunc) (*FileManagerService, *httptest.Server, error) {
	server := httptest.NewServer(handler)
	u, err := url.Parse(server.URL)
	if err != nil {
		server.Close()
		return nil, nil, err
	}

	client, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4("testaccesskey", "testsecretkey", ""),
		Secure: false,
		Region: "us-east-1",
	})
	if err != nil {
		server.Close()
		return nil, nil, err
	}

	svc := NewFileManagerService(client, "test-bucket")
	return svc, server, nil
}

func TestFileManagerService_UploadFile(t *testing.T) {
	t.Run("başarıyla dosya yüklemeli", func(t *testing.T) {
		putCalled := false
		svc, server, err := setupMockMinIO(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				putCalled = true
				w.Header().Set("ETag", "\"d41d8cd98f00b204e9800998ecf8427e\"")
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		require.NoError(t, err)
		defer server.Close()

		content := []byte("fake image content")
		reader := bytes.NewReader(content)

		uploadedName, err := svc.UploadFile("test.webp", int64(len(content)), "prefix", context.Background(), reader)
		assert.NoError(t, err)
		assert.Equal(t, "test.webp", uploadedName)
		assert.True(t, putCalled)
	})

	t.Run("storage sunucusu hata verirse UploadFile hata dönmeli", func(t *testing.T) {
		svc, server, err := setupMockMinIO(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		require.NoError(t, err)
		defer server.Close()

		content := []byte("fake image content")
		reader := bytes.NewReader(content)

		_, err = svc.UploadFile("error.webp", int64(len(content)), "prefix", context.Background(), reader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage'a yüklerken hatayla karşılaşıldı")
	})
}

func TestFileManagerService_DeleteFile(t *testing.T) {
	t.Run("başarıyla dosya silmeli", func(t *testing.T) {
		deleteCalled := false
		svc, server, err := setupMockMinIO(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				deleteCalled = true
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		require.NoError(t, err)
		defer server.Close()

		err = svc.DeleteFile(context.Background(), "image-123.webp")
		assert.NoError(t, err)
		assert.True(t, deleteCalled)
	})
}

func TestFileManagerService_DeleteMultipleFile(t *testing.T) {
	t.Run("başarıyla birden fazla dosya silmeli", func(t *testing.T) {
		svc, server, err := setupMockMinIO(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && strings.Contains(r.URL.RawQuery, "delete") {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><DeleteResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Deleted><Key>img1.webp</Key></Deleted><Deleted><Key>img2.webp</Key></Deleted></DeleteResult>`))
				return
			}
			w.WriteHeader(http.StatusOK)
		})
		require.NoError(t, err)
		defer server.Close()

		err = svc.DeleteMultipleFile(context.Background(), []string{"img1.webp", "img2.webp"})
		assert.NoError(t, err)
	})
}

func TestFileManagerService_GetPresignedFileObject(t *testing.T) {
	t.Run("başarıyla presigned URL üretmeli", func(t *testing.T) {
		svc, server, err := setupMockMinIO(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		require.NoError(t, err)
		defer server.Close()

		result, err := svc.GetPresignedFileObject()
		assert.NoError(t, err)
		require.NotNil(t, result.PresignedURL)
		assert.Contains(t, result.PresignedURL.String(), "a.jpeg")
		assert.NotEmpty(t, result.PresignedURL.Query().Get("X-Amz-Signature"))
	})
}
