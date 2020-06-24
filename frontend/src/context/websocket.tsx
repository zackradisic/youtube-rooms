import React, { createContext } from 'react'
import PropTypes from 'prop-types'
import { useDispatch } from 'react-redux'

import { setPlaying, setCurrent } from '../features/player/playerSlice'

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
    const con = new WebSocket('ws://localhost:3000/ws?roomName=zack%27s%20room&roomPassword=test123')
    wsManager = { ws: con }

    con.onopen = () => {
      console.log('WebSocket connection opened')
<<<<<<< HEAD
      const sample = JSON.stringify({ action: 'select-video', data: 'https://www.youtube.com/watch?v=dkrKp4nEe4w&t=22s' })
=======
      const sample = JSON.stringify({ action: 'select-video', data: 'https://www.youtube.com/watch?v=rWBSMsLG8po' })
>>>>>>> 91fdae3d3a99636e03aff907780afff0e78330cc
      setTimeout(() => con.send(sample), 3000)
    }

    con.onmessage = e => {
      const payload = JSON.parse(e.data)
      switch (payload.action) {
        case 'set-video': {
          const vid = {
            url: payload.data.url,
            title: 'test',
            requester: payload.data.requester
          }
          dispatch(setCurrent(vid))
          break
        }
      }
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
