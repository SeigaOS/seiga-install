use std::{env::current_exe, fs::File, io::Write};

use interface_proc::embed_schema;
use serde::{Deserialize, Serialize};
use toml::from_str;

pub fn get_schema() -> Schema {
    from_str(embed_schema!()).expect("invalid schema")
}

pub fn write2ron() {
    let mut path = current_exe().unwrap();
    path.pop();
    path.push("options.ron");
    let mut file = File::create(path).unwrap();
    file.write_all(ron::to_string(&get_schema()).unwrap().as_bytes()).unwrap();
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Schema {
    wheel_users: Vec<String>,
    normal_users: Vec<String>,
    keyboard: Keyboard,
    mirrors: Mirrors,
    fs: Fs,
    bootloader: Bootloader,
    swap: Swap,
    host: String,
    root_password: String,
    audio: Audio,
    kernel: Kernel,
    network: Network,
    time: Time
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct Time {
    timezone: String
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Keyboard {
    layout: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Mirrors {
    mirrors: Vec<String>,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Fs {
    fs: String, 
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Bootloader {
    bootloader: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Swap {
    on: bool,
    r#type: String,
    file: String,
    mnt: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Audio {
    audio: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Kernel {
    kernel: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Network {
    provider: String,
    possible: Vec<String>
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn check_schema() {
        get_schema();
    }

    #[test]
    fn write() {
        write2ron();
    }
}