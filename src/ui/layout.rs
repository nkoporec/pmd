use std::io::Stdout;

use tui::{
    backend::CrosstermBackend,
    layout::{Constraint, Direction, Layout, Rect},
    style::{Color, Style},
    text::Spans,
    widgets::{Block, Borders, Clear, List, ListItem, ListState, Paragraph},
    Frame,
};

use crate::{
    server::{Breakpoint, Callstack},
    ui::UiState,
};

// @todo: Use traits for this.
#[derive(Debug)]
pub struct BreakpointList {
    pub state: ListState,
    pub selected: Option<usize>,
    pub items: Vec<Breakpoint>,
}

#[derive(Debug)]
pub struct CallstackList {
    pub state: ListState,
    pub items: Vec<Callstack>,
}

impl BreakpointList {
    pub fn with_items(items: Vec<Breakpoint>) -> BreakpointList {
        BreakpointList {
            state: ListState::default(),
            selected: None,
            items,
        }
    }

    pub fn get_selected_filepath(&mut self) -> String {
        if self.items.len() == 0 {
            return "No data".to_string();
        }

        let selected = self.state.selected().unwrap();
        let data = &self.items[selected];
        return data.filepath.clone();
    }

    pub fn next(&mut self) {
        if self.items.len() == 0 {
            return;
        }

        let i = match self.state.selected() {
            Some(i) => {
                if i >= self.items.len() - 1 {
                    0
                } else {
                    i + 1
                }
            }
            None => 0,
        };
        self.state.select(Some(i));
    }

    pub fn previous(&mut self) {
        if self.items.len() == 0 {
            return;
        }

        let i = match self.state.selected() {
            Some(i) => {
                if i == 0 {
                    self.items.len() - 1
                } else {
                    i - 1
                }
            }
            None => 0,
        };
        self.state.select(Some(i));
    }

    pub fn get_selected_breakpoint(&mut self) -> Option<&Breakpoint> {
        if self.items.len() == 0 {
            return None;
        }

        let selected = self.state.selected().unwrap();
        let breakpoint = &self.items[selected];
        return Some(breakpoint);
    }

    pub fn get_value(&mut self) -> String {
        if self.items.len() == 0 {
            return "No variables".to_string();
        }

        let selected = self.state.selected().unwrap();
        let data = &self.items[selected];
        return data.payload.clone();
    }

    pub fn get_callstack(&mut self) -> Vec<Callstack> {
        if self.items.len() == 0 {
            return vec![];
        }

        let selected = self.state.selected();
        match selected {
            Some(selected_i) => {
                let data = &self.items[selected_i];
                return data.callstack.to_owned();
            }
            None => {
                return vec![];
            }
        }
    }

    pub fn unselect(&mut self) {
        self.state.select(None);
    }
}

impl CallstackList {
    pub fn with_items(items: Vec<Callstack>) -> CallstackList {
        CallstackList {
            state: ListState::default(),
            items,
        }
    }

    pub fn get_selected_callstack(&mut self) -> Option<&Callstack> {
        if self.items.len() == 0 {
            return None;
        }

        let selected = self.state.selected().unwrap();
        let callstack = &self.items[selected];
        return Some(callstack);
    }

    pub fn next(&mut self) {
        if self.items.len() == 0 {
            return;
        }

        let i = match self.state.selected() {
            Some(i) => {
                if i >= self.items.len() - 1 {
                    0
                } else {
                    i + 1
                }
            }
            None => 0,
        };
        self.state.select(Some(i));
    }

    pub fn previous(&mut self) {
        if self.items.len() == 0 {
            return;
        }

        let i = match self.state.selected() {
            Some(i) => {
                if i == 0 {
                    self.items.len() - 1
                } else {
                    i - 1
                }
            }
            None => 0,
        };
        self.state.select(Some(i));
    }

    pub fn unselect(&mut self) {
        self.state.select(None);
    }

    pub fn get_selected_filepath(&mut self) -> String {
        if self.items.len() == 0 {
            return "No data".to_string();
        }

        let selected = self.state.selected().unwrap();
        let data = &self.items[selected];
        return data.filepath.clone();
    }

    pub fn select_first(&mut self) {
        if self.items.len() == 0 {
            return;
        }

        self.state.select(Some(0));
    }
}

pub fn render_main(state: &mut UiState, f: &mut Frame<CrosstermBackend<Stdout>>) {
    let size = f.size();

    let main = Layout::default()
        .direction(Direction::Vertical)
        .constraints(
            [
                Constraint::Percentage(29),
                Constraint::Percentage(70),
                Constraint::Percentage(1),
            ]
            .as_ref(),
        )
        .split(size);

    let top_layout = Layout::default()
        .direction(Direction::Horizontal)
        .constraints([Constraint::Percentage(50), Constraint::Percentage(50)].as_ref())
        .split(main[0]);

    let bottom_layout = Layout::default()
        .direction(Direction::Horizontal)
        .constraints([Constraint::Percentage(100)].as_ref())
        .split(main[1]);

    let status_layout = Layout::default()
        .direction(Direction::Horizontal)
        .constraints([Constraint::Percentage(100)].as_ref())
        .split(main[2]);

    let breakpoints_layout = Block::default().title("Breakpoints").borders(Borders::ALL);
    f.render_widget(breakpoints_layout, top_layout[0]);

    let breakpoint_items: Vec<ListItem> = state
        .list_state
        .breakpoint
        .items
        .iter()
        .map(|i| {
            let filepath = &i.filepath;
            let line_num = &i.line;
            let title = format!("[{}] {}", line_num, filepath);
            let lines = vec![Spans::from(title)];
            ListItem::new(lines).style(Style::default().fg(Color::White))
        })
        .collect();

    let breakpoint_list = List::new(breakpoint_items)
        .block(Block::default().borders(Borders::ALL).title("Breakpoints"))
        .highlight_style(Style::default().bg(Color::LightGreen))
        .highlight_symbol(">> ");

    f.render_stateful_widget(
        breakpoint_list,
        top_layout[0],
        &mut state.list_state.breakpoint.state,
    );

    let callstack_layout = Block::default().title("Callstack").borders(Borders::ALL);
    f.render_widget(callstack_layout, top_layout[1]);

    let callstack_items: Vec<ListItem> = state
        .list_state
        .callstack
        .items
        .iter()
        .map(|i| {
            let filepath = &i.filepath;
            let line_num = &i.line;
            let title = format!("[{}] {}", line_num, filepath);
            let lines = vec![Spans::from(title)];
            ListItem::new(lines).style(Style::default().fg(Color::White))
        })
        .collect();

    let callstack_list = List::new(callstack_items)
        .block(Block::default().borders(Borders::ALL).title("Callstack"))
        .highlight_style(Style::default().bg(Color::LightGreen))
        .highlight_symbol(">> ");

    f.render_stateful_widget(
        callstack_list,
        top_layout[1],
        &mut state.list_state.callstack.state,
    );

    let vars = Paragraph::new(state.variables.clone())
        .block(Block::default().title("Variables").borders(Borders::ALL));
    f.render_widget(vars, bottom_layout[0]);

    let status_bar = Paragraph::new(state.status_bar.get_status().clone())
        .block(Block::default().borders(Borders::NONE));
    f.render_widget(status_bar, status_layout[0]);

    // popup element
    if state.popup.is_active() {
        let area = render_popup(100, 20, size);
        let text = state.popup.get_text();

        let mut paragraph_text = vec![];
        for i in text.into_iter() {
            paragraph_text.push(Spans::from(i));
        }

        let popup_block =
            Paragraph::new(paragraph_text).block(Block::default().borders(Borders::ALL));

        f.render_widget(Clear, area);
        f.render_widget(popup_block, area);
    }
}

/// helper function to create a centered rect using up certain percentage of the available rect `r`
pub fn render_popup(percent_x: u16, percent_y: u16, r: Rect) -> Rect {
    let popup_layout = Layout::default()
        .direction(Direction::Vertical)
        .constraints(
            [
                Constraint::Percentage((100 - percent_y) / 2),
                Constraint::Percentage(percent_y),
                Constraint::Percentage((100 - percent_y) / 2),
            ]
            .as_ref(),
        )
        .split(r);

    Layout::default()
        .direction(Direction::Horizontal)
        .constraints(
            [
                Constraint::Percentage((100 - percent_x) / 2),
                Constraint::Percentage(percent_x),
                Constraint::Percentage((100 - percent_x) / 2),
            ]
            .as_ref(),
        )
        .split(popup_layout[1])[1]
}
