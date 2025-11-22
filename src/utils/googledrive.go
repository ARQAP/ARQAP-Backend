package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	driveService *drive.Service
	driveOnce    sync.Once
)

// InitGoogleDriveService inicializa el servicio de Google Drive usando Service Account
func InitGoogleDriveService() error {
	var initErr error
	driveOnce.Do(func() {
		// Obtener la ruta del archivo de credenciales desde variable de entorno
		credentialsPath := os.Getenv("GOOGLE_DRIVE_CREDENTIALS_PATH")
		if credentialsPath == "" {
			// Intentar usar el contenido JSON directamente desde variable de entorno
			credentialsJSON := os.Getenv("GOOGLE_DRIVE_CREDENTIALS_JSON")
			if credentialsJSON == "" {
				initErr = fmt.Errorf("GOOGLE_DRIVE_CREDENTIALS_PATH o GOOGLE_DRIVE_CREDENTIALS_JSON debe estar configurado")
				return
			}

			// Usar credenciales desde JSON string
			ctx := context.Background()
			creds, err := google.CredentialsFromJSON(ctx, []byte(credentialsJSON), drive.DriveReadonlyScope)
			if err != nil {
				initErr = fmt.Errorf("error cargando credenciales desde JSON: %w", err)
				return
			}

			driveService, err = drive.NewService(ctx, option.WithCredentials(creds))
			if err != nil {
				initErr = fmt.Errorf("error creando servicio de Google Drive: %w", err)
				return
			}
		} else {
			// Usar archivo de credenciales
			ctx := context.Background()
			credsBytes, readErr := os.ReadFile(credentialsPath)
			if readErr != nil {
				initErr = fmt.Errorf("error leyendo archivo de credenciales: %w", readErr)
				return
			}
			creds, err := google.CredentialsFromJSON(ctx, credsBytes, drive.DriveReadonlyScope)
			if err != nil {
				initErr = fmt.Errorf("error cargando credenciales: %w", err)
				return
			}

			driveService, err = drive.NewService(ctx, option.WithCredentials(creds))
			if err != nil {
				initErr = fmt.Errorf("error creando servicio de Google Drive: %w", err)
				return
			}
		}

		log.Printf("[GOOGLE_DRIVE] Servicio inicializado correctamente")
	})
	return initErr
}

// GetGoogleDriveService retorna el servicio de Google Drive (inicializa si es necesario)
func GetGoogleDriveService() (*drive.Service, error) {
	if driveService == nil {
		if err := InitGoogleDriveService(); err != nil {
			return nil, err
		}
	}
	return driveService, nil
}

// ExtractFileIDFromURL extrae el ID del archivo de una URL de Google Drive
func ExtractFileIDFromURL(url string) (string, error) {
	// Patrones comunes de URLs de Google Drive
	patterns := []string{
		`/file/d/([a-zA-Z0-9_-]+)`,                     // /file/d/FILE_ID
		`id=([a-zA-Z0-9_-]+)`,                          // ?id=FILE_ID
		`/folders/([a-zA-Z0-9_-]+)`,                    // /folders/FOLDER_ID
		`drive\.google\.com/open\?id=([a-zA-Z0-9_-]+)`, // open?id=FILE_ID
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no se pudo extraer el ID del archivo de la URL: %s", url)
}

// DownloadFileFromGoogleDrive descarga un archivo de Google Drive usando la API
func DownloadFileFromGoogleDrive(fileID string) (io.ReadCloser, string, error) {
	service, err := GetGoogleDriveService()
	if err != nil {
		return nil, "", fmt.Errorf("error obteniendo servicio de Google Drive: %w", err)
	}

	log.Printf("[GOOGLE_DRIVE] Descargando archivo con ID: %s", fileID)

	// Obtener información del archivo
	file, err := service.Files.Get(fileID).Fields("id", "name", "mimeType", "size").Do()
	if err != nil {
		return nil, "", fmt.Errorf("error obteniendo información del archivo: %w", err)
	}

	log.Printf("[GOOGLE_DRIVE] Archivo encontrado: %s (tipo: %s, tamaño: %d bytes)", file.Name, file.MimeType, file.Size)

	// Verificar si es una carpeta
	if file.MimeType == "application/vnd.google-apps.folder" {
		return nil, "", fmt.Errorf("las carpetas de Google Drive no se pueden descargar directamente")
	}

	// Descargar el archivo
	resp, err := service.Files.Get(fileID).Download()
	if err != nil {
		return nil, "", fmt.Errorf("error descargando archivo: %w", err)
	}

	log.Printf("[GOOGLE_DRIVE] Archivo descargado exitosamente: %s", file.Name)

	return resp.Body, file.Name, nil
}

// IsGoogleDriveURL verifica si una URL es de Google Drive
func IsGoogleDriveURL(url string) bool {
	return regexp.MustCompile(`drive\.google\.com`).MatchString(url)
}
