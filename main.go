package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-xmlfmt/xmlfmt"
	"google.golang.org/api/healthcare/v1"
)

const (
	ProjectID             = "GCP_PROJECT"
	Location              = "GCLOUD_PROJECT_LOCATION"
	DatasetID             = "GCLOUD_PROJECT_DATASET_ID"
	SourceDicomStore      = "SOURCE_DICOM_STORE"
	DestinationDicomStore = "DESTINATION_DICOM_STORE"
)

// GoogleDicomAPI represents a healthcare implementation of dicom.DicomService
type GoogleDicomAPI struct {
	HealthcareService *healthcare.Service
	StoreService      *healthcare.ProjectsLocationsDatasetsDicomStoresService
	Dataset           *healthcare.Dataset
}

// NewGoogleDicomAPI returns a new instance of DicomAPI
func NewGoogleDicomAPI(ctx context.Context) (*GoogleDicomAPI, error) {

	p := os.Getenv(ProjectID)
	l := os.Getenv(Location)
	d := os.Getenv(DatasetID)

	datasetName := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", p, l, d)

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("healthcare.NewService: %v", err)
	}

	dicomStoreService := healthcareService.Projects.Locations.Datasets.DicomStores

	dicomAPI := &GoogleDicomAPI{
		HealthcareService: healthcareService,
		StoreService:      dicomStoreService,

		Dataset: &healthcare.Dataset{
			Name: datasetName,
		},
	}
	return dicomAPI, nil
}

// DicomStoreService represents a service for managing DicomStores
type DicomStoreService struct {
	GoogleDicomAPI *GoogleDicomAPI
}

// DicomService represents a service for managing Dicoms
type DicomService struct {
	googleDicomAPI *GoogleDicomAPI
}

// NewDicomService returns a new instance of DicomService
func NewDicomService(googleDicomAPI *GoogleDicomAPI) *DicomService {
	return &DicomService{
		googleDicomAPI: googleDicomAPI,
	}
}

// NewDicomStoreService returns a new instance of DicomStoreService
func NewDicomStoreService(googleDicomAPI *GoogleDicomAPI) *DicomStoreService {
	return &DicomStoreService{
		GoogleDicomAPI: googleDicomAPI,
	}
}

func (s *DicomStoreService) DeidentifyDicomStore(ctx context.Context) error {

	// Get env variables
	sourceDicomStore := os.Getenv(SourceDicomStore)
	destinationDicomStore := os.Getenv(DestinationDicomStore)

	datasetsService := s.GoogleDicomAPI.HealthcareService.Projects.Locations.Datasets.DicomStores

	req := &healthcare.DeidentifyDicomStoreRequest{
		DestinationStore: fmt.Sprintf("%s/dicomStores/%s", s.GoogleDicomAPI.Dataset.Name, destinationDicomStore),
		Config: &healthcare.DeidentifyConfig{
			Dicom: &healthcare.DicomConfig{
				FilterProfile: "MINIMAL_KEEP_LIST_PROFILE",
			},
			Image: &healthcare.ImageConfig{
				TextRedactionMode: "REDACT_SENSITIVE_TEXT",
			},
		},
	}

	sourceName := fmt.Sprintf("%s/dicomStores/%s", s.GoogleDicomAPI.Dataset.Name, sourceDicomStore)
	resp, err := datasetsService.Deidentify(sourceName, req).Do()
	if err != nil {
		return fmt.Errorf("Deidentify: %v", err)
	}

	// Wait for the deidentification operation to finish.
	operationService := s.GoogleDicomAPI.HealthcareService.Projects.Locations.Datasets.Operations
	for {
		op, err := operationService.Get(resp.Name).Do()
		if err != nil {
			return fmt.Errorf("operationService.Get: %v", err)
		}
		if !op.Done {
			time.Sleep(1 * time.Second)
			continue
		}
		if op.Error != nil {
			return fmt.Errorf("deidentify operation error: %v", *op.Error)
		}
		fmt.Printf("Created de-identified dataset %s from %s\n", resp.Name, sourceName)
		return nil
	}

}

// CreateDicomInstance creates dicom instances in the cloud within special abstractions called dicomStores
func (s *DicomService) CreateDicomInstance(ctx context.Context, dicomFilePath string) error {

	// must get env otherwise the operation fails is returned
	sourceDicomStore := os.Getenv(SourceDicomStore)
	if sourceDicomStore == "" {
		return fmt.Errorf("%v envar could not be found", SourceDicomStore)
	}

	dicomData, err := ioutil.ReadFile(dicomFilePath)

	if err != nil {
		return fmt.Errorf("ReadFile: %v", err)
	}

	parent := fmt.Sprintf("%s/dicomStores/%s", s.googleDicomAPI.Dataset.Name, sourceDicomStore)
	dicomWebPath := "studies"

	call := s.googleDicomAPI.StoreService.StoreInstances(parent, dicomWebPath, bytes.NewReader(dicomData))
	call.Header().Set("Content-Type", "application/dicom")
	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("StoreInstances: %v", err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("StoreInstances: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	x := xmlfmt.FormatXML(string(respBytes), "\t", "  ")
	print(x)

	return nil

}

func main() {

	ctx := context.Background()
	// initialise Google's Dicom API
	dicomApI, err := NewGoogleDicomAPI(ctx)
	if err != nil {
		fmt.Printf("unable to create DICOM Google API: %v \n", err)
		os.Exit(1)
	}

	// create a DICOM instance of the read dicom on the gcloud
	dicomService := NewDicomService(dicomApI)
	err = dicomService.CreateDicomInstance(ctx, "test-dicoms/case1_044.dcm")
	if err != nil {
		fmt.Printf("\nunable to Create a Dicom instance on gcloud: %v \n", err)
		os.Exit(1)
	}

	// initialise dicom store service and start the deidentification process
	dicomStoreService := NewDicomStoreService(dicomApI)
	err = dicomStoreService.DeidentifyDicomStore(ctx)
	if err != nil {
		fmt.Printf("\nunable to Deidentify DICOMs in DICOM Store: %v \n", err)
		os.Exit(1)
	}

	fmt.Printf("Dicom Successfully Deidentified")
}
