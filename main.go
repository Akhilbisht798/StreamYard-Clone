package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pion/webrtc/v3"
)

var peerConnection *webrtc.PeerConnection

func main() {

	var err error
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("client"))
	mux.Handle("/", fs)
	mux.HandleFunc("/connect", connectWebrtc)
	fmt.Print("Hello world")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		fmt.Printf("Internal Error in Server\n")
	}
}

func connectWebrtc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got Request in Connect")
	var offer webrtc.SessionDescription

	err := json.NewDecoder(r.Body).Decode(&offer)
	if err != nil {
		http.Error(w, "Invalid SDP offer", http.StatusBadRequest)
	}
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&webrtc.MediaEngine{}))

	var errPeerConnection error
	peerConnection, errPeerConnection = api.NewPeerConnection(webrtc.Configuration{})
	if errPeerConnection != nil {
		fmt.Printf("Error Creating a PeerConnection: %s \n", err)
		return
	}

	setRemoteErr := peerConnection.SetRemoteDescription(offer)
	if setRemoteErr != nil {
		http.Error(w, "Error in Setting Remote Description", http.StatusInternalServerError)
	}

	answer, errAnswer := peerConnection.CreateAnswer(nil)
	if errAnswer != nil {
		http.Error(w, "Failed to Create Answer", http.StatusInternalServerError)
	}

	errSetLocal := peerConnection.SetLocalDescription(answer)
	if errSetLocal != nil {
		http.Error(w, "Failed to set LocalDescription", http.StatusInternalServerError)
	}

	//answerJSON, errEncode := json.Marshal(answer)
	errEncode := json.NewEncoder(w).Encode(answer)
	if errEncode != nil {
		http.Error(w, "Failed to encode SDP answer", http.StatusInternalServerError)
		return
	}
	fmt.Println(answer)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(answer.SDP))
	if err != nil {
		http.Error(w, "Failed to send SDP answer", http.StatusInternalServerError)
	}
}
