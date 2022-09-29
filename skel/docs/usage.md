# Usage

Skel may be run using either command-line sub-commands or by using its minimal, but hopefully helpful, interactive mode.

* To start the interactive project generator, run the skel command with no arguments:
    ```bash
    skel
    ```

* To list the targets and options for the skel command, add the `-h` flag:
    ```bash
    skel -h
    ```

* To get help for a particular target:
    ```bash
    skel <target> -h
    ```

In addition to the numerous generation targets, there are the following utility targets:

- `help`: display the help text
- `version`: display the current, and most recent `skel` build versions
- `completion`: generate the BASH completion script for `skel`
