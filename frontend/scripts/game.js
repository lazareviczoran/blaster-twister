import * as Paper from 'paper';

const WIDTH = 300;
const HEIGHT = 400;
const SIDE_COUNT = 3;
const WEBSOCKET_PROTOCOL = window.location.hostname === 'localhost' ? 'ws' : 'wss';
const WEBSOCKET_BASE_URL = `${WEBSOCKET_PROTOCOL}://${document.location.host}/ws`;
const { Point, PointText, Path } = Paper;

const playerPos = {};
const currentPaths = {};
const gameId = document.location.pathname.substring(3);
const clientId = new Date().getTime();
let playerId;
let textItem;
let ws;

const createOrMoveTriangle = (pId, { x, y, rotation }) => {
  const playerTriangle = playerPos[pId];
  if (playerTriangle) {
    playerTriangle.position = new Point(x, y);
    playerTriangle.rotation = rotation + 90;
  } else {
    const strokeColor = pId === '0' ? 'purple' : 'aliceblue';
    const fillColor = pId === '0' ? 'skyblue' : 'yellow';
    const triangle = new Path.RegularPolygon({
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
      playerPath.add(new Point(x, y));
    } else {
      const path = new Path();
      path.strokeColor = pId === '0' ? 'green' : 'red';
      path.add(new Point(x, y));
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
const createMessage = (content) => {
  const text = new PointText(new Point(0, 0));
  text.visible = false;
  text.content = content;
  text.fontSize = 20;
  const itemSize = text.handleBounds;
  text.remove();
  return new PointText({
    point: [WIDTH / 2 - Math.round(itemSize.width / 2), HEIGHT / 2],
    content,
    fillColor: 'white',
    fontSize: 20,
  });
};
const drawWinner = (winnerId, actualPlayerId) => {
  let content;
  if (winnerId === actualPlayerId) {
    content = 'You won!! :)';
  } else {
    content = `Player ${winnerId + 1} won!`;
  }
  createMessage(content);
};

window.addEventListener('load', () => {
  const canvas = document.getElementById('canvas');
  canvas.width = WIDTH;
  canvas.height = HEIGHT;
  Paper.setup(canvas);
  ws = new WebSocket(`${WEBSOCKET_BASE_URL}/${gameId}`);
  ws.onopen = () => {
    ws.send(JSON.stringify({ clientId: clientId.toString() }));
  };
  ws.onclose = () => {
    ws = null;
  };
  ws.onmessage = (evt) => {
    const status = JSON.parse(evt.data);
    if (status.winner != null) {
      drawWinner(status.winner, playerId);
    } else if (status.countdown != null) {
      const content = `Game starts in ${status.countdown}`;
      if (!textItem) {
        textItem = createMessage(content);
      } else if (status.countdown) {
        textItem.content = content;
      } else {
        textItem.remove();
        textItem = null;
      }
    } else {
      const playerKeys = Object.keys(status.players);
      if (playerId == null) {
        const myPlayer = Object.values(status.players).find((p) => p.clientId === clientId);
        if (myPlayer) {
          playerId = parseInt(playerKeys[0], 10);
        }
      }
      movePlayers(status.players);
    }
  };
  ws.onerror = (evt) => {
    // eslint-disable-next-line no-console
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
