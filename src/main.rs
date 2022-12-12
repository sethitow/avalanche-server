#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use]
extern crate rocket;
use rocket::http;
use rocket::response::content;
use serde_json;

mod avalanche_integration;

#[get("/")]
fn index() -> &'static str {
    "usage: /forecast/<center>"
}

#[get("/forecast/<center>")]
fn forecast(center: &http::RawStr) -> Option<content::Json<String>> {
    println!("Looking up center {}.", center.as_str());
    let feature = avalanche_integration::get_forecast_response(center.to_string());
    match feature {
        Ok(x) => Some(content::Json(
            serde_json::to_string(&x).expect("Could not seralize response"),
        )),
        Err(x) => {
            println!("{:?}", x);
            None
        }
    }
}

fn main() {
    rocket::ignite()
        .mount("/", routes![index, forecast])
        .launch();
}
