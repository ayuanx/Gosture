# Gosture
**By AyuanX, 22-Aug-2018**
___________________

## What is Gosture
Gosture is a Linux X Window system productivity tool, written in Go language.  
Gosture implements system-wide customizable mouse gestures as well as keyboard shortcuts. 

## How to use Gosture
1. Prepare your Gosture configuration file at `~/.Gosture.cfg`.  
   An example of configuration file is provided as `Gosture_Config_Example.txt`.  
   You can use it as a starting point. E.g. run `cp Gosture_Config_Example.txt ~/.Gosture.cfg`
2. Tweak the configuration to your needs. (Check detailed instructions below.)
3. Run Gosture in background and enjoy the mouse gestures and keyboard shortcuts. E.g. run `nohup ./Gosture &`  
   After launching, you can control it through the icon in system tray.  
   Tip: you can add it into your X Window Startup Applications so that it automatically runs at your login.

The following operations are defined by default in the provided configuration example:

Operation | Action  
--------- | ------
Mouse Middle Button | Trigger a mouse gesture
Gesture ↙ | Minimize active window
Gesture ↗ | Maximize / Restore active window
Gesture ↓→ | Close active window
Gesture ↑ | Scroll to top (Equivalent to Home key)
Gesture ↓ | Scroll to bottom (Equivalent to End key)
Gesture ← | Copy selection to clipboard (Equivalent to Ctrl+Insert)
Gesture → | Paste from clipboard (Equivalent to Shift+Insert)
Gesture ↑↓ | Snap window to top edge (Equivalent to Super+Up)
Gesture ↓↑ | Snap window to bottom edge (Equivalent to Super+Down)
Gesture ←→ | Snap window to left edge (Equivalent to Super+Left)
Gesture →← | Snap window to right edge (Equivalent to Super+Right)
Super+Alt+Z | Run gedit *(Super key is also known as Windows key)*
Super+Alt+X | Run terminal *(Super key is also known as Windows key)*
Super+Alt+C | Run calculator *(Super key is also known as Windows key)*

## How to configure Gosture
Gosture configuration file `~/.Gosture.cfg` is a standard JSON file.

Option | Description
------ | -----------
`mouse-gesture-enable`  |	`true`: enable mouse gesture; `false`: disable mouse gesture
`mouse-gesture-trigger` |	can be a single **[Mouse Button]** like `2`; or a **[Modifier Key]-[Mouse Button]** combination, like `Control-2`

### Definition of mouse buttons:

Mouse Button | Description
--- | ---
`1` | Left Button
`2` | Middle Button
`3` | Right Button
`4` | Scroll Up
`5` | Scroll Down

### Definition of keys:

Modifier Key | Description
------------ | -----------
`Shift` | Shift Key
`Lock` | Caps Lock Key
`Control` | Ctrl Key
`Mod1` | Alt Key
`Mod2` | Num Lock Key
`Mod3` | *(Usually not mapped to any physical key)*
`Mod4` | Super Key (also known as Windows Key) 
`Mod5` | AltGr Key (usually absent on US keyboard) 
* To get an acurate list of all modifier keys in your system, run `xmodmap`.
* To find a specific key name, run `xev`.  
   You can also reference these documents:  
   http://xahlee.info/linux/linux_show_keycode_keysym.html  
   http://wiki.linuxquestions.org/wiki/List_of_Keysyms_Recognised_by_Xmodmap

### Definition of mouse gestures:
All eight directions are supported; directions are mapped to digits on **Num Pad**.

.       | .       | .
 ------ | ------- | ------
`7` (↖) | `8` (↑) | `9` (↗)
`4` (←) |         | `6` (→)
`1` (↙) | `2` (↓) | `3` (↘)
* For example: gesture of "↑→" is `86`; gesture of "↖↘" is `73`; gesture of "←↓→" is `426`.  
  Tip: Mixture of orthogonal stroke and diagonal stroke in one gesture (like "↗→" or "↙↓↘") is supported, but not recommended.


### Current supported actions:

Action | Description
------ | -----------
`minwin`    | Minimize active window
`maxwin`    | Maximize active window / Restore it if already maximized
`closewin` | Close active window
`key,[key1],[key2],...` | Send a key combination. Each key is delimited by comma
`cmd,[executable],[dir]`	| Run [executable], can be a program or script with arguments. [dir] is optional working directory

## Dependencies and credits
* https://github.com/BurntSushi/xgbutil
* https://github.com/mattn/go-shellwords
* https://github.com/getlantern/systray
* https://github.com/skratchdot/open-golang

