package polygon

// Order used for sort
type Order string

const Ascend Order = "asc"
const Decend Order = "desc"

// Timespan used for aggregation
type Timespan string

const Minute Timespan = "minute"
const Hour Timespan = "hour"
const Day Timespan = "day"
const Week Timespan = "week"
const Month Timespan = "month"
const Quarter Timespan = "quarter"
const Year Timespan = "year"

type MarketStatus string

const Open MarketStatus = "open"
const Closed MarketStatus = "closed"
const EarlyHours MarketStatus = "early_hours"
const AfterHours MarketStatus = "after_hours"
