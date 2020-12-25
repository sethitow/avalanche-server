#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use]
extern crate rocket;
use rocket::http;
use rocket::response::content;
use serde_json;
use std::fs;

mod avalanche_data_model;

#[get("/")]
fn index() -> &'static str {
    "Hello, world!"
}

#[get("/forecast/<center>")]
fn hello(
    center: &http::RawStr,
    v: rocket::State<avalanche_data_model::Root>,
) -> Option<content::Json<String>> {
    println!("Looking up center {}.", center.as_str());
    let feature = v.get_feature_by_center_id(center.to_string());
    match feature {
        Some(x) => Some(content::Json(
            serde_json::to_string(&x.properties).expect("Could not seralize response"),
        )),
        None => None,
    }
}

fn main() {
    let filename = "avalanche_data.json";
    let data = fs::read_to_string(filename).expect("Error reading file");

    let v: avalanche_data_model::Root = serde_json::from_str(&data).expect("Error parsing file");

    println!("{}", v.obj_type);

    // let res = v.get_feature_by_center_id("SAC".to_string());

    // match res {
    //     Some(x) => println!("Got {:?}", x),
    //     None => println!("Got None"),
    // }

    rocket::ignite()
        .mount("/", routes![index, hello])
        .manage(v)
        .launch();
}
