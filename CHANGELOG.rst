Changelog
=========


2.4-rc1
-------

Changes
~~~~~~~
- Refactor configuration-related code. [Jorge Javier Araya Navarro]

Other
~~~~~
- Merge branch 'shackra/refactor-config-loading' into 'master' [Jorge
  Javier Araya Navarro]

  chg: Refactor configuration-related code

  See merge request shackra/goimapnotify!13
- Add @ckardaris fix for concurrent map write error. [Jorge Javier Araya
  Navarro]
- Merge branch 'master' into 'master' [Jorge Javier Araya Navarro]

  Multiple hosts and scripts using single configuration file (addressing #16 and #7)

  See merge request shackra/goimapnotify!12
- Multiple hosts and scripts using single configuration file (addressing
  #16 and #7) [Charalampos Kardaris]


2.3.7 (2021-10-10)
------------------

Fix
~~~
- Remove syncing on start. [Jorge Javier Araya Navarro]

  - fixes shackra/goimapnotify#15


2.3.6 (2021-09-29)
------------------

Fix
~~~
- Add +build tag. [Jorge Javier Araya Navarro]

Other
~~~~~
- Merge branch '13-2-3-5-uses-go-build-without-a-build-comment' into
  'master' [Jorge Javier Araya Navarro]

  Resolve "2.3.5 uses go:build without a +build comment"

  Closes #13

  See merge request shackra/goimapnotify!11


2.3.5 (2021-09-28)
------------------

Fix
~~~
- Ensure that sync happens upon starting. [Jorge Javier Araya Navarro]


2.3.4 (2021-09-28)
------------------
- Merge branch '12-shell-command-is-not-executed-with-2-3-3-version'
  into 'master' [Jorge Javier Araya Navarro]

  Resolve "Shell command is not executed with 2.3.3 version"

  Closes #12

  See merge request shackra/goimapnotify!10
- Resolve "Shell command is not executed with 2.3.3 version" [Jorge
  Javier Araya Navarro]


2.3.3 (2021-09-26)
------------------

Changes
~~~~~~~
- Update years. [Jorge Javier Araya Navarro]
- Update readme. [Jorge Javier Araya Navarro]

Other
~~~~~
- Merge branch 'master' of gitlab.com:shackra/goimapnotify. [Jorge
  Javier Araya Navarro]
- Merge branch '10-fetcher-is-triggered-many-times-in-quick-succession-
  if-multiple-mails-arrive-at-the-same-time' into 'master' [Jorge Javier
  Araya Navarro]

  Resolve "Fetcher is triggered many times in quick succession if multiple mails arrive at the same time"

  Closes #10

  See merge request shackra/goimapnotify!9
- Resolve "Fetcher is triggered many times in quick succession if
  multiple mails arrive at the same time" [Jorge Javier Araya Navarro]


2.3.2 (2021-06-21)
------------------

Fix
~~~
- Resolve "xoauth2 not working as expected" [Jorge Javier Araya Navarro]

Other
~~~~~
- Merge branch '9-xoauth2-not-working-as-expected' into 'master' [Jorge
  Javier Araya Navarro]

  fix: Resolve "xoauth2 not working as expected"

  Closes #9

  See merge request shackra/goimapnotify!8


2.3.1 (2021-06-17)
------------------

Fix
~~~
- Resolve "Unable to specify 'sub-boxes'" [Jorge Javier Araya Navarro]

Other
~~~~~
- Merge branch '8-unable-to-specify-sub-boxes' into 'master' [Jorge
  Javier Araya Navarro]

  Resolve "Unable to specify 'sub-boxes'"

  Closes #8

  See merge request shackra/goimapnotify!7


2.3 (2021-05-25)
----------------

Fix
~~~
- Prevent multiple calls of `onNewMail` and `onNewMailPost` [Jorge
  Javier Araya Navarro]

Other
~~~~~
- Merge branch 'fix/shackra/goimapnotify#4' into 'master' [Jorge Javier
  Araya Navarro]

  fix: Prevent multiple calls of `onNewMail` and `onNewMailPost`

  Closes #4

  See merge request shackra/goimapnotify!4


2.2 (2021-04-19)
----------------

New
~~~
- Add support for xoauth2 authentication. [Jorge Javier Araya Navarro]

Other
~~~~~
- Merge branch 'add-xoauth2-support' into 'master' [Jorge Javier Araya
  Navarro]

  new: Add support for xoauth2 authentication

  See merge request shackra/goimapnotify!6


2.1.1 (2021-03-21)
------------------
- Merge branch 'add_systemd_unit' into 'master' [Jorge Javier Araya
  Navarro]

  Add systemd unit

  See merge request shackra/goimapnotify!5
- Add systemd unit. [Cyril Levis]
- Merge branch 'feat/moreCMD' into 'master' [Jorge Javier Araya Navarro]

  Be able to fetch username and host with a Cmd like passwordCmd

  See merge request shackra/goimapnotify!2
- Be able to fetch username and host with a Cmd like passwordCmd. [Cyril
  Levis]


2.1 (2021-03-19)
----------------

New
~~~
- Move to go.mod. [Jorge Javier Araya Navarro]

Fix
~~~
- Update Gitlab CI instructions. [Jorge Javier Araya Navarro]
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


