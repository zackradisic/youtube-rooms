import React from 'react'

import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom'

import WebSocketProvider, { WebSocketContext } from '../websocket/context'
import Player from '../features/player/Player'

import logo from './logo.svg'
import './App.css'

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
