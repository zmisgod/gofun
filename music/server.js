const express = require('express')
const app = express()
const port = 20050
const sign = require('./sign')

app.post('/getSign', (req, res) => {
    let result = sign.getSign(req.body)
    res.send({"code":200, "data":result})
})

app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`)
})