const taskForm = document.querySelector("#taskForm");
const taskTitle = document.querySelector("#taskTitle");
const taskDescription = document.querySelector("#taskDescription");
const tasksTableBody = document.querySelector("#tasksTableBody");
const message = document.querySelector("#message");
const refreshButton = document.querySelector("#refreshButton");

const statsTotal = document.querySelector("#statsTotal");
const statsTodo = document.querySelector("#statsTodo");
const statsInProgress = document.querySelector("#statsInProgress");
const statsDone = document.querySelector("#statsDone");

const statusLabels = {
  todo: "todo",
  in_progress: "in_progress",
  done: "done",
};

taskForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  clearMessage();

  const title = taskTitle.value.trim();
  const description = taskDescription.value.trim();
  if (!title) {
    showMessage("Введите название задачи");
    return;
  }

  try {
    await request("/api/tasks", {
      method: "POST",
      body: JSON.stringify({ title, description }),
    });
    taskForm.reset();
    showMessage("Задача создана", true);
    await loadPageData();
  } catch (error) {
    showMessage(error.message);
  }
});

refreshButton.addEventListener("click", () => {
  loadPageData();
});

tasksTableBody.addEventListener("click", async (event) => {
  const button = event.target.closest("button[data-action]");
  if (!button) {
    return;
  }

  clearMessage();
  const id = button.dataset.id;
  const action = button.dataset.action;

  try {
    if (action === "delete") {
      await request(`/api/tasks/${id}`, { method: "DELETE" });
      showMessage("Задача удалена", true);
    } else {
      await request(`/api/tasks/${id}/status`, {
        method: "PATCH",
        body: JSON.stringify({ status: action }),
      });
      showMessage("Статус обновлён", true);
    }
    await loadPageData();
  } catch (error) {
    showMessage(error.message);
  }
});

loadPageData();

async function loadPageData() {
  setLoading(true);
  clearMessage();

  try {
    const [tasks, stats] = await Promise.all([
      request("/api/tasks"),
      request("/api/stats"),
    ]);
    renderTasks(tasks);
    renderStats(stats);
  } catch (error) {
    showMessage(error.message);
  } finally {
    setLoading(false);
  }
}

async function request(url, options = {}) {
  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (response.status === 204) {
    return null;
  }

  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || "Ошибка запроса");
  }

  return data;
}

function renderTasks(tasks) {
  if (!tasks.length) {
    tasksTableBody.innerHTML = `<tr><td colspan="4" class="empty-state">Задач пока нет</td></tr>`;
    return;
  }

  tasksTableBody.innerHTML = tasks.map((task) => `
    <tr>
      <td>${escapeHtml(task.title)}</td>
      <td>${escapeHtml(task.description || "")}</td>
      <td><span class="status status-${task.status}">${statusLabels[task.status] || task.status}</span></td>
      <td>
        <div class="actions">
          ${statusButton(task.id, "in_progress", "В работу", task.status === "in_progress")}
          ${statusButton(task.id, "done", "Готово", task.status === "done")}
          <button class="delete-button" type="button" data-action="delete" data-id="${task.id}">Удалить</button>
        </div>
      </td>
    </tr>
  `).join("");
}

function statusButton(id, status, label, disabled) {
  return `<button type="button" data-action="${status}" data-id="${id}" ${disabled ? "disabled" : ""}>${label}</button>`;
}

function renderStats(stats) {
  statsTotal.textContent = stats.total ?? 0;
  statsTodo.textContent = stats.todo ?? 0;
  statsInProgress.textContent = stats.in_progress ?? 0;
  statsDone.textContent = stats.done ?? 0;
}

function setLoading(isLoading) {
  refreshButton.disabled = isLoading;
}

function showMessage(text, success = false) {
  message.textContent = text;
  message.classList.toggle("success", success);
}

function clearMessage() {
  showMessage("");
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}
