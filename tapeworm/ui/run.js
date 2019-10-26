const express = require('express')
const moment = require('moment')
const app = express()
const port = 3000

app.get('/', (req, res) => {
    const d = moment().format()
    res.send('Hello World! ' + d)
})
app.listen(port, () => console.log(`Example app listening on port ${port}!`))
