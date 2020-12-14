import * as Paper from 'paper';

const WIDTH = 500;
const HEIGHT = 600;
const LEFT = 'left';
const RIGHT = 'right';
const UP = 'up';
const DOWN = 'down';
const LEFT_KEY = 'ArrowLeft';
const RIGHT_KEY = 'ArrowRight';
const WEBSOCKET_PROTOCOL = window.location.hostname === 'localhost' ? 'ws' : 'wss';
const WEBSOCKET_BASE_URL = `${WEBSOCKET_PROTOCOL}://${window.location.host}/ws`;
const {
  Point, PointText, Path, Raster, Layer,
} = Paper;

const playerPos = {};
const currentPaths = {};
const gameId = window.location.pathname.substring(3);
const clientId = new Date().getTime();
let playerId;
let textItem;
let pathLayer;
let iconLayer;
let messageLayer;
let ws;

const createOrMoveTriangle = (pId, { x, y, rotation }) => {
  const playerTriangle = playerPos[pId];
  if (playerTriangle) {
    playerTriangle.position = new Point(x, y);
    playerTriangle.rotation = rotation + 90;
  } else {
    const icon = pId === '0' ? 'rocket1' : 'rocket2';
    const playerIcon = new Raster(icon);
    playerIcon.position = new Point(x, y);
    playerIcon.rotation = rotation + 90;
    playerIcon.scale(0.15);
    iconLayer.addChild(playerIcon);
    playerPos[pId] = playerIcon;
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
      path.strokeWidth = 2;
      path.add(new Point(x, y));
      pathLayer.addChild(path);
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
  const messageItem = createMessage(content);
  messageLayer.addChild(messageItem);
  document.getElementById('back').classList.remove('d-none');
};

window.addEventListener('load', () => {
  const canvas = document.getElementById('canvas');
  canvas.width = WIDTH;
  canvas.height = HEIGHT;
  Paper.setup(canvas);
  pathLayer = new Layer();
  iconLayer = new Layer();
  messageLayer = new Layer();

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
        messageLayer.addChild(textItem);
      } else if (status.countdown) {
        textItem.content = content;
      } else {
        textItem.remove();
        textItem = null;
      }
    } else {
      const playerKeys = Object.keys(status.players);
      const playerSpan = document.getElementById(`player${playerKeys[0]}`);
      let playerText = 'Opponent';
      if (playerId == null) {
        const myPlayer = Object.values(status.players).find((p) => p.clientId === clientId);
        if (myPlayer) {
          playerId = parseInt(playerKeys[0], 10);
          playerText = 'Me';
          playerSpan.innerHTML = playerText;
        }
      }
      if (playerSpan.innerHTML === '') {
        playerSpan.innerHTML = playerText;
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
      if (event.key === LEFT_KEY) {
        ws.send(JSON.stringify({ dir: DOWN, key: LEFT }));
      } else if (event.key === RIGHT_KEY) {
        ws.send(JSON.stringify({ dir: DOWN, key: RIGHT }));
      }
    }
  };
  document.onkeyup = (event) => {
    if (ws) {
      if (event.key === LEFT_KEY) {
        ws.send(JSON.stringify({ dir: UP, key: LEFT }));
      } else if (event.key === RIGHT_KEY) {
        ws.send(JSON.stringify({ dir: UP, key: RIGHT }));
      }
    }
  };
  const sendRotationEvent = (event, dir, key) => {
    if (event) {
      event.preventDefault();
    }
    if (ws) {
      ws.send(JSON.stringify({ dir, key }));
    }
  };
  document.getElementById(LEFT).addEventListener('mousedown', (e) => {
    sendRotationEvent(e, DOWN, LEFT);
  });
  document.getElementById(LEFT).addEventListener('mouseup', (e) => {
    sendRotationEvent(e, UP, LEFT);
  });
  document.getElementById(RIGHT).addEventListener('mousedown', (e) => {
    sendRotationEvent(e, DOWN, RIGHT);
  });
  document.getElementById(RIGHT).addEventListener('mouseup', (e) => {
    sendRotationEvent(e, UP, RIGHT);
  });
  document.getElementById(LEFT).addEventListener('touchstart', (e) => {
    sendRotationEvent(e, DOWN, LEFT);
  });
  document.getElementById(LEFT).addEventListener('touchend', (e) => {
    sendRotationEvent(e, UP, LEFT);
  });
  document.getElementById(RIGHT).addEventListener('touchstart', (e) => {
    sendRotationEvent(e, DOWN, RIGHT);
  });
  document.getElementById(RIGHT).addEventListener('touchend', (e) => {
    sendRotationEvent(e, UP, RIGHT);
  });
});
