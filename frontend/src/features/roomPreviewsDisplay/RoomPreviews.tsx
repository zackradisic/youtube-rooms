import axios from 'axios'
import React, { useEffect } from 'react'

import { RootState } from '../../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'

import { RoomPreview } from '../../api/youtube-rooms-API'
import { setRoomPreviews } from '../roomPreviewsDisplay/index'

const RoomPreviews = () => {
  const rooms = useSelector((state: RootState) => state.roomPreviews.rooms)
  const dispatch = useDispatch()

  useEffect(() => {
    const loadRoomPreviews = async () => {
      if (rooms.length === 0) {
        try {
          const res = await axios.get('http://localhost/api/rooms/')
          if (!res.data.rooms) return
          if (res.data.rooms.length === 0) return

          if ((res.data.rooms as RoomPreview[])[0].id) {
            dispatch(setRoomPreviews(res.data.rooms as RoomPreview[]))
          }
        } catch (err) {
          console.log(err)
        }
      }
    }

    loadRoomPreviews()
  }, [])

  const roomsDOM = rooms.map(e => <RoomPreviewContainer key={`room-preview-${e.id}`} room={e} />)
  return (
    <div>
      {roomsDOM}
    </div>
  )
}

const RoomPreviewContainer = ({ room }: { room: RoomPreview}) => {
  return (
    <div>
      <h1>{room.name}</h1>
      <p>{`Users: ${room.usersCount}`}</p>
    </div>
  )
}

export default RoomPreviews
