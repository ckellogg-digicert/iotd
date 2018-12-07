# IotD

Go-based application to enable setting your macOS wallpaper to match Bing's daily background image.


## Installation

Currently you have to build the binary manually with go. A installer script will be coming.

1. Build with go: `go build -o iotd main.go`
2. Move _iotd_ binary to _/usr/local/bin_: `mv iotd /usr/local/bin`
3. Copy plist file: `cp com.github.thoom.iotd.plist ~/Library/LaunchAgents`
4. Enable launch agent: `launchctl load ~/Library/LaunchAgents/com.github.thoom.iotd.plist`


The launch agent currently stores the output log to `/tmp/iotd.out`. Running the script manually outputs to _STDOUT_.