const bodyParser = require('body-parser');
const express = require("express");
const app = express();
app.use(bodyParser.urlencoded({extended: true}));

const server = app.listen(3000, function(){
    console.log("Node.js is listening to PORT:" + server.address().port);
});
app.set('view engine', 'ejs');
app.get("/", function(req, res, next){
    res.render("index", {});
});

const googlehome = require('./googlehome');
app.post("/music", function(req, res) {
  console.log('post');
  try {
    googlehome.playMusic(req.body.url);
    res.render("index", {}); // TODO: return appropriate html
  } catch (err) {
    console.log('error');
    console.log(err);
    res.render("index", {});  // TODO: return appropriate html
  }
});
