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

	a.registerAction("set-video", selectVideoAction)
	a.registerAction("set-video-playing", setVideoPlayingAction)
	a.registerAction("seek-to", seekToAction)
	a.registerAction("get-users", getUsersAction)
	a.registerAction("init-client", initClientAction)
	return a
}

type roomUserJSON struct {
	DiscordID           string `json:"discordID"`
	DiscordUsername     string `json:"discordUsername"`
	DiscordDisriminator string `json:"discordDiscriminator"`
}

func initClientAction(data interface{}, client *Client) (*HubMessage, error) {
	type currentVideoJSON struct {
		Title     string `json:"title"`
		URL       string `json:"url"`
		Requester string `json:"requester"`
	}

	type jsonResponse struct {
		Action  string           `json:"action"`
		Current currentVideoJSON `json:"data"`
	}

	jr := &jsonResponse{
		Action:  "init-client",
		Current: currentVideoJSON{},
	}

	title := ""
	url := ""
	requester := ""

	if client.user.CurrentRoom.Current != nil {
		title = client.user.CurrentRoom.Current.Title
		url = client.user.CurrentRoom.Current.URL
		requester = client.user.CurrentRoom.Current.Requester.DiscordHandle()
	}

	jr.Current.Title = title
	jr.Current.URL = url
	jr.Current.Requester = requester

	r, err := json.Marshal(jr)
	if err != nil {
		return nil, err
	}

	return NewHubMessage(r, nil, []*room.User{client.user}), nil
}

func getUsersJSON(room *room.Room) []*roomUserJSON {
	response := make([]*roomUserJSON, 0)

	users := room.GetUsers()
	for _, user := range users {
		ru := &roomUserJSON{
			DiscordID:           user.Model.DiscordID,
			DiscordUsername:     user.Model.LastDiscordUsername,
			DiscordDisriminator: user.Model.LastDiscordDiscriminator,
		}

		response = append(response, ru)
	}

	return response
}

func getUsersAction(data interface{}, client *Client) (*HubMessage, error) {
	type jsonResponse struct {
		Action string          `json:"action"`
		Users  []*roomUserJSON `json:"data"`
	}

	jr := &jsonResponse{
		Action: "get-users",
		Users:  make([]*roomUserJSON, 0),
	}

	jr.Users = getUsersJSON(client.user.CurrentRoom)

	r, err := json.Marshal(jr)
	if err != nil {
		return nil, err
	}

	return NewHubMessage(r, client.user.CurrentRoom, nil), nil
}

func setVideoPlayingAction(data interface{}, client *Client) (*HubMessage, error) {
	type jsonResponse struct {
		Action    string `json:"action"`
		IsPlaying bool   `json:"data"`
	}

	isPlaying, ok := data.(bool)
	if !ok {
		return nil, fmt.Errorf("Invalid data supplied")
	}

	fmt.Printf("Is playing? (%t)\n", client.user.CurrentRoom.GetIsPlaying())
	if client.user.CurrentRoom.GetIsPlaying() == isPlaying {
		return nil, fmt.Errorf("Room playing state is already: (%t)", isPlaying)
	}

	jr := &jsonResponse{
		Action:    "set-video-playing",
		IsPlaying: isPlaying,
	}

	r, err := json.Marshal(jr)
	if err != nil {
		return nil, err
	}

	client.user.CurrentRoom.SetIsPlaying(isPlaying)
	return NewHubMessage(r, client.user.CurrentRoom, nil), nil
}

func selectVideoAction(data interface{}, client *Client) (*HubMessage, error) {
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
	return NewHubMessage(j, client.user.CurrentRoom, nil), nil
}

func seekToAction(data interface{}, client *Client) (*HubMessage, error) {
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

	return NewHubMessage(j, client.user.CurrentRoom, nil), nil
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
