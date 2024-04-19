
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
    };
    xhr.send('userAnswer=' + encodeURIComponent(userAnswer));
});
