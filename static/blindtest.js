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

function removeAccents(str) {
    return str.normalize("NFD").replace(/[\u0300-\u036f]/g, "");
}

function verifyAnswer(guessTrack) {
    var resultDiv = document.getElementById("result");
    var correctTrackName = removeAccents(document.getElementById('blindtest').getAttribute('data-track'));

    if (removeAccents(guessTrack.toLowerCase()) === correctTrackName.toLowerCase()) {
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
