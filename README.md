# FlexiKiller

These are simple removal utilities for Windows and Mac versions of the FlexiSpy trojan.
The "good" part of FlexiSpy is that the bundles intalled on infected systems already contain appropriate uninstall utilities. They are obviously hidden and they are normally executed either through remote trigger or by typing a combination of keys on the computer.

These utilities should be able to find active (and inactive) FlexiSpy infections and locate the uninstall utilities and launch them.

The binary for the Windows utility is at:

    https://ops.securitywithoutborders.org/flexispy/FlexiKiller.exe

For Mac it is just a bash script which can be executed with (yep, curling bash scripts, YOLO):

    $ bash <(curl -s https://ops.securitywithoutborders.org/flexispy/kill.sh)
