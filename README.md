# Go IMAP notify [![pipeline status](https://gitlab.com/shackra/goimapnotify/badges/master/pipeline.svg)](https://gitlab.com/shackra/goimapnotify/commits/master) [![coverage report](https://gitlab.com/shackra/goimapnotify/badges/master/coverage.svg)](https://gitlab.com/shackra/goimapnotify/commits/master)

Execute scripts on IMAP mailbox changes (new/deleted/updated messages) using IDLE, golang version.

Please read the `CHANGELOG` file to know what's new.

This application is mostly compatible with the configuration of [imapnotify made with Python](https://github.com/a-sk/python-imapnotify) (be sure to change `password_eval` to `passwordCmd`, see [issue #3](https://gitlab.com/shackra/goimapnotify/issues/3)), the following are all options available for the configuration:

    {
      "host": "",
      "hostCmd": "",
      "port": 143,
      "tls": false,
      "tlsOptions": {
        "rejectUnauthorized": true
      },
      "username": "",
      "usernameCmd": "",
      "password": "",
      "passwordCmd": "",
      "xoauth2": false,
      "onNewMail": "",
      "onNewMailPost": "",
      "wait": 20,
      "boxes": [
        "INBOX"
      ]
    }

On first start, the application will run `onNewMail` and `onNewMailPost` and then wait for events from your IMAP server.

- `onNewMail`: is an executable or script to run when new mail has arrived.
- `onNewMailPost`: is an executable or script to run after `onNewMail` has ran.
- `hostCmd`: is an executable or script that retrieves your host from somewhere, we cannot pass arguments to this command from `Stdin`.
- `usernameCmd`: is an executable or script that retrieves your username from somewhere, we cannot pass arguments to this command from `Stdin`.
- `passwordCmd`: is an executable or script that retrieves your password from somewhere, we cannot pass arguments to this command from `Stdin`.
- `xoauth2`: is an option that allow us to login on your IMAP using OAuth2, **be aware**: the token is retrieve from `passwordCmd` (see shackra/goimapnotify#9).
- `wait`: is the delay in seconds before the mail syncing is trigger (see shackra/goimapnotify#10).

The application will use TLS as long as the IMAP server advertises this capability. If you use self-signed certificates or something, be sure to set `rejectUnauthorized` as `false`.

If your host do not offer IDLE, a sane default of checking every 15 minutes will take place instead.

You can also use xoauth2 instead of password based authentication by setting the `xoauth2` option to `true` and the output of a tool which can provide xoauth2 encoded tokens in `passwordCmd`. Examples: [Google oauth2l](https://github.com/google/oauth2l) or [xoauth2 fetcher for O365](https://github.com/harishkrupo/oauth2ms).

## Install

    go get -u gitlab.com/shackra/goimapnotify

## Usage

    Usage of goimapnotify:
    -conf string
        Configuration file (default "path/to/imapnotify.conf")
    -debug
        Output all network activity to the terminal (!! this may leak passwords !!)
    -list
        List all mailboxes and exit

As you can notice, `-list` can help you figure out the mailbox hierarchy of your mail server.
