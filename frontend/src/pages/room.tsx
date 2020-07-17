import React from 'react'

import { RootState } from '../app/rootReducer'
import { useSelector } from 'react-redux'

import WebSocketProvider, { WebSocketContext } from '../websocket/context'
import Player from '../features/player/Player'
import { useParams } from 'react-router-dom'

const Room = () => {
  const { roomName } = useParams()
  const { password } = useSelector((state: RootState) => state.roomCredentials)

  if (!password) {
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
