# MOUNT A BITLOCKER DRIVE

## Using dislocker

> ### Install dislocker
>
>     sudo apt-get install dislocker

> ### Create Folders to mount VHD
>
>     sudo mkdir -p /media/bitlocker; sudo mkdir -p /media/bitlockermount

> ### Configure File as Loop Device 
>
>     sudo losetup -f -P <FILE>.vhd

> ### Decrypt
>
>     sudo dislocker /dev/loop0p1 -u<PASSWORD> -- /media/bitlocker

> ### Mount
>
>     sudo mount -o loop /media/bitlocker/dislocker-file /media/bitlockermount

> ### Unmount
>
>     sudo umount /media/bitlocker; sudo umount /media/bitlockermount