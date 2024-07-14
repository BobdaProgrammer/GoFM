package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	instruction = ""

	p *tea.Program

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
	red = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)


type FM struct{
	dir    string
	files  []os.DirEntry
	quit   bool
	pos    int
	height int
	offset int
	maxH   int
	isFileLocked bool
}


type IconStyle struct {
	Icon  string
	Color string
}

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

type statB struct {
	instructions string
	currFileSize int32
	currFileName string
}

// Model defines the application's state
type model struct {
	fm			 FM
}

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

// Initialize the model
func initialModel() model {
    pwd, _ := os.Getwd()
    files, err := os.ReadDir(pwd)
    if err != nil {
        log.Fatal("couldn't fetch directory")
    }

    fm := FM{
        dir:    pwd,
        files:  files,
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
	case tea.WindowSizeMsg:
		m.fm.height = msg.Height - 3 // Adjust for header and footer
		m.fm.maxH = m.fm.height - 10
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q":
			m.fm.quit = true
			return m, tea.Quit
		case "up":
			if m.fm.pos > 0 {
				m.fm.pos--
				if m.fm.pos < m.fm.offset {
					m.fm.offset--
				}
				if m.isLocked(filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name())) {
					m.fm.isFileLocked = true
				} else {
					m.fm.isFileLocked = false
				}
			}
		case "down":
			if m.fm.pos < len(m.fm.files)-1 {
				m.fm.pos++
				if m.fm.pos >= m.fm.offset+m.fm.maxH {
					m.fm.offset++
				}
				if m.isLocked(filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name())) {
					m.fm.isFileLocked = true
				} else {
					m.fm.isFileLocked = false
				}
			}
		case "left":
			// Implement logic for going back a directory
			parentDir := filepath.Dir(m.fm.dir)
			files, err := os.ReadDir(parentDir)
			if err == nil {
				m.fm.dir = parentDir
				m.fm.files = files
				m.fm.pos = 0
				m.fm.offset = 0
			}
		case "right":
			if m.fm.files[m.fm.pos].IsDir() {
				nestedDir := filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name())
				files, err := os.ReadDir(nestedDir)
				if err == nil {
					m.fm.dir = nestedDir
					m.fm.pos = 0
					m.fm.offset = 0
					m.fm.files = files
				} else if m.isLocked(nestedDir) {
					instruction += "You don't have access to this folder"
				}
			} else {
				instruction += "this is a file not a folder"
			}
		}
	}
	return m, nil
}

func (m model) isLocked(filepath string) bool {
	_, err := os.ReadDir(filepath)
	return err != nil
}

// View is called to render the UI
func (m model) View() string {
	if m.fm.quit {
		p.RestoreTerminal()
		return ""
	}

	var s string

	// Calculate width of the screen
	width := m.fm.height + 4 // Adjust as necessary
	folderName := filepath.Base(m.fm.dir)
	border := generateBorder()
	border.Top = border_top + border_middle_top_left + " " + folderName + " " + border_middle_top_right + strings.Repeat(border_top, width-len(folderName))
	// Create a style for the box
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
	var text []string

	// Display files within the current window view
	for i := m.fm.offset; i < m.fm.offset+m.fm.maxH; i++ {
		if i < len(m.fm.files) {
			name := m.fm.files[i].Name()
			if len(m.fm.files[i].Name()) > width-6 {
				name = m.fm.files[i].Name()[:width-6] + "…"
			}
			style := fileStyle
			if i == m.fm.pos {
				style = currfileStyle
			}
			before := " "
			beforeStyle := style
			if !m.fm.files[i].IsDir() {
				color, Icon := getFileIcon(filepath.Join(m.fm.dir, m.fm.files[i].Name()))
				before = Icon+" "
				if color != "NONE"{
					beforeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
				}
			} else if m.fm.files[i].IsDir() {
				if m.isLocked(filepath.Join(m.fm.dir, m.fm.files[i].Name())) {
					before = "\uf023 "
					beforeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("F44336"))
				} else {
					before = "\uf115 "
				}
			}
			text = append(text, beforeStyle.Render(before)+style.Render(name))
		} else {
			text = append(text, "")
		}
	}
	combined := strings.Join(text, "\n")
	combined = boxStyle.Render(combined)
	s += combined
	var statbar statB
	statString := ""
	if len(m.fm.files) > 0 {
		statbar.currFileName = m.fm.files[m.fm.pos].Name()
		if !m.fm.files[m.fm.pos].IsDir() {
			fileStats, err := os.Stat(filepath.Join(m.fm.dir, m.fm.files[m.fm.pos].Name()))
			if err == nil {
				statbar.currFileSize = int32(fileStats.Size())
			} else {
				panic(err)
			}
		}
		statbar.instructions = "press q to quit"
		if instruction != "" {
			statbar.instructions += " | " + red.Render(instruction)
		}
		statString = "File name: " + statbar.currFileName + " | "
		if !m.fm.files[m.fm.pos].IsDir() {
			statString += "File size (bytes): " + fmt.Sprint(statbar.currFileSize) + " | "
		}
		statString += statbar.instructions
	} else {
		statString = "This folder is empty | press q to quit"
	}
	s += "\n" + statString
	instruction = ""
	return s
}

func main() {
	p = tea.NewProgram(initialModel())
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
