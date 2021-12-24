use chrono::DateTime;
use scraper::Html;
use serde::{Deserialize, Serialize};
use std::cmp;
use std::convert::TryFrom;

const QUERY_URL: &str =
    "https://api.avalanche.org/v2/public/product?type=forecast&center_id=SAC&zone_id=77";

#[derive(Serialize, Deserialize, Debug)]
pub struct ForecastResponse {
    pub danger_level: i8,
    pub upper_danger_level: i8,
    pub middle_danger_level: i8,
    pub lower_danger_level: i8,
    pub travel_advice: String,
    pub updated_at: i32,
    pub expires_at: i32,
}

#[derive(Serialize, Deserialize, Debug)]
struct Root {
    pub danger: Vec<DangerItem>,
    pub bottom_line: String,
    pub updated_at: String,
    pub expires_time: String,
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
    travel_advice = travel_advice.replace(|c: char| !c.is_ascii(), "");

    let high_danger = cmp::max(
        resp.danger[0].lower,
        cmp::max(resp.danger[0].middle, resp.danger[0].upper),
    );

    let updated_at = DateTime::parse_from_rfc3339(&resp.updated_at)?;
    let expires_at = DateTime::parse_from_rfc3339(&resp.expires_time)?;

    let r = ForecastResponse {
        danger_level: high_danger,
        upper_danger_level: resp.danger[0].upper,
        middle_danger_level: resp.danger[0].middle,
        lower_danger_level: resp.danger[0].lower,
        travel_advice: travel_advice,
        updated_at: i32::try_from(updated_at.timestamp())?,
        expires_at: i32::try_from(expires_at.timestamp())?,
    };
    return Ok(r);
}
