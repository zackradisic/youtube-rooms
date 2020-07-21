import React, { useEffect } from 'react'

import { RootState } from '../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'
import { useHistory, useParams } from 'react-router-dom'

import WebSocketProvider, { WebSocketContext } from '../websocket/context'
import Player from '../features/player/Player'
import { RoomPreview } from '../api/youtube-rooms-API'
import { fetchRoomPreview } from '../features/roomPreviewsDisplay'

const Room = () => {
  const dispatch = useDispatch()
  const history = useHistory()
  const { roomName } = useParams()
  const { name, password } = useSelector((state: RootState) => state.roomCredentials)
  const rooms = useSelector((state: RootState) => state.roomPreviews.rooms)
  const room = rooms.find((r: RoomPreview) => r.name === roomName)

  useEffect(() => {
    const a = async () => {
      if (!name || !room) {
        await dispatch(fetchRoomPreview(roomName))
      }
    }

    a()
  }, [])

  if (!password && room?.passwordProtected) {
    return (
      <div className="section">
        <div className="container">
          <h1 style={{
            fontSize: '36px',
            color: '#E3E3E3'
          }}>This room is password protected.</h1>
        </div>
      </div>
    )
  }

  return (
    <WebSocketProvider roomName={roomName} roomPassword={password}>
      <Player name={name} />
    </WebSocketProvider>
  )
}

export default Room
