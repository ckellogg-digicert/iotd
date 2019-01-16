package main

import (
	"fmt"
	"os"
)

func getCmd(path string) (cmd string) {
	switch os.Getenv("XDG_CURRENT_DESKTOP") {
	case "XFCE":
		// TODO: Check this scenario
		cmd = `xres=($(echo $(xfconf-query --channel xfce4-desktop --list | grep last-image)))
for x in "${xres[@]}"
do
	xfconf-query --channel xfce4-desktop --property $x --set %s
done`
	case "MATE":
		cmd = "gsettings set org.mate.background picture-filename %s"
	default:
		cmd = `if gsettings set org.gnome.desktop.background picture-uri "file://%s"; then
	gsettings set org.gnome.desktop.background picture-options "zoom"
else
	echo "$XDG_CURRENT_DESKTOP not supported."
	break
fi`
	}

	return fmt.Sprintf(cmd, path)
}
