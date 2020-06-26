import React, { useState, useEffect } from 'react'
import YoutubePlayer from 'youtube-player'

import { RootState } from '../../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'

import { setPlaying, setCurrent } from './playerSlice'
import { extractYoutubeID } from '../../util'

import { WebSocketContext } from '../../context/websocket'

interface PlayerInfo {
  id: string,
}

const Player = () => {
  const playerState = useSelector((state: RootState) => state.player)
  const ws = React.useContext(WebSocketContext)
  const [player, setPlayer] = useState(null as any)
  const [ready, setReady] = useState(false)
  const [playerInfo, setPlayerInfo] = useState({ id: 's36EMcPph00' } as PlayerInfo)
  const dispatch = useDispatch()

  useEffect(() => {
    if (!player) setPlayer(YoutubePlayer('player'))

    if (!player) return

    if (!ready) {
      player.on('ready', () => {
        let id = extractYoutubeID(playerState.current.url)
        if (!id) id = 's36EMcPph00'

        player.loadVideoById(id as string)
        player.playVideo()
        player.pauseVideo()
        setReady(true)
      })

      player.on('stateChange', (e: any) => {
        switch (e.data) {
          case 1:
            (ws.ws as WebSocket).send(JSON.stringify({
              action: 'set-video-playing',
              data: true
            }))
            return

          case 2:
            (ws.ws as WebSocket).send(JSON.stringify({
              action: 'set-video-playing',
              data: false
            }))
        }
      })

      return
    }

    let id = extractYoutubeID(playerState.current.url)
    if (!id) id = 's36EMcPph00'

    if (id !== playerInfo.id) {
      console.log('new video')
      player.loadVideoById(id as string)
      setPlayerInfo({ ...playerInfo, id: id })
    }

    playerState.isPlaying ? player.playVideo() : player.pauseVideo()
  })

  return (
    <>
      <h1>{playerState.current.title}</h1>
      <VideoInput url={playerState.current.url} ws={ws.ws}></VideoInput>
      <div id="player"></div>
      <TogglePlay isPlaying={playerState.isPlaying} ws={ws.ws}/>
    </>
  )
}

const TogglePlay = ({ isPlaying, ws }: { isPlaying: boolean, ws?: WebSocket }) => {
  const handleClick = () => {
    if (!ws) return
    (ws as WebSocket).send(JSON.stringify({
      action: 'set-video-playing',
      data: !isPlaying
    }))
  }

  return <button onClick={handleClick}>{isPlaying ? 'Pause' : 'Play'}</button>
}

const VideoInput = ({ url, ws }: { url: string, ws?: WebSocket }) => {
  const [val, setVal] = useState(url)

  const sendIfValid = (url: string) => {
    if (!extractYoutubeID(url)) return
    if (!ws) return

    console.log('Valid URL detected, sending to Websocket.');

    (ws as WebSocket).send(JSON.stringify({
      action: 'select-video',
      data: url
    }))
  }

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    e.preventDefault()
    setVal(e.currentTarget.value)
    sendIfValid(e.currentTarget.value)
  }

  return (
    <div>
      <input type="text" value={val} onChange={handleChange} placeholder="Enter a valid YouTube URL..." />
    </div>
  )
}

export default Player
