package ws

import (
	"encoding/json"
	"fmt"

	"github.com/zackradisic/youtube-rooms/internal/room"
)

// ActionInvoker manages and invokes actions
type ActionInvoker struct {
	actions map[string]action
}

type action func(data interface{}, client *Client) (*HubMessage, error)

// InvokeAction invokes an action specified by name
func (a *ActionInvoker) InvokeAction(message *ClientMessage, outbound chan *HubMessage) error {
	if ac, ok := a.actions[message.Action]; ok {
		res, err := ac(message.Data, message.Client)
		if err != nil {
			fmt.Println(err)
			return err
		}
		go func() { outbound <- res }()
		return nil
	}
	return fmt.Errorf("could not find action: (%s)", message.Action)
}

func (a *ActionInvoker) registerAction(name string, ac action) {
	a.actions[name] = ac
}

// NewActionInvoker returns a new ActionInvoker
func NewActionInvoker() *ActionInvoker {
	a := &ActionInvoker{
		actions: make(map[string]action),
	}

	a.registerAction("set-video", selectVideo)
	a.registerAction("set-video-playing", setVideoPlaying)
	a.registerAction("seek-to", seekTo)
	a.registerAction("get-users", getUsers)
	return a
}

func getUsers(data interface{}, client *Client) (*HubMessage, error) {
	type roomUser struct {
		DiscordID           string `json:"discordID"`
		DiscordUsername     string `json:"discordUsername"`
		DiscordDisriminator string `json:"discordDiscriminator"`
	}

	type jsonResponse struct {
		Action string      `json:"action"`
		Users  []*roomUser `json:"data"`
	}

	jr := &jsonResponse{
		Action: "get-users",
		Users:  make([]*roomUser, 0),
	}

	users := client.user.CurrentRoom.GetUsers()
	for _, user := range users {
		ru := &roomUser{
			DiscordID:           user.Model.DiscordID,
			DiscordUsername:     user.Model.LastDiscordUsername,
			DiscordDisriminator: user.Model.LastDiscordUsername,
		}

		jr.Users = append(jr.Users, ru)
	}

	r, err := json.Marshal(jr)
	if err != nil {
		return nil, err
	}

	return NewHubMessage(r, client.user.CurrentRoom), nil
}

func setVideoPlaying(data interface{}, client *Client) (*HubMessage, error) {
	type jsonResponse struct {
		Action    string `json:"action"`
		IsPlaying bool   `json:"data"`
	}

	isPlaying, ok := data.(bool)
	if !ok {
		return nil, fmt.Errorf("Invalid data supplied")
	}

	jr := &jsonResponse{
		Action:    "set-video-playing",
		IsPlaying: isPlaying,
	}

	r, err := json.Marshal(jr)
	if err != nil {
		return nil, err
	}

	client.user.CurrentRoom.IsPlaying = true
	return NewHubMessage(r, client.user.CurrentRoom), nil
}

func selectVideo(data interface{}, client *Client) (*HubMessage, error) {
	type jsonVideo struct {
		URL       string `json:"url"`
		Requester string `json:"requester"`
	}

	type jsonData struct {
		Action string     `json:"action"`
		Data   *jsonVideo `json:"data"`
	}
	url, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("Invalid data supplied")
	}

	jv := &jsonVideo{
		URL:       url,
		Requester: client.user.Model.LastDiscordUsername,
	}
	jd := &jsonData{
		Action: "set-video",
		Data:   jv,
	}
	j, err := json.Marshal(jd)
	if err != nil {
		return nil, err
	}

	video := room.NewVideo(url, client.user)
	client.user.CurrentRoom.SetCurrentVideo(video)
	return NewHubMessage(j, client.user.CurrentRoom), nil
}

func seekTo(data interface{}, client *Client) (*HubMessage, error) {
	type jsonResponse struct {
		Action string `json:"action"`
		Data   int    `json:"data"`
	}

	secondsFloat, ok := data.(float64)
	if !ok {
		return nil, fmt.Errorf("Invalid data: Expected float64")
	}

	seconds := int(secondsFloat)

	jr := &jsonResponse{
		Action: "seek-to",
		Data:   seconds,
	}

	j, err := json.Marshal(jr)

	if err != nil {
		return nil, err
	}

	return NewHubMessage(j, client.user.CurrentRoom), nil
}

// Maybe there is something we can do to make this code more DRY?
//
// 		1) Have the action functions return the JSON data as []byte
//		2) A surrounding function creates the HubMessage
//
// This makes the action function solely responsible for validating the input
// and marshalling the data into JSON form.
//
// I don't like how the JSON response structs are defined each time in the body of
// the action functions, though I am uncertain of the best way to resolve this problem.
//
