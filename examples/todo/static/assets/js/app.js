var loadTodos = function() {
  $.ajax({
    url: "/graphql?query={todoList{id,text,done}}"
  }).done(function(data) {
    console.log(data);
    var dataParsed = JSON.parse(data);
    var todos = dataParsed.data.todoList;

    if (!todos.length) {
      $('.todo-list-container').append('<p>There are no tasks for you today</p>');
    }

    $.each(todos, function(i, v) {
      var doneHtml = '<input id="' + v.id + '" type="checkbox"' + (v.done ? ' checked="checked"' : '') + '>';      
      var labelHtml = '<label for="' + v.id + '">' + doneHtml + v.text + '</label>';
      var itemHtml = '<div class="todo-item">' + labelHtml + '</div>';
      
      $('.todo-list-container').append(itemHtml);
    });
  });
};

$(document).ready(function() {
  loadTodos();
});
