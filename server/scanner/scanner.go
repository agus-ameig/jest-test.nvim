package scanner

import (
	"bufio"
	"io/fs"
	"jest/scanner/message"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

type Test struct {
  Name      string  `json:"name"`
  Line      int     `json:"line"`
}

type TestCase struct {
  Name      string  `json:"name"`
  Line      int     `json:"line"`
  Tests     []Test  `json:"tests"`
}

type TestTree struct {
  Name      string        `json:"name"`
  Path      string        `json:"path"`
  Type      string        `json:"type"`
  TestCases []TestCase    `json:"test_cases"`
  Children  []TestTree    `json:"children"`
}

func (t TestTree) IsDir() bool {
  return t.Type == "DIR"
}

type AdapterDetails struct {
  TestCasePattern   string
  TestPattern       string
}

var adapterDetails = AdapterDetails{
  TestCasePattern: "describe\\([\"']([^\"']+)[\"']",
  TestPattern: "(test|it)\\([\"']([^\"']+)[\"']",
}

var testTree TestTree

func SetUp(configuration message.Configuration) *TestTree {
  log.Println("Setting up...")
  d, err := os.Open(configuration.Dir)
  if err != nil {
    log.Println("Error opening directory", err)
    return nil
  }
  defer d.Close()

  files, err := d.ReadDir(-1)
  if err != nil {
    log.Println("Error reading directory contents:", err)
    return nil
  }
  tree, err := filter(files, configuration.Dir, configuration)
  var finalTree []TestTree
  for _, item := range tree {
    child := &item
    if item.IsDir() {
      child, err = findInDir(&item, configuration)
      if err != nil {
        log.Println("Error in children", err)
        return nil
      }
    }
    finalTree = append(finalTree, *child)
  }
  testTree = TestTree{
    Name: filepath.Base(configuration.Dir),
    Path: configuration.Dir,
    Type: "DIR",
    TestCases: []TestCase{},
    Children: finalTree,
  }
  log.Println("This is the test tree", testTree)
  return &testTree
}

func filter(files []fs.DirEntry, parentPath string, configuration message.Configuration) ([]TestTree, error) {
  var testTreeItems []TestTree

  for _, file := range files {
    if file.IsDir()  && !slices.Contains(configuration.Exclude, file.Name()){
      item := TestTree{
        Name: file.Name(),
        Path: filepath.Join(parentPath, file.Name()),
        Type: "DIR",
        TestCases: []TestCase{},
        Children: []TestTree{},
      }
      testTreeItems = append(testTreeItems, item)
    }

    m, err := regexp.MatchString(configuration.Pattern, file.Name());
    if  !file.IsDir() &&  err == nil && m {
      tests, err := findTestsInFile(filepath.Join(parentPath, file.Name()))
      if err != nil {
        return testTreeItems, err
      }

      item := TestTree {
        Name: file.Name(),
        Path: filepath.Join(parentPath, file.Name()),
        Type: "FILE",
        TestCases: tests,  //Find the Tests here
        Children: []TestTree{},
      }
      testTreeItems = append(testTreeItems, item)
    } else if err != nil {
      return testTreeItems, err
    }
  }

  return testTreeItems, nil
}

func findInDir(root *TestTree, configuration message.Configuration) (*TestTree, error) {
    d, err := os.Open(root.Path)
    if err != nil {
      return root, err
    }
    defer d.Close()

    files, err := d.ReadDir(-1)
    if err != nil {
      return root, err
    }
    children, err := filter(files, root.Path, configuration)
    if err != nil {
      return root, err
    }

    var finalChildren []TestTree
    for _, item := range children {
      child := &item
      if item.IsDir() {
        child, _  = findInDir(&item, configuration)
      }
      finalChildren = append(finalChildren, *child)
    }
    root.Children = finalChildren
    return root, nil
}



func findTestsInFile(path string) ([]TestCase, error) {
  file, err := os.Open(path)
  if err != nil {
    return []TestCase{}, err
  }
  defer file.Close()

  testCasePattern := regexp.MustCompile(adapterDetails.TestCasePattern)
  testPattern := regexp.MustCompile(adapterDetails.TestPattern)
  namePattern := regexp.MustCompile("[\"']([^\"']+)[\"']")
  result  := []TestCase{}
  var currentTestCase TestCase
  scanner := bufio.NewScanner(file)
  lineNumber := 1

  for scanner.Scan() {
    line := scanner.Text()

    if testCasePattern.MatchString(line) {
      if currentTestCase.Line != 0 {
        result = append(result, currentTestCase)
      }
      name := namePattern.FindString(line)
      currentTestCase = TestCase{
        Name: name[1: len(name) - 1],
        Line: lineNumber,
      }
    } else if testPattern.MatchString(line) {
      name := namePattern.FindString(line)
      test := Test{
        Name: name[1: len(name) - 1],
        Line: lineNumber,
      }
      currentTestCase.Tests = append(currentTestCase.Tests, test)
    }

    lineNumber ++
  }

  if currentTestCase.Line != 0 {
    result = append(result, currentTestCase)
  }

  return result, scanner.Err()
}
