package apns_test

import (
	"encoding/hex"
	"log"
	"math/rand"
	"time"

	"github.com/mdigger/apns3"
)

func Example() {
	cert, err := apns.LoadCertificate("cert.p12", "password")
	if err != nil {
		log.Fatalln("Error loading certificate:", err)
	}
	client := apns.New(apns.Development, cert)
	token, _ := hex.DecodeString(`883982D57CDC4138D71E16B5ACBCB5DEBE3E625AFCEEE809A0F32895D2EA9D51`)
	payload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": "Hello!",
			"badge": rand.Int31n(99),
		},
		"time": time.Now().Format(time.RFC3339Nano),
	}
	id, err := client.Push(token, payload, nil)
	if err != nil {
		log.Fatalln("Error push:", err)
	}
	log.Println("Sended:", id)
}
