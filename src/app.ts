import * as Express from "express";
import * as BodyParser from 'body-parser';

const app = Express();
app.use(BodyParser.urlencoded({extended: true}));

const server = app.listen(3000, () => console.log("Node.js is listening to PORT:" + server.address().port));
app.set('view engine', 'ejs');
app.get('/', (req, res, next) => res.render('index', {}));

const googlehome = require('./googlehome')
app.post('/music', (req, res) => {
  console.log('post')
  try {
    googlehome.playMusic(req.body.url);
    res.render('index', {}); //TODO: return appropriate html
  } catch (err) {
    console.log('error');
    console.log(err);
    res.render("index", {});  //TODO: return appropriate html
  }
});
