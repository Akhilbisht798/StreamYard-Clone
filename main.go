package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pion/webrtc/v3"
)

var peerConnection *webrtc.PeerConnection
var dc *webrtc.DataChannel

type WebRtcData struct {
	Offer        webrtc.SessionDescription `json:"offer"`
	ICECandidate []webrtc.ICECandidateInit `json:"iceCandidates"`
}

func main() {
	// Webrtc implementation
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"}, // STUN server for NAT traversal
			},
		},
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&webrtc.MediaEngine{}))

	var errPeerConnection error
	peerConnection, errPeerConnection = api.NewPeerConnection(config)
	if errPeerConnection != nil {
		fmt.Println("Error in creating peerConnection")
		return
	}
	// Data Channel
	peerConnection.OnDataChannel(func(ch *webrtc.DataChannel) {
		dc = ch
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Println("Message From Client" + string(msg.Data))
		})
		dc.OnOpen(func() {
			fmt.Println("Data Channel from client is opened")
		})
	})
	peerConnection.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		fmt.Println("Recived Stream from web")
	})

	var err error
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("client"))
	mux.Handle("/", fs)
	mux.HandleFunc("/connect", connectWebrtc)
	fmt.Print("Server Running in Localhost:3000\n")
	err = http.ListenAndServe(":3000", mux)
	dc.OnOpen(func() {
		fmt.Println("Connection Opened")
	})
	if err != nil {
		fmt.Printf("Internal Error in Server\n")
	}
}

func connectWebrtc(w http.ResponseWriter, r *http.Request) {
	var offer webrtc.SessionDescription
	err := json.NewDecoder(r.Body).Decode(&offer)
	if err != nil {
		http.Error(w, "Invalid SDP offer", http.StatusBadRequest)
	}
	fmt.Println(offer)

	// Exchanging of sdp
	setRemoteErr := peerConnection.SetRemoteDescription(offer)
	if setRemoteErr != nil {
		fmt.Printf("Error in Setting Remote Description\n%v\n", setRemoteErr)
		http.Error(w, "Error in Setting Remote Description", http.StatusInternalServerError)
	}

	answer, errAnswer := peerConnection.CreateAnswer(nil)
	if errAnswer != nil {
		fmt.Printf("Failed to Create Answer\n%v \n", err)
		http.Error(w, "Failed to Create Answer", http.StatusInternalServerError)
	}

	errSetLocal := peerConnection.SetLocalDescription(answer)
	if errSetLocal != nil {
		fmt.Printf("Failed to set LocalDescription\n%v \n", err)
		http.Error(w, "Failed to set LocalDescription", http.StatusInternalServerError)
	}

	answerJSON, errEncode := json.Marshal(answer)
	if errEncode != nil {
		fmt.Printf("Failed to encode SDP answer\n%v \n", err)
		http.Error(w, "Failed to encode SDP answer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(answerJSON)
	if err != nil {
		fmt.Printf("Failed to send SDP answer\n%v \n", err)
		http.Error(w, "Failed to send SDP answer", http.StatusInternalServerError)
	}
}
