package maxmind

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const geoLiteURLprefix = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key="
const geoLiteURLsuffix = "&suffix=zip"

const zipFileName = "db.zip"
const unzipDir = "/tmp/"
const csvFileNameIPv4 = "GeoLite2-Country-Blocks-IPv4.csv"
const csvFileNameIPv6 = "GeoLite2-Country-Blocks-IPv6.csv"
const csvFileNameCountry = "GeoLite2-Country-Locations-en.csv"

type GeoDB interface {
	GenerateDB() error
	IsWhitelisted(net.IP, []string) bool
}

type geoDB struct {
	geoLite2URL        string
	iPv4CSVList        []*GeoLite2CountryBlocksIPv4CSV
	iPv6CSVList        []*GeoLite2CountryBlocksIPv6CSV
	locationsEnCSVList []*GeoLite2CountryLocationsEnCSV
	countryCodeToGeo   map[string]int64
	geo2IPv4           map[int64][]net.IPNet
	geo2IPv6           map[int64][]net.IPNet
}

func (db *geoDB) IsWhitelisted(ip net.IP, countries []string) bool {
	for _, country := range countries {
		geonameID, ok := db.countryCodeToGeo[country]
		if !ok { // if goenameID is not found move to next
			continue
		}
		var whitelistedIPs []net.IPNet
		if isIPv6(ip) {
			whitelistedIPs = db.geo2IPv6[geonameID]

		} else {
			whitelistedIPs = db.geo2IPv4[geonameID]
		}
		for _, whitelistedIP := range whitelistedIPs {
			if whitelistedIP.Contains(ip) {
				return true
			}
		}
	}
	return false
}

// checks if the ip address is v6 or v4
func isIPv6(ip net.IP) bool {
	if strings.Contains(ip.String(), ":") {
		return true
	}
	return false
}

type NewParams struct {
	LicenseKey string
}

func New(p NewParams) GeoDB {
	return &geoDB{
		geoLite2URL: geoLiteURLprefix + p.LicenseKey + geoLiteURLsuffix,
	}
}

func (db *geoDB) GenerateDB() error {
	// Create zip file to write to
	out, err := os.Create("./" + zipFileName)
	if err != nil {
		return err
	}
	log.Println(out.Name())
	defer out.Close()

	log.Println(exec.Command("ls").Run())

	// Download zip file and write to file system
	resp, err := http.Get(db.geoLite2URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to fetch geoLite DB")
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	// Unzip zip file
	if err := unzip(zipFileName, unzipDir); err != nil {
		return err
	}

	return db.generateNetworkMap()
}

func (db *geoDB) generateNetworkMap() error {
	// Parse CSVs
	csvFile, err := os.Open(unzipDir + csvFileNameIPv4)
	if err != nil {
		return err
	}
	csvr := csv.NewReader(csvFile)
	db.iPv4CSVList, err = generateGeoLite2CountryBlocksIPv4List(csvr)
	if err != nil {
		return err
	}
	csvFile.Close()

	csvFile, err = os.Open(unzipDir + csvFileNameIPv6)
	if err != nil {
		return err
	}
	csvr = csv.NewReader(csvFile)
	db.iPv6CSVList, err = generateGeoLite2CountryBlocksIPv6CSVList(csvr)
	if err != nil {
		return err
	}
	csvFile.Close()

	csvFile, err = os.Open(unzipDir + csvFileNameIPv6)
	if err != nil {
		return err
	}
	csvr = csv.NewReader(csvFile)
	db.locationsEnCSVList, err = generateGeoLite2CountryLocationsEnCSVList(csvr)
	if err != nil {
		return err
	}
	csvFile.Close()

	// map the 2 character country code to its longer number
	db.countryCodeToGeo = map[string]int64{}
	// fill the map
	for _, c := range db.locationsEnCSVList {
		if c != nil {
			db.countryCodeToGeo[c.CountryISOCode] = c.GeonameID
		}
	}

	// Map the geocode to list of whitelisted ip masks
	db.geo2IPv4 = map[int64][]net.IPNet{}
	for _, geoIPv4 := range db.iPv4CSVList {
		if geoIPv4 != nil && geoIPv4.Network != nil {
			geoID := geoIPv4.GeonameID
			if db.geo2IPv4[geoID] == nil {
				db.geo2IPv4[geoID] = []net.IPNet{}
			}
			db.geo2IPv4[geoID] = append(db.geo2IPv4[geoID], *geoIPv4.Network)
		}
	}

	db.geo2IPv6 = map[int64][]net.IPNet{}
	for _, geoIPv6 := range db.iPv6CSVList {
		if geoIPv6 != nil && geoIPv6.Network != nil {
			geoID := geoIPv6.GeonameID
			if db.geo2IPv6[geoID] == nil {
				db.geo2IPv6[geoID] = []net.IPNet{}
			}
			db.geo2IPv6[geoID] = append(db.geo2IPv6[geoID], *geoIPv6.Network)
		}
	}

	return nil
}

// https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
