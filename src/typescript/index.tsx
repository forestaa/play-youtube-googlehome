import * as React from 'react';
import * as ReactDOM from 'react-dom';
 

interface MusicPlayerState {
  websocket: WebSocket,
  connected: boolean,
  inputURL: string,
  playlist: Array<string>
}

class MusicPlayer extends React.Component<{}, MusicPlayerState> {
  constructor(props: {}) {
    super(props);

    const ws = new WebSocket('ws://localhost:4000');
    ws.onopen = () => {this.onopen()};
    ws.onclose = (e) => {this.onclose(e)};
    ws.onerror = () => {this.onerror()};
    ws.onmessage = (e) => {this.onmessage(e)};

    this.state = {
      websocket: ws,
      connected: false,
      inputURL: '',
      playlist: []
    };
  }

  onopen() {
    console.log('open websocket')
    this.setState({connected: true});
  }

  onclose(e: CloseEvent) {
    console.log('close websocket')
    this.setState({connected: false});
  }

  onerror() {
    console.error('got websocket error');
  }

  onmessage(e: MessageEvent) {
    try {
      const message = JSON.parse(e.data);
      console.log('[Info] websocket.onmessage: e.data: ' + JSON.stringify(e.data))
      if (message.api === 'QUEUE_LOAD') {
        this.setState({playlist: []});
      } else if (message.api === 'QUEUE_INSERT') {
        this.setState({playlist: this.state.playlist.concat(message.title)})
      } else {
        console.error('unknown api: ' + message.api)
      }
      console.log('[Info] websocket.onmessage: this.state.playlist: ' + this.state.playlist)
    } catch(err) {
      console.error('JSON.parse failed: ' + e.data)
    }
  }

  handleChange(event: React.FormEvent<HTMLInputElement>) {
    this.setState({inputURL: event.currentTarget.value});
  }

  handleClick() {
    // const url = this.state.inputURL
    // const data = {url}
    const data = {url: this.state.inputURL};
    this.state.websocket.send(JSON.stringify(data));
  }

  render(): JSX.Element {
    return (
      <div>
        <Input
          value = {this.state.inputURL}
          handleChange = {(e) => this.handleChange(e)}
        />
        <Button handleClick = {() => this.handleClick()} />
        <Playlist playlist = {this.state.playlist} />
      </div>
    );
  }
}

interface InputProps {
  value: string;
  handleChange: (event: React.FormEvent<HTMLInputElement>) => void;
}
const Input: React.SFC<InputProps> = props => {
  const { value, handleChange }: InputProps = props;
  return (
    <input
      type="text"
      placeholder="Input Youtube URL"
      value={value}
      onChange={handleChange}
    />
  );
};
 
interface ButtonProps {
  handleClick: () => void;
}
const Button: React.SFC<ButtonProps> = props => {
  const { handleClick }: ButtonProps = props;
  return <button onClick={handleClick}>Say Hello</button>;
};

interface PlaylistProps {
  playlist: string[];
}
const Playlist: React.SFC<PlaylistProps> = props => {
  const { playlist }: PlaylistProps = props;
  const playlistElements = playlist.map(title => {
    return (
      <li key={title}>
        {title}
      </li>
    )
  });
  return (
    <div className="playlist">
      <ol>
        {playlistElements}
      </ol>
    </div>
  );
};
 
ReactDOM.render(
  <MusicPlayer />,
  document.getElementById('app')
  // document.querySelector('.app'),
);