
# **z-run** -- lightweight Go-based scripts library tool




## About

(TO-BE-CONTINUED...)




## FAQ




### Why `z-run` vs classical scripts?

What sets `z-run` aside from a "classical" scripts is:

* first it allows one to easily mix multiple languages;
for example one can easily write a Bash script that chains together a Python, Ruby and another Bash script,
all within the same file;
(with classical scripts one would have had to create multiple files, one for each language, and explicitly invoke them via the interpreter;)

* it allows one to easily bundle together unrelated scripts;
for example say one uses `z-run` to aid in the development process, then one could bundle together scripts related
to compiling, testing, deploying, and other single-shot shortcuts in the same file;
(again, with classical scripts, one would have usually put each of these in separate files,
or by using `case` or `function` try to manage them inside a reduced set of files;)

* it allows one to easily invoke scriptlets from within the same library, without caring what language they are written in, or where their source code is stored;
(with classical scripts one has to take care of getting the paths to other scripts right, especially when running from a different folder than the original one;)

However `z-run` takes it a step further by offering the following functionalities:

* environment variables management -- by allowing certain environment variables to be overriden, removed or appended (like in the case of `$PATH`);
(with classical scripts one has to manage these by themselves, resort to various "envdir" tools, or interact with shell magic;)

* hierarchical menus -- by simply calling `z-run` without any arguments, it presents one with a selectable menu of all the scriptlets in the library,
and after one selects a scriptlet it executes it;  moreover one can group various scriptlets under other sub-menus, thus allowing arbitrary nesting and navigation;
(when used from a terminal, `z-run` uses and embedded `fzf`-based UI;  when used from outside a terminal, `z-run` tries to use `z-run--select` (any system), `rofi` or `dmenu` (under Linux and BSD's) or `choose` (under OSX);)

* SSH-based remote execution -- one can easily execute a given scriptlet on a remote server without having to previously copy anything there;
moreover once the scriptlet is executing on the remote server, it can invoke other scriptlets from the library that are also to be executed remotely;
(the only requirement is having `z-run` installed on the remote machine;)

* scriptlets generation -- one can easily write a scriptlet that generates `z-run` compliant source, thus generating other scriptlets based on arbitrary criteria;
(for example one could write a scriptlet, that generates other scriptlets specific for each file in a given folder;)

* library compiling -- one can easily create a single file that contains the entire library, which can then be moved and used in another place;

* Go-based templates -- that are useful especially in generating other scriptlets;




### Caching

`z-run` is designed to support thousands (and in extreme cases tens of thousands) of scriptlets, especially when scriptlet generation is used.

Therefore `z-run` has built-in optimizations to cache the library contents (without any regeneration) unless any of the following conditions are met:

* any of the files that comprise the source code of the library are changed;  (based both on file timestamps and contents hashing;)
* any of the environment variables change;
* a different version of `z-run` is installed;  (this covers both upgrading and downgrading;)

However once a scriptlet is executed, it and any other invoked scriptlets will use exactly the same cached library contents.
Thus it is safe to change the source code of the library, while a scriptlet is executing.




## Notice (copyright and licensing)


### Notice -- short version

The code is licensed under GPL 3 or later.


### Notice -- long version

For details about the copyright and licensing, please consult the `notice.txt` file in the `documentation/licensing` folder.

If someone requires the sources and/or documentation to be released
under a different license, please send an email to the authors,
stating the licensing requirements, accompanied with the reasons
and other details; then, depending on the situation, the authors might
release the sources and/or documentation under a different license.

