package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func fetchAndParseFeed(url string) (*gtfs.FeedMessage, error) {
	fmt.Printf("Fetching GTFS feed from: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GTFS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch GTFS feed: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("failed to parse protobuf: %w", err)
	}

	return feed, nil
}

func writeToFile(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filename, err)
	}

	return nil
}

func formatFeedMessage(feed *gtfs.FeedMessage) string {
	var result string

	// Add header information
	result += fmt.Sprintf("GTFS-Realtime Feed Message\n")
	result += fmt.Sprintf("==========================\n")
	if feed.Header != nil {
		result += fmt.Sprintf("Header:\n")
		result += fmt.Sprintf("  Version: %s\n", feed.Header.GetGtfsRealtimeVersion())
		result += fmt.Sprintf("  Timestamp: %d\n", feed.Header.GetTimestamp())
		result += fmt.Sprintf("  Incrementality: %s\n", feed.Header.GetIncrementality().String())
		result += fmt.Sprintf("\n")
	}

	// Add entity information
	result += fmt.Sprintf("Entities (%d total):\n", len(feed.Entity))
	result += fmt.Sprintf("===================\n\n")

	for i, entity := range feed.Entity {
		result += fmt.Sprintf("Entity %d:\n", i+1)
		result += fmt.Sprintf("  ID: %s\n", entity.GetId())
		result += fmt.Sprintf("  IsDeleted: %t\n", entity.GetIsDeleted())

		// Check what type of entity this is
		if entity.TripUpdate != nil {
			result += fmt.Sprintf("  Type: Trip Update\n")
			result += fmt.Sprintf("  Trip Update: %s\n", entity.TripUpdate.String())
		}

		if entity.Vehicle != nil {
			result += fmt.Sprintf("  Type: Vehicle Position\n")
			result += fmt.Sprintf("  Vehicle Position: %s\n", entity.Vehicle.String())
		}

		if entity.Alert != nil {
			result += fmt.Sprintf("  Type: Alert\n")
			result += fmt.Sprintf("  Alert: %s\n", entity.Alert.String())
		}

		result += fmt.Sprintf("\n")
		result += fmt.Sprintf("---\n\n")
	}

	return result
}

func main() {
	// Get current timestamp for filenames
	timestamp := time.Now().Format("20060102_150405")

	// URLs
	tripUpdatesURL := "https://api.pugetsound.onebusaway.org/api/gtfs_realtime/trip-updates-for-agency/1.pb?key=org.onebusaway.iphone"
	vehiclePositionsURL := "https://api.pugetsound.onebusaway.org/api/gtfs_realtime/vehicle-positions-for-agency/1.pb?key=org.onebusaway.iphone"
	alertsURL := "https://api.pugetsound.onebusaway.org/api/gtfs_realtime/alerts-for-agency/1.pb?key=org.onebusaway.iphone"

	// Create output directory
	outputDir := "./gtfsRT-Analysis/gtfs_outputs"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Fetch and parse feeds
	tripUpdates, err := fetchAndParseFeed(tripUpdatesURL)
	if err != nil {
		log.Fatalf("Failed to fetch trip updates: %v", err)
	}

	vehiclePositions, err := fetchAndParseFeed(vehiclePositionsURL)
	if err != nil {
		log.Fatalf("Failed to fetch vehicle positions: %v", err)
	}

	alerts, err := fetchAndParseFeed(alertsURL)
	if err != nil {
		log.Fatalf("Failed to fetch alerts: %v", err)
	}

	// Write decoded (human-readable) text files
	tripUpdatesFile := filepath.Join(outputDir, fmt.Sprintf("trip_updates_%s.txt", timestamp))
	if err := writeToFile(tripUpdatesFile, formatFeedMessage(tripUpdates)); err != nil {
		log.Fatalf("Failed to write trip updates file: %v", err)
	}

	vehiclePositionsFile := filepath.Join(outputDir, fmt.Sprintf("vehicle_positions_%s.txt", timestamp))
	if err := writeToFile(vehiclePositionsFile, formatFeedMessage(vehiclePositions)); err != nil {
		log.Fatalf("Failed to write vehicle positions file: %v", err)
	}

	alertsFile := filepath.Join(outputDir, fmt.Sprintf("alerts_%s.txt", timestamp))
	if err := writeToFile(alertsFile, formatFeedMessage(alerts)); err != nil {
		log.Fatalf("Failed to write alerts file: %v", err)
	}

	fmt.Println("Successfully fetched and saved GTFS real-time feeds:")
	fmt.Printf("- Trip Updates: %s\n", tripUpdatesFile)
	fmt.Printf("- Vehicle Positions: %s\n", vehiclePositionsFile)
	fmt.Printf("- Alerts: %s\n", alertsFile)
}
