import { Dispatch } from '@reduxjs/toolkit'
import { setCurrent, setPlaying, seekTo } from '../features/player/playerSlice'

export enum Action {
    SetVideo = 'set-video',
    SetVideoPlaying = 'set-video-playing',
    SeekTo = 'seek-to'
}

interface WebsocketPayload {
    action: Action,
    data: any
}

interface SetVideoPayload extends WebsocketPayload {
    action: Action.SetVideo,
    data: { url: string, requester: string }
}

interface SetVideoPlayingPayload extends WebsocketPayload {
    action: Action.SetVideoPlaying,
    data: boolean
}

interface SeekToPayload extends WebsocketPayload {
    action: Action.SeekTo,
    data: number
}

interface ClientPayload {
  action: Action,
  data: any
}

interface ClientSetVideoPayload extends ClientPayload {
  action: Action.SetVideo,
  data: string
}

interface ClientSetVideoPlayingPayload extends ClientPayload {
  action: Action.SetVideoPlaying,
  data: boolean
}

interface ClientSeekToPayload extends ClientPayload {
  action: Action.SeekTo,
  data: number
}

export const parsePayload = (json: any, dispatch: Dispatch<any>) => {
  if ((json as SetVideoPayload).action === Action.SetVideo) {
    const payload = json as SetVideoPayload
    const video = { url: payload.data.url, requester: payload.data.requester, title: "Zack's room" }
    dispatch(setCurrent(video))
  } else if ((json as SetVideoPlayingPayload).action === Action.SetVideoPlaying) {
    const payload = json as SetVideoPlayingPayload
    dispatch(setPlaying(payload.data))
  } else if ((json as SeekToPayload).action === Action.SeekTo) {
    const payload = json as SeekToPayload
    dispatch(seekTo(payload.data as number))
  }
}

export const sendClientAction = (payload: ClientPayload, ws: WebSocket) => {
  ws.send(JSON.stringify(payload))
}
