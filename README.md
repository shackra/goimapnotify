# Go IMAP notify [![pipeline status](https://gitlab.com/shackra/goimapnotify/badges/master/pipeline.svg)](https://gitlab.com/shackra/goimapnotify/commits/master) [![coverage report](https://gitlab.com/shackra/goimapnotify/badges/master/coverage.svg)](https://gitlab.com/shackra/goimapnotify/commits/master) [![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/K3K1XEZCQ) [![Support me on Patreon](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Fshieldsio-patreon.vercel.app%2Fapi%3Fusername%3Dshackra%26type%3Dpatrons&style=flat)](https://patreon.com/shackra)

Execute scripts on IMAP mailbox changes (new/deleted/updated messages) using IDLE, Golang version.

Please read the `CHANGELOG` file to know what's new.

🗊 **You can also check what's new at ko-fi, read in English or Spanish**: https://ko-fi.com/shackra/posts

This application is mostly compatible with the configuration of [imapnotify made with Python](https://github.com/a-sk/python-imapnotify) (be sure to change `password_eval` to `passwordCMD`, see [issue #3](https://gitlab.com/shackra/goimapnotify/issues/3)), the following are all options available for the configuration:

```yaml
configurations:
    -
        host: example.com
        port: 143
        tls: true
        tlsOptions:
            rejectUnauthorized: false
            starttls: true
        idleLogoutTimeout: 15
        username: USERNAME
        alias: ExampleCOM
        password: PASSWORD
        xoAuth2: false
        boxes:
            -
                mailbox: INBOX
                onNewMail: 'mbsync examplecom:INBOX'
                onNewMailPost: SKIP
    -
        hostCMD: COMMAND_TO_RETRIEVE_HOST
        port: 993
        tls: true
        tlsOptions:
            rejectUnauthorized: true
            starttls: true
        username: ''
        usernameCMD: ''
        password: ''
        passwordCMD: ''
        xoAuth2: false
        onNewMail: ''
        onNewMailPost: ''
        onDeletedMail: ''
        onDeletedMailPost: ''
        boxes:
            -
                mailbox: INBOX
                onNewMail: 'mbsync examplenet:INBOX'
                onNewMailPost: SKIP
            -
                mailbox: Junk
                onNewMail: 'mbsync examplenet:Junk'
                onNewMailPost: SKIP
```

On first start, the application will run `onNewMail` and `onNewMailPost` and then wait for events from your IMAP server.

- `onNewMail`: is an executable or script to run when new mail has arrived.
- `onNewMailPost`: is an executable or script to run after `onNewMail` has ran.
- `onDeletedMail`: is an executable or script to run when mail has been delete.
- `onDeletedMailPost`: is an executable or script to run after `onDeletedMail` has ran.
- `hostCMD`: is an executable or script that retrieves your host from somewhere, we cannot pass arguments to this command from `Stdin`.
- `usernameCMD`: is an executable or script that retrieves your username from somewhere, we cannot pass arguments to this command from `Stdin`.
- `passwordCMD`: is an executable or script that retrieves your password from somewhere, we cannot pass arguments to this command from `Stdin`.
- `xoAuth2`: is an option that allow us to login on your IMAP using OAuth2, **be aware**: the token is retrieve from `passwordCMD` (see shackra/goimapnotify#9).
- `wait`: is the delay in seconds before the mail syncing is trigger (see shackra/goimapnotify#10).
- `boxes`: List of mailboxes. If none is defined, all will be monitored.
- `idleLogoutTimeout`: Change the time between restarts of the IDLE command (see shackra/goimapnotify#49)
- `enableIDCommand`: Tell goimapotify that your server needs (and supports!) the ID command (see shackra/goimapnotify#58 shackra/goimapnotify#57; the servers in those tickets did not support ID and they responded with a non-standard error message, causing goimapnotify to fail)

The application will use TLS as long as the IMAP server advertises this capability. If you use self-signed certificates or something, be sure to set `rejectUnauthorized` as `false`.
To enable TLS connection, set `tls` as `true` and `starttls` as `false`

If your host do not offer IDLE, a sane default of checking every 15 minutes will take place instead.

You can also use xoAuth2 instead of password based authentication by setting the `xoAuth2` option to `true` and the output of a tool which can provide xoAuth2 encoded tokens in `passwordCMD`. Examples: [Google oauth2l](https://github.com/google/oauth2l) or [xoauth2 fetcher for O365](https://github.com/harishkrupo/oauth2ms).

## Install

    go install gitlab.com/shackra/goimapnotify@latest

## Usage

    Usage of goimapnotify:
      -conf string
        	Configuration file (default "${HOME}/.config/goimapnotify/goimapnotify.yaml")
      -list
        	List all mailboxes and exit
      -log-level string
        	change the logging level, possible values: error, warning/warn, info/information, debug (default "info")
      -wait int
        	Period in seconds between IDLE event and execution of scripts (default 1)

As you can notice, `-list` can help you figure out the mailbox hierarchy of your mail server.

# Development
nix-flake is use for development, is great and [you should try it](https://github.com/DeterminateSystems/nix-installer?tab=readme-ov-file#the-determinate-nix-installer) too! Activate support for flake in your nix installation and the environment will setup ✨*automagically*✨ for you.

## Generating and editing the CHANGELOG
When I started this project, I was naive and inexperienced with the fundamentals of software development, that has make most commits in this project have inconsistent titles that make it harder for tools like [`git-chglog`](https://github.com/git-chglog/git-chglog) help with CHANGELOG generation. I generated an ["old" CHANGELOG](./CHANGELOG_old.md) that contains all information until tag `2.3.13`. So, from now on, generate the CHANGELOG from tag `2.3.14` onwards, please and thank you!.
