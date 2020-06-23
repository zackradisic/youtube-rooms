import React from 'react'
import YoutubePlayer from 'youtube-player'

import { RootState } from '../../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'
import { useEffect } from 'react'

import { setPlaying, setCurrent } from './playerSlice'

const Player = () => {
  const playerState = useSelector((state: RootState) => state.player)

  let player
  useEffect(() => {
    player = YoutubePlayer('player')
    player.loadVideoById('f5KsFMsd1qE')
    player.playVideo()
    console.log(playerState.current)
  })

  return (
    <div id="player"></div>
  )
}

export default Player
