const taskForm = document.getElementById("task-form");
const taskInput = document.getElementById("task-input");
const taskList = document.getElementById("task-list");

let tasks = [];

function escapeHTML(str) {
  const p = document.createElement("p");
  p.textContent = str;
  return p.innerHTML;
}

window.renderTasks = function () {
  taskList.innerHTML = "";

  tasks.forEach((task) => {
    const li = document.createElement("li");
    li.className = `task-item ${task.completed ? "completed" : ""}`;

    li.innerHTML = `
            <span>${escapeHTML(task.title)}</span>
            <div class="task-actions">
                <button class="btn-done" onclick="toggleTask(${
                  task.id
                })">✓</button>
                <button class="btn-delete" onclick="deleteTask(${
                  task.id
                })">✕</button>
            </div>
        `;
    taskList.appendChild(li);
  });
};

taskForm.addEventListener("submit", (e) => {
  e.preventDefault();
  const title = taskInput.value.trim();

  if (title) {
    const newTask = {
      id: Date.now(),
      title: title,
      completed: false,
    };

    tasks.push(newTask);
    renderTasks();
    taskInput.value = "";
  }
});

window.deleteTask = function (id) {
  tasks = tasks.filter((t) => t.id !== id);
  renderTasks();
};

window.toggleTask = function (id) {
  tasks = tasks.map((t) => {
    if (t.id === id) {
      return { ...t, completed: !t.completed };
    }
    return t;
  });
  renderTasks();
};

document.addEventListener("DOMContentLoaded", renderTasks);