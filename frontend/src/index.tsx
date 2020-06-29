import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import './styles.scss'
import * as serviceWorker from './serviceWorker'
import store from './app/store'
import { Provider } from 'react-redux'

const render = () => {
  const App = require('./app/App').default
  ReactDOM.render(
    <Provider store={store}>
      <React.StrictMode>
        <App />
      </React.StrictMode>
    </Provider>,
    document.getElementById('root')
  )
}

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister()

render()
