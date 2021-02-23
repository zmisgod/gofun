const express = require('express')
const app = express()
const port = 20050
const sign = require('./sign')
const bodyParser = require("body-parser");

app.use(bodyParser.urlencoded({extended: false}));

app.post('/getSign', function(req, res) {
    try{
        let jsonObj = JSON.parse(req.body["needSign"])
        if (jsonObj !== undefined) {
            let result = sign.getSign(jsonObj)
            res.send({"code":200, "data":result, "msg": "ok"})
        }else{
            res.send({"code":400, "data":"", "msg": "parse error"})
        }
    }catch (e) {
        res.send({"code":500, "data":"", "msg": e.toString()})
    }
})

app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`)
})