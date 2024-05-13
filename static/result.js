// Add an event listener to the form submission event
document.getElementById('songForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Prevent the form from submitting normally

    // Retrieve the value of the user input
    var userAnswer = document.getElementById('userAnswer').value;

    // Perform an AJAX request to check the answer on the server side
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/check-answer');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.onload = function() {
        // Display the result of the verification in the resultMessage div element
        document.getElementById('resultMessage').innerText = this.responseText;
    };
    xhr.send('userAnswer=' + encodeURIComponent(userAnswer));
});
