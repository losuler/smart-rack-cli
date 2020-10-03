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

In order to connect Smart Rack you are required to be connected to the Swinburne VPN. For all operating systems I recommend OpenConnect.

There's a [Windows client](https://openconnect.github.io/openconnect-gui/), [a package](https://formulae.brew.sh/formula/openconnect) in Homebrew for macOS and is packaged in most distros (as a CLI tool and in NetworkManager/GNOME upstream).

### Linux/macOS

On Linux and macOS `sshpass` is required to be installed.

### Windows

On Windows `putty` is required to be installed.

I recommend the official [Windows Terminal](https://github.com/microsoft/terminal) over the pre-installed Command Prompt.

## Usage

Run by extracting and executing the binary named `srs` from the [latest release](https://github.com/losuler/smart-rack-cli/releases/latest) archive for your system.

A kit will be randomly selected from those not in use. For each router or switch you would like to connect to, start a seperate process.

**Note:** Answering `y` to the prompt to shutdown and release will effect ALL routers and switches booked.

## Cisco

To close the `ssh` session, type `Ctrl+]` followed by `quit`.

To disable inline log messages, enter `no logging console` in `configure terminal` mode.
