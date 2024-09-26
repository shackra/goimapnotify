<a name="unreleased"></a>
## [Unreleased]


<a name="2.3.13"></a>
## [2.3.13] - 2024-04-27
### Merge Requests
- Merge branch 'fix-issue-24' into '2.3.x'
- Merge branch 'documentation-fix' into '2.3.x'


<a name="2.3.12"></a>
## [2.3.12] - 2024-03-05

<a name="2.3.11"></a>
## [2.3.11] - 2024-01-02
### Merge Requests
- Merge branch '2.3.x' into '2.3.x'


<a name="2.3.10"></a>
## [2.3.10] - 2023-12-17
### Merge Requests
- Merge branch 'implement-starttls' into '2.3.x'


<a name="2.3.9"></a>
## [2.3.9] - 2023-10-21

<a name="2.3.8"></a>
## [2.3.8] - 2023-10-21

<a name="2.4-rc4"></a>
## [2.4-rc4] - 2023-05-07
### Changes
- Update to Go 1.20
- Remove Terra Luna donation address
- Advertise donations through some stablecoins

### Fixes
- Display all errors when trying to load the configuration


<a name="2.4-rc3"></a>
## [2.4-rc3] - 2022-04-03
### Merge Requests
- Merge branch '25-bye-session-invalidated-accesstokenexpired' into 'master'


<a name="2.4-rc2"></a>
## [2.4-rc2] - 2022-03-02
### Merge Requests
- Merge branch '20-support-starting-when-network-is-not-available' into 'master'
- Merge branch 'list' into 'master'


<a name="2.4-rc1"></a>
## [2.4-rc1] - 2021-12-02
### Changes
- Refactor configuration-related code

### Merge Requests
- Merge branch 'shackra/refactor-config-loading' into 'master'
- Merge branch 'master' into 'master'


<a name="2.3.7"></a>
## [2.3.7] - 2021-10-10
### Changes
- Update CHANGELOG !minor

### Fixes
- Remove syncing on start


<a name="2.3.6"></a>
## [2.3.6] - 2021-09-29
### Fixes
- Add +build tag

### Merge Requests
- Merge branch '13-2-3-5-uses-go-build-without-a-build-comment' into 'master'


<a name="2.3.5"></a>
## [2.3.5] - 2021-09-27
### Fixes
- Ensure that sync happens upon starting


<a name="2.3.4"></a>
## [2.3.4] - 2021-09-27
### Changes
- Update changelog !minor

### Merge Requests
- Merge branch '12-shell-command-is-not-executed-with-2-3-3-version' into 'master'


<a name="2.3.3"></a>
## [2.3.3] - 2021-09-26
### Changes
- Update years
- Update readme

### Merge Requests
- Merge branch '10-fetcher-is-triggered-many-times-in-quick-succession-if-multiple-mails-arrive-at-the-same-time' into 'master'


<a name="2.3.2"></a>
## [2.3.2] - 2021-06-21
### Fixes
- Resolve "xoauth2 not working as expected"

### Merge Requests
- Merge branch '9-xoauth2-not-working-as-expected' into 'master'


<a name="2.3.1"></a>
## [2.3.1] - 2021-06-17
### Fixes
- Resolve "Unable to specify 'sub-boxes'"

### Merge Requests
- Merge branch '8-unable-to-specify-sub-boxes' into 'master'


<a name="2.3"></a>
## [2.3] - 2021-05-25
### Changes
- Update changelog !minor

### Fixes
- Prevent multiple calls of `onNewMail` and `onNewMailPost`

### Merge Requests
- Merge branch 'fix/shackra/goimapnotify[#4](https://gitlab.com/shackra/goimapnotify/issues/4)' into 'master'


<a name="2.2"></a>
## [2.2] - 2021-04-19
### Features
- Add support for xoauth2 authentication

### Merge Requests
- Merge branch 'add-xoauth2-support' into 'master'


<a name="2.1.1"></a>
## [2.1.1] - 2021-03-21
### Merge Requests
- Merge branch 'add_systemd_unit' into 'master'
- Merge branch 'feat/moreCMD' into 'master'


<a name="2.1"></a>
## [2.1] - 2021-03-19
### Changes
- Update CHANGELOG !minor
- Update README.md !minor

### Features
- Move to go.mod

### Fixes
- Update Gitlab CI instructions
- Fix misleading description on README.md


<a name="2.0"></a>
## [2.0] - 2019-04-26
### Changes
- Update CHANGELOG !minor
- Small changes !minor

### Features
- Enable debug flag that shows network events
- Updates code to use emersion's libraries

### Fixes
- Update the usage section in the "Read me" file !minor


<a name="1.1"></a>
## [1.1] - 2019-01-22
### Changes
- Change glide for dep
- Update copyright date
- Make port in configuration mandatory
- Always try to enable STARTTLS

### Features
- Say "Bye!" when all goroutines has finished !minor

### Fixes
- Fix logical error in code


<a name="1.0.1"></a>
## [1.0.1] - 2017-08-30
### Features
- Send the IDLE command again after 15 minutes


<a name="1.0"></a>
## 1.0 - 2017-08-25
### Changes
- Move the IMAP client creation to a function !refactor
- Execute OnNewMailPost command

### Features
- Add GPL3+ license to the project
- Add read me file
- Add read me file
- Add Gitlab Pipelines integration
- Makes the post-commit script executable !minor
- Pass TLS options to secure Dial to server
- Add post commit hook !minor
- Handles TLS options from the configuration file
- List mailboxes and exit
- First commit

### Fixes
- Update canonical URL for the application !minor
- Parse commands to execute them correctly
- Stop the application from hanging when close


[Unreleased]: https://gitlab.com/shackra/goimapnotify/compare/2.3.13...HEAD
[2.3.13]: https://gitlab.com/shackra/goimapnotify/compare/2.3.12...2.3.13
[2.3.12]: https://gitlab.com/shackra/goimapnotify/compare/2.3.11...2.3.12
[2.3.11]: https://gitlab.com/shackra/goimapnotify/compare/2.3.10...2.3.11
[2.3.10]: https://gitlab.com/shackra/goimapnotify/compare/2.3.9...2.3.10
[2.3.9]: https://gitlab.com/shackra/goimapnotify/compare/2.3.8...2.3.9
[2.3.8]: https://gitlab.com/shackra/goimapnotify/compare/2.4-rc4...2.3.8
[2.4-rc4]: https://gitlab.com/shackra/goimapnotify/compare/2.4-rc3...2.4-rc4
[2.4-rc3]: https://gitlab.com/shackra/goimapnotify/compare/2.4-rc2...2.4-rc3
[2.4-rc2]: https://gitlab.com/shackra/goimapnotify/compare/2.4-rc1...2.4-rc2
[2.4-rc1]: https://gitlab.com/shackra/goimapnotify/compare/2.3.7...2.4-rc1
[2.3.7]: https://gitlab.com/shackra/goimapnotify/compare/2.3.6...2.3.7
[2.3.6]: https://gitlab.com/shackra/goimapnotify/compare/2.3.5...2.3.6
[2.3.5]: https://gitlab.com/shackra/goimapnotify/compare/2.3.4...2.3.5
[2.3.4]: https://gitlab.com/shackra/goimapnotify/compare/2.3.3...2.3.4
[2.3.3]: https://gitlab.com/shackra/goimapnotify/compare/2.3.2...2.3.3
[2.3.2]: https://gitlab.com/shackra/goimapnotify/compare/2.3.1...2.3.2
[2.3.1]: https://gitlab.com/shackra/goimapnotify/compare/2.3...2.3.1
[2.3]: https://gitlab.com/shackra/goimapnotify/compare/2.2...2.3
[2.2]: https://gitlab.com/shackra/goimapnotify/compare/2.1.1...2.2
[2.1.1]: https://gitlab.com/shackra/goimapnotify/compare/2.1...2.1.1
[2.1]: https://gitlab.com/shackra/goimapnotify/compare/2.0...2.1
[2.0]: https://gitlab.com/shackra/goimapnotify/compare/1.1...2.0
[1.1]: https://gitlab.com/shackra/goimapnotify/compare/1.0.1...1.1
[1.0.1]: https://gitlab.com/shackra/goimapnotify/compare/1.0...1.0.1
