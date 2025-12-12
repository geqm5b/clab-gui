document.addEventListener('DOMContentLoaded', async () => {
    const gridContainer = document.getElementById('lab-grid');
    
    // Referencias para el Upload
    const fileInput = document.getElementById('file-upload');
    const fileNameSpan = document.getElementById('file-name');
    const uploadBtn = document.getElementById('btn-upload');

    // CARGA INICIAL 
    loadLabs();

    // MANEJO DE SELECCIÓN DE ARCHIVO 
    fileInput.addEventListener('change', () => {
        if (fileInput.files.length > 0) {
            fileNameSpan.textContent = fileInput.files[0].name;
            uploadBtn.disabled = false; // Habilitar botón subir
        } else {
            fileNameSpan.textContent = 'Ningún archivo seleccionado';
            uploadBtn.disabled = true;
        }
    });

    // LÓGICA SUBIR Y PROCESAR (UPLOAD) 
    uploadBtn.addEventListener('click', async () => {
        const file = fileInput.files[0];
        if (!file) return;

        const originalText = uploadBtn.textContent;
        uploadBtn.disabled = true;
        uploadBtn.textContent = 'Procesando...';

        const formData = new FormData();
        formData.append('file', file); // 'file' debe coincidir con backend en Go

        try {
            const res = await fetch('/api/upload', { 
                method: 'POST',
                body: formData
            });

            if (!res.ok) throw new Error(`HTTP ${res.status}`);

            // Éxito
            uploadBtn.textContent = '¡Conversión Exitosa!';
            uploadBtn.style.backgroundColor = '#28a745';
            
            // Limpiar input
            fileInput.value = '';
            fileNameSpan.textContent = 'Ningún archivo seleccionado';

            // RECARGAR LA LISTA DE LABS AUTOMÁTICAMENTE
            setTimeout(() => {
                uploadBtn.disabled = true; // Se deshabilita hasta seleccionar otro
                uploadBtn.textContent = originalText;
                uploadBtn.style.backgroundColor = '';
                loadLabs(); //  refrescar la grilla
            }, 2000);

        } catch (error) {
            console.error(error);
            uploadBtn.textContent = 'Error al subir';
            uploadBtn.style.backgroundColor = '#dc3545';
            setTimeout(() => {
                uploadBtn.disabled = false;
                uploadBtn.textContent = originalText;
                uploadBtn.style.backgroundColor = '';
            }, 3000);
        }
    });

    // DELEGACIÓN DE EVENTOS (BOTONES DEPLOY/DESTROY) ---
    gridContainer.addEventListener('click', (e) => {
        const btn = e.target.closest('.btn-action');
        if (btn) {
            const labName = btn.dataset.name;
            if (btn.classList.contains('btn-deploy')) {
                deployLab(labName, btn);
            } else if (btn.classList.contains('btn-destroy')) {
                destroyLab(labName, btn);
            }
        }
    });
});

// Función reutilizable para cargar la grilla
async function loadLabs() {
    const gridContainer = document.getElementById('lab-grid');
    try {
        const response = await fetch('/api/getLabs');
        if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);
        
        const data = await response.json();
        const labs = data.labs || [];

        if (labs.length > 0) {
            gridContainer.innerHTML = labs.map(lab => {
                const displayName = lab.name
                    .replace(/\.(clab\.yml|yml)$/, '')
                    .toUpperCase();

                return `
                <div class="lab-card">
                    <div class="lab-name">${displayName}</div>
                    
                    <button class="btn-action btn-deploy" data-name="${lab.name}">
                        Desplegar Lab
                    </button>
                    
                    <button class="btn-action btn-destroy" data-name="${lab.name}">
                        Terminar Lab
                    </button>
                </div>
                `;
            }).join('');
        } else {
            gridContainer.innerHTML = '<p style="text-align:center; grid-column: 1/-1; color: #666;">No se encontraron archivos .clab.yml</p>';
        }

    } catch (error) {
        console.error('Error cargando lista:', error);
        gridContainer.innerHTML = '<p style="color:#d73a49; text-align:center; grid-column: 1/-1;">Error de conexión con el servidor.</p>';
    }
}

async function deployLab(labName, btn) {
    const originalText = btn.textContent;
    btn.disabled = true;
    btn.textContent = 'Iniciando...';

    try {
        const res = await fetch('/api/deployLab', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: labName })
        });

        if (!res.ok) throw new Error(`HTTP ${res.status}`);

        btn.textContent = '¡Iniciado!';
        btn.style.backgroundColor = '#28a745'; 

        setTimeout(() => {
            btn.disabled = false;
            btn.textContent = originalText;
            btn.style.backgroundColor = ''; 
        }, 2000);

    } catch (error) {
        console.error(error);
        btn.textContent = 'Error';
        btn.style.backgroundColor = '#d73a49';
        setTimeout(() => {
            btn.disabled = false;
            btn.textContent = originalText;
            btn.style.backgroundColor = '';
        }, 3000);
    }
}

async function destroyLab(labName, btn) {
    const originalText = btn.textContent;
    btn.disabled = true;
    btn.textContent = 'Terminando...';
    btn.style.backgroundColor = '#d39e00'; 

    try {
        const res = await fetch('/api/destroyLab', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: labName })
        });

        if (!res.ok) throw new Error(`HTTP ${res.status}`);

        btn.textContent = '¡Terminado!';
        btn.style.backgroundColor = '#28a745'; 

        setTimeout(() => {
            btn.disabled = false;
            btn.textContent = originalText;
            btn.style.backgroundColor = ''; 
        }, 2000);

    } catch (error) {
        console.error(error);
        btn.textContent = 'Error';
        setTimeout(() => {
            btn.disabled = false;
            btn.textContent = originalText;
            btn.style.backgroundColor = '';
        }, 3000);
    }
}