import options as o
import tasks
import aya

if __name__ == "__main__":
    # Init logging
    log = aya.Logger(template="{TYPE} - {FILE}: {LINE} {MSG}",logfile="logs.txt")
    
    log.debug("Getting TOML options")
    options = o.Options("options.toml")
    
    log.info("Starting install")
    
    tasks.set_keyboard_layout(options.layout)
    log.info("Set keyboard layout to " + options.layout)
    
    uefi = tasks.check_uefi()
    log.debug("uefi on? " + uefi)
    
    tasks.connect_to_internet(options.ssid, options.passphrase)
    log.info("Successfully connected to internet")

    log.critical("SYSTEM IS ABOUT TO BE PARTITIONED. ALL DATA WILL BE DESTROYED")
    log.warn("Type \"YES\" to continue")
    confirm = input(">>> ")
    if confirm != "YES":
        exit()
    
    tasks.partition("/dev/sda") # Please don't put this in prod
    log.warn("Partition table created")
    
    tasks.create_fs(options.mounts, fs="ext4")
    log.info("Disks formated and mounted!")

    tasks.pacstrap(pkgs=["vim"], kernel=options.kernel)
    log.info("System pacstraped!")

    tasks.gen_fstab()
    log.debug("Fstab generated")

    tasks.chroot()
    log.warn("No longer in live env! System has chrooted!")

    tasks.set_timezone(options.tz)

    tasks.set_locale(layout=options.layout, locale="en_US.UTF-8")
    log.info("set locale")

    tasks.set_rootpass(options.rootpass)
    log.info("Root password correctly set")

    tasks.install_bootloader("grub", uefi, "/dev/sda") # PLEASE don't keep this in prod
    log.info("INSTALLATION COMPLETE! Please reboot!")
    

