document.addEventListener('DOMContentLoaded', async () => {
    const gridContainer = document.getElementById('lab-grid');

    // --- 1. DELEGACIÓN DE EVENTOS CENTRALIZADA ---
    gridContainer.addEventListener('click', (e) => {
        // Detectamos si se hizo clic en CUALQUIER botón de acción
        const btn = e.target.closest('.btn-action');
        
        if (btn) {
            const labName = btn.dataset.name; // Nombre real (ej: dhcp.clab.yml)

            // Identificamos QUÉ botón es y llamamos a la función correspondiente
            if (btn.classList.contains('btn-deploy')) {
                deployLab(labName, btn);
            } else if (btn.classList.contains('btn-destroy')) {
                destroyLab(labName, btn);
            }
        }
    });

    // --- 2. CARGA INICIAL (GET) ---
    try {
        const response = await fetch('/api/getLabs');
        
        if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);
        
        const data = await response.json();
        const labs = data.labs || [];

        if (labs.length > 0) {
            // Renderizamos las tarjetas
            gridContainer.innerHTML = labs.map(lab => {
                // Lógica visual: Quitar extensión y Mayúsculas
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
});


// --- 3. LÓGICA: DESPLEGAR (Deploy) ---
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

        // Feedback Visual Éxito
        btn.textContent = '¡Iniciado!';
        btn.style.backgroundColor = '#28a745'; 

        // Restaurar estado
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


// --- 4. LÓGICA: DESTRUIR (Destroy) ---
async function destroyLab(labName, btn) {
    const originalText = btn.textContent;
    btn.disabled = true;
    btn.textContent = 'Terminando...';

    // Feedback visual de proceso (Naranja oscuro)
    btn.style.backgroundColor = '#d39e00'; 

    try {
        const res = await fetch('/api/destroyLab', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: labName })
        });

        if (!res.ok) throw new Error(`HTTP ${res.status}`);

        // Feedback Visual Éxito
        btn.textContent = '¡Terminado!';
        btn.style.backgroundColor = '#28a745'; // Verde

        // Restaurar estado
        setTimeout(() => {
            btn.disabled = false;
            btn.textContent = originalText;
            btn.style.backgroundColor = ''; // Vuelve a rojo (CSS)
        }, 2000);

    } catch (error) {
        console.error(error);
        btn.textContent = 'Error';
        // Mantenemos el rojo de error un poco más
        setTimeout(() => {
            btn.disabled = false;
            btn.textContent = originalText;
            btn.style.backgroundColor = '';
        }, 3000);
    }
}