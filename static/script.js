function startCountdown() {
    let countDownSeconds = 15;
    let countdownInterval = setInterval(function() {
        document.getElementById('countdown').innerText = countDownSeconds;
        countDownSeconds--;
        if (countDownSeconds < 0) {
            clearInterval(countdownInterval);
            window.location.reload();
        }
    }, 1000);
}


window.onload = startCountdown;

let score = localStorage.getItem('score') ? parseInt(localStorage.getItem('score')) : 0;
document.getElementById('score').innerText = "Score: " + score;

let wrongAttempts = localStorage.getItem('wrongAttempts') ? parseInt(localStorage.getItem('wrongAttempts')) : 0;
document.getElementById('attempts').innerText = "Essais restants : " + (5 - wrongAttempts);

document.getElementById('songForm').addEventListener('submit', function(event) {
    event.preventDefault(); 

    let userAnswer = document.getElementById('userAnswer').value;

    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/check-answer');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.onload = function() {
        document.getElementById('resultMessage').innerText = this.responseText;
    
        if (this.responseText.trim() === "Bravo, vous avez deviné la bonne chanson !") {
            score++; 
            document.getElementById('score').innerText = "Score: " + score; 
            localStorage.setItem('score', score);

            document.getElementById('attempts').innerText = "Essais restants : 5";
        } else {
            wrongAttempts++;

            let remainingAttempts = 5 - wrongAttempts;
            document.getElementById('attempts').innerText = "Essais restants : " + remainingAttempts;

            localStorage.setItem('wrongAttempts', wrongAttempts);

            if (wrongAttempts === 5) {
                let losePageURL = "/lose";
                window.location.href = losePageURL; 
            }
        }

        if (score === 5) {
            let winPageURL = "/win?score=" + score;
            window.location.href = winPageURL;
        } else {
            setTimeout(function() {
                window.location.reload();
            }, 1000);
        }
    };    
    xhr.send('userAnswer=' + encodeURIComponent(userAnswer));
});

function resetScore() {
    localStorage.setItem('score', 0);
}

let quitButton = document.getElementById('quitButton');

quitButton.addEventListener('click', function() {
    resetScore();
    localStorage.setItem('wrongAttempts', 0);
    window.location.href = "/home";
});