package routes

import (
	"Webhook-Gdrive/types"
	"log"
	"net/http"
)

func Receive(w http.ResponseWriter, r *http.Request){
	log.Printf("New Webhook recieved")
	log.Printf(r.RemoteAddr)
	var reqData types.Request
	data := reqData.Unmarshal(r)
	log.Println(data)
	var call types.ErrorResponse
	if reqData.Channel != "" {
		call.Status = 200
		call.Send(w)
	}else {
		call.Status = 404
		call.Send(w)
	}
	//Google Drive stuff here ig?
}