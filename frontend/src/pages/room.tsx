import React, { useEffect } from 'react'

import { RootState } from '../app/rootReducer'
import { useSelector } from 'react-redux'
import { useHistory, useParams } from 'react-router-dom'

import WebSocketProvider, { WebSocketContext } from '../websocket/context'
import Player from '../features/player/Player'
import { RoomPreview } from '../api/youtube-rooms-API'

const Room = () => {
  const history = useHistory()
  const { roomName } = useParams()
  const { name, password } = useSelector((state: RootState) => state.roomCredentials)
  const room = useSelector((state: RootState) => state.roomPreviews.rooms).find((r: RoomPreview) => r.name === name)

  useEffect(() => {
    // In the future this should make an API call to retrieve room info if
    // it is not found in the state
    if (!name || !room) {
      history.push('/')
    }
  })

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
      <Player />
    </WebSocketProvider>
  )
}

export default Room
