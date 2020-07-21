import React, { createContext } from 'react'
import PropTypes from 'prop-types'
import { useDispatch } from 'react-redux'

import { Action, parsePayload } from './websocket'

export interface WSManager {
    ws?: WebSocket
}

export interface WSPayload {
  action: string,
  data: any
}

const WebSocketContext: React.Context<WSManager> = createContext({})

export { WebSocketContext }

const WebSocketProvider = ({ children, roomName, roomPassword }: { children: React.ReactChild | React.ReactChild[], roomName: string, roomPassword?: string}) => {
  let wsManager: WSManager

  const dispatch = useDispatch()

  if (window.WebSocket) {
    const queryParams = `roomName=${encodeURI(roomName)}` + (roomPassword ? `&roomPassword=${roomPassword}` : '')
    console.log(queryParams)
    const con = new WebSocket(`ws://api.theatre.theradisic.com/ws?${queryParams}`)
    wsManager = { ws: con }

    con.onopen = () => {
      console.log('WebSocket connection opened')
      const getUsers = JSON.stringify({ action: Action.GetUsers, data: null })
      const sample = JSON.stringify({ action: Action.SetVideo, data: 'https://www.youtube.com/watch?v=YT127qw8eQQ' })
      con.send(getUsers)
      setTimeout(() => con.send(sample), 3000)
    }

    con.onmessage = e => {
      const payload = JSON.parse(e.data)
      parsePayload(payload, dispatch)
    }

    return (
      <WebSocketContext.Provider value={wsManager}>
        {children}
      </WebSocketContext.Provider>
    )
  }

  return (
    <WebSocketContext.Provider value={{}}>

    </WebSocketContext.Provider>
  )
}

export default WebSocketProvider
