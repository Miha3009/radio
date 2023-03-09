const pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: 'stun:stun.l.google.com:19302'
    }
  ]
})

let localDescription = null

pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState)
pc.onicecandidate = event => {
  if (event.candidate === null) {
    localDescription = pc.localDescription
  }
}

pc.addTransceiver('audio')
pc.createOffer()
  .then(d => pc.setLocalDescription(d))
  .catch(console.log)

pc.ontrack = function (event) {
  const el = document.getElementById('audio1')
  el.srcObject = event.streams[0]
  el.autoplay = true
  el.controls = true
}

window.startSession = () => {
  let xhr = new XMLHttpRequest();
  xhr.open("POST", "http://localhost:8080/podcast/0/start");
  xhr.setRequestHeader('Content-Type', 'application/json; charset=utf-8');
  console.log(localDescription);
  xhr.send(JSON.stringify({sdp: localDescription.sdp}));
  xhr.onload = (e) => {
    response = JSON.parse(xhr.response)
    console.log(response)
    try {
      pc.setRemoteDescription(response)
    } catch (e) {
      alert(e)
    }
  }
}
