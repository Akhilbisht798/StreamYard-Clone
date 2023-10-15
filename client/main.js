const btn = document.getElementById("connect");
btn.addEventListener("click", async function () {
  console.log("trying to start connection");
  const peerConnection = new RTCPeerConnection();

  const offer = await peerConnection.createOffer();

  await peerConnection.setLocalDescription(offer);

  fetch("/connect", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(offer),
  })
    .then((res) => {
      console.log(res);
      return JSON.parse(res);
    })
    .then((ans) => {
      peerConnection.setRemoteDescription(ans);
    })
    .catch((err) => {
      console.log("Error Sending in SDP", err);
    });
});
