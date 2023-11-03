document.addEventListener("DOMContentLoaded", function() {
    // Configura el evento de escucha una sola vez aquí
    document.getElementById("cuentaButton").addEventListener("click", Register);
});

function Register() {

    const Nombre = document.getElementById("Nombre").value;
    const Apellido = document.getElementById("Apellido").value;
    const UserName = document.getElementById("Usuario").value;
    const Password = document.getElementById("Contraseña").value;
    const Email = document.getElementById("email").value;
    const Sexo =  document.getElementById("Sexo").value;
    const Nacionalidad = document.getElementById("Nacionalidad").value;
    const Provincia = document.getElementById("Provincia").value;
    const Ciudad = document.getElementById("Ciudad").value;
    const Domicilio = document.getElementById("Domicilio").value;
  
    const datos = JSON.stringify({
        Nombre: Nombre,
        Apellido: Apellido,
        UserName: UserName,
        Password: Password,
        Email:  Email,
        Sexo: Sexo,
        Nacionalidad: Nacionalidad,
        Provincia: Provincia,
        Ciudad: Ciudad,
        Domicilio: Domicilio,
    });
  
    fetch("/reg/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: datos,
    })
    .then(response => {
        console.log(response);
        if (!response.ok) { // Verifica si la respuesta es exitosa (código 200-299)
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
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
        console.log(error);
        console.error("Error:", error);
    });
  
  }


