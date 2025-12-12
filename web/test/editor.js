document.addEventListener('DOMContentLoaded', () => {
    const iframe = document.getElementById('drawio-frame');
    const btnSave = document.getElementById('btn-save');
    const inputName = document.getElementById('filename');

    // CONFIGURACION DEL BOTON ---
    btnSave.addEventListener('click', () => {
        const nombre = inputName.value.trim();
        
        if (!nombre) {
            alert("Por favor, escribe un nombre para el archivo.");
            inputName.focus();
            return;
        }

        // Feedback visual
        btnSave.disabled = true;
        btnSave.textContent = "Guardando...";

        // Pedimos el XML al iframe de Draw.io
        iframe.contentWindow.postMessage(JSON.stringify({
            action: 'export', 
            format: 'xml', 
            spin: 'Generando XML...'
        }), '*');
    });

    // ESCUCHAR MENSAJES DE DRAW.IO 
    window.addEventListener('message', function(evt) {
        if (!evt.data || evt.data.length < 1) return;

        try {
            var msg = JSON.parse(evt.data);

            // A. INICIALIZACION
            if (msg.event === 'init') {
                iframe.contentWindow.postMessage(JSON.stringify({
                    action: 'load', autosave: 1, xml: ''
                }), '*');
            }

            // B. RECEPCION DEL XML
            if (msg.event === 'export') {
                enviarAlBackend(msg.data);
            }

        } catch(e) {}
    });

    // --- 3. ENVIO AL BACKEND ---
    function enviarAlBackend(xmlString) {
        const nombre = inputName.value.trim();

        fetch('/api/save_xml', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ 
                name: nombre,
                xml: xmlString 
            })
        })
        .then(r => r.json())
        .then(data => {
            if(data.error) {
                alert("Error al guardar: " + data.error);
            } else {
                alert("Topologia guardada correctamente.");
            }
        })
        .catch(e => alert("Error de red: " + e))
        .finally(() => {
            btnSave.disabled = false;
            btnSave.textContent = "Guardar Topologia XML";
        });
    }
});