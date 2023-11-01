document.addEventListener("DOMContentLoaded", function() {

    document.getElementById('loginButton').addEventListener('click', function(event) {
        event.preventDefault();

        const loginButton = document.getElementById("loginButton");

        const UserName = document.getElementById("floatingInput").value;
        const Password = document.getElementById("floatingPassword").value;

        const datos = JSON.stringify({
            UserName: UserName,
            Password: Password,
        });

        loginButton.addEventListener("click", function() {
            fetch("/users/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                // Si necesitas enviar datos con tu POST, añádelos aquí. Por ejemplo:
                body: datos
            })
            .then(response => response.json())  // asume que el servidor responde con json
            .then(data => {
                console.log(data);
                if (data['x-jwt-token']) {
                    localStorage.setItem('x-jwt-token', data['x-jwt-token']);
                }            
                if (data.redirectTo) {
                    window.location.href = data.redirectTo;
                }
            })
            .catch(error => {
                console.error("Error:", error);
            });
        });
    });
});