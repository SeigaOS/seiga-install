use interface_proc::embed_schema;
use serde::Deserialize;
use toml::from_str;

pub fn get_schema() -> Schema {
    from_str(embed_schema!()).expect("invalid schema")
}

#[derive(Debug, Clone, Deserialize)]
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
    tz: String
}
#[derive(Debug, Clone, Deserialize)]
pub struct Keyboard {
    layout: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize)]
pub struct Mirrors {
    Mirrors: Vec<String>,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize)]
pub struct Fs {
    fs: String, 
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize)]
pub struct Bootloader {
    bootloader: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize)]
pub struct Swap {
    on: bool,
    r#type: String,
    file: String,
    mnt: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize)]
pub struct Audio {
    audio: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize)]
pub struct Kernel {
    kernel: String,
    possible: Vec<String>
}
#[derive(Debug, Clone, Deserialize)]
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
}