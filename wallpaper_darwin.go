package main

import "fmt"

func getCmd(path string) string {
	applScript := `/usr/bin/osascript<<END
tell application "System Events"
	set picture of every desktop to POSIX file "%s"
end tell
END`

	return fmt.Sprintf(applScript, path)
}
