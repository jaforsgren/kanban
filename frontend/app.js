const taskDivHtml = (task, borderColor) => {
    //<article class="card"></article>

    return `
    <div style="border-top: 5px solid ${borderColor};">
    <article class="task-card card" data-task-id=${task.id} data-board-id=${task.board_id} data-task-status="${task.status}" data-task-desc="${task.description}">
        <header>    
            <h3>${task.title}</h3>
        </header>
        <footer>
            <div class="flex two">
                <p>Status: ${task.status}</p> <button data-task-id=${task.id} data-board-id=${task.board_id} class="delete-task-btn">Delete</button>
            </div>
        </footer>
    </article>
    </div>
    `;
}

const createNewTaskDivHtml = () => {
    return `<div class="create-task-card">
        <h3>Create New Task</h3>
    </div>`;
}

const createNewBoardDivHtml = () => {
    return`
    <div class="create-board-card" >
    <article class="card create-board-link" >
        <header>    
            <p>New Board2</p>
        </header>
    </article>
    </div>
    `;

}

const BoardListHtml = (board) => {
    return `
    <div>
    <article class="card board-link" data-id="${board.id}" data-title="${board.title}">
        <header>    
            <h4>${board.title}</h4>
        </header>
     
    </article>
    </div>
    `;
}

const getSections = (tasks) => {
    const taskSections = {};
    tasks.forEach(function(task) {
        if (!taskSections[task.status]) {
            taskSections[task.status] = [];
        }
        taskSections[task.status].push(task);
    });
    
    if (!Object.keys(taskSections).length) {
        taskSections['To Do'] = []
    } 
    return taskSections
}

const populateTaskLists = (tasks) => {
    // Populate task lists
    const taskSections = getSections(tasks)
    let taskListHtml = '';

    // Ugly styling

    const startHue = -80;
    let hueShift = 80;

    for (let status in taskSections) {

        const borderColor = `hsl(${(startHue + hueShift) % 360}, 100%, 50%)`;
        hueShift += hueShift;
        
        taskListHtml += '<div>'
        taskListHtml += `<h2>${status}</h2>`;        
        // task Cards
        taskSections[status].forEach(function(task) {
            taskListHtml += taskDivHtml(task, borderColor)
        });
        taskListHtml += createNewTaskDivHtml()
        taskListHtml += '</div>'
    }

    $('#tasks-list').html(taskListHtml);

    $('.create-task-card').click(function(event) {
        event.stopPropagation(); // Prevent the event from bubbling up to the parent
        const modal = $("#createTaskModal");
        $(".create-task-btn").text("Create new task");
        modal.removeData('task-id');
        
        modal.show();
    });

    // Task div clicK 
    $('.task-card').click(function() {
        const taskId = $(this).data('task-id');
        const taskStatus = $(this).data('task-status');

        const taskDescription = $(this).data('task-desc') ?? "";

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
        const boardId = $(this).data('boardId');
       
        $.ajax({
            url: `http://localhost:9080/api/tasks/${taskId}`,
            method: "DELETE",
            success: function () {
                loadBoardDetails(boardId);
            },
            error: function () {
              alert("Failed to delete task.");
            },
          });
    });
}

function loadAllBoards() {
    $.get('http://localhost:9080/api/boards', function(data) {
        data.forEach(board => {
            $('#boards-list').append(BoardListHtml(board));
        });

        $('#boards-list').append(createNewBoardDivHtml())
    });
}

function loadBoardDetails(boardId, title) {
    $.get(`http://localhost:9080/api/boards/${boardId}`, function(data) {
        if (title) {
            $('#board-title').text(`${title} (${boardId})`);
        }
        $('#tasks-list').empty();
        populateTaskLists(data)
        $('#all-boards-page').hide();
        $('#board-page').show();
        const modal = $("#createTaskModal");
        modal.data('boardId',boardId)
    });
}


$(document).ready(function() {
    const modal = $("#createTaskModal");
    const createBoardModal = $("#createBoardModal");
    const closeButton = $(".close");
  
    // Close modal when the close button is clicked
    closeButton.on("click", function () {
      modal.hide();
      createBoardModal.hide()
    });
  
    // Close modal when clicking outside the modal content
    $(window).on("click", function (event) {
      if ($(event.target).is(modal)) {
        modal.hide();
        createBoardModal.hide()
      }
    });

     // Load all boards by default
     loadAllBoards();

     // Navigation event
     $(document).on('click', '.board-link', function(e) {
         e.preventDefault();
         const boardId = $(this).data('id');
         const boardTitle = $(this).data('title');
         loadBoardDetails(boardId, boardTitle);
     });

     $(document).on('click', '.create-board-card', function(event) {
        event.stopPropagation(); // Prevent the event from bubbling up to the parent
        const modal = $("#createBoardModal");
        $(".create-board-btn").text("Create new board");
        modal.show();
    });

     $("#createBoardForm").on("submit", function (event) {
        const title = $("#boardTitle").val();
        const createBoardModal = $("#createBoardModal");
        event.preventDefault();
        $.ajax({
            url: "http://localhost:9080/api/boards",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({ title, user_id:"06242ebf-5a47-48fb-982b-13982568e845" }),
            success: function (data) {
              loadBoardDetails(JSON.parse(data).id, title);
              createBoardModal.hide();
            },
            error: function () {
              alert("Failed to create task.");
            },
          });
    });

  
    // Handle create / update task form submission
    $("#createTaskForm").on("submit", function (event) {
      event.preventDefault();
      const title = $("#taskTitle").val();
      const description = $("#taskDescription").val();
      const status = $("#taskStatus").val();
      const taskId = modal.data('task-id');
      const boardId = modal.data('boardId');

      let method = 'POST'
      if (taskId) {
        method = 'PATCH'
      }

      $.ajax({
        url: "http://localhost:9080/api/tasks",
        method,
        contentType: "application/json",
        data: JSON.stringify({ title, description, status, id: taskId, board_id: boardId}),
        success: function () {
          loadBoardDetails(boardId);
          modal.hide();
        },
        error: function () {
          alert("Failed to create task.");
        },
      });
    });
});
