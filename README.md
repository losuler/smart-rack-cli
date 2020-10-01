<div align="center">
<p align="center">
  <p align="center">
    <h3 align="center">Smart Rack CLI</h3>
    <p align="center">
      A command line tool to connect to Swinburne's Smart Rack Solution.
    </p>
  </p>
</p>
<img src="preview.gif"/>
</div>
<br>

`srs` is a small command line tool for connecting to Swinburne University's [Smart Rack Solution](https://smartrack.ict.swin.edu.au/) which automates the the kit selection and booking, connection via `ssh` and powering off and release of devices.

## Requires

On Linux and macOS systems `sshpass` is required to be installed.

On Windows systems `putty` is required to be installed.


## Usage

Start by simply extracting and running the single binary in the [latest release](https://github.com/losuler/smart-rack-cli/releases/latest) for your system. For each router or switch you would like to connect to, start a seperate process (e.g. in another terminal window). 

**Note:** Answering `y` to the prompt to shutdown and release will effect ALL routers and switches booked.

## Cisco

To close the `ssh` session, type `Ctrl+]` followed by `quit`.

To disable inline log messages, enter `no logging console` in `configure terminal` mode.
