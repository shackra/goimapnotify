Changelog
=========


(unreleased)
------------

Fix
~~~
- Fix misleading description on README.md. [Jorge Javier Araya Navarro]

  fix issue #3

Other
~~~~~
- Fix typo, add missing arg to README. [Maxim Baz]


2.0 (2019-04-27)
----------------

New
~~~
- Enable debug flag that shows network events. [Jorge Araya Navarro]

  Requirement of some users that need to debug network issues with their IMAP servers. The debugging
  output starts right after goimapnotify was able to establish a connection with the IMAP server but
  not before the user credentials are sent
- Updates code to use emersion's libraries. [Jorge Araya Navarro]

  the past library was unmaintained and old


1.1 (2019-01-22)
----------------

Changes
~~~~~~~
- Change glide for dep. [Jorge Araya Navarro]
- Update copyright date. [Jorge Araya Navarro]
- Make port in configuration mandatory. [Jorge Araya Navarro]
- Always try to enable STARTTLS. [Jorge Araya Navarro]

Fix
~~~
- Fix logical error in code. [Jorge Araya Navarro]

  Helps with the following error `[ERR] Cannot connect to imap.mail.yahoo.com:993: EOF`


1.0.1 (2017-08-31)
------------------

New
~~~
- Send the IDLE command again after 15 minutes. [Jorge Araya Navarro]

  This avoid the server closing the connection


1.0 (2017-08-26)
----------------

New
~~~
- Add GPL3+ license to the project. [Jorge Araya Navarro]
- Add read me file. [Jorge Araya Navarro]

  Explains important things about the application
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
