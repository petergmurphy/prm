---
title: "Installing Tools"
description: "Learn how to install tools to PRM."
category: narrative
tags:
  - tools
  - usage
  - install tool
  - update tool
  - list tools
weight: 10
---

This page contains a guide on how to install tools to PRM. Tools will be installed to a default location unless a different directory is specified using the flag `--toolpath {alt_location}`.

### Local archive

Tool packages can be installed locally using the `prm install` command.

For example, this command:

```bash
prm install ~/my-tool-1.2.3.tar.gz
```

Will install the tool contained in `my-tool-1.2.3.tar.gz` to the default tool location.

### Remote archive

Tool packages stored remotely can be automatically downloaded and extracted with `prm install` so long as you know the URL to where the archive is.

For example, this command:

```bash
prm install https://packages.mycompany.com/prm/my-tool-1.2.3.tar.gz
```

Will attempt to download the PRM tool from the specified url and then afterward install it like any other locally available PRM tool archive.

### Remote git repository

**Git** must be installed for this feature to work. The git repository must contain only one tool and must be structured with the `prm-config.yml` file and the `content` directory in the root directory of the repository.

For example, this command:

```bash
prm install --git-uri https://github.com/myorg/myawesometool
```

This will attempt to clone the PRM tool from the git repository at the specified URI and install to the default tool location.


### Updating tools

Currently, tools can't be updated but rather newer versions of the tool can be installed. Also, currently only the latest version
of a selected tool will execute, the ability to select the version of the tool to be executed will be included in a future PRM version.

### List installed tools

Installed tools can be listed by running the command `prm exec --list`, displayed in the following format:

```
       DISPLAYNAME      |   AUTHOR   |        NAME        |                       PROJECT URL                       | VERSION
------------------------+------------+--------------------+---------------------------------------------------------+----------
  Embedded Puppet (EPP) | puppetlabs | epp                | https://puppet.com/docs/puppet/7/lang_template_epp.html | 0.1.0
  metadata-json-lint    | puppetlabs | metadata-json-lint | https://github.com/voxpupuli/metadata-json-lint         | 0.1.0
```

The `--toolpath {alt_location}` can also be added to list tools installed in an alternate location.
