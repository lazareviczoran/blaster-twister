window.addEventListener("load", function (evt) {
  const canvas = document.getElementById('canvas');
  canvas.width = 300;
  canvas.height = 400;
  const ctx = canvas.getContext('2d');
  var drawTriangle = function (pId, x, y, rotation) {
    var size = 5;
    var sideCount = 3;
    var strokeWidth=4;
    var strokeColor='purple';
    var fillColor='skyblue';
    var radians=-rotation*Math.PI/180;
    ctx.translate(x, y);
    ctx.rotate(radians);
    ctx.beginPath();
    ctx.moveTo (size * Math.cos(0), size * Math.sin(0));
    for (var i = 1; i <= sideCount;i += 1) {
        ctx.lineTo (size * Math.cos(i * 2 * Math.PI / sideCount), size * Math.sin(i * 2 * Math.PI / sideCount));
    }
    ctx.closePath();
    ctx.fillStyle=fillColor;
    ctx.strokeStyle = strokeColor;
    ctx.lineWidth = strokeWidth;
    ctx.stroke();
    ctx.fill();
    ctx.rotate(-radians);
    ctx.translate(-x, -y);
  };
  var movePlayers = function (players) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    Object.entries(players).forEach(function (entry) {
      var id = entry[0];
      var p = entry[1];
      drawTriangle(id, p.x, p.y, p.rotation);
    })
  }
  var gameId = document.location.pathname.substring(3);
  var ws = new WebSocket("ws://" + document.location.host + "/ws/" + gameId);
  var playerId;
  ws.onopen = function (evt) {
    console.log("OPEN");
  };
  ws.onclose = function (evt) {
    console.log("CLOSE");
    ws = null;
  };
  ws.onmessage = function (evt) {
    console.log("RESPONSE: " + evt.data);
    var status = JSON.parse(evt.data);
    var playerKeys = Object.keys(status.players);
    if (!playerId) {
      playerId = playerKeys.length === 1 ? playerKeys[0]:playerKeys[1]
    }
    movePlayers(status.players)
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