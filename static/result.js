// Fonction pour démarrer le compteur
function startCountdown() {
    // Démarrer le compteur (15 secondes)
    var countDownSeconds = 15;
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

// Appeler la fonction startCountdown lors du chargement de la page
window.onload = startCountdown;

// Initialiser le score à partir du stockage local
var score = localStorage.getItem('score') ? parseInt(localStorage.getItem('score')) : 0;
document.getElementById('score').innerText = "Score: " + score;

document.getElementById('songForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Empêcher le formulaire de se soumettre normalement

    // Récupérer la valeur de l'input de l'utilisateur
    var userAnswer = document.getElementById('userAnswer').value;

    // Effectuer une requête AJAX pour vérifier la réponse côté serveur
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/check-answer');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.onload = function() {
        // Afficher le résultat de la vérification dans la balise div
        document.getElementById('resultMessage').innerText = this.responseText;

        // Si la réponse est correcte, incrémenter le score et afficher le nouveau score
        if (this.responseText.trim() === "Bravo, vous avez deviné la bonne chanson !") {
            score++; // Incrémenter le score
            document.getElementById('score').innerText = "Score: " + score; // Afficher le nouveau score
            // Stocker le score dans le stockage local
            localStorage.setItem('score', score);
        }

        // Actualiser la page après 1 seconde
        setTimeout(function() {
            window.location.reload();
        }, 1000);
    };
    xhr.send('userAnswer=' + encodeURIComponent(userAnswer));
});


// Fonction pour réinitialiser le score
function resetScore() {
    // Réinitialiser le score à zéro
    localStorage.setItem('score', 0);
}

// Récupérer le bouton Quitter
var quitButton = document.getElementById('quitButton');

// Ajouter un écouteur d'événements au clic sur le bouton Quitter
quitButton.addEventListener('click', function() {
    // Réinitialiser le score avant de quitter
    resetScore();
    // Rediriger l'utilisateur vers la page d'accueil
    window.location.href = "/home"; // Mettez ici le chemin de votre page d'accueil
});