// Code JavaScript pour récupérer les paroles de la chanson et les afficher dans le conteneur des paroles
window.onload = function() {
    // Simulons des paroles de chanson
    var lyrics = "Les paroles de la chanson vont ici...";

    // Sélectionner l'élément qui contiendra les paroles de la chanson
    var lyricsContainer = document.getElementById("lyrics-container");

    // Mettre à jour le contenu de l'élément avec les paroles de la chanson
    lyricsContainer.innerText = lyrics;
};
