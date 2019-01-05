var googlehome = require('google-home-notifier');
var language = 'jp';

// googlehome.device('Living'); // Change to your Google Home name
// or if you know your Google Home IP
googlehome.ip('192.168.0.2', language);

// googlehome.notify('Hey Foo', function(res) {
//   console.log(res);
// });
googlehome.play('https://r1---sn-3qqp-ioqlz.googlevideo.com/videoplayback?c=WEB&keepalive=yes&sparams=clen%2Cdur%2Cei%2Cgcr%2Cgir%2Cid%2Cip%2Cipbits%2Citag%2Ckeepalive%2Clmt%2Cmime%2Cmm%2Cmn%2Cms%2Cmv%2Cpcm2cms%2Cpl%2Crequiressl%2Csource%2Cexpire&source=youtube&pl=47&txp=5511222&requiressl=yes&ipbits=0&fvip=1&mime=audio%2Fwebm&clen=81121123&gir=yes&mm=31%2C29&itag=251&key=yt6&signature=C6BA18374728AA6F71191D350042BF73C5ACF2ED.30CAFFE8F52D54475D60562BDCA95929315D2ABD&mn=sn-3qqp-ioqlz%2Csn-oguelne7&ms=au%2Crdu&mt=1546679270&mv=u&dur=5106.441&ei=8XgwXNngL8qy4wKqg7iAAg&gcr=jp&pcm2cms=yes&ip=240f%3A73%3A3a67%3A1%3Acaf7%3A33ff%3Afe80%3Aec10&lmt=1544144559286962&id=o-AILoPAsiHVAayyikMlZMl8CzJJhrVoU3bLEN7V9wNdV_&expire=1546702161&ratebypass=yes', res => {
  console.log(res);
});