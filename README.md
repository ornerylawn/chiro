# chiro

Install `chiro`:

```
$ go get github.com/ryanlbrown/chiro
```

## What can I do with `chiro`?

Create a single-page app (complete with backbone, underscore, jquery, requirejs, and sass):

```
$ mkdir todolist
$ cd todolist
$ chiro init
```

Add models and views:

```
$ chiro add model Todo
$ chiro add view Todo
$ chiro add view TodoList
```

Add fonts from Google Web Fonts:

```
$ chiro add font Lato
```

And you can remove all of those things too.

## Structure

The structure is very simple. There's an `index.html` and there are folders that contain exactly what they say:

* __js__ - javascript files
* __tmpl__ - underscore templates
* __sass__ - sass styles
* __css__ - css files generated from the sass files (and reset.css)

When you add a model, say Todo, `chiro` creates `js/todo_model.js`. It looks like this:

```
define([

  'underscore', 'base_model',

], function(_, BaseModel) {

  'use strict';

  var TodoModel = BaseModel.extend({

    defaults: {

    },
    
    initialize: function(options) {
    
    },

  });

  return TodoModel;

});
```

When you add a view, say Todo, `chiro` creates `sass/todo.js` and `tmpl/todo.js`, which are both empty, and `js/todo_view.js`. It looks like this:

```
define([

  'jquery', 'underscore', 'base_view', 'text!tmpl/todo.html',

], function($, _, BaseView, tmplText) {

  'use strict';

  var TodoView = BaseView.extend({

    tmpl: _.template(tmplText),
    tagName: 'div',
    className: 'todo',

    initialize: function(options) {

    },

    render: function() {
      this.$el.html(this.tmpl());
      return this;
    },

  });

  return TodoView;

});
```

`chiro` also adds a link tag in `index.html` that points to `css/todo.css`, which will be created automatically for you if you run this command:

```
$ sass --watch sass:css
```

You'll end up with an `index.html` that looks like this:

```
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>spapp</title>
    <link rel="stylesheet" href="http://fonts.googleapis.com/css?family=Open+Sans" />
    <link rel="stylesheet" href="css/reset.css" />
    <link rel="stylesheet" href="css/main.css" />
    <link rel="stylesheet" href="css/todo.css" />
    <script data-main="js/main" src="js/require.js"></script>
  </head>
  <body>
  </body>
</html>
```

As you can see, CSS is managed in `index.html`. Everything else
is brought in by `require.js`, starting with `js/main.js`.

`js/main.js` just tells `require.js` where the `tmpl` directory
is (so that the views can refer to it) and then renders the main
view. It adds the main view to the DOM once the `document` is `ready`.

```
require.config({ paths: {'tmpl': '../tmpl'} });

require(['jquery', 'main_view'], function($, MainView) {

  'use strict';

  var mainView = new MainView();
  mainView.render();

  $(document).ready(function() {
    $('body').append(mainView.el);
  });

});
```

Let's say you've just created a new view, TodoList, and you want to add it to the main view.
All you need to do is create a spot for it in the main view's template (`tmpl/main.html`):

```
<h1>Todo List</h1>
<div class="content">
  <div class="list-container"></div>
</div>
```

And then make sure that the main view renders a `TodoListView` into the element:

```
define([

  'jquery', 'underscore', 'base_view', 'text!tmpl/main.html', `todo_list_view`

], function($, _, BaseView, tmplText, TodoListView) {

  'use strict';

  var MainView = BaseView.extend({

    tmpl: _.template(tmplText),
    tagName: 'div',
    className: 'main',

    render: function() {
      this.$el.html(this.tmpl())
      
      var todoListView = new TodoListView();
      this.$('.list-container').html(todoListView.render().el);
      
      return this;
    },

  });

  return MainView;

});
```

## Todo

* Add a router.
* Add a `build` command that will prepare the app for release.
* Allow `init`ing from a custom template.
