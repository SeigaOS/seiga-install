import tomllib
class Options:
    layout: str
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
