package podcast

import (
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v3"
	"io"
	"netradio/internal/model"
	"netradio/internal/repository"
	"netradio/pkg/context"
)

func PodacstGetter(context context.Context) (any, error) {

	audioTrack, _ := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: webrtc.MimeTypeOpus,
	}, "audio", "pion")

	// load music
	go func() {
		musicChunks, err := repository.NewMusicDB().LoadMusicBatch(model.MusicInfo{})
		if err != nil {
			return
		}

		for batch := range musicChunks {
			audioTrack.Write(batch) // пока попробуем писать без кодека, придумаем если что-то не пойдет
		}
	}()

	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// offer for webrtc
	offer, err := readOffer(context.GetRequest().Body)
	if err != nil {
		return nil, err
	}

	pc, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		return nil, err
	}

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
	}
	pc.AddTrack(audioTrack)
	return nil, nil
}

func readOffer(reader io.Reader) (*webrtc.SessionDescription, error) {
	dec := json.NewDecoder(reader)
	var offer webrtc.SessionDescription
	err := dec.Decode(&offer)
	if err != nil {
		return nil, err
	}

	if offer.Type != webrtc.SDPTypeOffer {
		return nil, fmt.Errorf("received SDP is not an offer")
	}

	return &offer, nil
}
