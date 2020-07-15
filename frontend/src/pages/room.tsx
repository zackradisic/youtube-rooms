import React from 'react'

import { RootState } from '../app/rootReducer'
import { useSelector } from 'react-redux'

import WebSocketProvider, { WebSocketContext } from '../websocket/context'
import Player from '../features/player/Player'
import { useParams } from 'react-router-dom'

const Room = () => {
  const { roomName } = useParams()
  const { password } = useSelector((state: RootState) => state.roomCredentials)
  return (
    <WebSocketProvider roomName={roomName} roomPassword={password}>
      <Player />
    </WebSocketProvider>
  )
}

export default Room
