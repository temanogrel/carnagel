package myfreecams

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"sync"

	"net/url"

	"git.misc.vee.bz/carnagel/coordinator/pkg"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
	"github.com/sirupsen/logrus"
)

type Scraper struct {
	app *coordinator.Application
	log logrus.FieldLogger

	performers map[int64]*coordinator.MyFreeCamsPerformer
	wsSend     chan []byte
	mtx        sync.RWMutex
	sessionId  int64
}

func NewScraper(app *coordinator.Application) *Scraper {
	return &Scraper{
		app: app,
		log: app.Logger.
			WithField("component", "scraper").
			WithField("site", "mfc"),

		mtx:        sync.RWMutex{},
		performers: make(map[int64]*coordinator.MyFreeCamsPerformer),
		wsSend:     make(chan []byte),
	}
}

func (mfc *Scraper) GetPerformers() ([]*coordinator.MyFreeCamsPerformer, int) {
	mfc.mtx.Lock()
	defer mfc.mtx.Unlock()

	performers := []*coordinator.MyFreeCamsPerformer{}

	for _, performer := range mfc.performers {
		if performer.IsComplete() {
			performers = append(performers, performer)
		}
	}

	return performers, len(mfc.performers)
}

func (mfc *Scraper) Scrape(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			mfc.log.Warn("Received close from parent context")
			return nil

		default:
			// Run this in a loop so it restart if the websocket gets disconnected
			uri, domain, err := mfc.selectServerUri()
			if err != nil {
				mfc.log.WithError(err).Errorf("Failed to select server to connect to")

				return errors.Wrap(err, "Failed to select a server to connect to")
			}

			mfc.log.WithFields(logrus.Fields{
				"uri":    uri,
				"domain": domain,
			}).Debug("Selected server to connect to")

			mfc.run(ctx, uri)
		}
	}
}

func (mfc *Scraper) GetSessionId() int64 {
	return mfc.sessionId
}

func (mfc *Scraper) run(ctx context.Context, uri string) {
	header := http.Header{}
	header.Set("Origin", wsOrigin)
	header.Set("User-Agent", wsUserAgent)

	dialer := websocket.Dialer{
		Proxy: mfc.app.ScraperProxyService.GetProxy,
	}

	client, resp, err := dialer.Dial(uri, header)
	if err != nil {
		logger := mfc.log.WithError(err)

		if resp != nil {
			if data, e := ioutil.ReadAll(resp.Body); e == nil {
				logger.WithField("message", string(data))
			}
		}

		logger.Error("Error occurred on connection")
		return
	}

	defer resp.Body.Close()
	defer client.Close()

	// Randomly generate a timeout between 30 and 90 minutes
	timeout := (rand.Intn(60) + 1) + 30

	// Timeout the connection automatically after an hour
	ctx, cancel := context.WithTimeout(ctx, time.Minute*time.Duration(timeout))

	mfc.log.WithField("timeout", timeout).Info("Created context with timeout")

	go func() {
		for {
			msgType, message, err := client.ReadMessage()
			if err != nil {
				mfc.log.WithError(err).Error("Read error occurred")
				cancel()
				return
			}

			if msgType != websocket.TextMessage {
				mfc.log.WithField("msgType", msgType).Warn("Received message with unsupported type")
				return
			}

			go mfc.handleCommand(message)
		}
	}()

	timer := time.NewTimer(time.Second * 5)

	// Meet and greet with the server
	client.WriteMessage(websocket.TextMessage, []byte(greetMessage))
	client.WriteMessage(websocket.TextMessage, []byte(loginMessage))

	for {
		select {
		case <-ctx.Done():
			mfc.log.Info("Received ctx.Done() closing websocket")
			return

		case <-timer.C:
			mfc.log.Debug("Sending ping message")

			// send a custom ping message
			client.WriteMessage(websocket.TextMessage, []byte(pingMessage))
			timer.Reset(time.Second * 5)

		case msg := <-mfc.wsSend:
			mfc.log.WithField("message", strings.Trim(string(msg), "\n")).Debug("Sending message")

			client.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (mfc *Scraper) handleCommand(msg []byte) {

	args := strings.SplitN(string(msg), " ", 6)
	// First 4 digits in sequence represents the number of characters to read from msg to get command + args
	msgLength, err := strconv.ParseInt(args[0][:4], 10, 64)
	if err != nil {
		mfc.log.WithError(err).Errorf("Unable to calculate message length")
		return
	}

	// Check if the message is MORE than one command (plus the 4 digits for size of command)
	if int64(len(msg)) > msgLength+4 {
		// Break down into one command and call handleCommand again
		mfc.handleCommand(msg[:msgLength+4])
		mfc.handleCommand(msg[msgLength+4:])
		return
	}

	cmd, err := strconv.ParseInt(args[0][4:], 10, 64)
	if err != nil {
		mfc.log.WithError(err).Errorf("Failed to parse command")
		return
	}

	fromId, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		mfc.log.WithError(err).Errorf("Failed to parse fromId")
		return
	}

	toId, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		mfc.log.WithError(err).Errorf("Failed to parse toId")
		return
	}

	arg1, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		mfc.log.WithError(err).Errorf("Failed to parse arg1")
		return
	}

	arg2, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		mfc.log.WithError(err).Errorf("Failed to parse arg2")
		return
	}

	payload := ""

	if len(args) == 6 {
		if str, err := url.QueryUnescape(args[5]); err != nil {
			mfc.log.WithError(err).Errorf("Failed to parse json payload")
			return
		} else {
			payload = str
		}
	}

	command := &wsCommand{
		Command: IncomingCommand(cmd),
		FromId:  fromId,
		ToId:    toId,
		Arg1:    arg1,
		Arg2:    arg2,
		Payload: payload,
	}

	switch command.Command {
	case CommandRxData:

		// we must respond to this request else we don't get data
		if command.Arg1 == 1 {
			mfc.wsSend <- []byte(fmt.Sprintf("1 0 0 20071025 0 %s@1/guest:guest\n\x00", command.Payload))
		}

	case FCTYPE_ROOMDATA:
		mfc.updateRoomData(command)

	case FCTYPE_SESSIONSTATE:
		mfc.updateSession(command)

	case FCTYPE_LOGIN:
		mfc.onLogin(command)
	}
}

func (mfc *Scraper) onLogin(command *wsCommand) {

	// Subscribe to roomdata events
	mfc.wsSend <- []byte(fmt.Sprintf("%d %d 0 1 0\n\x00", FCTYPE_ROOMDATA, command.ToId))

	// Store our session id
	mfc.sessionId = command.ToId
}

func (mfc *Scraper) updateRoomData(command *wsCommand) {

	viewerCount := map[string]int64{}

	if err := json.Unmarshal([]byte(command.Payload), &viewerCount); err != nil {
		mfc.log.
			WithError(err).
			WithField("command", command).
			Errorf("Failed to unmarshal json in update room event")

		return
	}

	for uidAsString, count := range viewerCount {

		uid, err := strconv.ParseInt(uidAsString, 10, 64)
		if err != nil {
			mfc.log.WithError(err).Errorf("Failed parse uid as int")
			return
		}

		mfc.mtx.RLock()
		if performer, ok := mfc.performers[uid]; ok {
			performer.ViewerCount = count
		}
		mfc.mtx.RUnlock()
	}
}

func (mfc *Scraper) updateSession(command *wsCommand) {
	mfc.mtx.Lock()
	defer mfc.mtx.Unlock()

	if performer, ok := mfc.performers[command.Arg2]; ok {

		// Update performer with payload
		if err := json.Unmarshal([]byte(command.Payload), performer); err != nil {
			mfc.log.
				WithError(err).
				WithField("command", command).
				Errorf("Failed to unmarshal json in session update event")

			return
		}
	} else {

		performer := &coordinator.MyFreeCamsPerformer{}

		if err := json.Unmarshal([]byte(command.Payload), performer); err != nil {
			mfc.log.
				WithError(err).
				WithField("command", command).
				Errorf("Failed to unmarshal json in session update event")

			return
		}

		mfc.performers[command.Arg2] = performer
	}
}

func (mfc *Scraper) getServerConfig() (*ServerConfig, error) {

	request, err := http.NewRequest(http.MethodGet, "http://www.myfreecams.com/mfc2/data/serverconfig.js", nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create request object")
	}

	proxy, err := mfc.app.ScraperProxyService.GetProxy(request)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get proxy")
	}

	// Create a http proxy client
	client := pester.NewExtendedClient(&http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}})

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query the remote api")
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.Errorf("Unexpected response code, received %d", response.StatusCode)
	}

	serverConfig := &ServerConfig{}

	if err := json.NewDecoder(response.Body).Decode(serverConfig); err != nil {
		return nil, errors.Wrap(err, "Failed to decode response of getServerConfig")
	}

	return serverConfig, nil
}

func (mfc *Scraper) selectServerUri() (string, string, error) {

	serverConfig, err := mfc.getServerConfig()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get the server config")
	}

	servers := []string{}

	for hostname, protocol := range serverConfig.WebsocketServers {
		if protocol == "rfc6455" {
			servers = append(servers, hostname)
		}
	}

	if len(servers) == 0 {
		return "", "", errors.New("No websocket servers available")
	}

	selectedServer := servers[rand.Int()%len(servers)]

	// Server uri, selected hostname
	return fmt.Sprintf("ws://%s.myfreecams.com:8080/fcsl", selectedServer), selectedServer, nil
}
