package apns

import (
	"fmt"
	"time"
)

// Error описывает ответ от сервера
type Error struct { // для разбора ответа с ошибкой
	Code      int    // статус кода ответа сервера
	Reason    string `json:"reason"`    // ошибка
	Timestamp int64  `json:"timestamp"` // временная метка
}

// Error возвращает описание ошибки.
func (r Error) Error() string {
	return fmt.Sprintf("[%d] %s", r.Code, r.Reason)
}

// Time возвращает временную метку приведенной к формату времени.
func (r Error) Time() time.Time {
	return time.Unix(r.Timestamp, 0)
}

// Status возвращает код HTTP-ответа сервера.
func (r Error) Status() int {
	return r.Code
}
