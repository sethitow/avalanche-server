use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone)]
pub struct Root {
    #[serde(rename = "type")]
    pub obj_type: String,
    pub features: Vec<Feature>,
}

impl Root {
    pub fn get_feature_by_center_id(&self, center_id: String) -> Option<&Feature> {
        for f in &self.features {
            if f.properties.center_id == center_id {
                return Some(f);
            }
        }
        return None;
    }
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Feature {
    #[serde(rename = "type")]
    pub obj_type: String,
    pub id: i64,
    pub properties: FeatureProperties,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct FeatureProperties {
    pub name: String,
    pub center: String,
    pub center_id: String,
    pub danger: String,
    pub danger_level: i8,
    pub travel_advice: String,
}
