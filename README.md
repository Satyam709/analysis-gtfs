# GTFS Real-time Fetcher (Go Version)

This Go program fetches GTFS real-time feeds from the Puget Sound OneBusAway API and saves them as human-readable text files.

## Features

- Fetches three types of GTFS real-time feeds:
  - Trip Updates
  - Vehicle Positions  
  - Alerts
- Parses protobuf format and converts to readable text
- Shows repeated entities clearly with detailed formatting
- Saves files with timestamps for easy tracking

## Prerequisites

- Go 1.21 or higher
- Internet connection to fetch feeds

## Installation

1. Clone or download this project
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

Run the program:
```bash
go run main.go
```

This will:
1. Create a `gtfsRT-Analysis/gtfs_outputs` directory
2. Fetch the three real-time feeds
3. Save them as timestamped text files:
   - `trip_updates_YYYYMMDD_HHMMSS.txt`
   - `vehicle_positions_YYYYMMDD_HHMMSS.txt`
   - `alerts_YYYYMMDD_HHMMSS.txt`

## Output Format

The output files show:
- Feed header information (version, timestamp, incrementality)
- Total number of entities
- Each entity with:
  - Entity ID
  - Type (Trip Update, Vehicle Position, or Alert)
  - Detailed entity data
  - Clear separation between entities

## Dependencies

- `google.golang.org/protobuf` - Protocol buffer support
- `github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs` - GTFS real-time Go bindings

## Example Output

```
GTFS-Realtime Feed Message
==========================
Header:
  Version: 2.0
  Timestamp: 1704067200
  Incrementality: FULL_DATASET

Entities (150 total):
===================

Entity 1:
  ID: 1_12345678
  IsDeleted: false
  Type: Trip Update
  Trip Update: [detailed trip update information]

---

Entity 2:
  ID: 1_23456789
  IsDeleted: false
  Type: Vehicle Position
  Vehicle Position: [detailed vehicle position information]

---
```

## Comparison with Python Version

This Go version provides the same functionality as the Python version but with:
- Better performance and lower memory usage
- Static typing for better reliability
- Clear entity formatting showing repeated fields
- Cross-platform single binary distribution 