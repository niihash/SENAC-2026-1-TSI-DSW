//auth
if (!SessionManager.getUserId()) {
    window.location.replace('login.html');
}

const taskInput = document.getElementById('task-input');
const addBtn = document.getElementById('add-btn');
const taskList = document.getElementById('task-list');

let allTasks = [];

// Carrega as tarefas do banco
window.loadTasks = async function() {
    const userId = SessionManager.getUserId();

    if (!userId) {
        window.location.href = 'login.html';
        return; // para aqui e não executa mais nada
    }

    const email = SessionManager.getUserEmail();
    const titulo = document.querySelector('h1');
    if (titulo && email) {
        titulo.innerText = `Tarefas de: ${email}`;
    }

   try {
        const data = await api.getTasks();
        if (Array.isArray(data)) {
            allTasks = data;
        } else {
            allTasks = [];
        }
        render();
    } catch (error) {
        console.error("Erro ao carregar:", error);
    }
};

// Adiciona nova tarefa
async function adicionarTarefa() {
    const title = taskInput.value.trim();
    if (!title) return;
    
    try {
        const savedTask = await api.createTask(title);
        // Garante que pega a tarefa retornada pelo back ou recarrega a lista
        if (savedTask && (savedTask.id || savedTask.ID)) {
            allTasks.push(savedTask);
            taskInput.value = '';
            render();
        } else {
            window.loadTasks();
            taskInput.value = '';
        }
    } catch (error) {
        alert('Erro ao adicionar: ' + error.message);
    }
}

// Inverte o status
window.alternarStatus = async function(taskId) {
    try {
        const index = allTasks.findIndex(t => (t.id || t.ID) === taskId);
        if (index > -1) {
            const statusAtual = allTasks[index].done || allTasks[index].Done;
            const novoStatus = !statusAtual;
            
            // Atualiza a memória local
            allTasks[index].done = novoStatus;
            allTasks[index].Done = novoStatus;
            
            // Envia para o banco
            await api.updateTask(allTasks[index]);
            render();
        }
    } catch (error) {
        alert('Erro ao atualizar no banco.');
        window.loadTasks(); // Se der erro no back, desfaz na tela
    }
};

// Exclui a tarefa
window.excluirTarefa = async function(taskId) {
    if (!confirm('Deseja excluir esta tarefa?')) return;
    try {
        await api.deleteTask(taskId);
        allTasks = allTasks.filter(t => (t.id || t.ID) !== taskId);
        render();
    } catch (error) {
        alert('Erro ao excluir: ' + error.message);
    }
};

// Renderiza na tela SEGUINDO SEU CSS
function render() {
    if (!taskList) return;
    taskList.innerHTML = '';
    
    allTasks.forEach(t => {
        const taskId = t.id || t.ID;
        const isDone = t.done || t.Done;
        const title = t.title || t.Title;
        
        // <li> principal
        const li = document.createElement('li');
        if (isDone) li.classList.add('completed'); 

        // <span> do texto da tarefa
        const span = document.createElement('span');
        span.className = 'task-text'; 
        span.textContent = title;

        // <div> das ações
        const divActions = document.createElement('div');
        divActions.className = 'actions'; 

        // Botão Concluir
        const btnDone = document.createElement('button');
        btnDone.className = 'btn-done'; 
        btnDone.textContent = '✓';
        btnDone.onclick = () => window.alternarStatus(taskId);

        // Botão Excluir
        const btnDelete = document.createElement('button');
        btnDelete.className = 'btn-delete'; 
        btnDelete.textContent = 'X';
        btnDelete.onclick = () => window.excluirTarefa(taskId);

        divActions.append(btnDone, btnDelete);
        li.append(span, divActions);
        taskList.appendChild(li);
    });
}

// Eventos de clique e enter
if (addBtn) addBtn.onclick = adicionarTarefa;
if (taskInput) {
    taskInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') adicionarTarefa();
    });
}

// Inicia as tarefas
if (taskList) window.loadTasks();