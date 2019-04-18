## Android Media Backups

Sometimes there are connection problems with getting android file transfer to work (i.e. Google Pixel).

1. Enable USB debugging on the device
1. Download [Android Debug Bridge (adb)](https://developer.android.com/studio/command-line/adb)
3. Install [Golang (Go)](http://golang.org)


### Run

Using the command line, you can check out this project, and run it using Go:

    $ git clone git@github.com:Xeoncross/android_media_backup.git
    $ cd android_media_backup
    $ go run android_media_backup.go

By default the script downloads everything from `DCIM/Camera`, but you can change this when calling the script.

    $ go run android_media_backup.go -dir="WhatsApp"

### Warning

Files are removed from device after being transferred to the computer. This required when dealing with phones that continue to disconnect in the middle of transfers.


### More about adb

Android Debug Bridge (adb) can be used directly from the terminal. You can use it to access the files manually for coping by starting a new shell instance to browse the device folders:

    adb shell

or run a command remotely

    adb shell ls /sdcard/DCIM

Then you can use `adb pull ...` to download the files locally.

    adb pull -a /sdcard/DCIM/Camera/ ./
    adb pull -a /sdcard/Snapseed ./
    adb pull -a /sdcard/WhatsApp ./
