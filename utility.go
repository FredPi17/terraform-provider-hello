package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"

	"github.com/armon/circbuf"
)

func readEnvironmentVariables(ev map[string]interface{}) []string {
	var variables []string
	if ev != nil {
		for k, v := range ev {
			variables = append(variables, k+"="+v.(string))
		}
	}
	return variables
}

func printStackTrace(stack []string) {
	log.Printf("-------------------------")
	log.Printf("[DEBUG] Current stack:")
	for _, v := range stack {
		log.Printf("[DEBUG] -- %s", v)
	}
	log.Printf("-------------------------")
}

func NewState(environment []string, output map[string]string) *State {
	return &State{Environment: environment, Output: output}
}

func runCommand(command string, state *State, environment []string, workingDirectory string) (*State, error) {
	const maxBufSize = 8 * 1024
	// Execute the command using a shell
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

	// Setup the command
	command = fmt.Sprintf("cd %s && %s", workingDirectory, command)
	cmd := exec.Command(shell, flag, command)
	input, _ := json.Marshal(state.Output)
	stdin := bytes.NewReader(input)
	cmd.Stdin = stdin
	environment = append(environment, os.Environ()...)
	cmd.Env = environment
	stdout, _ := circbuf.NewBuffer(maxBufSize)
	stderr, _ := circbuf.NewBuffer(maxBufSize)
	cmd.Stderr = io.Writer(stderr)
	cmd.Stdout = io.Writer(stdout)
	pr, pw, err := os.Pipe()
	cmd.ExtraFiles = []*os.File{pw}

	log.Printf("[DEBUG] shell script command old state: \"%v\"", state)

	// Output what we're about to run
	log.Printf("[DEBUG] shell script going to execute: %s %s \"%s\"", shell, flag, command)

	// Run the command to completion
	err = cmd.Run()
	pw.Close()
	log.Printf("[DEBUG] Command execution completed. Reading from output pipe: >&3")

	//read back diff output from pipe
	buffer := new(bytes.Buffer)
	for {
		tmpdata := make([]byte, maxBufSize)
		bytecount, _ := pr.Read(tmpdata)
		if bytecount == 0 {
			break
		}
		buffer.Write(tmpdata)
	}
	log.Printf("[DEBUG] shell script command stdout: \"%s\"", stdout.String())
	log.Printf("[DEBUG] shell script command stderr: \"%s\"", stderr.String())
	log.Printf("[DEBUG] shell script command output: \"%s\"", buffer.String())

	if err != nil {
		return nil, fmt.Errorf("Error running command: '%v'", err)
	}

	output, err := parseJSON(buffer.Bytes())
	if err != nil {
		log.Printf("[DEBUG] Unable to unmarshall data to map[string]string: '%v'", err)
		return nil, nil
	}
	newState := NewState(environment, output)
	log.Printf("[DEBUG] shell script command new state: \"%v\"", newState)
	return newState, nil
}

func parseJSON(b []byte) (map[string]string, error) {
	os.Stdout.Write(b)
	tb := bytes.Trim(b, "\x00")
	s := string(tb)
	var f map[string]interface{}
	err := json.Unmarshal([]byte(s), &f)
	output := make(map[string]string)
	for k, v := range f {
		output[k] = v.(string)
	}
	return output, err
}

func read(d *schema.ResourceData, meta interface{}, stack []string) error {
	os.Stdout.WriteString("Reading shell script resource")
	log.Printf("[DEBUG] Reading shell script resource...")
	printStackTrace(stack)
	l := d.Get("lifecycle_commands").([]interface{})
	c := l[0].(map[string]interface{})
	command := c["read"].(string)

	//if read is not set then do nothing. assume something either create or update is setting the state
	if len(command) == 0 {
		os.Stdout.WriteString("No command provided")
		return nil
	}

	vars := d.Get("environment").(map[string]interface{})
	environment := readEnvironmentVariables(vars)
	workingDirectory := d.Get("working_directory").(string)
	o := d.Get("output").(map[string]interface{})
	output := make(map[string]string)
	for k, v := range o {
		output[k] = v.(string)
	}

	//obtain exclusive lock
	//shellMutexKV.Lock(shellScriptMutexKey)

	state := NewState(environment, output)
	newState, err := runCommand(command, state, environment, workingDirectory)
	if err != nil {
		os.Stderr.WriteString(err)
		return err
	}

	//shellMutexKV.Unlock(shellScriptMutexKey)
	if newState == nil {
		os.Stdout.WriteString("State from read operation was nil. Marking resource for deletion.")
		log.Printf("[DEBUG] State from read operation was nil. Marking resource for deletion.")
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] output:|%v|", output)
	log.Printf("[DEBUG] new output:|%v|", newState.Output)
	isStateEqual := reflect.DeepEqual(output, newState.Output)
	isNewResource := d.IsNewResource()
	isUpdatedResource := stack[0] == "update"
	if !isStateEqual && !isNewResource && !isUpdatedResource {
		log.Printf("[DEBUG] Previous state not equal to new state. Marking resource as dirty to trigger update.")
		d.Set("dirty", true)
		return nil
	}

	d.Set("output", newState.Output)

	return nil
}
