package Webhook_Gdrive

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"Webhook-Gdrive/routes"
	"Webhook-Gdrive/utils"
)

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/", routes.Receive).Methods("POST")
	log.Println("Starting Webhook listener")
	port := utils.GetPort()
	log.Printf("Started Weblistener on  port %s\n", port)
	_ = http.ListenAndServe("0.0.0.0:"+port, router)

}