package metrics

import "time"

// Copied from https://github.com/arp242/goatcounter/blob/master/handlers/api.go

type APICountRequest struct {
	// By default it's an error to send pageviews that don't have either a
	// Session or UserAgent and IP set. This avoids accidental errors.
	//
	// When this is set it will just continue without recording sessions for
	// pageviews that don't have these parameters set.
	NoSessions bool `json:"no_sessions"`

	// Filter pageviews; accepted values:
	//
	//   ip     Ignore requests coming from IP addresses listed in "Settings → Ignore IP". Requires the IP field to be set.
	//
	// ["ip"] is used if this field isn't sent; send an empty array ([]) to not
	// filter anything.
	//
	// The X-Goatcounter-Filter header will be set to a list of indexes if any
	// pageviews are filtered; for example:
	//
	//    X-Goatcounter-Filter: 5, 10
	//
	// This header will be omitted if nothing is filtered.
	Filter []string `json:"filter"`

	// Hits is the list of pageviews.
	Hits []APICountRequestHit `json:"hits"`
}

type APICountRequestHit struct {
	// Path of the pageview, or the event name. {required}
	Path string `json:"path" query:"p"`

	// Page title, or some descriptive event title.
	Title string `json:"title" query:"t"`

	// Is this an event?
	Event bool `json:"event" query:"e"`

	// Referrer value, can be an URL (i.e. the Referal: header) or any
	// string.
	Ref string `json:"ref" query:"r"`

	// Screen size as "x,y,scaling"
	// Size goatcounter.Floats `json:"size" query:"s"`

	// Query parameters for this pageview, used to get campaign parameters.
	Query string `json:"query" query:"q"`

	// Hint if this should be considered a bot; should be one of the JSBot*`
	// constants from isbot; note the backend may override this if it
	// detects a bot using another method.
	// https://github.com/zgoat/isbot/blob/master/isbot.go#L28
	Bot int `json:"bot" query:"b"`

	// User-Agent header.
	UserAgent string `json:"user_agent"`

	// Location as ISO-3166-1 alpha2 string (e.g. NL, ID, etc.)
	Location string `json:"location"`

	// IP to get location from; not used if location is set. Also used for
	// session generation.
	IP string `json:"ip"`

	// Time this pageview should be recorded at; this can be in the past,
	// but not in the future.
	CreatedAt time.Time `json:"created_at"`

	// Normally a session is based on hash(User-Agent+IP+salt), but if you don't
	// send the IP address then we can't determine the session.
	//
	// In those cases, you can store your own session identifiers and send them
	// along. Note these will not be stored in the database as the sessionID
	// (just as the hashes aren't), they're just used as a unique grouping
	// identifier.
	Session string `json:"session"`
}
