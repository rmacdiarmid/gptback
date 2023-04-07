document.addEventListener("DOMContentLoaded", function () {
    // Attach event listener to the form
    var form = document.getElementById("taskForm");
    if (form) {
      form.addEventListener("submit", submitTaskForm);
    }
  
    // Fetch the task list if the page contains the task list table
    if (document.getElementById("task-list")) {
      fetchTaskList();
    }
  });
  
  function submitTaskForm(event) {
    event.preventDefault();
  
    const title = document.getElementById("title").value;
    const description = document.getElementById("description").value;
    const task = { title: title, description: description };
  
    fetch("/tasks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(task),
    })
      .then((response) => {
        if (response.ok) {
          alert("Task created successfully!");
          document.getElementById("taskForm").reset();
          fetchTaskList(); // Refresh the task list
        } else {
          alert("Error creating task: " + response.statusText);
        }
      })
      .catch((error) => {
        alert("Error creating task: " + error);
      });
  }
  
  function fetchTaskList() {
    fetch("/tasks")
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error("Error fetching task list: " + response.statusText);
        }
      })
      .then((tasks) => {
        displayTaskList(tasks);
      })
      .catch((error) => {
        alert(error);
      });
  }
  
  function displayTaskList(tasks) {
    const taskList = document.getElementById("task-list");
    taskList.innerHTML = "";
  
    tasks.forEach((task) => {
      const row = document.createElement("tr");
  
      const idCell = document.createElement("td");
      idCell.textContent = task.id;
      row.appendChild(idCell);
  
      const titleCell = document.createElement("td");
      titleCell.textContent = task.title;
      row.appendChild(titleCell);
  
      const descriptionCell = document.createElement("td");
      descriptionCell.textContent = task.description;
      row.appendChild(descriptionCell);
  
      taskList.appendChild(row);
    });
  }
  