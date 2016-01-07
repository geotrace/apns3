// +build go1.6
package apns

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Адреса северов Apple Push Notification.
const (
	Development = "https://api.development.push.apple.com"
	Production  = "https://api.push.apple.com"
)

// Service описывает сервис для отправки уведомлений на устройства.
type Service struct {
	url    string       // адрес для запроса
	client *http.Client // клиент для доступа к сервису
}

// New инициализирует и возвращает сервис для отправки уведомлений.
func New(url string, cert tls.Certificate) *Service {
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	config.BuildNameToCertificate()
	return &Service{
		url: fmt.Sprintf("%v/3/device/", url), // добавляем адрес для отправки токена
		client: &http.Client{ // инициализируем HTTP Client с сертификатом
			Transport: &http.Transport{TLSClientConfig: config},
		},
	}
}

// Push отправляет на сервер уведомление payload для устройства с токеном token. Так же можно
// опционально указать дополнительные параметры для отправки сообщения в options.
// В ответ возвращается уникальный идентификатор сообщения или ошибка.
func (s *Service) Push(token []byte, payload interface{}, options *Options) (id string, err error) {
	var (
		data []byte // данные для отправки
		ok   bool   // флаг, что данные уже в нужном формате
	)
	if data, ok = payload.([]byte); !ok { // проверяем, что данные уже готовы для отправки
		data, err = json.Marshal(payload) // переводим содержимое для отправки в формат JSON
		if err != nil {
			return
		}
	}
	// формируем запрос с данными
	req, err := http.NewRequest("POST", s.url+string(token), bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json") // устанавливаем тип данных
	// устанавливаем дополнительные заголовки запроса, если они определены
	if options != nil {
		if options.ID != "" {
			req.Header.Set("apns-id", options.ID)
		}
		if !options.Expire.IsZero() {
			req.Header.Set("apns-expiration", strconv.FormatInt(options.Expire.Unix(), 10))
		}
		if options.LowPriority { // если не установлено, то приоритет считается 10
			req.Header.Set("apns-priority", "5")
		}
		if options.Topic != "" {
			req.Header.Set("apns-topic", options.Topic)
		}
	}
	resp, err := s.client.Do(req) // выполняем отправку на сервер
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		id = resp.Header.Get("apns-id")
		return // все хорошо — возвращаем идентификатор
	}
	var response ErrResponse                     // описание ошибки
	json.NewDecoder(resp.Body).Decode(&response) // декодируем описание ошибки и возвращаем его
	response.Code = resp.StatusCode              // добавляем код статуса ответа
	err = response
	return
}

// Options добавляет к отправляемому сообщению дополнительные необязательные данные, которые
// уточняют поведение сервиса.
type Options struct {
	ID          string    // назначить уникальный идентификатор
	Expire      time.Time // доставить до или не доставлять уже
	LowPriority bool      // низкий приоритет доставки
	Topic       string    // тема для сертификата с поддержкой нескольких тем
}
