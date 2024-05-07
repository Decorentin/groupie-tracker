// Déclaration de la variable pour suivre l'état de la soumission
var submitted = false;

// Fonction pour démarrer le compteur
function startCountdown() {
    // Démarrer le compteur (15 secondes)
    var countDownSeconds = 5;
    var countdownInterval = setInterval(function() {
        document.getElementById('countdown').innerText = countDownSeconds;
        countDownSeconds--;
        if (countDownSeconds < 0) {
            clearInterval(countdownInterval);
            // Actualiser la page après que le compteur soit écoulé
            window.location.reload();
        }
    }, 1000); // 1000 millisecondes = 1 seconde
}

// Démarrer le compteur dès que la page est chargée
startCountdown();

// Initialiser le score à partir du stockage local
var score = localStorage.getItem('score') ? parseInt(localStorage.getItem('score')) : 0;
var faults = localStorage.getItem('faults') ? parseInt(localStorage.getItem('faults')) : 0;
document.getElementById('score').innerText = "Score: " + score;
document.getElementById('faults').innerText = "Fautes: " + faults;

// Ajouter les événements nécessaires pour le formulaire de blind test
document.getElementById('blindTestForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Empêcher le formulaire de se soumettre normalement

    // Vérifier si la soumission a déjà été effectuée
    if (!submitted) {
        // Récupérer la valeur de la réponse de l'utilisateur
        var guessTrack = document.getElementById('guessTrack').value;

        // Vérifier la réponse de l'utilisateur
        verifyAnswer(guessTrack);

        // Marquer la soumission comme effectuée
        submitted = true;
    }
});

// Fonction pour réinitialiser l'état de la soumission
function resetSubmission() {
    submitted = false;
}

// Fonction pour vérifier la réponse de l'utilisateur
function verifyAnswer(guessTrack) {
    var resultDiv = document.getElementById("result");
    var correctTrackName = removeAccents(document.getElementById('blindtest').getAttribute('data-track'));

    if (removeAccents(guessTrack.trim().toLowerCase()) === correctTrackName.trim().toLowerCase()) {
        resultDiv.innerHTML = "Félicitations ! Vous avez deviné la chanson correctement : " + correctTrackName;
        // Si la réponse est correcte, incrémenter le score et afficher le nouveau score
        score++;
        document.getElementById('score').innerText = "Score: " + score;
        localStorage.setItem('score', score);

        if (score >= 3) { // Rediriger vers le scoreboard lorsque le score atteint 3
            window.location.href = "/scoreboard";
        }
    } else {
        resultDiv.innerHTML = "Désolé, votre réponse est incorrecte.";
        // Si la réponse est incorrecte, incrémenter le compteur de fautes
        faults++;
        document.getElementById('faults').innerText = "Fautes: " + faults;
        localStorage.setItem('faults', faults);
        // Vérifier si le joueur a atteint 3 fautes ou 5 tentatives
        if (score >= 3 || faults >= 5) {
            resultDiv.innerHTML = "Trop de fautes ou de tentatives. Vous avez perdu !";
            // Rediriger vers la page de défaite
            window.location.href = "/loose";
        }
        // Si le score atteint 5, réinitialiser le score pour ne pas atteindre 5 fautes
        if (score >= 5) {
            score = 0;
            localStorage.setItem('score', score);
        }
    }

    // Réinitialiser l'état de la soumission pour permettre une nouvelle vérification
    resetSubmission();
}

// Fonction pour redémarrer le jeu
function restartGame() {
    localStorage.removeItem('faults'); // Réinitialiser le compteur de fautes
    location.reload();
}

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

// Fonction pour réinitialiser le score et le compteur de fautes à zéro
function resetScore() {
    localStorage.setItem('score', 0);
    localStorage.setItem('faults', 0);
    location.reload();
}

// Récupérer le bouton Quitter
var quitButton = document.getElementById('quitButton');

// Ajouter un écouteur d'événements au clic sur le bouton Quitter
quitButton.addEventListener('click', function() {
    // Réinitialiser le score et les fautes avant de quitter
    resetScore();
    // Rediriger l'utilisateur vers la page d'accueil
    window.location.href = "/home"; // Mettez ici le chemin de votre page d'accueil
});
