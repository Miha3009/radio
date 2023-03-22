package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"netradio/internal/controller/requests"
	"netradio/internal/model"
	"netradio/internal/service"
	"netradio/pkg/context"
	"netradio/pkg/files"
	"netradio/pkg/handlers"
	"netradio/pkg/log"
	webrtchelper "netradio/pkg/webrtc"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pion/webrtc/v3"
	"golang.org/x/net/websocket"
)

func HandleGetChannels(ctx context.Context) (any, error) {
	res, err := service.NewChannelService().GetChannels()
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

func HandleGetChannel(ctx context.Context) (any, error) {
	var request requests.GetChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	res, err := service.NewChannelService().GetChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	if !res.Found {
		ctx.GetResponseWriter().WriteHeader(http.StatusNotFound)
		return res, nil
	}

	return res, nil
}

func HandleCreateChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.CreateChannelRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = service.NewChannelService().CreateChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleDeleteChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.DeleteChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewChannelService().DeleteChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleUpdateChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.UpdateChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	res, err := service.NewChannelService().UpdateChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	if !res.Found {
		ctx.GetResponseWriter().WriteHeader(http.StatusNotFound)
		return res, nil
	}

	return nil, nil
}

func HandleStartChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.StartChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewChannelService().StartChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleStopChannel(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.StopChannelRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	err := service.NewChannelService().StopChannel(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleUploadLogo(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	path, err := files.Save(ctx.GetRequest())
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	var request requests.UploadLogoRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")
	request.Logo = path

	err = service.NewChannelService().UploadLogo(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleAddTrack(ctx context.Context) (any, error) {
	user := ctx.GetUser()
	if user.Role != model.UserAdministrator {
		ctx.GetResponseWriter().WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var request requests.AddTrackRequest
	decoder := json.NewDecoder(ctx.GetRequest().Body)
	err := decoder.Decode(&request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	err = service.NewChannelService().AddTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return nil, nil
}

func HandleGetCurrentTrack(ctx context.Context) (any, error) {
	var request requests.GetCurrentTrackRequest
	request.ID = chi.URLParam(ctx.GetRequest(), "id")

	user := ctx.GetUser()
	if user.Role != model.UserGuest {
		userId := strconv.Itoa(user.ID)
		request.UserID = &userId
	}

	res, err := service.NewChannelService().GetCurrentTrack(request)
	if err != nil {
		ctx.GetResponseWriter().WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return res, nil
}

type ConnectStruct struct {
	ID string
}

func (s *ConnectStruct) HandleConnectChannel(ws *websocket.Conn) {
	logger := log.NewLogger()
	peerConnection, err := webrtc.NewPeerConnection(webrtchelper.GetPeerConfig())
	if err != nil {
		logger.Error(err)
		return
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		outbound, marshalErr := json.Marshal(c.ToJSON())
		if marshalErr != nil {
			logger.Error(marshalErr)
			return
		}

		if _, err = ws.Write(outbound); err != nil {
			logger.Error(err)
			return
		}
	})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnOpen(func() {
			for range time.Tick(time.Second * 3) {
				if err = d.SendText(time.Now().String()); err != nil {
					logger.Error(err)
					return
				}
			}
		})
	})

	track, err := webrtchelper.GetAudioTrack(s.ID)
	if err != nil {
		logger.Error(err)
		return
	}

	rtpSender, err := peerConnection.AddTrack(track)
	if err != nil {
		logger.Error(err)
		return
	}

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	for {
		buf := make([]byte, 10000)

		n, err := ws.Read(buf)
		if err != nil {
			logger.Error(err)
			return
		}

		var (
			candidate webrtc.ICECandidateInit
			offer     webrtc.SessionDescription
		)

		switch {
		case json.Unmarshal(buf[:n], &offer) == nil && offer.SDP != "":
			if err = peerConnection.SetRemoteDescription(offer); err != nil {
				logger.Error(err)
				return
			}

			answer, answerErr := peerConnection.CreateAnswer(nil)
			if answerErr != nil {
				logger.Error(answerErr)
				return
			}

			if err = peerConnection.SetLocalDescription(answer); err != nil {
				logger.Error(err)
				return
			}

			outbound, marshalErr := json.Marshal(answer)
			if marshalErr != nil {
				logger.Error(marshalErr)
				return
			}

			if _, err = ws.Write(outbound); err != nil {
				logger.Error(err)
				return
			}
		case json.Unmarshal(buf[:n], &candidate) == nil && candidate.Candidate != "":
			if err = peerConnection.AddICECandidate(candidate); err != nil {
				logger.Error(err)
				return
			}
		default:
			fmt.Println("Unknown message")
			return
		}
	}
}

func RouteChannelPaths(
	core handlers.Core,
	router chi.Router,
) {
	router.MethodFunc("GET", "/channel", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetChannels), core))
	router.MethodFunc("GET", "/channel/{id}", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetChannel), core))
	router.MethodFunc("PUT", "/channel", handlers.MakeHandler(HandleCreateChannel, core))
	router.MethodFunc("DELETE", "/channel/{id}", handlers.MakeHandler(HandleDeleteChannel, core))
	router.MethodFunc("PATCH", "/channel/{id}", handlers.MakeHandler(HandleUpdateChannel, core))
	router.MethodFunc("POST", "/channel/{id}/start", handlers.MakeHandler(HandleStartChannel, core))
	router.MethodFunc("POST", "/channel/{id}/stop", handlers.MakeHandler(HandleStopChannel, core))
	router.MethodFunc("POST", "/channel/{id}/upload-logo", handlers.MakeHandler(HandleUploadLogo, core))
	router.MethodFunc("POST", "/channel/{id}/add-track", handlers.MakeHandler(HandleAddTrack, core))
	router.MethodFunc("GET", "/channel/{id}/track", handlers.MakeHandler(handlers.MakeJSONWrapper(HandleGetCurrentTrack), core))

	router.HandleFunc("/channel/{id}/connect",
		func(w http.ResponseWriter, req *http.Request) {
			conn := ConnectStruct{chi.URLParam(req, "id")}
			s := websocket.Server{Handler: websocket.Handler(conn.HandleConnectChannel)}
			s.ServeHTTP(w, req)
		})
}
