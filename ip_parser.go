package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"os"
	"net"
	"github.com/oschwald/geoip2-golang"
	"io"
	"compress/gzip"
	"net/http"
	"time"
)

const (
	CITY_FILE_URL        = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"
	LOCAL_CITY_FILE_PATH = "Keep Your Local System Path Here"
	CITY_FILE_PREFIX     = "GeoLite2-City"
	GZ_EXT               = ".gz"
	MMDB_EXT             = ".mmdb"
)

// our main function
func main() {
	fmt.Println("Starting Application")

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Next()
	})

	// Download mmdb file
	fmt.Println("Download Started for MMDB File")
	if err := downloadFile(CITY_FILE_URL, LOCAL_CITY_FILE_PATH+CITY_FILE_PREFIX+MMDB_EXT+GZ_EXT); err != nil {
		fmt.Println("Error While Downloading MMDB File")
	}

	// Unzip mmdb file
	if err := copyMmdbFileFromGzip(LOCAL_CITY_FILE_PATH+CITY_FILE_PREFIX+MMDB_EXT+GZ_EXT, LOCAL_CITY_FILE_PATH+CITY_FILE_PREFIX+MMDB_EXT); err != nil {
		fmt.Println("Error While Extracting MMDB File")
	}
	fmt.Println("Download Completed for MMDB File")

	// Create mmdb instance
	mmdbReader, err := NewMMDb(LOCAL_CITY_FILE_PATH+CITY_FILE_PREFIX+MMDB_EXT)
	if err !=nil {
		fmt.Println("Unable to create mmdb instance")
	}
	defer mmdbReader.Close()
	router.GET("/test", mmdbReader.ResolveIp)
	router.Run(":31001")

}


type MMDb struct {
	cityReader           *geoip2.Reader
	cityFile              string
}

// This returns pointer to a MMDb instance or error if any
func NewMMDb(cityFilePath string) (*MMDb, error) {
	var cityReader *geoip2.Reader
	var err error

	if cityReader, err = geoip2.Open(cityFilePath); err != nil {
		return nil, err
	}

	return &MMDb{
		cityReader:           cityReader,
		cityFile:              cityFilePath,
	}, nil
}


func (mmdb *MMDb) getLocationInfo(ip net.IP) (*geoip2.City, error) {
	return mmdb.cityReader.City(ip)
}


// This function extracts required values from ip
func (mmdb *MMDb) ResolveIp(c *gin.Context) {
	t := time.Now()
	var city *geoip2.City
	var err error
	var ip = net.ParseIP(c.ClientIP())

	if city, err = mmdb.getLocationInfo(ip); err != nil {
		fmt.Println("Error while getting location info:"+err.Error())
	} else {
		fmt.Println("Country: "+city.Country.Names["en"])
	}

	latency := time.Since(t)
	fmt.Println("Latency for above request: "+latency.String())
}


// This function downloads file from given url to given file path
func downloadFile(url, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error while creating file at: "+filepath)
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading file from url: "+url)
		return err
	}

	_, err = io.Copy(f, res.Body)
	defer res.Body.Close()

	if err != nil {
		fmt.Println("Error while copying response to filepath: ", filepath)
		return err
	}

	defer f.Close()
	return nil
}


// This function copies gzip archive content to output file
func copyMmdbFileFromGzip(gzipFile, outputFile string) error {
	var out, f *os.File
	var gzReader io.Reader
	var err error

	if out, err = os.Create(outputFile); err != nil {
		fmt.Println("Error while creating output file: "+outputFile)
		return err
	}
	defer out.Close()

	if f, err = os.Open(gzipFile); err != nil {
		fmt.Println("Error while opening gzip file: "+gzipFile)
		return err
	}
	defer f.Close()
	if gzReader, err = gzip.NewReader(f); err != nil {
		fmt.Println("Error while reading gzip file: "+gzipFile)
		return err
	}

	io.Copy(out, gzReader)
	return nil
}


// This function closes all readers
func (mmdb *MMDb) Close() {
	mmdb.cityReader.Close()
	fmt.Println("Successfully closed MMDb reader")
}



