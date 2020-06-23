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
		fmt.Println(res)
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

	a.registerAction("select-video", selectVideo)
	return a
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
