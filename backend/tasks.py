import command


def set_keyboard_layout(layout: str) -> None:
    command.run(["loadkeys", layout])


def check_uefi():
    out = command.run(["cat", "/sys/firmware/efi/fw_platform_size"])
    return int(out) == 64


def connect_to_internet(ssid: str, passphrase: str, device="wlan0") -> None:
    command.run(["ip", "link", "set", device, "up"])
    command.run(["iwctl", "station", device, "scan"])
    command.run(
        ["iwctl", "--passphrase", passphrase, "station", device, "connect", ssid]
    )


# def partition_disks():
#     # """
#     # We will likely change this once the frontend devs pick a way to partition
#     # for now cfdisk will be used as a proof of concept

#     # for those who don't know what cfdisk is, cfdisk is an interactive cli that will partition the drives for you
#     # """
#     command.run("cfdisk")


def partition(disk: str, fs: str = "btrfs") -> None:
    # """
    # Format is like this disks = ['/dev/sda', '/dev/sdb']
    # """

    # command.run(f"parted {disk} mklabel gpt")
    # command.run(f"parted {disk} mkpart primary fat32 1MiB 512MiB")
    # command.run(f"parted {disk} set 1 esp on")
    # command.run(f"parted {disk} mkpart primary linux-swap 512MiB 17GiB")
    # command.run(f"parted {disk} mkpart primary {fs} 17GiB 100%")
    command.run(["parted", disk, "mklabel", "gpt"])
    command.run(["parted", disk, "mkpart", "primary", "fat32", "1MiB", "512MiB"])
    command.run(["parted", disk, "set", "1", "esp", "on"])
    command.run(["parted", disk, "mkpart", "primary", "linux-swap", "512MiB", "17GiB"])
    command.run(["parted", disk, "mkpart", "primary", fs, "17GiB", "100%"])


def create_fs(mounts: list[str], fs: str = "btrfs") -> None:
    # """
    # Format is like this mounts = ['/', '/dev/sda1', '/boot', '/dev/sda2', 'SWAP', '/dev/sda3', '/home', '/dev/sda4']
    # """
    for mount_point, device in zip(mounts[::2], mounts[1::2]):
        if mount_point != "/" and mount_point != "SWAP":
            command.run(["mkdir", "/mnt" + mount_point])

        if mount_point != "SWAP":
            command.run(["mkfs." + fs, device])
            command.run(["mount", device, "/mnt" + mount_point])
        else:
            command.run(["mkswap", device])
            command.run(["swapon", mount_point])
    # for i in range(0, len(mounts), 2):
    #     mount_point = mounts[i]
    #     device = mounts[i + 1]
    #     if mount_point != "/" and mount_point != "SWAP":
    #         command.run(f"mkdir /mnt{mount_point}")

    #     if mount_point != "SWAP":
    #         command.run(f"mkfs.ext4 {device}")
    #         command.run(f"mount {device} /mnt{mount_point}")
    #     else:
    #         command.run(f"mkswap {device}")
    #         command.run(f"swapon {mount_point}")


# def select_mirrors(mirrors):
#     raise NotImplementedError(
#         "This feature is not implemented in this early version of the installer yet!"
#     )


def pacstrap(pkgs: list[str], kernel: str = "linux") -> None:
    BASE_PKGS = [
        kernel,
        kernel + "-firmware",
        kernel + "-headers",
        "base",
        "base-devel",
    ]
    # command.run(
    #     "pactrap -K /mnt base base-devel linux linux-firmware neovim"
    # )  # Sneaky vim shill... jk it's important to be able edit files
    command.run(["pacstrap", "/mnt", *BASE_PKGS, *pkgs])
    # Bro we use neovim here


def gen_fstab() -> None:
    fstab = command.run(["genfstab", "-U", "/mnt"])
    with open("/mnt/etc/fstab", "w+") as f:
        f.write(fstab)

    # command.run(["genfstab", "-U", "/mnt", ">>", "/mnt/etc/fstab"])
    # command.run(["genfstab, .-U /mnt >> /mnt/etc/fstab")


def chroot() -> None:
    command.run(["arch-chroot", "/mnt"])


def set_timezone(timezone: str) -> None:
    # command.run(f"ln -sf /usr/share/zoneinfo/{timezone}")
    # command.run("hwclock --systohc")
    command.run(["ln", "-sf", f"/usr/share/zoneinfo/{timezone}", "/etc/localtime"])
    command.run(["hwclock", "--systohc"])


def set_locale(layout: str, locale: str = "en_US.UTF-8") -> None:
    with open("/etc/locale.conf", "w") as f:
        f.write(f"LANG={locale}\n")  # FIX LATER
    with open("/etc/vconsole.conf", "w") as f:
        f.write(f"KEYMAP={layout}")


def set_hostname(host: str) -> None:
    with open("/etc/hostname", "w") as f:
        f.write(host)


def set_rootpass(password: str) -> None:
    command.run(["chpasswd", "-R", "/mnt"], input="root:" + password)


def install_bootloader(bootloader: str, uefi: bool, partition: str):
    if bootloader != "grub":
        raise NotImplementedError("This install script only supports grub currently")
    # command.run(
    #     pacman -S networkmanager grub"
    # )  # Users will be pissed if network is down at boot
    # command.run("systemctl enable NetworkManager")
    command.run(["pacman", "-S", "networkmanager", bootloader])
    command.run(["systemctl", "enable", "NetworkManager"])

    if uefi:
        # command.run(
        #     "grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=GRUB"
        # )
        command.run(
            [
                "grub-install",
                "--target=x86_64-efi",
                "--efi-directory=/boot/efi",
                "--bootloader-id=GRUB",
                partition,
            ]
        )
    else:
        # command.run(
        #     f"grub-install {partition}"
        # )  # Keep in mind that this must be the whole partition not the /boot partition
        command.run(["grub-install", partition])

    # result = command.run("grub-mkconfig -o /boot/grub/grub.cfg")
    out = command.run(["grub-mkconfig", "-o", "/boot/grub/grub.cfg"])
    # This will just error?
    if "Found linux image:" not in out:
        raise RuntimeError("Linux kernel not found! Installation aborted!")
    if "Found initrd image" not in out:
        raise RuntimeError("Initramfs not found! Installation aborted!")
