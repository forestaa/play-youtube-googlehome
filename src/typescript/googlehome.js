import { spawn } from 'child_process';

var mdns = require('multicast-dns')();
var Client = require('castv2-client').Client;
var DefaultMediaReceiver = require('castv2-client').DefaultMediaReceiver;

var host = '';
mdns.on('response', function (response) {
    let name = '';
    for (const additional of response.additionals) {
        if (additional.type == 'TXT') name = additional.name;
        if (additional.type == 'A') host = additional.data;
    }
    if (name != '') console.log(`${name}(${host})`);
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


export async function playMusic(url) {
  console.log(`ipaddress = ${host}`); // TODO: check whether the address is valid

  const infos = await getInfos(url);
  const items = infos.map((info) => { 
    return {media: {
      contentId: info.url,
      streamType: 'BUFFERED',
      metadata: {
        metadataType: 3,
        title: info.title,
        images: [{url: info.thumbnail}],
      }
    }}; 
  });
  const client = new Client();
  client.connect(host, () => {
    console.log('connected! Launching DefaultMediaReceiver ...');
    client.launch(DefaultMediaReceiver, (err, player) => {
      if(err) throw err; 
      console.log('app "%s" launched, loading media ...', player.session.displayName);
      player.on('status', function(status) {
        console.log('status broadcast playerState=%s', status.playerState);
      });
      
      player.queueLoad(items, {}, (err, status) => {
        if (err) throw err;
        console.log(status)
      });
    });
  });
  client.on('error', function(err) {
    client.close();
    throw err;
  });
}

async function getInfos(url) {
  const youtubedl = spawn('youtube-dl', ['--get-title', '--get-url', '--get-thumbnail', '-x', url]);
  var n = 0;
  var info = {};
  var infos = [];
  for await (const data of youtubedl.stdout) {
    console.log(n + ': ' + data);
    const lines = data.toString().split('\n').filter((line) => line != "");
    for (const line of lines) {
      switch (n % 3) {
        case 0:
          info.title = line;
          break;
        case 1:
          info.url = line;
          break;
        default:
          info.thumbnail = line;
          infos.push(info);
          info = {};
          break;
      }
      n += 1;
    }
  }
  console.log(infos);
  return infos;
}
