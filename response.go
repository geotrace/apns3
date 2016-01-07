package apns

import (
	"fmt"
	"net/http"
	"time"
)

// ErrResponse описывает ответ от сервера
type ErrResponse struct { // для разбора ответа с ошибкой
	Code      int    // статус кода ответа сервера
	Reason    string `json:"reason"`    // ошибка
	Timestamp int64  `json:"timestamp"` // временная метка
}

// Error возвращает описание ошибки.
func (r ErrResponse) Error() string {
	return fmt.Sprintf("[%d] %s", r.Code, r.Reason)
}

// Time возвращает временную метку приведенной к формату времени.
func (r ErrResponse) Time() time.Time {
	return time.Unix(r.Timestamp, 0)
}

// Status возвращает код HTTP-ответа сервера.
func (r ErrResponse) Status() int {
	return r.Code
}

// Status возвращает текстовое описание кода HTTP-ответа сервера.
func (r ErrResponse) StatusText() string {
	return http.StatusText(r.Code)
}
