import pandas as pd

# Load GTFS core files
calendar = pd.read_csv("./data/calendar.txt")
calendar_dates = pd.read_csv("./data/calendar_dates.txt")
trips = pd.read_csv("./data/trips.txt")
stop_times = pd.read_csv("./data/stop_times.txt")
stops = pd.read_csv("./data/stops.txt")
routes = pd.read_csv("./data/routes.txt")
agency = pd.read_csv("./data/agency.txt")

print("Data loaded successfully.")

# Merge trips with calendar and calendar_dates
trips = trips.merge(calendar, on="service_id", how="left")
calendar_dates['type'] = calendar_dates['exception_type'].map({1: 'added', 2: 'removed'})

# Merge trips with stop_times to get stop_id and arrival time
trip_stops = trips.merge(stop_times, on="trip_id", how="inner")

# Merge with stops to get stop_name
trip_stops = trip_stops.merge(stops[['stop_id', 'stop_name']], on='stop_id', how='left')

# Merge with routes to get route_name and agency_id
trip_stops = trip_stops.merge(routes[['route_id', 'route_short_name', 'agency_id']], on='route_id', how='left')

# Merge with agency to get agency_name
trip_stops = trip_stops.merge(agency[['agency_id', 'agency_name']], on='agency_id', how='left')

# Extract weekday-only, weekend-only, and holiday-modified services
weekday_only = trip_stops[(trip_stops['monday'] == 1) & (trip_stops['saturday'] == 0) & (trip_stops['sunday'] == 0)]
weekend_only = trip_stops[(trip_stops['saturday'] == 1) & (trip_stops['monday'] == 0) & (trip_stops['sunday'] == 1)]
holiday_services = calendar_dates[calendar_dates['type'] == 'added'].merge(trips, on='service_id')

# Merge holiday_services with stop_times and stops
holiday_services = holiday_services.merge(stop_times, on='trip_id')
holiday_services = holiday_services.merge(stops[['stop_id', 'stop_name']], on='stop_id')
holiday_services = holiday_services.merge(routes[['route_id', 'route_short_name', 'agency_id']], on='route_id')
holiday_services = holiday_services.merge(agency[['agency_id', 'agency_name']], on='agency_id')
holiday_services['type'] = 'holiday'

# Add "type" column to distinguish trip types
weekday_only['type'] = 'weekday_only'
weekend_only['type'] = 'weekend_only'

# Select useful columns
useful_columns = [
    'trip_id', 'stop_id', 'stop_name', 'arrival_time', 'departure_time',
    'route_id', 'route_short_name', 'agency_name', 'service_id', 'type'
]

# Combine all testable cases
testable = pd.concat([
    weekday_only[useful_columns],
    weekend_only[useful_columns],
    holiday_services[useful_columns]
])

# Sort by stop_name and type for easier testing
testable = testable.sort_values(['stop_name', 'type'])

# Show a sample
print(testable.head(20))

# Optionally, save to CSV for inspection
testable.to_csv("testable_arrivals_by_day.csv", index=False)
