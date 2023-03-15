package podcast

import (
	"encoding/json"
	"fmt"
	"io"
	"netradio/internal/controller/requests"
	"netradio/pkg/context"

	"github.com/pion/webrtc/v3"
)

func PodacstGetter(context context.Context) (any, error) {
	/*
		offer, err := readOffer(context.GetRequest().Body)
		if err != nil {
			return nil, err
		}

		pc, err := webrtc.NewPeerConnection(webrtchelper.GetPeerConfig())
		if err != nil {
			return nil, err
		}

		rtpSender, err := pc.AddTrack(webrtchelper.GetAudioTrack())
		if err != nil {
			panic(err)
		}

		go func() {
			rtcpBuf := make([]byte, 1500)
			for {
				if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
					return
				}
			}
		}()

		err = pc.SetRemoteDescription(*offer)
		if err != nil {
			return nil, err
		}

		answer, err := pc.CreateAnswer(nil)
		if err != nil {
			return nil, err
		}

		err = pc.SetLocalDescription(answer)
		if err != nil {
			return nil, err
		}

		ansJson, err := json.Marshal(answer)
		if err != nil {
			return nil, err
		}
		_, err = context.GetResponseWriter().Write(ansJson)
		if err != nil {
			return nil, err
		}*/
	return nil, nil
}

func readOffer(reader io.Reader) (*webrtc.SessionDescription, error) {
	var request requests.PodcastStartRequest
	dec := json.NewDecoder(reader)
	err := dec.Decode(&request)
	if err != nil {
		return nil, err
	}

	var offer webrtc.SessionDescription
	offer.SDP = request.SDP
	offer.Type = webrtc.SDPTypeOffer
	if offer.Type != webrtc.SDPTypeOffer {
		return nil, fmt.Errorf("received SDP is not an offer")
	}

	return &offer, nil
}
