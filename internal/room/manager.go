package room

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/zackradisic/youtube-rooms/internal/models"
)

// Manager manages rooms and users
type Manager struct {
	rooms     map[string]*Room
	DB        *gorm.DB
	mux       sync.Mutex
	saveVideo chan *SaveVideoRequest
}

// NewManager returns a new manager
func NewManager(db *gorm.DB) *Manager {
	return &Manager{
		rooms:     make(map[string]*Room),
		DB:        db,
		saveVideo: make(chan *SaveVideoRequest),
	}
}

// GetRooms returns an array of all Rooms currently being managed
func (m *Manager) GetRooms() []*Room {
	m.mux.Lock()
	defer m.mux.Unlock()

	rooms := []*Room{}
	for _, v := range m.rooms {
		rooms = append(rooms, v)
	}

	return rooms
}

// GetRoom returns a room specified by the given name and caches it
func (m *Manager) GetRoom(name string) (*Room, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	// First check the internal cache to see if the room exists
	if room, ok := m.rooms[name]; ok {
		return room, nil
	}

	// If not found lets check the DB
	room := &models.Room{
		Name: name,
	}
	if err := m.DB.First(room, "name = ?", room.Name).Error; err != nil {
		return nil, fmt.Errorf("could not find that room (%s)", name)
	}

	// Cache it
	newRoom := NewRoom(room)
	newRoom.saveVideo = m.saveVideo

	lastVideo := &models.RoomVideo{}
	if err := m.DB.Where("room_id = ?", room.ID).Last(lastVideo).Error; err != nil {
		newRoom.LastVideo = lastVideo
	} else {
		return nil, err
	}

	m.rooms[room.Name] = newRoom

	return newRoom, nil
}

func (m *Manager) getLatestVideo(roomID uint) (*models.Video, error) {
	type result struct {
		Title     string
		YoutubeID string
	}

	r := &result{}
	stmt := "SELECT v.title, v.youtube_id FROM room_videos rv, videos v WHERE rv.video_id = v.id AND rv.room_id = $1 ORDER BY rv.created_at DESC LIMIT 1"
	err := m.DB.Raw(stmt, roomID).Scan(r).Error
	if err != nil {
		return nil, err
	}

	return &models.Video{Title: r.Title, YoutubeID: r.YoutubeID}, nil
}

// RemoveUser removes a user from the room it is in, if the user is not in a room it does
// nothing
func (m *Manager) RemoveUser(user *User) {
	// m.mux.Lock()
	// defer m.mux.Unlock()
	for _, room := range m.rooms {
		if room.Model.Name == user.CurrentRoom.Model.Name {
			room.RemoveUser(user)
		}
	}
}

// ListenForVideoSave listens for calls to save a video to DB
func (m *Manager) ListenForVideoSave() {
	go func() {
		for {
			select {
			case video := <-m.saveVideo:
				m.SaveVideo(video)
			}
		}
	}()
}

// SaveVideo saves a video to DB
func (m *Manager) SaveVideo(request *SaveVideoRequest) error {
	video := &models.Video{}
	tx := m.DB.Begin()
	if err := tx.Where(&models.Video{YoutubeID: request.Video.ExtractID()}).First(video).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			video.Title = request.Video.Title
			video.YoutubeID = request.Video.ExtractID()

			if err = tx.Create(video).Error; err != nil {
				tx.Rollback()
				return err
			}

		} else {
			tx.Rollback()
			return err
		}
	}

	roomVideo := &models.RoomVideo{}
	roomVideo.VideoID = video.ID
	roomVideo.RoomID = request.Room.Model.ID
	roomVideo.RequesterID = request.Video.Requester.UserID
	if err := tx.Create(roomVideo).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
