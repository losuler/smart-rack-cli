<div align="center">
<p align="center">
  <p align="center">
    <h3 align="center">Smart Rack CLI</h3>
    <p align="center">
      A command line tool to connect to Swinburne's Smart Rack Solution.
    </p>
  </p>
</p>
<img src="img/preview.gif"/>
</div>
<br>

`srs` is a small command line tool for connecting to Swinburne University's [Smart Rack Solution](https://smartrack.ict.swin.edu.au/) which automates the the kit selection and booking, connection via `ssh` and powering off and release of devices.

## Requires

In order to connect to Smart Rack you are required to be connected to the Swinburne VPN. For all operating systems it's recommended to use OpenConnect (an open source client for Cisco's AnyConnect SSL VPN).

It is packaged in [most](https://repology.org/project/openconnect/versions) Linux distros (as a CLI tool and in NetworkManager upstream), as [a package](https://ports.macports.org/port/openconnect/summary) in MacPorts and [Homebrew](https://formulae.brew.sh/formula/openconnect) for macOS and as a [client](https://openconnect.github.io/openconnect-gui/) on Windows.

### Linux

```
sshpass
```

`sshpass` has been packaged on [most](https://repology.org/project/sshpass/versions) Linux distros.

### macOS

```
sshpass
```

On MacPorts, `sshpass` has been [packaged](https://ports.macports.org/port/sshpass/summary) in the offical repo.

On Homebrew `sshpass` [is not](https://github.com/Homebrew/brew/commit/04dfdd972c7fca25e86e9e2ff7767b9f5b789f20) in the official repo. However a Homebrew [Tap](https://docs.brew.sh/Taps) exists, the source of which is [here](https://github.com/hudochenkov/homebrew-sshpass/blob/master/sshpass.rb). This can be installed by running `brew install hudochenkov/sshpass/sshpass`.

### Windows

```
putty
```

It's also recommended to use the official [Windows Terminal](https://github.com/microsoft/terminal) over the pre-installed Command Prompt.

## Usage

Run by extracting and executing the binary named `srs` from the [latest release](https://github.com/losuler/smart-rack-cli/releases/latest) archive for your system.

A kit will be randomly selected from those not in use. For each router or switch you would like to connect to, start a seperate process.

**Note:** Answering `y` to the prompt to shutdown and release will effect ALL routers and switches booked.

## Cisco

To close the `ssh` session, type `Ctrl+]` followed by `quit`.

To disable inline log messages, enter `no logging console` in `configure terminal` mode.
