import React, { useState, useEffect } from 'react'
import YoutubePlayer from 'youtube-player'

import { RootState } from '../../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'

import { setPlaying, setCurrent, seekTo } from './playerSlice'
import { extractYoutubeID } from '../../util'

import { WebSocketContext } from '../../websocket/context'
import { Action, sendClientAction } from '../../websocket/websocket'

interface PlayerInfo {
  id: string,
}

interface VideoSeekInput {
  minutes: number,
  seconds: number
}

const styles = {
  title: {
    color: '#E3E3E3',
    fontWeight: 'normal',
    fontSize: '28px'
  } as React.CSSProperties
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
            sendClientAction({
              action: Action.SetVideoPlaying,
              data: true
            }, ws.ws as WebSocket)
            return

          case 2:
            sendClientAction({
              action: Action.SetVideoPlaying,
              data: false
            }, ws.ws as WebSocket)
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

    if (playerState.seekTo !== -1) {
      console.log('Player.tsx -> ' + playerState.seekTo)
      player.seekTo(playerState.seekTo, true)
      dispatch(seekTo(-1))
    }

    playerState.isPlaying ? player.playVideo() : player.pauseVideo()
  })

  return (
    <div className="section">
      <h1 style={styles.title}>{playerState.current.title}</h1>
      <VideoInput url={playerState.current.url} ws={ws.ws}></VideoInput>
      <div id="player"></div>

      <div className="controls">
        <TogglePlay isPlaying={playerState.isPlaying} ws={ws.ws}/>
        <SeekControls player={player} ws={ws.ws}/>
      </div>
    </div>
  )
}

const TogglePlay = ({ isPlaying, ws }: { isPlaying: boolean, ws?: WebSocket }) => {
  const handleClick = () => {
    if (!ws) return
    sendClientAction({
      action: Action.SetVideoPlaying,
      data: !isPlaying
    }, ws)
  }

  return <button onClick={handleClick}>{isPlaying ? 'Pause' : 'Play'}</button>
}

const SeekControls = ({ player, ws }: { player: any, ws?: WebSocket }) => {
  const [timeInput, setTimeInput] = useState('')
  const offset = 2

  const seekTo = (time: number) => {
    if (!ws) return
    sendClientAction({
      action: Action.SeekTo,
      data: Math.floor(time)
    }, ws)
  }

  const seekMinSec = (minutes: number, seconds: number) => seekTo((60 * minutes) + seconds)
  const seekOffset = async (offset: number) => { seekTo(await player.getCurrentTime() + offset) }

  const handleChange = (e: any) => {
    setTimeInput(e.target.value)
  }

  const handleKeyDown = (e: any) => {
    if (e.key !== 'Enter') return
    const [minutes, seconds] = timeInput.split(':')
    if (!minutes || !seconds) return
    if (isNaN(+minutes) || isNaN(+seconds)) return
    seekMinSec(+minutes, +seconds)
  }
  return (
    <>
      <h1><a onClick={() => seekOffset(offset * -1)}>👈</a> <a onClick={() => seekOffset(offset)}>👉</a></h1>
      <input type="text" onKeyDown={handleKeyDown} onChange={handleChange} value={timeInput} placeholder="Enter a time..." />
    </>
  )
}

const VideoInput = ({ url, ws }: { url: string, ws?: WebSocket }) => {
  const [val, setVal] = useState(url)

  const sendIfValid = (url: string) => {
    if (!extractYoutubeID(url)) return
    if (!ws) return

    console.log('Valid URL detected, sending to Websocket.')

    sendClientAction({
      action: Action.SetVideo,
      data: url
    }, ws)
  }

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    e.preventDefault()
    setVal(e.currentTarget.value)
    sendIfValid(e.currentTarget.value)
  }

  return (
    <div>
      <input className="video-input" type="text" value={val} onChange={handleChange} placeholder="Enter a valid YouTube URL..." />
    </div>
  )
}

export default Player
