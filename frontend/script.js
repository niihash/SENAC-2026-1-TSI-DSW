const taskForm = document.getElementById("task-form");
const taskInput = document.getElementById("task-input");
const taskList = document.getElementById("task-list");

const API_URL = "/api/v1/tasks";
let tasks = [];

function escapeHTML(str) {
  const p = document.createElement("p");
  p.textContent = str;
  return p.innerHTML;
}

async function loadTasks() {
  const response = await fetch(API_URL);
  if (!response.ok) {
    throw new Error("Erro ao carregar tarefas");
  }

  tasks = await response.json();
  renderTasks();
}

window.renderTasks = function () {
  taskList.innerHTML = "";

  tasks.forEach((task) => {
    const li = document.createElement("li");
    li.className = `task-item ${task.completed ? "completed" : ""}`;

    li.innerHTML = `
            <span>${escapeHTML(task.title)}</span>
            <div class="task-actions">
                <button class="btn-done" onclick="toggleTask(${task.id})">&check;</button>
                <button class="btn-delete" onclick="deleteTask(${task.id})">&times;</button>
            </div>
        `;
    taskList.appendChild(li);
  });
};

taskForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  const title = taskInput.value.trim();

  if (title) {
    const response = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title }),
    });

    if (!response.ok) {
      alert("Nao foi possivel adicionar a tarefa.");
      return;
    }

    const newTask = await response.json();
    tasks.push(newTask);
    renderTasks();
    taskInput.value = "";
  }
});

window.deleteTask = async function (id) {
  const response = await fetch(`${API_URL}/${id}`, { method: "DELETE" });
  if (!response.ok) {
    alert("Nao foi possivel excluir a tarefa.");
    return;
  }

  tasks = tasks.filter((t) => t.id !== id);
  renderTasks();
};

window.toggleTask = async function (id) {
  const task = tasks.find((t) => t.id === id);
  if (!task) {
    return;
  }

  const completed = !task.completed;
  const response = await fetch(`${API_URL}/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ completed }),
  });

  if (!response.ok) {
    alert("Nao foi possivel atualizar a tarefa.");
    return;
  }

  tasks = tasks.map((t) => {
    if (t.id === id) {
      return { ...t, completed };
    }
    return t;
  });
  renderTasks();
};

document.addEventListener("DOMContentLoaded", () => {
  loadTasks().catch(() => {
    alert("Nao foi possivel conectar ao backend.");
  });
});