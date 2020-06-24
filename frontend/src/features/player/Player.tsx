import React, { useState } from 'react'
import YoutubePlayer from 'youtube-player'

import { RootState } from '../../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'
import { useEffect } from 'react'

import { setPlaying, setCurrent } from './playerSlice'
import { extractYoutubeID } from '../../util'

const Player = () => {
  const playerState = useSelector((state: RootState) => state.player)
  const [ready, setReady] = useState(false)
  const [player, setPlayer] = useState(null as any)

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

      return
    }

    let id = extractYoutubeID(playerState.current.url)
    if (!id) id = 's36EMcPph00';

    player.loadVideoById(id as string)
  })

  return (
    <div id="player"></div>
  )
}

export default Player
