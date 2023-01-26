use std::{
    io::Stdout,
    process,
    time::{Duration, Instant},
};

use crossterm::event::{self, Event, KeyCode};
use serde::Deserialize;
use tui::{backend::CrosstermBackend, Terminal};

use crate::ui::{layout::tree, UiState};

#[derive(Debug, Deserialize, Clone, Eq, Hash, PartialEq)]
pub enum InputMode {
    Normal,
    Inspection,
    Visual,
}

#[derive(Debug, Deserialize, Clone, Eq, Hash, PartialEq)]
pub enum ListMode {
    Breakpoint,
    Callstack,
}

impl InputMode {
    pub fn to_string(&self) -> String {
        match self {
            InputMode::Normal => {
                return String::from("Normal");
            }
            InputMode::Inspection => {
                return String::from("Inspection");
            }
            InputMode::Visual => {
                return String::from("Visual");
            }
        }
    }
}

pub fn handle_keymap_event(state: &mut UiState, terminal: &mut Terminal<CrosstermBackend<Stdout>>) {
    if let Event::Key(key) = event::read().unwrap() {
        match state.input_mode {
            InputMode::Normal => {
                normal_keymap(key.code, state, terminal);
            }
            InputMode::Inspection => {
                inspection_keymap(key.code, state, terminal);
            }
            InputMode::Visual => {
                visual_keymap(key.code, state, terminal);
            }
        }
    }
}

pub fn normal_keymap(
    key_code: KeyCode,
    state: &mut UiState,
    terminal: &mut Terminal<CrosstermBackend<Stdout>>,
) {
    let leader_key_timeout = Duration::from_millis(1000);
    let leader_key = KeyCode::Char(state.config.keymap.leader_key);

    match key_code {
        // leader key.
        leader if leader == leader_key => {
            state.leader_tick_time = Instant::now();
        }
        KeyCode::Char('l') => {
            if leader_key_timeout >= state.leader_tick_time.elapsed() {
                state.list_state.breakpoint.selected = state.list_state.breakpoint.state.selected();
                state.list_state.callstack.select_first();
                state.list_mode = ListMode::Callstack;
            }
        }
        KeyCode::Char('h') => {
            if leader_key_timeout >= state.leader_tick_time.elapsed() {
                state.list_state.callstack.unselect();
                state.list_mode = ListMode::Breakpoint;
                state
                    .list_state
                    .breakpoint
                    .state
                    .select(state.list_state.breakpoint.selected);
            }
        }
        KeyCode::Char('q') => {
            terminal.clear().unwrap();
            process::exit(0);
        }
        KeyCode::Char('j') => match state.list_mode {
            ListMode::Breakpoint => {
                state.list_state.breakpoint.next();
                state.tree_state.items =
                    tree::build_tree_items(state.list_state.breakpoint.get_value());
                state.list_state.callstack.items = state.list_state.breakpoint.get_callstack();
            }
            ListMode::Callstack => {
                state.list_state.callstack.next();
            }
        },
        KeyCode::Down => match state.list_mode {
            ListMode::Breakpoint => {
                state.list_state.breakpoint.next();
                state.tree_state.items =
                    tree::build_tree_items(state.list_state.breakpoint.get_value());
                state.list_state.callstack.items = state.list_state.breakpoint.get_callstack();
            }
            ListMode::Callstack => {
                state.list_state.callstack.next();
            }
        },
        KeyCode::Char('k') => match state.list_mode {
            ListMode::Breakpoint => {
                state.list_state.breakpoint.previous();
                state.tree_state.items =
                    tree::build_tree_items(state.list_state.breakpoint.get_value());
                state.list_state.callstack.items = state.list_state.breakpoint.get_callstack();
            }
            ListMode::Callstack => {
                state.list_state.callstack.previous();
            }
        },
        KeyCode::Up => match state.list_mode {
            ListMode::Breakpoint => {
                state.list_state.breakpoint.previous();
                state.tree_state.items =
                    tree::build_tree_items(state.list_state.breakpoint.get_value());
                state.list_state.callstack.items = state.list_state.breakpoint.get_callstack();
            }
            ListMode::Callstack => {
                state.list_state.callstack.previous();
            }
        },
        KeyCode::Char('v') => {
            state.popup.show();
            match state.list_mode {
                ListMode::Breakpoint => {
                    let breakpoint = state
                        .list_state
                        .breakpoint
                        .get_selected_breakpoint()
                        .unwrap();

                    let mut popup_text = vec![];

                    let breakpoint_file = String::from(format!("[FILE]: {}", breakpoint.filepath));
                    popup_text.push(breakpoint_file);

                    let breakpoint_line = String::from(format!("[LINE]: {}", breakpoint.line));
                    popup_text.push(breakpoint_line);

                    let breakpoint_timestamp =
                        String::from(format!("[TIME]: {}", breakpoint.timestamp));
                    popup_text.push(breakpoint_timestamp);

                    let breakpoint_connector =
                        String::from(format!("[CONNECTOR]: {}", breakpoint.connector_type));
                    popup_text.push(breakpoint_connector);

                    state.popup.set_text(popup_text);
                }
                ListMode::Callstack => {
                    let mut popup_text = vec![];
                    let callstack = state.list_state.callstack.get_selected_callstack().unwrap();

                    let callstack_file = String::from(format!("[FILE]: {}", callstack.filepath));
                    popup_text.push(callstack_file);

                    let callstack_line = String::from(format!("[LINE]: {}", callstack.line));
                    popup_text.push(callstack_line);

                    state.popup.set_text(popup_text);
                }
            }
            state.input_mode = InputMode::Visual;
            state
                .status_bar
                .set_status(format!(":{}", state.input_mode.to_string()));
        }
        KeyCode::Char('i') => {
            state.input_mode = InputMode::Inspection;
            state
                .status_bar
                .set_status(format!(":{}", state.input_mode.to_string()));
        }
        KeyCode::Esc => {
            if state.popup.is_active() {
                state.popup.hide();
            }
        }
        _ => {}
    }
}

pub fn inspection_keymap(
    key_code: KeyCode,
    state: &mut UiState,
    terminal: &mut Terminal<CrosstermBackend<Stdout>>,
) {
    match key_code {
        KeyCode::Esc => {
            state.input_mode = InputMode::Normal;
            state
                .status_bar
                .set_status(format!(":{}", state.input_mode.to_string()));
        }
        KeyCode::Char('l') => {
            state.tree_state.right();
        }
        KeyCode::Left => {
            state.tree_state.left();
        }
        KeyCode::Char('h') => {
            state.tree_state.left();
        }
        KeyCode::Right => {
            state.tree_state.left();
        }
        KeyCode::Char('j') => {
            state.tree_state.down();
        }
        KeyCode::Down => {
            state.tree_state.down();
        }
        KeyCode::Char('k') => {
            state.tree_state.up();
        }
        KeyCode::Up => {
            state.tree_state.up();
        }
        KeyCode::Char('q') => {
            terminal.clear().unwrap();
            process::exit(0);
        }
        _ => {}
    }
}

pub fn visual_keymap(
    key_code: KeyCode,
    state: &mut UiState,
    terminal: &mut Terminal<CrosstermBackend<Stdout>>,
) {
    match key_code {
        KeyCode::Esc => {
            if state.popup.is_active() {
                state.popup.hide();
                state.input_mode = InputMode::Normal;
                state
                    .status_bar
                    .set_status(format!(":{}", state.input_mode.to_string()));
            }
        }
        KeyCode::Char('q') => {
            terminal.clear().unwrap();
            process::exit(0);
        }
        _ => {}
    }
}
