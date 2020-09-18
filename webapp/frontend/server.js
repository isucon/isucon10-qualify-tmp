const express = require('express')
const next = require('next')
const { createProxyMiddleware } = require('http-proxy-middleware')

const port = parseInt(process.env.PORT, 10) || 3000
const dev = process.env.NODE_ENV !== 'production'
const app = next({ dev })
const handle = app.getRequestHandler()
const apiHost = process.env.API_HOST || 'http://localhost'
const apiPort = process.env.API_PORT || '1323'

app.prepare().then(() => {
  const server = express()

  server.use(
    '/api',
    createProxyMiddleware({ target: `${apiHost}:${apiPort}` })
  )

  server.all('*', (req, res) => handle(req, res))

  server.listen(port, err => {
    if (err) throw err
    console.log(`> Ready on http://localhost:${port}`)
  })
})
