use std::{io, sync::{Arc, Mutex}, time::{Duration, Instant}};

use crossterm::{terminal::enable_raw_mode, event::poll};
use tokio::sync::broadcast::Receiver;
use tui::{backend::CrosstermBackend, Terminal, text::Spans};

use crate::{server::Breakpoint, ui::{keymap::{InputMode, ListMode}, layout::{BreakpointList, CallstackList}}, config::Config};

mod keymap;
mod layout;

// Our shared state
#[derive(Debug)]
pub struct UiState {
    pub list_state: ListState,
    pub list_mode: ListMode, 
    pub input_mode: InputMode, 
    pub variables: String,
    pub popup: Popup,
    pub status_bar: StatusBar,
    pub leader_tick_time: Instant,
    pub config: Config,
}

#[derive(Debug)]
pub struct ListState {
    pub breakpoint: BreakpointList,
    pub callstack: CallstackList,
}

impl ListState {
    pub fn new() -> ListState {
        ListState {
            breakpoint: BreakpointList::with_items(vec![]),
            callstack: CallstackList::with_items(vec![]),
        }
    }
}


#[derive(Debug)]
pub struct Popup {
    pub show: bool,
    pub text: Vec<String>,
}

impl Popup {
    pub fn new() -> Popup {
        Popup {
            show: false,
            text: vec![],
        }
    }

    pub fn is_active(&self) -> bool {
        return self.show;
    }

    pub fn show(&mut self) {
        self.show = true;
    }

    pub fn hide(&mut self) {
        self.show = false;
        self.set_text(vec![]);
    }

    pub fn set_text(&mut self, text: Vec<String>) {
        self.text.clear();
        self.text = text;
    }

    pub fn get_text(&self) -> Vec<String> {
        return self.text.clone();
    }
}

#[derive(Debug)]
pub struct StatusBar {
    pub text: String,
}

impl StatusBar {
    pub fn new() -> StatusBar {
        StatusBar {
            text: String::from(""),
        }
    }

    pub fn set_status(&mut self, text: String) {
        self.text = text;
    }

    pub fn get_status(&mut self) -> String {
        return self.text.clone();
    }
}

impl UiState {
    pub fn new(config: Config) -> UiState {
        UiState {
            list_state: ListState::new(),
            variables: "".to_string(),
            popup: Popup::new(),
            status_bar: StatusBar::new(),
            input_mode: InputMode::Normal,
            leader_tick_time: Instant::now(),
            list_mode: ListMode::Breakpoint,
            config,
        }
    }
}

#[tokio::main(worker_threads = 1)]
#[allow(unused_must_use)]
pub async fn render(mut rx: Receiver<Breakpoint>, config: Config) {
    // Set up terminal output
    enable_raw_mode().unwrap();
    let stdout = io::stdout();
    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend).unwrap();

    // Clear the terminal before first draw.
    terminal.clear().unwrap();

    let data = Arc::new(Mutex::new(vec![]));
    let server_msgs = Arc::clone(&data);
    let tui_msgs = Arc::clone(&data);

    tokio::spawn(async move {
        loop {
            let msg = rx.recv().await.unwrap();
            server_msgs.lock().unwrap().push(msg);
        }
    });

    // Create a new ui state.
    let mut state = UiState::new(config);
    loop {
        if poll(Duration::from_millis(1)).unwrap() {
            // It's guaranteed that `read` won't block, because `poll` returned
            // `Ok(true)`.
            keymap::handle_keymap_event(&mut state, &mut terminal);
        }

        // Lock the terminal and start a drawing session.
        terminal.draw(|f| {
            state.list_state.breakpoint.items = tui_msgs.lock().unwrap().to_vec();
            layout::render_main(&mut state, f);
        }).unwrap();
    }
}
