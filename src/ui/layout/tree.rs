use serde_json::Value;
use tui::{
    style::{Color, Modifier, Style},
    widgets::{Block, Borders},
};
use tui_tree_widget::{Tree, TreeItem};

pub fn render_tree(items: Vec<TreeItem>) -> Tree {
    let items = Tree::new(items)
        .block(
            Block::default()
                .borders(Borders::ALL)
                .title("Inspection".to_string()),
        )
        .highlight_style(
            Style::default()
                .fg(Color::Black)
                .bg(Color::LightGreen)
                .add_modifier(Modifier::BOLD),
        )
        .highlight_symbol(">> ");

    items
}

pub fn build_tree_items(payload: String) -> Vec<TreeItem<'static>> {
    let mut items = vec![];

    if payload.is_empty() || payload == "No variables" {
        return items;
    }

    let data: Value = serde_json::from_str(&payload).unwrap();
    if data.is_string() {
        items.push(TreeItem::new_leaf(data.to_string()));
    }

    if data.is_object() {
        let obj = data.as_object().unwrap();
        for item in obj {
            let key = item.0;
            let value = item.1;

            // Strings and numbers.
            if value.is_string() || value.is_number() {
                items.push(TreeItem::new_leaf(value.to_string()))
            }

            if value.is_object() {
                let leaf = TreeItem::new_leaf(key.to_string());
                let recursive_tree = recursive_flatten(value, &mut leaf.clone());
                items.push(recursive_tree);
            }
        }
    }

    items
}

// Recursive flatten an object, so we can extract only strings and numbers.
fn recursive_flatten<'a>(value: &Value, result: &mut TreeItem<'a>) -> TreeItem<'a> {
    let obj = value.as_object().unwrap();

    for item in obj {
        let key = item.0;
        let value = item.1;

        // Strings and numbers.
        if value.is_string() || value.is_number() {
            result.add_child(TreeItem::new(
                key.to_string(),
                vec![TreeItem::new_leaf(value.to_string())],
            ));
        }

        if value.is_object() {
            recursive_flatten(value, result);
        }
    }

    result.clone()
}
