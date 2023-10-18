const btn = document.getElementById("connect");
const iceCandidate = []
  const config = [
    {
      urls: "stun:stun.l.google.com:19302",
    }
]

const peerConnection = new RTCPeerConnection({ iceServers: config});
peerConnection.onicecandidate = (event) => {
  if (event.candidate) {
     iceCandidate.push(event.candidate)
     console.log(iceCandidate)
   
  }
}
const dc = peerConnection.createDataChannel("data-channel", { reliable: true})
dc.onopen = e =>  {
    console.log("Data Channel Opened")
}
dc.onmessage = e => console.log("Got Message ", e.data)

btn.addEventListener("click", async function () {
  console.log("trying to start connection");
  try {
    const offer = await peerConnection.createOffer();
    await peerConnection.setLocalDescription(offer);
    sendOfferAndIce(offer, iceCandidate)
  } catch(error) {
    console.log("Error Creating a webrtc connection", error)
  }
});

const sendOfferAndIce = (offer, iceCandidates) => {
  fetch("/connect", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(offer),
  })
    .then((res) => {
      return res.json();
    })
    .then((ans) => {
      peerConnection.setRemoteDescription(ans);
    })
   .catch((err) => {
      console.log("Error Sending in SDP", err);
    });
}


