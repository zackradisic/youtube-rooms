import React from 'react'

import { RootState } from '../app/rootReducer'
import { useSelector } from 'react-redux'

import WebSocketProvider, { WebSocketContext } from '../websocket/context'
import Player from '../features/player/Player'
import { useParams } from 'react-router-dom'

const Room = () => {
  const { name } = useParams()
  const { password } = useSelector((state: RootState) => state.roomCredentials)
  return (
    <WebSocketProvider roomName={name} roomPassword={password}>
      <Player />
    </WebSocketProvider>
  )
}

export default Room
