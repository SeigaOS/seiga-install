import tomllib
import subprocess

# Only to be used as a struct
class Options:
    def __init__(self, filename):
        with open(filename, 'rb') as f:
            data = tomllib.load(f)

        self.layout = data['Keyboard']['layout']
        self.mirrors = data['Mirrors']['mirrors']
        #self.partitions = data['Partitions']['partitions']
        self.mounts = data['Mount']['mount']
        self.swapon = data['Swap']['on']
        self.swapfile = data['Swap']['file']
        self.host = data['Host']['host']
        self.rootpass = data['Root']['password']
        self.wheel_users = data['Wheel-Users']['wheel-users']
        self.normal_users = data['Normal-Users']['normal-users']
        self.audio = data['Audio']['audio']
        self.kernel = data['Kernel']['kernel']
        self.ssid = data['Network']['ssid']
        self.passphrase = data['Network']['passphrase']
        self.tz = data['Timezone']['timezone']

def run(cmd):
    result = subprocess.run(cmd, shell=True, check=True)
    if result.returncode == 0:
        output_str = result.stdout
    else:
        raise RuntimeError(f"Command: {cmd} failed with return code {result.returncode}\n{result.stderr}")

    if output_str is None:
        raise RuntimeError(f"Expected output from command {cmd} but recived 'None'.\nIf no output was expeted use function 'void_run' instead")
    return output_str

# Will not return stdout
def void_run(cmd):
    subprocess.run(cmd, shell=True, check=True)

# This is static method class 
class Install:
    def set_keyboard_layout(layout):
        void_run(f"loadkeys {layout}")
    
    def check_uefi():
        booted_bits = run("cat /sys/firmware/efi/fw_platform_size")
        if booted_bits = 64:
            return True
        else:
            return False

    def connect_to_internet(ssid, passphrase):
        void_run("iwctl device wlan0 set-property Powered on")
        void_run("iwctl station wlan0 scan")
        void_run(f"iwctl --passphrase {passphrase} station wlan0 connect {ssid}")

    def partition_disks():
        """
        We will likley change this once the frontend devs pick a way to partition
        for now cfdisk will be used as a proof of concept

        for those who don't know what cfdisk is, cfdisk is an interactive cli that will partition the drives for you
        """
        void_run("cfdisk")

    def mount_fs(mounts):
        """
        Format is like this mounts = ['/', '/dev/sda1', '/boot', '/dev/sda2', 'SWAP', '/dev/sda3', '/home', '/dev/sda4']
        """
        for i in range(0, len(mounts), 2):
            mount_point = mounts[i]
            device = mounts[i+1]
            if mount_point != "/" and mount_point != "SWAP":
                void_run(f"mkdir /mnt{mount_point}")
            
            if mount_point != "SWAP":
                void_run(f"mkfs.ext4 {device}")
                void_run(f"mount {device} /mnt{mount_point}")
            else:
                void_run(f"mkswap {device}")
                void_run(f"swapon {mount_point}")

    def select_mirrors(mirrors):
        raise NotImplementedError("This feature is not implemented in this early version of the installer yet!")

    def pactrap():
        void_run("pactrap -K /mnt base base-devel linux linux-firmware vim") # Sneaky vim shill... jk it's important to be able edit files

    def gen_fstab():
        void_run("genfstab -U /mnt >> /mnt/etc/fstab")

    def chroot():
        void_run("arch-chroot /mnt")

    def set_timezone(timezone):
        void_run(f"ln -sf /usr/share/zoneinfo/{timezone}")
        void_run("hwclock --systohc")

    def set_locale(layout):
        with open("/etc/locale.conf", 'w') as f:
            f.write("LANG=en_US.UTF-8\n") # FIX LATER
        with open("/etc/vconsole.conf", 'w') as f:
            f.write(f"KEYMAP={layout}")

    def set_hostname(host):
        with open("/etc/hostname", 'w') as f:
            f.write(host)

    def set_rootpass(password):
        void_run(f"echo '{password}' | passwd")
    
    def install_bootloader(bootloader, isuefi, partition):
        if bootloader != "grub":
            raise NotImplementedError("This install script only supports grub currently")
        void_run("pacman -S networkmanager grub") # Users will be pissed if network is down at boot
        void_run("systemctl enable NetworkManager")

        if isuefi:
            void_run("grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=GRUB")
        else:
            void_run(f"grub-install {partition}") # Keep in mind that this must be the whole partition not the /boot partition
        
        result = run("grub-mkconfig -o /boot/grub/grub.cfg")
        if "Found linux image:" not in result:
            raise RuntimeError("Linux kernel not found! Installation aborted!")
        if "Found initrd image" not in result:
            raise RuntimeError("Initramfs not found! Installation aborted!")


if __name__ == '__main__':
    options = Options("options.toml")
    
