const googlehome = require('google-home-notifier');
const language = 'jp';
googlehome.ip('192.168.0.2', language);
// googlehome.device('Living'); // Change to your Google Home name

const ytdl = require('ytdl-core');

exports.playMusic = function(url) {
  ytdl.getInfo(url, (err, info) => {
    if (err) throw err;
    // if (err) return;
    let audioFormat = ytdl.filterFormats(info.formats, 'audioonly');
    if (audioFormat.length > 0) {
      // console.log(audioFormat);
      let hoge = audioFormat.pop().url;
      googlehome.play(hoge, res => {
        console.log(res);
      });
    }
  }); 
}
