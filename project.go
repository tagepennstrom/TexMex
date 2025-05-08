package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Project struct {
	Name      string     `json:"name"`
	Documents []Document `json:"documents"`
}

type Document struct {
	Name string `json:"name"`
}

const projectsDir = "projects"
const projectDocumentsDir = "documents"
const filename = "document"

func getProjectFromFilesystem(name string) (Project, error) {
	projectDir := fmt.Sprintf("%s/%s", projectsDir, name)
	projectdirEntries, err := os.ReadDir(projectDir)
	documents := []Document{}
	for _, entry := range projectdirEntries {
		document := Document{
			Name: entry.Name(),
		}
		documents = append(documents, document)
	}
	return Project{
		Name:      name,
		Documents: documents,
	}, err
}

func getProjectsFromFilesystem() []Project {
	projects := []Project{}
	projectsdirEntries, _ := os.ReadDir(projectsDir)
	for _, entry := range projectsdirEntries {
		project, err := getProjectFromFilesystem(entry.Name())
		if err == nil {
			projects = append(projects, project)
		}
	}
	return projects
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	projects := getProjectsFromFilesystem()

	projectsJson, err := json.Marshal(projects)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to convert projects to json: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(projectsJson)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write projects: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getProject(w, r)
	case "POST":
		createProject(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getProject(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	project, err := getProjectFromFilesystem(name)
	if err != nil {
		log.Println(fmt.Sprintf("Project doesn't exist: %s", err))
		http.Error(w, "Project doesn't exist", http.StatusBadRequest)
		return
	}

	projectJson, err := json.Marshal(project)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to convert project to json: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(projectJson)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write project: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
}

func createProject(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	projectPath := fmt.Sprintf("%s/%s", projectsDir, name)
	err := os.Mkdir(projectPath, os.FileMode(0775))
	if err != nil {
		log.Println(fmt.Sprintf("Unable to create project: %s", err))
		http.Error(w, "Unable to create project", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func projectDocumentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getProjectDocument(w, r)
	case "PUT":
		saveProjectDocument(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getProjectDocument(w http.ResponseWriter, r *http.Request) {
	projectName := r.PathValue("projectName")
	documentName := r.PathValue("documentName")

	documentPath := fmt.Sprintf("%s/%s/%s", projectsDir, projectName, documentName)
	document, err := os.ReadFile(documentPath)
	if err != nil {
		errorMessage := fmt.Sprintf("Document doesn't exist: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	_, err = w.Write(document)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write document: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
}

func saveProjectDocument(w http.ResponseWriter, r *http.Request) {
	projectName := r.PathValue("projectName")
	documentName := r.PathValue("documentName")

	updatedDocument, err := io.ReadAll(r.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("Error reading request body: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	r.Body.Close()

	log.Println("Saving document...")
	documentPath := fmt.Sprintf("%s/%s/%s", projectsDir, projectName, documentName)
	err = os.WriteFile(documentPath, updatedDocument, os.FileMode(0600))
	if err != nil {
		errorMessage := fmt.Sprintf("Error when saving file: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	log.Println("Document saved successfully")
}

func getProjectPdf(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	project, err := getProjectFromFilesystem(name)
	if err != nil {
		log.Println(fmt.Sprintf("Project doesn't exist: %s", err))
		http.Error(w, "Project doesn't exist", http.StatusBadRequest)
		return
	}
	filenameLatex := fmt.Sprintf("%s.tex", filename)
	projectPath := fmt.Sprintf("%s/%s", projectsDir, project.Name)
	latexFilepath := fmt.Sprintf("%s/%s", projectPath, filenameLatex)
	cmd := exec.Command(
		"pdflatex",
		"-interaction=nonstopmode",
		fmt.Sprintf("-output-directory=%s", projectPath),
		latexFilepath,
	)
	err = cmd.Run()
	if err != nil {
		errorMessage := fmt.Sprintf("Error compiling LaTeX file: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	filenamePdf := fmt.Sprintf("%s.pdf", filename)
	pdfFilepath := fmt.Sprintf("%s/%s", projectPath, filenamePdf)
	pdf, err := os.ReadFile(pdfFilepath)
	if err != nil {
		errorMessage := fmt.Sprintf("Error reading pdf file: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(pdf)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write pdf: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
}
