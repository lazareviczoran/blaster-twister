const WIDTH = 300;
const HEIGHT = 400;
const TRIANGLE_SIZE = 2;
const SIDE_COUNT = 3;
const STROKE_WIDTH = 4;

window.addEventListener("load", function (evt) {
  const board = new Array(0);
  const canvas = document.getElementById('canvas');
  canvas.width = WIDTH;
  canvas.height = HEIGHT;
  const ctx = canvas.getContext('2d');
  const drawTriangle = function (pId, x, y, rotation) {
    markFieldAsUsed(pId, x, y)
    const strokeColor= pId === '0'?'purple':'aliceblue';
    const fillColor= pId === '0'?'skyblue':'yellow';
    const radians=rotation*Math.PI/180;
    ctx.translate(x, y);
    ctx.rotate(radians);
    ctx.beginPath();
    ctx.moveTo (TRIANGLE_SIZE * Math.cos(0), TRIANGLE_SIZE * Math.sin(0));
    for (let i = 1; i <= SIDE_COUNT;i += 1) {
        ctx.lineTo (TRIANGLE_SIZE * Math.cos(i * 2 * Math.PI / SIDE_COUNT), TRIANGLE_SIZE * Math.sin(i * 2 * Math.PI / SIDE_COUNT));
    }
    ctx.closePath();
    ctx.fillStyle=fillColor;
    ctx.strokeStyle = strokeColor;
    ctx.lineWidth = STROKE_WIDTH;
    ctx.stroke();
    ctx.fill();
    ctx.rotate(-radians);
    ctx.translate(-x, -y);
  };
  const movePlayers = function (players) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    drawVisitedPositions();

    Object.entries(players).forEach(function (entry) {
      const id = entry[0];
      const p = entry[1];
      drawTriangle(id, p.x, p.y, p.rotation);
    })
  }
  const markFieldAsUsed = function (pId, x, y) {
    board.push({x, y, pId});
  };
  const drawVisitedPositions = function () {
    board.forEach(function (pos) {
      ctx.fillStyle = pos.pId === "0"?"green":"red";
      ctx.fillRect(pos.x,pos.y,1,1);
    });
  };
  const drawWinner = function (pId) {
    console.log(pId);
  }
  const gameId = document.location.pathname.substring(3);
  let ws = new WebSocket("ws://" + document.location.host + "/ws/" + gameId);
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
    if (status.winner) {
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