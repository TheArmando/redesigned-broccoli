package maxmind

// CSV imports were pulled from // https://stackoverflow.com/questions/24999079/reading-csv-file-in-go

import (
	"encoding/csv"
	"errors"
	"io"
	"net"
	"strconv"
)

type GeoLite2CountryBlocksIPv4CSV struct {
	Network                     *net.IPNet
	GeonameID                   int64
	RegisteredCountryGeonameID  int64
	RepresentedCountryGeonameID int64
	IsAnonymousProxy            bool
	IsSatelliteProvider         bool
}

// TODO: Dry these csv loading functions into a single multi-use one
func generateGeoLite2CountryBlocksIPv4List(csvr *csv.Reader) ([]*GeoLite2CountryBlocksIPv4CSV, error) {
	s := []*GeoLite2CountryBlocksIPv4CSV{}
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return s, err
		}

		val := &GeoLite2CountryBlocksIPv4CSV{}
		if _, val.Network, err = net.ParseCIDR(row[0]); err != nil {
			return nil, err
		}

		if val.GeonameID, err = strconv.ParseInt(row[1], 10, 64); err != nil {
			return nil, err
		}

		if val.RegisteredCountryGeonameID, err = strconv.ParseInt(row[2], 10, 64); err != nil {
			return nil, err
		}

		if val.RepresentedCountryGeonameID, err = strconv.ParseInt(row[3], 10, 64); err != nil {
			return nil, err
		}

		if val.IsAnonymousProxy, err = strconv.ParseBool(row[4]); err != nil {
			return nil, err
		}

		if val.IsSatelliteProvider, err = strconv.ParseBool(row[5]); err != nil {
			return nil, err
		}

		s = append(s, val)
	}
	return nil, errors.New("failed to parse, csv was empty")
}

type GeoLite2CountryBlocksIPv6CSV struct {
	Network                     *net.IPNet
	GeonameID                   int64
	RegisteredCountryGeonameID  int64
	RepresentedCountryGeonameID int64
	IsAnonymousProxy            bool
	IsSatelliteProvider         bool
}

func generateGeoLite2CountryBlocksIPv6CSVList(csvr *csv.Reader) ([]*GeoLite2CountryBlocksIPv6CSV, error) {
	s := []*GeoLite2CountryBlocksIPv6CSV{}
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return s, err
		}

		val := &GeoLite2CountryBlocksIPv6CSV{}
		if _, val.Network, err = net.ParseCIDR(row[0]); err != nil {
			return nil, err
		}

		if val.GeonameID, err = strconv.ParseInt(row[1], 10, 64); err != nil {
			return nil, err
		}

		if val.RegisteredCountryGeonameID, err = strconv.ParseInt(row[2], 10, 64); err != nil {
			return nil, err
		}

		if val.RepresentedCountryGeonameID, err = strconv.ParseInt(row[3], 10, 64); err != nil {
			return nil, err
		}

		if val.IsAnonymousProxy, err = strconv.ParseBool(row[4]); err != nil {
			return nil, err
		}

		if val.IsSatelliteProvider, err = strconv.ParseBool(row[5]); err != nil {
			return nil, err
		}

		s = append(s, val)
	}
	return nil, errors.New("failed to parse, csv was empty")
}

type GeoLite2CountryLocationsEnCSV struct {
	GeonameID         int64
	LocaleCode        string
	ContinentCode     string
	ContinentName     string
	CountryISOCode    string
	CountryName       string
	IsInEuropeanUnion bool
}

func generateGeoLite2CountryLocationsEnCSVList(csvr *csv.Reader) ([]*GeoLite2CountryLocationsEnCSV, error) {
	s := []*GeoLite2CountryLocationsEnCSV{}
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return s, err
		}

		val := &GeoLite2CountryLocationsEnCSV{}
		if val.GeonameID, err = strconv.ParseInt(row[0], 10, 64); err != nil {
			return nil, err
		}

		val.LocaleCode = row[1]
		val.ContinentCode = row[2]
		val.ContinentName = row[3]
		val.CountryISOCode = row[4]
		val.CountryName = row[5]

		if val.IsInEuropeanUnion, err = strconv.ParseBool(row[6]); err != nil {
			return nil, err
		}

		s = append(s, val)
	}
	return nil, errors.New("failed to parse, csv was empty")
}
