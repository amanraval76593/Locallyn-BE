package elasticsearch

const feedPostsMapping = `{
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "user_id": {
        "type": "keyword"
      },
      "incident_id": {
        "type": "keyword"
      },
      "incident_title": {
        "type": "text"
      },
      "incident_category": {
        "type": "keyword"
      },
      "incident_trust_score": {
        "type": "float"
      },
      "content": {
        "type": "text"
      },
      "location": {
        "type": "geo_point"
      },
      "radius": {
        "type": "integer"
      },
      "identity_type": {
        "type": "keyword"
      },
      "post_type": {
        "type": "keyword"
      },
      "trust_score": {
        "type": "float"
      },
      "media_urls": {
        "type": "keyword"
      },
      "created_at": {
        "type": "date"
      },
      "expires_at": {
        "type": "date"
      },
      "is_deleted": {
        "type": "boolean"
      },
      "is_flagged": {
        "type": "boolean"
      }
    }
  }
}`
