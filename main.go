//A simple file manager in go

package main

import (
	"bufio"
	"fmt"
	"github.com/yorukot/ansichroma"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	//actions
	deleting bool = false

	//any info to tell the user
	instruction = ""

	p *tea.Program

	//border characters
	border_top              = "─"
	border_bottom           = "─"
	border_left             = "│"
	border_right            = "│"
	border_top_left         = "╭"
	border_top_right        = "╮"
	border_bottom_left      = "╰"
	border_bottom_right     = "╯"
	border_middle_top_left  = "┤"
	border_middle_top_right = "├"

	//colors
	red       = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	lightGray = lipgloss.NewStyle().Foreground(lipgloss.Color("#a1a6a2"))
	gray      = lipgloss.NewStyle().Foreground(lipgloss.Color("#484a48"))
	white     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	grayCol   = lipgloss.Color("#484a48")
)

// the file manager
type FM struct {
	dir          string
	files        []os.DirEntry
	quit         bool
	pos          int
	height       int
	fullWidth    int
	offset       int
	maxH         int
	isFileLocked bool
	fileContent  string
}

// Icon
type IconStyle struct {
	Icon  string
	Color string
}

//Icons for special directories
var DirIcons = map[string]IconStyle{
	".config":     {Icon: "", Color: "NONE"},
	".git":        {Icon: "", Color: "NONE"},
	"Desktop":     {Icon: "", Color: "NONE"},
	"Development": {Icon: "", Color: "NONE"},
	"Documents":   {Icon: "", Color: "NONE"},
	"Downloads":   {Icon: "", Color: "NONE"},
	"Library":     {Icon: "", Color: "NONE"},
	"Movies":      {Icon: "", Color: "NONE"},
	"Music":       {Icon: "", Color: "NONE"},
	"Pictures":    {Icon: "", Color: "NONE"},
	"Public":      {Icon: "", Color: "NONE"},
	"Videos":      {Icon: "", Color: "NONE"},
	"Folder":	   {Icon: "\uf115", Color: "NONE" },
}

// Icons for folder and file images
var Icons = map[string]IconStyle{
	"ai": {
		Icon:  "",
		Color: "#ce6f14",
	},
	"android":      {Icon: "", Color: "#a7c83f"},
	"apple":        {Icon: "", Color: "#78909c"},
	"asm":          {Icon: "󰘚", Color: "#ff7844"},
	"audio":        {Icon: "", Color: "#ee524f"},
	"binary":       {Icon: "", Color: "#ff7844"},
	"c":            {Icon: "", Color: "#0188d2"},
	"cfg":          {Icon: "", Color: "#8B8B8B"},
	"clj":          {Icon: "", Color: "#68b338"},
	"conf":         {Icon: "", Color: "#8B8B8B"},
	"cpp":          {Icon: "", Color: "#0188d2"},
	"css":          {Icon: "", Color: "#2d53e5"},
	"dart":         {Icon: "", Color: "#03589b"},
	"db":           {Icon: "", Color: "#FF8400"},
	"deb":          {Icon: "", Color: "#ab0836"},
	"doc":          {Icon: "", Color: "#295394"},
	"dockerfile":   {Icon: "󰡨", Color: "#099cec"},
	"ebook":        {Icon: "", Color: "#67b500"},
	"env":          {Icon: "", Color: "#eed645"},
	"f":            {Icon: "󱈚", Color: "#8e44ad"},
	"file":         {Icon: "\uf15b", Color: "NONE"},
	"font":         {Icon: "\uf031", Color: "#3498db"},
	"fs":           {Icon: "\ue7a7", Color: "#2ecc71"},
	"gb":           {Icon: "\ue272", Color: "#f1c40f"},
	"gform":        {Icon: "\uf298", Color: "#9b59b6"},
	"git":          {Icon: "\ue702", Color: "#e67e22"},
	"go":           {Icon: "", Color: "#6ed8e5"},
	"graphql":      {Icon: "\ue662", Color: "#e74c3c"},
	"glp":          {Icon: "󰆧", Color: "#3498db"},
	"groovy":       {Icon: "\ue775", Color: "#2ecc71"},
	"gruntfile.js": {Icon: "\ue74c", Color: "#3498db"},
	"gulpfile.js":  {Icon: "\ue610", Color: "#e67e22"},
	"gv":           {Icon: "\ue225", Color: "#9b59b6"},
	"h":            {Icon: "\uf0fd", Color: "#3498db"},
	"haml":         {Icon: "\ue664", Color: "#9b59b6"},
	"hs":           {Icon: "\ue777", Color: "#2980b9"},
	"html":         {Icon: "\uf13b", Color: "#e67e22"},
	"hx":           {Icon: "\ue666", Color: "#e74c3c"},
	"ics":          {Icon: "\uf073", Color: "#f1c40f"},
	"image":        {Icon: "\uf1c5", Color: "#e74c3c"},
	"iml":          {Icon: "\ue7b5", Color: "#3498db"},
	"ini":          {Icon: "󰅪", Color: "#f1c40f"},
	"ino":          {Icon: "\ue255", Color: "#2ecc71"},
	"iso":          {Icon: "󰋊", Color: "#f1c40f"},
	"jade":         {Icon: "\ue66c", Color: "#9b59b6"},
	"java":         {Icon: "\ue738", Color: "#e67e22"},
	"jenkinsfile":  {Icon: "\ue767", Color: "#e74c3c"},
	"jl":           {Icon: "\ue624", Color: "#2ecc71"},
	"js":           {Icon: "\ue781", Color: "#f39c12"},
	"json":         {Icon: "\ue60b", Color: "#f1c40f"},
	"jsx":          {Icon: "\ue7ba", Color: "#e67e22"},
	"key":          {Icon: "\uf43d", Color: "#f1c40f"},
	"ko":           {Icon: "\uebc6", Color: "#9b59b6"},
	"kt":           {Icon: "\ue634", Color: "#2980b9"},
	"less":         {Icon: "\ue758", Color: "#3498db"},
	"lock":         {Icon: "\uf023", Color: "#f1c40f"},
	"log":          {Icon: "\uf18d", Color: "#7f8c8d"},
	"lua":          {Icon: "\ue620", Color: "#e74c3c"},
	"maintainers":  {Icon: "\uf0c0", Color: "#7f8c8d"},
	"makefile":     {Icon: "\ue20f", Color: "#3498db"},
	"md":           {Icon: "\uf48a", Color: "#7f8c8d"},
	"mjs":          {Icon: "\ue718", Color: "#f39c12"},
	"ml":           {Icon: "󰘧", Color: "#2ecc71"},
	"mustache":     {Icon: "\ue60f", Color: "#e67e22"},
	"nc":           {Icon: "󰋁", Color: "#f1c40"},
	"nim":          {Icon: "\ue677", Color: "#3498db"},
	"nix":          {Icon: "\uf313", Color: "#f39c12"},
	"npmignore":    {Icon: "\ue71e", Color: "#e74c3c"},
	"package":      {Icon: "󰏗", Color: "#9b59b6"},
	"passwd":       {Icon: "\uf023", Color: "#f1c40f"},
	"patch":        {Icon: "\uf440", Color: "#e67e22"},
	"pdf":          {Icon: "\uf1c1", Color: "#d35400"},
	"php":          {Icon: "\ue608", Color: "#9b59b6"},
	"pl":           {Icon: "\ue7a1", Color: "#3498db"},
	"prisma":       {Icon: "\ue684", Color: "#9b59b6"},
	"ppt":          {Icon: "\uf1c4", Color: "#c0392b"},
	"psd":          {Icon: "\ue7b8", Color: "#3498db"},
	"py":           {Icon: "\ue606", Color: "#3498db"},
	"r":            {Icon: "\ue68a", Color: "#9b59b6"},
	"rb":           {Icon: "\ue21e", Color: "#9b59b6"},
	"rdb":          {Icon: "\ue76d", Color: "#9b59b6"},
	"rpm":          {Icon: "\uf17c", Color: "#d35400"},
	"rs":           {Icon: "\ue7a8", Color: "#f39c12"},
	"rss":          {Icon: "\uf09e", Color: "#c0392b"},
	"rst":          {Icon: "󰅫", Color: "#2ecc71"},
	"rubydoc":      {Icon: "\ue73b", Color: "#e67e22"},
	"sass":         {Icon: "\ue603", Color: "#e74c3c"},
	"scala":        {Icon: "\ue737", Color: "#e67e22"},
	"shell":        {Icon: "\uf489", Color: "#2ecc71"},
	"shp":          {Icon: "󰙞", Color: "#f1c40f"},
	"sol":          {Icon: "󰡪", Color: "#3498db"},
	"sqlite":       {Icon: "\ue7c4", Color: "#27ae60"},
	"styl":         {Icon: "\ue600", Color: "#e74c3c"},
	"svelte":       {Icon: "\ue697", Color: "#ff3e00"},
	"swift":        {Icon: "\ue755", Color: "#ff6f61"},
	"tex":          {Icon: "\u222b", Color: "#9b59b6"},
	"tf":           {Icon: "\ue69a", Color: "#2ecc71"},
	"toml":         {Icon: "󰅪", Color: "#f39c12"},
	"ts":           {Icon: "󰛦", Color: "#2980b9"},
	"twig":         {Icon: "\ue61c", Color: "#9b59b6"},
	"txt":          {Icon: "\uf15c", Color: "#7f8c8d"},
	"vagrantfile":  {Icon: "\ue21e", Color: "#3498db"},
	"video":        {Icon: "\uf03d", Color: "#c0392b"},
	"vim":          {Icon: "\ue62b", Color: "#019833"},
	"vue":          {Icon: "\ue6a0", Color: "#41b883"},
	"windows":      {Icon: "\uf17a", Color: "#4a90e2"},
	"xls":          {Icon: "\uf1c3", Color: "#27ae60"},
	"xml":          {Icon: "\ue796", Color: "#3498db"},
	"yml":          {Icon: "\ue601", Color: "#f39c12"},
	"zig":          {Icon: "\ue6a9", Color: "#9b59b6"},
	"zip":          {Icon: "\uf410", Color: "#e74c3c"},
}

// The metadate we find from each file
type metaData struct {
	name    string
	size    int32
	modTime string
}

// Model defines the application's state, contains the parts of the app
type model struct {
	fm   FM
	meta metaData
}

// creates a border with our border characters
func generateBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         border_top,
		Bottom:      border_bottom,
		Left:        border_left,
		Right:       border_right,
		TopLeft:     border_top_left,
		TopRight:    border_top_right,
		BottomLeft:  border_bottom_left,
		BottomRight: border_bottom_right,
	}
}

// find the icon in the list for a filename
func getFileIcon(filename string) (string, string) {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	iconStyle, exists := Icons[ext]
	if !exists {
		iconStyle = Icons["file"] // Default icon if extension not found
	}
	color := iconStyle.Color
	if color == "NONE" {
		return color, iconStyle.Icon
	}
	return color, iconStyle.Icon
}

// find the icon in the list for a filename
func getDirIcon(filename string) (string, string) {
	iconStyle, exists := DirIcons[filename]
	if !exists {
		iconStyle = DirIcons["Folder"] // Default icon if extension not found
	}
	color := iconStyle.Color
	if color == "NONE" {
		return color, iconStyle.Icon
	}
	return color, iconStyle.Icon
}

// Initialize the model
func initialModel() model {
	pwd, _ := os.Getwd()
	if len(os.Args) > 1 {
		pwd, _ = filepath.Abs(os.Args[1])
	}
	files, err := os.ReadDir(pwd)
	if err != nil {
		log.Fatal("couldn't fetch directory")
	}

	fm := FM{
		dir:   pwd,
		files: files,
		// Initialize other fields of FM as needed
	}

	return model{fm: fm}
}

// Init is called when the program starts
func (m model) Init() tea.Cmd {
	m.fm.pos = 0
	return nil
}

// Update is called when messages are received
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//If resize and also at the start
	case tea.WindowSizeMsg:
		m.fm.height = msg.Height - 3 // Adjust for header and footer
		m.fm.maxH = m.fm.height - 10
		m.fm.fullWidth = msg.Width - 40
	//if a key is pressed
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q":
			//quit
			m.fm.quit = true
			return m, tea.Quit
		case "d", "D":
			//add instruction and tell model we may be deleting
			deleting = true
			instruction += "Do you want to delete this file (y/n)? | "
		case "y", "Y":
			if deleting {
				//if deleting, remove the file, refresh the screen and tell model we aren't deleting
				os.Remove(filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name()))
				m.RefreshFM()
				deleting = false
			}
		case "n":
			//cancel any actions
			if deleting {
				deleting = false
			}
		case "up":
			//if not at the first file move up
			if m.fm.pos > 0 {
				m.fm.pos--
				if m.fm.pos < m.fm.offset {
					m.fm.offset--
				}
				//check if the file is locked
				if m.isLocked(filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name())) {
					m.fm.isFileLocked = true
				} else {
					m.fm.isFileLocked = false
				}
			}
		case "down":
			//if not at bottom of the list of files, go down
			if m.fm.pos < len(m.fm.files)-1 {
				m.fm.pos++
				if m.fm.pos >= m.fm.offset+m.fm.maxH {
					m.fm.offset++
				}
				//check if file is locked
				if m.isLocked(filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name())) {
					m.fm.isFileLocked = true
				} else {
					m.fm.isFileLocked = false
				}
			}
		case "left":
			//go to parent directory
			m.GoToParentDir()
		case "right":
			//open directory if file is a folder if not a folder, tell the user
			if m.fm.files[m.fm.pos].IsDir() {
				nestedDir := filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name())
				files, err := os.ReadDir(nestedDir)
				if err == nil {
					m.fm.dir = nestedDir
					m.fm.pos = 0
					m.fm.offset = 0
					m.fm.files = files
				} else if m.isLocked(nestedDir) {
					instruction += "You don't have access to this folder | "
				}
			} else {
				instruction += "this is a file not a folder | "
			}
		case "enter":
			exec.Command("cmd", "/c", "start", filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name())).Output()
		}
	}
	return m, nil
}

// refreshes the screen
func (m *model) RefreshFM() {
	files, err := os.ReadDir(m.fm.dir)
	if err == nil {
		m.fm.files = files
	} else {
		m.GoToParentDir()
	}
	if m.fm.pos > len(m.fm.files)-1 {
		m.fm.pos = len(m.fm.files) - 1
	} else if m.fm.pos < 0 {
		m.fm.pos = 0
	}
}

func (m *model) GoToParentDir() {
	//find the directory above the current
	parentDir := filepath.Dir(m.fm.dir)
	//read directory
	files, err := os.ReadDir(parentDir)
	if err == nil {
		//reset most things on model
		m.fm.dir = parentDir
		m.fm.files = files
		m.fm.pos = 0
		m.fm.offset = 0
	}
}

func (m model) isLocked(filepath string) bool {
	//if an err the file is probably locked, that was the only error I encountered in testing
	_, err := os.ReadDir(filepath)
	return err != nil
}

// View is called to render the UI
func (m model) View() string {
	if m.fm.quit {
		//go back to default screen
		p.RestoreTerminal()
		return ""
	}
	//read directory and do saftey checks
	files, err := os.ReadDir(m.fm.dir)
	if err == nil {
		m.fm.files = files
	} else {
		m.GoToParentDir()
	}
	if m.fm.pos > len(m.fm.files)-1 {
		m.fm.pos = len(m.fm.files) - 1
	} else if m.fm.pos < 0 {
		m.fm.pos = 0
	}
	//set widths
	var s string
	width := m.fm.height + 4
	var maxWidth int = 70
	folderName := filepath.Base(m.fm.dir)
	border := generateBorder()
	// Ensure repeat count is non-negative
	repeatCount := width - len(folderName) - 6
	if repeatCount < 0 {
		repeatCount = 0
	}
	//make borders and some styles
	border.Top = border_top + border_middle_top_left + " " + folderName + " " + border_middle_top_right + strings.Repeat(border_top, repeatCount)

	boxStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left).
		Border(border).
		BorderForeground(lipgloss.Color("103")).
		PaddingRight(1)

	fileStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("105"))
	currfileStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("115"))

	var Otext []string

	//Get the contents of the directory onto the screen
	for i := m.fm.offset; i < m.fm.offset+m.fm.maxH; i++ {
		if i < len(m.fm.files) {
			name := m.fm.files[i].Name()
			if len(name) > width-6 {
				name = name[:width-6] + "…"
			}
			style := fileStyle
			if i == m.fm.pos {
				style = currfileStyle
			}
			before := " "
			beforeStyle := style
			if !m.fm.files[i].IsDir() {
				color, Icon := getFileIcon(filepath.Join(m.fm.dir, m.fm.files[i].Name()))
				before = Icon + " "
				if color != "NONE" {
					beforeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
				}
			} else if m.fm.files[i].IsDir() {
				if m.isLocked(filepath.Join(m.fm.dir, m.fm.files[i].Name())) {
					before = "\uf023 "
					beforeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("F44336"))
				} else {
					color, Icon := getDirIcon(m.fm.files[i].Name())
					before = Icon + " "
					if color != "NONE" {
						beforeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
					}
				}
			}
			Otext = append(Otext, beforeStyle.Render(before)+style.Render(name))
		} else {
			Otext = append(Otext, "")
		}
	}
	combined := strings.Join(Otext, "\n")
	combined = boxStyle.Render(combined)

	//previewing the file
	var filePrev string
	if len(m.fm.files) > 0 {
		//styles+borders
		prevBord := generateBorder()
		if len(m.fm.files[m.fm.pos].Name()) > maxWidth-5 {
			prevBord.Top = border_top + border_middle_top_left + " " + m.fm.files[m.fm.pos].Name() + " " + border_middle_top_right
		} else {

			prevBord.Top = border_top + border_middle_top_left + " " + m.fm.files[m.fm.pos].Name() + " " + border_middle_top_right + strings.Repeat(border_top, maxWidth-len(m.fm.files[m.fm.pos].Name())-5)
		}
		prevStyle := lipgloss.NewStyle().
			Border(prevBord).
			BorderForeground(lipgloss.Color("103")).
			MaxWidth(maxWidth + 3).
			MarginLeft(2).
			PaddingRight(1)
		//use this to cut the text if it is too long
		cutter := lipgloss.NewStyle().MaxWidth(maxWidth - 5).MaxHeight(m.fm.height)

		//if what we are on isnt a directory
		if !m.fm.files[m.fm.pos].IsDir() {
			filePath := filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name())
			//check if the file is readable (not binary)
			if !isBinaryFile(filePath) {
				//open file
				file, err := os.Open(filePath)
				if err == nil {
					defer file.Close()
					//get contents
					fileConts, err := os.ReadFile(filePath)
					if err != nil {
						m.fm.fileContent = string(fileConts[:])
					}

					var lines []string
					//see if we can highlight the text
					highlightedText, err := ansichroma.HighlightFromFile(filePath, m.fm.height, "witchhazel", "")
					if err == nil {
						//cut the text to right width and height if it is too long then render
						filePrev = prevStyle.Render(cutter.Render(highlightedText))
					} else {
						//read the file as a stream, wrap it, then render it
						scanner := bufio.NewScanner(file)

						for scanner.Scan() {
							wrappedLines := wrapText(scanner.Text(), maxWidth-5)
							lines = append(lines, wrappedLines...)
							if len(lines) >= m.fm.height {
								lines = lines[:m.fm.height]
								break
							}
						}
						filePrev = prevStyle.Render(lightGray.Render(strings.Join(lines, "\n")))
					}
				}
			}
		}
	}
	//initialize variables
	var metaData metaData
	var final string
	var MBord = generateBorder()
	// Ensure repeat count is non-negative
	repeatCountMeta := width - 8
	if repeatCountMeta < 0 {
		repeatCountMeta = 0
	}
	//make border
	MBord.Top = border_top + border_middle_top_left + " Metadata " + border_middle_top_right + strings.Repeat(border_top, repeatCountMeta)
	//if not in a empty directory
	if len(m.fm.files) > 0 {
		//get name and wrap
		metaData.name = wrapTextSingleLine(m.fm.files[m.fm.pos].Name(), width-1-len(" file name: "))
		//get the last mod time and size if it is a file
		if !m.fm.files[m.fm.pos].IsDir() {
			fileStats, err := os.Stat(filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name()))
			if err == nil {
				metaData.size = int32(fileStats.Size())
				metaData.modTime = fileStats.ModTime().Format("02-01-2006 15:04:05")
			} else {
				panic(err)
			}
		}
		final = gray.Render("File name: ") + white.Render(metaData.name) + "\n"
		//render
		if !m.fm.files[m.fm.pos].IsDir() {
			metaData.modTime = "\n" + strings.ReplaceAll(metaData.modTime, " ", "\n")
			final += gray.Render("File size (bytes): ") + white.Render(fmt.Sprint(metaData.size)) + "\n" + gray.Render("Date modified: ") + white.Render(metaData.modTime) + "\n"
		}
		lines := strings.Count(final, "\n")
		//add blank space to fill in any gap
		if lines < 7 {
			final += strings.Repeat("\n", 7-lines)
		}
	} else {
		//say the folder is empty and fill in the space
		final = "This folder is empty"
		lines := strings.Count(final, "\n")
		if lines < 7 {
			final += strings.Repeat("\n", 7-lines)
		}
	}
	//instruction handling
	Finstruction := gray.Render("press ") + lightGray.Render("q") + gray.Render(" to ") + lightGray.Render("quit")
	if instruction != "" {
		Finstruction += gray.Render(" | ") + red.Render(instruction)
	}
	//border
	metaStyle := lipgloss.NewStyle().
		Width(width).
		Border(MBord).
		BorderForeground(grayCol).
		PaddingRight(1).
		Align(lipgloss.Left).
		PaddingLeft(1)
	//add it all together
	s += lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.JoinVertical(lipgloss.Center, combined, metaStyle.Render(final)), filePrev)
	s += "\n" + Finstruction
	instruction = ""
	return s
}

func isBinaryFile(filePath string) bool {
	//create a buffer
	const bufferSize = 8000
	//open file
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()
	//read the file as a stream
	reader := bufio.NewReader(file)
	//make empty buffer the size of buffer size (8000)
	buffer := make([]byte, bufferSize)
	n, err := reader.Read(buffer)
	if err != nil {
		return false
	}

	for i := 0; i < n; i++ {
		if buffer[i] > 0 && buffer[i] < 32 && buffer[i] != 9 && buffer[i] != 10 && buffer[i] != 13 {
			return true
		}
	}

	return false
}

func wrapText(text string, maxWidth int) []string {
	//if to long make a new part of the array to simulate new line
	var wrapped []string
	for len(text) > maxWidth {
		wrapped = append(wrapped, text[:maxWidth])
		text = text[maxWidth:]
	}
	wrapped = append(wrapped, text)
	return wrapped
}

func wrapTextSingleLine(text string, maxWidth int) string {
	//if too long, truncate and add: …
	if maxWidth <= 0 {
		return "…"
	}
	if len(text) > maxWidth {
		return text[:maxWidth-1] + "…"
	}
	return text
}

func main() {
	//make a new program
	p = tea.NewProgram(initialModel())
	// I don't care if it is deprecated EnterAltScreen was the only thing that worked
	p.EnterAltScreen()
	//same with p.Start() it was just easier
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
