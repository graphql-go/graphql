var updateTodo = function(id, isDone){
  $.ajax({
    url: '/graphql?query=mutation+_{updateTodo(id:"' + id + '",done:' + isDone + '){id,text,done}}'
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    var updatedTodo = dataParsed.data.updateTodo;
    if (updatedTodo.done) {
      $('#' + updatedTodo.id).parent().parent().addClass('todo-done');
    } else {
      $('#' + updatedTodo.id).parent().parent().removeClass('todo-done');
    } 
  });
};

var handleTodoList = function(object) {
  var todos = object;

  if (!todos.length) {
    $('.todo-list-container').append('<p>There are no tasks for you today</p>');
  }

  $.each(todos, function(i, v) {
    var doneHtml = '<input id="' + v.id + '" type="checkbox"' + (v.done ? ' checked="checked"' : '') + '>';      
    var labelHtml = '<label for="' + v.id + '">' + doneHtml + ' ' + v.text + '</label>';
    var itemHtml = '<div class="todo-item">' + labelHtml + '</div>';
    
    $('.todo-list-container').append(itemHtml);
    $('#' + v.id).click(function(){
      var id = $(this).prop('id');
      var isDone = $(this).prop('checked');
      updateTodo(id, isDone);
    });
  });
};

var loadTodos = function() {
  $.ajax({
    url: "/graphql?query={todoList{id,text,done}}"
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    handleTodoList(dataParsed.data.todoList);
  });
};

var addTodo = function(todoText) {
  if (!todoText || todoText === "") {
    alert('Please specify a task');
    return;
  }

  $.ajax({
    url: '/graphql?query=mutation+_{createTodo(text:"' + todoText + '"){id,text,done}}'
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    var todoList = [dataParsed.data.createTodo];
    handleTodoList(todoList);
  });
};

$(document).ready(function() {
  $('.todo-add-form').submit(function(e){
    e.preventDefault();
    addTodo($('.todo-add-form #task').val());
    $('.todo-add-form #task').val('');
  });

  loadTodos();
});
