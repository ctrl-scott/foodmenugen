// Cargo.toml
// [package] name = "menuapi" version = "0.1.0" edition = "2021"
// [dependencies]
// axum = "0.7"
// serde = { version = "1", features = ["derive"] }
// serde_json = "1"
// tokio = { version = "1", features = ["rt-multi-thread", "macros"] }
// askama = "0.12"
// parking_lot = "0.12"
// once_cell = "1"

use axum::{routing::post, Router, Json, http::StatusCode};
use once_cell::sync::Lazy;
use parking_lot::RwLock;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
struct Menu {
    id: String,
    restaurant: Restaurant,
    settings: Settings,
    sections: Vec<Section>,
}
#[derive(Debug, Clone, Serialize, Deserialize)]
struct Restaurant {
    name: String,
    address: String,
    phone: String,
    website: String,
    currency: String,
    note: String,
}
#[derive(Debug, Clone, Serialize, Deserialize)]
struct Settings { priceDecimals: u8, showTags: bool, showDescriptions: bool }
#[derive(Debug, Clone, Serialize, Deserialize)]
struct Section { id: String, title: String, note: String, items: Vec<Item> }
#[derive(Debug, Clone, Serialize, Deserialize)]
struct Item { id: String, name: String, description: String, price: f64, tags: Vec<String>, available: Option<bool> }

static DB: Lazy<RwLock<HashMap<String, Menu>>> = Lazy::new(|| RwLock::new(HashMap::new()));

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/api/menu", post(store_menu))
        .route("/api/render", post(render_menu));
    println!("listening on 0.0.0.0:8080");
    axum::Server::bind(&"0.0.0.0:8080".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn store_menu(Json(menu): Json<Menu>) -> StatusCode {
    DB.write().insert(menu.id.clone(), menu);
    StatusCode::NO_CONTENT
}

#[derive(Deserialize)]
struct RenderReq { id: String }

#[derive(askama::Template)]
#[template(source = r#"
<!doctype html><html><head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{ menu.restaurant.name }}</title>
<style>{{ css }}</style>
</head><body><main class="right"><section class="menu" id="preview">
<header>
  <h1 style="font-size:28px; margin:0;">{{ menu.restaurant.name }}</h1>
  <div class="meta">{{ menu.restaurant.address }}{% if menu.restaurant.phone %} · {{ menu.restaurant.phone }}{% endif %}{% if menu.restaurant.website %} · {{ menu.restaurant.website }}{% endif %}</div>
  {% if menu.restaurant.note %}<div class="meta">{{ menu.restaurant.note }}</div>{% endif %}
</header>
{% for sec in menu.sections %}
  <h2>{{ sec.title }}</h2>
  {% if sec.note %}<div class="meta">{{ sec.note }}</div>{% endif %}
  {% for it in sec.items %}
    {% if it.available == None || it.available == Some(true) %}
    <div class="item">
      <div>
        <div class="name">{{ it.name }}</div>
        {% if menu.settings.showDescriptions && it.description %}<div class="desc">{{ it.description }}</div>{% endif %}
        {% if menu.settings.showTags && it.tags.len() > 0 %}
          <div class="tags">{% for t in it.tags %}<span class="tag">{{ t }}</span>{% endfor %}</div>
        {% endif %}
      </div>
      <div class="price">{{ format!("{:.*}", menu.settings.priceDecimals as usize, it.price) }}</div>
    </div>
    {% endif %}
  {% endfor %}
{% endfor %}
</section></main></body></html>
"#, ext = "html")]
struct MenuPage<'a> { menu: &'a Menu, css: &'a str }

async fn render_menu(Json(req): Json<RenderReq>) -> Result<Json<serde_json::Value>, (StatusCode, String)> {
    let db = DB.read();
    let Some(menu) = db.get(&req.id) else { return Err((StatusCode::NOT_FOUND, "not found".into())); };

    let css = r#":root { --bg:#fff; --fg:#111; --muted:#666; --border:#e5e5e5;}
.right{padding:16px}.menu{max-width:900px;margin:0 auto}.menu h2{margin:24px 0 8px;border-bottom:2px solid #000;padding-bottom:6px;font-size:20px}
.meta{color:#666;font-size:12px;margin-bottom:16px}.item{display:grid;grid-template-columns:1fr auto;gap:8px;padding:8px 0;border-bottom:1px dashed #e5e5e5}
.item:last-child{border-bottom:0}.name{font-weight:600}.desc{color:#666;font-size:13px}.price{font-variant-numeric:tabular-nums;font-weight:600}
.tag{border:1px solid #e5e5e5;border-radius:999px;font-size:10px;padding:2px 6px;color:#666}"#;

    let page = MenuPage { menu, css };
    let html = page.render().map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))?;
    Ok(Json(serde_json::json!({ "html": html })))
}
