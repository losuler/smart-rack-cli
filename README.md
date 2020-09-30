# Smart Rack Solution

`srs` is a small command line program designed for Swinburne University Smart Rack Solution that quickly connects you to the router or switch of your choice.

## Dependencies

Linux/macOS:

```
sshpass
```

Windows:

```
putty
```

## Usage

Run by simply running the executable in releases. For each router or switch, start a seperate process. Answering yes to the prompt for shutdown and release will shutdown ALL the sessions created from each process.

In Linux/macOS:

```
Example output
```

In Windows:

```
Example output
```

## Notes

To quit the `ssh` session, type `Ctrl+]` followed by `quit`.

To disable inline log messages, enter `no logging console` in `configure terminal`.
