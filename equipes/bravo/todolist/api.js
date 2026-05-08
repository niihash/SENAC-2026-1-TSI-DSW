document.addEventListener('DOMContentLoaded', () => {
    carregarTarefas();
    const form = document.getElementById('form-tarefa');
    if (form) {
        form.addEventListener('submit', adicionarTarefa);
    }
});

// 1. ADICIONAR TAREFA (CREATE)
async function adicionarTarefa(event) {
    event.preventDefault();
    const input = document.getElementById('new-task');
    const nomeTarefa = input.value.trim();
    if (!nomeTarefa) return;

    const idTemporario = `temp-${Date.now()}`;

    // Renderiza com ID temporário
    renderizarTarefaNaTela({
        task_id: idTemporario,
        name: nomeTarefa,
        status: 'pending'
    });
    input.value = '';

    try {
        const response = await fetch('http://localhost:8080/api/v1/tasks', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: nomeTarefa, user_id: 1 })
        });

        const tarefaCriada = await response.json();
        const idReal = tarefaCriada.task_id;

        // Atualiza o elemento no DOM com o ID real
        const elemento = document.getElementById(`tarefa-${idTemporario}`);
        if (elemento && idReal) {
            elemento.id = `tarefa-${idReal}`;

            const checkbox = elemento.querySelector('input[type="checkbox"]');
            const botao = elemento.querySelector('button');

            if (checkbox) checkbox.setAttribute('onchange', `alternarStatus(${idReal}, this.checked)`);
            if (botao) botao.setAttribute('onclick', `deletarTarefa(${idReal})`);
        }

    } catch (error) {
        console.error("Erro ao salvar tarefa:", error);
    }
}

// 2. CARREGAR TAREFAS (READ)
async function carregarTarefas() {
    try {
        const response = await fetch('http://localhost:8080/api/v1/tasks');
        const tarefas = await response.json();
        const lista = document.getElementById('lista-tarefas');
        lista.innerHTML = '';

        if (Array.isArray(tarefas)) {
            tarefas.forEach(renderizarTarefaNaTela);
        }
    } catch (error) {
        console.error("Erro ao carregar tarefas:", error);
    }
}

// 3. FUNÇÃO QUE DESENHA NA TELA
function renderizarTarefaNaTela(tarefa) {
    const lista = document.getElementById('lista-tarefas');
    const li = document.createElement('li');
    li.id = `tarefa-${tarefa.task_id}`;
    li.className = "flex items-center justify-between bg-slate-50 border-2 border-slate-900 p-4 mb-2 shadow-[4px_4px_0px_0px_rgba(0,0,0,0.05)] transition-all";

    li.innerHTML = `
        <div class="flex items-center gap-3">
            <input type="checkbox" 
                   ${tarefa.status === 'completed' ? 'checked' : ''} 
                   onchange="alternarStatus(${tarefa.task_id}, this.checked)"
                   class="checkbox-minimal w-5 h-5 border-2 border-slate-900 cursor-pointer">
            <span class="task-text font-medium text-slate-700 ${tarefa.status === 'completed' ? 'line-through opacity-50' : ''}">
                ${tarefa.name}
            </span>
        </div>
        <button onclick="deletarTarefa(${tarefa.task_id})" class="text-red-500 hover:text-red-700 font-bold text-sm">
            EXCLUIR
        </button>
    `;
    lista.appendChild(li);
}

// 4. ATUALIZAR STATUS (UPDATE)
async function alternarStatus(id, concluida) {
    const status = concluida ? 'completed' : 'pending';
    const span = document.querySelector(`#tarefa-${id} .task-text`);

    if (!span) return;

    span.classList.toggle('line-through', concluida);
    span.classList.toggle('opacity-50', concluida);

    try {
        await fetch(`http://localhost:8080/api/v1/tasks/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ status })
        });
    } catch (error) {
        // Reverte o estado visual
        span.classList.toggle('line-through', !concluida);
        span.classList.toggle('opacity-50', !concluida);
        console.error("Erro ao atualizar status:", error);
    }
}

// 5. DELETAR TAREFA (DELETE)
async function deletarTarefa(id) {
    const elemento = document.getElementById(`tarefa-${id}`);
    if (elemento) elemento.remove();

    try {
        await fetch(`http://localhost:8080/api/v1/tasks/${id}`, {
            method: 'DELETE'
        });
    } catch (error) {
        console.error("Erro ao deletar tarefa:", error);
    }
}