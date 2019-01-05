const googlehome = require('google-home-notifier');
const language = 'jp';

// googlehome.device('Living'); // Change to your Google Home name
// or if you know your Google Home IP
googlehome.ip('192.168.0.2', language);

const ytdl = require('ytdl-core');
const url = 'https://www.youtube.com/watch?v=94zKvCMy0-M';
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