import axios from 'axios'

const https = require('https')

console.log(process.env.NODE_ENV !== 'development')
const instance = axios.create()

export default instance
