const taskDiv = (task) => {
    return `<div class="task-card" data-task-id=${task.id} data-task-status="${task.status}" data-task-desc="${task.description}">
        <h3>${task.title}</h3>
        <p>Status: ${task.status}</p> <button data-task-id=${task.id} class="delete-task-btn">Delete</button>
    </div>`;
}

const getSections = (tasks) => {
    const taskSections = {};
    tasks.forEach(function(task) {
        if (!taskSections[task.status]) {
            taskSections[task.status] = [];
        }
        taskSections[task.status].push(task);
    });
    return taskSections
}

const populateTaskLists = (tasks) => {
    // Populate task lists
    const taskSections = getSections(tasks)
    let taskListHtml = '';
    for (let status in taskSections) {
        taskListHtml += `<h2>${status}</h2>`;
        taskListHtml += `<button id="create-new-task-${status.toLowerCase().replace(' ', '-')}" class="create-new-task-btn">Create Task</button>`;
        
        // task Cards
        taskSections[status].forEach(function(task) {
            taskListHtml += taskDiv(task)
        });
    }
    $('#task-list').html(taskListHtml);


    // Create task card
    $('.create-new-task-btn').click(function() {
        let status = $(this).prev('h2').text();
        const modal = $("#createTaskModal");
        $(".create-task-btn").text("Create new task");
        modal.show();
    });

    // Task div clicK 
    $('.task-card').click(function() {
        const taskId = $(this).data('task-id');
        const taskStatus = $(this).data('task-status');
        const taskDescription = $(this).data('task-description') ?? "not a rela value";
        const modal = $("#createTaskModal");
        $(".create-task-btn").text("Save");
        modal.show();
        $("#taskTitle").val($(this).find('h3').text());
        $("#taskDescription").text(taskDescription);
        $("#taskStatus").val(taskStatus);

        modal.data('task-id',taskId); // set the id of the task to the modal
    });


    // delete btn click
    $('.delete-task-btn').click(function(event) {
        event.stopPropagation(); // Prevent the event from bubbling up to the parent
        const taskId = $(this).data('task-id');
        $.ajax({
            url: `http://localhost:9080/api/tasks/${taskId}`,
            method: "DELETE",
            success: function () {
              getPopulateTaskLists();
            },
            error: function () {
              alert("Failed to delete task.");
            },
          });
    });
}

const getPopulateTaskLists = () => {
    $.ajax({
        url: 'http://localhost:9080/api/tasks/list',
        method: 'GET',
        success: function(tasks) {      
            populateTaskLists(tasks)
        },
        error: function(err) {
            console.error('Error fetching tasks:', err);
        }
    });
}


$(document).ready(function() {
    const modal = $("#createTaskModal");
    const closeButton = $(".close");

    getPopulateTaskLists();
  
    // Close modal when the close button is clicked
    closeButton.on("click", function () {
      modal.hide();
    });
  
    // Close modal when clicking outside the modal content
    $(window).on("click", function (event) {
      if ($(event.target).is(modal)) {
        modal.hide();
      }
    });
  
    // Handle create / update task form submission
    $("#createTaskForm").on("submit", function (event) {
      event.preventDefault();
      const title = $("#taskTitle").val();
      const description = $("#taskDescription").val();
      const status = $("#taskStatus").val();
      const taskId = modal.data('task-id');  
      $.ajax({
        url: "http://localhost:9080/api/tasks/create",
        method: "POST",
        contentType: "application/json",
        data: JSON.stringify({ title, description, status, id: taskId }),
        success: function () {
          getPopulateTaskLists();
          modal.hide();
        },
        error: function () {
          alert("Failed to create task.");
        },
      });
    });
});
