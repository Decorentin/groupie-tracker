document.getElementById('songForm').addEventListener('submit', function(event) {
    event.preventDefault(); 

    var userAnswer = document.getElementById('userAnswer').value;

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/check-answer');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.onload = function() {
        document.getElementById('resultMessage').innerText = this.responseText;
    };
    xhr.send('userAnswer=' + encodeURIComponent(userAnswer));
});
