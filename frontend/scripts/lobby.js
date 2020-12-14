const BASE_URL = window.location.origin;
const message = 'Waiting for players to join';
const unsuccessfulMessage = 'No available players at the moment.';
const WEBSOCKET_PROTOCOL = window.location.hostname === 'localhost' ? 'ws' : 'wss';
const WEBSOCKET_BASE_URL = `${WEBSOCKET_PROTOCOL}://${window.location.host}/ws`;
let ws;

window.addEventListener('load', () => {
  const progressContainer = document.getElementById('lobby');
  const progress = document.getElementById('progress-bar');
  let progressStatus = 0;
  document.getElementById('lobby-message').innerHTML = message;
  document.getElementById('unsuccessful-message').innerHTML = unsuccessfulMessage;
  document.getElementById('retry').setAttribute('href', '/join');

  setTimeout(() => {
    ws = new WebSocket(`${WEBSOCKET_BASE_URL}/lobby`);
    ws.onopen = () => {};
    ws.onclose = () => {
      ws = null;
    };
    ws.onmessage = (evt) => {
      if (!evt.data) {
        return;
      }
      if (evt.data === 'ping') {
        ws.send('success');
      } else {
        window.location = `${BASE_URL}/g/${evt.data}`;
      }
    };
    // eslint-disable-next-line no-console
    ws.onerror = console.error;
  }, 1200);

  const intervalId = setInterval(() => {
    progressStatus += 1;
    if (progressStatus > 60) {
      clearInterval(intervalId);
      document.getElementById('unsuccessful').classList.remove('d-none');
      progressContainer.classList.add('d-none');
      ws = null;
    }
    progress.setAttribute('aria-valuenow', `${progressStatus}`);
    progress.style.width = `${Math.round(progressStatus / 0.6)}%`;
  }, 1000);
});
