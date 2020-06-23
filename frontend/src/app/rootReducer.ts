import { combineReducers } from '@reduxjs/toolkit'

import videosDisplayReducer from '../features/videosDisplay'
import usersDisplayReducer from '../features/usersDisplay'
import playerReducer from '../features/player/playerSlice'

const rootReducer = combineReducers({
  videosDisplay: videosDisplayReducer,
  usersDisplay: usersDisplayReducer,
  player: playerReducer
})

export type RootState = ReturnType<typeof rootReducer>

export default rootReducer
