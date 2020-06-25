import React, { useState } from 'react'
import YoutubePlayer from 'youtube-player'

import { RootState } from '../../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'
import { useEffect } from 'react'

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
        if (!id) id = 's36EMcPph00';

        player.loadVideoById(id as string)
        player.playVideo()
        player.pauseVideo()
        setReady(true)
      })

      player.on('stateChange', (e: any) => {
        switch(e.data) {
            case 1:
              (ws.ws as WebSocket).send(JSON.stringify({
                action: 'set-video-playing',
                data: true
              }))
                return;

            case 2:
              (ws.ws as WebSocket).send(JSON.stringify({
                action: 'set-video-playing',
                data: false
              }))
              console.log('oooga')
              return;

        }
    })

      return
    }

    let id = extractYoutubeID(playerState.current.url)
    if (!id) id = 's36EMcPph00';

    if (id !== playerInfo.id) {
      console.log('new video')
      player.loadVideoById(id as string)
      setPlayerInfo({ ...playerInfo, id: id })
    }

    playerState.isPlaying ? player.playVideo() : player.pauseVideo()
  })

  const handleClick = () => {
    (ws.ws as WebSocket).send(JSON.stringify({
      action: 'set-video-playing',
      data: !playerState.isPlaying
    }))
  }
  return (
    <>
      <h1>{playerState.current.title}</h1>
      <div id="player"></div>
      <button onClick={handleClick}>toggle play</button>
    </>
  )
}

export default Player
