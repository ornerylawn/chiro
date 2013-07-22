# chiro

As in... chiropractor. As in... it helps you with your Backbone.

Install chiro:

```
$ go install github.com/ryanlbrown/chiro
```

## What can I do with chiro?

Create a single-page app (complete with backbone, underscore, jquery, requirejs, and sass):

```
$ mkdir todolist
$ cd todolist
$ chiro init

Downloading template...
Created: css
Created: css/reset.css
Created: index.html
Created: js
Created: js/backbone.js
Created: js/base_model.js
Created: js/base_view.js
Created: js/jquery.js
Created: js/main.js
Created: js/main_view.js
Created: js/require.js
Created: js/text.js
Created: js/underscore.js
Created: sass
Created: sass/main.sass
Created: tmpl
Created: tmpl/main.html

```

Add models and views:

```
$ chiro add model Todo

Created: js/todo_model.js

$ chiro add view Todo

Created: js/todo_view.js
Created: sass/todo.sass
Created: tmpl/todo.html
Modified: index.html

$ chiro add view TodoList

Created: js/todo_list_view.js
Created: sass/todo_list.sass
Created: tmpl/todo_list.html
Modified: index.html

```

Add fonts from Google Web Fonts:

```
$ chiro add font Lato

Modified: index.html

```

There are commands for removing things too:

```
$ chiro

 Usage: chiro <command>

  command:

  init - create a new app

  add <subcommand>

    subcommand:

    view <view name> - create a new view with the given name
    model <model name> - create a new model with the given name
    font <font string> - add a font from Google Web Fonts

  remove <subcommand>

    subcommand:

    view <view name> - remove the view with the given name
    model <model name> - remove the model with the given name
    font <font string> - remove a font

```

## Structure

Let's start with CSS. Don't write CSS. Write Sass bro! This command
will take all of your Sass in `sass/` and write your CSS in `css/`.

```
$ sass --watch sass:css
```

These CSS files need to be included in `index.html` (apparently it is
too difficult for `require.js` to manage CSS correctly for all
browsers). `chiro` automatically does this when it creates the
Sass file for your view.

Let me just show you what `index.html` looks like after adding a `Todo` view:

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

As you can see, CSS (including fonts) is the only thing that
is managed in `index.html`. Everything else
is brought in by `require.js`, starting with `js/main.js`.

`js/main.js` looks like this:

```
require.config({

  paths: {
    'tmpl': '../tmpl',
  },

});

require([

  'jquery', 'main_view',

], function($, MainView) {

  'use strict';

  var mainView = new MainView();
  mainView.render();

  $(document).ready(function() {
    $('body').append(mainView.el);
  });

});
```

The `main` view is rendered right away, and appended to `body` when `document` is ready.

`js/main_view.js` looks like this:

```
define([

  'jquery', 'underscore', 'base_view', 'text!tmpl/main.html',

], function($, _, BaseView, tmplText) {

  'use strict';

  var MainView = BaseView.extend({

    tmpl: _.template(tmplText),
    tagName: 'div',
    className: 'main',

    render: function() {
      this.$el.html(this.tmpl())
      return this;
    },

  });

  return MainView;

});
```

It just renders its template inside of a `div`. Assuming you've already
run `chiro add view TodoList` you can add it the `main` view like so:

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

But we'll also need to make a place for it in the template (`tmpl/main.html`):

```
<h1>Todo List</h1>
<div class="content">
  <div class="list-container"></div>
</div>
```

The templates are underscore templates. The configuration for
it is in `js/base_view.js`, which is a good place to add any common view code.
