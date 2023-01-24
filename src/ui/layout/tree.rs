use serde_json::Value;
use tui::{
    style::{Color, Modifier, Style},
    widgets::{Block, Borders},
};
use tui_tree_widget::{Tree, TreeItem, TreeState};

pub fn render_tree(tree_state: TreeState, items: Vec<TreeItem>) -> Tree {
    let items = Tree::new(items)
        .block(
            Block::default()
                .borders(Borders::ALL)
                .title(format!("Tree Widget {:?}", tree_state)),
        )
        .highlight_style(
            Style::default()
                .fg(Color::Black)
                .bg(Color::LightGreen)
                .add_modifier(Modifier::BOLD),
        )
        .highlight_symbol(">> ");

    return items;
}

pub fn build_tree_items(payload: String) -> Vec<TreeItem<'static>> {
    let mut items = vec![];

    if payload == "" {
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
                items.push(TreeItem::new(
                    key.to_string(),
                    vec![TreeItem::new_leaf(value.to_string())],
                ));
            }

            if value.is_object() {
                let recursive_tree = recursive(value, key.to_string());
                items.push(recursive_tree);
            }
        }
    }

    return items;
}

fn recursive(value: &Value, key: String) -> TreeItem<'static> {
    let mut result = TreeItem::new_leaf(key);
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
            recursive(value, key.to_string());
        }
    }

    return result;
}
