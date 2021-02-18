const express = require('express')
const app = express()
const port = 20050
const sign = require('./sign')
const bodyParser = require("body-parser");

app.use(bodyParser.urlencoded({extended: false}));

app.post('/getSign', function(req, res) {
    let result = sign.getSign(req.body["needSign"])
    console.log(result)
    res.send({"code":200, "data":result})
})

app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`)
})