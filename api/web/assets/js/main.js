document.addEventListener("DOMContentLoaded", function() {

    document.getElementById('loginButton').addEventListener('click', function(event) {
        event.preventDefault();

        const loginButton = document.getElementById("loginButton");

        let UserName = document.getElementById("floatingInput").value;
        let Password = document.getElementById("floatingPassword").value;
        let alertJS = document.getElementById("alertJS");

        let datos = JSON.stringify({
            UserName: UserName,
            Password: Password,
        });

        fetch("/users/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: datos,
        })
        .then(response => response.json())
        .then(data => {
            if (data['x-jwt-token']) {
                localStorage.setItem('x-jwt-token', data['x-jwt-token']);
            }
            if (data.redirectTo) {
                window.location.href = data.redirectTo;
            }
        })
        .catch(error => {
            console.log("error");
            alertJS.hidden = false;
            console.error("Error:", error);
        });
    });
});