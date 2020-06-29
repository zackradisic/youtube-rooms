import React, { createContext } from 'react'
import PropTypes from 'prop-types'
import { useDispatch } from 'react-redux'

import { parsePayload } from './websocket'

export interface WSManager {
    ws?: WebSocket
}

export interface WSPayload {
  action: string,
  data: any
}

const WebSocketContext: React.Context<WSManager> = createContext({})

export { WebSocketContext }

const WebSocketProvider = (props: any) => {
  let wsManager: WSManager

  const dispatch = useDispatch()

  if (window.WebSocket) {
    const con = new WebSocket('ws://localhost/ws?roomName=zack%27s%20room&roomPassword=test123')
    wsManager = { ws: con }

    con.onopen = () => {
      console.log('WebSocket connection opened')
      const sample = JSON.stringify({ action: 'select-video', data: 'https://www.youtube.com/watch?v=dkrKp4nEe4w&t=22s' })
      setTimeout(() => con.send(sample), 10000)
    }

    con.onmessage = e => {
      const payload = JSON.parse(e.data)
      parsePayload(payload, dispatch)
    }

    return (
      <WebSocketContext.Provider value={wsManager}>
        {props.children}
      </WebSocketContext.Provider>
    )
  }

  return (
    <WebSocketContext.Provider value={{}}>

    </WebSocketContext.Provider>
  )
}

export default WebSocketProvider
