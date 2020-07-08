/* eslint-disable react/no-children-prop */
import React from 'react'

import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom'

import Room from '../pages/room'

import logo from './logo.svg'
import './App.css'

const App = () => {
  return (
    <Router>
      <Switch>
        <Route exact path="/">
          <h1>hi</h1>
        </Route>

        <Route path="/room/:roomName" children={<Room />} />
      </Switch>
    </Router>
  )
}

export default App
