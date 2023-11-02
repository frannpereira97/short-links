function crearShort() {
    // Obteniendo la URL del campo de entrada del formulario
    const urlIntroducida = document.getElementById('urlInput').value;

    const shortIntroducido = document.getElementById('shortInput').value;
    
    const token = localStorage.getItem('x-jwt-token');
    // Datos que quieres enviar al servidor
    const datos = JSON.stringify({
        Pagina: urlIntroducida,
        Short: shortIntroducido,
        Expiry: 30,
        FechaCreaion:'',
        Abierto:'',
        UserID: 1,        
    });

    fetch('/users/Shorten', {
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
    .then(data => {
        // Aquí manejamos la respuesta del servidor
        // const resultado = JSON.stringify(data, null, 2);
        document.getElementById('respuestaShort').hidden = false;
        const respuesta = document.getElementById('respuesta');
        respuesta.href = data.Short;
        respuesta.textContent = data.Short;
        respuesta.target = '_blank';

    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('respuesta').textContent = error;
    });
    listarShorts();
}

function listarShorts() {
    const token = localStorage.getItem('x-jwt-token');
    fetch('/shorts/list', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'x-jwt-token': token,
        },
    })
    .then(response => response.json())
    .then(data => {
        // Llama a la función para rellenar la tabla con los datos obtenidos.
        fillTable(data);
    })
    .catch(error => {
        console.error("Error al obtener los datos:", error);
    });
}

function fillTable(data) {
    var table = document.getElementById("shorts-table");

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

        var link = document.createElement("a");
        link.href = item.short;
        link.textContent = item.short;
        link.target = '_blank';

        cell2.appendChild(link);
        cell3.textContent = item.pagina;
    });
}

window.onload = function() {
    listarShorts();
};

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
