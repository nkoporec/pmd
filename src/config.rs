use std::{
    fs::{self, read_to_string},
    path::PathBuf,
};

use clap::Parser;
use home::home_dir;
use rust_embed::RustEmbed;
use serde::Deserialize;

#[derive(Debug)]
pub struct Config {
    pub port: u16,
    pub keymap: Keymap,
}

#[derive(Debug)]
pub struct Keymap {
    pub leader_key: char,
}

impl Keymap {
    pub fn new() -> Keymap {
        Keymap { leader_key: ',' }
    }
}

impl Config {
    pub fn new() -> Config {
        Config {
            port: 6969,
            keymap: Keymap::new(),
        }
    }
}

#[derive(Debug, Deserialize)]
struct TomlConfig {
    server: Option<ServerTomlConfig>,
    keymap: Option<KeymapTomlConfig>,
}

#[derive(Debug, Deserialize)]
struct ServerTomlConfig {
    port: Option<u16>,
}

#[derive(Debug, Deserialize)]
struct KeymapTomlConfig {
    leader_key: Option<char>,
}

#[derive(RustEmbed)]
#[folder = "stubs/"]
struct Stub;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
pub struct PmdCli {
    /// Sets a custom config file
    #[arg(short, long, value_name = "FILE")]
    config: Option<PathBuf>,

    /// Sets a listening port.
    #[arg(short, long, value_name = "PORT")]
    port: Option<u16>,
}

pub fn parse(config: &mut Config) {
    let cli = PmdCli::parse();

    // Check if we have the default config, otherwise create it.
    let config_file_path = get_config_path();
    if config_file_path.exists() == false {
        create_default_config_file(config_file_path.clone());
    }

    // Check if we pass the config file at runtime
    // otherwise try to load it from the default locations.
    if cli.config.is_some() {
        let config_path = cli.config.expect("Missing config path.");
        parse_toml_config(config, config_path);
    } else {
        parse_toml_config(config, config_file_path);
    }

    // cmd line args have higher prio.
    if cli.port.is_some() {
        config.port = cli.port.unwrap();
    }
}

fn parse_toml_config(config: &mut Config, path: PathBuf) {
    let toml_content = read_to_string(path).expect("Can't read the toml config path.");

    let toml_config: TomlConfig =
        toml::from_str(&toml_content).expect("Can't parse the toml config file.");

    // server config
    if toml_config.server.is_some() {
        let server_config = toml_config.server.unwrap();
        if server_config.port.is_some() {
            let server_port_config = server_config.port.unwrap();
            config.port = server_port_config.clone();
        }
    }

    // keymap config
    if toml_config.keymap.is_some() {
        let keymap_config = toml_config.keymap.unwrap();
        if keymap_config.leader_key.is_some() {
            let leader_key_config = keymap_config.leader_key.unwrap();
            config.keymap.leader_key = leader_key_config.clone();
        }
    }
}

fn get_config_path() -> PathBuf {
    match std::env::consts::OS {
        "linux" => {
            return home_dir().unwrap().join(".config/pmd/config.toml");
        }
        "macos" => {
            return home_dir().unwrap().join(".config/pmd/config.toml");
        }
        "windows" => {
            return home_dir().unwrap().join(".pmd/config.toml");
        }
        _ => {
            panic!("Unsupported operating system")
        }
    }
}

fn create_default_config_file(path: PathBuf) {
    if path.exists() == false {
        let default_config_embed = Stub::get("config.toml").unwrap();
        let default_config = std::str::from_utf8(default_config_embed.data.as_ref()).unwrap();

        fs::write(path, default_config).expect("Could not write the default config");
    }
}
