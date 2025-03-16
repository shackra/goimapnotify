<a name="unreleased"></a>
## [Unreleased]

### Changes
- Don't suppress stderr of executed commands

### Fixes
- "Limit number of restarts for systemd job"
- "Update README.md to reflect extension of `-conf` flag"


<a name="2.4"></a>
## [2.4] - 2024-09-30
### Changes
- Merge branch `2.3.x` into `master`
- Change configuration format in readme to YAML
- Change installation instructions
- Remove privacy alert in command description
- Display a message asking for donations :D
- Update to Go 1.20
- Remove Terra Luna donation address
- Advertise donations through some stablecoins
- Refactor configuration-related code

### Features
- Support for monitoring all mailboxes
- Add YAML support
- Send IMAP ID
- Add Code of Conduct
- Keep privacy of users censoring credentials

### Fixes
- Switch from TLS1.3 to TLS1.2
- Display donation message on stderr
- Ensure the List call from `printDelimiter` finishes before returning
- Resolve ambiguous TLS and config extension
- Improve conditional check at watch.go
- Missing commas in the JSON example in README.md
- Resolve "Improve/re-do changelog"
- Display all errors when trying to load the configuration

### Reverts
- Merge branch '2.3.x' into 'master'


<a name="2.3.16"></a>
## [2.3.16] - 2024-09-01
### Changes
- Change configuration format in readme to YAML
- Change installation instructions
- Remove privacy alert in command description
- Display a message asking for donations :D

### Features
- Support for monitoring all mailboxes
- Add YAML support
- Send IMAP ID
- Add Code of Conduct

### Fixes
- Display donation message on stderr
- Ensure the List call from `printDelimiter` finishes before returning
- Resolve ambiguous TLS and config extension
- Improve conditional check at watch.go
- Missing commas in the JSON example in README.md


<a name="2.3.15"></a>
## [2.3.15] - 2024-04-27
### Features
- Keep privacy of users censoring credentials


<a name="2.3.14"></a>
## [2.3.14] - 2024-04-27
### Fixes
- Resolve "Improve/re-do changelog"


[Unreleased]: https://gitlab.com/shackra/goimapnotify/compare/2.4...HEAD
[2.4]: https://gitlab.com/shackra/goimapnotify/compare/2.3.16...2.4
[2.3.16]: https://gitlab.com/shackra/goimapnotify/compare/2.3.15...2.3.16
[2.3.15]: https://gitlab.com/shackra/goimapnotify/compare/2.3.14...2.3.15
[2.3.14]: https://gitlab.com/shackra/goimapnotify/compare/2.3.13...2.3.14
