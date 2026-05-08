// --- COMMIT 3: Base do Frontend e Variáveis ---
const API_URL = "/api/v1";

window.onload = () => {
    if (localStorage.getItem("token")) {
        showTodoSection();
        loadTasks();
    }
}

// --- COMMIT 4: Funcionalidade de Autenticação (Login/Register) ---
async function login() {
    const user = document.getElementById("username").value;
    const pass = document.getElementById("password").value;

    const res = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: user, password: pass })
    });

    if (res.ok) {
        const data = await res.json();
        localStorage.setItem("token", data.token); // Salva o JWT no navegador
        showTodoSection();
        loadTasks();
    } else {
        alert("Credenciais inválidas!");
    }
}

async function register() {
    const user = document.getElementById("username").value;
    const pass = document.getElementById("password").value;

    const res = await fetch(`${API_URL}/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: user, password: pass })
    });

    if (res.ok) alert("Usuário criado! Agora faça login.");
}

function logout() {
    localStorage.removeItem("token");
    document.getElementById("auth-section").style.display = "block";
    document.getElementById("todo-section").style.display = "none";
}

function showTodoSection() {
    document.getElementById("auth-section").style.display = "none";
    document.getElementById("todo-section").style.display = "block";
}

// --- COMMIT 5: Funcionalidade de Listagem de Tarefas (Read) ---
async function loadTasks() {
    const res = await fetch(`${API_URL}/tasks`, {
        headers: { "Authorization": "Bearer " + localStorage.getItem("token") }
    });
    
    if (res.ok) {
        const tasks = await res.json();
        const list = document.getElementById("taskList");
        list.innerHTML = ""; 

        if(tasks){
            tasks.forEach(task => {
                const li = document.createElement("li");
                if (task.completed) li.classList.add("completed");

                const span = document.createElement("span");
                span.textContent = task.title;

                const check = document.createElement("input");
                check.type = "checkbox";
                check.checked = task.completed;
                check.onclick = () => toggleTask(task.id, task.title, !task.completed);

                const editBtn = document.createElement("button");
                editBtn.textContent = "Editar";
                editBtn.onclick = () => editTask(task.id, task.title, task.completed);

                const delBtn = document.createElement("button");
                delBtn.textContent = "Excluir";
                delBtn.onclick = () => deleteTask(task.id);

                li.appendChild(check);
                li.appendChild(span);
                li.appendChild(editBtn);
                li.appendChild(delBtn);
                list.appendChild(li);
            });
        }
    } else {
        logout();
    }
}

// --- COMMIT 6: Funcionalidade de Criação de Tarefas (Create) ---
async function addTask() {
    const title = document.getElementById("newTask").value;
    if (!title) return;

    await fetch(`${API_URL}/tasks`, {
        method: "POST",
        headers: { 
            "Content-Type": "application/json",
            "Authorization": "Bearer " + localStorage.getItem("token")
        },
        body: JSON.stringify({ title: title })
    });

    document.getElementById("newTask").value = "";
    loadTasks();
}

// --- COMMIT 7: Funcionalidade de Edição e Status (Update) ---
async function toggleTask(id, title, completed) {
    await fetch(`${API_URL}/tasks`, {
        method: "PUT",
        headers: { 
            "Content-Type": "application/json",
            "Authorization": "Bearer " + localStorage.getItem("token")
        },
        body: JSON.stringify({ id, title, completed })
    });
    loadTasks();
}

async function editTask(id, currentTitle, completed) {
    const newTitle = prompt("Novo título:", currentTitle);
    if (newTitle === null || newTitle === currentTitle) return;

    await fetch(`${API_URL}/tasks`, {
        method: "PUT",
        headers: { 
            "Content-Type": "application/json",
            "Authorization": "Bearer " + localStorage.getItem("token")
        },
        body: JSON.stringify({ id, title: newTitle, completed })
    });
    loadTasks();
}

// --- COMMIT 8: Funcionalidade de Exclusão (Delete) ---
async function deleteTask(id) {
    if (!confirm("Tem certeza que deseja excluir?")) return;

    await fetch(`${API_URL}/tasks?id=${id}`, {
        method: "DELETE",
        headers: { 
            "Authorization": "Bearer " + localStorage.getItem("token")
        }
    });
    loadTasks();
}