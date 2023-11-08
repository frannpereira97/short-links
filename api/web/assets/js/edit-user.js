document.addEventListener("DOMContentLoaded", function() {
    const token = localStorage.getItem('x-jwt-token');
    if (token) {
        fetch('/users/validate', {
            method: 'POST', 
            headers: {
                'Content-Type': 'application/json',
                'x-jwt-token': token 
            }
        })
        .then(response => {
            if (response.ok) {
                getUserData()
                getDatosUser()
            } else {
                // El token no es válido o ha ocurrido un error
                window.location.href = '/';
                throw new Error('Token no válido o error en la solicitud');
            }
        })
        .catch(error => {
            console.error('Error en la validación del token:', error);
        });
    } else {
        window.location.href = '/';
    }
});

function getUserData() {
    const token = localStorage.getItem('x-jwt-token');
    fetch('/users/list', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'x-jwt-token': token,
        },
    })
    .then(response => response.json())
    .then(data => {
        if (data.admin == false) {
            document.getElementById("user-access").style.display = "block";
        } else{
            document.getElementById("admin-access").style.display = "block";
            fillTable(data);
        }
    })
    .catch(error => {
        console.error("Error al obtener los datos:", error);
    });
}

function fillTable(data) {
    var table = document.getElementById("users-table");

    while (table.rows.length > 1) {
        table.deleteRow(1);
    }
    // Itera a través de los datos y crea filas en la tabla.
    data.forEach(item => {
        var row = table.insertRow();
        var cell1 = row.insertCell(0);
        var cell2 = row.insertCell(1);
        var cell3 = row.insertCell(2);
        cell1.textContent = item.id;
        cell2.textContent = item.username;
        cell3.textContent = item.permisos;
    });
}

function getDatosUser() {
    const token = localStorage.getItem('x-jwt-token');
    fetch('/users/data', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'x-jwt-token': token,
        },
    })
    .then(response => response.json())
    .then(data => {
        console.log(data);
        fillUserData(data);
    })
    .catch(error => {
        console.error("Error al obtener los datos:", error);
    });
}

function fillUserData(data) {
    document.getElementById("modal_nombre").value = data.nombre;
    document.getElementById("modal_apellido").value = data.apellido;
    document.getElementById("modal_email").value = data.email;
    document.getElementById("modal_sexo").value = data.sexo;
    document.getElementById("modal_nacionalidad").value = data.nacionalidad;
    document.getElementById("modal_provincia").value = data.provincia;
    document.getElementById("modal_ciudad").value = data.ciudad;
    document.getElementById("modal_domicilio").value = data.Domicilio;
}

function logout() {
    const token = localStorage.getItem('x-jwt-token');

    fetch('/users/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'x-jwt-token': token,
        },
        body: ''
    })
    .then(response => {
        // Verificando si la respuesta es exitosa
        if (!response.ok) {
            throw new Error('Error en la red o en el servidor');
        }
        return response.json();
    })
    .then(data => {
        // Aquí manejamos la respuesta del servidor
        if (data.redirectTo) {
            window.location.href = data.redirectTo;
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function changePW() {
    const token = localStorage.getItem('x-jwt-token');
    const pass_pw = document.getElementById("pass_pw").value;
    const new_pw = document.getElementById("pass_newpw").value;
    const new_pw2 = document.getElementById("pass_vnewpw").value;

    const datos = JSON.stringify({
        Password: pass_pw,
        NewPass: new_pw,
        VnewPass: new_pw2,
});
    if (new_pw == new_pw2) {
        fetch('/users/changepw', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'x-jwt-token': token,
            },
            body: datos
        })
        .then(response => {
            // Verificando si la respuesta es exitosa
            if (!response.ok) {
                throw new Error('Error en la red o en el servidor');
            }
            return response.json();
        })
        .catch(error => {
            console.error('Error:', error);
        });
    } else {
        console.error('Error:', "Las contraseñas no coinciden");
    }
    

}