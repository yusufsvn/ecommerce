package Utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ReadJSON: Request body'yi okur ve gönderilen struct işaretçisine (pointer) yazar.
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// 1. Güvenlik: Maksimum boyut limiti (Örnek: 1MB)
	// Bu, kötü niyetli büyük payload'ların sunucuyu kilitlemesini engeller.
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	// 2. Strict Mode: Struct'ta olmayan bir alan JSON'da varsa hata fırlat.
	// İstemcinin yanlış veri göndermesini engellemek için önemlidir.
	dec.DisallowUnknownFields()

	// 3. Decode İşlemi
	err := dec.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// JSON formatı bozuksa
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// Yanlış tip gönderildiyse (örn: int beklenen yere string gelmesi)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// Body boşsa
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Bilinmeyen bir alan varsa (DisallowUnknownFields nedeniyle)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Dosya boyutu limitini aştıysa
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// Kod hatası (Pointer gönderilmediyse)
		case errors.As(err, &invalidUnmarshalError):
			return fmt.Errorf("error unmarshalling json: %s", err.Error())

		default:
			return err
		}
	}

	// 4. Tek bir JSON objesi olduğundan emin ol (örn: {"a":1}{"a":1} gönderilmemeli)
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// Opsiyonel: Hata mesajlarını standart bir formatta dönmek için yardımcı fonksiyon
func WriteError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{
		Error:   true,
		Message: err.Error(),
	})
}
