use scraper::Html;
use serde::{Deserialize, Serialize};
use std::cmp;

const QUERY_URL: &str =
    "https://api.avalanche.org/v2/public/product?type=forecast&center_id=SAC&zone_id=77";

#[derive(Serialize, Deserialize, Debug)]
pub struct ForecastResponse {
    pub danger_level: i8,
    pub travel_advice: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct Root {
    pub danger: Vec<DangerItem>,
    pub bottom_line: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct DangerItem {
    pub lower: i8,
    pub upper: i8,
    pub middle: i8,
    pub valid_day: String,
}

pub fn get_forecast_response(
    center_id: String,
) -> Result<ForecastResponse, Box<dyn std::error::Error>> {
    if center_id != "SAC" {
        return Err("Unknown Avalanche Center".into());
    }

    let resp = reqwest::blocking::get(QUERY_URL)?.json::<Root>()?;
    println!("{:?}", resp);

    let bottom_line = Html::parse_fragment(resp.bottom_line.as_str());
    let mut travel_advice = String::new();
    for elem in bottom_line.root_element().text() {
        travel_advice.push_str(elem)
    }
    travel_advice = travel_advice
        .to_string()
        .replace(|c: char| !c.is_ascii(), "");

    let high_danger = cmp::max(
        resp.danger[0].lower,
        cmp::max(resp.danger[0].middle, resp.danger[0].upper),
    );

    let r = ForecastResponse {
        danger_level: high_danger,
        travel_advice: travel_advice,
    };
    return Ok(r);
}
