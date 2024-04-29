function playMusic(trackPreviewURL) {
    var audioPlayer = document.getElementById('musicPlayer');
    audioPlayer.src = trackPreviewURL;
    audioPlayer.play();
}