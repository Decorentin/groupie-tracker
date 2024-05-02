function playMusic(trackPreviewURL) {
    var audioPlayer = document.getElementById('musicPlayer');
    audioPlayer.src = trackPreviewURL;
    audioPlayer.play();
}

function rewindMusic() {
    var audioPlayer = document.getElementById('musicPlayer');
    audioPlayer.currentTime -= 5;
}

function forwardMusic() {
    var audioPlayer = document.getElementById('musicPlayer');
    audioPlayer.currentTime += 5;
}

function togglePlayPause() {
    var audioPlayer = document.getElementById('musicPlayer');
    if (audioPlayer.paused) {
        audioPlayer.play();
    } else {
        audioPlayer.pause();
    }
}

function verifyAnswer(guessTrack) {
    var resultDiv = document.getElementById("result");

    var correctTrackName = document.getElementById('blindtest').getAttribute('data-track');

    if (guessTrack.toLowerCase() === correctTrackName.toLowerCase()) {
        resultDiv.innerHTML = "Félicitations ! Vous avez deviné la chanson correctement : " + correctTrackName;
    } else {
        resultDiv.innerHTML = "Désolé, votre réponse est incorrecte.";
    }
}

document.getElementById("blindTestForm").addEventListener("submit", function(event) {
    event.preventDefault();
    var formData = new FormData(event.target);
    var guessTrack = formData.get("guessTrack");

    verifyAnswer(guessTrack);
});

function restartGame() {
    location.reload();
}
