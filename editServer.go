// Server for TextAid.  Listens on localhost:8080/edit, reads a file from request, spawns $EDITOR, writes file back.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	editor := os.Getenv("EDITOR")
	if len(editor) ==0 {
		fmt.Fprintf(os.Stderr, "Need to define $EDITOR\n")
		os.Exit(1)
	}
	http.HandleFunc("/edit", func (w http.ResponseWriter, r *http.Request) {
		if err := editFile(w,r, editor); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			fmt.Fprintf(w, "%s\n", err)
		}
	})
	log.Println("listening")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func editFile(w http.ResponseWriter, r *http.Request, editor string) error {
	f, err := ioutil.TempFile("", "edit")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	defer r.Body.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		return err
	}
	
	cmd := exec.Command(editor, f.Name())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s running %#v", err, cmd)
	}
	
	f.Seek(0,0)
	_, err = io.Copy(w, f)
	return err
}
