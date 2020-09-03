package Webhook_Gdrive

import (
	"Webhook-Gdrive/helpers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/", helpers.Receive).Methods("POST")
	log.Println("Starting Webhook listener")
	port := helpers.GetPort()
	log.Printf("Started Weblistener on  port %s\n", port)
	_ = http.ListenAndServe("0.0.0.0:"+port, router)

}