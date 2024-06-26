/*
File: define_selector.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-10 13:33:59

Description: 定义选择器

- Init, Update, View 是不可或缺的方法
*/

package general

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/color"
)

const (
	selectKey = " "     // 默认的选择键
	enterKey  = "enter" // 默认的确认键
	quietKey  = "q"     // 默认的退出键

	selectorType = "repository name" // 选择器主题
)

// 实际按键和显示文本的映射
var keyMap = map[string]string{
	selectKey: func() string {
		key := strings.TrimSpace(selectKey)
		if key == "" {
			key = "space"
		}
		return UpperFirstChar(key)
	}(),
	enterKey: UpperFirstChar(enterKey),
	quietKey: quietKey,
}

// model 结构体，选择器的数据模型
type model struct {
	choices   []string         // 所有选项
	hlChoices []string         // 高亮选项
	cursor    int              // 光标当前所在选项的索引
	selected  map[int]struct{} // 已选中选项，key 为选项 choices 的索引。使用 map 便于判断指定选项是否已被选中
	negatives string           // 希望选择器在运行后输出的信息
	ready     bool             // 模型是否准备好
	viewport  viewport.Model   // 视图窗口
	builder   strings.Builder  // 用于构建字符串
}

// initialModel 初始化选择器数据模型
//
// 参数：
//   - choices: 可选项
//   - highlights: 高亮项
//   - negatives: 希望选择器在运行后输出的信息
//
// 返回：
//   - 初始化后的选择器数据模型
func initialModel(choices, highlights []string, negatives string) *model {
	allChoices := make([]string, 0)
	allChoices = append(allChoices, color.Sprintf("%s%s", SelectAllFlag, FgLightYellowText(SelectAllTips)))
	allChoices = append(allChoices, choices...)

	hlChoices := make([]string, 0)
	hlChoices = append(hlChoices, highlights...)

	return &model{
		choices:   allChoices,
		hlChoices: hlChoices,
		cursor:    0,
		selected:  make(map[int]struct{}),
		negatives: negatives,
	}
}

// Init 选择器数据模型的初始化方法，是 BubbleTea 框架中的一个特殊方法
//
// 返回：
//   - 一个 I/O 操作，完成后会返回一条消息，如果为 nil 则被视为无 I/O 操作
func (m *model) Init() tea.Cmd {
	return nil
}

// Update 选择器数据模型的更新方法，是 BubbleTea 框架中的一个特殊方法
//
// 参数：
//   - msg: 包含来自 I/O 操作结果的数据，触发数据模型更新，并以此触发 UI 绘制
//
// 返回：
//   - 更新后的数据模型
//   - 一个 I/O 操作，完成后会返回一条消息，如果为 nil 则被视为无操作
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg: // 监控按键事件
		// 对按下的相应按键做出对应反应
		switch msg.String() {
		case quietKey, "ctrl+c", "esc": // 按下退出键
			// 退出前取消所有选中
			m.selected = make(map[int]struct{})
			return m, tea.Quit
		case "up", "k": // 按下光标上移键
			// 光标向上移动，到达上边界后循环
			m.cursor--
			m.fixCursor()
			m.fixViewport(false)
		case "down", "j": // 按下光标下移键
			// 光标向下移动，到达下边界后循环
			m.cursor++
			m.fixCursor()
			m.fixViewport(false)
		case "pgup", "u":
			m.viewport.LineUp(1)
			m.fixViewport(true)
		case "pgdown", "d":
			m.viewport.LineDown(1)
			m.fixViewport(true)
		case selectKey: // 按下选择键
			if m.cursor == 0 { // 选中“全选”项
				// 在全选和取消全选之间切换
				if len(m.selected) == len(m.choices)-1 { // 取消全选
					m.selected = make(map[int]struct{})
				} else { // 全选
					m.selected = make(map[int]struct{})
					for i := 1; i < len(m.choices); i++ {
						m.selected[i] = struct{}{} // 选中所有非 "Select All" 选项
					}
				}
			} else { // 选中其他选项
				// 判断当前光标所在选项是否被选中
				_, ok := m.selected[m.cursor]
				if ok { // 已选中，取消选中
					delete(m.selected, m.cursor)
				} else { // 未选中，选中
					m.selected[m.cursor] = struct{}{}
				}
			}
		case enterKey: // 按下确认键
			// 回车执行
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.Model{
				Width:  msg.Width,
				Height: m.getOptimalViewportHeight(msg),
			}
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = m.getOptimalViewportHeight(msg)
			m.fixViewport(true)
		}
	}

	// 将更新后的选择器数据模型返回给 BubbleTea 进行处理
	return m, nil
}

// View 选择器数据模型的视图方法，是 BubbleTea 框架中的一个特殊方法
//
// 返回：
//   - 所有绘制内容
func (m *model) View() string {
	// 构建显示内容
	var header, body, footer string

	// Header
	header = m.headerView()
	// Body
	m.viewport.SetContent(m.content())
	body = m.viewport.View()
	// Footer
	footer = m.footerView()

	return color.Sprintf("%s\n%s\n%s", header, body, footer)
}

// headerView 构建头部视图
//
// 返回：
//   - 头部内容
func (m *model) headerView() string {
	s := strings.Builder{}
	s.WriteString(m.negatives)
	s.WriteString(color.Sprintf("%s\n", strings.Repeat(Separator1st, len(MultiSelectTips)+len(selectorType))))
	s.WriteString(color.Sprintf(QuestionText(MultiSelectTips), selectorType))
	s.WriteString(color.Sprintf(SecondaryText(KeyTips), keyMap[selectKey], keyMap[enterKey], keyMap[quietKey]))
	s.WriteString(color.Sprintf("%s", strings.Repeat(Separator1st, len(MultiSelectTips)+len(selectorType))))

	return s.String()
}

// content 构建选择器选项内容
//
// 返回：
//   - 选择器选项内容
func (m *model) content() string {
	defer m.builder.Reset()

	// 对所有选项进行迭代
	SelectedFlag = SuccessText(SelectedFlag) // 选中状态标识
	selectedCount := 1                       // 已选中选项计数
	for i, choice := range m.choices {
		// 检查当前选项是否在高亮项中
		if i != 0 { // 排除索引为 0 的 "Select All"
			hiFlag := MeanFlag // 不在高亮项中
			for _, hlChoice := range m.hlChoices {
				if choice == hlChoice {
					hiFlag = NiceFlag // 在高亮项中
					break
				}
			}
			choice = color.Sprintf("%s %s", hiFlag, choice)
		}
		// 检查光标是否指向当前选项，默认未指向
		cursorFlag := CursorOffFlag // 未指向当前选项
		if m.cursor == i {
			cursorFlag = CursorOnFlag                          // 指向当前选项
			choice = color.Sprintf("\x1b[7m%s\x1b[0m", choice) // 光标所在选项着色
		}
		// 检查当前选项是否被选中
		checked := UnselectedFlag
		if i == 0 { // 如果当前选项索引为 0 即是 "Select All"
			if len(m.selected) == len(m.choices)-1 { // 所有选项都已选中（不包括 "Select All" 这个特殊选项）
				checked = SelectedFlag // 已选中
			}
		} else { // 如果当前选项索引不为 0 即非 "Select All"
			if _, ok := m.selected[i]; ok { // 当前选项的索引在已选中选项中
				checked = SelectedFlag                                                                           // 已选中
				choice = color.Sprintf("%s %s", FgLightBlueText(choice), SecondaryText("(", selectedCount, ")")) // 已选中选项着色
				selectedCount++
			}
		}
		m.builder.WriteString(color.Sprintf("%s [%s] %s\n", cursorFlag, checked, choice))
	}

	content := strings.TrimSuffix(m.builder.String(), "\n")

	return content
}

// footerView 构建底部视图
//
// 返回：
//   - 底部内容
func (m *model) footerView() string {
	s := strings.Builder{}
	s.WriteString(color.Sprintf("%s", strings.Repeat(Separator1st, len(MultiSelectTips)+len(selectorType))))
	return s.String()
}

// fixCursor 修正光标位置，防止越界
func (m *model) fixCursor() {
	if m.cursor > len(m.choices)-1 {
		m.cursor = 0
	} else if m.cursor < 0 {
		m.cursor = len(m.choices) - 1
	}
}

// fixViewport 修正视图窗口位置，防止越界
func (m *model) fixViewport(moveCursor bool) {
	top := m.viewport.YOffset                            // 当前视图窗口顶端选项索引
	bottom := m.viewport.Height + m.viewport.YOffset - 1 // 当前视图窗口底端选项索引

	if moveCursor {
		if m.cursor < top {
			m.cursor = top
		} else if m.cursor > bottom {
			m.cursor = bottom
		}
		return
	}

	if m.cursor < top {
		m.viewport.LineUp(top - m.cursor)
	} else if m.cursor > bottom {
		m.viewport.LineDown(m.cursor - bottom)
	}
}

// getOptimalViewportHeight 优化视图窗口高度
//
// 参数：
//   - msg: 终端尺寸
//
// 返回：
//   - 最佳视图窗口高度
func (m *model) getOptimalViewportHeight(msg tea.WindowSizeMsg) int {
	totalHeight := msg.Height
	contentHeight := lipgloss.Height(m.content())   // viewport 内容高度
	headerHeight := lipgloss.Height(m.headerView()) // viewport 头部高度
	footerHeight := lipgloss.Height(m.footerView()) // viewport 底部高度
	verticalMarginHeight := headerHeight + footerHeight

	if viewHeight := totalHeight - verticalMarginHeight; viewHeight < contentHeight {
		return viewHeight
	}
	return contentHeight
}

// MultipleSelectionFilter 多选筛选器，接受一个可选项切片，返回一个已选项切片，允许全选
//
// 参数：
//   - choices: 可选项
//   - highlights: 高亮项
//   - negatives: 希望选择器在运行后输出的信息
//
// 返回：
//   - 已选项
//   - 错误信息
func MultipleSelectionFilter(choices, highlights []string, negatives string) ([]string, error) {
	program := tea.NewProgram(
		initialModel(choices, highlights, negatives),
		tea.WithAltScreen(), // 启动程序时启用备用屏幕缓冲区，即程序以全窗口模式启动
	)

	// 返回选择器数据模型
	initModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	selectedChoices := []string{}
	// 将选择器数据模型断言为自定义类型
	if myModel, ok := initModel.(*model); ok {
		if len(myModel.choices)-1 > 0 { // -1 除去 "Select All" 的干扰
			for index := range myModel.selected {
				selectedChoices = append(selectedChoices, myModel.choices[index])
			}
		}
	} else {
		return nil, fmt.Errorf("initModel is not model")
	}
	return selectedChoices, nil
}
