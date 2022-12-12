use chrono::DateTime;
use serde::{Deserialize, Serialize};
use std::convert::TryFrom;

const QUERY_BASE_URL: &str =
    "https://api.avalanche.org/v2/public/products/map-layer/";

#[derive(Serialize, Deserialize, Debug)]
pub struct ForecastResponse {
    pub danger_level: i8,
    pub travel_advice: String,
    pub updated_at: i32,
    pub expires_at: i32,
    pub off_season: bool
}

#[derive(Serialize, Deserialize, Debug)]
struct Root {
    #[serde(rename(serialize = "type", deserialize = "type"))]
    pub object_type: String, 
    pub features: Vec<Feature>,
}

#[derive(Serialize, Deserialize, Debug)]
struct Feature {
    #[serde(rename(serialize = "type", deserialize = "type"))]
    pub object_type: String,
    pub properties: Properties,
}

#[derive(Serialize, Deserialize, Debug)]
struct Properties {
    pub danger_level: i8,
    pub travel_advice: String,
    pub start_date: String,
    pub end_date: String,
    pub off_season: bool
}

pub fn get_forecast_response(
    center_id: String,
) -> Result<ForecastResponse, Box<dyn std::error::Error>> {
    let query_url = String::from(QUERY_BASE_URL) + &center_id;

    let resp = reqwest::blocking::get(query_url)?.json::<Root>()?;
    println!("{:?}", resp);

    let feature = &resp.features[0];

    let parsable_start_date = feature.properties.start_date.clone() + "Z";
    let parsable_end_date = feature.properties.end_date.clone() + "Z";
    let updated_at = DateTime::parse_from_rfc3339(&parsable_start_date)?;
    let expires_at = DateTime::parse_from_rfc3339(&parsable_end_date)?;

    let r = ForecastResponse {
        danger_level: feature.properties.danger_level,
        travel_advice: feature.properties.travel_advice.clone(),
        updated_at: i32::try_from(updated_at.timestamp())?,
        expires_at: i32::try_from(expires_at.timestamp())?,
        off_season: feature.properties.off_season
    };
    return Ok(r);
}
