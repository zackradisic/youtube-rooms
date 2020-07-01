import { Dispatch } from '@reduxjs/toolkit'
import { setCurrent, setPlaying, seekTo } from '../features/player/playerSlice'

export enum Action {
    SetVideo = 'set-video',
    SetVideoPlaying = 'set-video-playing',
    SeekTo = 'seek-to'
}

interface WebsocketPayload<T extends Action, K> {
    action: T,
    data: K
}

interface WebsocketClientPayload<T extends Action, K> extends WebsocketPayload<T, K> {}

export type SetVideoPayload = WebsocketPayload<Action.SetVideo, { url: string, requester: string }>
export type SetVideoPlayingPayload = WebsocketPayload<Action.SetVideoPlaying, boolean>
export type SeekToPayload = WebsocketPayload<Action.SeekTo, number>

export type SetVideoPayloadClient = WebsocketClientPayload<Action.SetVideo, string>
export type SetVideoPlayingPayloadClient = WebsocketClientPayload<Action.SetVideoPlaying, boolean>
export type SeekToPayloadClient = WebsocketClientPayload<Action.SeekTo, number>

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

export const sendClientAction = <T extends Action, K> (payload: WebsocketClientPayload<T, K>, ws: WebSocket) => {
  ws.send(JSON.stringify(payload))
}
