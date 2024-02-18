const express = require('express')
const cors = require('cors');
const fs = require('fs');

const port = 3000

const app = express();
app.use(cors());

const WORDS_LENGTH = 161744;

let words = [];

function getRandomArbitrary(min, max) {
  return Math.floor(Math.random() * (max - min) + min);
}

app.get('/word', (req, res) => {

  const randomWord = words[getRandomArbitrary(0, WORDS_LENGTH)];
  console.log(randomWord);
  res.json({ randomWord: randomWord.toLowerCase() })

});

app.listen(port, () => {
  fs.readFile('../words.txt', function (err, data) {
    if (err) throw err;
    words = data.toString().split("\n");
  });
  console.log(`Example app listening on port ${port}`)
})