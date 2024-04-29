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
