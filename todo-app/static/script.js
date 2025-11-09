async function fetchTodos() {
  console.log("Fetching todos from ", window.BACKEND_URL);
  const response = await fetch(window.BACKEND_URL);
  const todos = await response.json();
  const list = document.getElementById('todo-list');
  list.innerHTML = '';

  todos && todos.forEach(todo => {
    const li = document.createElement('li');
    li.textContent = todo.Task;
    list.appendChild(li);
  });
}

async function createTodo() {
  const input = document.getElementById('todo-input');
  const text = input.value.trim();
  if (!text) return;

  try {
    await fetch(window.BACKEND_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ task: text }),
  });

  input.value = '';
  fetchTodos();
  } catch (error) {
    console.error("error occured when creating new todo: ", error);
  }
}

document.getElementById('create-btn').addEventListener('click', createTodo);

fetchTodos(); // Load on page start
