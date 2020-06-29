import React from 'react'
import logo from './logo.svg'
import './App.css'

import WebSocketProvider, { WebSocketContext } from '../websocket/context'
import Player from '../features/player/Player'

const App = () => {
  return (
    <div className="App">
      <WebSocketProvider>
        <Player />
      </WebSocketProvider>
    </div>
  )
}

export default App
