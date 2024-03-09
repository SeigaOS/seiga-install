# Put functions here
import os

def connect:
    """ Connects machine to WiFI network"""
    os.system("iwctl device wlan0 set-property Powered on")
    os.system("iwctl station wlan0 scan")
    os.system("iwctl staTion wlan0 get-networks")
    ssid = input("Please enter your SSID >>> ")
    os.system(f"iwctl station wlan0 connect {ssid}")
