const WIDTH = 300;
const HEIGHT = 400;
const TRIANGLE_SIZE = 2;
const CLEAR_TRIANGLE_SIZE = 3;
const SIDE_COUNT = 3;
const STROKE_WIDTH = 4;
const CLEAR_STROKE_WIDTH = 5;

window.addEventListener("load", function (evt) {
  const board = new Array(0);
  const previousTrianglePos = {};
  const canvas = document.getElementById('canvas');
  canvas.width = WIDTH;
  canvas.height = HEIGHT;
  const ctx = canvas.getContext('2d');
  const drawTriangle = function (pId, {x, y, rotation}, clear) {
    const radians=rotation*Math.PI/180;
    const size = clear?CLEAR_TRIANGLE_SIZE:TRIANGLE_SIZE
    let strokeColor= pId === '0'?'purple':'aliceblue';
    let fillColor= pId === '0'?'skyblue':'yellow';
    if (clear) {
      strokeColor = 'black'
      fillColor = 'black'
    }
    ctx.translate(x, y);
    ctx.rotate(radians);
    ctx.beginPath();
    ctx.moveTo (size * Math.cos(0), size * Math.sin(0));
    for (let i = 1; i <= SIDE_COUNT;i += 1) {
      ctx.lineTo(
        size * Math.cos(i * 2 * Math.PI / SIDE_COUNT),
        size * Math.sin(i * 2 * Math.PI / SIDE_COUNT)
      );
    }
    ctx.closePath();
    ctx.fillStyle=fillColor;
    ctx.strokeStyle = strokeColor;
    ctx.lineWidth = clear?CLEAR_STROKE_WIDTH:STROKE_WIDTH;
    ctx.stroke();
    ctx.fill();
    ctx.rotate(-radians);
    ctx.translate(-x, -y);
  };
  const clearTriangle = function (pId) {
    const prevPos = previousTrianglePos[pId];
    if (prevPos) {
      drawTriangle(pId, prevPos, true)
    }
  }
  const movePlayers = function (players) {
    Object.entries(players).forEach(function (entry) {
      const [id] = entry;
      clearTriangle(id)
    });

    drawVisitedPositions();

    Object.entries(players).forEach(function (entry) {
      const [id, p] = entry;
      markFieldAsUsed(id, p)
      drawTriangle(id, p);
    })
  }
  const markFieldAsUsed = function (pId, {x, y, trace, rotation}) {
    previousTrianglePos[pId] = {x, y, rotation}
    board.push({x, y, pId, trace});
  };
  const drawVisitedPositions = function () {
    board.forEach(function (pos) {
      if (pos.trace) {
        ctx.fillStyle = pos.pId === "0"?"green":"red";
        ctx.fillRect(pos.x,pos.y,1,1);
      }
    });
  };
  const drawWinner = function (pId) {
    console.log(pId);
  }
  const gameId = document.location.pathname.substring(3);
  let ws = new WebSocket("wss://" + document.location.host + "/ws/" + gameId);
  let playerId;
  ws.onopen = function (evt) {
    console.log("OPEN");
  };
  ws.onclose = function (evt) {
    console.log("CLOSE");
    ws = null;
  };
  ws.onmessage = function (evt) {
    const status = JSON.parse(evt.data);
    if (status.winner != null) {
      drawWinner(status.winner);
    } else {
      const playerKeys = Object.keys(status.players);
      if (!playerId) {
        playerId = playerKeys.length === 1 ? playerKeys[0]:playerKeys[1]
      }
      movePlayers(status.players);
    }
  };
  ws.onerror = function (evt) {
    console.log("ERROR: " + evt.data);
  };
  document.onkeydown = function(event){
    if (ws) {
      if (event.repeat) { return }
      if (event.keyCode === 37) {
        ws.send(JSON.stringify({dir: 'down', key: 'left'}));
      } else if (event.keyCode === 39) {
        ws.send(JSON.stringify({dir: 'down', key: 'right'}));
      }
    }
  };
  document.onkeyup = function(event){
    if (ws) {
      if (event.keyCode === 37) {
        ws.send(JSON.stringify({dir: 'up', key: 'left'}));
      } else if (event.keyCode === 39) {
        ws.send(JSON.stringify({dir: 'up', key: 'right'}));
      }
    }
  };
});