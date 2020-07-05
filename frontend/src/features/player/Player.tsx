import React, { useState, useEffect } from 'react'
import YoutubePlayer from 'youtube-player'

import { RootState } from '../../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'

import { setPlaying, setCurrent, seekTo } from './playerSlice'
import { User } from '../../api/youtube-rooms-API'
import { extractYoutubeID } from '../../util'

import { WebSocketContext } from '../../websocket/context'
import {
  Action,
  sendClientAction,
  SetVideoPayloadClient,
  SetVideoPlayingPayloadClient,
  SeekToPayloadClient
} from '../../websocket/websocket'

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
    fontSize: '28px',
    paddingBottom: '1rem'
  } as React.CSSProperties
}

const Player = () => {
  const users = useSelector((state: RootState) => state.usersDisplay)
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
            } as SetVideoPlayingPayloadClient, ws.ws as WebSocket)
            return

          case 2:
            sendClientAction({
              action: Action.SetVideoPlaying,
              data: false
            } as SetVideoPlayingPayloadClient, ws.ws as WebSocket)
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
      <div className="columns">
        <div className="column">
          <div className="columns is-multiline">
            <div className="column is-6 is-12-mobile">
              <VideoInput url={playerState.current.url} ws={ws.ws}></VideoInput>
            </div>

          </div>

          <div className="columns is-multiline" style={{ minHeight: '60vh' }}>
            <div className="column is-8">
              <div style={{ width: '100%', height: '100%' }} id="player"></div>
            </div>

            <div className="column is-4">
              <PlayerSidebar users={users.users}/>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

const PlayerSidebar = ({ users }: { users: User[]}) => {
  const u = users.map(user => <li key={`viewers-${user.discordID}`}>{`${user.discordUsername}#${user.discordDiscriminator}`}</li>)
  return (
    <div className="columns is-multiline">
      <div className="column is-full">
        <h1 style={{ letterSpacing: '0.305em', fontSize: '16px', fontWeight: 600, color: '#9C9C9C' }}>VIEWERS</h1>
        <ul className="viewer-list">
          {u}
        </ul>
      </div>
    </div>
  )
}

const PlayerControls = ({ player, isPlaying, ws }: { player: any, isPlaying: boolean, ws?: WebSocket }) => {
  return (
    <div className="columns is-multiline" style={{ paddingLeft: '1rem' }}>
      <div className="column is-12">
        <TogglePlay isPlaying={isPlaying} ws={ws}/>
      </div>

      <div className="column is-12">
        <SeekControls player={player} ws={ws}/>
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
    } as SetVideoPlayingPayloadClient, ws)
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
    } as SeekToPayloadClient, ws)
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
    } as SetVideoPayloadClient, ws)
  }

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    e.preventDefault()
    setVal(e.currentTarget.value)
    sendIfValid(e.currentTarget.value)
  }

  return (
    <div>
      <input style={{ width: '100%' }} className="video-input" type="text" value={val} onChange={handleChange} placeholder="Enter a valid YouTube URL..." />
    </div>
  )
}

export default Player
