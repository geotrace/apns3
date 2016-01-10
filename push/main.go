// Send Apple Push notification
//
//  ./push [-params] <token> [<token2> [...]]
//    -a    use production service
//    -b badge
//          badge number
//    -c certificate
//          push certificate (default "cert.p12")
//    -f file
//          JSON file with push message
//    -p password
//          certificate password
//    -t text
//          message text (default "Hello!")
//
//  Sample JSON file:
//    {
//      "payload": {
//        "aps": {
//          "alert": "message",
//          "badge": 0
//        }
//      }
//    }
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/geotrace/apns3"
	"github.com/mdigger/p12"
)

func main() {
	certFileName := flag.String("c", "cert.p12", "push `certificate`")
	password := flag.String("p", "", "certificate `password`")
	production := flag.Bool("a", false, "use production service")
	notificationFileName := flag.String("f", "", "JSON `file` with push message")
	alert := flag.String("t", "Hello!", "message `text`")
	badge := flag.Uint("b", 0, "`badge` number")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Send Apple Push notification\n")
		fmt.Fprintf(os.Stderr, "%s [-params] <token> [<token2> [...]]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\n"+`Sample JSON file:
  { 
    "payload": {
      "aps": {
        "alert": "message",
        "badge": 0 
      }
    }
  }`)
	}
	flag.Parse()
	log.SetFlags(0)

	if flag.NArg() < 1 {
		log.Fatalln("Error: no tokens")
	}
	tokens := make([][]byte, flag.NArg())
	for i, token := range flag.Args() {
		btoken, err := hex.DecodeString(token)
		if err != nil || len(btoken) != 32 {
			log.Fatalln("Bad token:", token)
		}
		tokens[i] = btoken
	}
	var payload = make(map[string]interface{})
	if *notificationFileName != "" {
		data, err := ioutil.ReadFile(*notificationFileName)
		if err != nil {
			log.Fatalln("Error loading push file:", err)
		}
		err = json.Unmarshal(data, &payload)
		if err != nil {
			log.Fatalln("Error parsing push file:", err)
		}
	} else if *alert != "" {
		payload["aps"] = map[string]interface{}{
			"alert": *alert,
			"badge": *badge,
		}
	} else {
		log.Fatalln("Nothing to send")
	}
	cert, err := p12.Load(*certFileName, *password)
	if err != nil {
		log.Fatalln("Error loading certificate:", err)
	}
	url := apns.Development
	if *production {
		url = apns.Production
	}
	client := apns.New(url, cert)
	for _, token := range tokens {
		id, err := client.Push(token, payload, nil)
		if err != nil {
			log.Println("Error:", err)
			break
		}
		log.Println("Sended:", id)
	}
	log.Println("Complete!")
}
