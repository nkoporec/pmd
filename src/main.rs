use std::thread;

use tokio::sync::broadcast;

use crate::{server::Breakpoint, config::Config};

mod server;
mod config;
mod ui;

fn main() {
    let mut config = Config::new();
    config::parse(&mut config);

    let config_port = config.port.clone();

    let (tx, rx) = broadcast::channel::<Breakpoint>(100);
    thread::spawn(move || { server::run(tx, config_port); });
    ui::render(rx, config);
}
