package main

import (
	"context"
	"testing"
)

func TestDicomService_CreateDicomInstances(t *testing.T) {
	type fields struct {
		googleDicomAPI *GoogleDicomAPI
	}
	type args struct {
		ctx           context.Context
		dicomFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DicomService{
				googleDicomAPI: tt.fields.googleDicomAPI,
			}
			if err := s.CreateDicomInstances(tt.args.ctx, tt.args.dicomFilePath); (err != nil) != tt.wantErr {
				t.Errorf("DicomService.CreateDicomInstances() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDicomStoreService_DeidentifyDicomStore(t *testing.T) {
	type fields struct {
		GoogleDicomAPI *GoogleDicomAPI
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DicomStoreService{
				GoogleDicomAPI: tt.fields.GoogleDicomAPI,
			}
			if err := s.DeidentifyDicomStore(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("DicomStoreService.DeidentifyDicomStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
