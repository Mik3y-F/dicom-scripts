# Google Healthcare DICOM Deidentification

This script handles uploading DICOM files to Google Healthcare's API and subsequently deidentifying the uploaded data. DICOM (Digital Imaging and Communications in Medicine) is a standard for storing and transmitting medical images.

## Getting Started

To run this project, you will need to have the Go programming language installed on your system. You can download Go from the official website: [https://golang.org/dl/](https://golang.org/dl/)

After that, you can clone this repository using the following command:

```bash
git clone https://github.com/Mik3y-F/dicom-scripts.git
```

## Setting Up Environment Variables

The script requires several environment variables to be set:

- `GCP_PROJECT`: Your Google Cloud Project ID.
- `GCLOUD_PROJECT_LOCATION`: The location of your Google Cloud Project.
- `GCLOUD_PROJECT_DATASET_ID`: The ID of your dataset within your Google Cloud Project.
- `SOURCE_DICOM_STORE`: The name of the DICOM store that holds the original DICOM instances.
- `DESTINATION_DICOM_STORE`: The name of the DICOM store where the deidentified DICOM instances will be placed.

You can set these variables in your `shell`, or you can use a `.env` file if you prefer.

## Running the Script

To run this script, navigate to the directory containing the script and execute the following commands:

```bash
go get -u ./...
```

```bash
go run main.go
```

## Understanding the Script

This script has three primary functions:

- `NewGoogleDicomAPI()`: Initializes a new instance of Google's Healthcare DICOM API.
- `CreateDicomInstance()`: Uploads a DICOM file from a specified path to the source DICOM store in Google Cloud.
- `DeidentifyDicomStore()`: Deidentifies the DICOM instances in the source DICOM store and moves them to the destination DICOM store.

## License

This project is is licensed under the MIT license. See [LICENSE](LICENSE) for the full license text.
