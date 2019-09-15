import * as Paper from 'paper';

const WIDTH = 300;
const HEIGHT = 400;
const SIDE_COUNT = 3;
const WEBSOCKET_PROTOCOL = window.location.hostname === 'localhost' ? 'ws' : 'wss';
const WEBSOCKET_BASE_URL = `${WEBSOCKET_PROTOCOL}://${document.location.host}/ws`;

window.addEventListener('load', () => {
  const playerPos = {};
  const currentPaths = {};
  const canvas = document.getElementById('canvas');
  canvas.width = WIDTH;
  canvas.height = HEIGHT;
  Paper.setup(canvas);
  const createOrMoveTriangle = (pId, { x, y, rotation }) => {
    const playerTriangle = playerPos[pId];
    if (playerTriangle) {
      playerTriangle.position = new Paper.Point(x, y);
      playerTriangle.rotation = rotation + 90;
    } else {
      const strokeColor = pId === '0' ? 'purple' : 'aliceblue';
      const fillColor = pId === '0' ? 'skyblue' : 'yellow';
      const triangle = new Paper.Path.RegularPolygon({
        center: [x, y],
        sides: SIDE_COUNT,
        radius: 5,
        fillColor,
        strokeColor,
        applyMatrix: false,
      });
      triangle.rotation = rotation + 90;
      playerPos[pId] = triangle;
    }
  };
  const markFieldAsUsed = (pId, { x, y, trace }) => {
    const playerPath = currentPaths[pId];
    if (trace) {
      if (playerPath) {
        playerPath.add(new Paper.Point(x, y));
      } else {
        const path = new Paper.Path();
        path.strokeColor = pId === '0' ? 'green' : 'red';
        path.add(new Paper.Point(x, y));
        currentPaths[pId] = path;
      }
    } else if (playerPath) {
      currentPaths[pId] = null;
    }
  };
  const movePlayers = (players) => {
    Object.entries(players).forEach((entry) => {
      const [id, p] = entry;
      markFieldAsUsed(id, p);
      createOrMoveTriangle(id, p);
    });
  };
  const drawWinner = (winnerId, actualPlayerId) => {
    if (winnerId === actualPlayerId) {
      alert('You won!! :)');
    } else {
      alert('You lost!! :(');
    }
  };
  const gameId = document.location.pathname.substring(3);
  let ws = new WebSocket(`${WEBSOCKET_BASE_URL}/${gameId}`);
  let playerId;
  ws.onopen = () => {
    console.log('OPEN');
  };
  ws.onclose = () => {
    console.log('CLOSE');
    ws = null;
  };
  ws.onmessage = (evt) => {
    const status = JSON.parse(evt.data);
    if (status.winner != null) {
      drawWinner(status.winner, playerId);
    } else if (status.countdown != null) {
      console.log('Game starts in ', status.countdown);
    } else {
      const playerKeys = Object.keys(status.players);
      if (playerId == null) {
        playerId = parseInt(playerKeys[0], 10);
      }
      movePlayers(status.players);
    }
  };
  ws.onerror = (evt) => {
    console.log(`ERROR: ${evt.data}`);
  };
  document.onkeydown = (event) => {
    if (ws) {
      if (event.repeat) { return; }
      if (event.keyCode === 37) {
        ws.send(JSON.stringify({ dir: 'down', key: 'left' }));
      } else if (event.keyCode === 39) {
        ws.send(JSON.stringify({ dir: 'down', key: 'right' }));
      }
    }
  };
  document.onkeyup = (event) => {
    if (ws) {
      if (event.keyCode === 37) {
        ws.send(JSON.stringify({ dir: 'up', key: 'left' }));
      } else if (event.keyCode === 39) {
        ws.send(JSON.stringify({ dir: 'up', key: 'right' }));
      }
    }
  };
});
