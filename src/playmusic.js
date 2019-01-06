const bodyParser = require('body-parser');
const express = require("express");
const app = express();
app.use(bodyParser.urlencoded({extended: true}));

const server = app.listen(3000, function(){
    console.log("Node.js is listening to PORT:" + server.address().port);
});

const googlehome = require('google-home-notifier');
const language = 'jp';
googlehome.ip('192.168.0.2', language);
// googlehome.device('Living'); // Change to your Google Home name

const ytdl = require('ytdl-core');

function playMusic(url) {
  ytdl.getInfo(url, (err, info) => {
    if (err) throw err;
    let audioFormat = ytdl.filterFormats(info.formats, 'audioonly');
    if (audioFormat.length > 0) {
      console.log(audioFormat);
      let hoge = audioFormat.pop().url;
      googlehome.play(hoge, res => {
        console.log(res);
      });
    }
  }); 
}

app.post("/", function(req, res) {
  console.log('test');
  playMusic(req.body.url);
});

app.set('view engine', 'ejs');
app.get("/", function(req, res, next){
    res.render("index", {});
});