const BASE_URL = window.location.origin;
const isNew = window.location.pathname.match(/new/);
const message = isNew ? 'Waiting for players to join' : 'Searching for available games';
const unsuccessfulMessage = isNew ? 'No available players at the moment.' : 'No available games at the moment.';

window.addEventListener('load', () => {
  const progressContainer = document.getElementById('lobby');
  const progress = document.getElementById('progress-bar');
  let progressStatus = 0;
  document.getElementById('lobby-message').innerHTML = message;
  document.getElementById('unsuccessful-message').innerHTML = unsuccessfulMessage;
  document.getElementById('retry').setAttribute('href', isNew ? '/new' : '/join');

  const intervalId = setInterval(() => {
    progressStatus += 1;
    if (progressStatus > 60) {
      clearInterval(intervalId);
    }
    progress.setAttribute('aria-valuenow', `${progressStatus}`);
    progress.style.width = `${Math.round(progressStatus / 0.6)}%`;
  }, 1000);

  fetch(`${BASE_URL}/api/${window.location.pathname.substring(1)}`)
    .then((response) => {
      clearInterval(intervalId);
      if (response.status === 200) {
        const responseObj = response.json();
        window.location = `${BASE_URL}/g/${responseObj.gameId}`;
      } if (response.status === 422) {
        document.getElementById('unsuccessful').classList.remove('d-none');
        progressContainer.classList.add('d-none');
      } else {
        throw new Error('Something went wrong on api server!');
      }
    }).catch((error) => {
      console.error(error);
    });
});
