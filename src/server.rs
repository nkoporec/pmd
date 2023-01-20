use std::{sync::{Arc, Mutex}, net::SocketAddr};

use axum::{Router, routing::post, extract::State, Json};
use serde::Deserialize;
use tokio::sync::broadcast::{Sender, self};

// Our shared state
#[derive(Debug)]
pub struct AppState {
    pub data: Mutex<Vec<Breakpoint>>,
    // Channel used to send messages from server to client.
    pub tx: broadcast::Sender<Breakpoint>,
}

#[derive(Debug, Deserialize, Clone, Eq, Hash, PartialEq)]
pub struct Breakpoint {
    pub filepath: String,
    pub line: String,
    pub connector_type: String,
    pub payload: String,
    pub timestamp: String,
    pub callstack: Vec<Callstack>,
}

#[derive(Debug, Deserialize, Clone, Eq, Hash, PartialEq)]
pub struct Callstack {
    pub filepath: String,
    pub line: String,
}


#[tokio::main(worker_threads = 1)]
#[allow(unused_must_use)]
pub async fn run(tx: Sender<Breakpoint>, port: u16) {
    tokio::spawn(async move {
        let app_state = Arc::new(AppState { data: Mutex::new(vec![]), tx });
        let app = Router::new().route("/dump", post(dump)).with_state(app_state);

        let addr = SocketAddr::from(([127, 0, 0, 1], port));
        axum::Server::bind(&addr)
            .serve(app.into_make_service())
            .await
            .unwrap();
    }).await;

    println!("Server has stopped");
}

async fn dump(State(state): State<Arc<AppState>>, data: Json<Breakpoint>) {
    let breakpoint = Breakpoint{
        filepath: data.filepath.clone(),
        line: data.line.clone(),
        connector_type: data.connector_type.clone(),
        payload: data.payload.clone(),
        timestamp: data.timestamp.clone(),
        callstack: data.callstack.clone()
    };
    state.tx.send(breakpoint).unwrap();
}
