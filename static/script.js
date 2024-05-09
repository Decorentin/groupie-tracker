// Fonction pour démarrer le compteur
function startCountdown() {
    // Démarrer le compteur (15 secondes)
    let countDownSeconds = 15;
    let countdownInterval = setInterval(function() {
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
let score = localStorage.getItem('score') ? parseInt(localStorage.getItem('score')) : 0;
document.getElementById('score').innerText = "Score: " + score;

// Initialiser le compteur d'essais incorrects à partir du stockage local s'il existe
let wrongAttempts = localStorage.getItem('wrongAttempts') ? parseInt(localStorage.getItem('wrongAttempts')) : 0;
document.getElementById('attempts').innerText = "Essais restants : " + (5 - wrongAttempts);

document.getElementById('songForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Empêcher le formulaire de se soumettre normalement

    // Récupérer la valeur de l'input de l'utilisateur
    let userAnswer = document.getElementById('userAnswer').value;

    // Effectuer une requête AJAX pour vérifier la réponse côté serveur
    let xhr = new XMLHttpRequest();
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

            // Réinitialiser le compteur d'essais incorrects et mettre à jour le stockage local
            document.getElementById('attempts').innerText = "Essais restants : 5";
        } else {
            // Si la réponse est incorrecte, incrémenter le compteur d'essais incorrects
            wrongAttempts++;

            // Afficher le nombre d'essais restants
            let remainingAttempts = 5 - wrongAttempts;
            document.getElementById('attempts').innerText = "Essais restants : " + remainingAttempts;

            // Mettre à jour le stockage local avec le nouveau compteur d'essais incorrects
            localStorage.setItem('wrongAttempts', wrongAttempts);

            // Vérifier si le joueur a épuisé ses essais
            if (wrongAttempts === 5) {
                // Construire l'URL de redirection vers la page de défaite
                let losePageURL = "/lose";
                window.location.href = losePageURL; // Rediriger l'utilisateur vers la page de défaite
            }
        }

        // Vérifier si le score atteint 5
        if (score === 5) {
            // Construire l'URL de redirection vers la page de victoire avec le score
            let winPageURL = "/win?score=" + score;
            window.location.href = winPageURL; // Rediriger l'utilisateur vers la page de victoire
        } else {
            // Attendre 1 seconde avant de recharger la page
            setTimeout(function() {
                window.location.reload();
            }, 1000);
        }
    };    
    xhr.send('userAnswer=' + encodeURIComponent(userAnswer));
});

// Fonction pour réinitialiser le score
function resetScore() {
    // Réinitialiser le score à zéro
    localStorage.setItem('score', 0);
}

// Récupérer le bouton Quitter
let quitButton = document.getElementById('quitButton');

// Ajouter un écouteur d'événements au clic sur le bouton Quitter
quitButton.addEventListener('click', function() {
    // Réinitialiser le score avant de quitter
    resetScore();
    // Réinitialiser le compteur d'essais incorrects
    localStorage.setItem('wrongAttempts', 0);
    // Rediriger l'utilisateur vers la page d'accueil
    window.location.href = "/home"; // Mettez ici le chemin de votre page d'accueil
});
