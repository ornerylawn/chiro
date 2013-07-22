# chiro

As in... chiropractor. As in... it helps you with your Backbone.

Install chiro:

```
$ go install github.com/ryanlbrown/chiro
```

Create a single-page app (complete with backbone, underscore, jquery, and requirejs):

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
