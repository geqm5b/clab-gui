// Espera a que la página (index.html) se cargue completamente
document.addEventListener('DOMContentLoaded', () => {
    
    // 1. Busca el 'ul' (lista) que creamos en el HTML
    const listaLabs = document.getElementById('lista-labs');

    // 2. Llama a la API de Backend 
    fetch('/api/labs')
        .then(response => response.json()) // Convierte la respuesta a JSON
        .then(data => {
            console.log(data)
            // Limpiamos el "Cargando..."
            listaLabs.innerHTML = ''; 
            // 3. Verifica si la lista "labs" tiene algo
            if (data.labs && data.labs.length > 0) {
                
                // 4. Recorre la lista y crea un <li> por cada lab
                data.labs.forEach(lab => {
                    const item = document.createElement('button');
                    item.textContent = lab.name; // lab.name (viene del JSON)
                    listaLabs.appendChild(item);
                });

            } else {
                // Como dijiste: "que no haya archivos no es un fallo"
                listaLabs.innerHTML = '<li>No se encontraron laboratorios</li>';
            }
        })
        .catch(error => {
            // 5. Maneja errores (si la API se cae o hay un error 500)
            console.error('Error al llamar a la API:', error);
            listaLabs.innerHTML = '<li>Error al cargar los laboratorios.</li>';
        });
});