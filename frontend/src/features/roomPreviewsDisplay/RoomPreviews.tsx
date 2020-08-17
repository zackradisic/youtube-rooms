import axios from 'axios'
import React, { useState, useEffect } from 'react'

import { RootState } from '../../app/rootReducer'
import { useSelector, useDispatch } from 'react-redux'
import { useHistory } from 'react-router-dom'

import { RoomPreview } from '../../api/youtube-rooms-API'
import { setRoomPreviews, fetchRoomPreview } from '../roomPreviewsDisplay/index'
import { setCredentials } from '../roomCredentials/index'

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { IconDefinition } from '@fortawesome/fontawesome-common-types'
import { faLock } from '@fortawesome/free-solid-svg-icons'
import { Dispatch } from '@reduxjs/toolkit'
import { verifyRoomPassword } from '../../util'

const style = {
  roomPreview: {
    backgroundColor: '#272727',
    color: '#AFAFAF',
    borderRadius: '5px',
    padding: '0.5rem'
  },
  title: {
    fontSize: '36px',
    color: '#E3E3E3'
  }
}

const RoomPreviews = () => {
  const history = useHistory()
  const dispatch = useDispatch()

  const rooms = useSelector((state: RootState) => state.roomPreviews.rooms)

  let fetched = false
  const [modalActive, setModalActive] = useState(false)
  const [selectedRoom, setSelectedRoom] = useState<RoomPreview | null>(null)

  useEffect(() => {
    const loadRoomPreviews = async () => {
      if (!fetched) {
        await dispatch(fetchRoomPreview())
        fetched = true
      }
    }

    loadRoomPreviews()
  }, [])

  const onClick = (room: RoomPreview) => {
    console.log(room)
    if (room.passwordProtected) {
      setSelectedRoom(room)
      setModalActive(true)
      return
    }

    dispatch(setCredentials({ name: room.name }))
    history.push('/rooms/' + encodeURI(room.name))
  }

  const toggleModal = () => setModalActive(!modalActive)

  const enterPassword = async (e: any) => {
    if (selectedRoom) {
      const password = (document.getElementById('password-input') as HTMLInputElement).value
      try {
        const valid = await verifyRoomPassword(password, selectedRoom.name)
        if (!valid) {
          document.getElementById('password-context')!.innerText = 'That password is invalid.'
          return
        }
        dispatch(setCredentials({
          name: selectedRoom.name,
          password: password
        }))
      } catch (err) {
        return
      }

      history.push('/rooms/' + encodeURI(selectedRoom.name))
    }
  }

  const roomsDOM = rooms.map(e => <RoomPreviewContainer onClick={() => onClick(e)} key={`room-preview-${e.id}`} room={e} />)
  return (
    <div className="section">
      <div className="container">
        <div className="columns is-multiline">
          <div className="column is-6">
            <h1 style={style.title}>Rooms</h1>
          </div>
          <div className="column is-6 has-text-right" style={{ verticalAlign: 'text-bottom' }}>
            <h1 style={{ color: 'rgb(22, 22, 22)' }}>h</h1>
            <h1 style={{ ...style.title, fontSize: '16px' }}>Logged in as: zack#6593</h1>
          </div>
          {roomsDOM}
          <PasswordPrompt enterPassword={enterPassword} dispatch={dispatch} isActive={modalActive} toggleModal={toggleModal} />
        </div>
      </div>
    </div>
  )
}

const RoomPreviewContainer = ({ room, onClick }: { room: RoomPreview, onClick: () => void}) => {
  const icon = room.passwordProtected ? <FontAwesomeIcon icon={faLock} /> : <> </>
  return (
    <div className="column is-4" onClick={onClick}>
      <div style={style.roomPreview} className="room-preview">
        <div className="columns is-mobile">
          <div className="column is-1">
            {icon}
          </div>
          <h1 className="column">{room.name}</h1>
          <div className="column has-text-right">
            {`${room.usersCount}`}
          </div>
        </div>
      </div>
    </div>
  )
}

const PasswordPrompt = ({ dispatch, isActive, toggleModal, enterPassword }: {
  dispatch: Dispatch<any>, isActive: boolean, toggleModal: () => void,
  enterPassword: (e: any) => void }) => {
  console.log(isActive)
  return (
    <div className={`modal ${isActive ? 'is-active' : ''}`}>
      <div className="modal-background"></div>
      <div className="modal-card">
        <header className="modal-card-head">
          <p className="modal-card-title">This room is password protected</p>
          <button onClick={toggleModal} className="button is-light is-small">X</button>
        </header>
        <section className="modal-card-body">
          <input id="password-input" style={{ width: '100%' }} className="video-input" placeholder="Enter the password here..." type="password" />
        </section>
        <footer className="modal-card-foot">
          <button onClick={enterPassword} className="button is-success">Enter</button>
          <p id="password-context" className="has-text-danger"></p>
        </footer>
      </div>
    </div>
  )
}

export default RoomPreviews
