function enviarJSON() {
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
        document.getElementById('respuesta').textContent = data.Short
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('respuesta').textContent = error;
    });
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