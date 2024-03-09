import os

def connect() -> None:
    """ Connects machine to WiFI network"""
    os.system("iwctl device wlan0 set-property Powered on")
    os.system("iwctl station wlan0 scan")
    os.system("iwctl staTion wlan0 get-networks")
    ssid = input("Please enter your SSID >>> ")
    os.system(f"iwctl station wlan0 connect {ssid}")



def partition() -> None:
    """ Partitions the disk"""
    os.system("lsblk -o NAME")
    disk = input("Please enter the disk you want to partition >>> ")
    os.system(f"cfdisk /dev/{disk}")