// Variable to track if the form has been submitted
let submitted = false;

// Function to start the countdown timer
function startCountdown() {
    let countDownSeconds = 15;
    let countdownInterval = setInterval(function() {
        document.getElementById('countdown').innerText = countDownSeconds;
        countDownSeconds--;
        if (countDownSeconds < 0) {
            clearInterval(countdownInterval);
            // Reload the page after the countdown is over
            window.location.reload();
        }
    }, 1000); // 1000 milliseconds = 1 second
}

// Start the countdown timer
startCountdown();

// Initialize the score and faults from local storage
let score = localStorage.getItem('score') ? parseInt(localStorage.getItem('score')) : 0;
let faults = localStorage.getItem('faults') ? parseInt(localStorage.getItem('faults')) : 0;
document.getElementById('score').innerText = "Score: " + score;
document.getElementById('faults').innerText = "Faults: " + faults;

// Add event listener for the blind test form submission
document.getElementById('blindTestForm').addEventListener('submit', function(event) {
    event.preventDefault(); 

    // Ensure the form is submitted only once
    if (!submitted) {
        let guessTrack = document.getElementById('guessTrack').value;

        // Verify the user's answer
        verifyAnswer(guessTrack);

        // Mark form as submitted
        submitted = true;
    }
});

// Function to reset form submission status
function resetSubmission() {
    submitted = false;
}

// Function to verify the user's answer
function verifyAnswer(guessTrack) {
    let resultDiv = document.getElementById("result");
    let correctTrackName = removeAccents(document.getElementById('blindtest').getAttribute('data-track'));

    // Check if the user's guess is correct
    if (removeAccents(guessTrack.trim().toLowerCase()) === correctTrackName.trim().toLowerCase()) {
        resultDiv.innerHTML = "Congratulations! You guessed the song correctly: " + correctTrackName;

        // Increment score and update UI
        score++;
        document.getElementById('score').innerText = "Score: " + score;
        localStorage.setItem('score', score);

        // Redirect to scoreboard if score is equal or greater than 3
        if (score >= 3) { 
            window.location.href = "/scoreboard" + "?score=" + score;
        }
    } else {
        // Handle incorrect guess
        resultDiv.innerHTML = "Sorry, your answer is incorrect.";

        // Increment faults and update UI
        faults++;
        document.getElementById('faults').innerText = "Faults: " + faults;
        localStorage.setItem('faults', faults);
      
        // Redirect to "lose" page if score or faults exceeds the threshold
        if (score >= 3 || faults >= 5) {
            resultDiv.innerHTML = "Too many faults or attempts. You lost!";
            window.location.href = "/lose";
        }

        // Reset score if it reaches 5
        if (score >= 5) {
            score = 0;
            localStorage.setItem('score', score);
        }
    }

    // Reset form submission status
    resetSubmission();
}

// Function to restart the game by resetting score and faults
function restartGame() {
    localStorage.removeItem('faults'); // Remove faults from local storage
    location.reload(); // Reload the page
}

// Function to play music
function playMusic(trackPreviewURL) {
    let audioPlayer = document.getElementById('musicPlayer');
    audioPlayer.src = trackPreviewURL;
    audioPlayer.play();
}

// Function to rewind music by 5 seconds
function rewindMusic() {
    let audioPlayer = document.getElementById('musicPlayer');
    audioPlayer.currentTime -= 5;
}

// Function to forward music by 5 seconds
function forwardMusic() {
    let audioPlayer = document.getElementById('musicPlayer');
    audioPlayer.currentTime += 5;
}

// Function to toggle play/pause of music
function togglePlayPause() {
    let audioPlayer = document.getElementById('musicPlayer');
    if (audioPlayer.paused) {
        audioPlayer.play();
    } else {
        audioPlayer.pause();
    }
}

// Function to remove accents from a string
function removeAccents(str) {
    return str.normalize("NFD").replace(/[\u0300-\u036f]/g, "");
}

// Function to reset score and faults
function resetScore() {
    localStorage.setItem('score', 0); // Reset score
    localStorage.setItem('faults', 0); // Reset faults
    location.reload(); // Reload the page
}

// Event listener for the quit button
let quitButton = document.getElementById('quitButton');
quitButton.addEventListener('click', function() {
    // Reset score and faults
    resetScore();
    // Redirect to home page
    window.location.href = "/home";
});
