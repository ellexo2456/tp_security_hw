package main

import (
	"log"

	"github.com/ellexo2456/tp_security_hw/src/handler"
	"github.com/ellexo2456/tp_security_hw/src/proxy"
)

//
//import (
//	"github.com/ellexo2456/tp_security_hw/proxy"
//	"log"
//	"net/http"
//)
//
//const crt = "./list/ca.crt"
//const key = "./list/ca.key"
//
//func main() {
//	p := proxy.Server{}
//
//	proxy := &http.Server{
//		Addr:    ":8080",
//		Handler: p,
//		//TLSNextProto: make(map[string]func(*http.Server, *list.Conn, http.Handler)),
//	}
//
//	if err := proxy.ListenAndServe(); err != nil {
//		log.Fatalf(err.Error())
//	}
//	if err := proxy.ListenAndServeTLS(crt, key); err != nil {
//		log.Fatalf(err.Error())
//	}
//	//output, err := exec.Command("./gen_cert.sh", "mail.ru", strconv.Itoa(rand.Intn(1000000000000))).Output()
//	//if err != nil {
//	//	fmt.Println(err)
//	//}
//	//
//	//err = os.WriteFile("mail.ru", output, 0644)
//	//if err != nil {
//	//	fmt.Println(err)
//	//}
//
//}

const port = 8080

func main() {
	h, err := handler.New()
	if err != nil {
		log.Fatal(err)
	}
	proxy.StartServer(port, h.Handle)
}
