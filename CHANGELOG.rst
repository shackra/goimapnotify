Changelog
=========


(unreleased)
------------

New
~~~
- Add read me file. [Jorge Araya Navarro]

  Explains important things about the application
- Add Gitlab Pipelines integration. [Jorge Araya Navarro]

  Ensures the health of the code of the application
- Pass TLS options to secure Dial to server. [Jorge Araya Navarro]
- Handles TLS options from the configuration file. [Jorge Araya Navarro]
- List mailboxes and exit. [Jorge Araya Navarro]

  Gives a better panoram to the user regarding the hierarchy of his mailboxes

  http://busylog.net/telnet-imap-commands-note/

Changes
~~~~~~~
- Execute OnNewMailPost command. [Jorge Araya Navarro]

Fix
~~~
- Parse commands to execute them correctly. [Jorge Araya Navarro]

  Golang applications are not Unix shells
- Stop the application from hanging when close. [Jorge Araya Navarro]

  Avoid `kill`ing the application because the hang.

  http://www.tapirgames.com/blog/golang-channel-closing


