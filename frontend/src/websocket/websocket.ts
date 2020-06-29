import { Dispatch } from '@reduxjs/toolkit'
import { setCurrent, setPlaying, seekTo } from '../features/player/playerSlice'

enum Action {
    SetVideo = 'set-video',
    SetVideoPlaying = 'set-video-playing',
    SeekTo = 'seek-to'
}

interface Payload {
    action: Action,
    data: any
}

interface SetVideoPayload extends Payload {
    action: Action.SetVideo,
    data: { url: string, requester: string }
}

interface SetVideoPlayingPayload extends Payload {
    action: Action.SetVideoPlaying,
    data: boolean
}

interface SeekToPayload extends Payload {
    action: Action.SeekTo,
    data: number
}

export const parsePayload = (json: any, dispatch: Dispatch<any>) => {
  if ((json as SetVideoPayload).action === Action.SetVideo) {
    const payload = json as SetVideoPayload
    const video = { url: payload.data.url, requester: payload.data.requester, title: 'test' }
    dispatch(setCurrent(video))
  } else if ((json as SetVideoPlayingPayload).action === Action.SetVideoPlaying) {
    const payload = json as SetVideoPlayingPayload
    dispatch(setPlaying(payload.data))
  } else if ((json as SeekToPayload).action === Action.SeekTo) {
    const payload = json as SeekToPayload
    dispatch(seekTo(payload.data as number))
  }
}
