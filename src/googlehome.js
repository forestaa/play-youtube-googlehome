const googlehome = require('google-home-notifier');
const language = 'jp';
var deviceAddress = '';
var mdns = require('multicast-dns')();
mdns.on('response', function (response) {
    let name = '';
    for (const additional of response.additionals) {
        if (additional.type == 'TXT') name = additional.name;
        if (additional.type == 'A') deviceAddress = additional.data;
    }
    if (name != '') console.log(`${name}(${deviceAddress})`);
    googlehome.ip(deviceAddress, language);
    mdns.destroy();
});

mdns.query({
    questions: [{
        name: '_googlecast._tcp.local',
        type: 'PTR'
    }]
});

setTimeout(function() {
  mdns.destroy();
}, 2500);


const ytdl = require('ytdl-core');
exports.playMusic = function(url) {
  console.log(`ipaddress = ${deviceAddress}`);
  ytdl.getInfo(url, (err, info) => {
    if (err) throw err;
    let audioFormat = ytdl.filterFormats(info.formats, 'audioonly');
    if (audioFormat.length > 0) {
      let url = audioFormat.pop().url;
      googlehome.play(url, res => {
        console.log(res);
      });
    }
  }); 
}
